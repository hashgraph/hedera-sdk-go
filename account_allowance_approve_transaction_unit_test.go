//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUnitAccountAllowanceApproveTransaction(t *testing.T) {
	t.Parallel()

	tokenID1 := TokenID{Token: 1}
	tokenID2 := TokenID{Token: 141}
	serialNumber1 := int64(3)
	serialNumber2 := int64(4)
	nftID1 := tokenID2.Nft(serialNumber1)
	nftID2 := tokenID2.Nft(serialNumber2)
	owner := AccountID{Account: 10}
	spenderAccountID1 := AccountID{Account: 7}
	spenderAccountID2 := AccountID{Account: 7890}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	hbarAmount := HbarFromTinybar(100)
	tokenAmount := int64(101)

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountAllowanceApproveTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		ApproveHbarAllowance(owner, spenderAccountID1, hbarAmount).
		ApproveTokenAllowance(tokenID1, owner, spenderAccountID1, tokenAmount).
		ApproveTokenNftAllowance(nftID1, owner, spenderAccountID1).
		ApproveTokenNftAllowance(nftID2, owner, spenderAccountID1).
		ApproveTokenNftAllowance(nftID2, owner, spenderAccountID2).
		AddAllTokenNftApproval(tokenID1, spenderAccountID1).
		Freeze()
	require.NoError(t, err)

	data := transaction._Build()

	switch d := data.Data.(type) {
	case *services.TransactionBody_CryptoApproveAllowance:
		require.Equal(t, d.CryptoApproveAllowance.CryptoAllowances, []*services.CryptoAllowance{
			{
				Spender: spenderAccountID1._ToProtobuf(),
				Owner:   owner._ToProtobuf(),
				Amount:  hbarAmount.AsTinybar(),
			},
		})
		require.Equal(t, d.CryptoApproveAllowance.NftAllowances, []*services.NftAllowance{
			{
				TokenId:           tokenID2._ToProtobuf(),
				Spender:           spenderAccountID1._ToProtobuf(),
				Owner:             owner._ToProtobuf(),
				SerialNumbers:     []int64{serialNumber1, serialNumber2},
				ApprovedForAll:    &wrapperspb.BoolValue{Value: false},
				DelegatingSpender: nil,
			},
			{
				TokenId:           tokenID2._ToProtobuf(),
				Spender:           spenderAccountID2._ToProtobuf(),
				Owner:             owner._ToProtobuf(),
				SerialNumbers:     []int64{serialNumber2},
				ApprovedForAll:    &wrapperspb.BoolValue{Value: false},
				DelegatingSpender: nil,
			},
			{
				TokenId:           tokenID1._ToProtobuf(),
				Spender:           spenderAccountID1._ToProtobuf(),
				Owner:             nil,
				SerialNumbers:     []int64{},
				ApprovedForAll:    &wrapperspb.BoolValue{Value: true},
				DelegatingSpender: nil,
			},
		})
		require.Equal(t, d.CryptoApproveAllowance.TokenAllowances, []*services.TokenAllowance{
			{
				TokenId: tokenID1._ToProtobuf(),
				Owner:   owner._ToProtobuf(),
				Spender: spenderAccountID1._ToProtobuf(),
				Amount:  tokenAmount,
			},
		})
	}
}

func TestUnitAccountAllowanceApproveTransactionGet(t *testing.T) {
	t.Parallel()

	tokenID1 := TokenID{Token: 1}
	tokenID2 := TokenID{Token: 141}
	serialNumber1 := int64(3)
	serialNumber2 := int64(4)
	nftID1 := tokenID2.Nft(serialNumber1)
	nftID2 := tokenID2.Nft(serialNumber2)
	owner := AccountID{Account: 10}
	spenderAccountID1 := AccountID{Account: 7}
	spenderAccountID2 := AccountID{Account: 7890}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	hbarAmount := HbarFromTinybar(100)
	tokenAmount := int64(101)

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountAllowanceApproveTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		ApproveHbarAllowance(owner, spenderAccountID1, hbarAmount).
		ApproveTokenAllowance(tokenID1, owner, spenderAccountID1, tokenAmount).
		ApproveTokenNftAllowance(nftID1, owner, spenderAccountID1).
		ApproveTokenNftAllowance(nftID2, owner, spenderAccountID1).
		ApproveTokenNftAllowance(nftID2, owner, spenderAccountID2).
		AddAllTokenNftApproval(tokenID1, spenderAccountID1).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)
	transaction.GetTokenNftAllowances()
	transaction.GetHbarAllowances()
	transaction.GetTokenAllowances()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
}

func TestUnitAccountAllowanceApproveTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountAllowanceApproveTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)
	transaction.GetTokenNftAllowances()
	transaction.GetHbarAllowances()
	transaction.GetTokenAllowances()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
}

func TestUnitAccountAllowanceDeleteTransactionSetNothing(t *testing.T) {
	t.Parallel()

	token := TokenID{Token: 3}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountAllowanceDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		DeleteAllTokenNftAllowances(token.Nft(4), &AccountID{Account: 3}).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)
	transaction.GetAllTokenNftDeleteAllowances()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
}

