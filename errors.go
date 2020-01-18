package hedera

import (
	"fmt"
	"reflect"
)

type ErrMaxQueryPaymentExceeded struct {
	// The cost of the query that was attempted as returned by QueryBuilder.Cost
	QueryCost Hbar
	// The limit for a single automatic query payment, set by
	// Client.SetMaxQueryPayment(int64) or QueryBuilder.SetMaxQueryPayment(uint64).
	MaxQueryPayment Hbar
	// Name of the query builder class used for output
	query string
}

func newErrorMaxQueryPaymentExceeded(builder *QueryBuilder, queryCost Hbar, maxQueryPayment Hbar) ErrMaxQueryPaymentExceeded {
	return ErrMaxQueryPaymentExceeded{
		QueryCost:       queryCost,
		MaxQueryPayment: maxQueryPayment,
		query:           reflect.TypeOf(*builder).Name(),
	}
}

func (e ErrMaxQueryPaymentExceeded) Error() string {
	return fmt.Sprintf("cost of %s (%d) without explicit payment is greater than Client.MaxQueryPayment (%d)",
		e.query,
		e.QueryCost.tinybar,
		e.MaxQueryPayment.tinybar)
}

type ErrBadKey struct {
	message string
}

func newErrBadKeyf(format string, a ...interface{}) ErrBadKey {
	return ErrBadKey{fmt.Sprintf(format, a...)}
}

func (e ErrBadKey) Error() string {
	return e.message
}

type ErrHederaNetwork struct {
	error error
}

func newErrHederaNetwork(e error) ErrHederaNetwork {
	return ErrHederaNetwork{error: e}
}

func (e ErrHederaNetwork) Error() string {
	return fmt.Sprintf("transport error occurred while accessing the Hedera network: %s", e.Error())
}

type ErrHederaStatus struct {
	Status Status
}

func newErrHederaStatus(status Status) ErrHederaStatus {
	// note [2020-01-15]: in the Java sdk the constructor of HederaStatusException checks if the status code is actually
	// exceptional and throws an invalid argument exception

	return ErrHederaStatus{
		Status: status,
	}
}

func (e ErrHederaStatus) Error() string {
	return e.Status.String()
}

// ErrLocalValidation is returned by TransactionBuilder.Build(*Client) and QueryBuilder.Execute(*Client)
// if the constructed transaction or query fails local sanity checks.
type ErrLocalValidation struct {
	message string
}

func newErrLocalValidation(message string) ErrLocalValidation {
	return ErrLocalValidation{message}
}

func (e ErrLocalValidation) Error() string {
	return e.message
}

// Note: an Out of Range error for Hbar units as provided in the other SDKs does not have a clean translation to go.
// 		 it would require all conversions and hbar constructors to return both the object and error resulting in a worst
//       api usage experience.
