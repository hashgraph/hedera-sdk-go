//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"

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

// Helpers
const hexData = `608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`

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
		SetDecimals(3).
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

func createAccountWithInvalidMaxAutoAssociationsHelper(t *testing.T, env *IntegrationTestEnv, maxAutoAssociations int32) {
}

// Limited max auto association tests

func TestLimitedMaxAutoAssociationsFungibleTokensFlow(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create token1 with 1 mil supply
	tokenID1 := createFungibleTokenHelper(t, &env)
	// create token2 with 1 mil supply
	tokenID2 := createFungibleTokenHelper(t, &env)

	// account create with 1 max auto associations
	accountID1, _ := createAccountHelper(t, &env, 1)

	// account update with 1 max auto associations
	accountID2, newKey := createAccountHelper(t, &env, 100)
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// contract create with 1 max auto associations
	// contract update with 1 max auto associations
	// contract create flow with 1 max auto associations

	// transfer token1 to all some tokens
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, accountID1, 10).
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, accountID2, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer token2 to all should fail with NO_REMAINING_AUTOMATIC_ASSOCIATIONS
	tokenTransferTransaction2, err := NewTransferTransaction().
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, accountID1, 10).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, accountID2, 10).
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
	accountID1, _ := createAccountHelper(t, &env, 1)

	// account update with 1 max auto associations
	accountID2, newKey := createAccountHelper(t, &env, 1)

	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// contract create with 1 max auto associations
	// hexBytecode, err := hex.DecodeString(hexData)
	// require.NoError(t, err)

	// contractCreate, err := NewContractCreateTransaction().
	// 	SetAdminKey(env.Client.GetOperatorPublicKey()).
	// 	SetMaxAutomaticTokenAssociations(1).
	// 	SetGas(10000000).
	// 	SetBytecode(hexBytecode).
	// 	Execute(env.Client)
	// require.NoError(t, err)

	// receipt, err = contractCreate.SetValidateStatus(true).GetReceipt(env.Client)
	// require.NoError(t, err)

	// contractID1, err := AccountIDFromString(receipt.ContractID.String())
	// require.NoError(t, err)
	// // contract update with 1 max auto associations

	// contractCreate, err = NewContractCreateTransaction().
	// 	SetAdminKey(env.Client.GetOperatorPublicKey()).
	// 	SetMaxAutomaticTokenAssociations(100).
	// 	SetGas(10000000).
	// 	SetBytecode(hexBytecode).
	// 	Execute(env.Client)
	// require.NoError(t, err)

	// receipt, err = contractCreate.SetValidateStatus(true).GetReceipt(env.Client)
	// require.NoError(t, err)

	// contractID2, err := AccountIDFromString(receipt.ContractID.String())
	// require.NoError(t, err)

	// contractUpdate, err := NewContractUpdateTransaction().
	// 	SetContractID(*receipt.ContractID).
	// 	SetMaxAutomaticTokenAssociations(1).
	// 	Execute(env.Client)
	// require.NoError(t, err)

	// _, err = contractUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	// require.NoError(t, err)

	// // contract create flow with 1 max auto associations
	// contractCreateFlow, err := NewContractCreateFlow().
	// 	SetBytecode(hexBytecode).
	// 	SetAdminKey(env.Client.GetOperatorPublicKey()).
	// 	SetMaxAutomaticTokenAssociations(1).
	// 	SetGas(200000).
	// 	SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
	// 	SetContractMemo("[e2e::ContractCreateFlow]").
	// 	Execute(env.Client)
	// require.NoError(t, err)

	// receipt, err = contractCreateFlow.SetValidateStatus(true).GetReceipt(env.Client)
	// require.NoError(t, err)

	// contractCreateFlowID, err := AccountIDFromString(receipt.ContractID.String())
	// require.NoError(t, err)

	// transfer nft1 to all, 2 for each
	tokenTransferTransaction, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, accountID2).
		Execute(env.Client)

	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nft2 to all should fail with NO_REMAINING_AUTOMATIC_ASSOCIATIONS
	tokenTransferTransaction2, err := NewTransferTransaction().
		AddNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, accountID1).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction2.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "NO_REMAINING_AUTOMATIC_ASSOCIATIONS")

	tokenTransferTransaction3, err := NewTransferTransaction().
		AddNftTransfer(nftID2.Nft(serials[1]), env.OperatorID, accountID2).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction3.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "NO_REMAINING_AUTOMATIC_ASSOCIATIONS")
}

// HIP-904 Unlimited max auto association tests

