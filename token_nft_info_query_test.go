package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenNftGetInfo_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(TokenTypeNonFungibleUnique).
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
	metaData := []byte{50}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAmount(1).
		SetMetadata(metaData).
		Execute(env.Client)
	assert.NoError(t, err)

	mintReceipt, err := mint.GetReceipt(env.Client)
	assert.NoError(t, err)

	nftID := tokenID.GenerateNft(mintReceipt.SerialNumbers[0])

	info, err := NewTokenNftInfoQuery().
		ByNftID(nftID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	value := false
	for _, nftInfo := range info {
		if tokenID == nftInfo.NftID.TokenID {
			value = true
		}
	}
	assert.Truef(t, value, fmt.Sprintf("token nft transfer transaction failed"))
	assert.Equal(t, len(info), 1)
	assert.Equal(t, info[0].NftID, nftID)
	assert.Equal(t, info[0].Metadata[0], byte(50))
}

func TestTokenNftGetInfoByTokenID_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(TokenTypeNonFungibleUnique).
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
	metaData := [][]byte{{50}, {50}}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAmount(100).
		SetMetadatas(metaData).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = mint.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewTokenNftInfoQuery().
		ByTokenID(tokenID).
		SetStart(0).
		SetEnd(2).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, len(info), 2)

	assert.Equal(t, info[0].NftID.TokenID, tokenID)
	assert.Equal(t, info[0].AccountID, env.Client.GetOperatorAccountID())
}
