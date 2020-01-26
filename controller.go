package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	spin "github.com/tj/go-spin"
	"github.com/xlab/closer"
	"github.com/xlab/termtables"

	"github.com/InjectiveLabs/dexterm/ethfw/keystore"
	"github.com/InjectiveLabs/dexterm/ethfw/manager"
	relayer "github.com/InjectiveLabs/injective-core/api/gen/relayer"
)

func init() {
	closer.Bind(func() {
		logrus.Println("Bye!")
	})
}

type AppController struct {
	cfg           *toml.Tree
	configPath    string
	relayerClient *RelayerClient

	ethGasPrice *big.Int
	ethClient   *EthClient

	keystorePath string
	keystore     keystore.EthKeyStore
}

func NewAppController(configPath string) (*AppController, error) {
	cfg, err := loadOrCreateConfig(configPath)
	if err != nil {
		return nil, err
	}

	ctl := &AppController{
		cfg:        cfg,
		configPath: configPath,
	}

	clientCfg := &RelayerClientConfig{
		Endpoint: ctl.mustConfigValue("relayer.endpoint"),
	}

	if relayerClient, err := NewRelayerClient(clientCfg); err != nil {
		logrus.WithError(err).Warningln("running in offline mode")
	} else {
		ctl.relayerClient = relayerClient
	}

	keystorePath := ctl.mustConfigValue("accounts.keystore")
	keystorePath, _ = homedir.Expand(keystorePath)

	if err := os.MkdirAll(keystorePath, 0700); err != nil {
		return nil, err
	}

	ctl.keystorePath = keystorePath

	kb, err := keystore.New(keystorePath)
	if err != nil {
		return nil, err
	} else {
		ctl.keystore = kb
	}

	if ctl.takeFirstAccountAsDefault() {
		saveConfig(ctl.configPath, ctl.cfg)
	}

	if ctl.selectDefaultNetwork() {
		saveConfig(ctl.configPath, ctl.cfg)
	}

	if err := ctl.initEthClient(); err != nil {
		logrus.WithError(err).Warningln("failed to init Ethereum client")
	}

	return ctl, nil
}

func (ctl *AppController) ActionAbout() {
	fmt.Printf(
		"  ___  ____ _  _ ___ ____ ____ _  _\n" +
			"  |  \\ |___  \\/   |  |___ |__/ |\\/|\n" +
			"  |__/ |___ _/\\_  |  |___ |  \\ |  |\n" +
			"  Copyright (c) 2019-2020 Injective Protocol\n" +
			"  https://github.com/InjectiveLabs/dexterm\n",
	)
}

func (ctl *AppController) ActionQuit() {
	closer.Close()
}

type TradeMakeOrderArgs struct {
	MakerToken   string
	TakerToken   string
	MakerAmount  string
	Price        string
	SignPassword string
}

