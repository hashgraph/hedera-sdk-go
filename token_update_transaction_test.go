package hedera

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationTokenUpdateTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
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

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.Equalf(t, "A", info.Symbol, "token failed to update")

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

func TestUnitTokenUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	tokenUpdate := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenUpdate._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTokenUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	tokenUpdate := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenUpdate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch; some IDs have different networks set", err.Error())
	}
}

func TestIntegrationTokenUpdateTransactionDifferentKeys(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 5)
	pubKeys := make([]PublicKey, 5)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(pubKeys[0]).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

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

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffffc").
		SetTokenID(tokenID).
		SetTokenSymbol("K").
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(pubKeys[1]).
		SetWipeKey(pubKeys[2]).
		SetKycKey(pubKeys[3]).
		SetSupplyKey(pubKeys[4]).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, "K", info.Symbol)
	assert.Equal(t, "ffffc", info.Name)
	if info.FreezeKey != nil {
		freezeKey := info.FreezeKey
		assert.Equal(t, pubKeys[1].String(), freezeKey.String())
	}
	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

func TestIntegrationTokenUpdateTransactionNoTokenID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
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

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp2, err := NewTokenUpdateTransaction().
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOKEN_ID received for transaction %s", resp2.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	assert.NoError(t, err)
}
