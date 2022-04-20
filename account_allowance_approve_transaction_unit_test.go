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

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUnitAccountAllowanceApproveTransaction(t *testing.T) {
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
