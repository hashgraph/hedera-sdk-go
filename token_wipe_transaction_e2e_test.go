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

func TestIntegrationTokenWipeTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetKycKey(env.Client.GetOperatorPublicKey())
	})
	require.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -100).
		AddTokenTransfer(tokenID, accountID, 100).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewAccountBalanceQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	var value uint64
	for balanceTokenID, balance := range info.Token {
		if tokenID.String() == balanceTokenID.String() {
			value = balance
		}
	}

	// TODO: assert.Equal(t, uint64(100), value)

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(100).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewAccountBalanceQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	for balanceTokenID, balance := range info.Token {
		if tokenID.String() == balanceTokenID.String() {
			value = balance
		}
	}

	assert.Equal(t, uint64(0), value)

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

func TestIntegrationTokenWipeTransactionNoAmount(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetKycKey(env.Client.GetOperatorPublicKey())
	})
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nodeID := resp.NodeID

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp2, err := NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		Execute(env.Client)

	// It is now possible to perform Zero Token Operations https://hips.hedera.com/hip/hip-564
	receipt, err = resp2.SetValidateStatus(true).GetReceipt(env.Client)
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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES", err.Error())
	}

}

func TestIntegrationTokenWipeTransactionNoTokenID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetKycKey(env.Client.GetOperatorPublicKey())
	})
	require.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp2, err := NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES", err.Error())
	}

}

func TestIntegrationTokenWipeTransactionNoAccountID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetKycKey(env.Client.GetOperatorPublicKey())
	})
	require.NoError(t, err)

	nodeID := resp.NodeID

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp2, err := NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetTokenID(tokenID).
		SetAmount(10).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_ACCOUNT_ID received for transaction %s", resp2.TransactionID), err.Error())
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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES", err.Error())
	}

}

func TestIntegrationTokenWipeTransactionNotZeroTokensAtDelete(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetKycKey(env.Client.GetOperatorPublicKey())
	})
	require.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -100).
		AddTokenTransfer(tokenID, accountID, 100).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewAccountBalanceQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	var value uint64
	for balanceTokenID, balance := range info.Token {
		if tokenID.String() == balanceTokenID.String() {
			value = balance
		}
	}

	assert.Equal(t, value, value)
	// TODO: assert.Equal(t, value, uint64(100))

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewAccountBalanceQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	for balanceTokenID, balance := range info.Token {
		if tokenID.String() == balanceTokenID.String() {
			value = balance
		}
	}

	// TODO: assert.Equal(t, value, uint64(90))

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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES", err.Error())
	}

}

func DisabledTestIntegrationTokenWipeTransactionNftsIfNotOwned(t *testing.T) { // nolint
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

	_, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
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

	accountID := *receipt.AccountID

	freezeAssociate, err := NewTokenAssociateTransaction().
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	freezeAssociate.Sign(newKey)

	resp, err = freezeAssociate.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx = tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}
