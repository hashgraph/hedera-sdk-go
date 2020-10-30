package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenInfoQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTokenCreateTransaction().
		SetName("ffff").
		SetSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasury(client.GetOperatorID()).
		SetAdminKey(client.GetOperatorKey()).
		SetFreezeKey(client.GetOperatorKey()).
		SetWipeKey(client.GetOperatorKey()).
		SetKycKey(client.GetOperatorKey()).
		SetSupplyKey(client.GetOperatorKey()).
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

	assert.Equal(t, info.TokenID, tokenID)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(3))
	assert.Equal(t, info.Treasury, client.GetOperatorID())
	assert.Equal(t, (*info.AdminKey).String(), client.GetOperatorKey().String())
	assert.Equal(t, (*info.KycKey).String(), client.GetOperatorKey().String())
	assert.Equal(t, (*info.FreezeKey).String(), client.GetOperatorKey().String())
	assert.Equal(t, (*info.WipeKey).String(), client.GetOperatorKey().String())
	assert.Equal(t, (*info.SupplyKey).String(), client.GetOperatorKey().String())
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
