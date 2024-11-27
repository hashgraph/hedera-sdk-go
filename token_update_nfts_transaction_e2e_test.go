//go:build all || e2e
// +build all e2e

package hiero

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestTokenUpdateNftsUpdatesMetadata(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create supply key
	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	supplyPublicKey := supplyKey.PublicKey()

	// create metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	metadataPublicKey := metadataKey.PublicKey()

	nftCount := 4
	// create metadata list for all NFTs
	mintMetadata := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := []byte{6, 9}
	// create updated metadata list  for all NFTs
	updatedMetadataList := generateNftMetadata(updatedMetadata, nftCount/2)

	// create token with metadata key
	tokenID := createTokenWithMetadataKey(t, &env, &metadataPublicKey, &supplyPublicKey)

	// mint tokens using the metadata list
	tokenMintTx, _ := NewTokenMintTransaction().
		SetMetadatas(mintMetadata).
		SetTokenID(tokenID).
		FreezeWith(env.Client)
	tx, err := tokenMintTx.Sign(supplyKey).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.GetReceipt(env.Client)
	require.NoError(t, err)
	nftSerials := receipt.SerialNumbers

	// verify the metadata is set in the new tokens
	metadataListAfterMint := getMetadataList(t, &env, tokenID, nftSerials)

	assert.True(t, reflect.DeepEqual(metadataListAfterMint, mintMetadata), "metadata after minting should match initial metadata")

	// update nft metadata for half of the NFTs
	_, err = updateNftMetadata(t, &env, tokenID, nftSerials[:nftCount/2], updatedMetadata, &metadataKey)
	require.NoError(t, err)

	// verify the metadata is updated
	metadataListAfterUpdate := getMetadataList(t, &env, tokenID, nftSerials[:nftCount/2])
	assert.True(t, reflect.DeepEqual(metadataListAfterUpdate, updatedMetadataList), "updated metadata should match expected updated metadata")

	// verify the remaining NFTs' metadata is unchanged
	metadataList := getMetadataList(t, &env, tokenID, nftSerials[nftCount/2:])

	assert.True(t, reflect.DeepEqual(metadataList, mintMetadata[nftCount/2:]), "remaining NFTs' metadata should remain unchanged")
}

func TestCanUpdateEmptyNFTMetadata(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create supply key
	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	supplyPublicKey := supplyKey.PublicKey()

	// Generate metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	metadataPublicKey := metadataKey.PublicKey()

	nftCount := 4
	mintMetadata := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := make([]byte, 0)
	updatedMetadataList := make([][]byte, 4)

	// Create a token with a metadata key
	tokenID := createTokenWithMetadataKey(t, &env, &metadataPublicKey, &supplyPublicKey)

	// Mint tokens
	tokenMintTx, _ := NewTokenMintTransaction().
		SetMetadatas(mintMetadata).
		SetTokenID(tokenID).
		FreezeWith(env.Client)
	tx, err := tokenMintTx.Sign(supplyKey).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.GetReceipt(env.Client)
	require.NoError(t, err)
	nftSerials := receipt.SerialNumbers

	metadataListAfterMint := getMetadataList(t, &env, tokenID, nftSerials)
	assert.Equal(t, mintMetadata, metadataListAfterMint, "metadata after minting should match initial metadata")

	// Update metadata for all NFTs
	_, err = updateNftMetadata(t, &env, tokenID, nftSerials, updatedMetadata, &metadataKey)
	require.NoError(t, err)

	// Verify updated metadata for NFTs
	metadataListAfterUpdate := getMetadataList(t, &env, tokenID, nftSerials)
	assert.Equal(t, updatedMetadataList, metadataListAfterUpdate, "updated metadata should match expected updated metadata")
}