func (ctl *AppController) ActionTradeMakeOrder(args interface{}) {
	makeOrderArgs := args.(*TradeMakeOrderArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tradePair, err := ctl.relayerClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var makerAssetData []byte
	var takerAssetData []byte
	var isBid bool

	for _, pair := range tradePair {
		parts := strings.Split(pair.Name, "/")
		if len(parts) != 2 {
			continue
		}

		if parts[0] == makeOrderArgs.MakerToken && parts[1] == makeOrderArgs.TakerToken {
			makerAssetData = common.FromHex(pair.MakerAssetData)
			takerAssetData = common.FromHex(pair.TakerAssetData)
			break
		} else if parts[0] == makeOrderArgs.TakerToken && parts[1] == makeOrderArgs.MakerToken {
			makerAssetData = common.FromHex(pair.TakerAssetData)
			takerAssetData = common.FromHex(pair.MakerAssetData)
			isBid = true
			break
		}
	}

	if len(makerAssetData) == 0 || len(takerAssetData) == 0 {
		logrus.WithFields(logrus.Fields{
			"makerToken": makeOrderArgs.MakerToken,
			"takerToken": makeOrderArgs.TakerToken,
		}).Errorln("found no trading pair for tokens")

		return
	}

	var makerAmount *big.Int
	var price decimal.Decimal
	var takerAmount *big.Int

	makerAmountDec, err := decimal.NewFromString(makeOrderArgs.MakerAmount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse maker amount")
		return
	} else if makerAmountDec.LessThan(decimal.RequireFromString("0.0000001")) {
		logrus.Errorln("maker amount is too small, must be at least 0.0000001")
		return
	} else {
		makerAmount = dec2big(makerAmountDec)
	}

	if price, err = decimal.NewFromString(makeOrderArgs.Price); err != nil {
		logrus.WithError(err).Errorln("failed to parse price")
		return
	} else if isBid {
		takerAmount = dec2big(makerAmountDec.Div(price))
	} else {
		takerAmount = dec2big(makerAmountDec.Mul(price))
	}

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: makeOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	exchangeAddress := ctl.ethClient.contractAddresses[EthContractExchange]
	signedOrder, err := ctl.ethClient.CreateAndSignOrder(
		callArgs,
		makerAssetData,
		takerAssetData,
		makerAmount,
		takerAmount,
		exchangeAddress,
	)
	if err != nil {
		logrus.WithError(err).Errorln("unable to sign order")
		return
	}

	orderHash, err := ctl.relayerClient.PostMakeOrder(ctx, signedOrder)
	if err != nil {
		logrus.WithError(err).Errorln("unable to post make order")
		return
	}

	fmt.Println(orderHash)
}

func dec2big(d decimal.Decimal) *big.Int {
	v, _ := big.NewInt(0).SetString(d.Truncate(9).Shift(18).String(), 10)

	return v
}

type TradeFillOrderArgs struct {
	TakerToken   string
	MakerToken   string
	FillAmount   string
	OrderHash    string
	SignPassword string
}

func (ctl *AppController) ActionTradeFillOrder(args interface{}) {
	// meh
}

type TradeBuyArgs struct {
	Market       string
	Amount       string
	SignPassword string
}

func (ctl *AppController) ActionTradeBuy(args interface{}) {
	fmt.Println("Sorry, automatic order matching is not ready yet.")
}

type TradeSellArgs struct {
	Market       string
	Amount       string
	SignPassword string
}

func (ctl *AppController) ActionTradeSell(args interface{}) {
	fmt.Println("Sorry, automatic order matching is not ready yet.")
}

type TradeOrderbookArgs struct {
	Market string
}

func (ctl *AppController) ActionTradeOrderbook(args interface{}) {
	orderbookArgs := args.(*TradeOrderbookArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	bids, asks, err := ctl.relayerClient.Orderbook(ctx, orderbookArgs.Market)
	if err != nil {
		logrus.WithField("tradePair", orderbookArgs.Market).
			WithError(err).Errorln("unable to get orderbook for this market")
		return
	}

	pair := strings.Split(orderbookArgs.Market, "/")
	baseAsset := pair[0]
	quoteAsset := pair[1]

	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddTitle("ORDERBOOK")
	table.AddHeaders(
		fmt.Sprintf("Price (%s)", quoteAsset),
		fmt.Sprintf("Amount (%s)", baseAsset),
		"Notes",
	)

	if len(asks.Records) == 0 {
		table.AddRow(color.RedString("No asks."), "", "")
	} else {
		for _, ask := range asks.Records {
			var notes string
			if isMakerOf(ask.Order, defaultAccount) {
				notes = "⭑ owner"
			}

			price, vol := calcOrderPrice(ask.Order, false)
			table.AddRow(
				color.RedString("%s", price.StringFixed(9)),
				color.RedString("%s", vol.Shift(-18).StringFixed(9)),
				notes,
			)
		}
	}

	table.AddSeparator()

	if len(bids.Records) == 0 {
		table.AddRow(color.GreenString("No bids."), "", "")
	} else {
		for _, bid := range bids.Records {
			var notes string
			if isMakerOf(bid.Order, defaultAccount) {
				notes = "⭑ owner"
			}

			price, vol := calcOrderPrice(bid.Order, true)
			table.AddRow(
				color.GreenString("%s", price.StringFixed(9)),
				color.GreenString("%s", vol.Shift(-18).StringFixed(9)),
				notes,
			)
		}
	}

	fmt.Println(table.Render())
}

func isMakerOf(order *relayer.Order, address common.Address) bool {
	return bytes.Compare(
		common.HexToAddress(order.MakerAddress).Bytes(),
		address.Bytes(),
	) == 0
}

func (ctl *AppController) getTokenNamesAndAssets(ctx context.Context) (tokenNames []string, assets []common.Address, err error) {
	if ctl.relayerClient == nil {
		err := errors.New("client in offline mode")
		return nil, nil, err
	}

	pairs, err := ctl.relayerClient.TradePairs(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to list trade pairs")
		return tokenNames, assets, err
	}

	tokenMap := make(map[string]common.Address, len(pairs))

	for _, pair := range pairs {
		parts := strings.Split(pair.Name, "/")
		if len(parts) != 2 {
			continue
		}

		tokenMap[parts[0]] = common.HexToAddress("0x" + pair.MakerAssetData[len(pair.MakerAssetData)-40:])
		tokenMap[parts[1]] = common.HexToAddress("0x" + pair.TakerAssetData[len(pair.TakerAssetData)-40:])
	}

	if ctl.ethClient != nil {
		// always override WETH with client-side configured address
		tokenMap["WETH"] = ctl.ethClient.contractAddresses[EthContractWETH9]
	}

	tokenNames = make([]string, 0, len(tokenMap))
	for name := range tokenMap {
		tokenNames = append(tokenNames, name)
	}

	sort.Strings(tokenNames)

	assets = make([]common.Address, 0, len(tokenNames))

	for _, name := range tokenNames {
		assets = append(assets, tokenMap[name])
	}

	return tokenNames, assets, nil
}

func (ctl *AppController) ActionTradeTokens() {
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tokenNames, assets, err := ctl.getTokenNamesAndAssets(ctx)
	if err == ErrClientUnavailable {
		logrus.Errorln("Ethereum client is not initialized")
		return
	} else if err != nil {
		logrus.Errorln(err)
		return
	}

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	ethBalance, err := ctl.ethClient.EthBalance(ctx, defaultAccount)
	if err != nil {
		logrus.WithError(err).WithField("account", defaultAccount.Hex()).Warningln("Unable to get Eth balance")
	}

	var ethBalanceStr string = "-"
	if ethBalance != nil {
		ethBalanceDec := decimal.NewFromBigInt(ethBalance, 0)
		ethBalanceDec = ethBalanceDec.Div(decimal.New(1, 18))
		ethBalanceStr = ethBalanceDec.StringFixed(8)
	}

	networkName := ctl.mustConfigValue("networks.default")
	proxyAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.erc20proxy_address", networkName))

	allowances := ctl.ethClient.AllowancesMap(ctx, defaultAccount, common.HexToAddress(proxyAddressHex), assets)
	balances := ctl.ethClient.BalancesMap(ctx, defaultAccount, assets)

	if len(balances) == 0 && len(allowances) == 0 {
		fmt.Println("No token info available.")
		return
	}

	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddTitle(
		fmt.Sprintf("Account %s (%s ETH)", defaultAccount.Hex(), ethBalanceStr),
	)
	table.AddHeaders("Token", "Address", "Balance", "Unlocked")

	for idx, name := range tokenNames {
		addr := assets[idx]

		var balanceStr string = "-"
		var unlockedStr string = " "

		if balances[addr] != nil {
			balanceDec := decimal.NewFromBigInt(balances[addr], 0)
			balanceDec = balanceDec.Div(decimal.New(1, 18))
			balanceStr = balanceDec.StringFixed(8)
		}

		if allowances[addr] != nil {
			isUnlocked := (allowances[addr].Cmp(UnlimitedAllowance) == 0)
			if isUnlocked {
				unlockedStr = "x"
			}
		}

		table.AddRow(
			name,
			addr.Hex(),
			balanceStr,
			fmt.Sprintf("[%s]", unlockedStr),
		)
	}

	fmt.Println(table.Render())
}

func (ctl *AppController) ActionTradePairs() {
	ctx := context.Background()
	pairs, err := ctl.relayerClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("failed to list trade pairs")
		return
	}

	for idx, pair := range pairs {
		fmt.Printf("%d) %s\n", idx+1, pair.Name)
	}
}

type UtilTokenLockArgs struct {
	TokenName string
	Password  string
}

func (ctl *AppController) ActionUtilLock(args interface{}) {
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tokenNames, assets, err := ctl.getTokenNamesAndAssets(ctx)
	if err == ErrClientUnavailable {
		logrus.Errorln("Ethereum client is not initialized")
		return
	} else if err != nil {
		logrus.Errorln(err)
		return
	}

	tokenLockArgs := args.(*UtilTokenLockArgs)
	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: tokenLockArgs.Password,
		GasPrice: ctl.ethGasPrice,
	}

	var asset common.Address
	tokenName := strings.ToLower(tokenLockArgs.TokenName)
	for idx, name := range tokenNames {
		if strings.ToLower(name) == tokenName {
			asset = assets[idx]
			break
		}
	}

	txHash, err := ctl.ethClient.TokenLock(callArgs, asset)
	if err != nil {
		logrus.WithError(err).Errorln("unable to lock token")
		return
	}

	fmt.Println(ctl.formatTxLink(txHash))
	ctl.checkTx(txHash)
}

