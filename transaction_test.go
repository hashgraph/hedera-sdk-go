package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionSerializationDeserialization(t *testing.T) {
	transaction, err := newMockTransaction()

	assert.NoError(t, err)

	_ = transaction.GetTransactionHash()

	assert.NoError(t, err)

	txBytes, err := transaction.MarshalBinary()

	assert.NoError(t, err)

	var deserializedTX Transaction
	err = deserializedTX.UnmarshalBinary(txBytes)

	assert.NoError(t, err)

	assert.Equal(t, transaction.String(), deserializedTX.String())
}
