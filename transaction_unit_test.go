//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTransactionSerializationDeserialization(t *testing.T) {
	transaction, err := _NewMockTransaction()
	assert.NoError(t, err)

	_, err = transaction.Freeze()
	assert.NoError(t, err)

	_, err = transaction.GetSignatures()
	assert.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	assert.NoError(t, err)

	txBytes, err := transaction.ToBytes()
	assert.NoError(t, err)

	deserializedTX, err := TransactionFromBytes(txBytes)
	assert.NoError(t, err)

	var deserializedTXTyped TransferTransaction
	switch tx := deserializedTX.(type) {
	case TransferTransaction:
		deserializedTXTyped = tx
	default:
		panic("Transaction was not TransferTransaction")
	}

	assert.Equal(t, transaction.String(), deserializedTXTyped.String())
}
