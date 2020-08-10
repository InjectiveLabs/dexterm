package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cli "github.com/jawher/mow.cli"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var configPath = app.String(cli.StringOpt{
	Name:   "C config",
	Desc:   "Specify config path",
	EnvVar: "DEXTERM_CONFIG",
	Value:  "~/.dexterm/config.toml",
})

var (
	logDebugSet bool
	logDebugOpt = cli.StringOpt{
		Name:      "d debug",
		Desc:      "Enable debug logging.",
		EnvVar:    "DEXTERM_DEBUG",
		Value:     "false",
		SetByUser: &logDebugSet,
	}
)

var (
	relayerEndpointSet bool
	relayerEndpointOpt = cli.StringOpt{
		Name:      "R relayer-endpoint",
		Desc:      "Specify DEX API endpoint for the session.",
		EnvVar:    "DEXTERM_API_ENDPOINT",
		Value:     "http://localhost:4444",
		SetByUser: &relayerEndpointSet,
	}
)

//var (
//	derivativesEndpointSet bool
//	derivativesEndpointOpt = cli.StringOpt{
//		Name:      "D derivatives-endpoint",
//		Desc:      "Specify DEX derivatives API endpoint for the session.",
//		EnvVar:    "DEXTERM_DERIVATIVES_API_ENDPOINT",
//		Value:     "http://localhost:4444",
//		SetByUser: &derivativesEndpointSet,
//	}
//)

var (
	accountsKeystoreSet bool
	accountsKeystoreOpt = cli.StringOpt{
		Name:      "K keystore",
		Desc:      "Specify path for keystore, a dir containing private keys.",
		EnvVar:    "DEXTERM_KEYSTORE_PATH",
		Value:     "~/.dexterm/keystore",
		SetByUser: &accountsKeystoreSet,
	}
)

var (
	accountsDefaultSet bool
	accountsDefaultOpt = cli.StringOpt{
		Name:      "A account",
		Desc:      "Specify account address to use for all transactions within the session.",
		EnvVar:    "DEXTERM_ACCOUNT_ADDRESS",
		Value:     "",
		SetByUser: &accountsDefaultSet,
	}
)

