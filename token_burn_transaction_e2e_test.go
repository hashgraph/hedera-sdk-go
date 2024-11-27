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

func TestIntegrationTokenBurnTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
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

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	resp, err = NewTokenBurnTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAmount(10).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, uint64(999990), info.TotalSupply)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTokenBurnTransactionNoAmount(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
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

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	resp2, err := NewTokenBurnTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(env.Client)
	receipt, err = resp2.SetValidateStatus(true).GetReceipt(env.Client)
	// It is now possible to perform Zero Token Operations https://hips.hedera.com/hip/hip-564
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTokenBurnTransactionNoTokenID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
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

	resp2, err := NewTokenBurnTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAmount(10).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOKEN_ID received for transaction %s", resp2.TransactionID), err.Error())
	}

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func DisabledTestIntegrationTokenBurnTransactionTreasuryMustOwnBurnedNft(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	metaData := []byte{50, 50}

	mint, err := NewTokenMintTransaction().
		SetTokenID(tokenID).
		SetMetadata(metaData).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	freezeAssociate, err := NewTokenAssociateTransaction().
		SetAccountID(*receipt.AccountID).
		SetTokenIDs(tokenID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	freezeAssociate.Sign(newKey)

	resp, err = freezeAssociate.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenBurnTransaction().
		SetSerialNumbers(mintReceipt.SerialNumbers).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func DisabledTestIntegrationTokenBurnTransactionInvalidMetadata(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	resp, err := NewTokenBurnTransaction().
		SetAmount(1).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_TOKEN_BURN_METADATA", err.Error())
	}
}
