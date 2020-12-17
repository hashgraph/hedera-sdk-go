package hedera

import (
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
		SetQueryPayment(NewHbar(2)).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, info.TokenID, tokenID)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(3))
	assert.Equal(t, info.Treasury, client.GetOperatorAccountID())
	assert.NotNil(t, *info.AdminKey)
	assert.NotNil(t, *info.KycKey)
	assert.NotNil(t, *info.FreezeKey)
	assert.NotNil(t, *info.WipeKey)
	assert.NotNil(t, *info.SupplyKey)
	assert.Equal(t, (*info.AdminKey).String(), client.GetOperatorPublicKey().String())
	assert.Equal(t, (*info.KycKey).String(), client.GetOperatorPublicKey().String())
	assert.Equal(t, (*info.FreezeKey).String(), client.GetOperatorPublicKey().String())
	assert.Equal(t, (*info.WipeKey).String(), client.GetOperatorPublicKey().String())
	assert.Equal(t, (*info.SupplyKey).String(), client.GetOperatorPublicKey().String())
	assert.False(t, *info.DefaultFreezeStatus)
	assert.False(t, *info.DefaultKycStatus)

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
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
		SetQueryPayment(NewHbar(2)).
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

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_TokenInfo_NoTokenID(t *testing.T) {
	client := newTestClient(t)

	_, err := NewTokenInfoQuery().
		SetQueryPayment(NewHbar(2)).
		Execute(client)
	assert.Error(t, err)
}