var appConfigMap = map[string]*string{
	"log.debug": app.String(logDebugOpt),

	"relayer.endpoint": app.String(relayerEndpointOpt),
	//"derivatives.endpoint": app.String(derivativesEndpointOpt),

	"accounts.keystore": app.String(accountsKeystoreOpt),
	"accounts.default":  app.String(accountsDefaultOpt),

	"networks.allow_gas_oracles": app.String(networksAllowGasOraclesOpt),
	"networks.default":           app.String(networksDefaultOpt),

	"networks.mainnet.endpoint":            app.String(networksMainnetEndpointOpt),
	"networks.mainnet.explorer":            app.String(networksMainnetExplorerOpt),
	"networks.mainnet.gas_price":           app.String(networksMainnetGasPriceOpt),
	"networks.mainnet.weth9_address":       app.String(networksMainnetWETH9Opt),
	"networks.mainnet.erc20proxy_address":  app.String(networksMainnetERC20ProxyOpt),
	"networks.mainnet.exchange_address":    app.String(networksMainnetExchangeOpt),
	"networks.mainnet.futures_address":     app.String(networksMainnetFuturesOpt),
	"networks.mainnet.coordinator_address": app.String(networksMainnetCoordinatorOpt),

	"networks.ropsten.endpoint":            app.String(networksRopstenEndpointOpt),
	"networks.ropsten.explorer":            app.String(networksRopstenExplorerOpt),
	"networks.ropsten.gas_price":           app.String(networksRopstenGasPriceOpt),
	"networks.ropsten.weth9_address":       app.String(networksRopstenWETH9Opt),
	"networks.ropsten.erc20proxy_address":  app.String(networksRopstenERC20ProxyOpt),
	"networks.ropsten.exchange_address":    app.String(networksRopstenExchangeOpt),
	"networks.ropsten.futures_address":     app.String(networksRopstenFuturesOpt),
	"networks.ropsten.coordinator_address": app.String(networksRopstenCoordinatorOpt),

	"networks.kovan.endpoint":            app.String(networksKovanEndpointOpt),
	"networks.kovan.explorer":            app.String(networksKovanExplorerOpt),
	"networks.kovan.gas_price":           app.String(networksKovanGasPriceOpt),
	"networks.kovan.weth9_address":       app.String(networksKovanWETH9Opt),
	"networks.kovan.erc20proxy_address":  app.String(networksKovanERC20ProxyOpt),
	"networks.kovan.exchange_address":    app.String(networksKovanExchangeOpt),
	"networks.kovan.futures_address":     app.String(networksKovanFuturesOpt),
	"networks.kovan.coordinator_address": app.String(networksKovanCoordinatorOpt),

	"networks.devnet.endpoint":            app.String(networksDevnetEndpointOpt),
	"networks.devnet.explorer":            app.String(networksDevnetExplorerOpt),
	"networks.devnet.gas_price":           app.String(networksDevnetGasPriceOpt),
	"networks.devnet.weth9_address":       app.String(networksDevnetWETH9Opt),
	"networks.devnet.erc20proxy_address":  app.String(networksDevnetERC20ProxyOpt),
	"networks.devnet.exchange_address":    app.String(networksDevnetExchangeOpt),
	"networks.devnet.futures_address":     app.String(networksDevnetFuturesOpt),
	"networks.devnet.coordinator_address": app.String(networksDevnetCoordinatorOpt),

	"networks.injective.endpoint":            app.String(networksInjectiveEndpointOpt),
	"networks.injective.explorer":            app.String(networksInjectiveExplorerOpt),
	"networks.injective.gas_price":           app.String(networksInjectiveGasPriceOpt),
	"networks.injective.weth9_address":       app.String(networksInjectiveWETH9Opt),
	"networks.injective.erc20proxy_address":  app.String(networksInjectiveERC20ProxyOpt),
	"networks.injective.exchange_address":    app.String(networksInjectiveExchangeOpt),
	"networks.injective.futures_address":     app.String(networksInjectiveFuturesOpt),
	"networks.injective.coordinator_address": app.String(networksInjectiveCoordinatorOpt),

	"networks.matic.endpoint":            app.String(networksMaticEndpointOpt),
	"networks.matic.explorer":            app.String(networksMaticExplorerOpt),
	"networks.matic.gas_price":           app.String(networksMaticGasPriceOpt),
	"networks.matic.weth9_address":       app.String(networksMaticWETH9Opt),
	"networks.matic.erc20proxy_address":  app.String(networksMaticERC20ProxyOpt),
	"networks.matic.exchange_address":    app.String(networksMaticExchangeOpt),
	"networks.matic.futures_address":     app.String(networksMaticFuturesOpt),
	"networks.matic.coordinator_address": app.String(networksMaticCoordinatorOpt),
}

