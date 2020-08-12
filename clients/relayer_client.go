package clients

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
	"time"

	zeroex "github.com/InjectiveLabs/zeroex-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"

	sraHTTP "github.com/InjectiveLabs/dexterm/gen/http/relayer_api/client"
	sraAPI "github.com/InjectiveLabs/dexterm/gen/relayer_api"
	restAPI "github.com/InjectiveLabs/dexterm/gen/rest_api"
)

type SRAClient struct {
	cfg        *SRAClientConfig
	client     *sraAPI.Client
	restClient *RESTClient
}

type SRAClientConfig struct {
	Endpoint string
	Timeout  time.Duration
	Debug    bool
}

func (c *SRAClientConfig) check() *SRAClientConfig {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	return c
}

func NewSRAClient(restClient *RESTClient, cfg *SRAClientConfig) (*SRAClient, error) {
	u, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		err = errors.Wrap(err, "failed to parse endpoint URL")
		return nil, err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New("endpoint must have http:// or https:// scheme")
		return nil, err
	}

	cli := &SRAClient{
		cfg:        cfg.check(),
		client:     newSRAClient(u.Scheme, u.Host, cfg.Timeout, cfg.Debug),
		restClient: restClient,
	}

	return cli, nil
}

var ErrClientUnavailable = errors.New("offline mode: client is not available")

func (c *SRAClient) Orderbook(
	ctx context.Context,
	pairName string,
) (bids, asks []*sraAPI.OrderRecord, err error) {
	if c.client == nil {
		err = ErrClientUnavailable
		return
	}

	var pair *restAPI.TradePair

	if pair, err = c.getTradePairByName(ctx, pairName); err != nil {
		return
	} else if pair == nil {
		err = errors.New("trade pair not found")
		return
	}

	res, err := c.client.Orderbook(ctx, &sraAPI.OrderbookPayload{
		BaseAssetData:  pair.MakerAssetData,
		QuoteAssetData: pair.TakerAssetData,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to get orderbook")
		return
	}

	return res.Bids.Records, res.Asks.Records, nil
}

func (c *SRAClient) DerivativeOrders(
	ctx context.Context,
	assetData string,
) (bids, asks []*sraAPI.OrderRecord, err error) {
	if c.client == nil {
		err = ErrClientUnavailable
		return
	}
	emptyAssetStr := "0x000000000000000000000000000000000000000000000000000000000000000000000000"
	longRes, err := c.client.Orders(ctx, &sraAPI.OrdersPayload{
		MakerAssetData:    &assetData,
		TakerAssetData:    &emptyAssetStr,
		MakerFeeAssetData: &emptyAssetStr,
		TakerFeeAssetData: &emptyAssetStr,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to get derivative orders")
		return
	}

	shortRes, err := c.client.Orders(ctx, &sraAPI.OrdersPayload{
		MakerAssetData:    &emptyAssetStr,
		TakerAssetData:    &assetData,
		MakerFeeAssetData: &emptyAssetStr,
		TakerFeeAssetData: &emptyAssetStr,
	})

	if err != nil {
		err = errors.Wrap(err, "failed to get derivative orders")
		return
	}
	for _, order := range longRes.Records {
		bids = append(bids, order)
	}

	for _, order := range shortRes.Records {
		asks = append(asks, order)
	}
	return bids, asks, nil
}

func (c *SRAClient) PostOrder(
	ctx context.Context,
	order *zeroex.SignedOrder,
) (string, error) {
	orderHash, _ := order.ComputeOrderHash()

	orderPayload := &sraAPI.PostOrderPayload{
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

func (c *SRAClient) Order(ctx context.Context, orderHash string) (*sraAPI.Order, error) {
	if c.client == nil {
		return nil, errors.New("offline mode: SRA client is not available")
	}

	res, err := c.client.OrderByHash(ctx, &sraAPI.OrderByHashPayload{
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

func (c *SRAClient) FeeRecipients(ctx context.Context) (feeRecipients []common.Address, err error) {
	if c.client == nil {
		return nil, errors.New("offline mode: SRA client is not available")
	}

	res, err := c.client.FeeRecipients(ctx, &sraAPI.SRARequest{})
	if err != nil {
		err = errors.Wrap(err, "unable to get fee recipients")
		return nil, err
	}

	for _, addrHex := range res.List {
		feeRecipients = append(feeRecipients, common.HexToAddress(addrHex))
	}

	return feeRecipients, nil
}

func (c *SRAClient) getTradePairByName(ctx context.Context, pairName string) (*restAPI.TradePair, error) {
	var tradePairs []*restAPI.TradePair

	tradePairs, err := c.restClient.TradePairs(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to list supported trade pairs")
		return nil, err
	}

	var pair *restAPI.TradePair

	for _, currentPair := range tradePairs {
		if currentPair.Name == pairName {
			pair = currentPair
		}
	}

	return pair, nil
}

func newSRAClient(scheme, host string, timeout time.Duration, debug bool) *sraAPI.Client {
	var doer goahttp.Doer

	doer = &http.Client{
		Timeout: timeout,
	}

	if debug {
		doer = goahttp.NewDebugDoer(doer)
	}

	c := sraHTTP.NewClient(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
	)

	return sraAPI.NewClient(
		c.AssetPairs(),
		c.Orders(),
		c.OrderByHash(),
		c.Orderbook(),
		c.OrderConfig(),
		c.FeeRecipients(),
		c.PostOrder(),
	)
}
