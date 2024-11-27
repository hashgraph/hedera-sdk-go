//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCantTransferOnBehalfOfSpenderWithoutAllowanceApproval(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	spenderKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	spenderCreate, err := NewAccountCreateTransaction().SetKey(spenderKey).SetInitialBalance(NewHbar(2)).Execute(env.Client)
	require.NoError(t, err)
	spenderReceipt, err := spenderCreate.SetValidateStatus(true).GetReceipt(env.Client)
	spenderAccountId := spenderReceipt.AccountID
	receiverKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	receiverCreate, err := NewAccountCreateTransaction().SetKey(receiverKey).SetInitialBalance(NewHbar(2)).SetMaxAutomaticTokenAssociations(10).Execute(env.Client)
	require.NoError(t, err)
	receiverReceipt, err := receiverCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	receiverAccountId := receiverReceipt.AccountID
	tokenID, err := createNft(&env)
	require.NoError(t, err)
	frozenTxn, err := NewTokenAssociateTransaction().SetTokenIDs(tokenID).SetAccountID(*spenderAccountId).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().
		SetTokenID(tokenID).
		SetMetadata([]byte{0x01}).
		Execute(env.Client)
	require.NoError(t, err)
	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	serials := mintReceipt.SerialNumbers
	nft1 := NftID{TokenID: tokenID, SerialNumber: serials[0]}
	onBehalfOfTxId := TransactionIDGenerate(*spenderAccountId)

	frozenTransfer, err := NewTransferTransaction().AddApprovedNftTransfer(nft1, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTxId).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTransfer.Sign(spenderKey).Execute(env.Client)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.Error(t, err)
	require.Equal(t, "exceptional receipt status: SPENDER_DOES_NOT_HAVE_ALLOWANCE", err.Error())

}

func TestIntegrationCantTransferOnBehalfOfSpenderAfterRemovingTheAllowanceApproval(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
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

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	frozenTx, err := NewTokenAssociateTransaction().SetTokenIDs(tokenID).SetAccountID(*spenderAccountId).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTx.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	frozenTx, err = NewTokenAssociateTransaction().SetTokenIDs(tokenID).SetAccountID(*receiverAccountId).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(tokenID).SetMetadata([]byte{0x01}).SetMetadata([]byte{0x02}).Execute(env.Client)
	require.NoError(t, err)
	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	serials := mintReceipt.SerialNumbers
	nft1 := NftID{TokenID: tokenID, SerialNumber: serials[0]}
	nft2 := NftID{TokenID: tokenID, SerialNumber: serials[1]}
	approveTx, err := NewAccountAllowanceApproveTransaction().ApproveTokenNftAllowance(nft1, env.OperatorID, *spenderAccountId).
		ApproveTokenNftAllowance(nft2, env.OperatorID, *spenderAccountId).Execute(env.Client)
	require.NoError(t, err)
	_, err = approveTx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	deleteTx, err := NewAccountAllowanceDeleteTransaction().DeleteAllTokenNftAllowances(nft2, &env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	_, err = deleteTx.SetValidateStatus(true).GetReceipt(env.Client)

	onBehalfOfTxId := TransactionIDGenerate(*spenderAccountId)
	frozenTransfer, err := NewTransferTransaction().AddApprovedNftTransfer(nft1, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTxId).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTransfer.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenNftInfoQuery().SetNftID(nft1).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, *receiverAccountId, info[0].AccountID)

	onBehalfOfTxId2 := TransactionIDGenerate(*spenderAccountId)
	frozenTransfer2, err := NewTransferTransaction().AddApprovedNftTransfer(nft2, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTxId2).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTransfer2.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.Error(t, err)
	require.Equal(t, "exceptional receipt status: SPENDER_DOES_NOT_HAVE_ALLOWANCE", err.Error())

}

