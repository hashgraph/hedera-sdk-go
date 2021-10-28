//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenTransferTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	nftID, err := NftIDFromString("2@0.0.123-rmkyk")
	assert.NoError(t, err)

	tokenTransfer := NewTransferTransaction().
		AddTokenTransfer(tokenID, accountID, 1).
		AddNftTransfer(nftID, accountID, accountID)

	err = tokenTransfer._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTokenTransferTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)
	nftID, err := NftIDFromString("2@0.0.123-rmkykd")
	assert.NoError(t, err)

	tokenTransfer := NewTransferTransaction().
		AddTokenTransfer(tokenID, accountID, 1).
		AddNftTransfer(nftID, accountID, accountID)

	err = tokenTransfer._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