func TestCannotUpdateNFTMetadataWhenKeyIsNotSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create supply key
	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	supplyPublicKey := supplyKey.PublicKey()

	// Generate metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nftCount := 4
	mintMetadata := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := []byte{6, 9}

	// Create a token without a metadata key
	tokenID := createTokenWithMetadataKey(t, &env, nil, &supplyPublicKey)

	// Mint tokens
	tokenMintTx, _ := NewTokenMintTransaction().
		SetMetadatas(mintMetadata).
		SetTokenID(tokenID).
		FreezeWith(env.Client)
	tx, err := tokenMintTx.Sign(supplyKey).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.GetReceipt(env.Client)
	require.NoError(t, err)
	nftSerials := receipt.SerialNumbers

	metadataListAfterMint := getMetadataList(t, &env, tokenID, nftSerials)
	assert.Equal(t, mintMetadata, metadataListAfterMint, "metadata after minting should match initial metadata")

	// Update metadata for the first half of the NFTs
	_, err = updateNftMetadata(t, &env, tokenID, nftSerials[:nftCount/2], updatedMetadata, &metadataKey)
	require.Error(t, err)
}

func TestCannotUpdateNFTMetadataWhenTransactionIsNotSignedWithMetadataKey(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create supply key
	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	supplyPublicKey := supplyKey.PublicKey()

	// Generate metadata key
	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	metadataPublicKey := metadataKey.PublicKey()

	nftCount := 4
	mintMetadata := generateNftMetadata([]byte{4, 2, 0}, nftCount)
	updatedMetadata := []byte{6, 9}

	// Create a token with a metadata key
	tokenID := createTokenWithMetadataKey(t, &env, &metadataPublicKey, &supplyPublicKey)

	// Mint tokens
	tokenMintTx, _ := NewTokenMintTransaction().
		SetMetadatas(mintMetadata).
		SetTokenID(tokenID).
		FreezeWith(env.Client)
	tx, err := tokenMintTx.Sign(supplyKey).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.GetReceipt(env.Client)
	require.NoError(t, err)
	nftSerials := receipt.SerialNumbers

	// Assert this will fail
	_, err = updateNftMetadata(t, &env, tokenID, nftSerials, updatedMetadata, &env.OperatorKey)
	require.Error(t, err)

	_, err = updateNftMetadata(t, &env, tokenID, nftSerials, updatedMetadata, nil)
	require.Error(t, err)
}

// Utility functions
func createTokenWithMetadataKey(t *testing.T, env *IntegrationTestEnv, metadataKey *PublicKey, supplyKey *PublicKey) TokenID {
	var tokenCreateTx *TokenCreateTransaction
	if metadataKey == nil {
		tokenCreateTx = NewTokenCreateTransaction().
			SetNodeAccountIDs(env.NodeAccountIDs).
			SetTokenName("ffff").
			SetTokenSymbol("F").
			SetTokenType(TokenTypeNonFungibleUnique).
			SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
			SetAdminKey(env.Client.GetOperatorPublicKey()).
			SetSupplyKey(supplyKey)
	} else {
		tokenCreateTx = NewTokenCreateTransaction().
			SetNodeAccountIDs(env.NodeAccountIDs).
			SetTokenName("ffff").
			SetTokenSymbol("F").
			SetTokenType(TokenTypeNonFungibleUnique).
			SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
			SetAdminKey(env.Client.GetOperatorPublicKey()).
			SetSupplyKey(supplyKey).
			SetMetadataKey(metadataKey)
	}

	tx, err := tokenCreateTx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err := tx.GetReceipt(env.Client)
	require.NoError(t, err)

	return *receipt.TokenID
}

func updateNftMetadata(t *testing.T, env *IntegrationTestEnv, tokenID TokenID, serials []int64, updatedMetadata []byte, metadataKey *PrivateKey) (*TransactionReceipt, error) {
	var tokenUpdateNftsTx *TokenUpdateNfts
	if metadataKey == nil {
		tokenUpdateNftsTx = NewTokenUpdateNftsTransaction().
			SetTokenID(tokenID).
			SetSerialNumbers(serials).
			SetMetadata(updatedMetadata)
	} else {
		frozenTx, err := NewTokenUpdateNftsTransaction().
			SetTokenID(tokenID).
			SetSerialNumbers(serials).
			SetMetadata(updatedMetadata).
			FreezeWith(env.Client)
		require.NoError(t, err)

		tokenUpdateNftsTx = frozenTx.Sign(*metadataKey)
	}

	tx, err := tokenUpdateNftsTx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err := tx.GetReceipt(env.Client)

	return &receipt, err
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