type UtilTokenUnlockArgs struct {
	TokenName string
	Password  string
}

func (ctl *AppController) ActionUtilUnlock(args interface{}) {
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tokenNames, assets, err := ctl.getTokenNamesAndAssets(ctx)
	if err == ErrClientUnavailable {
		logrus.Errorln("Ethereum client is not initialized")
		return
	} else if err != nil {
		logrus.Errorln(err)
		return
	}

	tokenUnlockArgs := args.(*UtilTokenUnlockArgs)
	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: tokenUnlockArgs.Password,
		GasPrice: ctl.ethGasPrice,
	}

	var asset common.Address
	tokenName := strings.ToLower(tokenUnlockArgs.TokenName)
	for idx, name := range tokenNames {
		if strings.ToLower(name) == tokenName {
			asset = assets[idx]
			break
		}
	}

	txHash, err := ctl.ethClient.TokenUnlock(callArgs, asset)
	if err != nil {
		logrus.WithError(err).Errorln("unable to unlock token")
		return
	}

	fmt.Println(ctl.formatTxLink(txHash))
	ctl.checkTx(txHash)
}

type UtilWrapArgs struct {
	Amount   string
	Password string
}

func (ctl *AppController) ActionUtilWrap(args interface{}) {
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	ethWrapArgs := args.(*UtilWrapArgs)
	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: ethWrapArgs.Password,
		GasPrice: ctl.ethGasPrice,
	}

	var amount *big.Int

	amountDec, err := decimal.NewFromString(ethWrapArgs.Amount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse amount")
		return
	} else if amountDec.LessThan(decimal.RequireFromString("0.0000001")) {
		logrus.Errorln("amount is too small, must be at least 0.0000001 ETH")
		return
	} else {
		amount = dec2big(amountDec)
	}

	txHash, err := ctl.ethClient.EthWrap(callArgs, amount)
	if err != nil {
		logrus.WithError(err).Errorln("unable to wrap ETH")
		return
	}

	fmt.Println(ctl.formatTxLink(txHash))
	ctl.checkTx(txHash)
}

