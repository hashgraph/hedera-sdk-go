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

func createFungibleTokenHelper(t *testing.T, env *IntegrationTestEnv) TokenID {
	tokenCreate, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("asdf").
		SetDecimals(18).
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

func createAccountHelper(t *testing.T, env *IntegrationTestEnv) (AccountID, PrivateKey) {
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	accountCreate1, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetMaxAutomaticTokenAssociations(100).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := accountCreate1.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	return *receipt.AccountID, newKey
}

func TestIntegrationTokenRejectTransactionCanExecuteForFungible(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury
	tokenID := createFungibleTokenHelper(t, &env)
	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// when

	// transfer ft to the receiver
	_, err := NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		AddTokenID(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	// then

	// verify the balance of the receiver is 0
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID))

	// verify the tokens are transferred back to the treasury
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID))

	// verify the auto associations are not decremented
	accountInfo, err := NewAccountInfoQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint32(100), accountInfo.MaxAutomaticTokenAssociations)
}

func TestIntegrationTokenRejectTransactionCanExecuteForNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// given

	// create nft with treasury
	nftID := createNftHelper(t, &env)
	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers
	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// when

	// transfer nfts to the receiver
	_, err = NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[2]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)

	// reject one of the nfts
	frozenTxn, err := NewTokenRejectTransaction().SetOwnerID(receiver).AddNftID(nftID.Nft(serials[1])).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	// then

	// verify the balance is decremented by 1
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID))

	// verify the token is transferred back to the treasury
	nftBalance, err := NewTokenNftInfoQuery().SetNftID(nftID.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)

	// verify the auto associations are not decremented
	accountInfo, err := NewAccountInfoQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint32(100), accountInfo.MaxAutomaticTokenAssociations)
}

func TestIntegrationTokenRejectTransactionReceiverSigRequired(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create nft with treasury with receiver sig required
	// mint
	// create receiver account with auto associations

	// when

	// transfer nft to the receiver
	// reject the token

	// then

	// verify the balance is 0
	// verify the auto associations are not decremented

	// same test for fungible token
}

func TestIntegrationTokenRejectTransactionReceiverTokenFrozen(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create nft with treasury
	// mint
	// create receiver account with auto associations

	// when

	// transfer nft to the receiver
	// freeze the token

	// then

	// reject the token - should fail with TOKEN_IS_FROZEN

	// same test with fungible token
}

func TestIntegrationTokenRejectTransactionReceiverTokenPaused(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	//given

	// create nft with treasury
	// mint
	// create receiver account with auto associations

	// when

	// transfer nft to the receiver
	// pause the token

	// then

	// reject the token - should fail with TOKEN_IS_PAUSED

	// same test with fungible token
}

func TestIntegrationTokenRejectTransactionRemovesAllowance(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury
	// create receiver account with auto associations
	// create spender account to be approved

	// when

	// transfer ft to the receiver
	// approve allowance to the spender
	// verify the allowance with query
	// reject the token

	// then

	// verify the allowance - should be 0 , because the receiver is no longer the owner

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWithTokenReferenceRepeated(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury
	// create receiver account with auto associations

	// when

	// transfer ft to the receiver
	// reject the token

	// then

	// duplicate the reject - should fail with TOKEN_REFERENCE_REPEATED

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWhenOwnerHasNoBalance(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury
	// create receiver account with auto associations

	// when

	// skip the transfer
	// reject the token - should fail with INSUFFICIENT_TOKEN_BALANCE

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsTreasuryRejects(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury

	// when

	// skip the transfer
	// reject the token with the treasury - should fail with ACCOUNT_IS_TREASURY

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWithInvalidOwner(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury

	// when

	// reject the token with invalid owner - should fail with INVALID_OWNER_ID

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWithInvalidToken(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create receiver account with auto associations

	// when

	// reject the token with invalid token - should fail with INVALID_TOKEN_ID

	// same test for nft
}
