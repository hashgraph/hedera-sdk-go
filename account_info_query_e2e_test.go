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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountInfoQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)
	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	// sleep in order for mirror node information to update
	time.Sleep(3 * time.Second)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID, info.AccountID)
	assert.Equal(t, false, info.IsDeleted)
	assert.Equal(t, newKey.PublicKey(), info.Key)
	assert.Equal(t, newBalance.tinybar, info.Balance.tinybar)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	//err = CloseIntegrationTestEnv(env, nil)
	//require.NoError(t, err)
}

func TestIntegrationAccountInfoQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)
	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	// sleep in order for mirror node information to update
	time.Sleep(3 * time.Second)

	accountInfo := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := accountInfo.GetCost(env.Client)
	require.NoError(t, err)

	info, err := accountInfo.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID, info.AccountID)
	assert.Equal(t, false, info.IsDeleted)
	assert.Equal(t, newKey.PublicKey(), info.Key)
	assert.Equal(t, newBalance.tinybar, info.Balance.tinybar)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountInfoQueryInsufficientFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)
	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	accountInfo := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = accountInfo.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INSUFFICIENT_TX_FEE", err.Error())
	}

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountInfoQuerySetBigMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)
	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	// sleep in order for mirror node information to update
	time.Sleep(3 * time.Second)

	accountInfo := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetMaxQueryPayment(NewHbar(1000000)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = accountInfo.GetCost(env.Client)
	require.NoError(t, err)

	info, err := accountInfo.SetQueryPayment(NewHbar(1)).Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID, info.AccountID)
	assert.Equal(t, false, info.IsDeleted)
	assert.Equal(t, newKey.PublicKey(), info.Key)
	assert.Equal(t, newBalance.tinybar, info.Balance.tinybar)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountInfoQuerySetSmallMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)
	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	accountInfo := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := accountInfo.GetCost(env.Client)
	require.NoError(t, err)

	_, err = accountInfo.Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "cost of AccountInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tℏ", err.Error())
	}

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountInfoQueryNoAccountID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INVALID_ACCOUNT_ID", err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountInfoQueryTokenRelationshipStatuses(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)
	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxAutomaticTokenAssociations(10).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetMetadataKey(newKey).
		SetFreezeDefault(true).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	associateTxn, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(*receipt.TokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTxn.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// sleep in order for mirror node information to update
	time.Sleep(3 * time.Second)

	info, err := NewAccountInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, accountID, info.AccountID)
	assert.Equal(t, 1, len(info.TokenRelationships))
	assert.Equal(t, true, *info.TokenRelationships[0].FreezeStatus)
	assert.Equal(t, false, *info.TokenRelationships[0].KycStatus)

	unfreezeTxn, err := NewTokenUnfreezeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(*receipt.TokenID).
		FreezeWith(env.Client)

	require.NoError(t, err)
	resp, err = unfreezeTxn.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	kycUpdateTxn, err := NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(*receipt.TokenID).
		FreezeWith(env.Client)

	require.NoError(t, err)
	resp, err = kycUpdateTxn.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// sleep in order for mirror node information to update
	time.Sleep(3 * time.Second)

	info, err = NewAccountInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, accountID, info.AccountID)
	assert.Equal(t, 1, len(info.TokenRelationships))
	assert.Equal(t, false, *info.TokenRelationships[0].FreezeStatus)
	assert.Equal(t, true, *info.TokenRelationships[0].KycStatus)
}

func TestIntegrationAccountInfoQueryTokenRelationship(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)
	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxAutomaticTokenAssociations(10).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetMetadataKey(newKey).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	firstTokenID := *receipt.TokenID

	associateTxn, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(firstTokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTxn.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// sleep in order for mirror node information to update
	time.Sleep(3 * time.Second)

	info, err := NewAccountInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(info.TokenRelationships))
	assert.Equal(t, uint32(3), info.TokenRelationships[0].Decimals)

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(18).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetMetadataKey(newKey).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	secondTokenID := *receipt.TokenID

	associateTxn, err = NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(secondTokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTxn.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// sleep in order for mirror node information to update
	time.Sleep(3 * time.Second)

	info, err = NewAccountInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(info.TokenRelationships))
	assert.Equal(t, uint32(18), info.TokenRelationships[1].Decimals)

	dissociateTxn, err := NewTokenDissociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(secondTokenID, firstTokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = dissociateTxn.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// sleep in order for mirror node information to update
	time.Sleep(3 * time.Second)

	info, err = NewAccountInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, accountID, info.AccountID)
	assert.Equal(t, 0, len(info.TokenRelationships))
}

func TestIntegrationAccountInfoQueryWorksWithHollowAccountAlias(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// Create NFT
	cid := []string{"QmNPCiNA3Dsu3K5FxDPMG5Q3fZRwVTg14EXA92uqEeSRXn"}
	// Creating the transaction for token creation
	nftCreateTransaction, err := NewTokenCreateTransaction().
		SetTokenName("HIP-542 Example Collection").SetTokenSymbol("HIP-542").
		SetTokenType(TokenTypeNonFungibleUnique).SetDecimals(0).
		SetInitialSupply(0).SetMaxSupply(int64(len(cid))).
		SetTreasuryAccountID(env.OperatorID).SetSupplyType(TokenSupplyTypeFinite).
		SetAdminKey(env.OperatorKey).SetFreezeKey(env.OperatorKey).SetWipeKey(env.OperatorKey).SetSupplyKey(env.OperatorKey).FreezeWith(env.Client)

	// Sign the transaction with the operator key
	nftSignTransaction := nftCreateTransaction.Sign(env.OperatorKey)
	// Submit the transaction to the Hedera network
	nftCreateSubmit, err := nftSignTransaction.Execute(env.Client)
	require.NoError(t, err)

	// Get transaction receipt information
	nftCreateReceipt, err := nftCreateSubmit.GetReceipt(env.Client)
	require.NoError(t, err)

	// Get token id from the transaction
	nftTokenID := *nftCreateReceipt.TokenID

	nftCollection := []TransactionReceipt{}

	for _, s := range cid {
		mintTransaction, err := NewTokenMintTransaction().SetTokenID(nftTokenID).SetMetadata([]byte(s)).FreezeWith(env.Client)
		require.NoError(t, err)
		mintTransactionSubmit, err := mintTransaction.Sign(env.OperatorKey).Execute(env.Client)
		require.NoError(t, err)
		receipt, err := mintTransactionSubmit.GetReceipt(env.Client)
		require.NoError(t, err)
		nftCollection = append(nftCollection, receipt)
	}
	exampleNftId := nftTokenID.Nft(nftCollection[0].SerialNumbers[0])

	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	// Extract the ECDSA public key public key
	publicKey := privateKey.PublicKey()
	// Extract the Ethereum public address
	aliasAccountId := publicKey.ToAccountID(0, 0)

	nftTransferTransaction, err := NewTransferTransaction().AddNftTransfer(exampleNftId, env.OperatorID, *aliasAccountId).FreezeWith(env.Client)
	require.NoError(t, err)

	// Sign the transaction with the operator key
	nftTransferTransactionSign := nftTransferTransaction.Sign(env.OperatorKey)
	// Submit the transaction to the Hedera network
	nftTransferTransactionSubmit, err := nftTransferTransactionSign.Execute(env.Client)
	require.NoError(t, err)
	_, err = nftTransferTransactionSubmit.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Wait for mirror node to update
	time.Sleep(3 * time.Second)
	_, err = NewAccountInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(*aliasAccountId).
		Execute(env.Client)
	assert.NoError(t, err)
}
