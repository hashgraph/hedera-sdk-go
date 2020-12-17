package hedera

import (
	"fmt"
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

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
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
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)

	assert.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)

	assert.NoError(t, err)
}

func Test_AccountUpdate_NoSigning(t *testing.T) {
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

	_, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetKey(newKey2.PublicKey()).
		SetMaxTransactionFee(NewHbar(1)).
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
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)

	assert.NoError(t, err)

	txDelete.Sign(newKey)

	resp, err = txDelete.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_AccountUpdate_AccoundIDNotSet(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewAccountUpdateTransaction().
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_ACCOUNT_ID"), err.Error())
}
