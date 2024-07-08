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

var initialMetadataList = [][]byte{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}}

func createNftHelper(t *testing.T, env *IntegrationTestEnv) TokenID {
	tokenCreate, err := NewTokenCreateTransaction().
		SetTokenName("Example Collection").SetTokenSymbol("ABC").
		SetTokenType(TokenTypeNonFungibleUnique).SetDecimals(0).
		SetInitialSupply(0).SetMaxSupply(10).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).SetSupplyType(TokenSupplyTypeFinite).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetMetadataKey(env.Client.GetOperatorPublicKey()).
		Execute(env.Client)

	require.NoError(t, err)
	receipt, err := tokenCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	return *receipt.TokenID
}

func createFungibleTokenHelper(decimals uint, t *testing.T, env *IntegrationTestEnv) TokenID {
	tokenCreate, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("asdf").
		SetDecimals(decimals).
		SetInitialSupply(1_000_000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := tokenCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	return *receipt.TokenID

}

func createAccountHelper(t *testing.T, env *IntegrationTestEnv, maxAutoAssociations int32) (AccountID, PrivateKey) {
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	accountCreate1, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetMaxAutomaticTokenAssociations(maxAutoAssociations).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := accountCreate1.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	return *receipt.AccountID, newKey
}

// Limited max auto association tests
func TestLimitedMaxAutoAssociationsFungibleTokensFlow(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create token1
	tokenID1 := createFungibleTokenHelper(3, t, &env)

	// create token2
	tokenID2 := createFungibleTokenHelper(3, t, &env)

	// account create with 1 max auto associations
	receiver, _ := createAccountHelper(t, &env, 1)

	// transfer token1 to receiver account
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer token2 to the receiver should fail with NO_REMAINING_AUTOMATIC_ASSOCIATIONS
	tokenTransferTransaction2, err := NewTransferTransaction().
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction2.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "NO_REMAINING_AUTOMATIC_ASSOCIATIONS")
}

func TestLimitedMaxAutoAssociationsNFTsFlow(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create 2 NFT collections and mint 10 NFTs for each collection
	nftID1 := createNftHelper(t, &env)

	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftID2 := createNftHelper(t, &env)

	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// account create with 1 max auto associations
	receiver, _ := createAccountHelper(t, &env, 1)

	// transfer nftID1 nfts to receiver account
	tokenTransferTransaction, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, receiver).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nftID2 nft to receiver should fail with NO_REMAINING_AUTOMATIC_ASSOCIATIONS
	tokenTransferTransaction2, err := NewTransferTransaction().
		AddNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, receiver).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction2.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "NO_REMAINING_AUTOMATIC_ASSOCIATIONS")
}

func TestLimitedMaxAutoAssociationsFungibleTokensWithManualAssociate(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create token1
	tokenID1 := createFungibleTokenHelper(3, t, &env)

	// account create with 0 max auto associations
	receiver, key := createAccountHelper(t, &env, 0)

	frozenAssociateTxn, err := NewTokenAssociateTransaction().SetAccountID(receiver).AddTokenID(tokenID1).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenAssociateTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer token1 to receiver account
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 10
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(10), tokenBalance.Tokens.Get(tokenID1))
}

func TestLimitedMaxAutoAssociationsNFTsManualAssociate(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create NFT collection and mint 10
	nftID1 := createNftHelper(t, &env)

	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// account create with 1 max auto associations
	receiver, key := createAccountHelper(t, &env, 0)

	frozenAssociateTxn, err := NewTokenAssociateTransaction().SetAccountID(receiver).AddTokenID(nftID1).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenAssociateTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nftID1 nfts to receiver account
	tokenTransferTransaction, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, receiver).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

