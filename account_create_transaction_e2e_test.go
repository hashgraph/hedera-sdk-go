//go:build all || e2e
// +build all e2e

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountCreateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionCanFreezeModify(t *testing.T) {
	t.Parallel()
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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionNoKey(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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

	sig1, err := newKey.SignTransaction(&tx.Transaction)
	require.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	require.NoError(t, err)

	if newTx, ok := tx2.(AccountDeleteTransaction); ok {
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).Execute(env.Client)
		require.NoError(t, err)
	}

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func DisabledTestIntegrationAccountCreateTransactionSetProxyAccountID(t *testing.T) {
	t.Parallel()
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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionNetwork(t *testing.T) {
	t.Parallel()
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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionWithAliasFromAdminKey(t *testing.T) {
	// Tests the third row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionWithAliasFromAdminKeyWithReceiverSigRequired(t *testing.T) {
	// Tests the fourth row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionWithAliasFromAdminKeyWithReceiverSigRequiredWithoutSignature(t *testing.T) {
	// Tests the fourth row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionWithAlias(t *testing.T) {
	// Tests the fifth row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionWithAliasWithoutSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionWithAliasWithReceiverSigRequired(t *testing.T) {
	// Tests the sixth row of this table
	// https://github.com/hashgraph/hedera-improvement-proposal/blob/d39f740021d7da592524cffeaf1d749803798e9a/HIP/hip-583.md#signatures
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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
		SetReceiverSignatureRequired(true).
		SetKey(adminKey).
		SetAlias(evmAddress).
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

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionWithAliasWithReceiverSigRequiredWithoutSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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
		SetReceiverSignatureRequired(true).
		SetKey(adminKey).
		SetAlias(evmAddress).
		Sign(key).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_SIGNATURE", err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
