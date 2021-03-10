package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccountUpdateTransaction_Execute(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetExpirationTime(time.Now().Local().Add(time.Second * 5)).
		SetKey(newKey2.PublicKey()).
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
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(client)

	assert.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)

	assert.NoError(t, err)

}

func Test_AccountUpdate_NoSigning(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	_, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetKey(newKey2.PublicKey()).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(client)

	assert.NoError(t, err)

	txDelete.Sign(newKey)

	resp, err = txDelete.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_AccountUpdate_AccoundIDNotSet(t *testing.T) {
	client := newTestClient(t, false)

	resp, err := NewAccountUpdateTransaction().
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_ACCOUNT_ID"), err.Error())
	}
}

func TestAccountUpdateTransactionAddSignature_Execute(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetExpirationTime(time.Now().Local().Add(time.Second * 5)).
		SetTransactionID(TransactionIDGenerate(accountID)).
		SetKey(newKey2.PublicKey()).
		FreezeWith(client)
	assert.NoError(t, err)

	updateBytes, err := tx.ToBytes()
	assert.NoError(t, err)

	sig1, err := newKey.SignTransaction(&tx.Transaction)
	assert.NoError(t, err)
	sig2, err := newKey2.SignTransaction(&tx.Transaction)
	assert.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	assert.NoError(t, err)

	switch newTx := tx2.(type) {
	case AccountUpdateTransaction:
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).AddSignature(newKey2.PublicKey(), sig2).Execute(client)
		assert.NoError(t, err)
	}

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(client)

	assert.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)

	assert.NoError(t, err)

}