type UtilUnwrapArgs struct {
	Amount   string
	Password string
}

func (ctl *AppController) ActionUtilUnwrap(args interface{}) {
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	ethUnwrapArgs := args.(*UtilUnwrapArgs)
	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: ethUnwrapArgs.Password,
		GasPrice: ctl.ethGasPrice,
	}

	var amount *big.Int

	amountDec, err := decimal.NewFromString(ethUnwrapArgs.Amount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse amount")
		return
	} else if amountDec.LessThan(decimal.RequireFromString("0.0000001")) {
		logrus.Errorln("amount is too small, must be at least 0.0000001 WETH")
		return
	} else {
		amount = dec2big(amountDec)
	}

	txHash, err := ctl.ethClient.EthUnwrap(callArgs, amount)
	if err != nil {
		logrus.WithError(err).Errorln("unable to wrap WETH")
		return
	}

	fmt.Println(ctl.formatTxLink(txHash))
	ctl.checkTx(txHash)
}

func (ctl *AppController) ActionAccountsUse(args interface{}) {
	addr, err := ethParseAccount(args.(*AccountUseArgs))
	if err != nil {
		logrus.WithError(err).Errorln("failed to select default account")
		return
	}

	found := false
	allAccounts := ctl.keystore.Accounts()

	for _, acc := range allAccounts {
		if bytes.Equal(acc.Bytes(), addr.Bytes()) {
			found = true
			break
		}
	}

	if !found {
		logrus.WithField("address", addr.Hex()).Errorln("specified account not found in keystore")
		return
	}

	ctl.setConfigValue("accounts.default", addr.Hex())

	if err := saveConfig(ctl.configPath, ctl.cfg); err != nil {
		logrus.WithError(err).Errorln("failed to save config file")
	}

	if ctl.ethClient != nil {
		ctl.ethClient.SetDefaultFromAddress(addr)
	} else if err := ctl.initEthClient(); err != nil {
		logrus.WithError(err).Warningln("failed to init Ethereum client")
	}
}

