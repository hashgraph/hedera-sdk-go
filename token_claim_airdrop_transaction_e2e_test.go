//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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
	receiver, receiverKey := createAccountHelper(t, &env, 0)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver).
		AddTokenTransfer(tokenID, receiver, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
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

	require.Equal(t, uint64(transferAmount), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(2), reciverBalance.Tokens.Get(nftID))

	// Verify the operator does not hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000-transferAmount), operatorBalance.Tokens.Get(tokenID))
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

	// Create receiver1 with 0 auto associations
	receiver1, receiver1Key := createAccountHelper(t, &env, 0)

	// Create receiver2 with 0 auto associations
	receiver2, receiver2Key := createAccountHelper(t, &env, 0)

	// Airdrop the tokens to both
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver1).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver1).
		AddTokenTransfer(tokenID, receiver1, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
		AddNftTransfer(nftID.Nft(nftSerials[2]), env.OperatorID, receiver2).
		AddNftTransfer(nftID.Nft(nftSerials[3]), env.OperatorID, receiver2).
		AddTokenTransfer(tokenID, receiver2, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
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

	require.Equal(t, uint64(transferAmount), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(2), reciverBalance.Tokens.Get(nftID))

	// Verify the receiver2 holds the tokens via query
	reciverBalance, err = NewAccountBalanceQuery().
		SetAccountID(receiver2).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(transferAmount), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(2), reciverBalance.Tokens.Get(nftID))

	// Verify the operator does not hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000-transferAmount*2), operatorBalance.Tokens.Get(tokenID))
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

	// Create receiver with 0 auto associations
	receiver, receiverKey := createAccountHelper(t, &env, 0)

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
		AddTokenTransfer(tokenID, receiver, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
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

	require.Equal(t, uint64(transferAmount), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(2), reciverBalance.Tokens.Get(nftID))

	// Verify the operator does not hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000-transferAmount), operatorBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(8), operatorBalance.Tokens.Get(nftID))
}

func TestIntegrationTokenClaimAirdropCannotClaimNonExistingAirdrop(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create receiver with 0 auto associations
	receiver, _ := createAccountHelper(t, &env, 0)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Claim the tokens with the operator which does not have pending airdrops
	// Fails with INVALID_SIGNATURE
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

	// Create receiver with 0 auto associations
	receiver, receiverKey := createAccountHelper(t, &env, 0)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
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
	// Fails with INVALID_PENDING_AIRDROP_ID
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
	// Fails with EMPTY_PENDING_AIRDROP_ID_LIST
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

	// Create receiver with 0 auto associations
	receiver, receiverKey := createAccountHelper(t, &env, 0)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Claim the tokens with duplicate pending airdrop token ids
	// Fails with PENDING_AIRDROP_ID_REPEATED
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

	// Create receiver with 0 auto associations
	receiver, receiverKey := createAccountHelper(t, &env, 0)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
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
	// Fails with TOKEN_IS_PAUSED
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

	// Create receiver with 0 auto associations
	receiver, receiverKey := createAccountHelper(t, &env, 0)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
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
	// Fails with TOKEN_WAS_DELETED
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

	// Create receiver with 0 auto associations
	receiver, receiverKey := createAccountHelper(t, &env, 0)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, transferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -transferAmount).
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
	// Fails with ACCOUNT_FROZEN_FOR_TOKEN
	claimTx, err := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	claimResp, err := claimTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = claimResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_FROZEN_FOR_TOKEN")
}
