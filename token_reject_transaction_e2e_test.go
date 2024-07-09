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
		SetPauseKey(env.Client.GetOperatorPublicKey()).
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
		SetPauseKey(env.Client.GetOperatorPublicKey()).
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

	accountCreate, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		SetMaxAutomaticTokenAssociations(100).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := accountCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	return *receipt.AccountID, newKey
}

func TestIntegrationTokenRejectTransactionCanExecuteForFungibleToken(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create fungible tokens with treasury
	tokenID1 := createFungibleTokenHelper(t, &env)
	tokenID2 := createFungibleTokenHelper(t, &env)

	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// transfer fts to the receiver
	tx, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID1, tokenID2).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 0
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID1))
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID2))

	// verify the tokens are transferred back to the treasury
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID2))
}

func TestIntegrationTokenRejectTransactionCanExecuteForNFT(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create nft collections with treasury
	nftID1 := createNftHelper(t, &env)
	nftID2 := createNftHelper(t, &env)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(initialMetadataList).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(initialMetadataList).Execute(env.Client)
	require.NoError(t, err)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers
	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// transfer nfts to the receiver
	tx, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID2.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID2.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.GetReceipt(env.Client)
	require.NoError(t, err)

	// reject one of the nfts
	frozenTxn, err := NewTokenRejectTransaction().SetOwnerID(receiver).SetNftIDs(nftID1.Nft(serials[1]), nftID2.Nft(serials[1])).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance is decremented by 1
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID2))

	// verify the token is transferred back to the treasury
	nftBalance, err := NewTokenNftInfoQuery().SetNftID(nftID1.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)

	nftBalance, err = NewTokenNftInfoQuery().SetNftID(nftID2.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)
}

func TestIntegrationTokenRejectTransactionCanExecuteForFTAndNFTAtTheSameTime(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create fungible tokens with treasury
	tokenID1 := createFungibleTokenHelper(t, &env)
	tokenID2 := createFungibleTokenHelper(t, &env)

	// create nft collections with treasury
	nftID1 := createNftHelper(t, &env)
	nftID2 := createNftHelper(t, &env)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(initialMetadataList).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(initialMetadataList).Execute(env.Client)
	require.NoError(t, err)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// transfer fts to the receiver
	tx1, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx1.GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nfts to the receiver
	tx2, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID2.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID2.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx2.GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID1, tokenID2).
		SetNftIDs(nftID1.Nft(serials[1]), nftID2.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID1))
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID2))
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID2))

	// verify the tokens are transferred back to the treasury
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID2))

	nftBalance, err := NewTokenNftInfoQuery().SetNftID(nftID1.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)

	nftBalance, err = NewTokenNftInfoQuery().SetNftID(nftID2.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, env.OperatorID, nftBalance[0].AccountID)
}

func TestIntegrationTokenRejectTransactionReceiverSigRequired(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create nft with treasury with receiver sig required
	treasuryKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	accountCreateFrozen, err := NewAccountCreateTransaction().
		SetKey(treasuryKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(100).FreezeWith(env.Client)
	require.NoError(t, err)
	accountCreate, err := accountCreateFrozen.Sign(treasuryKey).Execute(env.Client)
	receipt, err := accountCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	treasury := *receipt.AccountID

	frozenTokenCreate, err := NewTokenCreateTransaction().
		SetTokenName("Example Collection").SetTokenSymbol("ABC").
		SetTokenType(TokenTypeNonFungibleUnique).SetDecimals(0).
		SetInitialSupply(0).
		SetMaxSupply(10).
		SetTreasuryAccountID(treasury).
		SetSupplyType(TokenSupplyTypeFinite).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetMetadataKey(env.Client.GetOperatorPublicKey()).FreezeWith(env.Client)
	require.NoError(t, err)
	tokenCreate, err := frozenTokenCreate.Sign(treasuryKey).Execute(env.Client)
	receipt, err = tokenCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftID := *receipt.TokenID

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(initialMetadataList).Execute(env.Client)
	require.NoError(t, err)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers
	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// transfer nft to the receiver
	frozenTransfer, err := NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), treasury, receiver).
		AddNftTransfer(nftID.Nft(serials[1]), treasury, receiver).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer, err := frozenTransfer.Sign(treasuryKey).Execute(env.Client)
	_, err = transfer.GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance is decremented by 1
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), tokenBalance.Tokens.Get(nftID))

	// verify the token is transferred back to the treasury
	nftBalance, err := NewTokenNftInfoQuery().SetNftID(nftID.Nft(serials[1])).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, treasury, nftBalance[0].AccountID)

	// same test for fungible token

	// create fungible token with treasury with receiver sig required
	frozenTokenCreate, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("asdf").
		SetDecimals(18).
		SetInitialSupply(1_000_000).
		SetTreasuryAccountID(treasury).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).FreezeWith(env.Client)
	require.NoError(t, err)

	tokenCreate, err = frozenTokenCreate.Sign(treasuryKey).Execute(env.Client)

	receipt, err = tokenCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	// transfer ft to the receiver
	frozenTransfer, err = NewTransferTransaction().
		AddTokenTransfer(tokenID, treasury, -10).
		AddTokenTransfer(tokenID, receiver, 10).
		FreezeWith(env.Client)
	transfer, err = frozenTransfer.Sign(treasuryKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// reject the token
	frozenTxn, err = NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 0
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Zero(t, tokenBalance.Tokens.Get(tokenID))

	// verify the tokens are transferred back to the treasury
	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(treasury).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), tokenBalance.Tokens.Get(tokenID))

}

