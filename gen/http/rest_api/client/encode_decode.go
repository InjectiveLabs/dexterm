// Code generated by goa v3.1.1, DO NOT EDIT.
//
// RestAPI HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	restapi "github.com/InjectiveLabs/dexterm/gen/rest_api"
	goahttp "goa.design/goa/v3/http"
)

// BuildGetActiveOrderRequest instantiates a HTTP request object with method
// and path set to call the "RestAPI" service "getActiveOrder" endpoint
func (c *Client) BuildGetActiveOrderRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetActiveOrderRestAPIPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("RestAPI", "getActiveOrder", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetActiveOrderRequest returns an encoder for requests sent to the
// RestAPI getActiveOrder server.
func EncodeGetActiveOrderRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*restapi.GetActiveOrderPayload)
		if !ok {
			return goahttp.ErrInvalidType("RestAPI", "getActiveOrder", "*restapi.GetActiveOrderPayload", v)
		}
		values := req.URL.Query()
		values.Add("orderHash", p.OrderHash)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetActiveOrderResponse returns a decoder for responses returned by the
// RestAPI getActiveOrder endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeGetActiveOrderResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "validation_error" (type *restapi.RESTValidationErrorResponse): http.StatusExpectationFailed
//	- error: internal error
func DecodeGetActiveOrderResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetActiveOrderResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getActiveOrder", err)
			}
			err = ValidateGetActiveOrderResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getActiveOrder", err)
			}
			res := NewGetActiveOrderResultOK(&body)
			return res, nil
		case http.StatusNotFound:
			var (
				body GetActiveOrderNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getActiveOrder", err)
			}
			err = ValidateGetActiveOrderNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getActiveOrder", err)
			}
			return nil, NewGetActiveOrderNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetActiveOrderInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getActiveOrder", err)
			}
			err = ValidateGetActiveOrderInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getActiveOrder", err)
			}
			return nil, NewGetActiveOrderInternal(&body)
		case http.StatusExpectationFailed:
			var (
				body GetActiveOrderValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getActiveOrder", err)
			}
			err = ValidateGetActiveOrderValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getActiveOrder", err)
			}
			return nil, NewGetActiveOrderValidationError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("RestAPI", "getActiveOrder", resp.StatusCode, string(body))
		}
	}
}

// BuildGetArchiveOrderRequest instantiates a HTTP request object with method
// and path set to call the "RestAPI" service "getArchiveOrder" endpoint
func (c *Client) BuildGetArchiveOrderRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetArchiveOrderRestAPIPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("RestAPI", "getArchiveOrder", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetArchiveOrderRequest returns an encoder for requests sent to the
// RestAPI getArchiveOrder server.
func EncodeGetArchiveOrderRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*restapi.GetArchiveOrderPayload)
		if !ok {
			return goahttp.ErrInvalidType("RestAPI", "getArchiveOrder", "*restapi.GetArchiveOrderPayload", v)
		}
		values := req.URL.Query()
		values.Add("orderHash", p.OrderHash)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetArchiveOrderResponse returns a decoder for responses returned by
// the RestAPI getArchiveOrder endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeGetArchiveOrderResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "validation_error" (type *restapi.RESTValidationErrorResponse): http.StatusExpectationFailed
//	- error: internal error
func DecodeGetArchiveOrderResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetArchiveOrderResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getArchiveOrder", err)
			}
			err = ValidateGetArchiveOrderResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getArchiveOrder", err)
			}
			res := NewGetArchiveOrderResultOK(&body)
			return res, nil
		case http.StatusNotFound:
			var (
				body GetArchiveOrderNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getArchiveOrder", err)
			}
			err = ValidateGetArchiveOrderNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getArchiveOrder", err)
			}
			return nil, NewGetArchiveOrderNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetArchiveOrderInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getArchiveOrder", err)
			}
			err = ValidateGetArchiveOrderInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getArchiveOrder", err)
			}
			return nil, NewGetArchiveOrderInternal(&body)
		case http.StatusExpectationFailed:
			var (
				body GetArchiveOrderValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getArchiveOrder", err)
			}
			err = ValidateGetArchiveOrderValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getArchiveOrder", err)
			}
			return nil, NewGetArchiveOrderValidationError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("RestAPI", "getArchiveOrder", resp.StatusCode, string(body))
		}
	}
}

// BuildListOrdersRequest instantiates a HTTP request object with method and
// path set to call the "RestAPI" service "listOrders" endpoint
func (c *Client) BuildListOrdersRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ListOrdersRestAPIPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("RestAPI", "listOrders", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeListOrdersRequest returns an encoder for requests sent to the RestAPI
// listOrders server.
func EncodeListOrdersRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*restapi.ListOrdersPayload)
		if !ok {
			return goahttp.ErrInvalidType("RestAPI", "listOrders", "*restapi.ListOrdersPayload", v)
		}
		body := NewListOrdersRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("RestAPI", "listOrders", err)
		}
		return nil
	}
}

