//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitTransactionID(t *testing.T) {
	txID := TransactionIDGenerate(AccountID{0, 0, 3, nil, nil})
	txID = txID.SetScheduled(true)
}

func TestUnitTransactionIDFromString(t *testing.T) {
	txID, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled")
	require.NoError(t, err)
	require.Equal(t, txID.AccountID.String(), "0.0.3")
	require.True(t, txID.scheduled)
}

func TestUnitTransactionIDFromStringNonce(t *testing.T) {
	txID, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled/4")
	require.NoError(t, err)
	require.Equal(t, *txID.Nonce, int32(4))
	require.Equal(t, txID.AccountID.String(), "0.0.3")
}
