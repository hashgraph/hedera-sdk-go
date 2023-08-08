//go:build all || e2e
// +build all e2e

package hedera

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

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCantTransferOnBehalfOfSpenderWithoutAllowanceApproval(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	spenderKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	spenderCreate, err := NewAccountCreateTransaction().SetKey(spenderKey).SetInitialBalance(NewHbar(2)).Execute(env.Client)
	require.NoError(t, err)
	spenderReceipt, err := spenderCreate.SetValidateStatus(true).GetReceipt(env.Client)
	spenderAccountId := spenderReceipt.AccountID
	receiverKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	receiverCreate, err := NewAccountCreateTransaction().SetKey(receiverKey).SetInitialBalance(NewHbar(2)).Execute(env.Client)
	require.NoError(t, err)
	receiverReceipt, err := receiverCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	receiverAccountId := receiverReceipt.AccountID
	tokenCreate, err := NewTokenCreateTransaction().SetTokenName("ffff").SetTokenSymbol("F").SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.OperatorID).SetKycKey(env.OperatorKey).SetFreezeKey(env.OperatorKey).
		SetWipeKey(env.OperatorKey).SetSupplyKey(env.OperatorKey).SetFreezeDefault(false).Execute(env.Client)
	require.NoError(t, err)
	tokenReceipt, err := tokenCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	tokenID := tokenReceipt.TokenID
	_, err = NewTokenAssociateTransaction().SetTokenIDs(*tokenID).SetAccountID(*spenderAccountId).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().
		SetTokenID(*tokenID).
		SetMetadata([]byte{0x01}).
		Execute(env.Client)
	require.NoError(t, err)
	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	serials := mintReceipt.SerialNumbers
	nft1 := NftID{TokenID: *tokenID, SerialNumber: serials[0]}
	onBehalfOfTxId := TransactionIDGenerate(*spenderAccountId)

	transfer, err := NewTransferTransaction().AddApprovedNftTransfer(nft1, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTxId).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.Error(t, err)
	require.Equal(t, "exceptional receipt status: SPENDER_DOES_NOT_HAVE_ALLOWANCE", err.Error())

}

func TestIntegrationCantTransferOnBehalfOfSpenderAfterRemovingTheAllowanceApproval(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	spenderKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	spenderCreate, err := NewAccountCreateTransaction().
		SetKey(spenderKey.PublicKey()).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)
	spenderAccountReceipt, err := spenderCreate.SetValidateStatus(true).GetReceipt(env.Client)
	spenderAccountId := spenderAccountReceipt.AccountID
	receiverKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	receiverCreate, err := NewAccountCreateTransaction().
		SetKey(receiverKey.PublicKey()).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)
	receiverAccountReceipt, err := receiverCreate.SetValidateStatus(true).GetReceipt(env.Client)
	receiverAccountId := receiverAccountReceipt.AccountID

	tokenCreate, err := NewTokenCreateTransaction().SetTokenName("ffff").SetTokenSymbol("F").SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.OperatorID).SetFreezeKey(env.OperatorKey).
		SetWipeKey(env.OperatorKey).SetSupplyKey(env.OperatorKey).SetFreezeDefault(false).Execute(env.Client)
	require.NoError(t, err)
	tokenReceipt, err := tokenCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	tokenID := tokenReceipt.TokenID

	_, err = NewTokenAssociateTransaction().SetTokenIDs(*tokenID).SetAccountID(*spenderAccountId).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = NewTokenAssociateTransaction().SetTokenIDs(*tokenID).SetAccountID(*receiverAccountId).Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(*tokenID).SetMetadata([]byte{0x01}).SetMetadata([]byte{0x02}).Execute(env.Client)
	require.NoError(t, err)
	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	serials := mintReceipt.SerialNumbers
	nft1 := NftID{TokenID: *tokenID, SerialNumber: serials[0]}
	nft2 := NftID{TokenID: *tokenID, SerialNumber: serials[1]}
	approveTx, err := NewAccountAllowanceApproveTransaction().ApproveTokenNftAllowance(nft1, env.OperatorID, *spenderAccountId).
		ApproveTokenNftAllowance(nft2, env.OperatorID, *spenderAccountId).Execute(env.Client)
	require.NoError(t, err)
	_, err = approveTx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	deleteTx, err := NewAccountAllowanceDeleteTransaction().DeleteAllTokenNftAllowances(nft2, &env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	_, err = deleteTx.SetValidateStatus(true).GetReceipt(env.Client)

	onBehalfOfTxId := TransactionIDGenerate(*spenderAccountId)
	transfer, err := NewTransferTransaction().AddApprovedNftTransfer(nft1, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTxId).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenNftInfoQuery().SetNftID(nft1).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, *receiverAccountId, info[0].AccountID)

	onBehalfOfTxId2 := TransactionIDGenerate(*spenderAccountId)
	transfer2, err := NewTransferTransaction().AddApprovedNftTransfer(nft2, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTxId2).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer2.SetValidateStatus(true).GetReceipt(env.Client)
	require.Error(t, err)
	require.Equal(t, "exceptional receipt status: SPENDER_DOES_NOT_HAVE_ALLOWANCE", err.Error())

}

