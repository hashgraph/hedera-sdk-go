package hedera

import "strings"

type ErrorTransactionValidation struct {
	Messages []string
	Err      error
}

func (e *ErrorTransactionValidation) Error() string {
	return "The following requirements were not met: \n" + strings.Join(e.Messages, "\n")
}

// fixme: this should probably be a struct with those functions implementat on it
type Transaction interface {
	Execute() error
	ExecuteForReceipt() (TransactionReceipt, error)
}

type TransactionBuilder interface {
	SetMaxTransactionFee(uint64) *TransactionBuilder
	SetMemo(string)
	validate() ErrorTransactionValidation
	Build() (Transaction, error)
}
