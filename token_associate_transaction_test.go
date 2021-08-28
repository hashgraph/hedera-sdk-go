package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationTokenAssociateTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
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
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.TokenID)
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

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	check := false
	for _, relation := range info.TokenRelationships {
		if tokenID.String() == relation.TokenID.String() {
			check = true
		}
	}
	assert.Truef(t, check, fmt.Sprintf("token associate transaction didnt work"))

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

func TestIntegrationTokenAssociateTransactionNoAccountID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_ACCOUNT_ID received for transaction %s", resp.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTokenAssociateTransactionNoTokenID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
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
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	nodeID := resp.NodeID

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	check := false
	for _, relation := range info.TokenRelationships {
		if tokenID.String() == relation.TokenID.String() {
			check = true
		}
	}
	assert.Falsef(t, check, fmt.Sprintf("token associate transaction somehow worked"))

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{nodeID}).
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

func TestIntegrationTokenAssociateTransactionAutoAssociate(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(20)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxAutomaticTokenAssociations(2).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	toke, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(2).
		SetTreasuryAccountID(accountID).
		SetAdminKey(newKey.PublicKey()).
		SetFreezeKey(newKey.PublicKey()).
		SetWipeKey(newKey.PublicKey()).
		SetKycKey(newKey.PublicKey()).
		SetSupplyKey(newKey.PublicKey()).
		SetFreezeDefault(false).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	toke.Sign(newKey)

	resp, err = toke.Execute(env.Client)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.TokenID)

	record, err := NewTransactionRecordQuery().
		SetTransactionID(resp.TransactionID).
		Execute(env.Client)
	assert.NoError(t, err)

	for _, s := range record.AutomaticTokenAssociations {
		println(s.AccountID.String(), s.TokenID.String())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
