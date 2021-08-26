package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIntegrationAccountUpdateTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetExpirationTime(time.Now().Add(time.Hour * 24 * 120)).
		SetKey(newKey2.PublicKey()).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	tx.Sign(newKey)
	tx.Sign(newKey2)

	resp, err = tx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)

	assert.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)

	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountUpdateTransactionNoSigning(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	_, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetKey(newKey2.PublicKey()).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, newKey.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)

	assert.NoError(t, err)

	txDelete.Sign(newKey)

	resp, err = txDelete.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountUpdateTransactionAccountIDNotSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewAccountUpdateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional receipt status: INVALID_ACCOUNT_ID"), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

//func TestAccountUpdateTransactionAddSignature_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	newKey2, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(2)
//
//	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	resp, err := NewAccountCreateTransaction().
//		SetKey(newKey.PublicKey()).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetInitialBalance(newBalance).
//		Execute(env.Client)
//
//	assert.NoError(t, err)
//
//	receipt, err := resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	accountID := *receipt.AccountID
//	assert.NoError(t, err)
//
//	tx, err := NewAccountUpdateTransaction().
//		SetAccountID(accountID).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetExpirationTime(time.Now().Add(time.Hour * 24 * 120)).
//		SetTransactionID(TransactionIDGenerate(accountID)).
//		SetKey(newKey2.PublicKey()).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	updateBytes, err := tx.ToBytes()
//	assert.NoError(t, err)
//
//	sig1, err := newKey.SignTransaction(&tx.Transaction)
//	assert.NoError(t, err)
//	sig2, err := newKey2.SignTransaction(&tx.Transaction)
//	assert.NoError(t, err)
//
//	tx2, err := TransactionFromBytes(updateBytes)
//	assert.NoError(t, err)
//
//	switch newTx := tx2.(type) {
//	case AccountUpdateTransaction:
//		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).AddSignature(newKey2.PublicKey(), sig2).Execute(env.Client)
//		assert.NoError(t, err)
//	}
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	info, err := NewAccountInfoQuery().
//		SetAccountID(accountID).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetMaxQueryPayment(NewHbar(1)).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())
//
//	txDelete, err := NewAccountDeleteTransaction().
//		SetAccountID(accountID).
//		SetTransferAccountID(env.Client.GetOperatorAccountID()).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		FreezeWith(env.Client)
//
//	assert.NoError(t, err)
//
//	txDelete.Sign(newKey2)
//
//	resp, err = txDelete.Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//
//	assert.NoError(t, err)
//
//	err = CloseIntegrationTestEnv(env, nil)
//	assert.NoError(t, err)
//}
