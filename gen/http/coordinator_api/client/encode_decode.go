// Code generated by goa v3.1.1, DO NOT EDIT.
//
// CoordinatorAPI HTTP client encoders and decoders
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

	coordinatorapi "github.com/InjectiveLabs/dexterm/gen/coordinator_api"
	goahttp "goa.design/goa/v3/http"
)

// BuildConfigurationRequest instantiates a HTTP request object with method and
// path set to call the "CoordinatorAPI" service "configuration" endpoint
func (c *Client) BuildConfigurationRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ConfigurationCoordinatorAPIPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("CoordinatorAPI", "configuration", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeConfigurationRequest returns an encoder for requests sent to the
// CoordinatorAPI configuration server.
func EncodeConfigurationRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*coordinatorapi.ConfigurationPayload)
		if !ok {
			return goahttp.ErrInvalidType("CoordinatorAPI", "configuration", "*coordinatorapi.ConfigurationPayload", v)
		}
		values := req.URL.Query()
		values.Add("chainId", fmt.Sprintf("%v", p.ChainID))
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeConfigurationResponse returns a decoder for responses returned by the
// CoordinatorAPI configuration endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeConfigurationResponse may return the following errors:
//	- "validation_error" (type *coordinatorapi.CoordinatorValidationErrorResponse): http.StatusExpectationFailed
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeConfigurationResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body ConfigurationResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "configuration", err)
			}
			err = ValidateConfigurationResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "configuration", err)
			}
			res := NewConfigurationResultOK(&body)
			return res, nil
		case http.StatusExpectationFailed:
			var (
				body ConfigurationValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "configuration", err)
			}
			err = ValidateConfigurationValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "configuration", err)
			}
			return nil, NewConfigurationValidationError(&body)
		case http.StatusNotFound:
			var (
				body ConfigurationNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "configuration", err)
			}
			err = ValidateConfigurationNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "configuration", err)
			}
			return nil, NewConfigurationNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body ConfigurationInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "configuration", err)
			}
			err = ValidateConfigurationInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "configuration", err)
			}
			return nil, NewConfigurationInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("CoordinatorAPI", "configuration", resp.StatusCode, string(body))
		}
	}
}

// BuildRequestTransactionRequest instantiates a HTTP request object with
// method and path set to call the "CoordinatorAPI" service
// "request_transaction" endpoint
func (c *Client) BuildRequestTransactionRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: RequestTransactionCoordinatorAPIPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("CoordinatorAPI", "request_transaction", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeRequestTransactionRequest returns an encoder for requests sent to the
// CoordinatorAPI request_transaction server.
func EncodeRequestTransactionRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*coordinatorapi.RequestTransactionPayload)
		if !ok {
			return goahttp.ErrInvalidType("CoordinatorAPI", "request_transaction", "*coordinatorapi.RequestTransactionPayload", v)
		}
		body := NewRequestTransactionRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("CoordinatorAPI", "request_transaction", err)
		}
		return nil
	}
}

// DecodeRequestTransactionResponse returns a decoder for responses returned by
// the CoordinatorAPI request_transaction endpoint. restoreBody controls
// whether the response body should be restored after having been read.
// DecodeRequestTransactionResponse may return the following errors:
//	- "validation_error" (type *coordinatorapi.CoordinatorValidationErrorResponse): http.StatusExpectationFailed
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeRequestTransactionResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body RequestTransactionResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "request_transaction", err)
			}
			err = ValidateRequestTransactionResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "request_transaction", err)
			}
			res := NewRequestTransactionResultOK(&body)
			return res, nil
		case http.StatusExpectationFailed:
			var (
				body RequestTransactionValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "request_transaction", err)
			}
			err = ValidateRequestTransactionValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "request_transaction", err)
			}
			return nil, NewRequestTransactionValidationError(&body)
		case http.StatusNotFound:
			var (
				body RequestTransactionNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "request_transaction", err)
			}
			err = ValidateRequestTransactionNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "request_transaction", err)
			}
			return nil, NewRequestTransactionNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body RequestTransactionInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "request_transaction", err)
			}
			err = ValidateRequestTransactionInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "request_transaction", err)
			}
			return nil, NewRequestTransactionInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("CoordinatorAPI", "request_transaction", resp.StatusCode, string(body))
		}
	}
}

// BuildSoftCancelsRequest instantiates a HTTP request object with method and
// path set to call the "CoordinatorAPI" service "soft_cancels" endpoint
func (c *Client) BuildSoftCancelsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: SoftCancelsCoordinatorAPIPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("CoordinatorAPI", "soft_cancels", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeSoftCancelsRequest returns an encoder for requests sent to the
// CoordinatorAPI soft_cancels server.
func EncodeSoftCancelsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*coordinatorapi.SoftCancelsPayload)
		if !ok {
			return goahttp.ErrInvalidType("CoordinatorAPI", "soft_cancels", "*coordinatorapi.SoftCancelsPayload", v)
		}
		values := req.URL.Query()
		values.Add("chainId", fmt.Sprintf("%v", p.ChainID))
		req.URL.RawQuery = values.Encode()
		body := NewSoftCancelsRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("CoordinatorAPI", "soft_cancels", err)
		}
		return nil
	}
}

