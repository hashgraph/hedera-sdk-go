package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenNftGetInfoByNftID_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

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
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID
	metaData := []byte{50}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadata(metaData).
		Execute(env.Client)
	assert.NoError(t, err)

	mintReceipt, err := mint.GetReceipt(env.Client)
	assert.NoError(t, err)

	nftID := tokenID.Nft(mintReceipt.SerialNumbers[0])

	info, err := NewTokenNftInfoQuery().
		ByNftID(nftID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	value := false
	for _, nftInfo := range info {
		if tokenID.String() == nftInfo.NftID.TokenID.String() {
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
		SetTokenType(TokenTypeNonFungibleUnique).
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
		SetMetadatas(metaData).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = mint.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewTokenNftInfoQuery().
		ByTokenID(tokenID).
		SetEnd(2).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.Equal(t, len(info), 2)
	assert.Equal(t, info[0].NftID.TokenID, tokenID)

	resp, err = NewTokenBurnTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetSerialNumber(2).
		SetTokenID(tokenID).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.Equal(t, len(info), 2)

	info, err = NewTokenNftInfoQuery().
		ByTokenID(tokenID).
		SetEnd(1).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, len(info), 1)
	assert.Equal(t, info[0].NftID.TokenID, tokenID)
	assert.Equal(t, info[0].AccountID, env.Client.GetOperatorAccountID())
}

func TestTokenNftGetInfoByAccountID_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

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
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID
	metaData := []byte{50}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadata(metaData).
		Execute(env.Client)
	assert.NoError(t, err)

	mintReceipt, err := mint.GetReceipt(env.Client)
	assert.NoError(t, err)

	nftID := tokenID.Nft(mintReceipt.SerialNumbers[0])

	info, err := NewTokenNftInfoQuery().
		ByAccountID(env.OperatorID).
		SetEnd(1).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	value := false
	for _, nftInfo := range info {
		if tokenID.String() == nftInfo.NftID.TokenID.String() {
			value = true
		}
	}
	assert.Truef(t, value, fmt.Sprintf("token nft transfer transaction failed"))
	assert.Equal(t, len(info), 1)
	assert.Equal(t, info[0].NftID, nftID)
	assert.Equal(t, info[0].Metadata[0], byte(50))
}