func TestIntegrationCantRemoveSingleSerialNumberAllowanceWhenAllowanceIsForAllSerials(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	spenderKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	spenderCreate, err := NewAccountCreateTransaction().
		SetKey(spenderKey.PublicKey()).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)
	spenderAccountReceipt, err := spenderCreate.SetValidateStatus(true).GetReceipt(env.Client)
	spenderAccountId := spenderAccountReceipt.AccountID

	receiverKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	receiverCreate, err := NewAccountCreateTransaction().
		SetKey(receiverKey.PublicKey()).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)
	receiverAccountReceipt, err := receiverCreate.SetValidateStatus(true).GetReceipt(env.Client)
	receiverAccountId := receiverAccountReceipt.AccountID

	tokenCreate, err := NewTokenCreateTransaction().SetTokenName("ffff").SetTokenSymbol("F").SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.OperatorID).SetFreezeKey(env.OperatorKey).
		SetWipeKey(env.OperatorKey).SetSupplyKey(env.OperatorKey).SetFreezeDefault(false).Execute(env.Client)
	require.NoError(t, err)
	tokenReceipt, err := tokenCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	tokenID := tokenReceipt.TokenID

	_, err = NewTokenAssociateTransaction().SetTokenIDs(*tokenID).SetAccountID(*spenderAccountId).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = NewTokenAssociateTransaction().SetTokenIDs(*tokenID).SetAccountID(*receiverAccountId).Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(*tokenID).SetMetadata([]byte{0x01}).SetMetadata([]byte{0x02}).Execute(env.Client)
	require.NoError(t, err)
	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	serials := mintReceipt.SerialNumbers
	nft1 := NftID{TokenID: *tokenID, SerialNumber: serials[0]}
	nft2 := NftID{TokenID: *tokenID, SerialNumber: serials[1]}
	approveTx, err := NewAccountAllowanceApproveTransaction().ApproveTokenNftAllowanceAllSerials(nft1.TokenID, env.OperatorID, *spenderAccountId).Execute(env.Client)
	require.NoError(t, err)
	_, err = approveTx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	onBehalfOfTransactionId := TransactionIDGenerate(*spenderAccountId)
	transfer, err := NewTransferTransaction().AddApprovedNftTransfer(nft1, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTransactionId).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	deleteTx, err := NewAccountAllowanceDeleteTransaction().DeleteAllTokenNftAllowances(nft2, &env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	_, err = deleteTx.SetValidateStatus(true).GetReceipt(env.Client)

	onBehalfOfTransactionId2 := TransactionIDGenerate(*spenderAccountId)
	transfer2, err := NewTransferTransaction().AddApprovedNftTransfer(nft2, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTransactionId2).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer2.SetValidateStatus(true).GetReceipt(env.Client)

	info, err := NewTokenNftInfoQuery().SetNftID(nft1).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, *receiverAccountId, info[0].AccountID)

	info2, err := NewTokenNftInfoQuery().SetNftID(nft2).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, *receiverAccountId, info2[0].AccountID)
}

