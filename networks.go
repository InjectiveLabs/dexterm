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
		Value:     "true",
		SetByUser: &networksAllowGasOraclesSet,
	}
)

var (
	networksDefaultSet bool
	networksDefaultOpt = cli.StringOpt{
		Name:      "N network",
		Desc:      "Specify network name to use for all transactions withing the session.",
		EnvVar:    "DEXTERM_NETWORK",
		Value:     "devnet",
		SetByUser: &networksDefaultSet,
	}
)

var (
	networksMainnetEndpointSet bool
	networksMainnetEndpointOpt = cli.StringOpt{
		Name:      "mainnet-endpoint",
		Desc:      "Specify endpoint for MainNet network",
		EnvVar:    "DEXTERM_MAINNET_ENDPOINT",
		Value:     "https://eth-mainnet.alchemyapi.io/v2/DqEv1TiHskO-G6JprqyhE25k1x0p3hpj",
		SetByUser: &networksMainnetEndpointSet,
	}
)

var (
	networksMainnetExplorerSet bool
	networksMainnetExplorerOpt = cli.StringOpt{
		Name:      "mainnet-explorer",
		Desc:      "Specify explorer prefix for transactions on MainNet network",
		EnvVar:    "DEXTERM_MAINNET_EXPLORER",
		Value:     "https://etherscan.io/tx/",
		SetByUser: &networksMainnetExplorerSet,
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
	networksMainnetExchangeSet bool
	networksMainnetExchangeOpt = cli.StringOpt{
		Name:      "mainnet-exchange",
		Desc:      "Specify address of 0x Exchange contract on MainNet network",
		EnvVar:    "DEXTERM_MAINNET_EXCHANGE",
		Value:     "",
		SetByUser: &networksMainnetExchangeSet,
	}
)

var (
	networksMainnetDevUtilsSet bool
	networksMainnetDevUtilsOpt = cli.StringOpt{
		Name:      "mainnet-devutils",
		Desc:      "Specify address of 0x Exchange DevUtils contract on Mainnet network",
		EnvVar:    "DEXTERM_MAINNET_DEVUTILS",
		Value:     "0x161793cdca4ff9e766a706c2c49c36ac1340bbcd",
		SetByUser: &networksMainnetDevUtilsSet,
	}
)

var (
	networksMainnetFuturesSet bool
	networksMainnetFuturesOpt = cli.StringOpt{
		Name:      "mainnet-futures",
		Desc:      "Specify address of Futures contract on MainNet network",
		EnvVar:    "DEXTERM_MAINNET_FUTURES",
		Value:     "",
		SetByUser: &networksMainnetFuturesSet,
	}
)

var (
	networksMainnetCoordinatorSet bool
	networksMainnetCoordinatorOpt = cli.StringOpt{
		Name:      "mainnet-coordinator",
		Desc:      "Specify address of Coordinator (Injective's Controller) contract on MainNet network",
		EnvVar:    "DEXTERM_MAINNET_COORDINATOR",
		Value:     "",
		SetByUser: &networksMainnetCoordinatorSet,
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
	networksRopstenExplorerSet bool
	networksRopstenExplorerOpt = cli.StringOpt{
		Name:      "ropsten-explorer",
		Desc:      "Specify explorer prefix for transactions on Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_EXPLORER",
		Value:     "https://ropsten.etherscan.io/tx/",
		SetByUser: &networksRopstenExplorerSet,
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
	networksRopstenExchangeSet bool
	networksRopstenExchangeOpt = cli.StringOpt{
		Name:      "ropsten-exchange",
		Desc:      "Specify address of Exchange (Injective's Controller) contract on Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_EXCHANGE",
		Value:     "",
		SetByUser: &networksRopstenExchangeSet,
	}
)

var (
	networksRopstenDevUtilsSet bool
	networksRopstenDevUtilsOpt = cli.StringOpt{
		Name:      "ropsten-devutils",
		Desc:      "Specify address of 0x Exchange DevUtils contract on Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_DEVUTILS",
		Value:     "0x161793cdca4ff9e766a706c2c49c36ac1340bbcd",
		SetByUser: &networksRopstenDevUtilsSet,
	}
)

var (
	networksRopstenFuturesSet bool
	networksRopstenFuturesOpt = cli.StringOpt{
		Name:      "ropsten-futures",
		Desc:      "Specify address of Futures contract on Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_FUTURES",
		Value:     "",
		SetByUser: &networksRopstenFuturesSet,
	}
)

var (
	networksRopstenCoordinatorSet bool
	networksRopstenCoordinatorOpt = cli.StringOpt{
		Name:      "ropsten-coordinator",
		Desc:      "Specify address of Coordinator (Injective's Controller) contract on Ropsten network",
		EnvVar:    "DEXTERM_ROPSTEN_COORDINATOR",
		Value:     "",
		SetByUser: &networksRopstenCoordinatorSet,
	}
)

var (
	networksKovanEndpointSet bool
	networksKovanEndpointOpt = cli.StringOpt{
		Name:      "kovan-endpoint",
		Desc:      "Specify endpoint for Kovan network",
		EnvVar:    "DEXTERM_KOVAN_ENDPOINT",
		Value:     "https://eth-kovan.alchemyapi.io/v2/9JhopJSP_O5NAgpUmbX7R09l39BLsZSh",
		SetByUser: &networksKovanEndpointSet,
	}
)

var (
	networksKovanExplorerSet bool
	networksKovanExplorerOpt = cli.StringOpt{
		Name:      "kovan-explorer",
		Desc:      "Specify explorer prefix for transactions on Kovan network",
		EnvVar:    "DEXTERM_KOVAN_EXPLORER",
		Value:     "https://kovan.etherscan.io/tx/",
		SetByUser: &networksKovanExplorerSet,
	}
)

var (
	networksKovanGasPriceSet bool
	networksKovanGasPriceOpt = cli.StringOpt{
		Name:      "kovan-gasprice",
		Desc:      "Specify min gasprice for Kovan network",
		EnvVar:    "DEXTERM_KOVAN_GASPRICE",
		Value:     "10000000000",
		SetByUser: &networksKovanGasPriceSet,
	}
)

var (
	networksKovanWETH9Set bool
	networksKovanWETH9Opt = cli.StringOpt{
		Name:      "kovan-weth9",
		Desc:      "Specify address of 0x WETH9 contract on Kovan network",
		EnvVar:    "DEXTERM_KOVAN_WETH9",
		Value:     "0xd0a1e359811322d97991e03f863a0c30c2cf029c",
		SetByUser: &networksKovanWETH9Set,
	}
)

var (
	networksKovanERC20ProxySet bool
	networksKovanERC20ProxyOpt = cli.StringOpt{
		Name:      "kovan-erc20proxy",
		Desc:      "Specify address of 0x ERC20Proxy contract on Kovan network",
		EnvVar:    "DEXTERM_KOVAN_ERC20PROXY",
		Value:     "0xf1ec01d6236d3cd881a0bf0130ea25fe4234003e",
		SetByUser: &networksKovanERC20ProxySet,
	}
)

var (
	networksKovanExchangeSet bool
	networksKovanExchangeOpt = cli.StringOpt{
		Name:      "kovan-exchange",
		Desc:      "Specify address of Exchange (Injective's Controller) contract on Kovan network",
		EnvVar:    "DEXTERM_KOVAN_EXCHANGE",
		Value:     "0x4eacd0af335451709e1e7b570b8ea68edec8bc97",
		SetByUser: &networksKovanExchangeSet,
	}
)

var (
	networksKovanFuturesSet bool
	networksKovanFuturesOpt = cli.StringOpt{
		Name:      "kovan-futures",
		Desc:      "Specify address of Futures contract on Kovan network",
		EnvVar:    "DEXTERM_KOVAN_FUTURES",
		Value:     "0x4eacd0af335451709e1e7b570b8ea68edec8bc97",
		SetByUser: &networksKovanFuturesSet,
	}
)

var (
	networksKovanDevUtilsSet bool
	networksKovanDevUtilsOpt = cli.StringOpt{
		Name:      "kovan-devutils",
		Desc:      "Specify address of 0x Exchange DevUtils contract on Kovan network",
		EnvVar:    "DEXTERM_KOVAN_DEVUTILS",
		Value:     "0x161793cdca4ff9e766a706c2c49c36ac1340bbcd",
		SetByUser: &networksKovanDevUtilsSet,
	}
)

var (
	networksKovanCoordinatorSet bool
	networksKovanCoordinatorOpt = cli.StringOpt{
		Name:      "kovan-coordinator",
		Desc:      "Specify address of Coordinator (Injective's Controller) contract on Kovan network",
		EnvVar:    "DEXTERM_KOVAN_COORDINATOR",
		Value:     "0x30493852999f5091d2430B6a1222Aa816237a486",
		SetByUser: &networksKovanCoordinatorSet,
	}
)

var (
	networksDevnetEndpointSet bool
	networksDevnetEndpointOpt = cli.StringOpt{
		Name:      "devnet-endpoint",
		Desc:      "Specify endpoint for Ganache network",
		EnvVar:    "DEXTERM_DEVNET_ENDPOINT",
		Value:     "http://localhost:8545",
		SetByUser: &networksDevnetEndpointSet,
	}
)

var (
	networksDevnetExplorerSet bool
	networksDevnetExplorerOpt = cli.StringOpt{
		Name:      "devnet-explorer",
		Desc:      "Specify explorer prefix for transactions on Ganache network",
		EnvVar:    "DEXTERM_DEVNET_EXPLORER",
		Value:     "",
		SetByUser: &networksDevnetExplorerSet,
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
	networksDevnetExchangeSet bool
	networksDevnetExchangeOpt = cli.StringOpt{
		Name:      "devnet-exchange",
		Desc:      "Specify address of Exchange (Injective's Controller) contract on Ganache network",
		EnvVar:    "DEXTERM_DEVNET_EXCHANGE",
		Value:     "0x48bacb9266a570d521063ef5dd96e61686dbe788",
		SetByUser: &networksDevnetExchangeSet,
	}
)

var (
	networksDevnetDevUtilsSet bool
	networksDevnetDevUtilsOpt = cli.StringOpt{
		Name:      "devnet-devutils",
		Desc:      "Specify address of 0x Exchange DevUtils contract on Ganache network",
		EnvVar:    "DEXTERM_DEVNET_DEVUTILS",
		Value:     "0xa31e64ea55b9b6bbb9d6a676738e9a5b23149f84",
		SetByUser: &networksDevnetDevUtilsSet,
	}
)

var (
	networksDevnetFuturesSet bool
	networksDevnetFuturesOpt = cli.StringOpt{
		Name:      "devnet-futures",
		Desc:      "Specify address of Futures contract on Ganache network",
		EnvVar:    "DEXTERM_DEVNET_FUTURES",
		Value:     "0xee233a2af6c7761443a951c0e942dff1f1686907",
		SetByUser: &networksDevnetFuturesSet,
	}
)

var (
	networksDevnetCoordinatorSet bool
	networksDevnetCoordinatorOpt = cli.StringOpt{
		Name:      "devnet-coordinator",
		Desc:      "Specify address of Coordinator (Injective's Controller) contract on Ganache network",
		EnvVar:    "DEXTERM_DEVNET_COORDINATOR",
		Value:     "0x4ef5b1E3dA5573466Fb1724D2Fca95290119B664",
		SetByUser: &networksDevnetCoordinatorSet,
	}
)

var (
	networksInjectiveEndpointSet bool
	networksInjectiveEndpointOpt = cli.StringOpt{
		Name:      "injective-endpoint",
		Desc:      "Specify endpoint for Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_ENDPOINT",
		Value:     "https://evm-us.injective.dev",
		SetByUser: &networksInjectiveEndpointSet,
	}
)

var (
	networksInjectiveExplorerSet bool
	networksInjectiveExplorerOpt = cli.StringOpt{
		Name:      "injective-explorer",
		Desc:      "Specify explorer prefix for transactions on Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_EXPLORER",
		Value:     "https://evm-explorer.injective.dev/tx/",
		SetByUser: &networksInjectiveExplorerSet,
	}
)

var (
	networksInjectiveGasPriceSet bool
	networksInjectiveGasPriceOpt = cli.StringOpt{
		Name:      "injective-gasprice",
		Desc:      "Specify min gasprice for Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_GASPRICE",
		Value:     "0",
		SetByUser: &networksInjectiveGasPriceSet,
	}
)

var (
	networksInjectiveWETH9Set bool
	networksInjectiveWETH9Opt = cli.StringOpt{
		Name:      "injective-weth9",
		Desc:      "Specify address of 0x WETH9 contract on Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_WETH9",
		Value:     "0x8a01cf83cb98dbfdab427ace0d276e8a77f5747b",
		SetByUser: &networksInjectiveWETH9Set,
	}
)

var (
	networksInjectiveERC20ProxySet bool
	networksInjectiveERC20ProxyOpt = cli.StringOpt{
		Name:      "injective-erc20proxy",
		Desc:      "Specify address of 0x ERC20Proxy contract on Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_ERC20PROXY",
		Value:     "0xdcf41118095d7bc0901d0215dae114a4cb2ad553", // TODO: THIS IS currently the futures contract address but this should be made separate in the future.
		SetByUser: &networksInjectiveERC20ProxySet,
	}
)

var (
	networksInjectiveExchangeSet bool
	networksInjectiveExchangeOpt = cli.StringOpt{
		Name:      "injective-exchange",
		Desc:      "Specify address of Exchange (Injective's Controller) contract on Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_EXCHANGE",
		Value:     "0x7f3e4297457740648b8262c4fc4bb3c67ff5d95c",
		SetByUser: &networksInjectiveExchangeSet,
	}
)

var (
	networksInjectiveDevUtilsSet bool
	networksInjectiveDevUtilsOpt = cli.StringOpt{
		Name:      "injective-devutils",
		Desc:      "Specify address of 0x Exchange DevUtils contract on Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_DEVUTILS",
		Value:     "0x988e8D73aCd4F9aB84cF4b0D6C264A6f0F275807",
		SetByUser: &networksInjectiveDevUtilsSet,
	}
)

var (
	networksInjectiveFuturesSet bool
	networksInjectiveFuturesOpt = cli.StringOpt{
		Name:      "injective-futures",
		Desc:      "Specify address of Futures contract on Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_FUTURES",
		Value:     "0x8f399baf9009a1466d9a3d8372703427c9f0c8cc",
		SetByUser: &networksInjectiveFuturesSet,
	}
)

var (
	networksInjectiveCoordinatorSet bool
	networksInjectiveCoordinatorOpt = cli.StringOpt{
		Name:      "injective-coordinator",
		Desc:      "Specify address of Coordinator (Injective's Controller) contract on Injective network",
		EnvVar:    "DEXTERM_INJECTIVE_COORDINATOR",
		Value:     "0x3b46ef40b11888b7353c764fca86a83ff89dc90c",
		SetByUser: &networksInjectiveCoordinatorSet,
	}
)

var (
	networksMaticEndpointSet bool
	networksMaticEndpointOpt = cli.StringOpt{
		Name:      "matic-endpoint",
		Desc:      "Specify endpoint for Matic network",
		EnvVar:    "DEXTERM_MATIC_ENDPOINT",
		Value:     "https://evm-rpc.injective.dev",
		SetByUser: &networksMaticEndpointSet,
	}
)

var (
	networksMaticExplorerSet bool
	networksMaticExplorerOpt = cli.StringOpt{
		Name:      "matic-explorer",
		Desc:      "Specify explorer prefix for transactions on Matic network",
		EnvVar:    "DEXTERM_MATIC_EXPLORER",
		Value:     "https://testnetv3-explorer.matic.network/tx/",
		SetByUser: &networksMaticExplorerSet,
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
		Value:     "0x5d8df0cc28e77cf206c901d7614aa4a98ef4925f",
		SetByUser: &networksMaticERC20ProxySet,
	}
)

var (
	networksMaticExchangeSet bool
	networksMaticExchangeOpt = cli.StringOpt{
		Name:      "matic-exchange",
		Desc:      "Specify address of Exchange (Injective's Controller) contract on Matic network",
		EnvVar:    "DEXTERM_MATIC_EXCHANGE",
		Value:     "0x50C655DD81B65D6B48D759F897881BD5ADd86E57",
		SetByUser: &networksMaticExchangeSet,
	}
)

var (
	networksMaticFuturesSet bool
	networksMaticFuturesOpt = cli.StringOpt{
		Name:      "matic-futures",
		Desc:      "Specify address of Futures contract on Matic network",
		EnvVar:    "DEXTERM_MATIC_FUTURES",
		Value:     "0x5d8df0cc28e77cf206c901d7614aa4a98ef4925f",
		SetByUser: &networksMaticFuturesSet,
	}
)

var (
	networksMaticCoordinatorSet bool
	networksMaticCoordinatorOpt = cli.StringOpt{
		Name:      "matic-coordinator",
		Desc:      "Specify address of Coordinator (Injective's Controller) contract on Matic network",
		EnvVar:    "DEXTERM_MATIC_COORDINATOR",
		Value:     "",
		SetByUser: &networksMaticCoordinatorSet,
	}
)
