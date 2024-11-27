//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestIntegrationTokenRejectTransactionCanExecuteForFungibleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible tokens with treasury
	tokenID1, err := createFungibleToken(&env)
	require.NoError(t, err)
	tokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer fts to the receiver
	tx, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID1, tokenID2).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 0
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID1))
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID2))

	// verify the tokens are transferred back to the treasury
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID2))
}

func TestIntegrationTokenRejectTransactionCanExecuteForNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create nft collections with treasury
	nftID1, err := createNft(&env)
	require.NoError(t, err)
	nftID2, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer nfts to the receiver
	tx, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID2.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID2.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject one of the nfts
	frozenTxn, err := NewTokenRejectTransaction().SetOwnerID(receiver).SetNftIDs(nftID1.Nft(serials[1]), nftID2.Nft(serials[1])).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance is decremented by 1
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID2))

	// verify the token is transferred back to the treasury
	nftBalance, err := NewTokenNftInfoQuery().SetNftID(nftID1.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)

	nftBalance, err = NewTokenNftInfoQuery().SetNftID(nftID2.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)
}

func TestIntegrationTokenRejectTransactionCanExecuteForFTAndNFTAtTheSameTime(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible tokens with treasury
	tokenID1, err := createFungibleToken(&env)
	require.NoError(t, err)
	tokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create nft collections with treasury
	nftID1, err := createNft(&env)
	require.NoError(t, err)
	nftID2, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer fts to the receiver
	tx1, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx1.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nfts to the receiver
	tx2, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID2.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID2.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx2.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID1, tokenID2).
		SetNftIDs(nftID1.Nft(serials[1]), nftID2.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID1))
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID2))
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID2))

	// verify the tokens are transferred back to the treasury
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID2))

	nftBalance, err := NewTokenNftInfoQuery().SetNftID(nftID1.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)

	nftBalance, err = NewTokenNftInfoQuery().SetNftID(nftID2.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)
}

func TestIntegrationTokenRejectTransactionReceiverSigRequired(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create nft with treasury with receiver sig required
	treasuryKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	accountCreateFrozen, err := NewAccountCreateTransaction().
		SetKey(treasuryKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(100).FreezeWith(env.Client)
	require.NoError(t, err)
	accountCreate, err := accountCreateFrozen.Sign(treasuryKey).Execute(env.Client)
	receipt, err := accountCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	treasury := *receipt.AccountID

	nftID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetTreasuryAccountID(treasury).FreezeWith(env.Client)
		transaction.Sign(treasuryKey)
	})
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer nft to the receiver
	frozenTransfer, err := NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), treasury, receiver).
		AddNftTransfer(nftID.Nft(serials[1]), treasury, receiver).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer, err := frozenTransfer.Sign(treasuryKey).Execute(env.Client)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance is decremented by 1
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID))

	// verify the token is transferred back to the treasury
	nftBalance, err := NewTokenNftInfoQuery().SetNftID(nftID.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, treasury, nftBalance[0].AccountID)

	// same test for fungible token

	// create fungible token with treasury with receiver sig required
	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetTreasuryAccountID(treasury).FreezeWith(env.Client)
		transaction.Sign(treasuryKey)
	})
	require.NoError(t, err)

	// transfer ft to the receiver
	frozenTransfer, err = NewTransferTransaction().
		AddTokenTransfer(tokenID, treasury, -10).
		AddTokenTransfer(tokenID, receiver, 10).
		FreezeWith(env.Client)
	transfer, err = frozenTransfer.Sign(treasuryKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err = NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 0
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID))

	// verify the tokens are transferred back to the treasury
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(treasury).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID))

}