func TestIntegrationAfterGivenAllowanceForAllSerialsCanGiveSingleSerialToOtherAccounts(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	spenderKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	spenderCreate, err := NewAccountCreateTransaction().
		SetKey(spenderKey.PublicKey()).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)
	spenderAccountReceipt, err := spenderCreate.SetValidateStatus(true).GetReceipt(env.Client)
	spenderAccountId := spenderAccountReceipt.AccountID

	delegateSpenderKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	delegateSpenderCreate, err := NewAccountCreateTransaction().SetKey(delegateSpenderKey.PublicKey()).SetInitialBalance(NewHbar(2)).Execute(env.Client)
	require.NoError(t, err)
	delegateSpenderAccountReceipt, err := delegateSpenderCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	delegateSpenderAccountId := delegateSpenderAccountReceipt.AccountID

	receiverKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	receiverCreate, err := NewAccountCreateTransaction().SetKey(receiverKey.PublicKey()).SetInitialBalance(NewHbar(2)).Execute(env.Client)
	require.NoError(t, err)
	receiverAccountReceipt, err := receiverCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	receiverAccountId := receiverAccountReceipt.AccountID

	tokenCreate, err := NewTokenCreateTransaction().SetTokenName("ffff").SetTokenSymbol("F").SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.OperatorID).SetFreezeKey(env.OperatorKey).SetWipeKey(env.OperatorKey).SetSupplyKey(env.OperatorKey).SetFreezeDefault(false).Execute(env.Client)
	require.NoError(t, err)
	tokenReceipt, err := tokenCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	tokenID := tokenReceipt.TokenID

	_, err = NewTokenAssociateTransaction().SetTokenIDs(*tokenID).SetAccountID(*spenderAccountId).Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = NewTokenAssociateTransaction().SetTokenIDs(*tokenID).SetAccountID(*receiverAccountId).Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(*tokenID).SetMetadata([]byte{0x01}).SetMetadata([]byte{0x02}).Execute(env.Client)
	require.NoError(t, err)
	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	serials := mintReceipt.SerialNumbers
	nft1 := NftID{TokenID: *tokenID, SerialNumber: serials[0]}
	nft2 := NftID{TokenID: *tokenID, SerialNumber: serials[1]}

	approveTx, err := NewAccountAllowanceApproveTransaction().
		ApproveTokenNftAllowanceAllSerials(*tokenID, env.OperatorID, *spenderAccountId).Execute(env.Client)
	require.NoError(t, err)
	_, err = approveTx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(*spenderAccountId, spenderKey)

	approveDelegateTx, err := NewAccountAllowanceApproveTransaction().
		ApproveTokenNftAllowanceWithDelegatingSpender(nft1, env.OperatorID, *delegateSpenderAccountId, *spenderAccountId).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = approveDelegateTx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	onBehalfOfTransactionId := TransactionIDGenerate(*delegateSpenderAccountId)
	transfer, err := NewTransferTransaction().AddApprovedNftTransfer(nft1, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTransactionId).Sign(delegateSpenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	onBehalfOfTransactionId2 := TransactionIDGenerate(*delegateSpenderAccountId)
	transfer2, err := NewTransferTransaction().AddApprovedNftTransfer(nft2, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTransactionId2).Sign(delegateSpenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer2.SetValidateStatus(true).GetReceipt(env.Client)
	require.Error(t, err)
	require.Equal(t, "exceptional receipt status: SPENDER_DOES_NOT_HAVE_ALLOWANCE", err.Error())

	info, err := NewTokenNftInfoQuery().SetNftID(nft1).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, *receiverAccountId, info[0].AccountID)
	info2, err := NewTokenNftInfoQuery().SetNftID(nft2).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, env.OperatorID, info2[0].AccountID)
}
