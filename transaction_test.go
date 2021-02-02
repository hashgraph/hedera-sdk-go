package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionSerializationDeserialization(t *testing.T) {
	transaction, err := newMockTransaction()
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

func TestTransactionAddSignature(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	updateBytes, err := tx.ToBytes()
	assert.NoError(t, err)

	sig1, err := newKey.SignTransaction(&tx.Transaction)
	assert.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	assert.NoError(t, err)

	switch newTx := tx2.(type) {
	case AccountDeleteTransaction:
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).Execute(client)
		assert.NoError(t, err)
	}

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
