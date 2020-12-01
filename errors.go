package hedera

import (
	"errors"
	"fmt"
	// "reflect"

	"google.golang.org/grpc/codes"
	status2 "google.golang.org/grpc/status"
)

type ErrMaxChunksExceeded struct {
	Chunks    uint64
	MaxChunks uint64
}

var errTransactionIsFrozen = errors.New("transaction is immutable; it has at least one signature or has been explicitly frozen")
var errNoClientOrTransactionID = errors.New("`client` must have an `operator` or `transactionId` must be set")
var errNoClientOrTransactionIDOrNodeId = errors.New("`client` must be provided or both `nodeId` and `transactionId` must be set")
var errClientOperatorSigning = errors.New("`client` must have an `operator` to sign with the operator")
var errNoClientProvided = errors.New("`client` must be provided and have an operator")
var errTransactionIsNotFrozen = errors.New("transaction is not frozen")
var errFailedToDeserializeBytes = errors.New("failed to deserialize bytes")
var errNoTransactionInBytes = errors.New("no transaction was found in bytes")
var errTransactionRequiresSingleNodeAccountID = errors.New("`PrivateKey.SignTransaction()` requires `Transaction` to have a single node `AccountID` set")
var errMaxRetryCountHit = errors.New("max retry count has been hit")

type ErrInvalidNodeAccountIDSet struct {
	NodeAccountID AccountID
}

func (err ErrInvalidNodeAccountIDSet) Error() string {
	return fmt.Sprintf("Invalid node AccountID was set for transaction: %v", err.NodeAccountID.String())
}

func (err ErrMaxChunksExceeded) Error() string {
	return fmt.Sprintf("Message requires %d chunks, but max chunks is %d", err.Chunks, err.MaxChunks)
}

// ErrMaxQueryPaymentExceeded is returned during query execution if the total cost of the query + estimated fees exceeds
// the max query payment threshold set on the client or QueryBuilder.
type ErrMaxQueryPaymentExceeded struct {
	// The cost of the query that was attempted as returned by QueryBuilder.GetCost
	QueryCost Hbar
	// The limit for a single automatic query payment, set by
	// Client.SetMaxQueryPayment(int64) or QueryBuilder.SetMaxQueryPayment(uint64).
	MaxQueryPayment Hbar
	// Name of the query transaction class used for output
	query string
}

// func newErrorMaxQueryPaymentExceeded(transaction *QueryBuilder, queryCost Hbar, maxQueryPayment Hbar) ErrMaxQueryPaymentExceeded {
// 	return ErrMaxQueryPaymentExceeded{
// 		QueryCost:       queryCost,
// 		MaxQueryPayment: maxQueryPayment,
// 		query:           reflect.TypeOf(*transaction).Name(),
// 	}
// }

// Error() implements the Error interface
func (e ErrMaxQueryPaymentExceeded) Error() string {
	return fmt.Sprintf("cost of %s (%d) without explicit payment is greater than the max query payment of %d",
		e.query,
		e.QueryCost.tinybar,
		e.MaxQueryPayment.tinybar)
}

// ErrBadKey is returned if a key is provided in an invalid format or structure
type ErrBadKey struct {
	message string
}

func newErrBadKeyf(format string, a ...interface{}) ErrBadKey {
	return ErrBadKey{fmt.Sprintf(format, a...)}
}

// Error() implements the Error interface
func (e ErrBadKey) Error() string {
	return e.message
}

// ErrHederaNetwork is returned in cases where the Hedera network cannot be reached or a network-side error occurs.
type ErrHederaNetwork struct {
	error error
	// GRPC Status Code
	StatusCode *codes.Code
}

func newErrHederaNetwork(e error) ErrHederaNetwork {
	if status, ok := status2.FromError(e); ok == true {
		statusCode := status.Code()
		return ErrHederaNetwork{error: e, StatusCode: &statusCode}
	}
	return ErrHederaNetwork{error: e}
}

// Error() implements the Error interface
func (e ErrHederaNetwork) Error() string {
	return fmt.Sprintf("transport error occurred while accessing the Hedera network: %s", e.error)
}

// ErrHederaPreCheckStatus is returned by Transaction.Execute and QueryBuilder.Execute if an exceptional status is
// returned during network side validation of the sent transaction.
type ErrHederaPreCheckStatus struct {
	TxID   TransactionID
	Status Status
}

func newErrHederaPreCheckStatus(id TransactionID, status Status) ErrHederaPreCheckStatus {
	return ErrHederaPreCheckStatus{TxID: id, Status: status}
}

// Error() implements the Error interface
func (e ErrHederaPreCheckStatus) Error() string {
	if e.TxID.AccountID.isZero() {
		return fmt.Sprintf("exceptional precheck status %s", e.Status.String())
	}
	return fmt.Sprintf("exceptional precheck status %s received for transaction %v", e.Status.String(), e.TxID)
}

// ErrHederaReceiptStatus is returned by TransactionID.GetReceipt if the status of the receipt is exceptional.
type ErrHederaReceiptStatus struct {
	TxID   TransactionID
	Status Status
}

func newErrHederaReceiptStatus(id TransactionID, status Status) ErrHederaReceiptStatus {
	return ErrHederaReceiptStatus{TxID: id, Status: status}
}

// Error() implements the Error interface
func (e ErrHederaReceiptStatus) Error() string {
	return fmt.Sprintf("exceptional status %s received for transaction %v", e.Status.String(), e.TxID)
}

// ErrHederaRecordStatus is returned by TransactionID.GetRecord if the status of the record is exceptional.
type ErrHederaRecordStatus struct {
	TxID   TransactionID
	Status Status
}

func newErrHederaRecordStatus(id TransactionID, status Status) ErrHederaRecordStatus {
	return ErrHederaRecordStatus{TxID: id, Status: status}
}

// Error() implements the Error interface
func (e ErrHederaRecordStatus) Error() string {
	return fmt.Sprintf("exceptional status %s received for transaction %v", e.Status.String(), e.TxID)
}

// ErrLocalValidation is returned by TransactionBuilder.Build(*Client) and QueryBuilder.Execute(*Client)
// if the constructed transaction or query fails local sanity checks.
type ErrLocalValidation struct {
	message string
}

func newErrLocalValidationf(format string, a ...interface{}) ErrLocalValidation {
	return ErrLocalValidation{fmt.Sprintf(format, a...)}
}

// Error() implements the Error interface
func (e ErrLocalValidation) Error() string {
	return e.message
}

// Note: an Out of Range error for Hbar units as provided in the other SDKs does not have a clean translation to go.
// 		 it would require all conversions and hbar constructors to return both the object and error resulting in a worst
//       api usage experience.

// ErrPingStatus is returned by client.Ping(AccountID) if an error occurs
type ErrPingStatus struct {
	error error
}

func newErrPingStatus(e error) ErrPingStatus {
	return ErrPingStatus{error: e}
}

// Error() implements the Error interface
func (e ErrPingStatus) Error() string {
	return fmt.Sprintf("error occured during ping")
}
