package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xlab/structwalk"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/InjectiveLabs/dexterm/ethereum/ethcore"
)

func NewAppState(controller *AppController) *AppState {
	return &AppState{
		root:        MenuMain,
		suggestions: mainSuggestions,
		controller:  controller,
	}
}

type AppState struct {
	root         MenuItem
	cmd          MenuItem
	suggestions  []prompt.Suggest
	argContainer *ArgContainer
	controller   *AppController
}

type MenuItem string

const (
	// Main menu handle
	MenuMain MenuItem = "main"

	// Main menu items
	MenuTrade       MenuItem = "trade-spot"
	MenuDerivatives MenuItem = "trade-derivatives"
	MenuAccounts    MenuItem = "keystore"
	MenuUtil        MenuItem = "util"

	// Trade menu items
	MenuTradeLimitBuy       MenuItem = "limitbuy"
	MenuTradeLimitSell      MenuItem = "limitsell"
	MenuTradeFillOrder      MenuItem = "fill"
	MenuTradeCancelOrder    MenuItem = "cancel"
	MenuTradeMarketBuy      MenuItem = "marketbuy"
	MenuTradeMarketSell     MenuItem = "marketsell"
	MenuTradeGenerateLimits MenuItem = "generatelimits"
	MenuTradeOrderbook      MenuItem = "orderbook"
	MenuTradeTokens         MenuItem = "tokens"
	MenuTradePairs          MenuItem = "pairs"

	// Derivatives menu items
	MenuDerivativesLimitLong  MenuItem = "limitlong"
	MenuDerivativesLimitShort MenuItem = "limitsshort"
	MenuDerivativesOrderbook MenuItem = "derivatives-orderbook"

	// Util menu items
	MenuUtilUnlock MenuItem = "unlock"
	MenuUtilLock   MenuItem = "lock"
	MenuUtilWrap   MenuItem = "wrap"
	MenuUtilUnwrap MenuItem = "unwrap"

	// Accounts menu items
	MenuAccountsUse           MenuItem = "use"
	MenuAccountsCreate        MenuItem = "create"
	MenuAccountsImport        MenuItem = "import"
	MenuAccountsImportPrivKey MenuItem = "privkey"
	MenuAccountsList          MenuItem = "list"

	// Actions in main menu
	MenuAbout MenuItem = "about"
	MenuQuit  MenuItem = "quit"
)

var mainSuggestions = []prompt.Suggest{
	{Text: "t/trade-spot", Description: "Start creating Buy and Sell orders with DEX trade mode."},
	{Text: "d/trade-derivatives", Description: "Start creating Derivatives orders with DEX trade mode."},
	{Text: "k/keystore", Description: "Manage Ethereum accounts and private keys."},
	{Text: "u/util", Description: "Misc utils for working with wallet balances."},
	{Text: "a/about", Description: "Print information about this app."},
	{Text: "q/quit", Description: "Quit from this app."},
}

var tradingSuggestions = []prompt.Suggest{
	{Text: "b/limitbuy", Description: "Create a Limit Buy order."},
	{Text: "s/limitsell", Description: "Create a Limit Sell order."},
	{Text: "f/fill", Description: "Fill an order (Take Order)."},
	{Text: "c/cancel", Description: "Cancel an order."},

	{Text: "mb/marketbuy", Description: "Create a Market Buy order."},
	{Text: "ms/marketsell", Description: "Create a Market Sell order."},

	{Text: "g/generatelimits", Description: "Generate many limit buy and sell orders to populate the orderbook"},

	{Text: "o/orderbook", Description: "View orderbook of a market."},
	{Text: "t/tokens", Description: "View your account token balances."},
	{Text: "p/pairs", Description: "View available pairs for trade."},
	// {Text: "h/history", Description: "Show historical data."},
	{Text: "q/quit", Description: "Quit from the trading menu."},
}

