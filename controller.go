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
	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	toml "github.com/pelletier/go-toml"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"

	"github.com/InjectiveLabs/dexterm/ethfw/keystore"
	"github.com/InjectiveLabs/dexterm/ethfw/manager"
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

func (ctl *AppController) ActionTradeSell(args interface{}) {

}

func (ctl *AppController) ActionTradeBuy(args interface{}) {

}

func (ctl *AppController) ActionTradeOrderbook(args interface{}) {

}

func (ctl *AppController) ActionTradeTokens() {
	if ctl.ethClient == nil {
		logrus.Errorln("Etherteum client is not initialized")
		return
	}

	ctx := context.Background()
	pairs, err := ctl.relayerClient.TradePairs(ctx)

	if err != nil {
		logrus.WithError(err).Errorln("failed to list trade pairs")
		return
	}

	tokenMap := make(map[string]common.Address, len(pairs))

	for _, pair := range pairs {
		parts := strings.Split(pairs[0].Name, "/")
		if len(parts) != 2 {
			continue
		}

		tokenMap[parts[0]] = common.HexToAddress("0x" + pair.MakerAssetData[len(pair.MakerAssetData)-40:])
		tokenMap[parts[1]] = common.HexToAddress("0x" + pair.TakerAssetData[len(pair.TakerAssetData)-40:])
	}

	// always override WETH with client-side configured address
	tokenMap["WETH"] = ctl.ethClient.contractAddresses[EthContractWETH9]

	tokenNames := make([]string, 0, len(tokenMap))
	for name := range tokenMap {
		tokenNames = append(tokenNames, name)
	}

	sort.Strings(tokenNames)

	assets := make([]common.Address, 0, len(tokenNames))

	for _, name := range tokenNames {
		assets = append(assets, tokenMap[name])
	}

	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

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

	fmt.Printf("Account %s\n", defaultAccount.Hex())
	fmt.Println("ETH available:", ethBalanceStr)

	networkName := ctl.mustConfigValue("networks.default")
	proxyAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.erc20proxy_address", networkName))

	allowances, err := ctl.ethClient.AllowancesMap(ctx, defaultAccount, common.HexToAddress(proxyAddressHex), assets)
	if err != nil {
		logrus.WithError(err).Errorln("Unable to quote allowances")
		return
	}

	balances, err := ctl.ethClient.BalancesMap(ctx, defaultAccount, assets)
	if err != nil {
		logrus.WithError(err).Errorln("Unable to quote balances")
		return
	}

	if len(balances) == 0 && len(allowances) == 0 {
		fmt.Println("No token info available.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Token", "Address", "Balance", "Unlocked"})

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

		table.Append([]string{
			name,
			addr.Hex(),
			balanceStr,
			fmt.Sprintf("[%s]", unlockedStr),
		})
	}

	table.Render()
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

func (ctl *AppController) ActionUtilUnlock(args interface{}) {

}

func (ctl *AppController) ActionUtilLock(args interface{}) {

}

func (ctl *AppController) ActionUtilWrap(args interface{}) {

}

func (ctl *AppController) ActionUtilUnwrap(args interface{}) {

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

func (ctl *AppController) initEthClient() error {
	networkName := ctl.mustConfigValue("networks.default")

	ethEndpoint := ctl.mustConfigValue(fmt.Sprintf("networks.%s.endpoint", networkName))
	ethGasPrice := ctl.mustConfigValue(fmt.Sprintf("networks.%s.gas_price", networkName))

	weth9AddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.weth9_address", networkName))
	erc20ProxyAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.erc20proxy_address", networkName))
	contractAddresses := map[EthContract]common.Address{
		EthContractWETH9:      common.HexToAddress(weth9AddressHex),
		EthContractERC20Proxy: common.HexToAddress(erc20ProxyAddressHex),
	}

	ethManager := manager.NewManager([]string{
		ethEndpoint,
	}, 6721975)

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
