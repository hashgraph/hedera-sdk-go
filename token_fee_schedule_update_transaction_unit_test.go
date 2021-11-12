//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenFeeScheduleUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	fee := NewCustomFixedFee().SetDenominatingTokenID(tokenID).SetFeeCollectorAccountID(accountID)
	assert.NoError(t, err)

	tokenFeeUpdate := NewTokenFeeScheduleUpdateTransaction().
		SetCustomFees([]Fee{fee}).
		SetTokenID(tokenID)

	err = tokenFeeUpdate._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTokenFeeScheduleUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	fee := NewCustomFixedFee().SetDenominatingTokenID(tokenID).SetFeeCollectorAccountID(accountID)
	assert.NoError(t, err)

	tokenFeeUpdate := NewTokenFeeScheduleUpdateTransaction().
		SetCustomFees([]Fee{fee}).
		SetTokenID(tokenID)

	err = tokenFeeUpdate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