var derivativesSuggestions = []prompt.Suggest{
	{Text: "l/limitlong", Description: "Create a Limit Long order."},
	{Text: "s/limitshort", Description: "Create a Limit Short order."},

	{Text: "o/orderbook", Description: "View orderbook of a market."},
	{Text: "q/quit", Description: "Quit from the trading menu."},
}

var accountsSuggestions = []prompt.Suggest{
	{Text: "u/use", Description: "Select account to use as default."},
	{Text: "c/create", Description: "Create a new account and generate a private key."},
	{Text: "i/import", Description: "Import an external keyfile into keystore."},
	{Text: "p/privkey", Description: "Import a private key into keystore."},
	{Text: "l/list", Description: "List all accounts in keystore."},
	{Text: "q/quit", Description: "Quit from the accounts menu."},
}

var utilSuggestions = []prompt.Suggest{
	{Text: "t/tokens", Description: "View your account token balances."},
	{Text: "u/unlock", Description: "Unlock a token and allow trading on the platform."},
	{Text: "l/lock", Description: "Lock a token from trade. Soft cancels all sell orders too."},
	{Text: "w/wrap", Description: "Wrap ETH into WETH ERC20 tokens."},
	{Text: "uw/unwrap", Description: "Unwrap WETH ERC20 tokens and receive ETH."},
	{Text: "q/quit", Description: "Quit from the util menu."},
}

func (a *AppState) LivePrefix() func() (prefix string, useLivePrefix bool) {
	return func() (prefix string, useLivePrefix bool) {
		if a.argContainer != nil {
			idx, name := a.argContainer.CurrentField()
			value := a.argContainer.CurrentFieldValue()

			prefix := fmt.Sprintf("%d) %s (%T) => ", idx+1, name, value)

			return prefix, true
		}

		switch a.root {
		case MenuMain, MenuAbout:
			return "", false
		default:
			return string(a.root) + " ∆ ", true
		}
	}
}

func (a *AppState) Completer() prompt.Completer {
	return func(d prompt.Document) []prompt.Suggest {
		switch {
		case a.argContainer != nil:
			if a.isCurrentFieldPassword() {
				return []prompt.Suggest{{
					Text:        "Passphrase",
					Description: "Sign using a private key, need to provide a passphrase to unlock it.",
				}}
			}

			return prompt.FilterFuzzy(a.argContainer.CurrentFieldSuggestions(), d.TextBeforeCursor(), true)
		case a.root == MenuMain,
			a.root == MenuUtil,
			a.root == MenuTrade,
			a.root == MenuAccounts:
			return prompt.FilterHasPrefix(a.suggestions, d.TextBeforeCursor(), true)
		default:
			return prompt.FilterFuzzy(a.suggestions, d.TextBeforeCursor(), true)
		}
	}
}

func (a *AppState) Executor() prompt.Executor {
	return func(cmd string) {
		switch a.root {
		case MenuMain:
			if isEmpty(cmd) {
				return
			}

			parts := strings.Fields(cmd)
			a.changeRoot(MenuItem(parts[0]))
		default:
			a.executeInRoot(cmd)
		}
	}
}

func (a *AppState) DiscardCmd() {
	a.cmd = ""
	a.argContainer = nil
	a.changeRoot(a.root)
}

func (a *AppState) Controller() *AppController {
	return a.controller
}

func (a *AppState) isCurrentFieldPassword() bool {
	_, fieldName := a.argContainer.CurrentField()
	return strings.Contains(strings.ToLower(fieldName), "password")
}

