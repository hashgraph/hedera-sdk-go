//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitTokenAirdropTransactionSetTokenTransferWithDecimals(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	senderAccountID := AccountID{Account: 2}
	amount := int64(10)
	decimals := uint32(5)

	transaction := NewTokenAirdropTransaction().
		AddTokenTransferWithDecimals(tokenID, senderAccountID, amount, decimals)

	require.Equal(t, transaction.GetTokenIDDecimals()[tokenID], decimals)
}

func TestUnitTokenAirdropTransactionAddApprovedTokenTransferWithDecimals(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	accountID := AccountID{Account: 2}
	amount := int64(100)
	decimals := uint32(5)

	transaction := NewTokenAirdropTransaction().
		AddApprovedTokenTransferWithDecimals(tokenID, accountID, amount, decimals, true)

	transfers := transaction.GetTokenTransfers()
	require.NotNil(t, transfers)
	require.Contains(t, transfers, tokenID)
	require.Len(t, transfers[tokenID], 1)
	assert.Equal(t, accountID, transfers[tokenID][0].AccountID)
	assert.Equal(t, amount, transfers[tokenID][0].Amount)
	assert.True(t, transfers[tokenID][0].IsApproved)
	assert.Equal(t, decimals, transaction.GetTokenIDDecimals()[tokenID])
}

func TestUnitTokenAirdropTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	nodeAccountIDs := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 123})

	checksum := "dmqui"

	token := TokenID{Token: 3, checksum: &checksum}
	nft := NftID{TokenID: TokenID{Token: 3, checksum: &checksum}, SerialNumber: 1}
	airdrop := NewTokenAirdropTransaction().
		AddTokenTransfer(token, accountID, 100).
		AddNftTransfer(nft, accountID, accountID).
		SetTransactionID(transactionID).SetNodeAccountIDs(nodeAccountIDs).
		SetMaxTransactionFee(HbarFromTinybar(100)).SetRegenerateTransactionID(true).
		SetTransactionMemo("go sdk unit test").SetTransactionValidDuration(time.Second * 120).
		SetMaxRetry(1).SetMaxBackoff(time.Second * 120).SetMinBackoff(time.Second * 1)

	err = airdrop.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenAirdropTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	airdrop := NewTokenAirdropTransaction().
		AddTokenTransfer(TokenID{Token: 1}, accountID, 100)

	err = airdrop.validateNetworkOnIDs(client)
	require.Error(t, err)
}

func TestUnitTokenAirdropTransactionAddTokenTransfer(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	accountID := AccountID{Account: 2}
	amount := int64(100)

	transaction := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, accountID, amount)

	transfers := transaction.GetTokenTransfers()
	require.NotNil(t, transfers)
	require.Contains(t, transfers, tokenID)
	require.Len(t, transfers[tokenID], 1)
	assert.Equal(t, accountID, transfers[tokenID][0].AccountID)
	assert.Equal(t, amount, transfers[tokenID][0].Amount)
}

func TestUnitTokenAirdropTransactionAddNftTransfer(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	serialNumber := int64(1)
	nftID := NftID{TokenID: tokenID, SerialNumber: serialNumber}
	sender := AccountID{Account: 2}
	receiver := AccountID{Account: 3}

	transaction := NewTokenAirdropTransaction().
		AddNftTransfer(nftID, sender, receiver)

	nftTransfers := transaction.GetNftTransfers()
	require.NotNil(t, nftTransfers)
	require.Contains(t, nftTransfers, tokenID)
	require.Len(t, nftTransfers[tokenID], 1)
	assert.Equal(t, sender, nftTransfers[tokenID][0].SenderAccountID)
	assert.Equal(t, receiver, nftTransfers[tokenID][0].ReceiverAccountID)
	assert.Equal(t, serialNumber, nftTransfers[tokenID][0].SerialNumber)
}

func TestUnitTokenAirdropTransactionSetTokenTransferApproval(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	accountID := AccountID{Account: 2}
	amount := int64(100)

	transaction := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, accountID, amount).
		SetTokenTransferApproval(tokenID, accountID, true)

	transfers := transaction.GetTokenTransfers()
	require.NotNil(t, transfers)
	require.Contains(t, transfers, tokenID)
	require.Len(t, transfers[tokenID], 1)
	assert.True(t, transfers[tokenID][0].IsApproved)
}

func TestUnitTokenAirdropTransactionSetNftTransferApproval(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	serialNumber := int64(1)
	nftID := NftID{TokenID: tokenID, SerialNumber: serialNumber}
	sender := AccountID{Account: 2}
	receiver := AccountID{Account: 3}

	transaction := NewTokenAirdropTransaction().
		AddNftTransfer(nftID, sender, receiver).
		SetNftTransferApproval(nftID, true)

	nftTransfers := transaction.GetNftTransfers()
	require.NotNil(t, nftTransfers)
	require.Contains(t, nftTransfers, tokenID)
	require.Len(t, nftTransfers[tokenID], 1)
	assert.True(t, nftTransfers[tokenID][0].IsApproved)
}