// DecodeListOrdersResponse returns a decoder for responses returned by the
// RestAPI listOrders endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeListOrdersResponse may return the following errors:
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "validation_error" (type *restapi.RESTValidationErrorResponse): http.StatusExpectationFailed
//	- error: internal error
func DecodeListOrdersResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ListOrdersResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listOrders", err)
			}
			err = ValidateListOrdersResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listOrders", err)
			}
			res := NewListOrdersResultOK(&body)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ListOrdersInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listOrders", err)
			}
			err = ValidateListOrdersInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listOrders", err)
			}
			return nil, NewListOrdersInternal(&body)
		case http.StatusExpectationFailed:
			var (
				body ListOrdersValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listOrders", err)
			}
			err = ValidateListOrdersValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listOrders", err)
			}
			return nil, NewListOrdersValidationError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("RestAPI", "listOrders", resp.StatusCode, string(body))
		}
	}
}

// BuildGetTradePairRequest instantiates a HTTP request object with method and
// path set to call the "RestAPI" service "getTradePair" endpoint
func (c *Client) BuildGetTradePairRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetTradePairRestAPIPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("RestAPI", "getTradePair", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetTradePairRequest returns an encoder for requests sent to the
// RestAPI getTradePair server.
func EncodeGetTradePairRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*restapi.GetTradePairPayload)
		if !ok {
			return goahttp.ErrInvalidType("RestAPI", "getTradePair", "*restapi.GetTradePairPayload", v)
		}
		values := req.URL.Query()
		if p.Name != nil {
			values.Add("name", *p.Name)
		}
		if p.Hash != nil {
			values.Add("hash", *p.Hash)
		}
		if p.MakerAssetData != nil {
			values.Add("makerAssetData", *p.MakerAssetData)
		}
		if p.TakerAssetData != nil {
			values.Add("takerAssetData", *p.TakerAssetData)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetTradePairResponse returns a decoder for responses returned by the
// RestAPI getTradePair endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeGetTradePairResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "validation_error" (type *restapi.RESTValidationErrorResponse): http.StatusExpectationFailed
//	- error: internal error
func DecodeGetTradePairResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetTradePairResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getTradePair", err)
			}
			err = ValidateGetTradePairResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getTradePair", err)
			}
			res := NewGetTradePairResultOK(&body)
			return res, nil
		case http.StatusNotFound:
			var (
				body GetTradePairNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getTradePair", err)
			}
			err = ValidateGetTradePairNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getTradePair", err)
			}
			return nil, NewGetTradePairNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetTradePairInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getTradePair", err)
			}
			err = ValidateGetTradePairInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getTradePair", err)
			}
			return nil, NewGetTradePairInternal(&body)
		case http.StatusExpectationFailed:
			var (
				body GetTradePairValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getTradePair", err)
			}
			err = ValidateGetTradePairValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getTradePair", err)
			}
			return nil, NewGetTradePairValidationError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("RestAPI", "getTradePair", resp.StatusCode, string(body))
		}
	}
}

// BuildListTradePairsRequest instantiates a HTTP request object with method
// and path set to call the "RestAPI" service "listTradePairs" endpoint
func (c *Client) BuildListTradePairsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ListTradePairsRestAPIPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("RestAPI", "listTradePairs", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeListTradePairsRequest returns an encoder for requests sent to the
// RestAPI listTradePairs server.
func EncodeListTradePairsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*restapi.ListTradePairsPayload)
		if !ok {
			return goahttp.ErrInvalidType("RestAPI", "listTradePairs", "*restapi.ListTradePairsPayload", v)
		}
		body := NewListTradePairsRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("RestAPI", "listTradePairs", err)
		}
		return nil
	}
}

// DecodeListTradePairsResponse returns a decoder for responses returned by the
// RestAPI listTradePairs endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeListTradePairsResponse may return the following errors:
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "validation_error" (type *restapi.RESTValidationErrorResponse): http.StatusExpectationFailed
//	- error: internal error
func DecodeListTradePairsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ListTradePairsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listTradePairs", err)
			}
			err = ValidateListTradePairsResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listTradePairs", err)
			}
			res := NewListTradePairsResultOK(&body)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ListTradePairsInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listTradePairs", err)
			}
			err = ValidateListTradePairsInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listTradePairs", err)
			}
			return nil, NewListTradePairsInternal(&body)
		case http.StatusExpectationFailed:
			var (
				body ListTradePairsValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listTradePairs", err)
			}
			err = ValidateListTradePairsValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listTradePairs", err)
			}
			return nil, NewListTradePairsValidationError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("RestAPI", "listTradePairs", resp.StatusCode, string(body))
		}
	}
}

