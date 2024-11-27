//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenAirdropTransactionTransfersTokensWhenAssociated(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create ft and nft
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	txResponse, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftSerials := receipt.SerialNumbers

	// Create receiver with unlimited auto associations and receiverSig = false
	receiver, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
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
	defer CloseIntegrationTestEnv(env, nil)

	// Create ft and nft
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	txResponse, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftSerials := receipt.SerialNumbers

	// Create receiver
	receiver, _, err := createAccount(&env)
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

func TestIntegrationTokenAirdropTransactionCreatesHollowAccount(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create ft and nft
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	txResponse, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftSerials := receipt.SerialNumbers

	// Create a ECDSA private key
	privateKey, err := PrivateKeyGenerateEcdsa()
	if err != nil {
		println(err.Error())
	}
	// Extract the ECDSA public key public key
	publicKey := privateKey.PublicKey()

	aliasAccountId := publicKey.ToAccountID(0, 0)

	// should lazy-create and transfer the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, *aliasAccountId).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, *aliasAccountId).
		AddTokenTransfer(tokenID, *aliasAccountId, 100).
		AddTokenTransfer(tokenID, env.OperatorID, -100).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Verify the receiver holds the tokens via query
	receiverAccountBalance, err := NewAccountBalanceQuery().
		SetAccountID(*aliasAccountId).
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

func TestIntegrationTokenAirdropTransactionWithCustomFees(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create receiver with unlimited auto associations and receiverSig = false
	receiver, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// create fungible token with custom fee another token
	customFeeTokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// make the custom fee to be paid by the sender and the fee collector to be the operator account
	fee := NewCustomFixedFee().
		SetFeeCollectorAccountID(env.OperatorID).
		SetDenominatingTokenID(customFeeTokenID).
		SetAmount(1).
		SetAllCollectorsAreExempt(true)

	txResponse, err := NewTokenCreateTransaction().
		SetTokenName("Test Fungible Token").
		SetTokenSymbol("TFT").
		SetTokenMemo("I was created for integration tests").
		SetDecimals(3).
		SetInitialSupply(1_000_000).
		SetMaxSupply(1_000_000).
		SetTreasuryAccountID(env.OperatorID).
		SetSupplyType(TokenSupplyTypeFinite).
		SetAdminKey(env.OperatorKey).
		SetFreezeKey(env.OperatorKey).
		SetSupplyKey(env.OperatorKey).
		SetMetadataKey(env.OperatorKey).
		SetPauseKey(env.OperatorKey).
		SetCustomFees([]Fee{fee}).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	tokenID := receipt.TokenID

	// create sender account with unlimited associations and send some tokens to it
	sender, senderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// associate the token to the sender
	frozenTxn, err := NewTokenAssociateTransaction().
		SetAccountID(sender).
		AddTokenID(*tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)
	txResponse, err = frozenTxn.Sign(senderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// send tokens to the sender
	txResponse, err = NewTransferTransaction().
		AddTokenTransfer(customFeeTokenID, sender, 100).
		AddTokenTransfer(customFeeTokenID, env.OperatorID, -100).
		AddTokenTransfer(*tokenID, sender, 100).
		AddTokenTransfer(*tokenID, env.OperatorID, -100).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// airdrop the tokens from the sender to the receiver
	frozenTx, err := NewTokenAirdropTransaction().
		AddTokenTransfer(*tokenID, receiver, 100).
		AddTokenTransfer(*tokenID, sender, -100).
		FreezeWith(env.Client)
	require.NoError(t, err)
	airdropTx, err := frozenTx.
		Sign(senderKey).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = airdropTx.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// verify the custom fee has been paid by the sender to the collector
	receiverAccountBalance, err := NewAccountBalanceQuery().
		SetAccountID(receiver).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), receiverAccountBalance.Tokens.Get(*tokenID))

	senderAccountBalance, err := NewAccountBalanceQuery().
		SetAccountID(sender).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), senderAccountBalance.Tokens.Get(*tokenID))
	assert.Equal(t, uint64(99), senderAccountBalance.Tokens.Get(customFeeTokenID))

	operatorBalance, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000-100), operatorBalance.Tokens.Get(*tokenID))
	assert.Equal(t, uint64(1_000_000-100+1), operatorBalance.Tokens.Get(customFeeTokenID))
}
func TestIntegrationTokenAirdropTransactionWithReceiverSigTrue(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible token
	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create nft
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	txResponse, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftSerials := receipt.SerialNumbers

	// create receiver with unlimited auto associations and receiverSig = true
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	accountCreateFrozen, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetInitialBalance(NewHbar(3)).
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(-1).
		FreezeWith(env.Client)
	require.NoError(t, err)
	accountCreate, err := accountCreateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)
	receipt, err = accountCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	receiver := *receipt.AccountID

	// airdrop the tokens
	airdropTx, err := NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, receiver).
		AddTokenTransfer(tokenID, receiver, 100).
		AddTokenTransfer(tokenID, env.OperatorID, -100).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = airdropTx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTokenAirdropTransactionWithNoBalanceFT(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible token
	tokenID, _ := createFungibleToken(&env)

	// create spender and approve to it some tokens
	spender, spenderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// create sender
	sender, senderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// transfer ft to sender
	txResponse, err := NewTransferTransaction().
		AddTokenTransfer(tokenID, sender, 100).
		AddTokenTransfer(tokenID, env.OperatorID, -100).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// approve allowance to the spender
	frozenTx, err := NewAccountAllowanceApproveTransaction().
		ApproveTokenAllowance(tokenID, sender, spender, 100).
		FreezeWith(env.Client)
	require.NoError(t, err)
	txResponse, err = frozenTx.Sign(senderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// airdrop the tokens from the sender to the spender via approval
	frozenTxn, err := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, spender, 100).
		AddApprovedTokenTransfer(tokenID, sender, -100, true).
		SetTransactionID(TransactionIDGenerate(spender)).
		FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.
		Sign(spenderKey).
		Execute(env.Client)
	assert.ErrorContains(t, err, "NOT_SUPPORTED")
}

