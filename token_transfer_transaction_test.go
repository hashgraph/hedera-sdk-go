package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenTransferTransaction_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	var value uint64
	for _, relation := range info.TokenRelationships {
		if tokenID.String() == relation.TokenID.String() {
			value = relation.Balance
		}
	}
	assert.Equalf(t, uint64(999990), value, fmt.Sprintf("token transfer transaction failed"))

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

func Test_TokenTransfer_NotZeroSum(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	nodeId := resp.NodeID

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp2, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status TRANSFERS_NOT_ZERO_SUM_FOR_TOKEN received for transaction %s", resp2.TransactionID), err.Error())
	}

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

//func TestTokenNftTransferTransaction_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(2)
//
//	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	resp, err := NewAccountCreateTransaction().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetKey(newKey.PublicKey()).
//		SetInitialBalance(newBalance).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	receipt, err := resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
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
//	assert.NoError(t, err)
//
//	receipt, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	tokenID := *receipt.TokenID
//	metaData := [][]byte{{50}, {50}}
//
//	mint, err := NewTokenMintTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetTokenID(tokenID).
//		SetMetadatas(metaData).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	mintReceipt, err := mint.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	transaction, err := NewTokenAssociateTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetAccountID(accountID).
//		SetTokenIDs(tokenID).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = transaction.
//		Sign(newKey).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = NewTokenGrantKycTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetAccountID(accountID).
//		SetTokenID(tokenID).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = NewTransferTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[0]), env.OperatorID, accountID).
//		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[1]), env.OperatorID, accountID).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	info, err := NewTokenNftInfoQuery().
//		ByNftID(tokenID.Nft(mintReceipt.SerialNumbers[0])).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	assert.Equal(t, accountID.String(), info[0].AccountID.String())
//
//	info, err = NewTokenNftInfoQuery().
//		ByNftID(tokenID.Nft(mintReceipt.SerialNumbers[1])).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	assert.Equal(t, accountID.String(), info[0].AccountID.String())
//
//	resp, err = NewTokenWipeTransaction().
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		SetTokenID(tokenID).
//		SetAccountID(accountID).
//		SetSerialNumbers(mintReceipt.SerialNumbers).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	tx, err := NewAccountDeleteTransaction().
//		SetAccountID(accountID).
//		SetTransferAccountID(env.Client.GetOperatorAccountID()).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = tx.
//		Sign(newKey).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//}

func TestTokenFeeScheduleUpdateRecursionDepthTransaction_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              100000000,
		DenominationTokenID: &tokenID,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, accountID, -10).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), 10).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	transferTx.Sign(newKey)

	resp, err = transferTx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional receipt status: CUSTOM_FEE_CHARGING_EXCEEDED_MAX_RECURSION_DEPTH"), err.Error())
	}
}

func TestTokenFeeScheduleUpdateHugeAmountTransaction_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(10).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1000,
		DenominationTokenID: nil,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -5).
		AddTokenTransfer(tokenID, accountID, 5).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transferTx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)
}

func TestTokenFeeScheduleUpdateHugeAmount1Transaction_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID2 := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(10).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &accountID2,
		},
		Amount:              1000,
		DenominationTokenID: nil,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -5).
		AddTokenTransfer(tokenID, accountID, 5).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transferTx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)
}

//func Test_TokenTransfer_Frozen(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(2)
//
//	assert.Equal(t, 2 * HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	resp, err := NewAccountCreateTransaction().
//		SetKey(newKey.PublicKey()).
//		SetInitialBalance(newBalance).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	receipt, err := resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	accountID := *receipt.AccountID
//
//	resp, err = NewTokenCreateTransaction().
//		SetTokenName("ffff").
//		SetTokenSymbol("F").
//		SetDecimals(3).
//		SetInitialSupply(1000000).
//		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		SetFreezeKey(env.Client.GetOperatorPublicKey()).
//		SetWipeKey(env.Client.GetOperatorPublicKey()).
//		SetKycKey(env.Client.GetOperatorPublicKey()).
//		SetSupplyKey(env.Client.GetOperatorPublicKey()).
//		SetFreezeDefault(false).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	receipt, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	tokenID := *receipt.TokenID
//
//	nodeId := resp.NodeID
//
//	transaction, err := NewTokenAssociateTransaction().
//		SetNodeAccountIDs([]AccountID{nodeId}).
//		SetAccountID(accountID).
//		SetTokenIDs(tokenID).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = transaction.
//		Sign(newKey).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = NewTokenGrantKycTransaction().
//		SetNodeAccountIDs([]AccountID{nodeId}).
//		SetAccountID(accountID).
//		SetTokenID(tokenID).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	transfer, err := NewTransferTransaction().
//		SetNodeAccountIDs([]AccountID{nodeId}).
//		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = transfer.
//		AddTokenTransfer(tokenID, accountID, 10).
//		Execute(env.Client)
//	assert.Error(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.Error(t, err)
//
//	resp, err = NewTokenWipeTransaction().
//		SetNodeAccountIDs([]AccountID{nodeId}).
//		SetTokenID(tokenID).
//		SetAccountID(accountID).
//		SetAmount(10).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	tx, err := NewAccountDeleteTransaction().
//		SetAccountID(accountID).
//		SetTransferAccountID(env.Client.GetOperatorAccountID()).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = tx.
//		Sign(newKey).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//}
