package hedera

import (
	"fmt"
	"reflect"
)

type errMaxQueryPaymentExceeded struct {
	// The cost of the query that was attempted as returned by QueryBuilder.Cost
	QueryCost Hbar
	// The limit for a single automatic query payment, set by
	// Client.SetMaxQueryPayment(tinyBars int64) or QueryBuilder.SetMaxQueryPayment(maxPayment uint64).
	MaxQueryPayment Hbar
	// Name of the query builder class used for output
	query string
}

func newErrorMaxQueryPaymentExceeded(builder *QueryBuilder, queryCost int64, maxQueryPayment int64) errMaxQueryPaymentExceeded {
	return errMaxQueryPaymentExceeded{
		QueryCost:       HbarFromTinybar(queryCost),
		MaxQueryPayment: HbarFromTinybar(maxQueryPayment),
		query:           reflect.TypeOf(*builder).Name(),
	}
}

func (e errMaxQueryPaymentExceeded) Error() string {
	return fmt.Sprintf("cost of %s (%d) without explicit payment is greater than Client.MaxQueryPayment (%d)",
		e.query,
		e.QueryCost.tinybar,
		e.MaxQueryPayment.tinybar)
}
