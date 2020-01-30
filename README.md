## DEXTerm

<img src="https://user-images.githubusercontent.com/477998/73288585-50271200-420c-11ea-8813-ebbbd6463ede.png" width="500px" />

Dexterm is the trading client for Injective Protocol's decentralized exchange built on [0x framework](http://0x.org). It allows users to interact with our protocol and decentralized relay network in several ways such as managing token balances and submitting orders.

**Version 1.0** is a preview release and expects to communicate with localhost by default, although this is a fully-featured client and can be configured to talk to staging and prod servers as well.

After initial launch, you can find its default config and keystore in `~/.dexterm`.

## Supported platforms

* MacOS 64-bit
* Windows 64-bit
* Linux 64-bit

## Features

* Supports persistent configuration
* Supports ENV variables and CLI flags for temporary config overrides
* Supports permanent encrypted keystore
* All commands have autocompletion and useful tips based on context

### Keystore

* Create an Ethereum wallet
* Import an Ethereum wallet
* Generate a new Ethereum wallet with private key
* List accounts
* Switch between default accounts

### Utils

* List all tokens and balances, also their unlock status
* Lock and Unlock tokens for trading within 0x
* Wrap ETH into WETH and Unwrap ETH from WETH, just inside the app

### Trading

* List all available markets
* View orderbook of a market (ask, bid orders, notes)
* Sign and post sell (ask) order
* Sign and post buy (bid) order
* Fill any order from the orderbook for variable amount

## License

[BSD 3-clause](/LICENSE)

## Usage Examples

First, you should have `relayerd` and `relayer-api` (see [`injective-core`](http://github.com/InjectiveLabs/injective-core)) running as well as an instance of Ganache if you are running the network locally. 

Then run `dexterm`. 
