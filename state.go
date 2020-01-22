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
	MenuTrade    MenuItem = "trade"
	MenuAccounts MenuItem = "keystore"
	MenuUtil     MenuItem = "util"

	// Trade menu items
	MenuTradeSell      MenuItem = "sell"
	MenuTradeBuy       MenuItem = "buy"
	MenuTradeOrderbook MenuItem = "orderbook"
	MenuTradeTokens    MenuItem = "tokens"
	MenuTradePairs     MenuItem = "pairs"

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
	{Text: "t/trade", Description: "Start creating Buy and Sell orders with DEX trade mode."},
	{Text: "k/keystore", Description: "Manage Ethereum accounts and private keys."},
	{Text: "u/util", Description: "Misc utils for working with wallet balances."},
	{Text: "a/about", Description: "Print information about this app."},
	{Text: "q/quit", Description: "Quit from this app."},
}

var tradingSuggestions = []prompt.Suggest{
	{Text: "s/sell", Description: "Create a market Sell order."},
	{Text: "b/buy", Description: "Create a market Buy order."},
	{Text: "o/orderbook", Description: "View orderbook of a market."},
	{Text: "t/tokens", Description: "View your account token balances."},
	{Text: "p/pairs", Description: "View available pairs for trade."},
	// {Text: "h/history", Description: "Show historical data."},
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
					Text: "Press enter and type password.",
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
			case oneOf(MenuItem(cmd), MenuTradeBuy, "b", "b/buy"):
				a.argContainer = NewArgContainer(&TradeBuyArgs{})
				a.cmd = MenuTradeBuy
				a.suggestions = nil

				// a.argContainer.AddSuggestions(1, []prompt.Suggest{
				// 	{Text: "USD/BTC"},
				// 	{Text: "KEK/LOL"},
				// })

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
		case MenuAccounts:
			switch {
			case oneOf(MenuItem(cmd), MenuAccountsUse, "u", "u/use"):
				a.argContainer = NewArgContainer(&AccountUseArgs{})
				a.cmd = MenuAccountsUse
				a.suggestions = nil

				a.argContainer.AddSuggestions(0, a.controller.SuggestAccounts())

				return
			case oneOf(MenuItem(cmd), MenuAccountsCreate, "c", "c/create"):
				a.argContainer = NewArgContainer(&AccountCreateArgs{})
				a.cmd = MenuAccountsCreate
				a.suggestions = nil

				return
			case oneOf(MenuItem(cmd), MenuAccountsImport, "i", "i/import"):
				a.argContainer = NewArgContainer(&AccountImportArgs{})
				a.cmd = MenuAccountsImport
				a.suggestions = nil

				return
			case oneOf(MenuItem(cmd), MenuAccountsImportPrivKey, "p", "p/privkey"):
				a.argContainer = NewArgContainer(&AccountImportPrivKeyArgs{})
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
			logrus.Warningf("unknown command: %s", cmd)
			return
		}
	}

	a.executeCmd(cmdArgs)
	a.changeRoot(a.root)
}

func (a *AppState) executeCmd(args interface{}) {
	switch a.cmd {
	case MenuTradeSell:
		a.controller.ActionTradeSell(args)
	case MenuTradeBuy:
		a.controller.ActionTradeBuy(args)
	case MenuTradeOrderbook:
		a.controller.ActionTradeOrderbook(args)
	case MenuTradeTokens:
		a.controller.ActionTradeTokens()
	case MenuTradePairs:
		a.controller.ActionTradePairs()
	case MenuUtilUnlock:
		a.controller.ActionUtilUnlock(args)
	case MenuUtilLock:
		a.controller.ActionUtilLock(args)
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
	case oneOf(newRoot, MenuTrade, "t", "t/trade"):
		a.root = MenuTrade
		a.suggestions = tradingSuggestions

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
	obj         interface{}
	fields      []string
	suggestions map[int][]prompt.Suggest
	offset      int
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

func (a *ArgContainer) CurrentFieldSuggestions() []prompt.Suggest {
	return a.suggestions[a.offset]
}

func (a *ArgContainer) AddSuggestions(index int, suggestions []prompt.Suggest) {
	if a.suggestions == nil {
		a.suggestions = make(map[int][]prompt.Suggest)
	}

	a.suggestions[index] = suggestions
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
