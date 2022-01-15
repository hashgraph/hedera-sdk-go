//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountStakersQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	stackersQuery := NewAccountStakersQuery().
		SetAccountID(accountID)

	err = stackersQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountStakersQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	stackersQuery := NewAccountStakersQuery().
		SetAccountID(accountID)

	err = stackersQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
