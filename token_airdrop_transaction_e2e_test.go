//go:build all || e2e
// +build all e2e

package hedera

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

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenAirdropTransactionTransfersTokensWhenAssociated(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create fungible and NFT token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	tx, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftSerials := receipt.SerialNumbers

	// Create receiver with unlimited auto associations and receiverSig = false
	receiver, _ := createAccountHelper(t, &env, -1)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver).
		AddTokenTransfer(tokenID, receiver, 100).
		AddTokenTransfer(tokenID, env.OperatorID, -100).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = airdropTx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the receiver holds the tokens via query
	receiverAccountBalance, err := NewAccountBalanceQuery().
		SetAccountID(receiver).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), receiverAccountBalance.Tokens.Get(tokenID))
	assert.Equal(t, uint64(2), receiverAccountBalance.Tokens.Get(nftID))

	// Verify the operator does not hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000-100), operatorBalance.Tokens.Get(tokenID))
	assert.Equal(t, uint64(8), operatorBalance.Tokens.Get(nftID))
}
func TestIntegrationTokenAirdropTransactionPendingTokensWhenNotAssociated(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	// defer CloseIntegrationTestEnv(env, nil)

	// Create fungible and NFT token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	tx, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftSerials := receipt.SerialNumbers

	// Create receiver with 0 auto associations and receiverSig = false
	receiver, _ := createAccountHelper(t, &env, 0)
	require.NoError(t, err)

	// Airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver).
		AddTokenTransfer(tokenID, receiver, 100).
		AddTokenTransfer(tokenID, env.OperatorID, -100).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// verify the pending airdrop record
	assert.Equal(t, 3, len(record.PendingAirdropRecords))
	assert.Equal(t, uint64(100), record.PendingAirdropRecords[0].pendingAirdropAmount)
	assert.Nil(t, record.PendingAirdropRecords[0].pendingAirdropId.nftID)
	assert.Equal(t, tokenID, *record.PendingAirdropRecords[0].pendingAirdropId.tokenID)

	assert.Equal(t, uint64(0), record.PendingAirdropRecords[1].pendingAirdropAmount)
	assert.Nil(t, record.PendingAirdropRecords[1].pendingAirdropId.tokenID)
	assert.Equal(t, nftID.Nft(nftSerials[0]), *record.PendingAirdropRecords[1].pendingAirdropId.nftID)

	assert.Equal(t, uint64(0), record.PendingAirdropRecords[2].pendingAirdropAmount)
	assert.Nil(t, record.PendingAirdropRecords[2].pendingAirdropId.tokenID)
	assert.Equal(t, nftID.Nft(nftSerials[1]), *record.PendingAirdropRecords[2].pendingAirdropId.nftID)

	// Verify the receiver does not hold the tokens via query
	receiverAccountBalance, err := NewAccountBalanceQuery().
		SetAccountID(receiver).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), receiverAccountBalance.Tokens.Get(tokenID))
	assert.Equal(t, uint64(0), receiverAccountBalance.Tokens.Get(nftID))

	// Verify the operator does hold the tokens
	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), operatorBalance.Tokens.Get(tokenID))
	assert.Equal(t, uint64(10), operatorBalance.Tokens.Get(nftID))
}
