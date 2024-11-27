package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// PrngTransaction is used to generate a random number in a given range
type PrngTransaction struct {
	*Transaction[*PrngTransaction]
	rang uint32
}

// NewPrngTransaction creates a PrngTransaction transaction which can be used to construct and execute
// a Prng Transaction.
func NewPrngTransaction() *PrngTransaction {
	tx := &PrngTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _PrngTransactionFromProtobuf(tx Transaction[*PrngTransaction], pb *services.TransactionBody) PrngTransaction {
	prgnTransaction := PrngTransaction{
		rang: uint32(pb.GetUtilPrng().GetRange()),
	}

	tx.childTransaction = &prgnTransaction
	prgnTransaction.Transaction = &tx
	return prgnTransaction
}

// SetPayerAccountID Sets an optional id of the account to be charged the service fee for the scheduled transaction at
// the consensus time that it executes (if ever); defaults to the ScheduleCreate payer if not
// given
func (tx *PrngTransaction) SetRange(r uint32) *PrngTransaction {
	tx._RequireNotFrozen()
	tx.rang = r

	return tx
}

// GetRange returns the range of the prng
func (tx *PrngTransaction) GetRange() uint32 {
	return tx.rang
}

// ----------- Overridden functions ----------------

func (tx PrngTransaction) getName() string {
	return "PrngTransaction"
}

func (tx PrngTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_UtilPrng{
			UtilPrng: tx.buildProtoBody(),
		},
	}
}

func (tx PrngTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_UtilPrng{
			UtilPrng: tx.buildProtoBody(),
		},
	}, nil
}

func (tx PrngTransaction) buildProtoBody() *services.UtilPrngTransactionBody {
	body := &services.UtilPrngTransactionBody{
		Range: int32(tx.rang),
	}

	return body
}

func (tx PrngTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetUtil().Prng,
	}
}
func (tx PrngTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx PrngTransaction) validateNetworkOnIDs(client *Client) error {
	return nil
}

func (tx PrngTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
