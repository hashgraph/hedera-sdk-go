package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountUpdateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetKey(newKey2.PublicKey()).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)

	assert.NoError(t, err)

	tx.Sign(newKey)
	tx.Sign(newKey2)

	resp, err = tx.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs(nodeIDs).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)

	assert.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)

	assert.NoError(t, err)
}

func TestAccountUpdateTransactionNoSigning_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	_, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetKey(newKey2.PublicKey()).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs(nodeIDs).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)

	assert.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)

	assert.NoError(t, err)
}
