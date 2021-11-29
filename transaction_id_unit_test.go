//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitTransactionID(t *testing.T) {
	txID := TransactionIDGenerate(AccountID{0, 0, 3, nil})
	txID = txID.SetScheduled(true)
}

func TestUnitTransactionIDFromString(t *testing.T) {
	_, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled")
	require.NoError(t, err)
}
