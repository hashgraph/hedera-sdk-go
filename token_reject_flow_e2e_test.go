//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestCustomizer(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	newKey, _ := GeneratePrivateKey()

	_, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		_, err := transaction.SetAdminKey(newKey).FreezeWith(env.Client)
		require.NoError(t, err)
		transaction.Sign(newKey)
	})

	require.NoError(t, err)
}
func TestIntegrationTokenRejectFlowCanExecuteForFungibleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create fungible tokens with treasury
	tokenID1, err := createFungibleToken(&env)
	require.NoError(t, err)
	tokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create receiver account with 0 auto associations
	receiver, key, err := createAccount(&env)
	require.NoError(t, err)

	// associate the tokens with the receiver
	frozenAssociateTxn, err := NewTokenAssociateTransaction().SetAccountID(receiver).AddTokenID(tokenID1).AddTokenID(tokenID2).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenAssociateTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

	// reject the token + dissociate
	frozenTxn, err := NewTokenRejectFlow().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID1, tokenID2).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
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

	// verify the tokens are not associated with the receiver
	tx, err = NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_NOT_ASSOCIATED_TO_ACCOUNT")

	tx, err = NewTransferTransaction().
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_NOT_ASSOCIATED_TO_ACCOUNT")
}

func TestIntegrationTokenRejectFlowCanExecuteForNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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

	// create receiver account
	receiver, key, err := createAccount(&env)
	require.NoError(t, err)

	// associate the tokens with the receiver
	frozenAssociateTxn, err := NewTokenAssociateTransaction().SetAccountID(receiver).AddTokenID(nftID1).AddTokenID(nftID2).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenAssociateTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

	// reject the token + dissociate
	frozenTxn, err := NewTokenRejectFlow().
		SetOwnerID(receiver).
		SetNftIDs(nftID1.Nft(serials[0]), nftID1.Nft(serials[1]), nftID2.Nft(serials[0]), nftID2.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance is decremented by 2
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(0), tokenBalance.Tokens.Get(nftID2))

	// verify the token is transferred back to the treasury
	nftBalance, err := NewTokenNftInfoQuery().SetNftID(nftID1.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)

	nftBalance, err = NewTokenNftInfoQuery().SetNftID(nftID2.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)

	// verify the tokens are not associated with the receiver
	tx, err = NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_NOT_ASSOCIATED_TO_ACCOUNT")
}

func TestIntegrationTokenRejectFlowFailsWhenNotRejectingAllNFTs(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

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

	// create receiver account
	receiver, key, err := createAccount(&env)
	require.NoError(t, err)

	// associate the tokens with the receiver
	frozenAssociateTxn, err := NewTokenAssociateTransaction().SetAccountID(receiver).AddTokenID(nftID1).AddTokenID(nftID2).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenAssociateTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
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

	// reject the token + dissociate
	frozenTxn, err := NewTokenRejectFlow().
		SetOwnerID(receiver).
		SetNftIDs(nftID1.Nft(serials[0]), nftID1.Nft(serials[1]), nftID2.Nft(serials[0])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_STILL_OWNS_NFTS")
}
