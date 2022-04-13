//go:build all || e2e
// +build all e2e

package hedera

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrationMultiAppTransfer(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	txID := TransactionIDGenerate(env.OperatorID)

	transaction, err := NewTransferTransaction().
		SetTransactionID(txID).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)
	signedTxBytes, err := signingService(txBytes, env.OperatorKey)
	require.NoError(t, err)

	var signedTx TransferTransaction
	tx, err := TransactionFromBytes(signedTxBytes)
	require.NoError(t, err)

	switch t := tx.(type) {
	case TransferTransaction:
		signedTx = t
	default:
		panic("Did not receive `TransferTransaction` back from signed bytes")
	}

	for _, s := range signedTx.GetNodeAccountIDs() {
		println(s.String())
	}

	response, err := signedTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = response.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func signingService(txBytes []byte, key PrivateKey) ([]byte, error) {
	var unsignedTx TransferTransaction
	tx, err := TransactionFromBytes(txBytes)
	if err != nil {
		return txBytes, err
	}

	switch t := tx.(type) {
	case TransferTransaction:
		unsignedTx = t
	default:
		panic("Did not receive `TransferTransaction` back from signed bytes")
	}

	return unsignedTx.
		Sign(key).
		ToBytes()
}
