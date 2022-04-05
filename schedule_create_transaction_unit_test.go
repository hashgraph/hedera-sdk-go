//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitScheduleCreateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	scheduleCreate := NewScheduleCreateTransaction().
		SetPayerAccountID(accountID)

	err = scheduleCreate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleCreateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	scheduleCreate := NewScheduleCreateTransaction().
		SetPayerAccountID(accountID)

	err = scheduleCreate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitScheduleInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	scheduleInfo := NewScheduleInfoQuery().
		SetScheduleID(scheduleID)

	err = scheduleInfo._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleInfoQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	scheduleInfo := NewScheduleInfoQuery().
		SetScheduleID(scheduleID)

	err = scheduleInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitScheduleSignTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	scheduleSign := NewScheduleSignTransaction().
		SetScheduleID(scheduleID)

	err = scheduleSign._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleSignTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	scheduleSign := NewScheduleSignTransaction().
		SetScheduleID(scheduleID)

	err = scheduleSign._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitScheduleDeleteTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	scheduleDelete := NewScheduleDeleteTransaction().
		SetScheduleID(scheduleID)

	err = scheduleDelete._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleDeleteTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	scheduleDelete := NewScheduleDeleteTransaction().
		SetScheduleID(scheduleID)

	err = scheduleDelete._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}
