package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenTransferTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2 * HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFreezeKey(client.GetOperatorPublicKey()).
		SetWipeKey(client.GetOperatorPublicKey()).
		SetKycKey(client.GetOperatorPublicKey()).
		SetSupplyKey(client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	var value uint64
	for _, relation := range info.TokenRelationships {
		if tokenID == relation.TokenID {
			value = relation.Balance
		}
	}
	assert.Equalf(t, uint64(999990), value, fmt.Sprintf("token transfer transaction failed"))

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_TokenTransfer_NotZeroSum(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2 * HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFreezeKey(client.GetOperatorPublicKey()).
		SetWipeKey(client.GetOperatorPublicKey()).
		SetKycKey(client.GetOperatorPublicKey()).
		SetSupplyKey(client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	nodeId := resp.NodeID

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp2, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -10).
		Execute(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status TRANSFERS_NOT_ZERO_SUM_FOR_TOKEN received for transaction %s", resp2.TransactionID), err.Error())

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

//func Test_TokenTransfer_Frozen(t *testing.T) {
//	client := newTestClient(t)
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
//		Execute(client)
//	assert.NoError(t, err)
//
//	receipt, err := resp.GetReceipt(client)
//	assert.NoError(t, err)
//
//	accountID := *receipt.AccountID
//
//	resp, err = NewTokenCreateTransaction().
//		SetTokenName("ffff").
//		SetTokenSymbol("F").
//		SetDecimals(3).
//		SetInitialSupply(1000000).
//		SetTreasuryAccountID(client.GetOperatorAccountID()).
//		SetAdminKey(client.GetOperatorPublicKey()).
//		SetFreezeKey(client.GetOperatorPublicKey()).
//		SetWipeKey(client.GetOperatorPublicKey()).
//		SetKycKey(client.GetOperatorPublicKey()).
//		SetSupplyKey(client.GetOperatorPublicKey()).
//		SetFreezeDefault(false).
//		Execute(client)
//	assert.NoError(t, err)
//
//	receipt, err = resp.GetReceipt(client)
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
//		FreezeWith(client)
//	assert.NoError(t, err)
//
//	resp, err = transaction.
//		Sign(newKey).
//		Execute(client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(client)
//	assert.NoError(t, err)
//
//	resp, err = NewTokenGrantKycTransaction().
//		SetNodeAccountIDs([]AccountID{nodeId}).
//		SetAccountID(accountID).
//		SetTokenID(tokenID).
//		Execute(client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(client)
//	assert.NoError(t, err)
//
//	transfer, err := NewTransferTransaction().
//		SetNodeAccountIDs([]AccountID{nodeId}).
//		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -10).
//		FreezeWith(client)
//	assert.NoError(t, err)
//
//	resp, err = transfer.
//		AddTokenTransfer(tokenID, accountID, 10).
//		Execute(client)
//	assert.Error(t, err)
//
//	_, err = resp.GetReceipt(client)
//	assert.Error(t, err)
//
//	resp, err = NewTokenWipeTransaction().
//		SetNodeAccountIDs([]AccountID{nodeId}).
//		SetTokenID(tokenID).
//		SetAccountID(accountID).
//		SetAmount(10).
//		Execute(client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(client)
//	assert.NoError(t, err)
//
//	tx, err := NewAccountDeleteTransaction().
//		SetAccountID(accountID).
//		SetTransferAccountID(client.GetOperatorAccountID()).
//		FreezeWith(client)
//	assert.NoError(t, err)
//
//	resp, err = tx.
//		Sign(newKey).
//		Execute(client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(client)
//	assert.NoError(t, err)
//}