func TestIntegrationTokenRejectTransactionTokenFrozen(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create nft with treasury
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer nft to the receiver
	tx, err := NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// freeze the token
	tokenFreeze, err := NewTokenFreezeTransaction().SetTokenID(nftID).SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	_, err = tokenFreeze.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token - should fail with ACCOUNT_FROZEN_FOR_TOKEN
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_FROZEN_FOR_TOKEN")

	// same test with fungible token

	// create fungible token with treasury
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// transfer ft to the receiver
	tx, err = NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// freeze the token
	tokenFreeze, err = NewTokenFreezeTransaction().SetTokenID(tokenID).SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	_, err = tokenFreeze.SetValidateStatus(true).GetReceipt(env.Client)

	// reject the token - should fail with ACCOUNT_FROZEN_FOR_TOKEN
	frozenTxn, err = NewTokenRejectTransaction().
		SetOwnerID(receiver).
		AddTokenID(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_FROZEN_FOR_TOKEN")

}

func TestIntegrationTokenRejectTransactionTokenPaused(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create nft with treasury
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer nft to the receiver
	tx, err := NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// pause the token
	tokenPause, err := NewTokenPauseTransaction().SetTokenID(nftID).Execute(env.Client)
	require.NoError(t, err)
	_, err = tokenPause.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token - should fail with TOKEN_IS_PAUSED
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_PAUSED")

	// same test with fungible token

	// create fungible token with treasury
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// transfer ft to the receiver
	tx, err = NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// pause the token
	tokenPause, err = NewTokenPauseTransaction().SetTokenID(tokenID).Execute(env.Client)
	require.NoError(t, err)
	_, err = tokenPause.SetValidateStatus(true).GetReceipt(env.Client)

	// reject the token - should fail with TOKEN_IS_PAUSED
	frozenTxn, err = NewTokenRejectTransaction().
		SetOwnerID(receiver).
		AddTokenID(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_PAUSED")
}

func TestIntegrationTokenRejectTransactionDoesNotRemoveAllowanceFT(t *testing.T) {
	t.Skip("Skipping test as this flow is currently not working as expected in services")
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible token with treasury
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)
	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)
	// create spender account to be approved
	spender, spenderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer ft to the receiver
	tx, err := NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// approve allowance to the spender
	frozenApprove, err := NewAccountAllowanceApproveTransaction().
		ApproveTokenAllowance(tokenID, receiver, spender, 10).FreezeWith(env.Client)
	require.NoError(t, err)
	approve, err := frozenApprove.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the spender has allowance
	env.Client.SetOperator(spender, spenderKey)
	frozenTx, err := NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID, receiver, -5, true).
		AddTokenTransfer(tokenID, spender, 5).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer, err := frozenTx.SignWith(spenderKey.PublicKey(), spenderKey.Sign).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		AddTokenID(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the allowance - should be 0 , because the receiver is no longer the owner
	env.Client.SetOperator(spender, spenderKey)
	frozenTx, err = NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID, receiver, -5, true).
		AddTokenTransfer(tokenID, spender, 5).FreezeWith(env.Client)
	tx, err = frozenTx.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INSUFFICIENT_TOKEN_BALANCE")
	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// transfer ft to the back to the receiver
	tx, err = NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the spender has allowance
	env.Client.SetOperator(spender, spenderKey)
	frozenTx, err = NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID, receiver, -5, true).
		AddTokenTransfer(tokenID, spender, 5).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer, err = frozenTx.SignWith(spenderKey.PublicKey(), spenderKey.Sign).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenRejectTransactionDoesNotRemoveAllowanceNFT(t *testing.T) {
	t.Skip("Skipping test as this flow is currently not working as expected in services")
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)
	// create spender account to be approved
	spender, spenderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)
	// create nft with treasury
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// transfer nfts to the receiver
	tx, err := NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[2]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[3]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// approve allowance to the spender
	frozenApprove, err := NewAccountAllowanceApproveTransaction().
		ApproveTokenNftAllowance(nftID.Nft(serials[0]), receiver, spender).
		ApproveTokenNftAllowance(nftID.Nft(serials[1]), receiver, spender).
		FreezeWith(env.Client)
	require.NoError(t, err)
	approve, err := frozenApprove.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the spender has allowance
	env.Client.SetOperator(spender, spenderKey)
	frozenTx, err := NewTransferTransaction().
		AddApprovedNftTransfer(nftID.Nft(serials[0]), receiver, spender, true).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer, err := frozenTx.SignWith(spenderKey.PublicKey(), spenderKey.Sign).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[1]), nftID.Nft(serials[2]), nftID.Nft(serials[3])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the allowance - should be 0 , because the receiver is no longer the owner
	env.Client.SetOperator(spender, spenderKey)
	frozenTx, err = NewTransferTransaction().
		AddApprovedNftTransfer(nftID.Nft(serials[1]), receiver, spender, true).
		FreezeWith(env.Client)
	tx, err = frozenTx.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "SPENDER_DOES_NOT_HAVE_ALLOWANCE")
	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// transfer nfts to the back to the receiver
	tx, err = NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[2]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[3]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the spender has allowance
	env.Client.SetOperator(spender, spenderKey)
	frozenTx, err = NewTransferTransaction().
		AddApprovedNftTransfer(nftID.Nft(serials[1]), receiver, spender, true).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer, err = frozenTx.SignWith(spenderKey.PublicKey(), spenderKey.Sign).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenRejectTransactionFailsWhenRejectingNFTWithTokenID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create nft with treasury
	nftID, err := createNft(&env)
	require.NoError(t, err)
	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer nfts to the receiver
	tx, err := NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[2]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the whole collection - should fail
	frozenTxn, err := NewTokenRejectTransaction().SetOwnerID(receiver).AddTokenID(nftID).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_AMOUNT_TRANSFERS_ONLY_ALLOWED_FOR_FUNGIBLE_COMMON")
}

