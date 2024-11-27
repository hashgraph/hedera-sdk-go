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

func TestIntegrationTokenTransferTransactionCanExecute(t *testing.T) {
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

	_, err = NewAccountBalanceQuery().
		SetAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	// TODO assert.Contains(t, balance.Tokens.balances, tokenID.String())
	// TODO amount := balance.Tokens.balances[tokenID.String()]
	// TODO assert.Equal(t, uint64(999990), amount, "token transfer transaction failed")
	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

func TestIntegrationTokenTransferTransactionNotZeroSum(t *testing.T) {
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

	resp2, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status TRANSFERS_NOT_ZERO_SUM_FOR_TOKEN received for transaction %s", resp2.TransactionID), err.Error())
	}

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)

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

func DisabledTestIntegrationNftTransferTransaction(t *testing.T) { // nolint
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

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	metaData := [][]byte{{50}, {50}}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadatas(metaData).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
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
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[0]), env.OperatorID, accountID).
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[1]), env.OperatorID, accountID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenNftInfoQuery().
		SetNftID(tokenID.Nft(mintReceipt.SerialNumbers[0])).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID.String(), info[0].AccountID.String())

	info, err = NewTokenNftInfoQuery().
		SetNftID(tokenID.Nft(mintReceipt.SerialNumbers[1])).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID.String(), info[0].AccountID.String())

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetSerialNumbers(mintReceipt.SerialNumbers).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

func DisabledTestIntegrationTokenFeeScheduleUpdateRecursionDepthTransaction(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(10)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              100000000,
		DenominationTokenID: nil,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, accountID, -1).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), 1).
		FreezeWith(env.Client)
	require.NoError(t, err)

	transferTx.Sign(newKey)

	resp, err = transferTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_FEE_CHARGING_EXCEEDED_MAX_RECURSION_DEPTH", err.Error())
	}
}

func TestIntegrationTokenFeeScheduleUpdateHugeAmountTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

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

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1000,
		DenominationTokenID: nil,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -5).
		AddTokenTransfer(tokenID, accountID, 5).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transferTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func DisabledTestIntegrationTokenFeeScheduleUpdateHugeAmountOneTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID2 := *receipt.AccountID

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &accountID2,
		},
		Amount:              1000,
		DenominationTokenID: nil,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -5).
		AddTokenTransfer(tokenID, accountID, 5).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transferTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationNftTransferTransactionAccountAmountTransfersOnlyForFungibleCommon(t *testing.T) {
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

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: ACCOUNT_AMOUNT_TRANSFERS_ONLY_ALLOWED_FOR_FUNGIBLE_COMMON", err.Error())
	}

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_WIPING_AMOUNT", err.Error())
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

func DisabledTestIntegrationNftAccountStillOwnsNfts(t *testing.T) { // nolint
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

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	metaData := [][]byte{{50}, {50}}

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

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadatas(metaData).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
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
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[0]), env.OperatorID, accountID).
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[1]), env.OperatorID, accountID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenNftInfoQuery().
		ByNftID(tokenID.Nft(mintReceipt.SerialNumbers[0])).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID.String(), info[0].AccountID.String())

	info, err = NewTokenNftInfoQuery().
		ByNftID(tokenID.Nft(mintReceipt.SerialNumbers[1])).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID.String(), info[0].AccountID.String())

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTokenNotAssociatedToFeeCollector(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

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
		transaction.SetKycKey(env.Client.GetOperatorPublicKey()).
			SetFeeScheduleKey(env.Client.GetOperatorPublicKey())
	})
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{NewCustomFixedFee().
			SetFeeCollectorAccountID(accountID).
			SetAmount(1).
			SetDenominatingTokenToSameToken()}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_TOKEN_ID_IN_CUSTOM_FEES", err.Error())
	}
}

func DisabledTestIntegrationTokenFeeScheduleUpdateCustomFeeChargingExceededMax(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID2 := *receipt.AccountID

	resp, err = NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID3 := *receipt.AccountID

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	tokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	tokenID3, err := createFungibleToken(&env)
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

	transaction, err = NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID2).
		SetTokenIDs(tokenID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	transaction, err = NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID3).
		SetTokenIDs(tokenID3).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := NewCustomFixedFee().SetDenominatingTokenToSameToken().SetAmount(1).SetFeeCollectorAccountID(accountID)

	feeArray := make([]Fee, 0)

	for i := 0; i < 10; i++ {
		feeArray = append(feeArray, customFee)
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees(feeArray).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee2 := NewCustomFixedFee().SetDenominatingTokenToSameToken().SetAmount(1).SetFeeCollectorAccountID(accountID2)

	feeArray2 := make([]Fee, 0)

	for i := 0; i < 10; i++ {
		feeArray2 = append(feeArray2, customFee2)
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID2).
		SetCustomFees(feeArray2).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee3 := NewCustomFixedFee().SetDenominatingTokenToSameToken().SetAmount(1).SetFeeCollectorAccountID(accountID3)

	feeArray3 := make([]Fee, 0)

	for i := 0; i < 10; i++ {
		feeArray3 = append(feeArray3, customFee3)
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID3).
		SetCustomFees(feeArray3).
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

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID2).
		SetTokenID(tokenID2).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID3).
		SetTokenID(tokenID3).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID3, env.Client.GetOperatorAccountID(), -1).
		AddTokenTransfer(tokenID, accountID, 1).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -1).
		AddTokenTransfer(tokenID2, accountID2, 1).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1).
		AddTokenTransfer(tokenID3, accountID3, 1).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transferTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}
