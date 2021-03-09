package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactionID_Execute(t *testing.T) {
	txID := TransactionIDGenerate(AccountID{0, 0, 3})

	txID = txID.SetScheduled(true)

	println(txID.String())

	txID2 := TransactionIDWithNonce([]byte("bam"))

	txID2 = txID2.SetScheduled(true)

	txID2.AccountID = &AccountID{0, 0, 3}

	println(txID2.String())

}

func TestTransactionIDFromString_Execute(t *testing.T) {
	txID, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled")
	assert.NoError(t, err)

	println(txID.String())

	txID2, err := TransactionIdFromString("62616d?scheduled")
	assert.NoError(t, err)

	println(txID2.String())
}
