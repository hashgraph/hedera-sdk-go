//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenTransferTransactionTransfers(t *testing.T) {
	amount := NewHbar(1)
	accountID1 := AccountID{Account: 3}
	accountID2 := AccountID{Account: 4}
	tokenID1 := TokenID{Token: 5}
	tokenID2 := TokenID{Token: 6}
	tokenID3 := TokenID{Token: 7}
	tokenID4 := TokenID{Token: 8}
	nftID1 := NftID{TokenID: tokenID3, SerialNumber: 9}
	nftID2 := NftID{TokenID: tokenID4, SerialNumber: 10}

	tokenTransfer := NewTransferTransaction().
		AddHbarTransfer(accountID1, amount).
		AddHbarTransfer(accountID2, amount.Negated()).
		AddTokenTransfer(tokenID1, accountID1, 10).
		AddTokenTransfer(tokenID1, accountID2, -10).
		AddTokenTransfer(tokenID2, accountID1, 10).
		AddTokenTransfer(tokenID2, accountID2, -10).
		AddNftTransfer(nftID1, accountID1, accountID2).
		AddNftTransfer(nftID2, accountID2, accountID1)

	hbarTransfers := tokenTransfer.GetHbarTransfers()
	tokenTransfers := tokenTransfer.GetTokenTransfers()
	nftTransfers := tokenTransfer.GetNftTransfers()

	assert.Equal(t, len(hbarTransfers), 2)
	assert.Equal(t, hbarTransfers[accountID1], amount)
	assert.Equal(t, hbarTransfers[accountID2], amount.Negated())

	assert.Equal(t, len(tokenTransfers), 2)
	assert.Equal(t, len(tokenTransfers[tokenID1]), 2)
	assert.Equal(t, tokenTransfers[tokenID1][0], TokenTransfer{AccountID: accountID1, Amount: 10})
	assert.Equal(t, tokenTransfers[tokenID1][1], TokenTransfer{AccountID: accountID2, Amount: -10})
	assert.Equal(t, len(tokenTransfers[tokenID2]), 2)
	assert.Equal(t, tokenTransfers[tokenID2][0], TokenTransfer{AccountID: accountID1, Amount: 10})
	assert.Equal(t, tokenTransfers[tokenID2][1], TokenTransfer{AccountID: accountID2, Amount: -10})

	assert.Equal(t, len(nftTransfers), 2)
	assert.Equal(t, len(nftTransfers[tokenID3]), 1)
	assert.Equal(t, nftTransfers[tokenID3][0], TokenNftTransfer{SenderAccountID: accountID1, ReceiverAccountID: accountID2, SerialNumber: 9})
	assert.Equal(t, len(nftTransfers[tokenID4]), 1)
	assert.Equal(t, nftTransfers[tokenID4][0], TokenNftTransfer{SenderAccountID: accountID2, ReceiverAccountID: accountID1, SerialNumber: 10})
}

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
