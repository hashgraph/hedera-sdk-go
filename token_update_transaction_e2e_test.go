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

func TestIntegrationTokenUpdateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
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

}

func TestIntegrationTokenUpdateTransactionDifferentKeys(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetKycKey(pubKeys[0])
	})
	require.NoError(t, err)

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

}

func TestIntegrationTokenUpdateTransactionNoTokenID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	resp2, err := NewTokenUpdateTransaction().
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOKEN_ID received for transaction %s", resp2.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func DisabledTestIntegrationTokenUpdateTransactionTreasury(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	tokenID, err := createNft(&env)
	require.NoError(t, err)

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

}

var newMetadata = []byte{3, 4, 5, 6}

func TestIntegrationTokenUpdateTransactionFungibleMetadata(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetTokenMetadata(initialMetadata)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, newMetadata, info.Metadata, "updated metadata did not match")

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionNFTMetadata(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetTokenMetadata(initialMetadata)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, newMetadata, info.Metadata, "updated metadata did not match")

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionMetadataImmutableFunbigleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	pubKey := metadataKey.PublicKey()

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetMetadataKey(pubKey).
			SetTokenMetadata(initialMetadata).
			SetAdminKey(nil)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")
	assert.Equalf(t, pubKey.String(), info.MetadataKey.String(), "metadata key did not match")
	assert.Equalf(t, nil, info.AdminKey, "admin key did not match")

	tx, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata(newMetadata).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := tx.Sign(metadataKey).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, newMetadata, info.Metadata, "updated metadata did not match")
}

func TestIntegrationTokenUpdateTransactionMetadataImmutableNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetMetadataKey(metadataKey.PublicKey()).
			SetAdminKey(nil).
			SetTokenMetadata(initialMetadata)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")
	assert.Equalf(t, metadataKey.PublicKey().String(), info.MetadataKey.String(), "metadata key did not match")
	assert.Equalf(t, nil, info.AdminKey, "admin key did not match")

	tx, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata(newMetadata).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := tx.Sign(metadataKey).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, newMetadata, info.Metadata, "updated metadata did not match")
}

func TestIntegrationTokenUpdateTransactionCannotUpdateMetadataFungible(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetTokenMetadata(initialMetadata)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMemo("asdf").
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionCannotUpdateMetadataNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetTokenMetadata(initialMetadata)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMemo("asdf").
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionEraseMetadataFungibleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetTokenMetadata(initialMetadata)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata([]byte{}).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, []byte(nil), info.Metadata, "metadata did not match")

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionEraseMetadataNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetTokenMetadata(initialMetadata)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, initialMetadata, info.Metadata, "metadata did not match")

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata([]byte{}).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equalf(t, []byte(nil), info.Metadata, "metadata did not match")

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionCannotUpdateMetadataWithoutKeyFungible(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	pubKey := metadataKey.PublicKey()

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetAdminKey(adminKey).
			SetMetadataKey(pubKey).
			FreezeWith(env.Client)

		transaction.Sign(adminKey)
	})
	require.NoError(t, err)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)

	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionCannotUpdateMetadataWithoutKeyNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	metadataKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	pubKey := metadataKey.PublicKey()

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetAdminKey(adminKey).
			SetSupplyKey(adminKey).
			SetMetadataKey(pubKey).
			FreezeWith(env.Client)

		transaction.Sign(adminKey)
	})
	require.NoError(t, err)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)

	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionCannotUpdateImmutableFungibleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetAdminKey(nil).
			SetMetadataKey(nil)
	})
	require.NoError(t, err)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_IMMUTABLE")
}

func TestIntegrationTokenUpdateTransactionCannotUpdateImmutableNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetAdminKey(nil).
			SetMetadataKey(nil)
	})
	require.NoError(t, err)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenMetadata(newMetadata).
		Execute(env.Client)
	assert.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_IMMUTABLE")
}
