//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenFeeScheduleUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	fee := NewCustomFixedFee().SetDenominatingTokenID(tokenID).SetFeeCollectorAccountID(accountID)
	require.NoError(t, err)

	tokenFeeUpdate := NewTokenFeeScheduleUpdateTransaction().
		SetCustomFees([]Fee{fee}).
		SetTokenID(tokenID)

	err = tokenFeeUpdate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenFeeScheduleUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	fee := NewCustomFixedFee().SetDenominatingTokenID(tokenID).SetFeeCollectorAccountID(accountID)
	require.NoError(t, err)

	tokenFeeUpdate := NewTokenFeeScheduleUpdateTransaction().
		SetCustomFees([]Fee{fee}).
		SetTokenID(tokenID)

	err = tokenFeeUpdate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}
