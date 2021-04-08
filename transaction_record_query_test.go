package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Record_Transaction(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		FreezeWith(client)
	assert.NoError(t, err)

	tx, err = tx.SignWithOperator(client)
	assert.NoError(t, err)

	resp, err := tx.Execute(client)
	assert.NoError(t, err)

	_, err = NewTransactionReceiptQuery().
		SetTransactionID(resp.TransactionID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	record, err := NewTransactionRecordQuery().
		SetTransactionID(resp.TransactionID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	accountID := *record.Receipt.AccountID
	assert.NotNil(t, accountID)

	transaction, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_Record_ReceiptPaymentZero_Transaction(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		FreezeWith(client)
	assert.NoError(t, err)

	tx, err = tx.SignWithOperator(client)
	assert.NoError(t, err)

	resp, err := tx.Execute(client)
	assert.NoError(t, err)

	_, err = NewTransactionReceiptQuery().
		SetTransactionID(resp.TransactionID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(HbarFromTinybar(0)).
		Execute(client)
	assert.NoError(t, err)

	record, err := NewTransactionRecordQuery().
		SetTransactionID(resp.TransactionID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	accountID := *record.Receipt.AccountID
	assert.NotNil(t, accountID)

	transaction, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_Record_Record_Insufficient_Transaction(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		FreezeWith(client)
	assert.NoError(t, err)

	tx, err = tx.SignWithOperator(client)
	assert.NoError(t, err)

	resp, err := tx.Execute(client)
	assert.NoError(t, err)

	receipt, err := NewTransactionReceiptQuery().
		SetTransactionID(resp.TransactionID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = NewTransactionRecordQuery().
		SetTransactionID(resp.TransactionID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(HbarFromTinybar(99999)).
		SetQueryPayment(HbarFromTinybar(1)).
		Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional receipt status INSUFFICIENT_TX_FEE"), err.Error())
	}

	accountID := receipt.AccountID
	assert.NotNil(t, accountID)

	transaction, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