// DecodeSoftCancelsResponse returns a decoder for responses returned by the
// CoordinatorAPI soft_cancels endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeSoftCancelsResponse may return the following errors:
//	- "validation_error" (type *coordinatorapi.CoordinatorValidationErrorResponse): http.StatusExpectationFailed
//	- "not_found" (type *goa.ServiceError): http.StatusNotFound
//	- "internal" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeSoftCancelsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body SoftCancelsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "soft_cancels", err)
			}
			err = ValidateSoftCancelsResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "soft_cancels", err)
			}
			res := NewSoftCancelsResultOK(&body)
			return res, nil
		case http.StatusExpectationFailed:
			var (
				body SoftCancelsValidationErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "soft_cancels", err)
			}
			err = ValidateSoftCancelsValidationErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "soft_cancels", err)
			}
			return nil, NewSoftCancelsValidationError(&body)
		case http.StatusNotFound:
			var (
				body SoftCancelsNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "soft_cancels", err)
			}
			err = ValidateSoftCancelsNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "soft_cancels", err)
			}
			return nil, NewSoftCancelsNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body SoftCancelsInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("CoordinatorAPI", "soft_cancels", err)
			}
			err = ValidateSoftCancelsInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("CoordinatorAPI", "soft_cancels", err)
			}
			return nil, NewSoftCancelsInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("CoordinatorAPI", "soft_cancels", resp.StatusCode, string(body))
		}
	}
}

// unmarshalCoordinatorValidationErrorResponseBodyToCoordinatorapiCoordinatorValidationError
// builds a value of type *coordinatorapi.CoordinatorValidationError from a
// value of type *CoordinatorValidationErrorResponseBody.
func unmarshalCoordinatorValidationErrorResponseBodyToCoordinatorapiCoordinatorValidationError(v *CoordinatorValidationErrorResponseBody) *coordinatorapi.CoordinatorValidationError {
	if v == nil {
		return nil
	}
	res := &coordinatorapi.CoordinatorValidationError{
		Code:   *v.Code,
		Reason: *v.Reason,
		Field:  v.Field,
	}

	return res
}

// marshalCoordinatorapiSignedTransactionToSignedTransactionRequestBody builds
// a value of type *SignedTransactionRequestBody from a value of type
// *coordinatorapi.SignedTransaction.
func marshalCoordinatorapiSignedTransactionToSignedTransactionRequestBody(v *coordinatorapi.SignedTransaction) *SignedTransactionRequestBody {
	res := &SignedTransactionRequestBody{
		Salt:                  v.Salt,
		SignerAddress:         v.SignerAddress,
		Data:                  v.Data,
		ExpirationTimeSeconds: v.ExpirationTimeSeconds,
		GasPrice:              v.GasPrice,
		Signature:             v.Signature,
	}
	if v.Domain != nil {
		res.Domain = marshalCoordinatorapiExchangeDomainToExchangeDomainRequestBody(v.Domain)
	}

	return res
}

// marshalCoordinatorapiExchangeDomainToExchangeDomainRequestBody builds a
// value of type *ExchangeDomainRequestBody from a value of type
// *coordinatorapi.ExchangeDomain.
func marshalCoordinatorapiExchangeDomainToExchangeDomainRequestBody(v *coordinatorapi.ExchangeDomain) *ExchangeDomainRequestBody {
	res := &ExchangeDomainRequestBody{
		VerifyingContract: v.VerifyingContract,
		ChainID:           v.ChainID,
	}

	return res
}

// marshalSignedTransactionRequestBodyToCoordinatorapiSignedTransaction builds
// a value of type *coordinatorapi.SignedTransaction from a value of type
// *SignedTransactionRequestBody.
func marshalSignedTransactionRequestBodyToCoordinatorapiSignedTransaction(v *SignedTransactionRequestBody) *coordinatorapi.SignedTransaction {
	res := &coordinatorapi.SignedTransaction{
		Salt:                  v.Salt,
		SignerAddress:         v.SignerAddress,
		Data:                  v.Data,
		ExpirationTimeSeconds: v.ExpirationTimeSeconds,
		GasPrice:              v.GasPrice,
		Signature:             v.Signature,
	}
	if v.Domain != nil {
		res.Domain = marshalExchangeDomainRequestBodyToCoordinatorapiExchangeDomain(v.Domain)
	}

	return res
}

// marshalExchangeDomainRequestBodyToCoordinatorapiExchangeDomain builds a
// value of type *coordinatorapi.ExchangeDomain from a value of type
// *ExchangeDomainRequestBody.
func marshalExchangeDomainRequestBodyToCoordinatorapiExchangeDomain(v *ExchangeDomainRequestBody) *coordinatorapi.ExchangeDomain {
	res := &coordinatorapi.ExchangeDomain{
		VerifyingContract: v.VerifyingContract,
		ChainID:           v.ChainID,
	}

	return res
}

// unmarshalFillSignaturesResponseBodyToCoordinatorapiFillSignatures builds a
// value of type *coordinatorapi.FillSignatures from a value of type
// *FillSignaturesResponseBody.
func unmarshalFillSignaturesResponseBodyToCoordinatorapiFillSignatures(v *FillSignaturesResponseBody) *coordinatorapi.FillSignatures {
	if v == nil {
		return nil
	}
	res := &coordinatorapi.FillSignatures{
		OrderHash:             *v.OrderHash,
		ExpirationTimeSeconds: *v.ExpirationTimeSeconds,
		TakerAssetFillAmount:  *v.TakerAssetFillAmount,
	}
	res.ApprovalSignatures = make([]string, len(v.ApprovalSignatures))
	for i, val := range v.ApprovalSignatures {
		res.ApprovalSignatures[i] = val
	}

	return res
}