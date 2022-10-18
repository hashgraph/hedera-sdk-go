//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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

	protobuf "google.golang.org/protobuf/proto"

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
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	nftID, err := NftIDFromString("2@0.0.123-esxsf")
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
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTransferTransactionGet(t *testing.T) {
	tokenID := TokenID{Token: 7}
	accountID := AccountID{Account: 3}
	nftID := tokenID.Nft(32)

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTransferTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetHbarTransferApproval(accountID, true).
		SetTokenTransferApproval(tokenID, accountID, true).
		SetNftTransferApproval(nftID, true).
		AddHbarTransfer(accountID, NewHbar(34)).
		AddTokenTransferWithDecimals(tokenID, accountID, 123, 12).
		AddTokenTransfer(tokenID, accountID, 123).
		AddNftTransfer(nftID, accountID, accountID).
		SetMaxTransactionFee(NewHbar(10)).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetHbarTransfers()
	transaction.GetTokenTransfers()
	transaction.GetNftTransfers()
	transaction.GetTokenIDDecimals()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTransferTransactionNothingSet(t *testing.T) {
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTransferTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetHbarTransfers()
	transaction.GetTokenTransfers()
	transaction.GetNftTransfers()
	transaction.GetTokenIDDecimals()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTransferTransactionMock(t *testing.T) {
	newKey, err := PrivateKeyFromStringEd25519("302e020100300506032b657004220420a869f4c6191b9c8c99933e7f6b6611711737e4b1a1a5a4cb5370e719a1f6df98")
	require.NoError(t, err)

	call := func(request *services.Transaction) *services.TransactionResponse {
		require.NotEmpty(t, request.SignedTransactionBytes)
		signedTransaction := services.SignedTransaction{}
		_ = protobuf.Unmarshal(request.SignedTransactionBytes, &signedTransaction)

		require.NotEmpty(t, signedTransaction.BodyBytes)
		transactionBody := services.TransactionBody{}
		_ = protobuf.Unmarshal(signedTransaction.BodyBytes, &transactionBody)

		require.NotNil(t, transactionBody.TransactionID)
		transactionId := transactionBody.TransactionID.String()
		require.NotEqual(t, "", transactionId)

		sigMap := signedTransaction.GetSigMap()
		require.NotNil(t, sigMap)

		return &services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		}
	}
	responses := [][]interface{}{{
		call,
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	freez, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}
