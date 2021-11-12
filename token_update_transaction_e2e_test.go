//+build all e2e

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

func DisabledTestIntegrationTokenUpdateTransactionTreasury(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	tokenCreate, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetMaxSupply(20).
		SetSupplyType(TokenSupplyTypeFinite).
		SetTreasuryAccountID(accountID).
		SetAdminKey(newKey.PublicKey()).
		SetFreezeKey(newKey.PublicKey()).
		SetWipeKey(newKey.PublicKey()).
		SetKycKey(newKey.PublicKey()).
		SetSupplyKey(newKey.PublicKey()).
		SetFreezeDefault(false).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	tokenCreate.Sign(newKey)

	resp, err = tokenCreate.Execute(env.Client)
	assert.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetMaxSupply(20).
		SetSupplyType(TokenSupplyTypeFinite).
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
	metaData := make([]byte, 50, 101)

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadata(metaData).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = mint.GetReceipt(env.Client)
	assert.NoError(t, err)

	update, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTreasuryAccountID(accountID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	update.Sign(newKey)

	resp, err = update.Execute(env.Client)
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
