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

func _TopicDeleteTransactionFromProtobuf(pb *services.TransactionBody) *TopicDeleteTransaction {
	return &TopicDeleteTransaction{
		topicID: _TopicIDFromProtobuf(pb.GetConsensusDeleteTopic().GetTopicID()),
	}
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

func (tx *TopicDeleteTransaction) getName() string {
	return "TopicDeleteTransaction"
}

func (tx *TopicDeleteTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TopicDeleteTransaction) build() *services.TransactionBody {
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

func (tx *TopicDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TopicDeleteTransaction) buildProtoBody() *services.ConsensusDeleteTopicTransactionBody {
	body := &services.ConsensusDeleteTopicTransactionBody{}
	if tx.topicID != nil {
		body.TopicID = tx.topicID._ToProtobuf()
	}

	return body
}

func (tx *TopicDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().DeleteTopic,
	}
}

func (tx *TopicDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *TopicDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction)
}

func (tx *TopicDeleteTransaction) setBaseTransaction(baseTx Transaction[TransactionInterface]) {
	tx.Transaction = castFromBaseToConcreteTransaction[*TopicDeleteTransaction](baseTx)
}
