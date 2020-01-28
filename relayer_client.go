package main

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
	"time"

	relayerHTTPClient "github.com/InjectiveLabs/injective-core/api/gen/http/relayer/client"
	relayer "github.com/InjectiveLabs/injective-core/api/gen/relayer"
	zeroex "github.com/InjectiveLabs/zeroex-go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

type RelayerClient struct {
	cfg         *RelayerClientConfig
	coreVersion string
	client      *relayer.Client
}

type RelayerClientConfig struct {
	Endpoint string
	Timeout  time.Duration
	Debug    bool
}

func (c *RelayerClientConfig) check() *RelayerClientConfig {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	return c
}

func NewRelayerClient(cfg *RelayerClientConfig) (*RelayerClient, error) {
	u, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		err = errors.Wrap(err, "failed to parse endpoint URL")
		return nil, err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New("endpoint must have http:// or https:// scheme")
		return nil, err
	}

	cli := &RelayerClient{
		cfg:    cfg.check(),
		client: newRelayerClient(u.Scheme, u.Host, cfg.Timeout, cfg.Debug),
	}

	ver, err := cli.getVersion(context.Background())
	if err != nil {
		err = errors.Wrap(err, "could not init relayer client")
		return nil, err
	}

	cli.coreVersion = ver

	logrus.WithFields(logrus.Fields{
		"version": ver,
	}).Debugf("Connected to relayer at %s", cfg.Endpoint)

	return cli, nil
}

func (c *RelayerClient) getVersion(ctx context.Context) (string, error) {
	res, err := c.client.Version(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to get version")
		return "", err
	}

	return res.Version, nil
}

var ErrClientUnavailable = errors.New("offline mode: client is not available")

func (c *RelayerClient) TradePairs(ctx context.Context) ([]*relayer.TradePair, error) {
	if c.client == nil {
		return nil, ErrClientUnavailable
	}

	res, err := c.client.ListTradePairs(ctx, &relayer.ListTradePairsPayload{})
	if err != nil {
		err = errors.Wrap(err, "failed to list trade pairs")
		return nil, err
	}

	return res.TradePairs, nil
}