func (ctl *AppController) ActionAccountsCreate(args interface{}) {
	acc, err := ethCreateAccount(ctl.keystorePath, args.(*AccountCreateArgs))
	if err != nil {
		logrus.WithError(err).Errorln("failed to create new account")
		return
	}

	if ctl.takeFirstAccountAsDefault() {
		saveConfig(ctl.configPath, ctl.cfg)
	}

	logrus.Infof("Created a new account: %s", acc.Address.Hex())
}

func (ctl *AppController) ActionAccountsImport(args interface{}) {
	addr, err := ethImportAccount(ctl.keystorePath, args.(*AccountImportArgs))
	if err != nil {
		logrus.WithError(err).Errorln("failed to import account")
		return
	}

	if ctl.takeFirstAccountAsDefault() {
		saveConfig(ctl.configPath, ctl.cfg)
	}

	logrus.Infof("Imported an existing keyfile: %s", addr.Hex())
}

func (ctl *AppController) ActionAccountsImportPrivKey(args interface{}) {
	addr, err := ethImportPrivKey(ctl.keystorePath, args.(*AccountImportPrivKeyArgs))
	if err != nil {
		logrus.WithError(err).Errorln("failed to import private key")
		return
	}

	if ctl.takeFirstAccountAsDefault() {
		saveConfig(ctl.configPath, ctl.cfg)
	}

	logrus.Infof("Imported private key for %s", addr.Hex())
}

func (ctl *AppController) ActionAccountsList() {
	allAccounts := ctl.keystore.Accounts()
	if len(allAccounts) == 0 {
		fmt.Printf("No accounts in %s\n", ctl.mustConfigValue("accounts.keystore"))
		return
	}

	for idx, acc := range allAccounts {
		fmt.Printf("%d) %s\n", idx+1, acc.Hex())
	}

	defaultAccount := ctl.mustConfigValue("accounts.default")

	fmt.Printf("\nUsing the default account: %s\n", defaultAccount)
}

