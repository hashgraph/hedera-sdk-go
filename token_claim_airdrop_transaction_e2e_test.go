//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

const tokenClaimAirdropTransferAmount = 100

func TestIntegrationTokenClaimAirdropCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create nft
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	txResponse, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	nftSerials := receipt.SerialNumbers

	// Create receiver with 0 auto associations
	receiver, receiverKey, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver).
		AddTokenTransfer(tokenID, receiver, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Claim the tokens with the receiver
	claimTx, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[1].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[2].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err := claimTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the receiver holds the tokens via query
	reciverBalance, err := NewAccountBalanceQuery().
		SetAccountID(receiver).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(tokenClaimAirdropTransferAmount), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(2), reciverBalance.Tokens.Get(nftID))

	// Verify the operator does not hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000-tokenClaimAirdropTransferAmount), operatorBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(8), operatorBalance.Tokens.Get(nftID))
}

func TestIntegrationTokenClaimAirdropMultipleReceivers(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create nft
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	txResponse, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	nftSerials := receipt.SerialNumbers

	// Create receiver1
	receiver1, receiver1Key, err := createAccount(&env)
	require.NoError(t, err)

	// Create receiver2
	receiver2, receiver2Key, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens to both
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver1).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver1).
		AddTokenTransfer(tokenID, receiver1, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		AddNftTransfer(nftID.Nft(nftSerials[2]), env.OperatorID, receiver2).
		AddNftTransfer(nftID.Nft(nftSerials[3]), env.OperatorID, receiver2).
		AddTokenTransfer(tokenID, receiver2, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Verify the txn record
	assert.Equal(t, 6, len(record.PendingAirdropRecords))

	// Claim the tokens signing with receiver1 and receiver2
	claimTx, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[1].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[2].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[3].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[4].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[5].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err := claimTx.Sign(receiver1Key).Sign(receiver2Key).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the receiver1 holds the tokens via query
	reciverBalance, err := NewAccountBalanceQuery().
		SetAccountID(receiver1).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(tokenClaimAirdropTransferAmount), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(2), reciverBalance.Tokens.Get(nftID))

	// Verify the receiver2 holds the tokens via query
	reciverBalance, err = NewAccountBalanceQuery().
		SetAccountID(receiver2).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(tokenClaimAirdropTransferAmount), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(2), reciverBalance.Tokens.Get(nftID))

	// Verify the operator does not hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000-tokenClaimAirdropTransferAmount*2), operatorBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(6), operatorBalance.Tokens.Get(nftID))

}

func TestIntegrationTokenClaimAirdropMultipleAirdropTxns(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create nft
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	txResponse, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	nftSerials := receipt.SerialNumbers

	// Create receiver
	receiver, receiverKey, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop some of the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver).
		Execute(env.Client)
	require.NoError(t, err)

	record1, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Airdrop some of the tokens
	airdropTx, err = NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver).
		Execute(env.Client)
	require.NoError(t, err)

	record2, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Airdrop some of the tokens
	airdropTx, err = NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record3, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Collect pending airdrop IDs into a slice
	pendingAirdrop1 := record1.PendingAirdropRecords[0].GetPendingAirdropId()
	pendingAirdrop2 := record2.PendingAirdropRecords[0].GetPendingAirdropId()
	pendingAirdrop3 := record3.PendingAirdropRecords[0].GetPendingAirdropId()
	pendingAirdropIDs := make([]*PendingAirdropId, 0)
	pendingAirdropIDs = append(pendingAirdropIDs, &pendingAirdrop1)
	pendingAirdropIDs = append(pendingAirdropIDs, &pendingAirdrop2)
	pendingAirdropIDs = append(pendingAirdropIDs, &pendingAirdrop3)

	// Claim the all the tokens with the receiver
	claimTx, err := NewTokenClaimAirdropTransaction().
		SetPendingAirdropIds(pendingAirdropIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err := claimTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the receiver holds the tokens via query
	reciverBalance, err := NewAccountBalanceQuery().
		SetAccountID(receiver).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(tokenClaimAirdropTransferAmount), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(2), reciverBalance.Tokens.Get(nftID))

	// Verify the operator does not hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000-tokenClaimAirdropTransferAmount), operatorBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(8), operatorBalance.Tokens.Get(nftID))
}

func TestIntegrationTokenClaimAirdropCannotClaimNonExistingAirdrop(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create receiver
	receiver, _, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Claim the tokens with the operator which does not have pending airdrops
	claimResp, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTokenClaimAirdropCannotClaimAlreadyClaimedAirdrop(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create receiver
	receiver, receiverKey, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Claim the tokens with the receiver
	claimTx, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err := claimTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Claim the tokens with the receiver again
	claimTx, err = NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err = claimTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_PENDING_AIRDROP_ID")
}

func TestIntegrationTokenClaimAirdropCannotClaimWithEmptyPendingAirdrops(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Claim the tokens with the receiver without setting pendingAirdropIds
	_, err := NewTokenClaimAirdropTransaction().
		Execute(env.Client)
	require.ErrorContains(t, err, "EMPTY_PENDING_AIRDROP_ID_LIST")
}

func TestIntegrationTokenClaimAirdropCannotClaimWithDupblicateEntries(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create receiver
	receiver, receiverKey, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Claim the tokens with duplicate pending airdrop token ids
	claimTx, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = claimTx.Sign(receiverKey).Execute(env.Client)
	require.ErrorContains(t, err, "PENDING_AIRDROP_ID_REPEATED")
}

func TestIntegrationTokenClaimAirdropCannotClaimWithPausedToken(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create receiver
	receiver, receiverKey, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Pause the token
	pauseResp, err := NewTokenPauseTransaction().SetTokenID(tokenID).Execute(env.Client)
	require.NoError(t, err)
	_, err = pauseResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Claim the tokens with receiver
	claimTx, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err := claimTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_PAUSED")
}

func TestIntegrationTokenClaimAirdropCannotClaimWithDeletedToken(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create receiver
	receiver, receiverKey, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Delete the token
	deleteResp, err := NewTokenDeleteTransaction().SetTokenID(tokenID).Execute(env.Client)
	require.NoError(t, err)
	_, err = deleteResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Claim the tokens with receiver
	claimTx, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err := claimTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_WAS_DELETED")
}

func TestIntegrationTokenClaimAirdropCannotClaimWithFrozenToken(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create receiver
	receiver, receiverKey, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, tokenClaimAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenClaimAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Associate the token
	associateTx, err := NewTokenAssociateTransaction().AddTokenID(tokenID).SetAccountID(receiver).FreezeWith(env.Client)
	require.NoError(t, err)
	associateResp, err := associateTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = associateResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Freeze the token
	freezeResp, err := NewTokenFreezeTransaction().SetTokenID(tokenID).SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	_, err = freezeResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Claim the tokens with receiver
	claimTx, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err := claimTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_FROZEN_FOR_TOKEN")
}
