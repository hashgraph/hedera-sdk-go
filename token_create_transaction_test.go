package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenCreateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTokenCreateTransaction().
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

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_TokenCreate_MultipleKeys(t *testing.T) {
	client := newTestClient(t)

	keys := make([]PrivateKey, 5)
	pubKeys := make([]PublicKey, 5)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	newBalance := NewHbar(2)

	assert.Equal(t, 2 * HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(pubKeys[0]).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFreezeKey(pubKeys[1]).
		SetWipeKey(pubKeys[2]).
		SetKycKey(pubKeys[3]).
		SetSupplyKey(pubKeys[4]).
		SetFreezeDefault(false).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_TokenCreate_NoKeys(t *testing.T) {
	client := newTestClient(t)

	keys := make([]PrivateKey, 6)
	pubKeys := make([]PublicKey, 6)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	newBalance := NewHbar(2)

	assert.Equal(t, 2 * HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(pubKeys[0]).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	info, err := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)

	assert.NoError(t, err)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(0))
	assert.Equal(t, info.TotalSupply, uint64(0))
	assert.Equal(t, info.Treasury.String(), client.GetOperatorAccountID().String())
	assert.Nil(t, info.AdminKey)
	assert.Nil(t, info.FreezeKey)
	assert.Nil(t, info.KycKey)
	assert.Nil(t, info.WipeKey)
	assert.Nil(t, info.SupplyKey)
	assert.Nil(t, info.DefaultFreezeStatus)
	assert.Nil(t, info.DefaultKycStatus)
	assert.NotNil(t, info.AutoRenewPeriod)
	assert.NotNil(t, info.ExpirationTime)

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status TOKEN_IS_IMMUTABLE"), err.Error())

}

func Test_TokenCreate_AdminSign(t *testing.T) {
	client := newTestClient(t)

	keys := make([]PrivateKey, 6)
	pubKeys := make([]PublicKey, 6)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	newBalance := NewHbar(2)

	assert.Equal(t, 2 * HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(pubKeys[0]).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenCreate, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(pubKeys[1]).
		SetFreezeKey(pubKeys[1]).
		SetWipeKey(pubKeys[2]).
		SetKycKey(pubKeys[3]).
		SetSupplyKey(pubKeys[4]).
		SetFreezeDefault(false).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tokenCreate.
		Sign(keys[0]).
		Sign(keys[1]).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.TokenID)
	tokenID := *receipt.TokenID

	tokenDelete, err := NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tokenDelete.
		Sign(keys[1]).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
