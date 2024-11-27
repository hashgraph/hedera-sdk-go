//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountCreateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		SetMaxAutomaticTokenAssociations(100).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountCreateTransactionCanFreezeModify(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
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

}

func TestIntegrationAccountCreateTransactionNoKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status KEY_REQUIRED received for transaction %s", resp.TransactionID), err.Error())
	}

}

func TestIntegrationAccountCreateTransactionAddSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)
	updateBytes, err := tx.ToBytes()
	require.NoError(t, err)

	sig1, err := newKey.SignTransaction(tx)
	require.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	require.NoError(t, err)

	if newTx, ok := tx2.(AccountDeleteTransaction); ok {
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).Execute(env.Client)
		require.NoError(t, err)
	}

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func DisabledTestIntegrationAccountCreateTransactionSetProxyAccountID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		SetProxyAccountID(accountID).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountCreateTransactionNetwork(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
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

}

func TestIntegrationAccountCreateTransactionWithAliasFromAdminKey(t *testing.T) {
	// Tests the third row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	evmAddress := adminKey.PublicKey().ToEvmAddress()

	// Create the admin account
	_, err = NewAccountCreateTransaction().
		SetKey(adminKey).
		Execute(env.Client)
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(adminKey).
		SetAlias(evmAddress).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.NotEmpty(t, info.AccountID)
	assert.Equal(t, evmAddress, info.ContractAccountID)
	assert.Equal(t, adminKey.PublicKey(), info.Key)

}

func TestIntegrationAccountCreateTransactionWithAliasFromAdminKeyWithReceiverSigRequired(t *testing.T) {
	// Tests the fourth row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	evmAddress := adminKey.PublicKey().ToEvmAddress()

	// Create the admin account
	_, err = NewAccountCreateTransaction().
		SetKey(adminKey).
		Execute(env.Client)
	require.NoError(t, err)

	frozenTxn, err := NewAccountCreateTransaction().
		SetReceiverSignatureRequired(true).
		SetKey(adminKey).
		SetAlias(evmAddress).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozenTxn.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.NotEmpty(t, info.AccountID)
	assert.Equal(t, evmAddress, info.ContractAccountID)
	assert.Equal(t, adminKey.PublicKey(), info.Key)

}

func TestIntegrationAccountCreateTransactionWithAliasFromAdminKeyWithReceiverSigRequiredWithoutSignature(t *testing.T) {
	// Tests the fourth row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	evmAddress := adminKey.PublicKey().ToEvmAddress()

	// Create the admin account
	_, err = NewAccountCreateTransaction().
		SetKey(adminKey).
		Execute(env.Client)
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetReceiverSignatureRequired(true).
		SetKey(adminKey).
		SetAlias(evmAddress).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_SIGNATURE", err.Error())
	}

}

func TestIntegrationAccountCreateTransactionWithAlias(t *testing.T) {
	// Tests the fifth row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create the admin account
	_, err = NewAccountCreateTransaction().
		SetKey(adminKey).
		Execute(env.Client)
	require.NoError(t, err)

	key, err := PrivateKeyGenerateEcdsa()
	evmAddress := key.PublicKey().ToEvmAddress()

	tx, err := NewAccountCreateTransaction().
		SetKey(adminKey).
		SetAlias(evmAddress).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := tx.
		Sign(key).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.NotEmpty(t, info.AccountID)
	assert.Equal(t, evmAddress, info.ContractAccountID)
	assert.Equal(t, adminKey.PublicKey(), info.Key)

}

func TestIntegrationAccountCreateTransactionWithAliasWithoutSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create the admin account
	_, err = NewAccountCreateTransaction().
		SetKey(adminKey).
		Execute(env.Client)
	require.NoError(t, err)

	key, err := PrivateKeyGenerateEcdsa()
	evmAddress := key.PublicKey().ToEvmAddress()

	resp, err := NewAccountCreateTransaction().
		SetKey(adminKey).
		SetAlias(evmAddress).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_SIGNATURE", err.Error())
	}

}

func TestIntegrationAccountCreateTransactionWithAliasWithReceiverSigRequired(t *testing.T) {
	// Tests the sixth row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create the admin account
	_, err = NewAccountCreateTransaction().
		SetKey(adminKey).
		Execute(env.Client)
	require.NoError(t, err)

	key, err := PrivateKeyGenerateEcdsa()
	evmAddress := key.PublicKey().ToEvmAddress()

	frozenTxn, err := NewAccountCreateTransaction().
		SetReceiverSignatureRequired(true).
		SetKey(adminKey).
		SetAlias(evmAddress).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozenTxn.
		Sign(key).
		Sign(adminKey).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.NotEmpty(t, info.AccountID)
	assert.Equal(t, evmAddress, info.ContractAccountID)
	assert.Equal(t, adminKey.PublicKey(), info.Key)

}

func TestIntegrationAccountCreateTransactionWithAliasWithReceiverSigRequiredWithoutSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create the admin account
	_, err = NewAccountCreateTransaction().
		SetKey(adminKey).
		Execute(env.Client)
	require.NoError(t, err)

	key, err := PrivateKeyGenerateEcdsa()
	evmAddress := key.PublicKey().ToEvmAddress()

	frozenTxn, err := NewAccountCreateTransaction().
		SetReceiverSignatureRequired(true).
		SetKey(adminKey).
		SetAlias(evmAddress).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozenTxn.
		Sign(key).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_SIGNATURE", err.Error())
	}

}

func TestIntegrationSerializeTransactionWithoutNodeAccountIdDeserialiseAndExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	transactionOriginal := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance)

	require.NoError(t, err)
	resp, _ := transactionOriginal.ToBytes()

	txFromBytes, err := TransactionFromBytes(resp)
	require.NoError(t, err)

	transaction := txFromBytes.(AccountCreateTransaction)
	_, err = transaction.
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)

	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionSetStakingNodeID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		SetStakedAccountID(env.OperatorID).
		SetStakedNodeID(0).
		SetMaxAutomaticTokenAssociations(100).
		Execute(env.Client)

	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}
