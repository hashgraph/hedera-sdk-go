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

func TestIntegrationAccountDeleteTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	accountID := receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransactionID(TransactionIDGenerate(*accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx = tx.Sign(newKey)

	assert.True(t, newKey.PublicKey().VerifyTransaction(tx))
	assert.False(t, env.Client.GetOperatorPublicKey().VerifyTransaction(tx))

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountDeleteTransactionNoTransferAccountID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	accountID := receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST received for transaction %s", resp.TransactionID), err.Error())
	}

}

func TestIntegrationAccountDeleteTransactionNoAccountID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST received for transaction %s", resp.TransactionID), err.Error())
	}

}

func TestIntegrationAccountDeleteTransactionNoSigning(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	accountID := receipt.AccountID
	require.NoError(t, err)

	acc := *accountID

	resp, err = NewAccountDeleteTransaction().
		SetAccountID(acc).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_SIGNATURE", err.Error())
	}

}

func TestIntegrationAccountDeleteTransactionCannotDeleteWithPendingAirdrops(t *testing.T) {
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

	// Try to delete the sender
	accountDeleteResp, err := NewAccountDeleteTransaction().
		SetAccountID(env.OperatorID).
		SetTransferAccountID(receiver).
		Execute(env.Client)
	receipt, err = accountDeleteResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "ACCOUNT_HAS_PENDING_AIRDROPS")
}
