package main

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/big"
	"os"
	"sort"
	"strings"
	"time"

	zeroex "github.com/InjectiveLabs/zeroex-go"
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

	"github.com/InjectiveLabs/dexterm/clients"
	"github.com/InjectiveLabs/dexterm/ethereum/ethcore"
	"github.com/InjectiveLabs/dexterm/ethereum/ethfw/keystore"
	"github.com/InjectiveLabs/dexterm/ethereum/ethfw/manager"
	sraAPI "github.com/InjectiveLabs/injective-core/api/gen/relayer_api"
	restAPI "github.com/InjectiveLabs/injective-core/api/gen/rest_api"
)

func init() {
	closer.Bind(func() {
		logrus.Println("Bye!")
	})
}

type AppController struct {
	cfg        *toml.Tree
	configPath string

	debugClient       *clients.DebugClient
	restClient        *clients.RESTClient
	sraClient         *clients.SRAClient
	sdaClient         *clients.SDAClient
	coordinatorClient *clients.CoordinatorClient

	ethGasPrice         *big.Int
	ethCore             *ethcore.EthClient
	feeRecipientAddress common.Address

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

	if debugClient, err := clients.NewDebugClient(&clients.DebugClientConfig{
		Endpoint: ctl.mustConfigValue("relayer.endpoint"),
	}); err != nil {
		logrus.WithError(err).Warningln("no debug API HTTP connection")
	} else {
		ctl.debugClient = debugClient
	}

	if restClient, err := clients.NewRESTClient(&clients.RESTClientConfig{
		Endpoint: ctl.mustConfigValue("relayer.endpoint"),
	}); err != nil {
		logrus.WithError(err).Warningln("no REST HTTP connection, running in offline mode")
	} else {
		ctl.restClient = restClient

		sraEndpoint := ctl.mustConfigValue("relayer.endpoint")

		if sraClient, err := clients.NewSRAClient(restClient, &clients.SRAClientConfig{
			Endpoint: sraEndpoint,
		}); err != nil {
			logrus.WithError(err).Warningln("no SRA HTTP connection")
		} else {
			ctl.sraClient = sraClient
		}

		sdaEndpoint := ctl.mustConfigValue("relayer.endpoint")

		if sdaClient, err := clients.NewSDAClient(restClient, &clients.SDAClientConfig{
			Endpoint: sdaEndpoint,
		}); err != nil {
			logrus.WithError(err).Warningln("no SDA HTTP connection")
		} else {
			ctl.sdaClient = sdaClient
		}

		coordinatorEndpoint := ctl.mustConfigValue("relayer.endpoint")

		if coordinatorClient, err := clients.NewCoordinatorClient(&clients.CoordinatorClientConfig{
			Endpoint: coordinatorEndpoint,
		}); err != nil {
			logrus.WithError(err).Warningln("no coordinator HTTP connection")
		} else {
			ctl.coordinatorClient = coordinatorClient
		}

		ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFn()

		feeRecipients, err := ctl.sraClient.FeeRecipients(ctx)
		if err != nil {
			logrus.Fatalln(err)
		} else if len(feeRecipients) == 0 {
			logrus.WithFields(logrus.Fields{
				"endpoint": sraEndpoint,
			}).Fatalln("no fee recipients fetched from SRA endpoint")
		}

		ctl.feeRecipientAddress = feeRecipients[0]
		logrus.WithField("address", ctl.feeRecipientAddress.Hex()).Println("SRA endpoint provided by staker")
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
	} else if _, ok := ctl.getConfigValue("accounts.default"); !ok {
		ctl.generateDefaultAccount()
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

type TradeLimitBuyOrderArgs struct {
	Market       string
	Amount       string
	Price        string
	SignPassword string
}

type TradeDerivativeLimitOrderArgs struct {
	Market       string
	Quantity     string
	Price        string
	SignPassword string
}

// keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}

func (ctl *AppController) ActionDerivativesLimitLong(args interface{}) {
	makeDerivativeOrderArgs := args.(*TradeDerivativeLimitOrderArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	markets, err := ctl.restClient.DerivativeMarkets(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var makerAssetData []byte
	var takerAssetData []byte

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))
	for _, market := range markets {
		if market.Ticker != makeDerivativeOrderArgs.Market {
			continue
		}
		// TODO: need to call getAccounts to get list of accountIDs to allow trader to select account to trade from earlier
		nonce := big.NewInt(1)
		makerAssetData = keccak256(defaultAccount.Bytes(), nonce.Bytes())
		takerAssetData = common.FromHex(market.MarketID)
	}

	if len(takerAssetData) == 0 {
		logrus.WithFields(logrus.Fields{
			"market": makeDerivativeOrderArgs.Market,
		}).Errorln("specified market not found")
		return
	}

	var takerAssetAmount *big.Int
	var makerAssetAmount *big.Int
	quantity, err := decimal.NewFromString(makeDerivativeOrderArgs.Quantity)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse buy amount")
		return
	} else if quantity.LessThan(decimal.RequireFromString("1.0")) {
		logrus.Errorln("Buy amount is too small, must be at least 1")
		return
	} else {
		takerAssetAmount, _ = big.NewInt(0).SetString(quantity.String(), 10)
		//takerAssetAmount = dec2big(quantity)
	}
	price, err := decimal.NewFromString(makeDerivativeOrderArgs.Price)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse buy price")
		return
	}
	makerAssetAmount = dec2big(price)

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: makeDerivativeOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	signedOrder, err := ctl.ethCore.CreateAndSignDerivativesOrder(
		callArgs,
		makerAssetData,
		takerAssetData,
		makerAssetAmount,
		takerAssetAmount,
		true,
	)
	if err != nil {
		logrus.WithError(err).Errorln("unable to sign order")
		return
	}

	orderHash, err := ctl.sdaClient.PostOrder(ctx, signedOrder)
	if err != nil {
		logrus.WithError(err).Errorln("unable to post order")
		return
	}

	fmt.Println(orderHash)
}