func (c *RelayerClient) Orderbook(
	ctx context.Context,
	pairName string,
) (bids, asks *relayer.OrderbookRecords, err error) {
	if c.client == nil {
		err = ErrClientUnavailable
		return
	}

	var tradePairs []*relayer.TradePair

	if tradePairs, err = c.TradePairs(ctx); err != nil {
		return
	}

	var pair *relayer.TradePair

	for _, currentPair := range tradePairs {
		if currentPair.Name == pairName {
			pair = currentPair
		}
	}

	if pair == nil {
		err = errors.New("trade pair not found")
		return
	}

	res, err := c.client.Orderbook(ctx, &relayer.OrderbookPayload{
		BaseAssetData:  pair.MakerAssetData,
		QuoteAssetData: pair.TakerAssetData,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to get orderbook")
		return
	}

	return res.Bids, res.Asks, nil
}

func (c *RelayerClient) Order(ctx context.Context, orderHash string) (*relayer.Order, error) {
	if c.client == nil {
		return nil, ErrClientUnavailable
	}

	res, err := c.client.GetActiveOrder(ctx, &relayer.GetActiveOrderPayload{
		OrderHash: orderHash,
	})
	if err != nil {
		if serviceError, ok := err.(*goa.ServiceError); ok {
			if serviceError.ErrorName() == "not_found" {
				err = errors.New("order not found")
				return nil, err
			}
		}

		err = errors.Wrap(err, "unable to get order")
		return nil, err
	}

	return res.Order, nil
}

func (c *RelayerClient) Orders(ctx context.Context, tradePairHash string) ([]*relayer.Order, error) {
	if c.client == nil {
		return nil, ErrClientUnavailable
	}

	var collection = "active"
	res, err := c.client.ListOrders(ctx, &relayer.ListOrdersPayload{
		Collection:    &collection,
		TradePairHash: &tradePairHash,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to list orders")
		return nil, err
	}

	return res.MakeOrders, nil
}

func (c *RelayerClient) PostOrder(
	ctx context.Context,
	order *zeroex.SignedOrder,
) (string, error) {
	orderHash, _ := order.ComputeOrderHash()

	orderPayload := &relayer.PostOrderPayload{
		ChainID: order.ChainID.Int64(),

		ExchangeAddress:     strings.ToLower(order.ExchangeAddress.Hex()),
		MakerAddress:        strings.ToLower(order.MakerAddress.Hex()),
		TakerAddress:        strings.ToLower(order.TakerAddress.Hex()),
		FeeRecipientAddress: strings.ToLower(order.FeeRecipientAddress.Hex()),
		SenderAddress:       strings.ToLower(order.SenderAddress.Hex()),

		MakerAssetAmount: order.MakerAssetAmount.String(),
		TakerAssetAmount: order.TakerAssetAmount.String(),
		MakerFee:         order.MakerFee.String(),
		TakerFee:         order.TakerFee.String(),

		ExpirationTimeSeconds: order.ExpirationTimeSeconds.String(),
		Salt:                  order.Salt.String(),

		MakerAssetData:    "0x" + hex.EncodeToString(order.MakerAssetData),
		TakerAssetData:    "0x" + hex.EncodeToString(order.TakerAssetData),
		MakerFeeAssetData: "0x" + hex.EncodeToString(order.MakerFeeAssetData),
		TakerFeeAssetData: "0x" + hex.EncodeToString(order.TakerFeeAssetData),
		Signature:         "0x" + hex.EncodeToString(order.Signature),
	}

	_, err := c.client.PostOrder(ctx, orderPayload)

	return orderHash.String(), err
}

func (c *RelayerClient) TakeOrder(
	ctx context.Context,
	makeOrders []*relayer.Order,
	takeOrder *zeroex.SignedOrder,
) (string, error) {
	takeOrderHash, _ := takeOrder.ComputeOrderHash()

	orderPayload := &relayer.TakeOrderPayload{
		MakeOrders: makeOrders,
		MakeOrderFillAmounts: []string{
			takeOrder.MakerAssetAmount.String(),
		},
		TakeOrder: &relayer.Order{
			ChainID:               takeOrder.ChainID.Int64(),
			ExchangeAddress:       strings.ToLower(takeOrder.ExchangeAddress.Hex()),
			MakerAddress:          strings.ToLower(takeOrder.MakerAddress.Hex()),
			TakerAddress:          strings.ToLower(takeOrder.TakerAddress.Hex()),
			FeeRecipientAddress:   strings.ToLower(takeOrder.FeeRecipientAddress.Hex()),
			SenderAddress:         strings.ToLower(takeOrder.SenderAddress.Hex()),
			MakerAssetAmount:      takeOrder.MakerAssetAmount.String(),
			TakerAssetAmount:      takeOrder.TakerAssetAmount.String(),
			MakerFee:              takeOrder.MakerFee.String(),
			TakerFee:              takeOrder.TakerFee.String(),
			ExpirationTimeSeconds: takeOrder.ExpirationTimeSeconds.String(),
			Salt:                  takeOrder.Salt.String(),
			MakerAssetData:        "0x" + hex.EncodeToString(takeOrder.MakerAssetData),
			TakerAssetData:        "0x" + hex.EncodeToString(takeOrder.TakerAssetData),
			MakerFeeAssetData:     "0x" + hex.EncodeToString(takeOrder.MakerFeeAssetData),
			TakerFeeAssetData:     "0x" + hex.EncodeToString(takeOrder.TakerFeeAssetData),
			Signature:             "0x" + hex.EncodeToString(takeOrder.Signature),
		},
	}

	_, err := c.client.TakeOrder(ctx, orderPayload)

	return takeOrderHash.String(), err
}

func calcOrderPrice(order *relayer.Order, bid bool) (price, vol decimal.Decimal) {
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

func newRelayerClient(scheme, host string, timeout time.Duration, debug bool) *relayer.Client {
	var doer goahttp.Doer

	doer = &http.Client{
		Timeout: timeout,
	}

	if debug {
		doer = goahttp.NewDebugDoer(doer)
	}

	c := relayerHTTPClient.NewClient(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
	)

	return relayer.NewClient(
		c.AssetPairs(),
		c.Orders(),
		c.OrderByHash(),
		c.Orderbook(),
		c.OrderConfig(),
		c.FeeRecipients(),
		c.PostOrder(),
		c.TakeOrder(),
		c.GetActiveOrder(),
		c.GetArchiveOrder(),
		c.ListOrders(),
		c.GetTradePair(),
		c.ListTradePairs(),
		c.GetAccount(),
		c.GetOnlineAccounts(),
		c.GetEthTransactions(),
		c.Version(),
	)
}
