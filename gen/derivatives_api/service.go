// Code generated by goa v3.1.1, DO NOT EDIT.
//
// DerivativesAPI service
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package derivativesapi

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// DerivativesAPI implements Injective Standard Derivatives API v0.
type Service interface {
	// Retrieves a list of orders given query parameters.
	Orders(context.Context, *OrdersPayload) (res *OrdersResult, err error)
	// Submit a signed derivative order to the relayer.
	PostOrder(context.Context, *PostOrderPayload) (res *PostOrderResult, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "DerivativesAPI"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"orders", "postOrder"}

// OrdersPayload is the payload type of the DerivativesAPI service orders
// method.
type OrdersPayload struct {
	// Specify page to return. Page numbering should be 1-indexed, not 0-indexed.
	Page int
	// Limit the amount of items returned per page.  If a query provides an
	// unreasonable perPage value, the API will return a validation error.
	PerPage int
	// ABIv2 encoded data marketID
	TakerAssetData *string
}

// OrdersResult is the result type of the DerivativesAPI service orders method.
type OrdersResult struct {
	// The maximum number of requests you're permitted to make per hour.
	RLimitLimit *int
	// The number of requests remaining in the current rate limit window.
	RLimitRemaining *int
	// The time at which the current rate limit window resets in UTC epoch seconds.
	RLimitReset *int
	// Total records found in collection.
	Total int
	// The page number, starts from 1.
	Page int
	// Records limit per each page.
	PerPage int
	// Derivative orders.
	Records []*DerivativeOrderRecord
}

// PostOrderPayload is the payload type of the DerivativesAPI service postOrder
// method.
type PostOrderPayload struct {
	// Specify chain ID.
	ChainID int64
	// Futures contract address?
	ExchangeAddress string
	// Address that created the order.
	MakerAddress string
	// Empty.
	TakerAddress string
	// Empty.
	FeeRecipientAddress string
	// Empty.
	SenderAddress string
	// The price of 1 contract denominated in base currency.
	MakerAssetAmount string
	// The quantity of contracts the maker seeks to obtain.
	TakerAssetAmount string
	// The direction of the contract. 1 for LONG, 2 for SHORT.
	MakerFee string
	// Empty.
	TakerFee string
	// Timestamp in seconds at which order expires.
	ExpirationTimeSeconds string
	// Arbitrary number to facilitate uniqueness of the order's hash.
	Salt string
	// The account ID of the account entering into the position. Must be an account
	// owned by makerAddress
	MakerAssetData string
	// The marketID of the market for the position
	TakerAssetData string
	// Empty.
	MakerFeeAssetData string
	// Empty.
	TakerFeeAssetData string
	// Order signature.
	Signature string
}

// PostOrderResult is the result type of the DerivativesAPI service postOrder
// method.
type PostOrderResult struct {
	// The maximum number of requests you're permitted to make per hour.
	RLimitLimit *int
	// The number of requests remaining in the current rate limit window.
	RLimitRemaining *int
	// The time at which the current rate limit window resets in UTC epoch seconds.
	RLimitReset *int
}

type DerivativeOrderRecord struct {
	// Derivative Order item.
	DerivativeOrder *DerivativeOrder
	// Additional meta data.
	MetaData map[string]string
}

// A valid signed derivative order based on the schema.
type DerivativeOrder struct {
	// Specify chain ID.
	ChainID int64
	// Futures contract address?
	ExchangeAddress string
	// Address that created the order.
	MakerAddress string
	// Empty.
	TakerAddress string
	// Empty.
	FeeRecipientAddress string
	// Empty.
	SenderAddress string
	// The price of 1 contract denominated in base currency.
	MakerAssetAmount string
	// The quantity of contracts the maker seeks to obtain.
	TakerAssetAmount string
	// The direction of the contract. 1 for LONG, 2 for SHORT.
	MakerFee string
	// Empty.
	TakerFee string
	// Timestamp in seconds at which order expires.
	ExpirationTimeSeconds string
	// Arbitrary number to facilitate uniqueness of the order's hash.
	Salt string
	// The account ID of the account entering into the position. Must be an account
	// owned by makerAddress
	MakerAssetData string
	// The marketID of the market for the position
	TakerAssetData string
	// Empty.
	MakerFeeAssetData string
	// Empty.
	TakerFeeAssetData string
	// Order signature.
	Signature string
}

// Error and description for bad requests.
type SDAValidationErrorResponse struct {
	// General error code
	Code int
	// Error reason description
	Reason string
	// A list of explained validation errors.
	ValidationErrors []*SDAValidationError
}

// Order validation error explained
type SDAValidationError struct {
	// Validation error code
	Code int
	// Validation error reason description
	Reason string
	// Field name
	Field *string
}

// Error returns an error description.
func (e *SDAValidationErrorResponse) Error() string {
	return "Error and description for bad requests."
}

// ErrorName returns "SDAValidationErrorResponse".
func (e *SDAValidationErrorResponse) ErrorName() string {
	return "validation_error"
}

// Error returns an error description.
func (e *SDAValidationError) Error() string {
	return "Order validation error explained"
}

// ErrorName returns "SDAValidationError".
func (e *SDAValidationError) ErrorName() string {
	return "SDAValidationError"
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not_found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeRateLimit builds a goa.ServiceError from an error.
func MakeRateLimit(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "rate_limit",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInternal builds a goa.ServiceError from an error.
func MakeInternal(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "internal",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNotImplemented builds a goa.ServiceError from an error.
func MakeNotImplemented(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not_implemented",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