func TestIntegrationCantRemoveSingleSerialNumberAllowanceWhenAllowanceIsForAllSerials(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
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

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	frozenTxn, err := NewTokenAssociateTransaction().SetTokenIDs(tokenID).SetAccountID(*spenderAccountId).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	frozenTxn, err = NewTokenAssociateTransaction().SetTokenIDs(tokenID).SetAccountID(*receiverAccountId).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(tokenID).SetMetadata([]byte{0x01}).SetMetadata([]byte{0x02}).Execute(env.Client)
	require.NoError(t, err)
	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	serials := mintReceipt.SerialNumbers
	nft1 := NftID{TokenID: tokenID, SerialNumber: serials[0]}
	nft2 := NftID{TokenID: tokenID, SerialNumber: serials[1]}
	approveTx, err := NewAccountAllowanceApproveTransaction().ApproveTokenNftAllowanceAllSerials(nft1.TokenID, env.OperatorID, *spenderAccountId).Execute(env.Client)
	require.NoError(t, err)
	_, err = approveTx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	onBehalfOfTransactionId := TransactionIDGenerate(*spenderAccountId)
	frozenTransfer, err := NewTransferTransaction().AddApprovedNftTransfer(nft1, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTransactionId).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTransfer.Sign(spenderKey).Execute(env.Client)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	deleteTx, err := NewAccountAllowanceDeleteTransaction().DeleteAllTokenNftAllowances(nft2, &env.OperatorID).Execute(env.Client)
	require.NoError(t, err)
	_, err = deleteTx.SetValidateStatus(true).GetReceipt(env.Client)

	onBehalfOfTransactionId2 := TransactionIDGenerate(*spenderAccountId)
	frozenTransfer2, err := NewTransferTransaction().AddApprovedNftTransfer(nft2, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTransactionId2).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTransfer2.Sign(spenderKey).Execute(env.Client)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)

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
	defer CloseIntegrationTestEnv(env, nil)
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

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	frozenTx, err := NewTokenAssociateTransaction().SetTokenIDs(tokenID).SetAccountID(*spenderAccountId).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTx.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	frozenTxn, err := NewTokenAssociateTransaction().SetTokenIDs(tokenID).SetAccountID(*receiverAccountId).FreezeWith(env.Client)
	require.NoError(t, err)
	_, err = frozenTxn.Sign(receiverKey).Execute(env.Client)

	mint, err := NewTokenMintTransaction().SetTokenID(tokenID).SetMetadata([]byte{0x01}).SetMetadata([]byte{0x02}).Execute(env.Client)
	require.NoError(t, err)
	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	serials := mintReceipt.SerialNumbers
	nft1 := NftID{TokenID: tokenID, SerialNumber: serials[0]}
	nft2 := NftID{TokenID: tokenID, SerialNumber: serials[1]}

	approveTx, err := NewAccountAllowanceApproveTransaction().
		ApproveTokenNftAllowanceAllSerials(tokenID, env.OperatorID, *spenderAccountId).Execute(env.Client)
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
	frozenTransfer, err := NewTransferTransaction().AddApprovedNftTransfer(nft1, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTransactionId).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenTransfer.Sign(delegateSpenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	onBehalfOfTransactionId2 := TransactionIDGenerate(*delegateSpenderAccountId)
	frozenTransfer2, err := NewTransferTransaction().AddApprovedNftTransfer(nft2, env.OperatorID, *receiverAccountId, true).SetTransactionID(onBehalfOfTransactionId2).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err = frozenTransfer2.Sign(delegateSpenderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.Error(t, err)
	require.Equal(t, "exceptional receipt status: SPENDER_DOES_NOT_HAVE_ALLOWANCE", err.Error())

	info, err := NewTokenNftInfoQuery().SetNftID(nft1).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, *receiverAccountId, info[0].AccountID)
	info2, err := NewTokenNftInfoQuery().SetNftID(nft2).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, env.OperatorID, info2[0].AccountID)
}