func (a *AppState) executeInRoot(cmd string) {
	var cmdArgs interface{}

	if a.argContainer != nil {
		fieldValue := a.argContainer.CurrentFieldValue()

		if a.isCurrentFieldPassword() {
			line, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				logrus.WithError(err).Warningln("failed to read input")
			} else {
				cmd = string(line)
			}
		}

		newFieldValue, err := parseFieldValue(cmd, fieldValue)
		if err != nil {
			logrus.WithError(err).Warningf("failed to parse argument for type %T", fieldValue)
			return
		}

		stop := a.argContainer.UpdateCurrentField(newFieldValue)
		if !stop {
			return
		}

		cmdArgs = a.argContainer.Object()
		a.argContainer = nil
	} else {
		switch {
		case isEmpty(cmd):
			return
		case oneOf(MenuItem(cmd), MenuQuit, "q", "q/quit", "exit"):
			a.changeRoot(MenuMain)
			return
		}
	}

	if cmdArgs == nil {
		switch a.root {
		case MenuTrade:
			switch {
			case oneOf(MenuItem(cmd), MenuTradeLimitBuy, "b", "b/limitbuy"):
				a.argContainer = NewArgContainer(&TradeLimitBuyOrderArgs{})
				a.cmd = MenuTradeLimitBuy
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())
				a.argContainer.AddSuggestions(1, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Amount must be entered as float. Minimum value is 0.0000001",
				}})
				a.argContainer.AddSuggestions(2, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Price must be entered as float. Minimum value is 0.0000001",
				}})

				return
			case oneOf(MenuItem(cmd), MenuTradeLimitSell, "s", "s/limitsell"):
				a.argContainer = NewArgContainer(&TradeLimitSellOrderArgs{})
				a.cmd = MenuTradeLimitSell
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())
				a.argContainer.AddSuggestions(1, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Amount must be entered as float. Minimum value is 0.0000001",
				}})
				a.argContainer.AddSuggestions(2, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Price must be entered as float. Minimum value is 0.0000001",
				}})

				return
			case oneOf(MenuItem(cmd), MenuTradeMarketBuy, "mb", "mb/marketbuy"):
				a.argContainer = NewArgContainer(&TradeMarketBuyOrderArgs{})
				a.cmd = MenuTradeMarketBuy
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())
				a.argContainer.AddSuggestions(1, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Amount must be entered as float. Minimum value is 0.0000001",
				}})

				return
			case oneOf(MenuItem(cmd), MenuTradeMarketSell, "ms", "ms/marketsell"):
				a.argContainer = NewArgContainer(&TradeMarketSellOrderArgs{})
				a.cmd = MenuTradeMarketSell
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())
				a.argContainer.AddSuggestions(1, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Amount must be entered as float. Minimum value is 0.0000001",
				}})

				return
			case oneOf(MenuItem(cmd), MenuTradeGenerateLimits, "g", "g/generatelimits"):
				a.argContainer = NewArgContainer(&TradeGenerateLimitOrdersArgs{})
				a.cmd = MenuTradeGenerateLimits
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())
				a.argContainer.AddSuggestions(1, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Amount must be entered as float. Minimum value is 0.0000001",
				}})

				return
			case oneOf(MenuItem(cmd), MenuTradeFillOrder, "f", "f/fill"):
				a.argContainer = NewArgContainer(&TradeFillOrderArgs{})
				a.cmd = MenuTradeFillOrder
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())
				a.argContainer.AddSuggestionsLazy(1, []int{0}, func(args ...interface{}) []prompt.Suggest {
					return a.controller.SuggestOrderToFill(args[0].(string))
				})
				a.argContainer.AddSuggestions(2, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Amount must be entered as float. Minimum value is 0.0000001",
				}})

				return
			case oneOf(MenuItem(cmd), MenuTradeCancelOrder, "c", "c/cancel"):
				a.argContainer = NewArgContainer(&TradeCancelOrderArgs{})
				a.cmd = MenuTradeCancelOrder
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())
				a.argContainer.AddSuggestionsLazy(1, []int{0}, func(args ...interface{}) []prompt.Suggest {
					return a.controller.SuggestOrderToCancel(args[0].(string))
				})

				return
			case oneOf(MenuItem(cmd), MenuTradeOrderbook, "o", "o/orderbook"):
				a.argContainer = NewArgContainer(&TradeOrderbookArgs{})
				a.cmd = MenuTradeOrderbook
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())

				return
			case oneOf(MenuItem(cmd), MenuTradePairs, "p", "p/pairs"):
				a.cmd = MenuTradePairs
				a.suggestions = nil
			case oneOf(MenuItem(cmd), MenuTradeTokens, "t", "t/tokens"):
				a.cmd = MenuTradeTokens
				a.suggestions = nil
			default:
				logrus.Warningf("unknown command: %s", cmd)
				return
			}
		case MenuDerivatives:
			switch {
			case oneOf(MenuItem(cmd), MenuDerivativesLimitLong, "l", "l/limitlong"):
				a.argContainer = NewArgContainer(&TradeDerivativeLimitOrderArgs{})
				a.cmd = MenuDerivativesLimitLong
				a.suggestions = nil
				a.argContainer.AddSuggestions(0, a.controller.SuggestDerivativesMarkets())

				a.argContainer.AddSuggestions(1, []prompt.Suggest{{
					Text:        "10",
					Description: "Quantity must be entered as positive integer. Minimum value is 1",
				}})
				a.argContainer.AddSuggestions(2, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Price must be entered as float. Minimum value is 0.0000001",
				}})

				return
			case oneOf(MenuItem(cmd), MenuDerivativesLimitShort, "s", "s/limitshort"):
				a.argContainer = NewArgContainer(&TradeDerivativeLimitOrderArgs{})
				a.cmd = MenuDerivativesLimitShort
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())
				a.argContainer.AddSuggestions(1, []prompt.Suggest{{
					Text:        "10",
					Description: "Quantity must be entered as positive integer. Minimum value is 1",
				}})
				a.argContainer.AddSuggestions(2, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Price must be entered as float. Minimum value is 0.0000001",
				}})

				return
			case oneOf(MenuItem(cmd), MenuDerivativesOrderbook, "o", "o/orderbook"):
				a.argContainer = NewArgContainer(&DerivativeOrderbookArgs{})
				a.cmd = MenuDerivativesOrderbook
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestMarkets())

				return
			default:
				logrus.Warningf("unknown command: %s", cmd)
				return
			}
		case MenuAccounts:
			switch {
			case oneOf(MenuItem(cmd), MenuAccountsUse, "u", "u/use"):
				a.argContainer = NewArgContainer(&ethcore.AccountUseArgs{})
				a.cmd = MenuAccountsUse
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestAccounts())

				return
			case oneOf(MenuItem(cmd), MenuAccountsCreate, "c", "c/create"):
				a.argContainer = NewArgContainer(&ethcore.AccountCreateArgs{})
				a.cmd = MenuAccountsCreate
				a.suggestions = nil

				return
			case oneOf(MenuItem(cmd), MenuAccountsImport, "i", "i/import"):
				a.argContainer = NewArgContainer(&ethcore.AccountImportArgs{})
				a.cmd = MenuAccountsImport
				a.suggestions = nil

				return
			case oneOf(MenuItem(cmd), MenuAccountsImportPrivKey, "p", "p/privkey"):
				a.argContainer = NewArgContainer(&ethcore.AccountImportPrivKeyArgs{})
				a.cmd = MenuAccountsImportPrivKey
				a.suggestions = nil

				return
			case oneOf(MenuItem(cmd), MenuAccountsList, "l", "l/list"):
				a.cmd = MenuAccountsList
				a.suggestions = nil
			default:
				logrus.Warningf("unknown command: %s", cmd)
				return
			}
		case MenuUtil:
			switch {
			case oneOf(MenuItem(cmd), MenuTradeTokens, "t", "t/tokens"):
				a.cmd = MenuTradeTokens
				a.suggestions = nil
			case oneOf(MenuItem(cmd), MenuUtilLock, "l", "l/lock"):
				a.argContainer = NewArgContainer(&UtilTokenLockArgs{})
				a.cmd = MenuUtilLock
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestTokens())

				return
			case oneOf(MenuItem(cmd), MenuUtilUnlock, "u", "u/unlock"):
				a.argContainer = NewArgContainer(&UtilTokenUnlockArgs{})
				a.cmd = MenuUtilUnlock
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestTokens())

				return
			case oneOf(MenuItem(cmd), MenuUtilWrap, "w", "w/wrap"):
				a.argContainer = NewArgContainer(&UtilWrapArgs{})
				a.cmd = MenuUtilWrap
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Amount must be entered as float. Minimum value is 0.0000001 ETH",
				}})

				return
			case oneOf(MenuItem(cmd), MenuUtilUnwrap, "uw", "uw/unwrap"):
				a.argContainer = NewArgContainer(&UtilUnwrapArgs{})
				a.cmd = MenuUtilUnwrap
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, []prompt.Suggest{{
					Text:        "1.00",
					Description: "Amount must be entered as float. Minimum value is 0.0000001 WETH",
				}})

				return
			default:
				logrus.Warningf("unknown command: %s", cmd)
				return
			}
		}
	}

	a.executeCmd(cmdArgs)
	a.changeRoot(a.root)
}