func (ctl *AppController) ActionDerivativesLimitShort(args interface{}) {
	makeDerivativeOrderArgs := args.(*TradeDerivativeLimitOrderArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	markets, err := ctl.restClient.DerivativeMarkets(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var makerAssetData []byte
	var takerAssetData []byte

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))
	for _, market := range markets {
		if market.Ticker != makeDerivativeOrderArgs.Market {
			continue
		}
		// TODO: need to call getAccounts to get list of accountIDs to allow trader to select account to trade from earlier
		nonce := big.NewInt(1)
		makerAssetData = keccak256(defaultAccount.Bytes(), nonce.Bytes())
		takerAssetData = common.FromHex(market.MarketID)
	}

	if len(takerAssetData) == 0 {
		logrus.WithFields(logrus.Fields{
			"market": makeDerivativeOrderArgs.Market,
		}).Errorln("specified market not found")
		return
	}

	var takerAssetAmount *big.Int
	var makerAssetAmount *big.Int
	quantity, err := decimal.NewFromString(makeDerivativeOrderArgs.Quantity)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse buy amount")
		return
	} else if quantity.LessThan(decimal.RequireFromString("1.0")) {
		logrus.Errorln("Buy amount is too small, must be at least 1")
		return
	} else {
		takerAssetAmount, _ = big.NewInt(0).SetString(quantity.String(), 10)
		//takerAssetAmount = dec2big(quantity)
	}
	price, err := decimal.NewFromString(makeDerivativeOrderArgs.Price)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse buy price")
		return
	}
	makerAssetAmount = dec2big(price)

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: makeDerivativeOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	signedOrder, err := ctl.ethCore.CreateAndSignDerivativesOrder(
		callArgs,
		makerAssetData,
		takerAssetData,
		makerAssetAmount,
		takerAssetAmount,
		false,
	)
	if err != nil {
		logrus.WithError(err).Errorln("unable to sign order")
		return
	}

	orderHash, err := ctl.sdaClient.PostOrder(ctx, signedOrder)
	if err != nil {
		logrus.WithError(err).Errorln("unable to post order")
		return
	}

	fmt.Println(orderHash)
}

