//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

const tokenCancelAirdropTransferAmount = 100

func TestIntegrationTokenCancelAirdropCanExecute(t *testing.T) {
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
	receiver, _, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver).
		AddTokenTransfer(tokenID, receiver, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Cancel the tokens with the operator
	cancelTx, err := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[1].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[2].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	cancelResp, err := cancelTx.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = cancelResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the operator does hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000), operatorBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(10), operatorBalance.Tokens.Get(nftID))
}

func TestIntegrationTokenCancelAirdropMultipleReceivers(t *testing.T) {
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
		AddTokenTransfer(tokenID, receiver1, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
		AddNftTransfer(nftID.Nft(nftSerials[2]), env.OperatorID, receiver2).
		AddNftTransfer(nftID.Nft(nftSerials[3]), env.OperatorID, receiver2).
		AddTokenTransfer(tokenID, receiver2, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Verify the txn record
	assert.Equal(t, 6, len(record.PendingAirdropRecords))

	// Cancel the tokens signing with receiver1 and receiver2
	cancelTx, err := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[1].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[2].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[3].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[4].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[5].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	cancelResp, err := cancelTx.Sign(receiver1Key).Sign(receiver2Key).Execute(env.Client)
	require.NoError(t, err)

	_, err = cancelResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the operator does hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000), operatorBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(10), operatorBalance.Tokens.Get(nftID))

}

func TestIntegrationTokenCancelAirdropMultipleAirdropTxns(t *testing.T) {
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
	receiver, _, err := createAccount(&env)
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
		AddTokenTransfer(tokenID, receiver, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
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

	// Cancel the all the tokens with the receiver
	cancelTx, err := NewTokenCancelAirdropTransaction().
		SetPendingAirdropIds(pendingAirdropIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	cancelResp, err := cancelTx.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = cancelResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the receiver does not hold the tokens via query
	reciverBalance, err := NewAccountBalanceQuery().
		SetAccountID(receiver).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(0), reciverBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(0), reciverBalance.Tokens.Get(nftID))

	// Verify the operator does hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	require.Equal(t, uint64(1_000_000), operatorBalance.Tokens.Get(tokenID))
	require.Equal(t, uint64(10), operatorBalance.Tokens.Get(nftID))
}

func TestIntegrationTokenCancelAirdropCannotCancelNonExistingAirdrop(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// Create receiver
	receiver, _, err := createAccount(&env)
	require.NoError(t, err)
	// Create random account
	randomAccount, _, err := createAccount(&env)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Cancel the tokens with the random account which has not created pending airdrops
	_, err = NewTokenCancelAirdropTransaction().
		SetTransactionID(TransactionIDGenerate(randomAccount)).
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		Execute(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTokenCancelAirdropCannotCancelAlreadyCanceledAirdrop(t *testing.T) {
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
		AddTokenTransfer(tokenID, receiver, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Cancel the tokens with the operator
	cancelTx, err := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	cancelResp, err := cancelTx.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = cancelResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Cancel the tokens with the operator again
	cancelTx, err = NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	cancelResp, err = cancelTx.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = cancelResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_PENDING_AIRDROP_ID")
}

func TestIntegrationTokenCancelAirdropCannotCancelWithEmptyPendingAirdrops(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Cancel the tokens with the operator without setting pendingAirdropIds
	_, err := NewTokenCancelAirdropTransaction().
		Execute(env.Client)
	require.ErrorContains(t, err, "EMPTY_PENDING_AIRDROP_ID_LIST")
}

func TestIntegrationTokenCancelAirdropCannotCancelWithDupblicateEntries(t *testing.T) {
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
		AddTokenTransfer(tokenID, receiver, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Cancel the tokens with duplicate pending airdrop token ids
	cancelTx, err := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = cancelTx.Sign(env.OperatorKey).Execute(env.Client)
	require.ErrorContains(t, err, "PENDING_AIRDROP_ID_REPEATED")
}

func TestIntegrationTokenCancelAirdropCanCancelWithPausedToken(t *testing.T) {
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
		AddTokenTransfer(tokenID, receiver, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Pause the token
	pauseResp, err := NewTokenPauseTransaction().SetTokenID(tokenID).Execute(env.Client)
	require.NoError(t, err)
	_, err = pauseResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Cancel the tokens with receiver
	cancelTx, err := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	cancelResp, err := cancelTx.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = cancelResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenCancelAirdropCanCancelWithDeletedToken(t *testing.T) {
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
		AddTokenTransfer(tokenID, receiver, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Delete the token
	deleteResp, err := NewTokenDeleteTransaction().SetTokenID(tokenID).Execute(env.Client)
	require.NoError(t, err)
	_, err = deleteResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Cancel the tokens with receiver
	cancelTx, err := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	cancelResp, err := cancelTx.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = cancelResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenCancelAirdropCanCancelWithFrozenToken(t *testing.T) {
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
		AddTokenTransfer(tokenID, receiver, tokenCancelAirdropTransferAmount).
		AddTokenTransfer(tokenID, env.OperatorID, -tokenCancelAirdropTransferAmount).
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
	freezeTx, err := NewTokenFreezeTransaction().SetTokenID(tokenID).SetAccountID(receiver).FreezeWith(env.Client)
	require.NoError(t, err)
	freezeResp, err := freezeTx.Sign(receiverKey).Execute(env.Client)
	_, err = freezeResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Cancel the tokens with receiver
	cancelTx, err := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	cancelResp, err := cancelTx.Sign(env.OperatorKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = cancelResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}
