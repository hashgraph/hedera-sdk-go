package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactionID_Execute(t *testing.T) {
	txID := TransactionIDGenerate(AccountID{0, 0, 3, nil, nil})
	txID = txID.SetScheduled(true)
}

func TestTransactionIDFromString_Execute(t *testing.T) {
	_, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled")
	assert.NoError(t, err)
}
