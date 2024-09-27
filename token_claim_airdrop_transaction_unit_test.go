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

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitTokenClaimAirdropTransactionSetPendingAirdropIds(t *testing.T) {
	t.Parallel()

	pendingAirdropId1 := &PendingAirdropId{tokenID: &TokenID{Token: 1}}
	pendingAirdropId2 := &PendingAirdropId{tokenID: &TokenID{Token: 2}}

	transaction := NewTokenClaimAirdropTransaction().
		SetPendingAirdropIds([]*PendingAirdropId{pendingAirdropId1, pendingAirdropId2})

	assert.Equal(t, []*PendingAirdropId{pendingAirdropId1, pendingAirdropId2}, transaction.GetPendingAirdropIds())
}

func TestUnitTokenClaimAirdropTransactionAddPendingAirdropId(t *testing.T) {
	t.Parallel()

	pendingAirdropId1 := PendingAirdropId{tokenID: &TokenID{Token: 1}}
	pendingAirdropId2 := PendingAirdropId{tokenID: &TokenID{Token: 2}}

	transaction := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(pendingAirdropId1).
		AddPendingAirdropId(pendingAirdropId2)

	assert.Equal(t, []*PendingAirdropId{&pendingAirdropId1, &pendingAirdropId2}, transaction.GetPendingAirdropIds())
}

func TestUnitTokenClaimAirdropTransactionFreeze(t *testing.T) {
	t.Parallel()

	pendingAirdropId := PendingAirdropId{tokenID: &TokenID{Token: 1}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})
	transaction := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(pendingAirdropId).
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID)

	_, err := transaction.Freeze()
	require.NoError(t, err)
}

func TestUnitTokenClaimAirdropTransactionToBytes(t *testing.T) {
	t.Parallel()

	pendingAirdropId := PendingAirdropId{tokenID: &TokenID{Token: 1}, sender: &AccountID{Account: 3}}

	transaction := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(pendingAirdropId)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)
}

func TestUnitTokenClaimAirdropTransactionFromBytes(t *testing.T) {
	t.Parallel()

	pendingAirdropId := PendingAirdropId{tokenID: &TokenID{Token: 1}}

	transaction := NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(pendingAirdropId)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	deserializedTransaction, err := TransactionFromBytes(bytes)
	require.NoError(t, err)

	switch tx := deserializedTransaction.(type) {
	case *TokenClaimAirdropTransaction:
		assert.Equal(t, transaction.GetPendingAirdropIds(), tx.GetPendingAirdropIds())
	default:
		t.Fatalf("expected TokenClaimAirdropTransaction, got %T", deserializedTransaction)
	}
}

func TestUnitTokenClaimAirdropTransactionScheduleProtobuf(t *testing.T) {
	t.Parallel()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	pendingAirdropId1 := PendingAirdropId{tokenID: &TokenID{Token: 1}}
	pendingAirdropId2 := PendingAirdropId{tokenID: &TokenID{Token: 2}}

	tx, err := NewTokenClaimAirdropTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		AddPendingAirdropId(pendingAirdropId1).
		AddPendingAirdropId(pendingAirdropId2).
		Freeze()
	require.NoError(t, err)

	expected := &services.SchedulableTransactionBody{
		TransactionFee: 100000000,
		Data: &services.SchedulableTransactionBody_TokenClaimAirdrop{
			TokenClaimAirdrop: &services.TokenClaimAirdropTransactionBody{
				PendingAirdrops: []*services.PendingAirdropId{
					pendingAirdropId1._ToProtobuf(),
					pendingAirdropId2._ToProtobuf(),
				},
			},
		},
	}

	actual, err := tx.buildScheduled()
	require.NoError(t, err)
	require.Equal(t, expected.String(), actual.String())
}

func TestUnitTokenClaimAirdropTransactionValidateNetworkOnIDs(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	checksum := "dmqui"
	pendingAirdropId := &PendingAirdropId{
		tokenID:  &TokenID{Token: 3, checksum: &checksum},
		sender:   &AccountID{Account: 3, checksum: &checksum},
		receiver: &AccountID{Account: 3, checksum: &checksum},
	}

	transaction := NewTokenClaimAirdropTransaction().
		SetPendingAirdropIds([]*PendingAirdropId{pendingAirdropId})

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
}