func (ctl *AppController) ActionTradeLimitBuy(args interface{}) {
	makeBuyOrderArgs := args.(*TradeLimitBuyOrderArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tradePairs, err := ctl.restClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var makerAssetData []byte
	var takerAssetData []byte

	for _, pair := range tradePairs {
		if pair.Name != makeBuyOrderArgs.Market {
			continue
		}

		// swapped because it's a bid
		makerAssetData = common.FromHex(pair.TakerAssetData)
		takerAssetData = common.FromHex(pair.MakerAssetData)
	}

	if len(makerAssetData) == 0 || len(takerAssetData) == 0 {
		logrus.WithFields(logrus.Fields{
			"market": makeBuyOrderArgs.Market,
		}).Errorln("specified trade pair not found")

		return
	}

	var takerAmount *big.Int
	var price decimal.Decimal
	var makerAmount *big.Int

	takerAmountDec, err := decimal.NewFromString(makeBuyOrderArgs.Amount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse buy amount")
		return
	} else if takerAmountDec.LessThan(decimal.RequireFromString("0.0000001")) {
		logrus.Errorln("Buy amount is too small, must be at least 0.0000001")
		return
	} else {
		takerAmount = dec2big(takerAmountDec)
	}

	if price, err = decimal.NewFromString(makeBuyOrderArgs.Price); err != nil {
		logrus.WithError(err).Errorln("failed to parse buy price")
		return
	}
	makerAmount = dec2big(takerAmountDec.Mul(price))

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: makeBuyOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	signedOrder, err := ctl.ethCore.CreateAndSignOrder(
		callArgs,
		ctl.feeRecipientAddress,
		makerAssetData,
		takerAssetData,
		makerAmount,
		takerAmount,
	)
	if err != nil {
		logrus.WithError(err).Errorln("unable to sign order")
		return
	}

	orderHash, err := ctl.sraClient.PostOrder(ctx, signedOrder)
	if err != nil {
		logrus.WithError(err).Errorln("unable to post order")
		return
	}

	fmt.Println(orderHash)
}

type TradeLimitSellOrderArgs struct {
	Market       string
	Amount       string
	Price        string
	SignPassword string
}

func (ctl *AppController) ActionTradeLimitSell(args interface{}) {
	makeSellOrderArgs := args.(*TradeLimitSellOrderArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tradePairs, err := ctl.restClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var makerAssetData []byte
	var takerAssetData []byte

	for _, pair := range tradePairs {
		if pair.Name != makeSellOrderArgs.Market {
			continue
		}

		makerAssetData = common.FromHex(pair.MakerAssetData)
		takerAssetData = common.FromHex(pair.TakerAssetData)
	}

	if len(makerAssetData) == 0 || len(takerAssetData) == 0 {
		logrus.WithFields(logrus.Fields{
			"market": makeSellOrderArgs.Market,
		}).Errorln("specified trade pair not found")

		return
	}

	var makerAmount *big.Int
	var price decimal.Decimal
	var takerAmount *big.Int

	makerAmountDec, err := decimal.NewFromString(makeSellOrderArgs.Amount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse sell amount")
		return
	} else if makerAmountDec.LessThan(decimal.RequireFromString("0.0000001")) {
		logrus.Errorln("Sell amount is too small, must be at least 0.0000001")
		return
	} else {
		makerAmount = dec2big(makerAmountDec)
	}

	if price, err = decimal.NewFromString(makeSellOrderArgs.Price); err != nil {
		logrus.WithError(err).Errorln("failed to parse sell price")
		return
	}
	takerAmount = dec2big(makerAmountDec.Mul(price))

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: makeSellOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	signedOrder, err := ctl.ethCore.CreateAndSignOrder(
		callArgs,
		ctl.feeRecipientAddress,
		makerAssetData,
		takerAssetData,
		makerAmount,
		takerAmount,
	)
	if err != nil {
		logrus.WithError(err).Errorln("unable to sign order")
		return
	}

	orderHash, err := ctl.sraClient.PostOrder(ctx, signedOrder)
	if err != nil {
		logrus.WithError(err).Errorln("unable to post order")
		return
	}

	fmt.Println(orderHash)
}

type TradeMarketBuyOrderArgs struct {
	Market       string
	Amount       string
	SignPassword string
}

type TradeMarketSellOrderArgs struct {
	Market       string
	Amount       string
	SignPassword string
}

func dec2big(d decimal.Decimal) *big.Int {
	v, _ := big.NewInt(0).SetString(d.Truncate(9).Shift(18).String(), 10)

	return v
}

type TradeFillOrderArgs struct {
	Market       string
	OrderHash    string
	FillAmount   string
	SignPassword string
}

func (ctl *AppController) ActionTradeFillOrder(args interface{}) {
	fillOrderArgs := args.(*TradeFillOrderArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tradePairs, err := ctl.restClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var tradePair *restAPI.TradePair
	for _, pair := range tradePairs {
		if pair.Name != fillOrderArgs.Market {
			continue
		} else if !pair.Enabled {
			continue
		}

		tradePair = pair
	}

	if tradePair == nil {
		logrus.WithFields(logrus.Fields{
			"market": fillOrderArgs.Market,
		}).Errorln("specified trade pair not found or is not enabled")

		return
	}

	makeOrder, err := ctl.sraClient.Order(ctx, fillOrderArgs.OrderHash)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"order": fillOrderArgs.OrderHash,
		}).Errorln(err)
		return
	}

	var isBid bool
	if makeOrder.TakerAssetData == tradePair.MakerAssetData {
		isBid = true
	}

	var fillAmount *big.Int
	var price decimal.Decimal

	fillAmountDec, err := decimal.NewFromString(fillOrderArgs.FillAmount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse fill amount")
		return
	} else if fillAmountDec.LessThan(decimal.RequireFromString("0.0000001")) {
		logrus.Errorln("fill amount is too small, must be at least 0.0000001")
		return
	} else {
		fillAmount = dec2big(fillAmountDec)
	}

	price, vol := calcOrderPrice(makeOrder, isBid)
	vol = vol.Shift(-18)

	if fillAmountDec.GreaterThan(vol) {
		err = fmt.Errorf("wrong fill amount: %s", fillAmountDec.StringFixed(9))
		logrus.WithError(err).Errorf("maximum fill amount: %s", vol.StringFixed(9))
		return
	}

	var makerAmount *big.Int
	var takerAmount *big.Int
	var makerAssetData []byte
	var takerAssetData []byte

	if isBid {
		// maker must have fillAmount / price of quote currency
		// taker must have fillAmount of base currency
		makerAmount = dec2big(fillAmountDec.Div(price))
		takerAmount = fillAmount

		makerAssetData = common.FromHex(tradePair.TakerAssetData)
		takerAssetData = common.FromHex(tradePair.MakerAssetData)
	} else {
		// maker must have fillAmount of base currency
		// taker must have fillAmount * price of quote currency
		makerAmount = fillAmount
		takerAmount = dec2big(fillAmountDec.Mul(price))

		makerAssetData = common.FromHex(tradePair.MakerAssetData)
		takerAssetData = common.FromHex(tradePair.TakerAssetData)
	}

	// not used for now
	_ = makerAmount
	_ = makerAssetData
	_ = takerAssetData

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: fillOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	exchangeAddress := ctl.ethCore.ContractAddress(ethcore.EthContractExchange)

	zeroExOrder, err := ro2zo(makeOrder)
	if err != nil {
		logrus.WithError(err).Errorln("failed to convert SRAv3 order into zeroex.SignedOrder")
		return
	}
	signedTx, err := ctl.ethCore.CreateAndSignTransaction_FillOrders(
		callArgs,
		exchangeAddress,
		[]*zeroex.SignedOrder{zeroExOrder},
		[]*big.Int{takerAmount},
	)
	if err != nil {
		logrus.WithError(err).Errorln("unable to create and sign transaction")
		return
	}

	approvals, expiryAt, err := ctl.coordinatorClient.GetCoordinatorApproval(ctx, signedTx, defaultAccount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get approval from Coordinator API")
		return
	} else if time.Now().After(expiryAt) {
		logrus.WithError(err).Errorln("issued approval from Coordinator API already expired")
		return
	}

	txHash, err := ctl.ethCore.ExecuteTransaction(callArgs, signedTx, approvals[0])
	if err != nil {
		logrus.WithError(err).Errorln("unable to execute Exchange transaction")
		return
	}

	fmt.Println(ctl.formatTxLink(txHash))
	ctl.checkTx(txHash)
}

