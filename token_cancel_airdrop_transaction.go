package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type TokenCancelAirdropTransaction struct {
	*Transaction[*TokenCancelAirdropTransaction]
	pendingAirdropIds []*PendingAirdropId
}

func NewTokenCancelAirdropTransaction() *TokenCancelAirdropTransaction {
	tx := &TokenCancelAirdropTransaction{
		pendingAirdropIds: make([]*PendingAirdropId, 0),
	}

	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(1))

	return tx
}

func _TokenCancelAirdropTransactionFromProtobuf(tx Transaction[*TokenCancelAirdropTransaction], pb *services.TransactionBody) TokenCancelAirdropTransaction {
	tokenCancelTransaction := TokenCancelAirdropTransaction{}

	for _, pendingAirdrops := range pb.GetTokenCancelAirdrop().PendingAirdrops {
		tokenCancelTransaction.pendingAirdropIds = append(tokenCancelTransaction.pendingAirdropIds, _PendingAirdropIdFromProtobuf(pendingAirdrops))
	}

	tx.childTransaction = &tokenCancelTransaction
	tokenCancelTransaction.Transaction = &tx
	return tokenCancelTransaction
}

// SetPendingAirdropIds sets the pending airdrop IDs for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) SetPendingAirdropIds(ids []*PendingAirdropId) *TokenCancelAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = ids
	return tx
}

// AddPendingAirdropId adds a pending airdrop ID to this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) AddPendingAirdropId(id PendingAirdropId) *TokenCancelAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = append(tx.pendingAirdropIds, &id)
	return tx
}

// GetPendingAirdropIds returns the pending airdrop IDs for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) GetPendingAirdropIds() []*PendingAirdropId {
	return tx.pendingAirdropIds
}

// ----------- Overridden functions ----------------

func (tx TokenCancelAirdropTransaction) getName() string {
	return "TokenCancelAirdropTransaction"
}

func (tx TokenCancelAirdropTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx TokenCancelAirdropTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenCancelAirdrop{
			TokenCancelAirdrop: tx.buildProtoBody(),
		},
	}
}

func (tx TokenCancelAirdropTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Data: &services.SchedulableTransactionBody_TokenCancelAirdrop{
			TokenCancelAirdrop: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenCancelAirdropTransaction) buildProtoBody() *services.TokenCancelAirdropTransactionBody {
	pendingAirdrops := make([]*services.PendingAirdropId, len(tx.pendingAirdropIds))
	for i, pendingAirdropId := range tx.pendingAirdropIds {
		pendingAirdrops[i] = pendingAirdropId._ToProtobuf()
	}

	return &services.TokenCancelAirdropTransactionBody{
		PendingAirdrops: pendingAirdrops,
	}
}

func (tx TokenCancelAirdropTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().CancelAirdrop,
	}
}

func (tx TokenCancelAirdropTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenCancelAirdropTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
