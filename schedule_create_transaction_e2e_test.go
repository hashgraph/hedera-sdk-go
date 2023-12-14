//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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

func TestIntegrationScheduleCreateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 2)
	pubKeys := make([]PublicKey, 2)

	for i := range keys {
		newKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	keyList := NewKeyList().
		AddAllPublicKeys(pubKeys)

	createResponse, err := NewAccountCreateTransaction().
		SetKey(keyList).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(10)).
		Execute(env.Client)
	require.NoError(t, err)

	transactionReceipt, err := createResponse.SetValidateStatus(true).GetReceipt(env.Client)

	transactionID := TransactionIDGenerate(env.OperatorID)
	newAccountID := *transactionReceipt.AccountID

	transferTx := NewTransferTransaction().
		SetTransactionID(transactionID).
		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), HbarFrom(1, HbarUnits.Hbar))

	scheduled, err := transferTx.Schedule()
	require.NoError(t, err)

	scheduleResponse, err := scheduled.
		SetExpirationTime(time.Now().Add(30 * time.Minute)).
		Execute(env.Client)
	require.NoError(t, err)

	scheduleRecord, err := scheduleResponse.GetRecord(env.Client)
	require.NoError(t, err)

	scheduleID := *scheduleRecord.Receipt.ScheduleID

	signTransaction, err := NewScheduleSignTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		FreezeWith(env.Client)

	signTransaction.Sign(keys[0])

	resp, err := signTransaction.Execute(env.Client)
	require.NoError(t, err)

	// Getting the receipt to make sure the signing executed properly
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Making sure the scheduled transaction executed properly with schedule info query
	info, err := NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	signTransaction, err = NewScheduleSignTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		FreezeWith(env.Client)

	// Signing the scheduled transaction
	signTransaction.Sign(keys[1])

	resp, err = signTransaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	require.NotNil(t, info.ExecutedAt)
}