func TestUnitAccountAllowanceApproveTransactionFromProtobuf(t *testing.T) {
	t.Parallel()

	tokenID1 := TokenID{Token: 1}
	tokenID2 := TokenID{Token: 141}
	serialNumber1 := int64(3)
	serialNumber2 := int64(4)
	nftID1 := tokenID2.Nft(serialNumber1)
	nftID2 := tokenID2.Nft(serialNumber2)
	owner := AccountID{Account: 10}
	spenderAccountID1 := AccountID{Account: 7}
	spenderAccountID2 := AccountID{Account: 7890}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	hbarAmount := HbarFromTinybar(100)
	tokenAmount := int64(101)

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	tx, err := NewAccountAllowanceApproveTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		ApproveHbarAllowance(owner, spenderAccountID1, hbarAmount).
		ApproveTokenAllowance(tokenID1, owner, spenderAccountID1, tokenAmount).
		ApproveTokenNftAllowance(nftID1, owner, spenderAccountID1).
		ApproveTokenNftAllowance(nftID2, owner, spenderAccountID1).
		ApproveTokenNftAllowance(nftID2, owner, spenderAccountID2).
		Freeze()
	require.NoError(t, err)

	txFromProto := _AccountAllowanceApproveTransactionFromProtobuf(tx.Transaction, tx._Build())
	require.Equal(t, tx, txFromProto)
}

func TestUnitAccountAllowanceApproveTransactionScheduleProtobuf(t *testing.T) {
	t.Parallel()

	tokenID1 := TokenID{Token: 1}
	tokenID2 := TokenID{Token: 141}
	serialNumber1 := int64(3)
	serialNumber2 := int64(4)
	nftID1 := tokenID2.Nft(serialNumber1)
	nftID2 := tokenID2.Nft(serialNumber2)
	owner := AccountID{Account: 10}
	spenderAccountID1 := AccountID{Account: 7}
	spenderAccountID2 := AccountID{Account: 7890}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	hbarAmount := HbarFromTinybar(100)
	tokenAmount := int64(101)

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	tx, err := NewAccountAllowanceApproveTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		ApproveHbarAllowance(owner, spenderAccountID1, hbarAmount).
		ApproveTokenAllowance(tokenID1, owner, spenderAccountID1, tokenAmount).
		ApproveTokenNftAllowance(nftID1, owner, spenderAccountID1).
		ApproveTokenNftAllowance(nftID2, owner, spenderAccountID1).
		ApproveTokenNftAllowance(nftID2, owner, spenderAccountID2).
		Freeze()
	require.NoError(t, err)

	expected := &services.SchedulableTransactionBody{
		TransactionFee: 200000000,
		Data: &services.SchedulableTransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: &services.CryptoApproveAllowanceTransactionBody{
				CryptoAllowances: []*services.CryptoAllowance{
					{
						Owner:   &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 10}},
						Spender: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 7}},
						Amount:  100,
					},
				},
				NftAllowances: []*services.NftAllowance{
					{
						TokenId:           &services.TokenID{TokenNum: 141},
						Owner:             &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 10}},
						Spender:           &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 7}},
						SerialNumbers:     []int64{3, 4},
						ApprovedForAll:    &wrapperspb.BoolValue{Value: false},
						DelegatingSpender: nil,
					},
					{
						TokenId:           &services.TokenID{TokenNum: 141},
						Owner:             &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 10}},
						Spender:           &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 7890}},
						SerialNumbers:     []int64{4},
						ApprovedForAll:    &wrapperspb.BoolValue{Value: false},
						DelegatingSpender: nil,
					},
				},
				TokenAllowances: []*services.TokenAllowance{
					{
						TokenId: &services.TokenID{TokenNum: 1},
						Owner:   &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 10}},
						Spender: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 7}},
						Amount:  101,
					},
				},
			},
		},
	}
	actual, err := tx._ConstructScheduleProtobuf()
	require.NoError(t, err)
	require.Equal(t, expected.String(), actual.String())
}

func TestUnitAccountAllowanceDeleteTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	token := TokenID{Token: 3, checksum: &checksum}
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewAccountAllowanceDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		DeleteAllTokenNftAllowances(token.Nft(4), &account).
		DeleteAllTokenNftAllowances(token.Nft(5), &account).
		SetMaxTransactionFee(NewHbar(3)).
		SetMaxRetry(3).
		DeleteAllHbarAllowances(&account).
		DeleteAllTokenAllowances(token, &account).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetTransactionMemo("no").
		SetTransactionValidDuration(time.Second * 30).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction._ValidateNetworkOnIDs(client)

	_, err = transaction.Schedule()
	require.NoError(t, err)
	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()
	transaction.GetMaxRetry()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxBackoff()
	transaction.GetMinBackoff()
	transaction.GetRegenerateTransactionID()
	byt, err := transaction.ToBytes()
	require.NoError(t, err)
	txFromBytesI, err := TransactionFromBytes(byt)
	require.NoError(t, err)
	txFromBytes, ok := txFromBytesI.(AccountAllowanceDeleteTransaction)
	require.Equal(t, true, ok)
	sig, err := newKey.SignTransaction(&transaction.Transaction)
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)
	transaction.GetAllTokenNftDeleteAllowances()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetAllHbarDeleteAllowances()
	transaction.GetAllTokenDeleteAllowances()
	transaction._GetLogID()
	txFromBytes.AddSignature(newKey.PublicKey(), sig)
}

func TestUnitAccountAllowanceDeleteTransactionMock(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	token := TokenID{Token: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 3}}
	transactionID := TransactionIDGenerate(AccountID{Account: 3})

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

	_, err := NewAccountAllowanceDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		DeleteAllTokenNftAllowances(token.Nft(4), &AccountID{Account: 3}).
		Execute(client)
	require.NoError(t, err)
}
