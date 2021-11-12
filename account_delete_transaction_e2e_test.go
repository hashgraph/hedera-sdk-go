//+build all e2e

package hedera

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationAccountDeleteTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransactionID(TransactionIDGenerate(*accountID)).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	tx = tx.Sign(newKey)

	assert.True(t, newKey.PublicKey().VerifyTransaction(tx.Transaction))
	assert.False(t, env.Client.GetOperatorPublicKey().VerifyTransaction(tx.Transaction))

	resp, err = tx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountDeleteTransactionNoTransferAccountID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST received for transaction %s", resp.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountDeleteTransactionNoAccountID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST received for transaction %s", resp.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountDeleteTransactionNoSigning(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := receipt.AccountID
	assert.NoError(t, err)

	acc := *accountID

	resp, err = NewAccountDeleteTransaction().
		SetAccountID(acc).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_SIGNATURE", err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