type TradeCancelOrderArgs struct {
	Market       string
	OrderHash    string
	SignPassword string
}

func (ctl *AppController) ActionTradeCancelOrder(args interface{}) {
	cancelOrderArgs := args.(*TradeCancelOrderArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tradePairs, err := ctl.restClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var tradePair *restAPI.TradePair
	for _, pair := range tradePairs {
		if pair.Name != cancelOrderArgs.Market {
			continue
		} else if !pair.Enabled {
			continue
		}

		tradePair = pair
	}

	if tradePair == nil {
		logrus.WithFields(logrus.Fields{
			"market": cancelOrderArgs.Market,
		}).Errorln("specified trade pair not found or is not enabled")

		return
	}

	makeOrder, err := ctl.sraClient.Order(ctx, cancelOrderArgs.OrderHash)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"order": cancelOrderArgs.OrderHash,
		}).Errorln(err)
		return
	}

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: cancelOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	exchangeAddress := ctl.ethCore.ContractAddress(ethcore.EthContractExchange)

	zeroExOrder, err := ro2zo(makeOrder)
	if err != nil {
		logrus.WithError(err).Errorln("failed to convert SRAv3 order into zeroex.SignedOrder")
		return
	}
	signedTx, err := ctl.ethCore.CreateAndSignTransaction_BatchCancelOrders(
		callArgs,
		exchangeAddress,
		[]*zeroex.SignedOrder{zeroExOrder},
	)
	if err != nil {
		logrus.WithError(err).Errorln("unable to create and sign transaction")
		return
	}

	cancellationSignatures, err := ctl.coordinatorClient.SendCoordinatorSoftCancelTransaction(ctx, signedTx, defaultAccount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get approval from Coordinator API")
		return
	}
	fmt.Println(cancellationSignatures)
	// No need to execute soft-cancel. Actually executing the soft-cancel on-chain would make it a hard cancel
	//txHash, err := ctl.ethCore.ExecuteTransaction(callArgs, signedTx, approvals[0])
	//if err != nil {
	//	logrus.WithError(err).Errorln("unable to execute Exchange transaction")
	//	return
	//}
	//
	//fmt.Println(ctl.formatTxLink(txHash))
	//ctl.checkTx(txHash)
}

func (ctl *AppController) ActionTradeMarketBuy(args interface{}) {
	marketOrderArgs := args.(*TradeMarketBuyOrderArgs)
	amountToBuy, _ := decimal.NewFromString(marketOrderArgs.Amount)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tradePairs, err := ctl.restClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var tradePair *restAPI.TradePair
	for _, pair := range tradePairs {
		if pair.Name != marketOrderArgs.Market {
			continue
		} else if !pair.Enabled {
			continue
		}

		tradePair = pair
		// swapped because it's a bid
	}

	if tradePair == nil {
		logrus.WithFields(logrus.Fields{
			"market": marketOrderArgs.Market,
		}).Errorln("specified trade pair not found or is not enabled")

		return
	}

	// get orderbook
	_, asks, err := ctl.sraClient.Orderbook(ctx, tradePair.Name)
	if err != nil {
		logrus.WithField("tradePair", tradePair.Name).
			WithError(err).Errorln("unable to get orderbook for trade pair")
		return
	}

	// sort asks
	sort.Slice(asks, func(i, j int) bool {
		a, _ := decimal.NewFromString(asks[i].Order.MakerAssetAmount)
		b, _ := decimal.NewFromString(asks[i].Order.TakerAssetAmount)
		c, _ := decimal.NewFromString(asks[j].Order.MakerAssetAmount)
		d, _ := decimal.NewFromString(asks[j].Order.TakerAssetAmount)
		// TODO: check if works, may be opposite
		return (a.Mul(d)).GreaterThan(b.Mul(c))
	})
	ordersToFill := []*zeroex.SignedOrder{}
	amountToFill := decimal.NewFromInt(0)
	for _, order := range asks {
		//makerAmount, _ := decimal.NewFromString(order.Order.MakerAssetAmount)
		takerAmount, _ := decimal.NewFromString(order.Order.TakerAssetAmount)
		sum := takerAmount.Add(amountToFill)
		if sum.GreaterThanOrEqual(amountToBuy) {
			zeroExOrder, err := ro2zo(order.Order)
			if err != nil {
				logrus.WithError(err).Errorln("failed to convert SRAv3 order into zeroex.SignedOrder")
				return
			}
			ordersToFill = append(ordersToFill, zeroExOrder)
			break
		}
	}

	// fill orders

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: marketOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	exchangeAddress := ctl.ethCore.ContractAddress(ethcore.EthContractExchange)

	signedTx, err := ctl.ethCore.CreateAndSignTransaction_MarketBuyOrders(
		callArgs,
		exchangeAddress,
		ordersToFill,
		dec2big(amountToBuy),
	)

	if err != nil {
		logrus.WithError(err).Errorln("unable to create and sign transaction")
		return
	}

	approvals, expiryAt, err := ctl.coordinatorClient.GetCoordinatorApproval(ctx, signedTx, defaultAccount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get approval from Coordinator API")
		return
	} else if time.Now().After(expiryAt) {
		logrus.WithError(err).Errorln("issued approval from Coordinator API already expired")
		return
	}

	txHash, err := ctl.ethCore.ExecuteTransaction(callArgs, signedTx, approvals[0])
	if err != nil {
		logrus.WithError(err).Errorln("unable to execute Exchange transaction")
		return
	}

	fmt.Println(ctl.formatTxLink(txHash))
	ctl.checkTx(txHash)
}

