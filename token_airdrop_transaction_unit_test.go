//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"testing"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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

func TestUnitTokenAirdropTransactionOrdered(t *testing.T) {
	t.Parallel()

	tokenID1, err := TokenIDFromString("1.1.1")
	require.NoError(t, err)
	tokenID2, err := TokenIDFromString("2.2.2")
	require.NoError(t, err)
	tokenID3, err := TokenIDFromString("3.3.3")
	require.NoError(t, err)
	tokenID4, err := TokenIDFromString("4.4.4")
	require.NoError(t, err)
	serialNum1 := int64(111111111)
	accountID1, err := AccountIDFromString("1.1.1")
	require.NoError(t, err)
	accountID2, err := AccountIDFromString("2.2.2")
	require.NoError(t, err)
	accountID3, err := AccountIDFromString("3.3.3")
	require.NoError(t, err)
	accountID4, err := AccountIDFromString("4.4.4")
	require.NoError(t, err)

	airdropTransaction, err := NewTokenAirdropTransaction().
		AddNftTransfer(tokenID1.Nft(serialNum1), accountID1, accountID2).
		AddNftTransfer(tokenID1.Nft(serialNum1), accountID1, accountID2).
		SetTransactionID(NewTransactionIDWithValidStart(AccountID{Shard: 3, Realm: 3, Account: 3, checksum: nil}, time.Unix(4, 4))).
		SetNodeAccountIDs([]AccountID{accountID4}).
		Freeze()
	require.NoError(t, err)

	airdropTransactionToBytes, err := airdropTransaction.ToBytes()
	require.NoError(t, err)

	airdropTransactionFromBytes, err := TransactionFromBytes(airdropTransactionToBytes)
	require.NoError(t, err)

	switch tx := airdropTransactionFromBytes.(type) {
	case TokenAirdropTransaction:
		require.Equal(t, tx.nftTransfers[tokenID1], airdropTransaction.nftTransfers[tokenID1])
	}

	airdropTransaction = NewTokenAirdropTransaction().
		AddNftTransfer(tokenID4.Nft(serialNum1), accountID2, accountID4).
		AddNftTransfer(tokenID4.Nft(serialNum1), accountID1, accountID3).
		AddNftTransfer(tokenID3.Nft(serialNum1), accountID1, accountID2).
		AddTokenTransfer(tokenID2, accountID4, -1).
		AddTokenTransfer(tokenID2, accountID3, 2).
		AddTokenTransfer(tokenID1, accountID2, -3).
		AddTokenTransfer(tokenID1, accountID1, -4).
		SetTransactionID(NewTransactionIDWithValidStart(accountID3, time.Unix(4, 4))).
		SetNodeAccountIDs([]AccountID{accountID4})

	data := airdropTransaction.build()

	switch d := data.Data.(type) {
	case *services.TransactionBody_CryptoTransfer:
		require.Equal(t, d.CryptoTransfer.Transfers.AccountAmounts, []*services.AccountAmount{
			{
				AccountID: accountID1._ToProtobuf(),
				Amount:    int64(100000000),
			},
			{
				AccountID: accountID2._ToProtobuf(),
				Amount:    int64(-100000000),
			},
		})

		require.Equal(t, d.CryptoTransfer.TokenTransfers, []*services.TokenTransferList{
			{
				Token: tokenID1._ToProtobuf(),
				Transfers: []*services.AccountAmount{
					{
						AccountID: accountID1._ToProtobuf(),
						Amount:    int64(-4),
					},
					{
						AccountID: accountID2._ToProtobuf(),
						Amount:    int64(-3),
					},
				}},
			{
				Token: tokenID2._ToProtobuf(),
				Transfers: []*services.AccountAmount{
					{
						AccountID: accountID3._ToProtobuf(),
						Amount:    int64(2),
					},
					{
						AccountID: accountID4._ToProtobuf(),
						Amount:    int64(-1),
					},
				}},
			{
				Token: tokenID3._ToProtobuf(),
				NftTransfers: []*services.NftTransfer{
					{
						SenderAccountID:   accountID1._ToProtobuf(),
						ReceiverAccountID: accountID2._ToProtobuf(),
						SerialNumber:      int64(111111111),
					},
				},
			},
			{
				Token: tokenID4._ToProtobuf(),
				NftTransfers: []*services.NftTransfer{
					{
						SenderAccountID:   accountID1._ToProtobuf(),
						ReceiverAccountID: accountID3._ToProtobuf(),
						SerialNumber:      int64(111111111),
					},
					{
						SenderAccountID:   accountID2._ToProtobuf(),
						ReceiverAccountID: accountID4._ToProtobuf(),
						SerialNumber:      int64(111111111),
					},
				},
			},
		})
	}
}
