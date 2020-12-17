package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountDeleteTransaction_Execute(t *testing.T) {
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

	accountID := receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(TransactionIDGenerate(*accountID)).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_AccountDelete_NoTransferAccountID(t *testing.T) {
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

	accountID := receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST received for transaction %s", resp.TransactionID), err.Error())

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("Invalid node AccountID was set for transaction: %s", resp.NodeID), err.Error())

}

func Test_AccountDelete_NoAccountID(t *testing.T) {
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

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST received for transaction %s", resp.TransactionID), err.Error())
}

func Test_AccountDelete_NoSinging(t *testing.T) {
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

	accountID := receipt.AccountID
	assert.NoError(t, err)

	acc := *accountID

	resp, err = NewAccountDeleteTransaction().
		SetAccountID(acc).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_SIGNATURE"), err.Error())
}