func TestUnitTokenAirdropTransactionAddApprovedTokenTransfer(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	accountID := AccountID{Account: 2}
	amount := int64(100)

	transaction := NewTokenAirdropTransaction().
		AddApprovedTokenTransfer(tokenID, accountID, amount, true)

	transfers := transaction.GetTokenTransfers()
	require.NotNil(t, transfers)
	require.Contains(t, transfers, tokenID)
	require.Len(t, transfers[tokenID], 1)
	assert.Equal(t, accountID, transfers[tokenID][0].AccountID)
	assert.Equal(t, amount, transfers[tokenID][0].Amount)
	assert.True(t, transfers[tokenID][0].IsApproved)
}

func TestUnitTokenAirdropTransactionAddApprovedNftTransfer(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	serialNumber := int64(1)
	nftID := NftID{TokenID: tokenID, SerialNumber: serialNumber}
	sender := AccountID{Account: 2}
	receiver := AccountID{Account: 3}

	transaction := NewTokenAirdropTransaction().
		AddApprovedNftTransfer(nftID, sender, receiver, true)

	nftTransfers := transaction.GetNftTransfers()
	require.NotNil(t, nftTransfers)
	require.Contains(t, nftTransfers, tokenID)
	require.Len(t, nftTransfers[tokenID], 1)
	assert.Equal(t, sender, nftTransfers[tokenID][0].SenderAccountID)
	assert.Equal(t, receiver, nftTransfers[tokenID][0].ReceiverAccountID)
	assert.Equal(t, serialNumber, nftTransfers[tokenID][0].SerialNumber)
	assert.True(t, nftTransfers[tokenID][0].IsApproved)
}

func TestUnitTokenAirdropTransactionToBytes(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	accountID := AccountID{Account: 2}
	amount := int64(100)

	transaction := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, accountID, amount)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)
}

func TestUnitTokenAirdropTransactionFromBytes(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	accountID := AccountID{Account: 2}
	amount := int64(100)

	transaction := NewTokenAirdropTransaction().
		AddTokenTransfer(tokenID, accountID, amount)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	deserializedTransaction, err := TransactionFromBytes(bytes)
	require.NoError(t, err)

	switch tx := deserializedTransaction.(type) {
	case TokenAirdropTransaction:
		assert.Equal(t, transaction.GetTokenTransfers(), tx.GetTokenTransfers())
	default:
		t.Fatalf("expected TokenAirdropTransaction, got %T", deserializedTransaction)
	}
}

func TestUnitTokenAirdropTransactionScheduleProtobuf(t *testing.T) {
	t.Parallel()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	tokenID1 := TokenID{Token: 1}
	tokenID2 := TokenID{Token: 2}
	accountID1 := AccountID{Account: 1}
	accountID2 := AccountID{Account: 2}
	amount1 := int64(100)
	amount2 := int64(200)
	serialNumber1 := int64(1)
	serialNumber2 := int64(2)
	nodeAccountID := []AccountID{{Account: 10}}

	tx, err := NewTokenAirdropTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		AddTokenTransfer(tokenID1, accountID1, amount1).
		AddTokenTransfer(tokenID2, accountID2, amount2).
		AddNftTransfer(tokenID1.Nft(serialNumber1), accountID1, accountID2).
		AddNftTransfer(tokenID2.Nft(serialNumber2), accountID2, accountID1).
		Freeze()
	require.NoError(t, err)

	expected := &services.SchedulableTransactionBody{
		TransactionFee: 100000000,
		Data: &services.SchedulableTransactionBody_TokenAirdrop{
			TokenAirdrop: &services.TokenAirdropTransactionBody{
				TokenTransfers: []*services.TokenTransferList{
					{
						Token: tokenID1._ToProtobuf(),
						Transfers: []*services.AccountAmount{
							{
								AccountID: accountID1._ToProtobuf(),
								Amount:    amount1,
							},
						},
					},
					{
						Token: tokenID2._ToProtobuf(),
						Transfers: []*services.AccountAmount{
							{
								AccountID: accountID2._ToProtobuf(),
								Amount:    amount2,
							},
						},
					},
					{
						Token: tokenID1._ToProtobuf(),
						NftTransfers: []*services.NftTransfer{
							{
								SenderAccountID:   accountID1._ToProtobuf(),
								ReceiverAccountID: accountID2._ToProtobuf(),
								SerialNumber:      serialNumber1,
							},
						},
					},
					{
						Token: tokenID2._ToProtobuf(),
						NftTransfers: []*services.NftTransfer{
							{
								SenderAccountID:   accountID2._ToProtobuf(),
								ReceiverAccountID: accountID1._ToProtobuf(),
								SerialNumber:      serialNumber2,
							},
						},
					},
				},
			},
		},
	}

	actual, err := tx.buildScheduled()
	require.NoError(t, err)
	require.ElementsMatch(t, expected.GetTokenAirdrop().TokenTransfers, actual.GetTokenAirdrop().TokenTransfers)
}
