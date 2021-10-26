package hedera

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationTokenDeleteTransactionCanExecute(t *testing.T) {
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

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	assert.NoError(t, err)
}

func TestUnitTokenDeleteTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	tokenDelete := NewTokenDeleteTransaction().
		SetTokenID(tokenID)

	err = tokenDelete._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTokenDeleteTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	tokenDelete := NewTokenDeleteTransaction().
		SetTokenID(tokenID)

	err = tokenDelete._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestIntegrationTokenDeleteTransactionNoKeys(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	tokenCreate, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.OperatorKey.PublicKey()).
		SetFreezeDefault(false).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	tokenCreate, err = tokenCreate.
		SignWithOperator(env.Client)
	assert.NoError(t, err)

	resp, err := tokenCreate.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	assert.NoError(t, err)
}

func TestIntegrationTokenDeleteTransactionNoTokenID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOKEN_ID received for transaction %s", resp.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
