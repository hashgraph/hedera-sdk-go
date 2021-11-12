//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTransactionRecordQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	transactionID := TransactionIDGenerate(accountID)
	assert.NoError(t, err)

	recordQuery := NewTransactionRecordQuery().
		SetTransactionID(transactionID)

	err = recordQuery._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTransactionRecordQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	transactionID := TransactionIDGenerate(accountID)
	assert.NoError(t, err)

	recordQuery := NewTransactionRecordQuery().
		SetTransactionID(transactionID)

	err = recordQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
