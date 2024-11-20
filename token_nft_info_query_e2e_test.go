//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenNftGetInfoByNftIDCanExecute(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	metaData := []byte{50}

	mint, err := NewTokenMintTransaction().
		SetTokenID(tokenID).
		SetMetadata(metaData).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftID := tokenID.Nft(mintReceipt.SerialNumbers[0])

	info, err := NewTokenNftInfoQuery().
		SetNftID(nftID).
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
	parsedInfo, err := TokenNftInfoFromBytes(info[0].ToBytes())
	assert.NoError(t, err)
	assert.Equal(t, parsedInfo, info[0])
}
