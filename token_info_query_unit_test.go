//+build all unit

package hedera

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	tokenInfo := NewTokenInfoQuery().
		SetTokenID(tokenID)

	err = tokenInfo._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTokenInfoQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	tokenInfo := NewTokenInfoQuery().
		SetTokenID(tokenID)

	err = tokenInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestUnitTokenInfoFromBytesBadBytes(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString("tfhyY++/Q4BycortAgD4cmMKACB/")
	assert.NoError(t, err)

	_, err = TokenInfoFromBytes(bytes)
	assert.NoError(t, err)
}

func TestUnitTokenInfoFromBytesNil(t *testing.T) {
	_, err := TokenRelationshipFromBytes(nil)
	assert.Error(t, err)
}

func TestUnitTokenInfoFromBytesEmptyBytes(t *testing.T) {
	_, err := TokenInfoFromBytes([]byte{})
	assert.NoError(t, err)
}
