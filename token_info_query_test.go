package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeTokenInfoQuery(t *testing.T) {
	query := NewTokenInfoQuery().
		SetQueryPayment(NewHbar(2)).
		SetTokenID(TokenID{Token: 3}).
		Query

	assert.Equal(t, `tokenGetInfo:{header:{}token:{tokenNum:3}}`, strings.ReplaceAll(strings.ReplaceAll(query.pb.String(), " ", ""), "\n", ""))
}

func TestTokenInfoQuery_Execute(t *testing.T) {
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

	info, err := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(2)).
		SetTokenID(tokenID).
		SetQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, info.TokenID, tokenID)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(3))
	assert.Equal(t, info.Treasury, client.GetOperatorAccountID())
	assert.NotNil(t, info.AdminKey)
	assert.NotNil(t, info.KycKey)
	assert.NotNil(t, info.FreezeKey)
	assert.NotNil(t, info.WipeKey)
	assert.NotNil(t, info.SupplyKey)
	assert.Equal(t, info.AdminKey.String(), client.GetOperatorPublicKey().String())
	assert.Equal(t, info.KycKey.String(), client.GetOperatorPublicKey().String())
	assert.Equal(t, info.FreezeKey.String(), client.GetOperatorPublicKey().String())
	assert.Equal(t, info.WipeKey.String(), client.GetOperatorPublicKey().String())
	assert.Equal(t, info.SupplyKey.String(), client.GetOperatorPublicKey().String())
	assert.False(t, *info.DefaultFreezeStatus)
	assert.False(t, *info.DefaultKycStatus)
}

func TestTokenInfoQueryCost_Execute(t *testing.T) {
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

	infoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(client)
	assert.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)
}

func TestTokenInfoQueryCost_BigMax_Execute(t *testing.T) {
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

	infoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1000000)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(client)
	assert.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)
}

func TestTokenInfoQueryCost_SmallMax_Execute(t *testing.T) {
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

	infoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(client)
	assert.NoError(t, err)

	_, err = infoQuery.Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("cost of TokenInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tħ"), err.Error())
	}
}

func TestTokenInfoQueryCost_InsufficientCost_Execute(t *testing.T) {
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

	infoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		SetTokenID(tokenID)

	_, err = infoQuery.GetCost(client)
	assert.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(HbarFromTinybar(1)).Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INSUFFICIENT_TX_FEE"), err.Error())
	}
}

func Test_TokenInfo_NoPayment(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFreezeKey(client.GetOperatorPublicKey()).
		SetKycKey(client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	info, err := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, info.TokenID, tokenID)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(3))
	assert.Equal(t, info.Treasury, client.GetOperatorAccountID())
	assert.False(t, *info.DefaultFreezeStatus)
	assert.False(t, *info.DefaultKycStatus)
}

func Test_TokenInfo_NoTokenID(t *testing.T) {
	client := newTestClient(t)

	_, err := NewTokenInfoQuery().
		SetQueryPayment(NewHbar(1)).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOKEN_ID"), err.Error())
	}
}