func TestIntegrationTokenRejectTransactionFailsWithTokenReferenceRepeated(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible token with treasury
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// transfer ft to the receiver
	tx, err := NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token with duplicate token id - should fail with TOKEN_REFERENCE_REPEATED
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID, tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.Sign(key).Execute(env.Client)
	require.ErrorContains(t, err, "TOKEN_REFERENCE_REPEATED")

	// same test for nft

	// create nft with treasury
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// transfer nft to the receiver
	tx, err = NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the nft with duplicate nft id - should fail with TOKEN_REFERENCE_REPEATED
	frozenTxn, err = NewTokenRejectTransaction().SetOwnerID(receiver).SetNftIDs(nftID.Nft(serials[0]), nftID.Nft(serials[0])).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.Sign(key).Execute(env.Client)
	require.ErrorContains(t, err, "TOKEN_REFERENCE_REPEATED")
}

func TestIntegrationTokenRejectTransactionFailsWhenOwnerHasNoBalance(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible token with treasury
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)
	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// skip the transfer
	// associate the receiver
	frozenAssociate, err := NewTokenAssociateTransaction().
		SetAccountID(receiver).
		AddTokenID(tokenID).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenAssociate.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	// reject the token - should fail with INSUFFICIENT_TOKEN_BALANCE
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	tx, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INSUFFICIENT_TOKEN_BALANCE")

	// same test for nft

	// create nft with treasury
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// skip the transfer
	// associate the receiver
	frozenAssociate, err = NewTokenAssociateTransaction().
		SetAccountID(receiver).
		AddTokenID(nftID).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenAssociate.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	// reject the nft - should fail with INVALID_OWNER_ID
	frozenTxn, err = NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[0])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	tx, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_OWNER_ID")

}

func TestIntegrationTokenRejectTransactionFailsTreasuryRejects(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible token with treasury
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// skip the transfer
	// reject the token with the treasury - should fail with ACCOUNT_IS_TREASURY
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(env.OperatorID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	tx, err := frozenTxn.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_IS_TREASURY")

	// same test for nft

	// create nft with treasury
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// skip the transfer
	// reject the nft with the treasury - should fail with ACCOUNT_IS_TREASURY
	frozenTxn, err = NewTokenRejectTransaction().
		SetOwnerID(env.OperatorID).
		SetNftIDs(nftID.Nft(serials[0])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	tx, err = frozenTxn.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_IS_TREASURY")
}

func TestIntegrationTokenRejectTransactionFailsWithInvalidToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// reject the token with invalid token - should fail with EMPTY_TOKEN_REFERENCE_LIST
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(env.OperatorID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.Sign(env.OperatorKey).Execute(env.Client)
	require.ErrorContains(t, err, "EMPTY_TOKEN_REFERENCE_LIST")
}

func TestIntegrationTokenRejectTransactionFailsWithReferenceSizeExceeded(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create receiver account with auto associations
	receiver, key, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// create fungible token with treasury
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create nft with treasury
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(mintMetadata).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// transfer the tokens to the receiver
	tx, err := NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[2]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[3]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[4]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[5]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[6]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[7]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[8]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[9]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token with 11 token references - should fail with TOKEN_REFERENCE_LIST_SIZE_LIMIT_EXCEEDED
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID).
		SetNftIDs(nftID.Nft(serials[0]), nftID.Nft(serials[1]),
			nftID.Nft(serials[2]), nftID.Nft(serials[3]),
			nftID.Nft(serials[4]), nftID.Nft(serials[5]),
			nftID.Nft(serials[6]), nftID.Nft(serials[7]),
			nftID.Nft(serials[8]), nftID.Nft(serials[9])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	tx, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_REFERENCE_LIST_SIZE_LIMIT_EXCEEDED")

}

func TestIntegrationTokenRejectTransactionFailsWithInvalidSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible token with treasury
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create receiver account with auto associations
	receiver, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// craete helper key
	otherKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// transfer ft to the receiver
	tx, err := NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token with different key - should fail with INVALID_SIGNATURE
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	tx, err = frozenTxn.Sign(otherKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}
