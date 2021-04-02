package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenWipeTransaction_Execute(t *testing.T) {
	client := newTestClient(t, true)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

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
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -100).
		AddTokenTransfer(tokenID, accountID, 100).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountBalanceQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	var value uint64
	for balanceTokenID, balance := range info.Token {
		if tokenID == balanceTokenID {
			value = balance
		}
	}

	assert.Equal(t, uint64(100), value)

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(100).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err = NewAccountBalanceQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	for balanceTokenID, balance := range info.Token {
		if tokenID == balanceTokenID {
			value = balance
		}
	}

	assert.Equal(t, uint64(0), value)

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

func Test_TokenWipe_NoAmount(t *testing.T) {
	client := newTestClient(t, true)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

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

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp2, err := NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_WIPING_AMOUNT received for transaction %s", resp2.TransactionID), err.Error())
	}

	receipt, err = resp2.GetReceipt(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("Invalid node AccountID was set for transaction: %s", resp2.NodeID), err.Error())
	}

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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES"), err.Error())
	}
}

func Test_TokenWipe_NoTokenID(t *testing.T) {
	client := newTestClient(t, true)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

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

	resp2, err := NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOKEN_ID received for transaction %s", resp2.TransactionID), err.Error())
	}

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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES"), err.Error())
	}
}

func Test_TokenWipe_NoAccountID(t *testing.T) {
	client := newTestClient(t, true)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

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

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp2, err := NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetTokenID(tokenID).
		SetAmount(10).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_ACCOUNT_ID received for transaction %s", resp2.TransactionID), err.Error())
	}

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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES"), err.Error())
	}
}

func TestTokenWipeTransaction_NotZeroTokensAtDelete_Execute(t *testing.T) {
	client := newTestClient(t, true)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

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
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -100).
		AddTokenTransfer(tokenID, accountID, 100).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountBalanceQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	var value uint64
	for balanceTokenID, balance := range info.Token {
		if tokenID == balanceTokenID {
			value = balance
		}
	}

	assert.Equal(t, uint64(100), value)

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err = NewAccountBalanceQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	for balanceTokenID, balance := range info.Token {
		if tokenID == balanceTokenID {
			value = balance
		}
	}

	assert.Equal(t, uint64(90), value)

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
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES"), err.Error())
	}
}
