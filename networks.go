package main

import (
	cli "github.com/jawher/mow.cli"
)

var (
	networksAllowGasOraclesSet bool
	networksAllowGasOraclesOpt = cli.StringOpt{
		Name:      "allow-gas-oracles",
		Desc:      "This option enables fetching gas prices from external APIs like EthGasStation.",
		EnvVar:    "DEXTERM_ALLOW_GAS_ORACLES",
		Value:     "on",
		SetByUser: &networksAllowGasOraclesSet,
	}
)

var (
	networksDefaultSet bool
	networksDefaultOpt = cli.StringOpt{
		Name:      "N network",
		Desc:      "Specify network name to use for all transactions withing the session.",
		EnvVar:    "DEXTERM_NETWORK",
		Value:     "mainnet",
		SetByUser: &networksDefaultSet,
	}
)

var (
	networksMainnetEndpointSet bool
	networksMainnetEndpointOpt = cli.StringOpt{
		Name:      "mainnet-endpoint",
		Desc:      "Specify endpoint for MainNet network",
		EnvVar:    "DEXTERM_MAINNET_ENDPOINT",
		Value:     "https://eth-mainnet.alchemyapi.io/jsonrpc/DqEv1TiHskO-G6JprqyhE25k1x0p3hpj",
		SetByUser: &networksMainnetEndpointSet,
	}
)

var (
	networksMainnetGasPriceSet bool
	networksMainnetGasPriceOpt = cli.StringOpt{
		Name:      "mainnet-gasprice",
		Desc:      "Specify min gasprice for MainNet network",
		EnvVar:    "DEXTERM_MAINNET_GASPRICE",
		Value:     "10000000000",
		SetByUser: &networksMainnetGasPriceSet,
	}
)

var (
	networksMainnetWETH9Set bool
	networksMainnetWETH9Opt = cli.StringOpt{
		Name:      "mainnet-weth9",
		Desc:      "Specify address of 0x WETH9 contract on MainNet network",
		EnvVar:    "DEXTERM_MAINNET_WETH9",
		Value:     "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		SetByUser: &networksMainnetWETH9Set,
	}
)

var (
	networksMainnetERC20ProxySet bool
	networksMainnetERC20ProxyOpt = cli.StringOpt{
		Name:      "mainnet-erc20proxy",
		Desc:      "Specify address of 0x ERC20Proxy contract on MainNet network",
		EnvVar:    "DEXTERM_MAINNET_ERC20PROXY",
		Value:     "0x95e6f48254609a6ee006f7d493c8e5fb97094cef",
		SetByUser: &networksMainnetERC20ProxySet,
	}
)

var (
	networksRopstenEndpointSet bool
	networksRopstenEndpointOpt = cli.StringOpt{
		Name:      "ropsten-endpoint",
		Desc:      "Specify endpoint for Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_ENDPOINT",
		Value:     "https://eth-ropsten.alchemyapi.io/jsonrpc/VNWdfCy8uyhRsbSveuAeyqU3VbPS55Sb",
		SetByUser: &networksRopstenEndpointSet,
	}
)

var (
	networksRopstenGasPriceSet bool
	networksRopstenGasPriceOpt = cli.StringOpt{
		Name:      "ropsten-gasprice",
		Desc:      "Specify min gasprice for Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_GASPRICE",
		Value:     "10000000000",
		SetByUser: &networksRopstenGasPriceSet,
	}
)

var (
	networksRopstenWETH9Set bool
	networksRopstenWETH9Opt = cli.StringOpt{
		Name:      "ropsten-weth9",
		Desc:      "Specify address of 0x WETH9 contract on Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_WETH9",
		Value:     "0xc778417e063141139fce010982780140aa0cd5ab",
		SetByUser: &networksRopstenWETH9Set,
	}
)

var (
	networksRopstenERC20ProxySet bool
	networksRopstenERC20ProxyOpt = cli.StringOpt{
		Name:      "ropsten-erc20proxy",
		Desc:      "Specify address of 0x ERC20Proxy contract on Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_ERC20PROXY",
		Value:     "0xb1408f4c245a23c31b98d2c626777d4c0d766caa",
		SetByUser: &networksRopstenERC20ProxySet,
	}
)

var (
	networksDevnetEndpointSet bool
	networksDevnetEndpointOpt = cli.StringOpt{
		Name:      "devnet-endpoint",
		Desc:      "Specify endpoint for Ganache network",
		EnvVar:    "DEXTERM_DEVNET_ENDPOINT",
		Value:     "devnet",
		SetByUser: &networksDevnetEndpointSet,
	}
)

var (
	networksDevnetGasPriceSet bool
	networksDevnetGasPriceOpt = cli.StringOpt{
		Name:      "devnet-gasprice",
		Desc:      "Specify min gasprice for Ganache network",
		EnvVar:    "DEXTERM_DEVNET_GASPRICE",
		Value:     "10000000000",
		SetByUser: &networksDevnetGasPriceSet,
	}
)

var (
	networksDevnetWETH9Set bool
	networksDevnetWETH9Opt = cli.StringOpt{
		Name:      "devnet-weth9",
		Desc:      "Specify address of 0x WETH9 contract on Ganache network",
		EnvVar:    "DEXTERM_DEVNET_WETH9",
		Value:     "0x0b1ba0af832d7c05fd64161e0db78e85978e8082",
		SetByUser: &networksDevnetWETH9Set,
	}
)

var (
	networksDevnetERC20ProxySet bool
	networksDevnetERC20ProxyOpt = cli.StringOpt{
		Name:      "devnet-erc20proxy",
		Desc:      "Specify address of 0x ERC20Proxy contract on Ganache network",
		EnvVar:    "DEXTERM_DEVNET_ERC20PROXY",
		Value:     "0x1dc4c1cefef38a777b15aa20260a54e584b16c48",
		SetByUser: &networksDevnetERC20ProxySet,
	}
)

var (
	networksMaticEndpointSet bool
	networksMaticEndpointOpt = cli.StringOpt{
		Name:      "matic-endpoint",
		Desc:      "Specify endpoint for Matic network",
		EnvVar:    "DEXTERM_MATIC_ENDPOINT",
		Value:     "matic",
		SetByUser: &networksMaticEndpointSet,
	}
)

var (
	networksMaticGasPriceSet bool
	networksMaticGasPriceOpt = cli.StringOpt{
		Name:      "matic-gasprice",
		Desc:      "Specify min gasprice for Matic network",
		EnvVar:    "DEXTERM_MATIC_GASPRICE",
		Value:     "5000000000",
		SetByUser: &networksMaticGasPriceSet,
	}
)

var (
	networksMaticWETH9Set bool
	networksMaticWETH9Opt = cli.StringOpt{
		Name:      "matic-weth9",
		Desc:      "Specify address of 0x WETH9 contract on Matic network",
		EnvVar:    "DEXTERM_MATIC_WETH9",
		Value:     "0x1d321b0bae75de3e4f5fb498e57d0276e73bfc0e",
		SetByUser: &networksMaticWETH9Set,
	}
)

var (
	networksMaticERC20ProxySet bool
	networksMaticERC20ProxyOpt = cli.StringOpt{
		Name:      "matic-erc20proxy",
		Desc:      "Specify address of 0x ERC20Proxy contract on Matic network",
		EnvVar:    "DEXTERM_MATIC_ERC20PROXY",
		Value:     "0x46fe55b51d24d269cccad63e2ab86f75751a39aa",
		SetByUser: &networksMaticERC20ProxySet,
	}
)