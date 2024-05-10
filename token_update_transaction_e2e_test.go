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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenUpdateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, "A", info.Symbol, "token failed to update")

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionDifferentKeys(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 5)
	pubKeys := make([]PublicKey, 5)

	for i := range keys {
		newKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(pubKeys[0]).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffffc").
		SetTokenID(tokenID).
		SetTokenSymbol("K").
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(pubKeys[1]).
		SetWipeKey(pubKeys[2]).
		SetKycKey(pubKeys[3]).
		SetSupplyKey(pubKeys[4]).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, "K", info.Symbol)
	assert.Equal(t, "ffffc", info.Name)
	if info.FreezeKey != nil {
		freezeKey := info.FreezeKey
		assert.Equal(t, pubKeys[1].String(), freezeKey.String())
	}
	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionNoTokenID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp2, err := NewTokenUpdateTransaction().
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOKEN_ID received for transaction %s", resp2.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func DisabledTestIntegrationTokenUpdateTransactionTreasury(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenCreate, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetMaxSupply(20).
		SetSupplyType(TokenSupplyTypeFinite).
		SetTreasuryAccountID(accountID).
		SetAdminKey(newKey.PublicKey()).
		SetFreezeKey(newKey.PublicKey()).
		SetWipeKey(newKey.PublicKey()).
		SetKycKey(newKey.PublicKey()).
		SetSupplyKey(newKey.PublicKey()).
		SetFreezeDefault(false).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenCreate.Sign(newKey)

	resp, err = tokenCreate.Execute(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetMaxSupply(20).
		SetSupplyType(TokenSupplyTypeFinite).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID
	metaData := make([]byte, 50, 101)

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadata(metaData).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	update, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTreasuryAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	update.Sign(newKey)

	resp, err = update.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, "A", info.Symbol, "token failed to update")

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

var initialMetadata = []byte{1, 2, 3}
var newMetadata = []byte{3, 4, 5, 6}

func TestIntegrationTokenUpdateTransactionFungibleMetadata(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(TokenTypeFungibleCommon).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetTokenMetadata(initialMetadata).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, newMetadata, info.Metadata, "updated metadata did not match")

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionNFTMetadata(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetTokenMetadata(initialMetadata).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, newMetadata, info.Metadata, "updated metadata did not match")

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionMetadataImmutableFunbigleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	pubKey := metadataKey.PublicKey()

	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMetadata(initialMetadata).
		SetDecimals(3).
		SetTokenType(TokenTypeFungibleCommon).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetMetadataKey(pubKey).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")
	assert.Equalf(t, pubKey.String(), info.MetadataKey.String(), "metadata key did not match")
	assert.Equalf(t, nil, info.AdminKey, "admin key did not match")

	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.Sign(metadataKey).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, newMetadata, info.Metadata, "updated metadata did not match")
}

func TestIntegrationTokenUpdateTransactionMetadataImmutableNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMetadata(initialMetadata).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetMetadataKey(metadataKey).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")
	assert.Equalf(t, metadataKey.PublicKey().String(), info.MetadataKey.String(), "metadata key did not match")
	assert.Equalf(t, nil, info.AdminKey, "admin key did not match")

	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.Sign(metadataKey).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, newMetadata, info.Metadata, "updated metadata did not match")
}

func TestIntegrationTokenUpdateTransactionCannotUpdateMetadataFungible(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	initialMetadata := make([]byte, 50, 101)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(TokenTypeFungibleCommon).
		SetTokenMetadata(initialMetadata).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMemo("asdf").
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionCannotUpdateMetadataNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	initialMetadata := make([]byte, 50, 101)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTokenMetadata(initialMetadata).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMemo("asdf").
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionEraseMetadataFungibleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(TokenTypeFungibleCommon).
		SetTokenMetadata(initialMetadata).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata([]byte{}).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, []byte(nil), info.Metadata, "metadata did not match")

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionEraseMetadataNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTokenMetadata(initialMetadata).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata([]byte{}).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, []byte(nil), info.Metadata, "metadata did not match")

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionCannotUpdateMetadataWithoutKeyFungible(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	pubKey := metadataKey.PublicKey()

	tx, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMetadata(initialMetadata).
		SetDecimals(3).
		SetTokenType(TokenTypeFungibleCommon).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(adminKey).
		SetMetadataKey(pubKey).
		SetFreezeDefault(false).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := tx.Sign(adminKey).Execute(env.Client)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)

	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionCannotUpdateMetadataWithoutKeyNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	pubKey := metadataKey.PublicKey()

	tx, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMetadata(initialMetadata).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(adminKey).
		SetSupplyKey(adminKey).
		SetMetadataKey(pubKey).
		SetFreezeDefault(false).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := tx.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)

	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionCannotUpdateImmutableFungibleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMetadata(initialMetadata).
		SetDecimals(3).
		SetTokenType(TokenTypeFungibleCommon).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_IMMUTABLE")
}

