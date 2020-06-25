// Code generated by goa v3.1.1, DO NOT EDIT.
//
// DerivativesAPI HTTP client encoders and decoders
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

	derivativesapi "github.com/InjectiveLabs/dexterm/gen/derivatives_api"
	goahttp "goa.design/goa/v3/http"
)

// BuildOrdersRequest instantiates a HTTP request object with method and path
// set to call the "DerivativesAPI" service "orders" endpoint
func (c *Client) BuildOrdersRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: OrdersDerivativesAPIPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("DerivativesAPI", "orders", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeOrdersRequest returns an encoder for requests sent to the
// DerivativesAPI orders server.
func EncodeOrdersRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*derivativesapi.OrdersPayload)
		if !ok {
			return goahttp.ErrInvalidType("DerivativesAPI", "orders", "*derivativesapi.OrdersPayload", v)
		}
		values := req.URL.Query()
		values.Add("page", fmt.Sprintf("%v", p.Page))
		values.Add("perPage", fmt.Sprintf("%v", p.PerPage))
		if p.TakerAssetData != nil {
			values.Add("takerAssetData", *p.TakerAssetData)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeOrdersResponse returns a decoder for responses returned by the
// DerivativesAPI orders endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeOrdersResponse may return the following errors:
//	- "validation_error" (type *derivativesapi.SDAValidationErrorResponse): http.StatusExpectationFailed
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "rate_limit" (type *goa.ServiceError): http.StatusTooManyRequests
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "not_implemented" (type *goa.ServiceError): http.StatusNotImplemented
//	- error: internal error
func DecodeOrdersResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body OrdersResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "orders", err)
			}
			err = ValidateOrdersResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "orders", err)
			}
			res := NewOrdersResultOK(&body)
			return res, nil
		case http.StatusExpectationFailed:
			var (
				body OrdersValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "orders", err)
			}
			err = ValidateOrdersValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "orders", err)
			}
			return nil, NewOrdersValidationError(&body)
		case http.StatusNotFound:
			var (
				body OrdersNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "orders", err)
			}
			err = ValidateOrdersNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "orders", err)
			}
			return nil, NewOrdersNotFound(&body)
		case http.StatusTooManyRequests:
			var (
				body OrdersRateLimitResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "orders", err)
			}
			err = ValidateOrdersRateLimitResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "orders", err)
			}
			return nil, NewOrdersRateLimit(&body)
		case http.StatusInternalServerError:
			var (
				body OrdersInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "orders", err)
			}
			err = ValidateOrdersInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "orders", err)
			}
			return nil, NewOrdersInternal(&body)
		case http.StatusNotImplemented:
			var (
				body OrdersNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "orders", err)
			}
			err = ValidateOrdersNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "orders", err)
			}
			return nil, NewOrdersNotImplemented(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("DerivativesAPI", "orders", resp.StatusCode, string(body))
		}
	}
}

// BuildPostOrderRequest instantiates a HTTP request object with method and
// path set to call the "DerivativesAPI" service "postOrder" endpoint
func (c *Client) BuildPostOrderRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: PostOrderDerivativesAPIPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("DerivativesAPI", "postOrder", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodePostOrderRequest returns an encoder for requests sent to the
// DerivativesAPI postOrder server.
func EncodePostOrderRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*derivativesapi.PostOrderPayload)
		if !ok {
			return goahttp.ErrInvalidType("DerivativesAPI", "postOrder", "*derivativesapi.PostOrderPayload", v)
		}
		body := NewPostOrderRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("DerivativesAPI", "postOrder", err)
		}
		return nil
	}
}

// DecodePostOrderResponse returns a decoder for responses returned by the
// DerivativesAPI postOrder endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodePostOrderResponse may return the following errors:
//	- "validation_error" (type *derivativesapi.SDAValidationErrorResponse): http.StatusExpectationFailed
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "rate_limit" (type *goa.ServiceError): http.StatusTooManyRequests
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- "not_implemented" (type *goa.ServiceError): http.StatusNotImplemented
//	- error: internal error
func DecodePostOrderResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusCreated:
			var (
				body PostOrderResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "postOrder", err)
			}
			res := NewPostOrderResultCreated(&body)
			return res, nil
		case http.StatusExpectationFailed:
			var (
				body PostOrderValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "postOrder", err)
			}
			err = ValidatePostOrderValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "postOrder", err)
			}
			return nil, NewPostOrderValidationError(&body)
		case http.StatusNotFound:
			var (
				body PostOrderNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "postOrder", err)
			}
			err = ValidatePostOrderNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "postOrder", err)
			}
			return nil, NewPostOrderNotFound(&body)
		case http.StatusTooManyRequests:
			var (
				body PostOrderRateLimitResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "postOrder", err)
			}
			err = ValidatePostOrderRateLimitResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "postOrder", err)
			}
			return nil, NewPostOrderRateLimit(&body)
		case http.StatusInternalServerError:
			var (
				body PostOrderInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "postOrder", err)
			}
			err = ValidatePostOrderInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "postOrder", err)
			}
			return nil, NewPostOrderInternal(&body)
		case http.StatusNotImplemented:
			var (
				body PostOrderNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("DerivativesAPI", "postOrder", err)
			}
			err = ValidatePostOrderNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("DerivativesAPI", "postOrder", err)
			}
			return nil, NewPostOrderNotImplemented(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("DerivativesAPI", "postOrder", resp.StatusCode, string(body))
		}
	}
}

// unmarshalDerivativeOrderRecordResponseBodyToDerivativesapiDerivativeOrderRecord
// builds a value of type *derivativesapi.DerivativeOrderRecord from a value of
// type *DerivativeOrderRecordResponseBody.
func unmarshalDerivativeOrderRecordResponseBodyToDerivativesapiDerivativeOrderRecord(v *DerivativeOrderRecordResponseBody) *derivativesapi.DerivativeOrderRecord {
	res := &derivativesapi.DerivativeOrderRecord{}
	res.DerivativeOrder = unmarshalDerivativeOrderResponseBodyToDerivativesapiDerivativeOrder(v.DerivativeOrder)
	res.MetaData = make(map[string]string, len(v.MetaData))
	for key, val := range v.MetaData {
		tk := key
		tv := val
		res.MetaData[tk] = tv
	}

	return res
}

// unmarshalDerivativeOrderResponseBodyToDerivativesapiDerivativeOrder builds a
// value of type *derivativesapi.DerivativeOrder from a value of type
// *DerivativeOrderResponseBody.
func unmarshalDerivativeOrderResponseBodyToDerivativesapiDerivativeOrder(v *DerivativeOrderResponseBody) *derivativesapi.DerivativeOrder {
	res := &derivativesapi.DerivativeOrder{
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

// unmarshalSDAValidationErrorResponseBodyToDerivativesapiSDAValidationError
// builds a value of type *derivativesapi.SDAValidationError from a value of
// type *SDAValidationErrorResponseBody.
func unmarshalSDAValidationErrorResponseBodyToDerivativesapiSDAValidationError(v *SDAValidationErrorResponseBody) *derivativesapi.SDAValidationError {
	if v == nil {
		return nil
	}
	res := &derivativesapi.SDAValidationError{
		Code:   *v.Code,
		Reason: *v.Reason,
		Field:  v.Field,
	}

	return res
}