// HIP-904 Unlimited max auto association tests
func TestUnlimitedMaxAutoAssociationsExecutes(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// account create with unlimited max auto associations - verify it executes
	createAccountHelper(t, &env, -1)

	accountID, newKey := createAccountHelper(t, &env, 100)

	// update the account with unlimited max auto associations
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFungibleTokens(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create token1
	tokenID1 := createFungibleTokenHelper(3, t, &env)

	// create token2
	tokenID2 := createFungibleTokenHelper(3, t, &env)

	// account create with unlimited max auto associations
	accountID1, _ := createAccountHelper(t, &env, -1)
	// create account with 100 max auto associations
	accountID2, newKey := createAccountHelper(t, &env, 100)

	// update the account with unlimited max auto associations
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer to both receivers some token1 tokens
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID1, 1000).
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID2, 1000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer to both receivers some token2 tokens
	tokenTransferTransaction, err = NewTransferTransaction().
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID1, 1000).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID2, 1000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receivers is 1000
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID1).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID2))

	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(accountID2).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID2))
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFungibleTokensWithDecimals(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create token1
	tokenID1 := createFungibleTokenHelper(10, t, &env)
	// create token2
	tokenID2 := createFungibleTokenHelper(10, t, &env)

	// account create with unlimited max auto associations
	accountID, _ := createAccountHelper(t, &env, -1)

	// transfer some token1 and token2 tokens
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransferWithDecimals(tokenID1, env.Client.GetOperatorAccountID(), -1000, 10).
		AddTokenTransferWithDecimals(tokenID1, accountID, 1000, 10).
		AddTokenTransferWithDecimals(tokenID2, env.Client.GetOperatorAccountID(), -1000, 10).
		AddTokenTransferWithDecimals(tokenID2, accountID, 1000, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 1000
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID2))
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFromFungibleTokens(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create spender account which will be approved to spend
	spender, spenderKey := createAccountHelper(t, &env, 10)

	// create token1
	tokenID1 := createFungibleTokenHelper(3, t, &env)

	// create token2
	tokenID2 := createFungibleTokenHelper(3, t, &env)

	// account create with unlimited max auto associations
	accountID, _ := createAccountHelper(t, &env, -1)

	// approve the spender
	approve, err := NewAccountAllowanceApproveTransaction().
		AddTokenApproval(tokenID1, spender, 2000).
		AddTokenApproval(tokenID2, spender, 2000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transferFrom some token1 and token2 tokens
	tokenTransferTransactionFrozen, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID, 1000).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID, 1000).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenTransferTransaction, err := tokenTransferTransactionFrozen.Sign(spenderKey).Execute(env.Client)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 1000
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID2))
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferNFTs(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create 2 NFT collections and mint 10 NFTs for each collection
	nftID1 := createNftHelper(t, &env)

	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftID2 := createNftHelper(t, &env)

	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// account create with unlimited max auto associations
	accountID1, _ := createAccountHelper(t, &env, -1)
	accountID2, newKey := createAccountHelper(t, &env, 100)

	// account update with unlimited max auto associations
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nft1 to both receivers, 2 for each
	tokenTransferTransaction, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, accountID2).
		Execute(env.Client)

	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nft2 to both receivers, 2 for each
	tokenTransferTransaction, err = NewTransferTransaction().
		AddNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID2.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID2.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID2.Nft(serials[3]), env.OperatorID, accountID2).
		Execute(env.Client)

	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receivers is 2
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID1).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID2))

	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(accountID2).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID2))
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFromNFTs(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create spender account which will be approved to spend
	spender, spenderKey := createAccountHelper(t, &env, 10)

	// create 2 NFT collections and mint 10 NFTs for each collection
	nftID1 := createNftHelper(t, &env)

	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftID2 := createNftHelper(t, &env)

	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// account create with unlimited max auto associations
	accountID, _ := createAccountHelper(t, &env, -1)

	// approve the spender
	approve, err := NewAccountAllowanceApproveTransaction().
		AddAllTokenNftApproval(nftID1, spender).
		AddAllTokenNftApproval(nftID1, spender).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transferFrom some nft1 nfts
	tokenTransferTransactionFrozen, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, accountID).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, accountID).
		AddNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, accountID).
		AddNftTransfer(nftID2.Nft(serials[1]), env.OperatorID, accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenTransferTransaction, err := tokenTransferTransactionFrozen.Sign(spenderKey).Execute(env.Client)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 2
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID2))
}

func TestUnlimitedMaxAutoAssociationsFailsWithInvalid(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// account create with -2 and with -1000 max auto associations
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	_, err = NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetMaxAutomaticTokenAssociations(-2).
		Execute(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTO_ASSOCIATIONS")

	_, err = NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetMaxAutomaticTokenAssociations(-1000).
		Execute(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTO_ASSOCIATIONS")

	// create account with 100 max auto associations
	accountID, newKey := createAccountHelper(t, &env, 100)

	// account update with -2 max auto associations - should fail
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-2).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTO_ASSOCIATIONS")

	// account update with -1000 max auto associations - should fail
	accountUpdateFrozen, err = NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1000).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err = accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTO_ASSOCIATIONS")
}