//
// func DisabledTestIntegrationScheduleCreateTransactionMultiSign(t *testing.T) {
// env := NewIntegrationTestEnv(t)
//
//	keys := make([]PrivateKey, 3)
//	pubKeys := make([]PublicKey, 3)
//
//	for i := range keys {
//		newKey, err := PrivateKeyGenerateEd25519()
//		require.NoError(t, err)
//
//		keys[i] = newKey
//		pubKeys[i] = newKey.PublicKey()
//	}
//
//	keyList := NewKeyList().
//		AddAllPublicKeys(pubKeys)
//
//	createResponse, err := NewAccountCreateTransaction().
//		SetKey(keyList).
//		SetInitialBalance(NewHbar(10)).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	transactionReceipt, err := createResponse.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	newAccountID := *transactionReceipt.AccountID
//
//	transferTx := NewTransferTransaction().
//		SetTransactionID(transactionID).
//		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
//		AddHbarTransfer(env.Client.GetOperatorAccountID(), HbarFrom(1, HbarUnits.Hbar))
//
//	scheduled, err := transferTx.Schedule()
//	require.NoError(t, err)
//
//	scheduleResponse, err := scheduled.Execute(env.Client)
//	require.NoError(t, err)
//
//	scheduleReceipt, err := scheduleResponse.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	scheduleID := *scheduleReceipt.ScheduleID
//
//	info, err := NewScheduleInfoQuery().
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		SetScheduleID(scheduleID).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	transfer, err := info.GetScheduledTransaction()
//	require.NoError(t, err)
//	assert.NotNil(t, transfer)
//
//	signTransaction, err := NewScheduleSignTransaction().
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		SetScheduleID(scheduleID).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	signTransaction.Sign(keys[0])
//	signTransaction.Sign(keys[1])
//	signTransaction.Sign(keys[2])
//
//	resp, err := signTransaction.Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	info2, err := NewScheduleInfoQuery().
//		SetScheduleID(scheduleID).
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		Execute(env.Client)
//	require.NoError(t, err)
//	assert.False(t, info2.ExecutedAt.IsZero())
//}
//
// func DisabledTestIntegrationScheduleDeleteTransactionCanExecute(t *testing.T) {
// env := NewIntegrationTestEnv(t)
//
//	key, err := GeneratePrivateKey()
//	key2, err := GeneratePrivateKey()
//	require.NoError(t, err)
//
//	createResponse, err := NewAccountCreateTransaction().
//		SetKey(key).
//		SetInitialBalance(NewHbar(10)).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	transactionReceipt, err := createResponse.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	newAccountID := *transactionReceipt.AccountID
//
//	transferTx := NewTransferTransaction().
//		SetTransactionID(transactionID).
//		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
//		AddHbarTransfer(env.Client.GetOperatorAccountID(), HbarFrom(1, HbarUnits.Hbar))
//
//	scheduled, err := transferTx.Schedule()
//	require.NoError(t, err)
//
//	fr, err := scheduled.SetAdminKey(key2).FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	scheduleResponse, err := fr.Sign(key2).Execute(env.Client)
//	require.NoError(t, err)
//
//	scheduleReceipt, err := scheduleResponse.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	scheduleID := *scheduleReceipt.ScheduleID
//
//	info, err := NewScheduleInfoQuery().
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		SetScheduleID(scheduleID).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	transfer, err := info.GetScheduledTransaction()
//	require.NoError(t, err)
//	assert.NotNil(t, transfer)
//	assert.Nil(t, info.ExecutedAt)
//	assert.Nil(t, info.DeletedAt)
//
//	tx2, err := NewScheduleDeleteTransaction().
//		SetScheduleID(scheduleID).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	resp, err := tx2.
//		Sign(key2).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	info2, err := NewScheduleInfoQuery().
//		SetScheduleID(scheduleID).
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		Execute(env.Client)
//	require.NoError(t, err)
//	assert.False(t, info2.DeletedAt.IsZero())
//}
//
// func DisabledTestIntegrationScheduleCreateTransactionCheckValidGetTransaction(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := PrivateKeyGenerateEd25519()
//	require.NoError(t, err)
//
//	newBalance := NewHbar(1)
//
//	assert.Equal(t, HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	tx := NewAccountCreateTransaction().
//		SetTransactionID(transactionID).
//		SetKey(newKey.PublicKey()).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetMaxTransactionFee(NewHbar(2)).
//		SetInitialBalance(newBalance)
//
//	require.NoError(t, err)
//
//	scheduleTx, err := tx.Schedule()
//	require.NoError(t, err)
//
//	resp, err := scheduleTx.
//		SetPayerAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	info, err := NewScheduleInfoQuery().
//		SetScheduleID(*receipt.ScheduleID).
//		SetQueryPayment(NewHbar(2)).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	infoTx, err := info.GetScheduledTransaction()
//	require.NoError(t, err)
//
//	assert.NotNil(t, infoTx)
//
//	switch createTx := infoTx.(type) {
//	case *AccountCreateTransaction:
//		assert.Equal(t, createTx.pbBody.GetCryptoCreateAccount().InitialBalance, uint64(NewHbar(1).tinybar))
//	}
//
//	tx2, err := NewScheduleDeleteTransaction().
//		SetScheduleID(*receipt.ScheduleID).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	resp, err = tx2.
//		Sign(newKey).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	assert.Error(t, err)
//	if err != nil {
//		assert.Equal(t, "exceptional receipt status: SCHEDULE_ALREADY_EXECUTED", err.Error())
//	}
//}
//
// func DisabledTestIntegrationScheduleCreateTransactionDuplicateFails(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	key, err := GeneratePrivateKey()
//	key2, err := GeneratePrivateKey()
//	require.NoError(t, err)
//
//	createResponse, err := NewAccountCreateTransaction().
//		SetKey(key).
//		SetInitialBalance(NewHbar(10)).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	transactionReceipt, err := createResponse.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	newAccountID := *transactionReceipt.AccountID
//
//	transferTx := NewTransferTransaction().
//		SetTransactionID(transactionID).
//		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
//		AddHbarTransfer(env.Client.GetOperatorAccountID(), HbarFrom(1, HbarUnits.Hbar))
//
//	scheduled, err := transferTx.Schedule()
//	require.NoError(t, err)
//
//	fr, err := scheduled.SetAdminKey(key2).FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	fr.Sign(key2)
//
//	scheduleResponse, err := fr.Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = scheduleResponse.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	resp, err := fr.Execute(env.Client)
//	assert.Error(t, err)
//	if err != nil {
//		assert.Equal(t, fmt.Sprintf("exceptional precheck status DUPLICATE_TRANSACTION received for transaction %s", resp.TransactionID), err.Error())
//	}
//}
//
// func DisabledTestIntegrationScheduleCreateTransactionWithTransferTransaction(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	key, err := GeneratePrivateKey()
//	require.NoError(t, err)
//
//	_Response, err := NewAccountCreateTransaction().
//		SetKey(key).
//		SetInitialBalance(NewHbar(2)).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err := _Response.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	accountID := *receipt.AccountID
//
//	tx := NewTransferTransaction().
//		AddHbarTransfer(accountID, NewHbar(1).Negated()).
//		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(1))
//
//	scheduleTx, err := tx.Schedule()
//	require.NoError(t, err)
//
//	scheduleTx = scheduleTx.
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetPayerAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID()))
//
//	_Response, err = scheduleTx.Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err = _Response.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	scheduleID := *receipt.ScheduleID
//
//	scheduleSignTx, err := NewScheduleSignTransaction().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetScheduleID(scheduleID).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	scheduleSignTx.Sign(key)
//
//	_Response, err = scheduleSignTx.Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = _Response.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	info, err := NewScheduleInfoQuery().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetScheduleID(scheduleID).
//		Execute(env.Client)
//	require.NoError(t, err)
//	assert.NotNil(t, info.ExecutedAt)
//
//	err = CloseIntegrationTestEnv(env, nil)
//	require.NoError(t, err)
//}
//
// func DisabledTestIntegrationScheduledTokenNftTransferTransaction(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	keys := make([]PrivateKey, 3)
//	pubKeys := make([]PublicKey, 3)
//
//	for i := range keys {
//		newKey, err := PrivateKeyGenerateEd25519()
//		require.NoError(t, err)
//
//		keys[i] = newKey
//		pubKeys[i] = newKey.PublicKey()
//	}
//
//	keyList := NewKeyList().
//		AddAllPublicKeys(pubKeys)
//
//	newBalance := NewHbar(2)
//
//	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)
//
//	resp, err := NewAccountCreateTransaction().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetKey(keyList).
//		SetInitialBalance(newBalance).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	accountID := *receipt.AccountID
//
//	resp, err = NewTokenCreateTransaction().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetTokenName("ffff").
//		SetTokenSymbol("F").
//		SetTokenType(TokenTypeNonFungibleUnique).
//		SetSupplyType(TokenSupplyTypeFinite).
//		SetMaxSupply(5).
//		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		SetFreezeKey(env.Client.GetOperatorPublicKey()).
//		SetWipeKey(env.Client.GetOperatorPublicKey()).
//		SetKycKey(env.Client.GetOperatorPublicKey()).
//		SetSupplyKey(env.Client.GetOperatorPublicKey()).
//		SetFreezeDefault(false).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	tokenID := *receipt.TokenID
//	metaData := [][]byte{{50}, {50}}
//
//	mint, err := NewTokenMintTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetTokenID(tokenID).
//		SetMetadatas(metaData).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	mintReceipt, err := mint.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	transaction, err := NewTokenAssociateTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetAccountID(accountID).
//		SetTokenIDs(tokenID).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	resp, err = transaction.
//		Sign(keys[0]).
//		Sign(keys[1]).
//		Sign(keys[2]).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	resp, err = NewTokenGrantKycTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetAccountID(accountID).
//		SetTokenID(tokenID).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	tx := NewTransferTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[0]), env.OperatorID, accountID).
//		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[1]), env.OperatorID, accountID)
//
//	scheduleTx, err := tx.Schedule()
//	require.NoError(t, err)
//
//	scheduleTx = scheduleTx.
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetPayerAccountID(accountID).
//		SetAdminKey(env.OperatorKey).
//		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID()))
//
//	resp, err = scheduleTx.Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	scheduleID := *receipt.ScheduleID
//
//	info, err := NewScheduleInfoQuery().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetScheduleID(scheduleID).
//		Execute(env.Client)
//	require.NoError(t, err)
//	assert.Equal(t, info.CreatorAccountID.String(), env.OperatorID.String())
//
//	signTransaction, err := NewScheduleSignTransaction().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetScheduleID(scheduleID).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	signTransaction.Sign(keys[0])
//	signTransaction.Sign(keys[1])
//	signTransaction.Sign(keys[2])
//
//	resp, err = signTransaction.Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	info2, err := NewScheduleInfoQuery().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetScheduleID(scheduleID).
//		Execute(env.Client)
//	require.NoError(t, err)
//	assert.NotNil(t, info2.ExecutedAt)
//
//	nftInfo, err := NewTokenNftInfoQuery().
//		ByNftID(tokenID.Nft(mintReceipt.SerialNumbers[0])).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	assert.Equal(t, accountID.String(), nftInfo[0].AccountID.String())
//
//	err = CloseIntegrationTestEnv(env, &tokenID)
//	require.NoError(t, err)
//}
//
// func DisabledTestIntegrationScheduledTokenNftTransferTransactionSigned(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	keys := make([]PrivateKey, 3)
//	pubKeys := make([]PublicKey, 3)
//
//	for i := range keys {
//		newKey, err := PrivateKeyGenerateEd25519()
//		require.NoError(t, err)
//
//		keys[i] = newKey
//		pubKeys[i] = newKey.PublicKey()
//	}
//
//	keyList := NewKeyList().
//		AddAllPublicKeys(pubKeys)
//
//	newBalance := NewHbar(2)
//
//	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)
//
//	resp, err := NewAccountCreateTransaction().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetKey(keyList).
//		SetInitialBalance(newBalance).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	accountID := *receipt.AccountID
//
//	resp, err = NewTokenCreateTransaction().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetTokenName("ffff").
//		SetTokenSymbol("F").
//		SetTokenType(TokenTypeNonFungibleUnique).
//		SetSupplyType(TokenSupplyTypeFinite).
//		SetMaxSupply(5).
//		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		SetFreezeKey(env.Client.GetOperatorPublicKey()).
//		SetWipeKey(env.Client.GetOperatorPublicKey()).
//		SetKycKey(env.Client.GetOperatorPublicKey()).
//		SetSupplyKey(env.Client.GetOperatorPublicKey()).
//		SetFreezeDefault(false).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	tokenID := *receipt.TokenID
//	metaData := [][]byte{{50}, {50}}
//
//	mint, err := NewTokenMintTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetTokenID(tokenID).
//		SetMetadatas(metaData).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	mintReceipt, err := mint.GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	transaction, err := NewTokenAssociateTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetAccountID(accountID).
//		SetTokenIDs(tokenID).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	resp, err = transaction.
//		Sign(keys[0]).
//		Sign(keys[1]).
//		Sign(keys[2]).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	resp, err = NewTokenGrantKycTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetAccountID(accountID).
//		SetTokenID(tokenID).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	tx := NewTransferTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[0]), env.OperatorID, accountID).
//		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[1]), env.OperatorID, accountID)
//
//	scheduleTx, err := tx.Schedule()
//	require.NoError(t, err)
//
//	scheduleTx, err = scheduleTx.
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetPayerAccountID(accountID).
//		SetAdminKey(env.OperatorKey).
//		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID())).
//		FreezeWith(env.Client)
//	require.NoError(t, err)
//
//	scheduleTx.Sign(keys[0])
//	scheduleTx.Sign(keys[1])
//	scheduleTx.Sign(keys[2])
//
//	resp, err = scheduleTx.Execute(env.Client)
//	require.NoError(t, err)
//
//	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
//	require.NoError(t, err)
//
//	scheduleID := *receipt.ScheduleID
//
//	info2, err := NewScheduleInfoQuery().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetScheduleID(scheduleID).
//		Execute(env.Client)
//	require.NoError(t, err)
//	assert.NotNil(t, info2.ExecutedAt)
//
//	nftInfo, err := NewTokenNftInfoQuery().
//		ByNftID(tokenID.Nft(mintReceipt.SerialNumbers[0])).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		Execute(env.Client)
//	require.NoError(t, err)
//
//	assert.Equal(t, accountID.String(), nftInfo[0].AccountID.String())
//
//	err = CloseIntegrationTestEnv(env, &tokenID)
//	require.NoError(t, err)
//}
