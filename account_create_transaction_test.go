package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountCreateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetInitialBalance(newBalance).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func TestAccountCreateTransaction_FreezeModify_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
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

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(client)
	assert.NoError(t, err)

	tx = tx.SetAccountID(accountID)

	resp, err = tx.
		Sign(newKey).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("transaction is immutable; it has at least one signature or has been explicitly frozen"), err.Error())
	}
}

func Test_AccountCreate_NoKey(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewAccountCreateTransaction().
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status KEY_REQUIRED received for transaction %s", resp.TransactionID), err.Error())
	}
}

func TestAccountCreateTransactionAddSignature(t *testing.T) {
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