func (a *AppState) executeCmd(args interface{}) {
	switch a.cmd {
	case MenuDerivativesLimitLong:
		a.controller.ActionDerivativesLimitLong(args)
	case MenuDerivativesLimitShort:
		a.controller.ActionDerivativesLimitShort(args)
	case MenuDerivativesOrderbook:
		a.controller.ActionDerivativesOrderbook(args)
	case MenuTradeLimitBuy:
		a.controller.ActionTradeLimitBuy(args)
	case MenuTradeLimitSell:
		a.controller.ActionTradeLimitSell(args)
	case MenuTradeMarketBuy:
		a.controller.ActionTradeMarketBuy(args)
	case MenuTradeMarketSell:
		a.controller.ActionTradeMarketSell(args)
	case MenuTradeGenerateLimits:
		a.controller.ActionTradeGenerateLimitOrders(args)
	case MenuTradeFillOrder:
		a.controller.ActionTradeFillOrder(args)
	case MenuTradeCancelOrder:
		a.controller.ActionTradeCancelOrder(args)
	case MenuTradeOrderbook:
		a.controller.ActionTradeOrderbook(args)
	case MenuTradeTokens:
		a.controller.ActionTradeTokens()
	case MenuTradePairs:
		a.controller.ActionTradePairs()
	case MenuUtilLock:
		a.controller.ActionUtilLock(args)
	case MenuUtilUnlock:
		a.controller.ActionUtilUnlock(args)
	case MenuUtilWrap:
		a.controller.ActionUtilWrap(args)
	case MenuUtilUnwrap:
		a.controller.ActionUtilUnwrap(args)
	case MenuAccountsUse:
		a.controller.ActionAccountsUse(args)
	case MenuAccountsCreate:
		a.controller.ActionAccountsCreate(args)
	case MenuAccountsImport:
		a.controller.ActionAccountsImport(args)
	case MenuAccountsImportPrivKey:
		a.controller.ActionAccountsImportPrivKey(args)
	case MenuAccountsList:
		a.controller.ActionAccountsList()
	}
}

