//go:build all || e2e
// +build all e2e

package hedera

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountCreateTransactionCanExecute(t *testing.T) {
	client, err := ClientFromConfigFile(os.Getenv("CONFIG_FILE"))
	require.NoError(t, err)
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	for i := 0; i < 100; i++ {
		resp, err := NewAccountCreateTransaction().
			SetKey(newKey).
			//SetNodeAccountIDs(env.NodeAccountIDs).
			SetInitialBalance(newBalance).
			SetMaxAutomaticTokenAssociations(100).
			Execute(client)

		require.NoError(t, err)

		receipt, err := resp.GetReceipt(client)
		require.NoError(t, err)

		accountID := *receipt.AccountID

		tx, err := NewAccountDeleteTransaction().
			SetNodeAccountIDs([]AccountID{resp.NodeID}).
			SetAccountID(accountID).
			SetTransferAccountID(client.GetOperatorAccountID()).
			SetTransactionID(TransactionIDGenerate(accountID)).
			FreezeWith(client)
		require.NoError(t, err)

		resp, err = tx.
			Sign(newKey).
			Execute(client)
		require.NoError(t, err)

		_, err = resp.GetReceipt(client)
		require.NoError(t, err)
	}

	//resp, err := NewAccountCreateTransaction().
	//	SetKey(newKey).
	//	//SetNodeAccountIDs(env.NodeAccountIDs).
	//	SetInitialBalance(newBalance).
	//	SetMaxAutomaticTokenAssociations(100).
	//	Execute(client)
	//
	//require.NoError(t, err)
	//
	//receipt, err := resp.GetReceipt(client)
	//require.NoError(t, err)
	//
	//accountID := *receipt.AccountID
	//
	//tx, err := NewAccountDeleteTransaction().
	//	SetNodeAccountIDs([]AccountID{resp.NodeID}).
	//	SetAccountID(accountID).
	//	SetTransferAccountID(client.GetOperatorAccountID()).
	//	SetTransactionID(TransactionIDGenerate(accountID)).
	//	FreezeWith(client)
	//require.NoError(t, err)
	//
	//resp, err = tx.
	//	Sign(newKey).
	//	Execute(client)
	//require.NoError(t, err)
	//
	//_, err = resp.GetReceipt(client)
	//require.NoError(t, err)

	//err = CloseIntegrationTestEnv(env, nil)
	//require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionCanFreezeModify(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx = tx.SetAccountID(accountID)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "transaction is immutable; it has at least one signature or has been explicitly frozen", err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionNoKey(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status KEY_REQUIRED received for transaction %s", resp.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionAddSignature(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	updateBytes, err := tx.ToBytes()
	require.NoError(t, err)

	sig1, err := newKey.SignTransaction(&tx.Transaction)
	require.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	require.NoError(t, err)

	if newTx, ok := tx2.(AccountDeleteTransaction); ok {
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).Execute(env.Client)
		require.NoError(t, err)
	}

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionSetProxyAccountID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		SetProxyAccountID(accountID).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID2 := *receipt.AccountID

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID2).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID.String(), info.ProxyAccountID.String())

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionNetwork(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	env.Client.SetAutoValidateChecksums(true)

	accountIDString, err := accountID.ToStringWithChecksum(ClientForMainnet())
	require.NoError(t, err)
	accountID, err = AccountIDFromString(accountIDString)
	require.NoError(t, err)

	_, err = NewAccountDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	assert.Error(t, err)

	env.Client.SetAutoValidateChecksums(false)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
