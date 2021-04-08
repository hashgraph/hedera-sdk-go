package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountDeleteTransaction_Execute(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
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

	accountID := receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransactionID(TransactionIDGenerate(*accountID)).
		FreezeWith(client)
	assert.NoError(t, err)

	tx = tx.Sign(newKey)

	assert.True(t, newKey.PublicKey().VerifyTransaction(tx.Transaction))
	assert.False(t, client.GetOperatorPublicKey().VerifyTransaction(tx.Transaction))

	resp, err = tx.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_AccountDelete_NoTransferAccountID(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
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

	accountID := receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST received for transaction %s", resp.TransactionID), err.Error())
	}
}

func Test_AccountDelete_NoAccountID(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
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
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST received for transaction %s", resp.TransactionID), err.Error())
	}
}

func Test_AccountDelete_NoSigning(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
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
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional receipt status INVALID_SIGNATURE"), err.Error())
	}
}
