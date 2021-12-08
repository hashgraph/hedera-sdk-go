//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTransactionSerializationDeserialization(t *testing.T) {
	transaction, err := _NewMockTransaction()
	require.NoError(t, err)

	_, err = transaction.Freeze()
	require.NoError(t, err)

	_, err = transaction.GetSignatures()
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	deserializedTX, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	var deserializedTXTyped *TransferTransaction
	switch tx := deserializedTX.(type) {
	case *TransferTransaction:
		deserializedTXTyped = tx
	default:
		panic("Transaction was not TransferTransaction")
	}

	assert.Equal(t, transaction.String(), deserializedTXTyped.String())
}