// BuildListDerivativeMarketsRequest instantiates a HTTP request object with
// method and path set to call the "RestAPI" service "listDerivativeMarkets"
// endpoint
func (c *Client) BuildListDerivativeMarketsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ListDerivativeMarketsRestAPIPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("RestAPI", "listDerivativeMarkets", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeListDerivativeMarketsResponse returns a decoder for responses returned
// by the RestAPI listDerivativeMarkets endpoint. restoreBody controls whether
// the response body should be restored after having been read.
// DecodeListDerivativeMarketsResponse may return the following errors:
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "validation_error" (type *restapi.RESTValidationErrorResponse): http.StatusExpectationFailed
//	- error: internal error
func DecodeListDerivativeMarketsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ListDerivativeMarketsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listDerivativeMarkets", err)
			}
			err = ValidateListDerivativeMarketsResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listDerivativeMarkets", err)
			}
			res := NewListDerivativeMarketsResultOK(&body)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ListDerivativeMarketsInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listDerivativeMarkets", err)
			}
			err = ValidateListDerivativeMarketsInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listDerivativeMarkets", err)
			}
			return nil, NewListDerivativeMarketsInternal(&body)
		case http.StatusExpectationFailed:
			var (
				body ListDerivativeMarketsValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "listDerivativeMarkets", err)
			}
			err = ValidateListDerivativeMarketsValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "listDerivativeMarkets", err)
			}
			return nil, NewListDerivativeMarketsValidationError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("RestAPI", "listDerivativeMarkets", resp.StatusCode, string(body))
		}
	}
}

// BuildGetAccountRequest instantiates a HTTP request object with method and
// path set to call the "RestAPI" service "getAccount" endpoint
func (c *Client) BuildGetAccountRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetAccountRestAPIPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("RestAPI", "getAccount", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetAccountRequest returns an encoder for requests sent to the RestAPI
// getAccount server.
func EncodeGetAccountRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*restapi.GetAccountPayload)
		if !ok {
			return goahttp.ErrInvalidType("RestAPI", "getAccount", "*restapi.GetAccountPayload", v)
		}
		values := req.URL.Query()
		values.Add("address", p.Address)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetAccountResponse returns a decoder for responses returned by the
// RestAPI getAccount endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeGetAccountResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "validation_error" (type *restapi.RESTValidationErrorResponse): http.StatusExpectationFailed
//	- error: internal error
func DecodeGetAccountResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetAccountResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getAccount", err)
			}
			err = ValidateGetAccountResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getAccount", err)
			}
			res := NewGetAccountResultOK(&body)
			return res, nil
		case http.StatusNotFound:
			var (
				body GetAccountNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getAccount", err)
			}
			err = ValidateGetAccountNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getAccount", err)
			}
			return nil, NewGetAccountNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetAccountInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getAccount", err)
			}
			err = ValidateGetAccountInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getAccount", err)
			}
			return nil, NewGetAccountInternal(&body)
		case http.StatusExpectationFailed:
			var (
				body GetAccountValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getAccount", err)
			}
			err = ValidateGetAccountValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getAccount", err)
			}
			return nil, NewGetAccountValidationError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("RestAPI", "getAccount", resp.StatusCode, string(body))
		}
	}
}

// BuildGetOnlineAccountsRequest instantiates a HTTP request object with method
// and path set to call the "RestAPI" service "getOnlineAccounts" endpoint
func (c *Client) BuildGetOnlineAccountsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetOnlineAccountsRestAPIPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("RestAPI", "getOnlineAccounts", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetOnlineAccountsRequest returns an encoder for requests sent to the
// RestAPI getOnlineAccounts server.
func EncodeGetOnlineAccountsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*restapi.GetOnlineAccountsPayload)
		if !ok {
			return goahttp.ErrInvalidType("RestAPI", "getOnlineAccounts", "*restapi.GetOnlineAccountsPayload", v)
		}
		values := req.URL.Query()
		if p.Version != nil {
			values.Add("version", *p.Version)
		}
		if p.Threshold != nil {
			values.Add("threshold", fmt.Sprintf("%v", *p.Threshold))
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetOnlineAccountsResponse returns a decoder for responses returned by
// the RestAPI getOnlineAccounts endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeGetOnlineAccountsResponse may return the following errors:
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "validation_error" (type *restapi.RESTValidationErrorResponse): http.StatusExpectationFailed
//	- error: internal error
func DecodeGetOnlineAccountsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetOnlineAccountsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getOnlineAccounts", err)
			}
			err = ValidateGetOnlineAccountsResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getOnlineAccounts", err)
			}
			res := NewGetOnlineAccountsResultOK(&body)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body GetOnlineAccountsInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getOnlineAccounts", err)
			}
			err = ValidateGetOnlineAccountsInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getOnlineAccounts", err)
			}
			return nil, NewGetOnlineAccountsInternal(&body)
		case http.StatusExpectationFailed:
			var (
				body GetOnlineAccountsValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("RestAPI", "getOnlineAccounts", err)
			}
			err = ValidateGetOnlineAccountsValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("RestAPI", "getOnlineAccounts", err)
			}
			return nil, NewGetOnlineAccountsValidationError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("RestAPI", "getOnlineAccounts", resp.StatusCode, string(body))
		}
	}
}

