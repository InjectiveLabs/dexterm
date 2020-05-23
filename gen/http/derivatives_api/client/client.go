// Code generated by goa v3.1.1, DO NOT EDIT.
//
// DerivativesAPI client HTTP transport
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package client

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the DerivativesAPI service endpoint HTTP clients.
type Client struct {
	// Orders Doer is the HTTP client used to make requests to the orders endpoint.
	OrdersDoer goahttp.Doer

	// PostOrder Doer is the HTTP client used to make requests to the postOrder
	// endpoint.
	PostOrderDoer goahttp.Doer

	// CORS Doer is the HTTP client used to make requests to the  endpoint.
	CORSDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the DerivativesAPI service
// servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		OrdersDoer:          doer,
		PostOrderDoer:       doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// Orders returns an endpoint that makes HTTP requests to the DerivativesAPI
// service orders server.
func (c *Client) Orders() goa.Endpoint {
	var (
		encodeRequest  = EncodeOrdersRequest(c.encoder)
		decodeResponse = DecodeOrdersResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildOrdersRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.OrdersDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("DerivativesAPI", "orders", err)
		}
		return decodeResponse(resp)
	}
}

// PostOrder returns an endpoint that makes HTTP requests to the DerivativesAPI
// service postOrder server.
func (c *Client) PostOrder() goa.Endpoint {
	var (
		encodeRequest  = EncodePostOrderRequest(c.encoder)
		decodeResponse = DecodePostOrderResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildPostOrderRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.PostOrderDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("DerivativesAPI", "postOrder", err)
		}
		return decodeResponse(resp)
	}
}