//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitFileContentsQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(fileID)

	err = fileContents._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitFileContentsQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(fileID)

	err = fileContents._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