func (ctl *AppController) ActionTradeMarketSell(args interface{}) {
	marketOrderArgs := args.(*TradeMarketSellOrderArgs)
	amountToSell, _ := decimal.NewFromString(marketOrderArgs.Amount)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tradePairs, err := ctl.restClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var tradePair *restAPI.TradePair
	for _, pair := range tradePairs {
		if pair.Name != marketOrderArgs.Market {
			continue
		} else if !pair.Enabled {
			continue
		}

		tradePair = pair
		// swapped because it's a bid
	}

	if tradePair == nil {
		logrus.WithFields(logrus.Fields{
			"market": marketOrderArgs.Market,
		}).Errorln("specified trade pair not found or is not enabled")

		return
	}

	// get orderbook
	bids, _, err := ctl.sraClient.Orderbook(ctx, tradePair.Name)
	if err != nil {
		logrus.WithField("tradePair", tradePair.Name).
			WithError(err).Errorln("unable to get orderbook for trade pair")
		return
	}

	// sort asks
	sort.Slice(bids, func(i, j int) bool {
		a, _ := decimal.NewFromString(bids[i].Order.MakerAssetAmount)
		b, _ := decimal.NewFromString(bids[i].Order.TakerAssetAmount)
		c, _ := decimal.NewFromString(bids[j].Order.MakerAssetAmount)
		d, _ := decimal.NewFromString(bids[j].Order.TakerAssetAmount)
		// TODO: check if works, may be opposite
		return (a.Mul(d)).GreaterThan(b.Mul(c))
	})
	ordersToFill := []*zeroex.SignedOrder{}
	amountToFill := decimal.NewFromInt(0)
	for _, order := range bids {
		//makerAmount, _ := decimal.NewFromString(order.Order.MakerAssetAmount)
		takerAmount, _ := decimal.NewFromString(order.Order.TakerAssetAmount)
		sum := takerAmount.Add(amountToFill)
		if sum.GreaterThanOrEqual(amountToSell) {
			zeroExOrder, err := ro2zo(order.Order)
			if err != nil {
				logrus.WithError(err).Errorln("failed to convert SRAv3 order into zeroex.SignedOrder")
				return
			}
			ordersToFill = append(ordersToFill, zeroExOrder)
			break
		}
	}

	// fill orders

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: marketOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	exchangeAddress := ctl.ethCore.ContractAddress(ethcore.EthContractExchange)

	signedTx, err := ctl.ethCore.CreateAndSignTransaction_MarketSellOrders(
		callArgs,
		exchangeAddress,
		ordersToFill,
		dec2big(amountToSell),
	)

	if err != nil {
		logrus.WithError(err).Errorln("unable to create and sign transaction")
		return
	}

	approvals, expiryAt, err := ctl.coordinatorClient.GetCoordinatorApproval(ctx, signedTx, defaultAccount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get approval from Coordinator API")
		return
	} else if time.Now().After(expiryAt) {
		logrus.WithError(err).Errorln("issued approval from Coordinator API already expired")
		return
	}

	txHash, err := ctl.ethCore.ExecuteTransaction(callArgs, signedTx, approvals[0])
	if err != nil {
		logrus.WithError(err).Errorln("unable to execute Exchange transaction")
		return
	}

	fmt.Println(ctl.formatTxLink(txHash))
	ctl.checkTx(txHash)
}

type TradeGenerateLimitOrdersArgs struct {
	Market       string
	Amount       string
	Price        string
	SignPassword string
}

