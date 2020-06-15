package clients

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	goahttp "goa.design/goa/v3/http"

	debugAPI "github.com/InjectiveLabs/dexterm/gen/debug_api"
	debugHTTP "github.com/InjectiveLabs/dexterm/gen/http/debug_api/client"
)

type DebugClient struct {
	cfg    *DebugClientConfig
	client *debugAPI.Client

	coreVersion string
}

type DebugClientConfig struct {
	Endpoint string
	Timeout  time.Duration
	Debug    bool
}

func (c *DebugClientConfig) check() *DebugClientConfig {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	return c
}

func NewDebugClient(cfg *DebugClientConfig) (*DebugClient, error) {
	u, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		err = errors.Wrap(err, "failed to parse endpoint URL")
		return nil, err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New("endpoint must have http:// or https:// scheme")
		return nil, err
	}

	cli := &DebugClient{
		cfg:    cfg.check(),
		client: newDebugClient(u.Scheme, u.Host, cfg.Timeout, cfg.Debug),
	}

	ver, err := cli.Version(context.Background())
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

func (c *DebugClient) Version(ctx context.Context) (string, error) {
	if c.client == nil {
		return "", ErrClientUnavailable
	}

	res, err := c.client.Version(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to get version")
		return "", err
	}

	return res.Version, nil
}

func newDebugClient(scheme, host string, timeout time.Duration, debug bool) *debugAPI.Client {
	var doer goahttp.Doer

	doer = &http.Client{
		Timeout: timeout,
	}

	if debug {
		doer = goahttp.NewDebugDoer(doer)
	}

	c := debugHTTP.NewClient(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
	)

	return debugAPI.NewClient(
		c.Version(),
	)
}