func TestIntegrationTokenRejectTransactionTokenFrozen(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// given

	// create nft with treasury
	nftID := createNftHelper(t, &env)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(initialMetadataList).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers
	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// when

	// transfer nft to the receiver
	tx, err := NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.GetReceipt(env.Client)
	require.NoError(t, err)

	// freeze the token
	tokenFreeze, err := NewTokenFreezeTransaction().SetTokenID(nftID).SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	_, err = tokenFreeze.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// then

	// reject the token - should fail with ACCOUNT_FROZEN_FOR_TOKEN
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_FROZEN_FOR_TOKEN")

	// same test with fungible token

	// create fungible token with treasury
	tokenID := createFungibleTokenHelper(t, &env)

	// transfer ft to the receiver
	tx, err = NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.GetReceipt(env.Client)
	require.NoError(t, err)

	// freeze the token
	tokenFreeze, err = NewTokenFreezeTransaction().SetTokenID(tokenID).SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	_, err = tokenFreeze.SetValidateStatus(true).GetReceipt(env.Client)

	// reject the token - should fail with ACCOUNT_FROZEN_FOR_TOKEN
	frozenTxn, err = NewTokenRejectTransaction().
		SetOwnerID(receiver).
		AddTokenID(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_FROZEN_FOR_TOKEN")

}

func TestIntegrationTokenRejectTransactionTokenPaused(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// given

	// create nft with treasury
	nftID := createNftHelper(t, &env)

	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(initialMetadataList).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers
	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// when

	// transfer nft to the receiver
	tx, err := NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.GetReceipt(env.Client)
	require.NoError(t, err)

	// pause the token
	tokenPause, err := NewTokenPauseTransaction().SetTokenID(nftID).Execute(env.Client)
	require.NoError(t, err)
	_, err = tokenPause.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// then

	// reject the token - should fail with TOKEN_IS_PAUSED
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[1])).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_PAUSED")

	// create fungible token with treasury
	tokenID := createFungibleTokenHelper(t, &env)

	// transfer ft to the receiver
	tx, err = NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.GetReceipt(env.Client)
	require.NoError(t, err)

	// pause the token
	tokenPause, err = NewTokenPauseTransaction().SetTokenID(tokenID).Execute(env.Client)
	require.NoError(t, err)
	_, err = tokenPause.SetValidateStatus(true).GetReceipt(env.Client)

	// reject the token - should fail with TOKEN_IS_PAUSED
	frozenTxn, err = NewTokenRejectTransaction().
		SetOwnerID(receiver).
		AddTokenID(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.ErrorContains(t, err, "TOKEN_IS_PAUSED")
}

func TestIntegrationTokenRejectTransactionRemovesAllowance(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury
	tokenID := createFungibleTokenHelper(t, &env)
	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)
	// create spender account to be approved
	spender, spenderKey := createAccountHelper(t, &env)

	// when

	// transfer ft to the receiver
	tx, err := NewTransferTransaction().
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.GetReceipt(env.Client)
	require.NoError(t, err)

	// approve allowance to the spender
	frozenApprove, err := NewAccountAllowanceApproveTransaction().
		ApproveTokenAllowance(tokenID, receiver, spender, 10).FreezeWith(env.Client)
	require.NoError(t, err)
	approve, err := frozenApprove.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the spender has allowance
	env.Client.SetOperator(spender, spenderKey)
	frozenTx, err := NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID, receiver, -5, true).
		AddTokenTransfer(tokenID, spender, 5).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer, err := frozenTx.SignWith(spenderKey.PublicKey(), spenderKey.Sign).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// reject the token
	frozenTxn, err := NewTokenRejectTransaction().
		SetOwnerID(receiver).
		AddTokenID(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// then

	// verify the allowance - should be 0 , because the receiver is no longer the owner
	env.Client.SetOperator(spender, spenderKey)
	frozenTx, err = NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID, receiver, -5, true).
		AddTokenTransfer(tokenID, spender, 5).FreezeWith(env.Client)
	tx, err = frozenTx.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.GetReceipt(env.Client)
	require.ErrorContains(t, err, "SPENDER_DOES_NOT_HAVE_ALLOWANCE")
	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

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

func TestIntegrationTokenRejectTransactionFailsWhenRejectingNFTWithTokenID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create nft with treasury
	nftID := createNftHelper(t, &env)
	// mint
	mint, err := NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(initialMetadataList).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers
	// create receiver account with auto associations
	receiver, key := createAccountHelper(t, &env)

	// transfer nfts to the receiver
	_, err = NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[1]), env.Client.GetOperatorAccountID(), receiver).
		AddNftTransfer(nftID.Nft(serials[2]), env.Client.GetOperatorAccountID(), receiver).
		Execute(env.Client)
	require.NoError(t, err)

	// reject the whole collection - should fail
	frozenTxn, err := NewTokenRejectTransaction().SetOwnerID(receiver).AddTokenID(nftID).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_AMOUNT_TRANSFERS_ONLY_ALLOWED_FOR_FUNGIBLE_COMMON")
}
