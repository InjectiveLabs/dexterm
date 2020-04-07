// Code generated by goa v3.1.1, DO NOT EDIT.
//
// ChronosAPI client HTTP transport
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

// Client lists the ChronosAPI service endpoint HTTP clients.
type Client struct {
	// SymbolInfo Doer is the HTTP client used to make requests to the symbolInfo
	// endpoint.
	SymbolInfoDoer goahttp.Doer

	// History Doer is the HTTP client used to make requests to the history
	// endpoint.
	HistoryDoer goahttp.Doer

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

// NewClient instantiates HTTP clients for all the ChronosAPI service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		SymbolInfoDoer:      doer,
		HistoryDoer:         doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// SymbolInfo returns an endpoint that makes HTTP requests to the ChronosAPI
// service symbolInfo server.
func (c *Client) SymbolInfo() goa.Endpoint {
	var (
		encodeRequest  = EncodeSymbolInfoRequest(c.encoder)
		decodeResponse = DecodeSymbolInfoResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildSymbolInfoRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.SymbolInfoDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("ChronosAPI", "symbolInfo", err)
		}
		return decodeResponse(resp)
	}
}

// History returns an endpoint that makes HTTP requests to the ChronosAPI
// service history server.
func (c *Client) History() goa.Endpoint {
	var (
		encodeRequest  = EncodeHistoryRequest(c.encoder)
		decodeResponse = DecodeHistoryResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildHistoryRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.HistoryDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("ChronosAPI", "history", err)
		}
		return decodeResponse(resp)
	}
}
