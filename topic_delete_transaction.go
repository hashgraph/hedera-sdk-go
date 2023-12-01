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
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// TopicDeleteTransaction is for deleting a topic on HCS.
type TopicDeleteTransaction struct {
	transaction
	topicID *TopicID
}

// NewTopicDeleteTransaction creates a TopicDeleteTransaction which can be used to construct
// and execute a Consensus Delete Topic transaction.
func NewTopicDeleteTransaction() *TopicDeleteTransaction {
	tx := TopicDeleteTransaction{
		transaction: _NewTransaction(),
	}

	tx.e = &tx
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _TopicDeleteTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *TopicDeleteTransaction {
	resultTx := &TopicDeleteTransaction{
		transaction: tx,
		topicID:     _TopicIDFromProtobuf(pb.GetConsensusDeleteTopic().GetTopicID()),
	}
	resultTx.e = resultTx
	return resultTx
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

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TopicDeleteTransaction) Sign(privateKey PrivateKey) *TopicDeleteTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TopicDeleteTransaction) SignWithOperator(client *Client) (*TopicDeleteTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TopicDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicDeleteTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TopicDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicDeleteTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TopicDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *TopicDeleteTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TopicDeleteTransaction) Freeze() (*TopicDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TopicDeleteTransaction) FreezeWith(client *Client) (*TopicDeleteTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TopicDeleteTransaction.
func (tx *TopicDeleteTransaction) SetMaxTransactionFee(fee Hbar) *TopicDeleteTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TopicDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TopicDeleteTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TopicDeleteTransaction.
func (tx *TopicDeleteTransaction) SetTransactionMemo(memo string) *TopicDeleteTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TopicDeleteTransaction.
func (tx *TopicDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *TopicDeleteTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TopicDeleteTransaction.
func (tx *TopicDeleteTransaction) SetTransactionID(transactionID TransactionID) *TopicDeleteTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TopicDeleteTransaction.
func (tx *TopicDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicDeleteTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TopicDeleteTransaction) SetMaxRetry(count int) *TopicDeleteTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TopicDeleteTransaction) SetMaxBackoff(max time.Duration) *TopicDeleteTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TopicDeleteTransaction) SetMinBackoff(min time.Duration) *TopicDeleteTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *TopicDeleteTransaction) SetLogLevel(level LogLevel) *TopicDeleteTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

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
		Memo:                     tx.transaction.memo,
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
		Memo:           tx.transaction.memo,
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
func (tx *TopicDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
