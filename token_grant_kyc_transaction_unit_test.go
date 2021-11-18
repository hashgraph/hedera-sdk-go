//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenGrantKycTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	tokenGrantKyc := NewTokenGrantKycTransaction().
		SetAccountID(accountID).
		SetTokenID(tokenID)

	err = tokenGrantKyc._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenGrantKycTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenGrantKyc := NewTokenGrantKycTransaction().
		SetAccountID(accountID).
		SetTokenID(tokenID)

	err = tokenGrantKyc._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
