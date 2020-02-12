package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactionSerializationDeserialization(t *testing.T) {

	transaction, err := newMockTransaction()

	assert.NoError(t, err)

	txBytes, err := transaction.Bytes()

	assert.NoError(t, err)

	deserializedTX, err := TransactionFromBytes(txBytes)

	assert.NoError(t, err)

	assert.Equal(t, transaction.String(), deserializedTX.String())
}
