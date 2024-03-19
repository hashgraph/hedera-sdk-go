//go:build all || e2e
// +build all e2e

package hedera

import (
	"reflect"
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

func TestTokenUpdateNftsUpdatesMetadata(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	// create metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nftCount := 4
	// create metadata list for all NFTs
	initialMetadataList := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := []byte{6, 9}
	// create updated metadata list  for all NFTs
	updatedMetadataList := generateNftMetadata(updatedMetadata, nftCount/2)

	// create token with metadata key
	tokenID := createTokenWithMetadataKey(t, &env, &metadataKey)
	nftSerials := make([]int64, 0)
	// mint tokens using the metadata list
	tokenMintTransactionReceipts := mintTokens(t, &env, tokenID, initialMetadataList)
	for _, receipt := range tokenMintTransactionReceipts {
		nftSerials = append(nftSerials, receipt.SerialNumbers[0])
	}

	// verify the metadata is set in the new tokens
	metadataListAfterMint := getMetadataList(t, &env, tokenID, nftSerials)

	assert.True(t, reflect.DeepEqual(metadataListAfterMint, initialMetadataList), "metadata after minting should match initial metadata")

	// update nft metadata for half of the NFTs
	tokenUpdateNftsTransactionReceipt, err := updateNftMetadata(t, &env, tokenID, nftSerials[:nftCount/2], updatedMetadata, &metadataKey)
	require.NoError(t, err)

	// verify the metadata is updated
	nftSerialsUpdated := tokenUpdateNftsTransactionReceipt.SerialNumbers[:nftCount/2]
	metadataListAfterUpdate := getMetadataList(t, &env, tokenID, nftSerialsUpdated)
	assert.True(t, reflect.DeepEqual(metadataListAfterUpdate, updatedMetadataList), "updated metadata should match expected updated metadata")

	// verify the remaining NFTs' metadata is unchanged
	nftSerialsSame := tokenUpdateNftsTransactionReceipt.SerialNumbers[nftCount/2:]
	metadataList := getMetadataList(t, &env, tokenID, nftSerialsSame)

	assert.True(t, reflect.DeepEqual(metadataList, initialMetadataList[nftCount/2:]), "remaining NFTs' metadata should remain unchanged")
}

func TestCanUpdateNFTMetadataAfterMetadataKeySet(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	// Generate metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nftCount := 4
	initialMetadataList := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := []byte{6, 9}
	updatedMetadataList := generateNftMetadata(updatedMetadata, nftCount/2)

	// Create a token without a metadata key
	tokenID := createTokenWithMetadataKey(t, &env, nil)

	// Update the token with a metadata key
	updateTokenMetadataKey(t, tokenID, &metadataKey, &env)

	// Mint tokens
	tokenMintTransactionReceipts := mintTokens(t, &env, tokenID, initialMetadataList)
	nftSerials := make([]int64, 0)
	for _, receipt := range tokenMintTransactionReceipts {
		nftSerials = append(nftSerials, receipt.SerialNumbers[0])
	}

	metadataListAfterMint := getMetadataList(t, &env, tokenID, nftSerials)
	assert.Equal(t, initialMetadataList, metadataListAfterMint, "metadata after minting should match initial metadata")

	// Update metadata for the first half of the NFTs
	tokenUpdateNftsTransactionReceipt, err := updateNftMetadata(t, &env, tokenID, nftSerials[:nftCount/2], updatedMetadata, &metadataKey)
	require.NoError(t, err)

	// Verify updated metadata for NFTs
	nftSerialsUpdated := tokenUpdateNftsTransactionReceipt.SerialNumbers[:nftCount/2]
	metadataListAfterUpdate := getMetadataList(t, &env, tokenID, nftSerialsUpdated)
	assert.Equal(t, updatedMetadataList, metadataListAfterUpdate, "updated metadata should match expected updated metadata")
}

func TestCannotUpdateNFTMetadataWhenKeyIsNotSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	// Generate metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nftCount := 4
	initialMetadataList := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := []byte{6, 9}

	// Create a token without a metadata key
	tokenID := createTokenWithMetadataKey(t, &env, nil)

	// Mint tokens
	tokenMintTransactionReceipts := mintTokens(t, &env, tokenID, initialMetadataList)
	nftSerials := make([]int64, 0)
	for _, receipt := range tokenMintTransactionReceipts {
		nftSerials = append(nftSerials, receipt.SerialNumbers[0])
	}

	metadataListAfterMint := getMetadataList(t, &env, tokenID, nftSerials)
	assert.Equal(t, initialMetadataList, metadataListAfterMint, "metadata after minting should match initial metadata")

	// Update metadata for the first half of the NFTs
	_, err = updateNftMetadata(t, &env, tokenID, nftSerials[:nftCount/2], updatedMetadata, &metadataKey)
	require.Error(t, err)
}

func TestCannotUpdateNFTMetadataWhenTransactionIsNotSignedWithMetadataKey(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	// Generate metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nftCount := 4
	initialMetadataList := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := []byte{6, 9}

	// Create a token with a metadata key
	tokenID := createTokenWithMetadataKey(t, &env, &metadataKey)

	// Mint tokens
	tokenMintTransactionReceipts := mintTokens(t, &env, tokenID, initialMetadataList)
	nftSerials := make([]int64, 0)
	for _, receipt := range tokenMintTransactionReceipts {
		nftSerials = append(nftSerials, receipt.SerialNumbers[0])
	}

	// Assert this will fail
	_, err = updateNftMetadata(t, &env, tokenID, nftSerials, updatedMetadata, &env.OperatorKey)
	require.Error(t, err)
}

