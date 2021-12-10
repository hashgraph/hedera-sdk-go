//+build all unit

package hedera

import (
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
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

	transactionID := TransactionIDGenerate(AccountID{Account: 1111})

	tokenTransfer := NewTransferTransaction().
		AddHbarTransfer(accountID1, amount).
		AddHbarTransfer(accountID2, amount.Negated()).
		AddTokenTransfer(tokenID1, accountID1, 10).
		AddTokenTransfer(tokenID1, accountID2, -10).
		AddTokenTransfer(tokenID2, accountID1, 10).
		AddTokenTransfer(tokenID2, accountID2, -10).
		SetTransactionID(transactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		AddNftTransfer(nftID1, accountID1, accountID2).
		AddNftTransfer(nftID2, accountID2, accountID1).
		_Build()

	require.Equal(t, tokenTransfer.GetCryptoTransfer().Transfers.AccountAmounts, []*services.AccountAmount{
		{
			AccountID: accountID1._ToProtobuf(),
			Amount:    amount.AsTinybar(),
		},
		{
			AccountID: accountID2._ToProtobuf(),
			Amount:    amount.Negated().AsTinybar(),
		},
	})

	require.Equal(t, tokenTransfer.GetCryptoTransfer().TokenTransfers, []*services.TokenTransferList{
		{
			Token: tokenID1._ToProtobuf(),
			Transfers: []*services.AccountAmount{
				{
					AccountID: accountID1._ToProtobuf(),
					Amount:    10,
				},
				{
					AccountID: accountID2._ToProtobuf(),
					Amount:    -10,
				},
			}},
		{
			Token: tokenID2._ToProtobuf(),
			Transfers: []*services.AccountAmount{
				{
					AccountID: accountID1._ToProtobuf(),
					Amount:    10,
				},
				{
					AccountID: accountID2._ToProtobuf(),
					Amount:    -10,
				},
			}},
		{
			Token: tokenID3._ToProtobuf(),
			NftTransfers: []*services.NftTransfer{
				{
					SenderAccountID:   accountID1._ToProtobuf(),
					ReceiverAccountID: accountID2._ToProtobuf(),
					SerialNumber:      int64(9),
				},
			},
		},
		{
			Token: tokenID4._ToProtobuf(),
			NftTransfers: []*services.NftTransfer{
				{
					SenderAccountID:   accountID2._ToProtobuf(),
					ReceiverAccountID: accountID1._ToProtobuf(),
					SerialNumber:      int64(10),
				},
			},
		},
	})
}

func TestUnitTokenTransferTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	nftID, err := NftIDFromString("2@0.0.123-rmkyk")
	require.NoError(t, err)

	tokenTransfer := NewTransferTransaction().
		AddTokenTransfer(tokenID, accountID, 1).
		AddNftTransfer(nftID, accountID, accountID)

	err = tokenTransfer._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenTransferTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	nftID, err := NftIDFromString("2@0.0.123-rmkykd")
	require.NoError(t, err)

	tokenTransfer := NewTransferTransaction().
		AddTokenTransfer(tokenID, accountID, 1).
		AddNftTransfer(nftID, accountID, accountID)

	err = tokenTransfer._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
