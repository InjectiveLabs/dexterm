// Code generated by goa v3.1.1, DO NOT EDIT.
//
// CoordinatorAPI client
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package coordinatorapi

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "CoordinatorAPI" service client.
type Client struct {
	ConfigurationEndpoint      goa.Endpoint
	RequestTransactionEndpoint goa.Endpoint
	SoftCancelsEndpoint        goa.Endpoint
}

// NewClient initializes a "CoordinatorAPI" service client given the endpoints.
func NewClient(configuration, requestTransaction, softCancels goa.Endpoint) *Client {
	return &Client{
		ConfigurationEndpoint:      configuration,
		RequestTransactionEndpoint: requestTransaction,
		SoftCancelsEndpoint:        softCancels,
	}
}

// Configuration calls the "configuration" endpoint of the "CoordinatorAPI"
// service.
func (c *Client) Configuration(ctx context.Context, p *ConfigurationPayload) (res *ConfigurationResult, err error) {
	var ires interface{}
	ires, err = c.ConfigurationEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ConfigurationResult), nil
}

// RequestTransaction calls the "request_transaction" endpoint of the
// "CoordinatorAPI" service.
func (c *Client) RequestTransaction(ctx context.Context, p *RequestTransactionPayload) (res *RequestTransactionResult, err error) {
	var ires interface{}
	ires, err = c.RequestTransactionEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*RequestTransactionResult), nil
}

// SoftCancels calls the "soft_cancels" endpoint of the "CoordinatorAPI"
// service.
func (c *Client) SoftCancels(ctx context.Context, p *SoftCancelsPayload) (res *SoftCancelsResult, err error) {
	var ires interface{}
	ires, err = c.SoftCancelsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SoftCancelsResult), nil
}