func TestIntegrationTokenAirdropTransactionWithNoBalanceNFT(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create nft
	nftID, err := createNft(&env)
	require.NoError(t, err)

	// Mint some NFTs
	txResponse, err := NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(mintMetadata).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftSerials := receipt.SerialNumbers

	// create spender and approve to it some tokens
	spender, spenderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// create sender
	sender, senderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// transfer ft to sender
	txResponse, err = NewTransferTransaction().
		AddNftTransfer(nftID.Nft(nftSerials[0]), env.OperatorID, sender).
		AddNftTransfer(nftID.Nft(nftSerials[1]), env.OperatorID, sender).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// approve allowance to the spender
	frozenTx, err := NewAccountAllowanceApproveTransaction().
		ApproveTokenNftAllowance(nftID.Nft(nftSerials[0]), sender, spender).
		ApproveTokenNftAllowance(nftID.Nft(nftSerials[1]), sender, spender).
		FreezeWith(env.Client)
	require.NoError(t, err)
	txResponse, err = frozenTx.Sign(senderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = txResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// airdrop the tokens from the sender to the spender via approval
	frozenTxn, err := NewTokenAirdropTransaction().
		AddApprovedNftTransfer(nftID.Nft(nftSerials[0]), sender, spender, true).
		AddApprovedNftTransfer(nftID.Nft(nftSerials[1]), sender, spender, true).
		SetTransactionID(TransactionIDGenerate(spender)).
		FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.
		Sign(spenderKey).
		Execute(env.Client)
	assert.ErrorContains(t, err, "NOT_SUPPORTED")
}

func TestIntegrationTokenAirdropTransactionWithInvalidBody(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create fungible token
	tokenID, _ := createFungibleToken(&env)

	// create receiver
	receiver, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	_, err = NewTokenAirdropTransaction().
		Execute(env.Client)
	require.ErrorContains(t, err, "EMPTY_TOKEN_TRANSFER_BODY")

	_, err = NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, receiver, 100).
		AddTokenTransfer(tokenID, receiver, 100).
		Execute(env.Client)
	require.ErrorContains(t, err, "INVALID_TRANSACTION_BODY")
}