func (a *AppState) changeRoot(newRoot MenuItem) {
	switch {
	case isEmpty(string(newRoot)):
		return
	case oneOf(newRoot, MenuMain):
		a.root = MenuMain
		a.suggestions = mainSuggestions
	case oneOf(newRoot, MenuAbout, "a", "a/about"):
		a.controller.ActionAbout()
	case oneOf(newRoot, MenuQuit, "q", "q/quit", "exit"):
		a.controller.ActionQuit()
	case oneOf(newRoot, MenuTrade, "t", "t/trade-spot"):
		a.root = MenuTrade
		a.suggestions = tradingSuggestions

		_, ok := a.controller.getConfigValue("accounts.default")
		if !ok {
			logrus.Warningln("Default account is not set, go to keystore menu first.")
		}
	case oneOf(newRoot, MenuDerivatives, "d", "t/trade-derivatives"):
		a.root = MenuDerivatives
		a.suggestions = derivativesSuggestions

		_, ok := a.controller.getConfigValue("accounts.default")
		if !ok {
			logrus.Warningln("Default account is not set, go to keystore menu first.")
		}
	case oneOf(newRoot, MenuAccounts, "k", "k/keystore"):
		a.root = MenuAccounts
		a.suggestions = accountsSuggestions
	case oneOf(newRoot, MenuUtil, "u", "u/util"):
		a.root = MenuUtil
		a.suggestions = utilSuggestions
	default:
		logrus.Warningf("unknown command: %+v", newRoot)
	}
}

