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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountUpdateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetExpirationTime(time.Now().Add(time.Hour * 24 * 120)).
		SetKey(newKey2.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)
	tx.Sign(newKey2)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)

	require.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)

	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountUpdateTransactionNoSigning(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	_, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetKey(newKey2.PublicKey()).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, newKey.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)

	require.NoError(t, err)

	txDelete.Sign(newKey)

	resp, err = txDelete.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountUpdateTransactionAccountIDNotSet(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewAccountUpdateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: ACCOUNT_ID_DOES_NOT_EXIST", err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

// func TestAccountUpdateTransactionAddSignature_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := PrivateKeyGenerateEd25519()
//	require.NoError(t, err)
//
//	newKey2, err := GeneratePrivateKey()
//	require.NoError(t, err)
//
//	newBalance := NewHbar(2)
//
//	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)
//
//	resp, err := NewAccountCreateTransaction().
//		SetKey(newKey.PublicKey()).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetInitialBalance(newBalance).
//		Execute(env.Client)
//
//	require.NoError(t, err)
//
//	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	accountID := *receipt.AccountID
//	require.NoError(t, err)
//
//	tx, err := NewAccountUpdateTransaction().
//		SetAccountID(accountID).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetExpirationTime(time.Now().Add(time.Hour * 24 * 120)).
//		SetTransactionID(TransactionIDGenerate(accountID)).
//		SetKey(newKey2.PublicKey()).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	updateBytes, err := tx.ToBytes()
//	require.NoError(t, err)
//
//	sig1, err := newKey.SignTransaction(&tx.Transaction)
//	require.NoError(t, err)
//	sig2, err := newKey2.SignTransaction(&tx.Transaction)
//	require.NoError(t, err)
//
//	tx2, err := TransactionFromBytes(updateBytes)
//	require.NoError(t, err)
//
//	switch newTx := tx2.(type) {
//	case AccountUpdateTransaction:
//		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).AddSignature(newKey2.PublicKey(), sig2).Execute(env.Client)
//		require.NoError(t, err)
//	}
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	info, err := NewAccountInfoQuery().
//		SetAccountID(accountID).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetMaxQueryPayment(NewHbar(1)).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())
//
//	txDelete, err := NewAccountDeleteTransaction().
//		SetAccountID(accountID).
//		SetTransferAccountID(env.Client.GetOperatorAccountID()).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		FreezeWith(env.Client)
//
//	require.NoError(t, err)
//
//	txDelete.Sign(newKey2)
//
//	resp, err = txDelete.Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//
//	require.NoError(t, err)
//
//	err = CloseIntegrationTestEnv(env, nil)
//	require.NoError(t, err)
//}
