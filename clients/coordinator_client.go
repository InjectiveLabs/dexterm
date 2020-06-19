package clients

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	zeroex "github.com/InjectiveLabs/zeroex-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	goahttp "goa.design/goa/v3/http"

	coordinatorAPI "github.com/InjectiveLabs/dexterm/gen/coordinator_api"
	coordinatorHTTP "github.com/InjectiveLabs/dexterm/gen/http/coordinator_api/client"
)

type CoordinatorClient struct {
	cfg    *CoordinatorClientConfig
	client *coordinatorAPI.Client
}

type CoordinatorClientConfig struct {
	Endpoint string
	Timeout  time.Duration
	Debug    bool
}

func (c *CoordinatorClientConfig) check() *CoordinatorClientConfig {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	return c
}

func NewCoordinatorClient(cfg *CoordinatorClientConfig) (*CoordinatorClient, error) {
	u, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		err = errors.Wrap(err, "failed to parse endpoint URL")
		return nil, err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New("endpoint must have http:// or https:// scheme")
		return nil, err
	}

	cli := &CoordinatorClient{
		cfg:    cfg.check(),
		client: newCoordinatorClient(u.Scheme, u.Host, cfg.Timeout, cfg.Debug),
	}

	return cli, nil
}

func (c *CoordinatorClient) GetCoordinatorApproval(
	ctx context.Context,
	tx *zeroex.SignedTransaction,
	txOrigin common.Address,
) (approvals [][]byte, expiryAt time.Time, err error) {
	resp, err := c.client.RequestTransaction(ctx, &coordinatorAPI.RequestTransactionPayload{
		SignedTransaction: ztx2ctx(tx),
		TxOrigin:          txOrigin.Hex(),
	})
	if err != nil {
		err = errors.Wrap(err, "failed to request transaction using Coordinator API")
		return
	}

	if resp.ExpirationTimeSeconds != nil {
		ts, _ := strconv.ParseInt(*resp.ExpirationTimeSeconds, 10, 64)
		expiryAt = time.Unix(ts, 0).UTC()
	}

	approvals = make([][]byte, len(resp.Signatures))
	for idx, sig := range resp.Signatures {
		approvals[idx] = common.FromHex(sig)
	}

	return approvals, expiryAt, nil
}

func (c *CoordinatorClient) SendCoordinatorSoftCancelTransaction(
	ctx context.Context,
	tx *zeroex.SignedTransaction,
	txOrigin common.Address,
) (cancellationSigs []string, err error) {
	resp, err := c.client.RequestTransaction(ctx, &coordinatorAPI.RequestTransactionPayload{
		SignedTransaction: ztx2ctx(tx),
		TxOrigin:          txOrigin.Hex(),
	})
	if err != nil {
		err = errors.Wrap(err, "failed to request transaction using Coordinator API")
		return
	}
	return resp.CancellationSignatures, nil
}

func ztx2ctx(tx *zeroex.SignedTransaction) *coordinatorAPI.SignedTransaction {
	ctx := &coordinatorAPI.SignedTransaction{
		Salt:                  tx.Salt.String(),
		SignerAddress:         tx.SignerAddress.Hex(),
		ExpirationTimeSeconds: tx.ExpirationTimeSeconds.String(),
		GasPrice:              tx.GasPrice.String(),

		Domain: &coordinatorAPI.ExchangeDomain{
			VerifyingContract: tx.Domain.VerifyingContract.Hex(),
			ChainID:           tx.Domain.ChainID.String(),
		},

		Data:      common.ToHex(tx.Data),
		Signature: common.ToHex(tx.Signature),
	}

	return ctx
}

func newCoordinatorClient(scheme, host string, timeout time.Duration, debug bool) *coordinatorAPI.Client {
	var doer goahttp.Doer

	doer = &http.Client{
		Timeout: timeout,
	}

	if debug {
		doer = goahttp.NewDebugDoer(doer)
	}

	c := coordinatorHTTP.NewClient(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
	)

	return coordinatorAPI.NewClient(
		c.Configuration(),
		c.RequestTransaction(),
		c.SoftCancels(),
	)
}
