//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenNftGetInfoByNftIDValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	nftID, err := NftIDFromString("2@0.0.123-rmkyk")
	assert.NoError(t, err)

	nftInfo := NewTokenNftInfoQuery().
		SetNftID(nftID)

	err = nftInfo._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTokenNftGetInfoByNftIDValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	nftID, err := NftIDFromString("2@0.0.123-rmkykd")
	assert.NoError(t, err)

	nftInfo := NewTokenNftInfoQuery().
		SetNftID(nftID)

	err = nftInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