func (ctl *AppController) ActionTradeGenerateLimitOrders(args interface{}) {
	makeBuyOrderArgs := args.(*TradeGenerateLimitOrdersArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	tradePairs, err := ctl.restClient.TradePairs(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("unable to fetch trade pairs")
		return
	}

	var makerAssetData []byte
	var takerAssetData []byte

	for _, pair := range tradePairs {
		if pair.Name != makeBuyOrderArgs.Market {
			continue
		}

		// swapped because it's a bid
		makerAssetData = common.FromHex(pair.TakerAssetData)
		takerAssetData = common.FromHex(pair.MakerAssetData)
	}

	if len(makerAssetData) == 0 || len(takerAssetData) == 0 {
		logrus.WithFields(logrus.Fields{
			"market": makeBuyOrderArgs.Market,
		}).Errorln("specified trade pair not found")

		return
	}

	var takerAmount *big.Int
	var price decimal.Decimal
	var makerAmount *big.Int

	takerAmountDec, err := decimal.NewFromString(makeBuyOrderArgs.Amount)
	if err != nil {
		logrus.WithError(err).Errorln("failed to parse buy amount")
		return
	} else if takerAmountDec.LessThan(decimal.RequireFromString("0.0000001")) {
		logrus.Errorln("Buy amount is too small, must be at least 0.0000001")
		return
	} else {
		takerAmount = dec2big(takerAmountDec)
	}

	if price, err = decimal.NewFromString(makeBuyOrderArgs.Price); err != nil {
		logrus.WithError(err).Errorln("failed to parse buy price")
		return
	}
	makerAmount = dec2big(takerAmountDec.Mul(price))

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
		Context:  ctx,
		From:     defaultAccount,
		FromPass: makeBuyOrderArgs.SignPassword,
		GasPrice: ctl.ethGasPrice,
	}

	// 5 make orders
	var i int64
	for i = 2; i < 7; i++ {
		scale := decimal.NewFromInt(i)
		newTakerAmountDec := takerAmountDec.Mul(scale)
		newTakerAmount := dec2big(newTakerAmountDec)
		signedOrder, err := ctl.ethCore.CreateAndSignOrder(
			callArgs,
			ctl.feeRecipientAddress,
			makerAssetData,
			takerAssetData,
			makerAmount,
			newTakerAmount,
		)
		if err != nil {
			logrus.WithError(err).Errorln("unable to sign order")
			return
		}

		orderHash, err := ctl.sraClient.PostOrder(ctx, signedOrder)
		if err != nil {
			logrus.WithError(err).Errorln("unable to post order")
			return
		}

		fmt.Println(orderHash)
	}

	// 5 take orders
	for i = 2; i < 7; i++ {
		scale := decimal.NewFromInt(i)
		makerAmountDec := takerAmountDec.Mul(price)
		newMakerAmountDec := makerAmountDec.Mul(scale)
		newMakerAmount := dec2big(newMakerAmountDec)
		signedOrder, err := ctl.ethCore.CreateAndSignOrder(
			callArgs,
			ctl.feeRecipientAddress,
			takerAssetData,
			makerAssetData,
			newMakerAmount,
			takerAmount,
		)
		if err != nil {
			logrus.WithError(err).Errorln("unable to sign order")
			return
		}

		orderHash, err := ctl.sraClient.PostOrder(ctx, signedOrder)
		if err != nil {
			logrus.WithError(err).Errorln("unable to post order")
			return
		}

		fmt.Println(orderHash)
	}
}

type DerivativeOrderbookArgs struct {
	Market string
}

