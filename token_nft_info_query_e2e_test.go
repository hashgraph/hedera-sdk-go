//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func DisabledTestIntegrationTokenNftGetInfoByNftIDCanExecute(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID
	metaData := []byte{50}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadata(metaData).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mint.GetReceipt(env.Client)
	require.NoError(t, err)

	nftID := tokenID.Nft(mintReceipt.SerialNumbers[0])

	info, err := NewTokenNftInfoQuery().
		SetNftID(nftID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	value := false
	for _, nftInfo := range info {
		if tokenID.String() == nftInfo.NftID.TokenID.String() {
			value = true
		}
	}
	assert.Truef(t, value, "token nft transfer transaction failed")
	assert.Equal(t, len(info), 1)
	assert.Equal(t, info[0].NftID, nftID)
	assert.Equal(t, info[0].Metadata[0], byte(50))
}
