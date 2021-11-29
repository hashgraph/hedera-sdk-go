//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitContractCreateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	contractCreate := NewContractCreateTransaction().
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractCreate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitContractCreateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	contractCreate := NewContractCreateTransaction().
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractCreate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