func (ctl *AppController) SuggestAccounts() []prompt.Suggest {
	allAccounts := ctl.keystore.Accounts()
	suggestions := make([]prompt.Suggest, len(allAccounts))

	for i, addr := range allAccounts {
		suggestions[i].Text = addr.Hex()
	}

	return suggestions
}

func (ctl *AppController) SuggestTokens() []prompt.Suggest {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	tokenNames, _, err := ctl.getTokenNamesAndAssets(ctx)
	if err != nil {
		return nil
	}

	suggestions := make([]prompt.Suggest, len(tokenNames))
	for i, name := range tokenNames {
		suggestions[i].Text = name
	}

	return suggestions
}

func (ctl *AppController) SuggestMarkets() []prompt.Suggest {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	pairs, err := ctl.relayerClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Warningln("failed to fetch trade pairs")
		return nil
	}

	suggestions := make([]prompt.Suggest, 0, len(pairs))
	for _, pair := range pairs {
		if !pair.Enabled {
			continue
		}

		suggestions = append(suggestions, prompt.Suggest{
			Text: pair.Name,
		})
	}

	return suggestions
}

func (ctl *AppController) takeFirstAccountAsDefault() bool {
	_, ok := ctl.getConfigValue("accounts.default")
	if !ok {
		allAccounts := ctl.keystore.Accounts()
		if len(allAccounts) > 0 {
			ctl.setConfigValue("accounts.default", allAccounts[0].Hex())
			return true
		} else {
			logrus.WithField("keystore", ctl.keystorePath).Infoln("No accounts found in keystore yet")
			return false
		}
	} else {
		// already set
		return false
	}

	return false
}

func (ctl *AppController) selectDefaultNetwork() bool {
	_, ok := ctl.getConfigValue("networks.default")
	if !ok {
		if _, hasMainnet := ctl.getConfigValue("networks.mainnet.endpoint"); hasMainnet {
			ctl.setConfigValue("networks.default", "mainnet")
			return true
		}

		if _, hasDevnet := ctl.getConfigValue("networks.devnet.endpoint"); hasDevnet {
			ctl.setConfigValue("networks.default", "devnet")
			return true
		}

		logrus.Infoln("No default network found in config")
		return false
	} else {
		// already set
		return false
	}

	return false
}

var (
	ErrTxTimeout   = errors.New("timeout while waiting for tx confirmation")
	ErrTxBadStatus = errors.New("tx execution ended with failing status code")
)

func (ctl *AppController) checkTx(txHash common.Hash) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 2*time.Minute)
	spinDone := makeSpin(ctx, "checking tx")

	if err := ctl.awaitTx(ctx, txHash); err != nil {
		switch err {
		case ErrTxTimeout:
			logrus.Warningln("unable to check tx confirmation, use explorer link above to check manually")
		case ErrTxBadStatus:
			logrus.Errorln("transaction has failed with bad status, check logs using explorer link above")
		default:
			logrus.WithError(err).Warningln("unable to check tx confirmation")
		}
	}

	cancelFn()
	<-spinDone
}

func (ctl *AppController) awaitTx(ctx context.Context, txHash common.Hash) error {
	tx, err := ctl.ethClient.ethManager.TransactionByHash(ctx, txHash.Hex())
	if err != nil {
		return err
	}

	t := time.NewTimer(time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			tx, err = ctl.ethClient.ethManager.TransactionByHash(ctx, txHash.Hex())
			if err == nil && tx.BlockNumber != nil {
				receipt, err := ctl.ethClient.ethManager.TransactionReceiptByHash(ctx, txHash.Hex())
				if err != nil {
					err = errors.Wrap(err, "failed to get tx receipt")
					return err
				} else if status := receipt.Status; status == 0 {
					return ErrTxBadStatus
				}

				// finally a transaction receipt,
				// with a successful status
				return nil
			} else if err != nil {
				logrus.WithError(err).Warningln("error while checking, retry in 10 seconds")
				t.Reset(10 * time.Second)

				continue
			}

			t.Reset(time.Second)
		case <-ctx.Done():
			if ctx.Err() != context.Canceled {
				return ErrTxTimeout
			}

			return nil
		}
	}

	return nil
}

