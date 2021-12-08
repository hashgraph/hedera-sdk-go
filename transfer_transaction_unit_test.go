//+build all unit

package hedera

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/require"
)

func TestUnitTransferTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	transfer := NewTransferTransaction().
		AddHbarTransfer(accountID, HbarFromTinybar(1))

	err = transfer._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransferTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	transfer := NewTransferTransaction().
		AddHbarTransfer(accountID, HbarFromTinybar(1))

	err = transfer._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestUnitTransferTransactionOrdered(t *testing.T) {
	tokenID1, err := TokenIDFromString("1.1.1")
	require.NoError(t, err)
	tokenID2, err := TokenIDFromString("2.2.2")
	require.NoError(t, err)
	tokenID3, err := TokenIDFromString("3.3.3")
	require.NoError(t, err)
	tokenID4, err := TokenIDFromString("4.4.4")
	require.NoError(t, err)
	serialNum1 := int64(111111111)
	accoundID1, err := AccountIDFromString("1.1.1")
	require.NoError(t, err)
	accoundID2, err := AccountIDFromString("2.2.2")
	require.NoError(t, err)
	accoundID3, err := AccountIDFromString("3.3.3")
	require.NoError(t, err)
	accoundID4, err := AccountIDFromString("4.4.4")
	require.NoError(t, err)

	expectedHbar := float64(3)
	accountID, err := AccountIDFromString("0.0.0")
	require.NoError(t, err)

	transfer := NewTransferTransaction().
		AddHbarTransfer(accountID, NewHbar(1)).
		AddHbarTransfer(accountID, NewHbar(1))

	transfer.AddHbarTransfer(accountID, NewHbar(1))
	transfer.AddHbarTransfer(AccountID{Account: 1}, NewHbar(1))

	require.Equal(t, transfer.hbarTransfers[accountID].As("hbar"), NewHbar(expectedHbar).As("hbar"))

	transferTransaction, err := NewTransferTransaction().
		AddNftTransfer(tokenID1.Nft(serialNum1), accoundID1, accoundID2).
		AddNftTransfer(tokenID1.Nft(serialNum1), accoundID1, accoundID2).
		SetTransactionID(NewTransactionIDWithValidStart(AccountID{Shard: 3, Realm: 3, Account: 3, checksum: nil}, time.Unix(4, 4))).
		SetNodeAccountIDs([]AccountID{accoundID4}).
		Freeze()
	require.NoError(t, err)

	transferTransactionToBytes, err := transferTransaction.ToBytes()
	require.NoError(t, err)

	transferTransactionFromBytes, err := TransactionFromBytes(transferTransactionToBytes)
	require.NoError(t, err)

	switch tx := transferTransactionFromBytes.(type) {
	case TransferTransaction:
		require.Equal(t, tx.nftTransfers[tokenID1], transferTransaction.nftTransfers[tokenID1])
	}

	transferTransaction = NewTransferTransaction().
		AddNftTransfer(tokenID4.Nft(serialNum1), accoundID2, accoundID4).
		AddNftTransfer(tokenID4.Nft(serialNum1), accoundID1, accoundID3).
		AddNftTransfer(tokenID3.Nft(serialNum1), accoundID1, accoundID2).
		AddTokenTransfer(tokenID2, accoundID4, -1).
		AddTokenTransfer(tokenID2, accoundID3, 2).
		AddTokenTransfer(tokenID1, accoundID2, -3).
		AddTokenTransfer(tokenID1, accoundID1, -4).
		AddHbarTransfer(accoundID2, NewHbar(-1)).
		AddHbarTransfer(accoundID1, NewHbar(1)).
		SetTransactionID(NewTransactionIDWithValidStart(accoundID3, time.Unix(4, 4))).
		SetNodeAccountIDs([]AccountID{accoundID4})

	data := transferTransaction._Build()

	switch d := data.Data.(type) {
	case *services.TransactionBody_CryptoTransfer:
		require.Equal(t, d.CryptoTransfer.Transfers.AccountAmounts, []*services.AccountAmount{
			{
				AccountID: accoundID1._ToProtobuf(),
				Amount:    int64(100000000),
			},
			{
				AccountID: accoundID2._ToProtobuf(),
				Amount:    int64(-100000000),
			},
		})

		require.Equal(t, d.CryptoTransfer.TokenTransfers, []*services.TokenTransferList{
			{
				Token: tokenID1._ToProtobuf(),
				Transfers: []*services.AccountAmount{
					{
						AccountID: accoundID1._ToProtobuf(),
						Amount:    int64(-4),
					},
					{
						AccountID: accoundID2._ToProtobuf(),
						Amount:    int64(-3),
					},
				}},
			{
				Token: tokenID2._ToProtobuf(),
				Transfers: []*services.AccountAmount{
					{
						AccountID: accoundID3._ToProtobuf(),
						Amount:    int64(2),
					},
					{
						AccountID: accoundID4._ToProtobuf(),
						Amount:    int64(-1)},
				}},
			{
				Token: tokenID3._ToProtobuf(),
				NftTransfers: []*services.NftTransfer{
					{
						SenderAccountID:   accoundID1._ToProtobuf(),
						ReceiverAccountID: accoundID2._ToProtobuf(),
						SerialNumber:      int64(111111111),
					},
				},
			},
			{
				Token: tokenID4._ToProtobuf(),
				NftTransfers: []*services.NftTransfer{
					{
						SenderAccountID:   accoundID1._ToProtobuf(),
						ReceiverAccountID: accoundID3._ToProtobuf(),
						SerialNumber:      int64(111111111),
					},
					{
						SenderAccountID:   accoundID2._ToProtobuf(),
						ReceiverAccountID: accoundID4._ToProtobuf(),
						SerialNumber:      int64(111111111),
					},
				},
			},
		})
	}
}
