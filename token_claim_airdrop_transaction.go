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
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TokenClaimAirdropTransaction struct {
	*Transaction[*TokenClaimAirdropTransaction]
	pendingAirdropIds []*PendingAirdropId
}

func NewTokenClaimAirdropTransaction() *TokenClaimAirdropTransaction {
	tx := &TokenClaimAirdropTransaction{
		pendingAirdropIds: make([]*PendingAirdropId, 0),
	}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(1))

	return tx
}

func _TokenClaimAirdropTransactionFromProtobuf(pb *services.TransactionBody) *TokenClaimAirdropTransaction {
	tokenClaimTransaction := &TokenClaimAirdropTransaction{}

	for _, pendingAirdrops := range pb.GetTokenClaimAirdrop().PendingAirdrops {
		tokenClaimTransaction.pendingAirdropIds = append(tokenClaimTransaction.pendingAirdropIds, _PendingAirdropIdFromProtobuf(pendingAirdrops))
	}

	return tokenClaimTransaction
}

// SetPendingAirdropIds sets the pending airdrop IDs for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) SetPendingAirdropIds(ids []*PendingAirdropId) *TokenClaimAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = ids
	return tx
}

// AddPendingAirdropId adds a pending airdrop ID to this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) AddPendingAirdropId(id PendingAirdropId) *TokenClaimAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = append(tx.pendingAirdropIds, &id)
	return tx
}

// GetPendingAirdropIds returns the pending airdrop IDs for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) GetPendingAirdropIds() []*PendingAirdropId {
	return tx.pendingAirdropIds
}

// ----------- Overridden functions ----------------

func (tx *TokenClaimAirdropTransaction) getName() string {
	return "TokenClaimAirdropTransaction"
}

func (tx *TokenClaimAirdropTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, pendingAirdropId := range tx.pendingAirdropIds {
		if pendingAirdropId.sender != nil {
			if err := pendingAirdropId.sender.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.receiver != nil {
			if err := pendingAirdropId.receiver.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.nftID != nil {
			if err := pendingAirdropId.nftID.Validate(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.tokenID != nil {
			if err := pendingAirdropId.tokenID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tx *TokenClaimAirdropTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenClaimAirdrop{
			TokenClaimAirdrop: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenClaimAirdropTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Data: &services.SchedulableTransactionBody_TokenClaimAirdrop{
			TokenClaimAirdrop: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenClaimAirdropTransaction) buildProtoBody() *services.TokenClaimAirdropTransactionBody {
	pendingAirdrops := make([]*services.PendingAirdropId, len(tx.pendingAirdropIds))
	for i, pendingAirdropId := range tx.pendingAirdropIds {
		pendingAirdrops[i] = pendingAirdropId._ToProtobuf()
	}

	return &services.TokenClaimAirdropTransactionBody{
		PendingAirdrops: pendingAirdrops,
	}
}

func (tx *TokenClaimAirdropTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().ClaimAirdrop,
	}
}

func (tx *TokenClaimAirdropTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *TokenClaimAirdropTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction)
}

func (tx *TokenClaimAirdropTransaction) setBaseTransaction(baseTx Transaction[TransactionInterface]) {
	tx.Transaction = castFromBaseToConcreteTransaction[*TokenClaimAirdropTransaction](baseTx)
}
