//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenDeleteTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	tokenDelete := NewTokenDeleteTransaction().
		SetTokenID(tokenID)

	err = tokenDelete._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenDeleteTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenDelete := NewTokenDeleteTransaction().
		SetTokenID(tokenID)

	err = tokenDelete._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
