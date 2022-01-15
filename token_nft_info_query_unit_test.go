//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenNftGetInfoByNftIDValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	nftID, err := NftIDFromString("2@0.0.123-rmkyk")
	require.NoError(t, err)

	nftInfo := NewTokenNftInfoQuery().
		SetNftID(nftID)

	err = nftInfo._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenNftGetInfoByNftIDValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	nftID, err := NftIDFromString("2@0.0.123-rmkykd")
	require.NoError(t, err)

	nftInfo := NewTokenNftInfoQuery().
		SetNftID(nftID)

	err = nftInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
