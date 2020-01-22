package main

import (
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

	"accounts.keystore": app.String(accountsKeystoreOpt),
	"accounts.default":  app.String(accountsDefaultOpt),

	"networks.allow_gas_oracles": app.String(networksAllowGasOraclesOpt),
	"networks.default":           app.String(networksDefaultOpt),

	"networks.mainnet.endpoint":           app.String(networksMainnetEndpointOpt),
	"networks.mainnet.gas_price":          app.String(networksMainnetGasPriceOpt),
	"networks.mainnet.weth9_address":      app.String(networksMainnetWETH9Opt),
	"networks.mainnet.erc20proxy_address": app.String(networksMainnetERC20ProxyOpt),

	"networks.ropsten.endpoint":           app.String(networksRopstenEndpointOpt),
	"networks.ropsten.gas_price":          app.String(networksRopstenGasPriceOpt),
	"networks.ropsten.weth9_address":      app.String(networksRopstenWETH9Opt),
	"networks.ropsten.erc20proxy_address": app.String(networksRopstenERC20ProxyOpt),

	"networks.devnet.endpoint":           app.String(networksDevnetEndpointOpt),
	"networks.devnet.gas_price":          app.String(networksDevnetGasPriceOpt),
	"networks.devnet.weth9_address":      app.String(networksDevnetWETH9Opt),
	"networks.devnet.erc20proxy_address": app.String(networksDevnetERC20ProxyOpt),

	"networks.matic.endpoint":           app.String(networksMaticEndpointOpt),
	"networks.matic.gas_price":          app.String(networksMaticGasPriceOpt),
	"networks.matic.weth9_address":      app.String(networksMaticWETH9Opt),
	"networks.matic.erc20proxy_address": app.String(networksMaticERC20ProxyOpt),
}

var appConfigSetMap = map[string]cli.StringOpt{
	"log.debug": logDebugOpt,

	"relayer.endpoint": relayerEndpointOpt,

	"accounts.keystore": accountsKeystoreOpt,
	"accounts.default":  accountsDefaultOpt,

	"networks.allow_gas_oracles": networksAllowGasOraclesOpt,
	"networks.default":           networksDefaultOpt,

	"networks.mainnet.endpoint":           networksMainnetEndpointOpt,
	"networks.mainnet.gas_price":          networksMainnetGasPriceOpt,
	"networks.mainnet.weth9_address":      networksMainnetWETH9Opt,
	"networks.mainnet.erc20proxy_address": networksMainnetERC20ProxyOpt,

	"networks.ropsten.endpoint":           networksRopstenEndpointOpt,
	"networks.ropsten.gas_price":          networksRopstenGasPriceOpt,
	"networks.ropsten.weth9_address":      networksRopstenWETH9Opt,
	"networks.ropsten.erc20proxy_address": networksRopstenERC20ProxyOpt,

	"networks.devnet.endpoint":           networksDevnetEndpointOpt,
	"networks.devnet.gas_price":          networksDevnetGasPriceOpt,
	"networks.devnet.weth9_address":      networksDevnetWETH9Opt,
	"networks.devnet.erc20proxy_address": networksDevnetERC20ProxyOpt,

	"networks.matic.endpoint":           networksMaticEndpointOpt,
	"networks.matic.gas_price":          networksMaticGasPriceOpt,
	"networks.matic.weth9_address":      networksMaticWETH9Opt,
	"networks.matic.erc20proxy_address": networksMaticERC20ProxyOpt,
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

func toBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "1", "t", "yes":
		return true
	default:
		return false
	}
}