var appConfigSetMap = map[string]cli.StringOpt{
	"log.debug": logDebugOpt,

	"relayer.endpoint": relayerEndpointOpt,

	"accounts.keystore": accountsKeystoreOpt,
	"accounts.default":  accountsDefaultOpt,

	"networks.allow_gas_oracles": networksAllowGasOraclesOpt,
	"networks.default":           networksDefaultOpt,

	"networks.mainnet.endpoint":            networksMainnetEndpointOpt,
	"networks.mainnet.explorer":            networksMainnetExplorerOpt,
	"networks.mainnet.gas_price":           networksMainnetGasPriceOpt,
	"networks.mainnet.weth9_address":       networksMainnetWETH9Opt,
	"networks.mainnet.erc20proxy_address":  networksMainnetERC20ProxyOpt,
	"networks.mainnet.exchange_address":    networksMainnetExchangeOpt,
	"networks.mainnet.devutils_address":    networksMainnetDevUtilsOpt,
	"networks.mainnet.futures_address":     networksMainnetFuturesOpt,
	"networks.mainnet.coordinator_address": networksMainnetCoordinatorOpt,

	"networks.ropsten.endpoint":            networksRopstenEndpointOpt,
	"networks.ropsten.explorer":            networksRopstenExplorerOpt,
	"networks.ropsten.gas_price":           networksRopstenGasPriceOpt,
	"networks.ropsten.weth9_address":       networksRopstenWETH9Opt,
	"networks.ropsten.erc20proxy_address":  networksRopstenERC20ProxyOpt,
	"networks.ropsten.exchange_address":    networksRopstenExchangeOpt,
	"networks.ropsten.coordinator_address": networksRopstenCoordinatorOpt,

	"networks.kovan.endpoint":            networksKovanEndpointOpt,
	"networks.kovan.explorer":            networksKovanExplorerOpt,
	"networks.kovan.gas_price":           networksKovanGasPriceOpt,
	"networks.kovan.weth9_address":       networksKovanWETH9Opt,
	"networks.kovan.erc20proxy_address":  networksKovanERC20ProxyOpt,
	"networks.kovan.exchange_address":    networksKovanExchangeOpt,
	"networks.kovan.devutils_address":    networksKovanDevUtilsOpt,
	"networks.kovan.futures_address":     networksKovanFuturesOpt,
	"networks.kovan.coordinator_address": networksKovanCoordinatorOpt,

	"networks.devnet.endpoint":            networksDevnetEndpointOpt,
	"networks.devnet.explorer":            networksDevnetExplorerOpt,
	"networks.devnet.gas_price":           networksDevnetGasPriceOpt,
	"networks.devnet.weth9_address":       networksDevnetWETH9Opt,
	"networks.devnet.erc20proxy_address":  networksDevnetERC20ProxyOpt,
	"networks.devnet.exchange_address":    networksDevnetExchangeOpt,
	"networks.devnet.devutils_address":    networksDevnetDevUtilsOpt,
	"networks.devnet.futures_address":     networksDevnetFuturesOpt,
	"networks.devnet.coordinator_address": networksDevnetCoordinatorOpt,

	"networks.injective.endpoint":            networksInjectiveEndpointOpt,
	"networks.injective.explorer":            networksInjectiveExplorerOpt,
	"networks.injective.gas_price":           networksInjectiveGasPriceOpt,
	"networks.injective.weth9_address":       networksInjectiveWETH9Opt,
	"networks.injective.erc20proxy_address":  networksInjectiveERC20ProxyOpt,
	"networks.injective.exchange_address":    networksInjectiveExchangeOpt,
	"networks.injective.devutils_address":    networksInjectiveDevUtilsOpt,
	"networks.injective.futures_address":     networksInjectiveFuturesOpt,
	"networks.injective.coordinator_address": networksInjectiveCoordinatorOpt,

	"networks.matic.endpoint":            networksMaticEndpointOpt,
	"networks.matic.explorer":            networksMaticExplorerOpt,
	"networks.matic.gas_price":           networksMaticGasPriceOpt,
	"networks.matic.weth9_address":       networksMaticWETH9Opt,
	"networks.matic.erc20proxy_address":  networksMaticERC20ProxyOpt,
	"networks.matic.exchange_address":    networksMaticExchangeOpt,
	"networks.matic.coordinator_address": networksMaticCoordinatorOpt,
}

func loadOrCreateConfig(configPath string) (*toml.Tree, error) {
	var cfg *toml.Tree

	cfgFile, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if cfg, err = writeInitialConfig(configPath, appConfigMap); err != nil {
				err = errors.Wrap(err, "failed to load or create initial config")
				return nil, err
			}

			return cfg, nil
		}

		err = errors.Wrap(err, "failed to open config for reading")
		return nil, err
	}

	defer cfgFile.Close()

	if cfg, err = toml.LoadReader(cfgFile); err != nil {
		err = errors.Wrap(err, "failed to parse config")
		return nil, err
	}

	return cfg, nil
}

func writeInitialConfig(configPath string, configMapping map[string]*string) (*toml.Tree, error) {
	baseDir := filepath.Dir(configPath)
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		err = errors.Wrap(err, "unable to mkdir for initial config")

		return nil, err
	}

	cfgFile, err := os.Create(filepath.Join(baseDir, filepath.Base(configPath)))
	if err != nil {
		err = errors.Wrap(err, "unable to open config file")

		return nil, err
	}
	defer cfgFile.Close()

	tree, _ := toml.TreeFromMap(map[string]interface{}{})

	for path, val := range configMapping {
		if len(*val) > 0 {
			tree.Set(path, *val)
		}
	}

	if _, err = tree.WriteTo(cfgFile); err != nil {
		logrus.WithError(err).Warningln("failed to write config contents")
	}

	return tree, nil
}

func saveConfig(configPath string, config *toml.Tree) error {
	cfgFile, err := os.OpenFile(configPath, os.O_EXCL|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		err = errors.Wrap(err, "unable to open config file")

		return err
	}
	defer cfgFile.Close()

	if _, err = config.WriteTo(cfgFile); err != nil {
		return err
	}

	return nil
}

func configSubsections(config *toml.Tree, section string) ([]string, bool) {
	if !config.Has(section) {
		return nil, false
	}

	v := config.Get(section)
	fmt.Printf("subsection type %T contents: %+v\n", v, v)

	return nil, true
}

func toBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "1", "t", "yes", "on":
		return true
	default:
		return false
	}
}
