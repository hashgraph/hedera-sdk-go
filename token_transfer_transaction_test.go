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
		if *tokenID.Checksum == *relation.TokenID.Checksum {
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
