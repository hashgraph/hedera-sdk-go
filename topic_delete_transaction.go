package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TopicDeleteTransaction is for deleting a topic on HCS.
type TopicDeleteTransaction struct {
	*Transaction[*TopicDeleteTransaction]
	topicID *TopicID
}

// NewTopicDeleteTransaction creates a TopicDeleteTransaction which can be used to construct
// and execute a Consensus Delete Topic Transaction.
func NewTopicDeleteTransaction() *TopicDeleteTransaction {
	tx := &TopicDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _TopicDeleteTransactionFromProtobuf(tx Transaction[*TopicDeleteTransaction], pb *services.TransactionBody) TopicDeleteTransaction {
	topicDeleteTransaction := TopicDeleteTransaction{
		topicID: _TopicIDFromProtobuf(pb.GetConsensusDeleteTopic().GetTopicID()),
	}

	tx.childTransaction = &topicDeleteTransaction
	topicDeleteTransaction.Transaction = &tx
	return topicDeleteTransaction
}

// SetTopicID sets the topic IDentifier.
func (tx *TopicDeleteTransaction) SetTopicID(topicID TopicID) *TopicDeleteTransaction {
	tx._RequireNotFrozen()
	tx.topicID = &topicID
	return tx
}

// GetTopicID returns the topic IDentifier.
func (tx *TopicDeleteTransaction) GetTopicID() TopicID {
	if tx.topicID == nil {
		return TopicID{}
	}

	return *tx.topicID
}

// ----------- Overridden functions ----------------

func (tx TopicDeleteTransaction) getName() string {
	return "TopicDeleteTransaction"
}

func (tx TopicDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.topicID != nil {
		if err := tx.topicID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TopicDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: tx.buildProtoBody(),
		},
	}
}

func (tx TopicDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TopicDeleteTransaction) buildProtoBody() *services.ConsensusDeleteTopicTransactionBody {
	body := &services.ConsensusDeleteTopicTransactionBody{}
	if tx.topicID != nil {
		body.TopicID = tx.topicID._ToProtobuf()
	}

	return body
}

func (tx TopicDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().DeleteTopic,
	}
}

func (tx TopicDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TopicDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