type ArgContainer struct {
	obj             interface{}
	fields          []string
	lazySuggestions map[int]lazySuggestion
	suggestions     map[int][]prompt.Suggest
	offset          int
}

func NewArgContainer(obj interface{}) *ArgContainer {
	fields := structwalk.FieldListNoSort(obj)
	if len(fields) == 0 {
		return nil
	}

	return &ArgContainer{
		obj:    obj,
		fields: fields,
	}
}

func (a *ArgContainer) Object() interface{} {
	return a.obj
}

func (a *ArgContainer) CurrentField() (index int, name string) {
	return a.offset, a.fields[a.offset]
}

func (a *ArgContainer) CurrentFieldValue() interface{} {
	v, found := structwalk.FieldValue(a.fields[a.offset], a.obj)
	if !found {
		return nil
	}

	return v
}

type lazySuggestion struct {
	Fields []int
	Fn     LazySuggestFn
}

func (a *ArgContainer) CurrentFieldSuggestions() []prompt.Suggest {
	if lazySuggstion, ok := a.lazySuggestions[a.offset]; ok {
		args := make([]interface{}, 0, len(lazySuggstion.Fields))
		for _, fieldOffset := range lazySuggstion.Fields {
			v, _ := structwalk.FieldValue(a.fields[fieldOffset], a.obj)
			args = append(args, v)
		}

		return lazySuggstion.Fn(args...)
	}

	return a.suggestions[a.offset]
}

func (a *ArgContainer) AddSuggestions(index int, suggestions []prompt.Suggest) {
	if a.suggestions == nil {
		a.suggestions = make(map[int][]prompt.Suggest)
	}

	a.suggestions[index] = suggestions
}

type LazySuggestFn func(args ...interface{}) []prompt.Suggest

func (a *ArgContainer) AddSuggestionsLazy(index int, fields []int, fn LazySuggestFn) {
	if a.lazySuggestions == nil {
		a.lazySuggestions = make(map[int]lazySuggestion)
	}

	a.lazySuggestions[index] = lazySuggestion{
		Fields: fields,
		Fn:     fn,
	}
}

func (a *ArgContainer) UpdateCurrentField(v interface{}) (stop bool) {
	structwalk.SetFieldValue(a.fields[a.offset], v, a.obj)

	a.offset++

	return a.offset >= len(a.fields)
}

func isEmpty(cmd string) bool {
	return len(strings.TrimSpace(cmd)) == 0
}

func oneOf(item MenuItem, mainOption MenuItem, options ...string) bool {
	if item == mainOption {
		return true
	}

	for _, opt := range options {
		if opt == string(item) {
			return true
		}
	}

	return false
}

func parseFieldValue(input string, v interface{}) (interface{}, error) {
	fieldType := reflect.TypeOf(v)
	if fieldType.Kind() == reflect.String {
		return input, nil
	} else if len(input) == 0 {
		err := errors.New("provided empty string")
		return nil, err
	}

	container := reflect.New(fieldType)

	if err := json.Unmarshal([]byte(input), container.Interface()); err != nil {
		return nil, err
	}

	fieldValue := container.Elem().Interface()

	return fieldValue, nil
}