func (ctl *AppController) formatTxLink(txHash common.Hash) string {
	networkName := ctl.mustConfigValue("networks.default")
	explorerEndpoint, ok := ctl.getConfigValue(fmt.Sprintf("networks.%s.explorer", networkName))
	if ok && len(explorerEndpoint) > 0 {
		return explorerEndpoint + txHash.Hex()
	}

	return txHash.Hex()
}

func (ctl *AppController) initEthClient() error {
	networkName := ctl.mustConfigValue("networks.default")

	ethEndpoint := ctl.mustConfigValue(fmt.Sprintf("networks.%s.endpoint", networkName))
	ethGasPrice := ctl.mustConfigValue(fmt.Sprintf("networks.%s.gas_price", networkName))

	exchangeAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.exchange_address", networkName))
	weth9AddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.weth9_address", networkName))
	erc20ProxyAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.erc20proxy_address", networkName))
	contractAddresses := map[EthContract]common.Address{
		EthContractWETH9:      common.HexToAddress(weth9AddressHex),
		EthContractERC20Proxy: common.HexToAddress(erc20ProxyAddressHex),
		EthContractExchange:   common.HexToAddress(exchangeAddressHex),
	}

	ethManager := manager.NewManager([]string{
		ethEndpoint,
	}, 100000) // allow only 100k gas per tx

	defaultFromAddress := common.HexToAddress(ctl.mustConfigValue("accounts.default"))
	allowGasOracles, _ := ctl.getConfigValue("networks.allow_gas_oracles", "true")

	var gasPriceParsed bool
	ctl.ethGasPrice, gasPriceParsed = big.NewInt(0).SetString(ethGasPrice, 10)
	if !gasPriceParsed {
		ctl.ethGasPrice = nil
	}

	ethClient, err := NewEthClient(
		ctl.keystore,
		ethManager,
		defaultFromAddress,
		contractAddresses,
		toBool(allowGasOracles),
	)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"chain_id": ethManager.ChainID(),
	}).Debugf("Connected to %s", networkName)

	ctl.ethClient = ethClient

	return nil
}

func (ctl *AppController) setConfigValue(path string, v interface{}) {
	if ctl.cfg == nil {
		ctl.cfg, _ = toml.TreeFromMap(map[string]interface{}{})
	}

	ctl.cfg.Set(path, v)
}

func (ctl *AppController) mustConfigValue(path string) string {
	val, ok := ctl.getConfigValue(path)
	if !ok {
		logrus.WithField("path", path).Fatalf("config value not found but required")
	}

	return val
}

func (ctl *AppController) getConfigValue(path string, fallback ...interface{}) (string, bool) {
	optVal, optOk := appConfigSetMap[path]
	if optOk && optVal.SetByUser != nil && *optVal.SetByUser {
		return *appConfigMap[path], true
	}

	var v interface{}

	if len(fallback) != 0 {
		v = ctl.cfg.GetDefault(path, fallback[0])
	} else if optOk && len(optVal.Value) > 0 {
		v = ctl.cfg.GetDefault(path, optVal.Value)
	} else {
		hasPath := ctl.cfg.Has(path)
		if !hasPath {
			return "", false
		}

		v = ctl.cfg.Get(path)
	}

	return v.(string), true
}

func makeSpin(ctx context.Context, label string) <-chan struct{} {
	s := spin.New()
	doneC := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != context.Canceled {
					fmt.Printf("\r  \033[36m%s\033[m timeout!\n", label)
					close(doneC)
					return
				}
				fmt.Printf("\r  \033[36m%s\033[m done\n", label)
				close(doneC)
				return
			default:
				fmt.Printf("\r  \033[36m%s\033[m %s", label, s.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return doneC
}