// unmarshalOrderResponseBodyToRestapiOrder builds a value of type
// *restapi.Order from a value of type *OrderResponseBody.
func unmarshalOrderResponseBodyToRestapiOrder(v *OrderResponseBody) *restapi.Order {
	if v == nil {
		return nil
	}
	res := &restapi.Order{
		ChainID:               *v.ChainID,
		ExchangeAddress:       *v.ExchangeAddress,
		MakerAddress:          *v.MakerAddress,
		TakerAddress:          *v.TakerAddress,
		FeeRecipientAddress:   *v.FeeRecipientAddress,
		SenderAddress:         *v.SenderAddress,
		MakerAssetAmount:      *v.MakerAssetAmount,
		TakerAssetAmount:      *v.TakerAssetAmount,
		MakerFee:              *v.MakerFee,
		TakerFee:              *v.TakerFee,
		ExpirationTimeSeconds: *v.ExpirationTimeSeconds,
		Salt:                  *v.Salt,
		MakerAssetData:        *v.MakerAssetData,
		TakerAssetData:        *v.TakerAssetData,
		MakerFeeAssetData:     *v.MakerFeeAssetData,
		TakerFeeAssetData:     *v.TakerFeeAssetData,
		Signature:             *v.Signature,
	}

	return res
}

// unmarshalRESTValidationErrorResponseBodyToRestapiRESTValidationError builds
// a value of type *restapi.RESTValidationError from a value of type
// *RESTValidationErrorResponseBody.
func unmarshalRESTValidationErrorResponseBodyToRestapiRESTValidationError(v *RESTValidationErrorResponseBody) *restapi.RESTValidationError {
	if v == nil {
		return nil
	}
	res := &restapi.RESTValidationError{
		Code:   *v.Code,
		Reason: *v.Reason,
		Field:  v.Field,
	}

	return res
}

// unmarshalTradePairResponseBodyToRestapiTradePair builds a value of type
// *restapi.TradePair from a value of type *TradePairResponseBody.
func unmarshalTradePairResponseBodyToRestapiTradePair(v *TradePairResponseBody) *restapi.TradePair {
	if v == nil {
		return nil
	}
	res := &restapi.TradePair{
		Name:           *v.Name,
		MakerAssetData: *v.MakerAssetData,
		TakerAssetData: *v.TakerAssetData,
		Hash:           *v.Hash,
		Enabled:        *v.Enabled,
	}

	return res
}

// unmarshalDerivativeMarketResponseBodyToRestapiDerivativeMarket builds a
// value of type *restapi.DerivativeMarket from a value of type
// *DerivativeMarketResponseBody.
func unmarshalDerivativeMarketResponseBodyToRestapiDerivativeMarket(v *DerivativeMarketResponseBody) *restapi.DerivativeMarket {
	if v == nil {
		return nil
	}
	res := &restapi.DerivativeMarket{
		Ticker:       *v.Ticker,
		Oracle:       *v.Oracle,
		BaseCurrency: *v.BaseCurrency,
		Nonce:        *v.Nonce,
		MarketID:     *v.MarketID,
		Enabled:      *v.Enabled,
	}

	return res
}

// unmarshalRelayerAccountResponseBodyToRestapiRelayerAccount builds a value of
// type *restapi.RelayerAccount from a value of type
// *RelayerAccountResponseBody.
func unmarshalRelayerAccountResponseBodyToRestapiRelayerAccount(v *RelayerAccountResponseBody) *restapi.RelayerAccount {
	if v == nil {
		return nil
	}
	res := &restapi.RelayerAccount{
		Address:       *v.Address,
		StakerAddress: v.StakerAddress,
		PublicKey:     *v.PublicKey,
		LastSeen:      *v.LastSeen,
		LastVersion:   *v.LastVersion,
		IsOnline:      *v.IsOnline,
	}

	return res
}
