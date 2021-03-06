// Code generated by goa v3.1.1, DO NOT EDIT.
//
// CoordinatorAPI service
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package coordinatorapi

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// CoordinatorAPI implements a Standard Coordinator API v2.
type Service interface {
	// Retrieves the current coordinator configuration
	Configuration(context.Context, *ConfigurationPayload) (res *ConfigurationResult, err error)
	// Submit a signed 0x transaction encoding either a 0x fill or cancellation. If
	// the 0x transaction encodes a fill, the sender is requesting a Coordinator
	// signature required to fill the order(s) on-chain. If the 0x transaction
	// encodes an order(s) cancellation request, the sender is requesting the
	// included order(s) to be soft-cancelled by the Coordinator.
	RequestTransaction(context.Context, *RequestTransactionPayload) (res *RequestTransactionResult, err error)
	// Within the Coordinator model, the Coordinator server is the source-of-truth
	// when it comes to determining whether an order has been soft-cancelled. This
	// endpoint can be used to query whether a set of orders have been
	// soft-cancelled. The response returns the subset of orders that have been
	// soft-cancelled.
	SoftCancels(context.Context, *SoftCancelsPayload) (res *SoftCancelsResult, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "CoordinatorAPI"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [3]string{"configuration", "request_transaction", "soft_cancels"}

// ConfigurationPayload is the payload type of the CoordinatorAPI service
// configuration method.
type ConfigurationPayload struct {
	// Specify Ethereum chain ID
	ChainID int64
}

// ConfigurationResult is the result type of the CoordinatorAPI service
// configuration method.
type ConfigurationResult struct {
	// Duration of validity of coordinator approval in seconds
	ExpirationDurationSeconds uint32
	// Duration of selective delay in milliseconds
	SelectiveDelayMs uint32
	// Supported Ethereum chain IDs
	SupportedChainIds []uint32
}

// RequestTransactionPayload is the payload type of the CoordinatorAPI service
// request_transaction method.
type RequestTransactionPayload struct {
	// Signed 0x Transaction
	SignedTransaction *SignedTransaction
	// Address of Ethereum transaction signer that is allowed to execute this 0x
	// transaction
	TxOrigin string
}

// RequestTransactionResult is the result type of the CoordinatorAPI service
// request_transaction method.
type RequestTransactionResult struct {
	// when the signatures will expire and no longer be valid
	ExpirationTimeSeconds *string
	// the Coordinator signatures required to submit the 0x transaction
	Signatures []string
	// Information about the outstanding signatures to fill the order(s) that have
	// been soft-cancelled.
	OutstandingFillSignatures []*FillSignatures
	// An approval signature of the cancellation 0x transaction submitted to the
	// Coordinator (with the expiration hard-coded to 0 -- although these never
	// expire). These signatures can be used to prove that a soft-cancel was
	// granted for these order(s).
	CancellationSignatures []string
}

// SoftCancelsPayload is the payload type of the CoordinatorAPI service
// soft_cancels method.
type SoftCancelsPayload struct {
	// Specify Ethereum chain ID
	ChainID int64
	// The hashes of orders to be checked whether they can be soft-cancelled
	OrderHashes []string
}

// SoftCancelsResult is the result type of the CoordinatorAPI service
// soft_cancels method.
type SoftCancelsResult struct {
	// The subset of orders that have been soft-cancelled
	OrderHashes []string
}

type SignedTransaction struct {
	// Arbitrary number to facilitate uniqueness of the transactions's hash.
	Salt string
	// Address of transaction signer
	SignerAddress string
	// The calldata that is to be executed. This must call an Exchange contract
	// method.
	Data string
	// Timestamp in seconds at which transaction expires.
	ExpirationTimeSeconds string
	// gasPrice that transaction is required to be executed with.
	GasPrice string
	// Exchange Domain specific values.
	Domain *ExchangeDomain
	// Signature of the 0x Transaction
	Signature string
}

type ExchangeDomain struct {
	// Address of the Injective Coordinator Contract.
	VerifyingContract string
	// Ethereum Chain ID of the transaction
	ChainID string
}

type FillSignatures struct {
	// EIP712 hash of order (see LibOrder.getTypedDataHash)
	OrderHash string
	// Array of signatures that correspond to the required signatures to execute
	// each order in the transaction
	ApprovalSignatures []string
	// Timestamp in seconds at which approval expires
	ExpirationTimeSeconds string
	// Desired amount of takerAsset to sell
	TakerAssetFillAmount string
}

// Error and description for bad requests.
type CoordinatorValidationErrorResponse struct {
	// General error code
	Code int
	// Error reason description
	Reason string
	// A list of explained validation errors.
	ValidationErrors []*CoordinatorValidationError
}

// Order validation error explained
type CoordinatorValidationError struct {
	// Validation error code
	Code int
	// Validation error reason description
	Reason string
	// Field name
	Field *string
}

// Error returns an error description.
func (e *CoordinatorValidationErrorResponse) Error() string {
	return "Error and description for bad requests."
}

// ErrorName returns "CoordinatorValidationErrorResponse".
func (e *CoordinatorValidationErrorResponse) ErrorName() string {
	return "validation_error"
}

// Error returns an error description.
func (e *CoordinatorValidationError) Error() string {
	return "Order validation error explained"
}

// ErrorName returns "CoordinatorValidationError".
func (e *CoordinatorValidationError) ErrorName() string {
	return "CoordinatorValidationError"
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not_found",
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
