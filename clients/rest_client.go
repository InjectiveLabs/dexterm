package clients

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"

	restHTTP "github.com/InjectiveLabs/dexterm/gen/http/rest_api/client"
	restAPI "github.com/InjectiveLabs/dexterm/gen/rest_api"
)

type RESTClient struct {
	cfg    *RESTClientConfig
	client *restAPI.Client
}

type RESTClientConfig struct {
	Endpoint string
	Timeout  time.Duration
	Debug    bool
}

func (c *RESTClientConfig) check() *RESTClientConfig {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	return c
}

func NewRESTClient(cfg *RESTClientConfig) (*RESTClient, error) {
	u, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		err = errors.Wrap(err, "failed to parse endpoint URL")
		return nil, err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New("endpoint must have http:// or https:// scheme")
		return nil, err
	}

	cli := &RESTClient{
		cfg:    cfg.check(),
		client: newRESTClient(u.Scheme, u.Host, cfg.Timeout, cfg.Debug),
	}

	return cli, nil
}

func (c *RESTClient) Order(ctx context.Context, orderHash string) (*restAPI.Order, error) {
	if c.client == nil {
		return nil, errors.New("offline mode: REST client is not available")
	}

	res, err := c.client.GetActiveOrder(ctx, &restAPI.GetActiveOrderPayload{
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

func (c *RESTClient) Orders(ctx context.Context, tradePairHash string) ([]*restAPI.Order, error) {
	if c.client == nil {
		return nil, ErrClientUnavailable
	}

	var collection = "active"
	res, err := c.client.ListOrders(ctx, &restAPI.ListOrdersPayload{
		Collection:    &collection,
		TradePairHash: &tradePairHash,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to list orders")
		return nil, err
	}

	return res.Orders, nil
}


func (c *RESTClient) TradePairs(ctx context.Context) ([]*restAPI.TradePair, error) {
	if c.client == nil {
		return nil, errors.New("offline mode: REST client is not available")
	}

	res, err := c.client.ListTradePairs(ctx, &restAPI.ListTradePairsPayload{})
	if err != nil {
		err = errors.Wrap(err, "failed to list trade pairs")
		return nil, err
	}

	return res.TradePairs, nil
}

func (c *RESTClient) DerivativeMarkets(ctx context.Context) ([]*restAPI.DerivativeMarket, error) {
	if c.client == nil {
		return nil, errors.New("offline mode: REST client is not available")
	}

	res, err := c.client.ListDerivativeMarkets(ctx, &restAPI.ListDerivativeMarketsPayload{})
	if err != nil {
		err = errors.Wrap(err, "failed to list trade pairs")
		return nil, err
	}

	return res.Markets, nil
}

func newRESTClient(scheme, host string, timeout time.Duration, debug bool) *restAPI.Client {
	var doer goahttp.Doer

	doer = &http.Client{
		Timeout: timeout,
	}

	if debug {
		doer = goahttp.NewDebugDoer(doer)
	}

	c := restHTTP.NewClient(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
	)

	return restAPI.NewClient(
		c.GetActiveOrder(),
		c.GetArchiveOrder(),
		c.ListOrders(),
		c.GetTradePair(),
		c.ListTradePairs(),
		c.ListDerivativeMarkets(),
		c.GetAccount(),
		c.GetOnlineAccounts(),
	)
}
