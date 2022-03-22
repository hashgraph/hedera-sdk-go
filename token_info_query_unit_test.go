//go:build all || unit
// +build all unit

package hedera

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tokenInfo := NewTokenInfoQuery().
		SetTokenID(tokenID)

	err = tokenInfo._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenInfoQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenInfo := NewTokenInfoQuery().
		SetTokenID(tokenID)

	err = tokenInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenInfoFromBytesBadBytes(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString("tfhyY++/Q4BycortAgD4cmMKACB/")
	require.NoError(t, err)

	_, err = TokenInfoFromBytes(bytes)
	require.NoError(t, err)
}

func TestUnitTokenInfoFromBytesNil(t *testing.T) {
	_, err := TokenRelationshipFromBytes(nil)
	assert.Error(t, err)
}

func TestUnitTokenInfoFromBytesEmptyBytes(t *testing.T) {
	_, err := TokenInfoFromBytes([]byte{})
	require.NoError(t, err)
}
