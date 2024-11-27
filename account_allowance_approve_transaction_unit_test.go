//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var tokenID1 = TokenID{Token: 1}
var tokenID2 = TokenID{Token: 141}
var serialNumber1 = int64(3)
var serialNumber2 = int64(4)
var nftID1 = tokenID2.Nft(serialNumber1)
var nftID2 = tokenID2.Nft(serialNumber2)
var owner = AccountID{Account: 10}
var spenderAccountID1 = AccountID{Account: 7}
var spenderAccountID2 = AccountID{Account: 7890}
var nodeAccountID = []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
var hbarAmount = HbarFromTinybar(100)
var tokenAmount = int64(101)

func TestUnitAccountAllowanceApproveTransaction(t *testing.T) {
	t.Parallel()

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

	data := transaction.build()

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
func TestUnitvalidateNetworkOnIDs(t *testing.T) {
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

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)

	e := transaction.validateNetworkOnIDs(client)
	require.NoError(t, e)
}
func TestUnitAccountAllowanceApproveTransactionGet(t *testing.T) {
	t.Parallel()

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

func TestUnitAccountAllowanceApproveTransactionFromProtobuf(t *testing.T) {
	t.Parallel()

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

	txFromProto := _AccountAllowanceApproveTransactionFromProtobuf(*tx.Transaction, tx.build())
	require.Equal(t, tx, &txFromProto)
}

func TestUnitAccountAllowanceApproveTransactionScheduleProtobuf(t *testing.T) {
	t.Parallel()

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
	actual, err := tx.buildScheduled()
	require.NoError(t, err)
	require.Equal(t, expected.String(), actual.String())
}