func TestCannotUpdateNFTMetadataWhenMetadataKeyWasRemoved(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	// Generate metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nftCount := 4
	initialMetadataList := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := []byte{6, 9}

	// Create a token with a metadata key
	tokenID := createTokenWithMetadataKey(t, &env, &metadataKey)

	// Check metadata key before removal
	tokenInfoBeforeUpdate := getTokenInfo(t, &env, tokenID)
	assert.Equal(t, metadataKey.PublicKey().String(), tokenInfoBeforeUpdate._ToProtobuf().MetadataKey.String(), "metadata key should match before removal")

	// Remove metadata key
	updateTokenMetadataKey(t, tokenID, nil, &env)

	// Check metadata key after removal
	tokenInfoAfterUpdate := getTokenInfo(t, &env, tokenID)
	assert.Nil(t, tokenInfoAfterUpdate._ToProtobuf().MetadataKey, "metadata key should be nil after removal")

	// Mint tokens
	tokenMintTransactionReceipts := mintTokens(t, &env, tokenID, initialMetadataList)
	nftSerials := make([]int64, 0)
	for _, receipt := range tokenMintTransactionReceipts {
		nftSerials = append(nftSerials, receipt.SerialNumbers[0])
	}

	// Check NFTs' metadata can't be updated when a metadata key is not set
	_, err = updateNftMetadata(t, &env, tokenID, nftSerials, updatedMetadata, &metadataKey)
	require.Error(t, err)
}

// Utility functions
func createTokenWithMetadataKey(t *testing.T, env *IntegrationTestEnv, metadataKey *PrivateKey) TokenID {
	var tokenCreateTx *TokenCreateTransaction
	if metadataKey == nil {
		tokenCreateTx = NewTokenCreateTransaction().
			SetNodeAccountIDs(env.NodeAccountIDs).
			SetTokenName("ffff").
			SetTokenSymbol("F").
			SetTokenType(TokenTypeNonFungibleUnique).
			SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
			SetAdminKey(env.Client.GetOperatorPublicKey()).
			SetSupplyKey(env.Client.GetOperatorPublicKey())
	} else {
		tokenCreateTx = NewTokenCreateTransaction().
			SetNodeAccountIDs(env.NodeAccountIDs).
			SetTokenName("ffff").
			SetTokenSymbol("F").
			SetTokenType(TokenTypeNonFungibleUnique).
			SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
			SetAdminKey(env.Client.GetOperatorPublicKey()).
			SetSupplyKey(env.Client.GetOperatorPublicKey()).
			SetMetadataKey(metadataKey)
	}

	tx, err := tokenCreateTx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err := tx.GetReceipt(env.Client)
	require.NoError(t, err)

	return *receipt.TokenID
}

func updateTokenMetadataKey(t *testing.T, tokenID TokenID, metadataKey *PrivateKey, env *IntegrationTestEnv) *TokenID {
	tokenCreateTx := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMetadataKey(metadataKey)

	tx, err := tokenCreateTx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err := tx.GetReceipt(env.Client)
	require.NoError(t, err)

	return receipt.TokenID
}

func mintTokens(t *testing.T, env *IntegrationTestEnv, tokenID TokenID, metadataList [][]byte) []*TransactionReceipt {
	var receipts = make([]*TransactionReceipt, len(metadataList))
	for i, metadata := range metadataList {
		tokenMintTx := NewTokenMintTransaction().
			SetMetadata(metadata).
			SetTokenID(tokenID)

		tx, err := tokenMintTx.Execute(env.Client)
		require.NoError(t, err)

		receipt, err := tx.GetReceipt(env.Client)
		require.NoError(t, err)

		receipts[i] = &receipt
	}
	return receipts
}

func updateNftMetadata(t *testing.T, env *IntegrationTestEnv, tokenID TokenID, serials []int64, updatedMetadata []byte, metadataKey *PrivateKey) (*TransactionReceipt, error) {
	var tokenUpdateNftsTx *TokenUpdateNfts
	if metadataKey == nil {
		tokenUpdateNftsTx = NewTokenUpdateNfts().
			SetTokenID(tokenID).
			SetSerialNumbers(serials)
	} else {
		tokenUpdateNftsTx = NewTokenUpdateNfts().
			SetTokenID(tokenID).
			SetSerialNumbers(serials).
			SetMetadata(updatedMetadata).
			Sign(*metadataKey)
	}

	tx, err := tokenUpdateNftsTx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err := tx.GetReceipt(env.Client)
	require.NoError(t, err)

	return &receipt, err
}

func getTokenInfo(t *testing.T, env *IntegrationTestEnv, tokenID TokenID) TokenInfo {

	tokenInfoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).SetTokenID(tokenID)

	tokenInfo, err := tokenInfoQuery.Execute(env.Client)
	require.NoError(t, err)
	return tokenInfo
}

func getMetadataList(t *testing.T, env *IntegrationTestEnv, tokenID TokenID, nftSerials []int64) [][]byte {
	var metadataList [][]byte

	for _, serial := range nftSerials {
		nftID := NftID{
			TokenID:      tokenID,
			SerialNumber: serial,
		}

		tokenNftInfoQuery := NewTokenNftInfoQuery().
			SetNodeAccountIDs(env.NodeAccountIDs).
			SetNftID(nftID)

		nftInfo, err := tokenNftInfoQuery.Execute(env.Client)
		require.NoError(t, err)
		metadataList = append(metadataList, nftInfo[0].Metadata)
	}

	return metadataList
}

func generateNftMetadata(data []byte, count int) [][]byte {
	var metadataList [][]byte

	for i := 0; i < count; i++ {
		metadataList = append(metadataList, data)
	}

	return metadataList
}
