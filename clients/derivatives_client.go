package clients

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
	"time"

	zeroex "github.com/InjectiveLabs/zeroex-go"
	//"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	goahttp "goa.design/goa/v3/http"
	//goa "goa.design/goa/v3/pkg"

	sdaHTTP "github.com/InjectiveLabs/dexterm/gen/http/derivatives_api/client"
	sdaAPI "github.com/InjectiveLabs/dexterm/gen/derivatives_api"
	//restAPI "github.com/InjectiveLabs/dexterm/gen/rest_api"
)

type SDAClient struct {
	cfg        *SDAClientConfig
	client     *sdaAPI.Client
	restClient *RESTClient
}

type SDAClientConfig struct {
	Endpoint string
	Timeout  time.Duration
	Debug    bool
}

func (c *SDAClientConfig) check() *SDAClientConfig {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	return c
}

func NewSDAClient(restClient *RESTClient, cfg *SDAClientConfig) (*SDAClient, error) {
	u, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		err = errors.Wrap(err, "failed to parse endpoint URL")
		return nil, err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New("endpoint must have http:// or https:// scheme")
		return nil, err
	}

	cli := &SDAClient{
		cfg:        cfg.check(),
		client:     newSDAClient(u.Scheme, u.Host, cfg.Timeout, cfg.Debug),
		restClient: restClient,
	}

	return cli, nil
}

func (c *SDAClient) PostOrder(
	ctx context.Context,
	order *zeroex.SignedOrder,
) (string, error) {
	orderHash, _ := order.ComputeOrderHash()

	orderPayload := &sdaAPI.PostOrderPayload{
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

//func (c *SDAClient) Orderbook(
//	ctx context.Context,
//	pairName string,
//) (bids, asks []*sdaAPI.DerivativeOrderRecord, err error) {
//	if c.client == nil {
//		err = ErrClientUnavailable
//		return
//	}
//
//	var pair *restAPI.TradePair
//
//	if pair, err = c.getTradePairByName(ctx, pairName); err != nil {
//		return
//	} else if pair == nil {
//		err = errors.New("trade pair not found")
//		return
//	}
//
//	res, err := c.client.Orderbook(ctx, &sraAPI.OrderbookPayload{
//		BaseAssetData:  pair.MakerAssetData,
//		QuoteAssetData: pair.TakerAssetData,
//	})
//	if err != nil {
//		err = errors.Wrap(err, "failed to get orderbook")
//		return
//	}
//
//	return res.Bids.Records, res.Asks.Records, nil
//}


func newSDAClient(scheme, host string, timeout time.Duration, debug bool) *sdaAPI.Client {
	var doer goahttp.Doer

	doer = &http.Client{
		Timeout: timeout,
	}

	if debug {
		doer = goahttp.NewDebugDoer(doer)
	}

	c := sdaHTTP.NewClient(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
	)

	return sdaAPI.NewClient(
		c.Orders(),
		c.PostOrder(),
	)
}