func (ctl *AppController) ActionDerivativesOrderbook(args interface{}) {
	derivativeOrderbookArgs := args.(*DerivativeOrderbookArgs)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))
	markets, err := ctl.restClient.DerivativeMarkets(ctx)
	var takerAssetDataString string
	for _, market := range markets {
		if market.Ticker == derivativeOrderbookArgs.Market {
			takerAssetDataString = market.MarketID
		}
	}


	bids, asks, err := ctl.sraClient.DerivativeOrders(ctx, takerAssetDataString)
	if err != nil {
		logrus.WithField("market", derivativeOrderbookArgs.Market).
			WithError(err).Errorln("unable to get orderbook for market")
		return
	}


	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddTitle("ORDERBOOK")
	table.AddHeaders(
		fmt.Sprintf("Price"),
		fmt.Sprintf("Contracts"),
		"Notes",
	)

	if len(asks) == 0 {
		table.AddRow(color.RedString("No asks."), "", "")
	} else {
		for _, ask := range asks {
			var notes string
			if isMakerOf(ask.Order, defaultAccount) {
				notes = "⭑ owner"
			}


			price := decimal.RequireFromString(ask.Order.MakerAssetAmount)
			quantity := decimal.RequireFromString(ask.Order.TakerAssetAmount)

			table.AddRow(
				color.RedString("%s", price.StringFixed(9)),
				color.RedString("%s", quantity),
				notes,
			)
		}
	}

	table.AddSeparator()

	if len(bids) == 0 {
		table.AddRow(color.GreenString("No bids."), "", "")
	} else {
		for _, bid := range bids {
			var notes string
			if isMakerOf(bid.Order, defaultAccount) {
				notes = "⭑ owner"
			}

			price := decimal.RequireFromString(bid.Order.MakerAssetAmount)
			quantity := decimal.RequireFromString(bid.Order.TakerAssetAmount)

			table.AddRow(
				color.GreenString("%s", price.StringFixed(9)),
				color.GreenString("%s", quantity),
				notes,
			)
		}
	}

	fmt.Println(table.Render())
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

	bids, asks, err := ctl.sraClient.Orderbook(ctx, orderbookArgs.Market)
	if err != nil {
		logrus.WithField("tradePair", orderbookArgs.Market).
			WithError(err).Errorln("unable to get orderbook for trade pair")
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

	if len(asks) == 0 {
		table.AddRow(color.RedString("No asks."), "", "")
	} else {
		for _, ask := range asks {
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

	if len(bids) == 0 {
		table.AddRow(color.GreenString("No bids."), "", "")
	} else {
		for _, bid := range bids {
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

func isMakerOf(order *sraAPI.Order, address common.Address) bool {
	return bytes.Compare(
		common.HexToAddress(order.MakerAddress).Bytes(),
		address.Bytes(),
	) == 0
}

func (ctl *AppController) getTokenNamesAndAssets(ctx context.Context) (tokenNames []string, assets []common.Address, err error) {
	if ctl.sraClient == nil {
		err := errors.New("client in offline mode")
		return nil, nil, err
	}

	pairs, err := ctl.restClient.TradePairs(ctx)
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

	if ctl.ethCore != nil {
		// always override WETH with client-side configured address
		tokenMap["WETH"] = ctl.ethCore.ContractAddress(ethcore.EthContractWETH9)
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
	if err == clients.ErrClientUnavailable {
		logrus.Errorln("Ethereum client is not initialized")
		return
	} else if err != nil {
		logrus.Errorln(err)
		return
	}

	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	ethBalance, err := ctl.ethCore.EthBalance(ctx, defaultAccount)
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

	allowances := ctl.ethCore.AllowancesMap(ctx, defaultAccount, common.HexToAddress(proxyAddressHex), assets)
	balances := ctl.ethCore.BalancesMap(ctx, defaultAccount, assets)

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
			isUnlocked := (allowances[addr].Cmp(ethcore.UnlimitedAllowance) == 0)
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

	pairs, err := ctl.restClient.TradePairs(ctx)
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
	if err == clients.ErrClientUnavailable {
		logrus.Errorln("Ethereum client is not initialized")
		return
	} else if err != nil {
		logrus.Errorln(err)
		return
	}

	tokenLockArgs := args.(*UtilTokenLockArgs)
	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
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

	txHash, err := ctl.ethCore.TokenLock(callArgs, asset)
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
	if err == clients.ErrClientUnavailable {
		logrus.Errorln("Ethereum client is not initialized")
		return
	} else if err != nil {
		logrus.Errorln(err)
		return
	}

	tokenUnlockArgs := args.(*UtilTokenUnlockArgs)
	defaultAccount := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	callArgs := &ethcore.CallArgs{
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

	txHash, err := ctl.ethCore.TokenUnlock(callArgs, asset)
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

	callArgs := &ethcore.CallArgs{
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

	txHash, err := ctl.ethCore.EthWrap(callArgs, amount)
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

	callArgs := &ethcore.CallArgs{
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

	txHash, err := ctl.ethCore.EthUnwrap(callArgs, amount)
	if err != nil {
		logrus.WithError(err).Errorln("unable to wrap WETH")
		return
	}

	fmt.Println(ctl.formatTxLink(txHash))
	ctl.checkTx(txHash)
}

func (ctl *AppController) ActionAccountsUse(args interface{}) {
	addr, err := ethcore.ParseAccount(args.(*ethcore.AccountUseArgs))
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

	if ctl.ethCore != nil {
		ctl.ethCore.SetDefaultFromAddress(addr)
	} else if err := ctl.initEthClient(); err != nil {
		logrus.WithError(err).Warningln("failed to init Ethereum client")
	}
}

func (ctl *AppController) ActionAccountsCreate(args interface{}) {
	acc, err := ethcore.CreateAccount(ctl.keystorePath, args.(*ethcore.AccountCreateArgs))
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
	addr, err := ethcore.ImportAccount(ctl.keystorePath, args.(*ethcore.AccountImportArgs))
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
	addr, err := ethcore.ImportPrivKey(ctl.keystorePath, args.(*ethcore.AccountImportPrivKeyArgs))
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

	pairs, err := ctl.restClient.TradePairs(ctx)
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

func (ctl *AppController) SuggestDerivativesMarkets() []prompt.Suggest {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	markets, err := ctl.restClient.DerivativeMarkets(ctx)
	if err != nil {
		logrus.WithError(err).Warningln("failed to fetch trade pairs")
		return nil
	}

	suggestions := make([]prompt.Suggest, 0, len(markets))
	for _, market := range markets {
		if !market.Enabled {
			continue
		}

		suggestions = append(suggestions, prompt.Suggest{
			Text: market.Ticker,
		})
	}

	return suggestions
}

func (ctl *AppController) SuggestOrderToFill(pairName string) []prompt.Suggest {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	bids, asks, err := ctl.sraClient.Orderbook(ctx, pairName)
	if err != nil {
		logrus.WithError(err).Warningln("failed to fetch orderbook")
		return nil
	}

	suggestions := make([]prompt.Suggest, 0, len(bids)+len(asks))
	for _, ask := range asks {
		// TODO: ignore own orders

		zxOrder, _ := ro2zo(ask.Order)
		orderHash, _ := zxOrder.ComputeOrderHash()
		price, vol := calcOrderPrice(ask.Order, false)
		vol = vol.Shift(-18)

		suggestions = append(suggestions, prompt.Suggest{
			Text:        orderHash.Hex(),
			Description: fmt.Sprintf("[ASK] %s %s", price.StringFixed(6), vol.StringFixed(6)),
		})
	}

	for _, bid := range bids {
		// TODO: ignore own orders

		zxOrder, _ := ro2zo(bid.Order)
		orderHash, _ := zxOrder.ComputeOrderHash()
		price, vol := calcOrderPrice(bid.Order, true)
		vol = vol.Shift(-18)

		suggestions = append(suggestions, prompt.Suggest{
			Text:        orderHash.Hex(),
			Description: fmt.Sprintf("[BID] %s %s", price.StringFixed(6), vol.StringFixed(6)),
		})
	}

	return suggestions
}

func (ctl *AppController) SuggestOrderToCancel(pairName string) []prompt.Suggest {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	owner := common.HexToAddress(ctl.mustConfigValue("accounts.default"))

	bids, asks, err := ctl.sraClient.Orderbook(ctx, pairName)
	if err != nil {
		logrus.WithError(err).Warningln("failed to fetch orderbook")
		return nil
	}

	suggestions := make([]prompt.Suggest, 0, len(bids)+len(asks))
	for _, ask := range asks {
		zxOrder, _ := ro2zo(ask.Order)
		if zxOrder.MakerAddress != owner {
			continue
		}

		orderHash, _ := zxOrder.ComputeOrderHash()
		price, vol := calcOrderPrice(ask.Order, false)
		vol = vol.Shift(-18)

		suggestions = append(suggestions, prompt.Suggest{
			Text:        orderHash.Hex(),
			Description: fmt.Sprintf("[ASK] %s %s", price.StringFixed(6), vol.StringFixed(6)),
		})
	}

	for _, bid := range bids {
		// TODO: ignore own orders

		zxOrder, _ := ro2zo(bid.Order)
		if zxOrder.MakerAddress != owner {
			continue
		}

		orderHash, _ := zxOrder.ComputeOrderHash()
		price, vol := calcOrderPrice(bid.Order, true)
		vol = vol.Shift(-18)

		suggestions = append(suggestions, prompt.Suggest{
			Text:        orderHash.Hex(),
			Description: fmt.Sprintf("[BID] %s %s", price.StringFixed(6), vol.StringFixed(6)),
		})
	}

	return suggestions
}

func getBidCurrency(market string) string {
	parts := strings.Split(market, "/")
	return parts[0]
}

func getAskCurrency(market string) string {
	parts := strings.Split(market, "/")
	return parts[1]
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

func (ctl *AppController) generateDefaultAccount() {
	const defaultPassword = "12345678"
	acc, err := ethcore.CreateAccount(ctl.keystorePath, &ethcore.AccountCreateArgs{
		Password:       defaultPassword,
		PasswordRepeat: defaultPassword,
	})
	if err != nil {
		logrus.WithError(err).Errorln("failed to generate default account")
		return
	}

	ctl.setConfigValue("accounts.default", acc.Address.Hex())

	logrus.WithFields(logrus.Fields{
		"account":    acc.Address.Hex(),
		"passphrase": defaultPassword,
	}).Infoln("Created a new default account, encrypting with insecure password.")

	logrus.Infoln("To import, create or switch your own accounts use keystore menu.")
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
	tx, err := ctl.ethCore.Ethereum().TransactionByHash(ctx, txHash.Hex())
	if err != nil {
		return err
	}

	t := time.NewTimer(time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			tx, err = ctl.ethCore.Ethereum().TransactionByHash(ctx, txHash.Hex())
			if err == nil && tx.BlockNumber != nil {
				receipt, err := ctl.ethCore.Ethereum().TransactionReceiptByHash(ctx, txHash.Hex())
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

	weth9AddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.weth9_address", networkName))
	erc20ProxyAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.erc20proxy_address", networkName))
	exchangeAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.exchange_address", networkName))
	futuresAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.futures_address", networkName))
	coordinatorAddressHex := ctl.mustConfigValue(fmt.Sprintf("networks.%s.coordinator_address", networkName))
	contractAddresses := map[ethcore.EthContract]common.Address{
		ethcore.EthContractWETH9:       common.HexToAddress(weth9AddressHex),
		ethcore.EthContractERC20Proxy:  common.HexToAddress(erc20ProxyAddressHex),
		ethcore.EthContractExchange:    common.HexToAddress(exchangeAddressHex),
		ethcore.EthContractFutures:     common.HexToAddress(futuresAddressHex),
		ethcore.EthContractCoordinator: common.HexToAddress(coordinatorAddressHex),
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

	ethCore, err := ethcore.New(
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

	ctl.ethCore = ethCore

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

func calcOrderPrice(order *sraAPI.Order, bid bool) (price, vol decimal.Decimal) {
	makerAmount := decimal.RequireFromString(order.MakerAssetAmount)
	takerAmount := decimal.RequireFromString(order.TakerAssetAmount)

	if bid { // i.e. buy
		price = makerAmount.DivRound(takerAmount, 9)
		vol = decimal.RequireFromString(order.TakerAssetAmount)
		return
	}

	price = takerAmount.DivRound(makerAmount, 9)
	vol = decimal.RequireFromString(order.MakerAssetAmount)
	return
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