func TestIntegrationTokenUpdateTransactionCannotUpdateImmutableNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMetadata(initialMetadata).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_IMMUTABLE")
}

const zeroKeyString = "0000000000000000000000000000000000000000000000000000000000000000"

func TestIntegrationTokenUpdateTransactionUpdateAdminKeyToEmptyKeyList(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// Create NFT with supply and admin keys
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetAdminKey(env.OperatorKey).
		SetSupplyKey(env.OperatorKey).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update admin key to empty key list
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetAdminKey(NewKeyList()).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)

	// Sign with the current admin key
	resp, err = tx.Sign(env.OperatorKey).Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTokenUpdateTransactionUpdateAdminKeyWithoutAlreadySetAdminKeyFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	zeroAdminKey, err := PrivateKeyFromString(zeroKeyString)
	require.NoError(t, err)

	// Create NFT with no admin key
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeFungibleCommon).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Cannot update token with no admin key
	resp, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetAdminKey(zeroAdminKey).
		SetKeyVerificationMode(NO_VALIDATION).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: TOKEN_IS_IMMUTABLE")
}

func TestIntegrationTokenUpdateTransactionUpdateAdminKeyWithoutSignFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	zeroAdminKey, err := PrivateKeyFromString(zeroKeyString)
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with supply and admin keys
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetSupplyKey(supplyKey).
		SetAdminKey(env.OperatorKey).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update admin key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetAdminKey(zeroAdminKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the supply key should fail
	resp, err = tx.Sign(supplyKey).Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	zeroSupplyKey, err := PrivateKeyFromString(zeroKeyString)
	require.NoError(t, err)

	validSupplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with admin and supply keys
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetSupplyKey(supplyKey).
		SetAdminKey(env.OperatorKey).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(zeroSupplyKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the supply key
	resp, err = tx.Sign(supplyKey).Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update supply key
	tx, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(validSupplyKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the admin key
	resp, err = tx.Sign(env.OperatorKey).Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyToEmptyKeylist(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	validSupplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with admin and supply keys
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetSupplyKey(supplyKey).
		SetAdminKey(env.OperatorKey).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(NewKeyList()).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the supply key
	resp, err = tx.Sign(supplyKey).Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err = NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(validSupplyKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the admin key
	resp, err = tx.Sign(env.OperatorKey).Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: TOKEN_HAS_NO_SUPPLY_KEY")
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyFullValidationWithAdminKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	validSupplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with admin and supply keys
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetSupplyKey(supplyKey).
		SetAdminKey(env.OperatorKey).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(validSupplyKey).
		SetKeyVerificationMode(FULL_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the admin key
	resp, err = tx.Sign(env.OperatorKey).Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyFullValidation(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	validSupplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with supply key
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetSupplyKey(supplyKey).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(validSupplyKey).
		SetKeyVerificationMode(FULL_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	_, err = validSupplyKey.SignTransaction(&tx.Transaction)
	assert.NoError(t, err)
	_, err = supplyKey.SignTransaction(&tx.Transaction)
	assert.NoError(t, err)

	// Sign with the supply key and new supply key
	resp, err = tx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyFullValidationFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	validSupplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with supply key
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetSupplyKey(supplyKey).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(validSupplyKey).
		SetKeyVerificationMode(FULL_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	_, err = supplyKey.SignTransaction(&tx.Transaction)
	assert.NoError(t, err)

	//Sign with the old supply key, should fail
	resp, err = tx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")

	// Update supply key
	tx2, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(validSupplyKey).
		SetKeyVerificationMode(FULL_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	_, err = validSupplyKey.SignTransaction(&tx2.Transaction)
	assert.NoError(t, err)

	//Sign with the new supply key, should fail
	resp, err = tx2.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyWithInvalidKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	var invalidKey _Ed25519PublicKey
	randomBytes := make([]byte, 32)
	keyData := [32]byte{
		0x01, 0x23, 0x45, 0x67,
		0x89, 0xab, 0xcd, 0xef,
		0xfe, 0xdc, 0xba, 0x98,
		0x76, 0x54, 0x32, 0x10,
		0x00, 0x11, 0x22, 0x33,
		0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb,
		0xcc, 0xdd, 0xee, 0xff,
	}
	randomBytes = keyData[:]
	copy(invalidKey.keyData[:], randomBytes)
	invalidSupplyKey := PublicKey{
		ed25519PublicKey: &invalidKey,
	}

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with supply key
	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetSupplyKey(supplyKey).
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(invalidSupplyKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	_, err = supplyKey.SignTransaction(&tx.Transaction)
	assert.NoError(t, err)

	//Sign with the old supply key
	_, err = tx.Execute(env.Client)
	assert.ErrorContains(t, err, "exceptional precheck status INVALID_SUPPLY_KEY")

}