func TestUnlimitedMaxAutoAssociationsExecutes(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// account create with unlimited max auto associations
	createAccountHelper(t, &env, -1)
	// account update with unlimited max auto associations
	accountID, newKey := createAccountHelper(t, &env, 100)
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// contract create
	// contract update
	// contract create flow

	// all should execute
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFungibleTokens(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create token1 with 1 mil supply
	tokenID1 := createFungibleTokenHelper(t, &env)
	// create token2 with 1 mil supply
	tokenID2 := createFungibleTokenHelper(t, &env)

	// account create with unlimited max auto associations
	accountID1, _ := createAccountHelper(t, &env, -1)
	// account update with unlimited max auto associations
	accountID2, newKey := createAccountHelper(t, &env, 100)
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// contract create
	// contract update
	// contract create flow

	// transfer to both accounts some token1 tokens
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID1, 1000).
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID2, 1000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer to both accounts some token2 tokens
	tokenTransferTransaction, err = NewTransferTransaction().
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID1, 1000).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID2, 1000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFungibleTokensWithDecimals(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// create token1 with 1 mil supply
	tokenID1 := createFungibleTokenHelper(t, &env)
	// create token2 with 1 mil supply
	tokenID2 := createFungibleTokenHelper(t, &env)

	// account create
	accountID1, _ := createAccountHelper(t, &env, -1)
	// account update
	accountID2, newKey := createAccountHelper(t, &env, 100)
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// contract create
	// contract update
	// contract create flow

	// transfer to both accounts some token1 tokens
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransferWithDecimals(tokenID1, env.Client.GetOperatorAccountID(), -1000, 10).
		AddTokenTransferWithDecimals(tokenID1, accountID1, 1000, 10).
		AddTokenTransferWithDecimals(tokenID1, env.Client.GetOperatorAccountID(), -1000, 10).
		AddTokenTransferWithDecimals(tokenID1, accountID2, 1000, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer to both accounts some token2 tokens
	tokenTransferTransaction, err = NewTransferTransaction().
		AddTokenTransferWithDecimals(tokenID2, env.Client.GetOperatorAccountID(), -1000, 10).
		AddTokenTransferWithDecimals(tokenID2, accountID1, 1000, 10).
		AddTokenTransferWithDecimals(tokenID2, env.Client.GetOperatorAccountID(), -1000, 10).
		AddTokenTransferWithDecimals(tokenID2, accountID2, 1000, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFromFungibleTokens(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	// create spender account which will be approved to spend
	spender, spenderKey := createAccountHelper(t, &env, 10)

	// create token1 with 1 mil supply
	tokenID1 := createFungibleTokenHelper(t, &env)
	// create token2 with 1 mil supply
	tokenID2 := createFungibleTokenHelper(t, &env)

	// account create with unlimited max auto associations
	accountID1, _ := createAccountHelper(t, &env, -1)
	// account update with unlimited max auto associations
	accountID2, newKey := createAccountHelper(t, &env, 100)
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// contract create
	// contract update
	// contract create flow

	// approve the spender
	approve, err := NewAccountAllowanceApproveTransaction().
		AddTokenApproval(tokenID1, spender, 2000).
		AddTokenApproval(tokenID2, spender, 2000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transferFrom to both accounts some token1 tokens
	tokenTransferTransactionFrozen, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID1, 1000).
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID2, 1000).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenTransferTransaction, err := tokenTransferTransactionFrozen.Sign(spenderKey).Execute(env.Client)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transferFrom token2 to all some tokens
	tokenTransferTransactionFrozen, err = NewTransferTransaction().
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID1, 1000).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID2, 1000).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenTransferTransaction, err = tokenTransferTransactionFrozen.Sign(spenderKey).Execute(env.Client)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// all should execute
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

	// account update with unlimited max auto associations
	accountID2, newKey := createAccountHelper(t, &env, 100)

	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// contract create
	// contract update
	// contract create flow

	// transfer nft1 to both accounts, 2 for each
	tokenTransferTransaction, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, accountID2).
		Execute(env.Client)

	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nft2 to both accounts, 2 for each
	tokenTransferTransaction, err = NewTransferTransaction().
		AddNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID2.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID2.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID2.Nft(serials[3]), env.OperatorID, accountID2).
		Execute(env.Client)

	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
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
	accountID1, _ := createAccountHelper(t, &env, -1)

	// account update with unlimited max auto associations
	accountID2, newKey := createAccountHelper(t, &env, 100)

	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// contract create
	// contract update
	// contract create flow

	// approve the spender
	approve, err := NewAccountAllowanceApproveTransaction().
		AddAllTokenNftApproval(tokenID1, spender).
		AddAllTokenNftApproval(tokenID2, spender).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transferFrom nft1 to all, 2 for each
	tokenTransferTransactionFrozen, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenTransferTransaction, err := tokenTransferTransactionFrozen.Sign(spenderKey).Execute(env.Client)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transferFrom nft2 to all, 2 for each
	tokenTransferTransactionFrozen, err = NewTransferTransaction().
		AddNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID2.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID2.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID2.Nft(serials[3]), env.OperatorID, accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenTransferTransaction, err = tokenTransferTransactionFrozen.Sign(spenderKey).Execute(env.Client)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// all should execute
}

func TestUnlimitedMaxAutoAssociationsFailsWithInvalid(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// account create with -2 and with -1000
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	accountCreate, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetMaxAutomaticTokenAssociations(-2).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = accountCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTOMATIC_ASSOCIATIONS")

	accountCreate, err = NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetMaxAutomaticTokenAssociations(-1000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = accountCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTOMATIC_ASSOCIATIONS")

	// account update with - 2 and with -1000
	accountID, newKey := createAccountHelper(t, &env, 100)
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-2).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTOMATIC_ASSOCIATIONS")

	accountUpdateFrozen, err = NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1000).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err = accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTOMATIC_ASSOCIATIONS")

	// contract create with -2 and with -1000
	// contract update with -2 and with -1000
	// contract create flow with -2 and with -1000

	// all fails with invalid max auto associations
}
