package main

import (
	"context"
	"net/http"
	"net/url"
	"time"

	relayerHTTPClient "github.com/InjectiveLabs/injective-core/api/gen/http/relayer/client"
	relayer "github.com/InjectiveLabs/injective-core/api/gen/relayer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	goahttp "goa.design/goa/v3/http"
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

type TradeBuyArgs struct {
	Foo    string
	Bar    string
	Baz    float64
	Params []int
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
