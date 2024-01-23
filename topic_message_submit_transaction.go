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
	"time"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

const chunkSize = 1024

// TopicMessageSubmitTransaction
// Sends a message/messages to the Topic ID
type TopicMessageSubmitTransaction struct {
	Transaction
	maxChunks uint64
	message   []byte
	topicID   *TopicID
}

// NewTopicMessageSubmitTransaction createsTopicMessageSubmitTransaction which
// sends a message/messages to the Topic ID
func NewTopicMessageSubmitTransaction() *TopicMessageSubmitTransaction {
	tx := TopicMessageSubmitTransaction{
		Transaction: _NewTransaction(),
		maxChunks:   20,
		message:     make([]byte, 0),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _TopicMessageSubmitTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TopicMessageSubmitTransaction {
	return &TopicMessageSubmitTransaction{
		Transaction: tx,
		maxChunks:   20,
		message:     pb.GetConsensusSubmitMessage().GetMessage(),
		topicID:     _TopicIDFromProtobuf(pb.GetConsensusSubmitMessage().GetTopicID()),
	}
}

// SetTopicID Sets the topic to submit message to.
func (tx *TopicMessageSubmitTransaction) SetTopicID(topicID TopicID) *TopicMessageSubmitTransaction {
	tx._RequireNotFrozen()
	tx.topicID = &topicID
	return tx
}

// GetTopicID returns the TopicID for this TopicMessageSubmitTransaction
func (tx *TopicMessageSubmitTransaction) GetTopicID() TopicID {
	if tx.topicID == nil {
		return TopicID{}
	}

	return *tx.topicID
}

// SetMessage Sets the message to be submitted.
func (tx *TopicMessageSubmitTransaction) SetMessage(message []byte) *TopicMessageSubmitTransaction {
	tx._RequireNotFrozen()
	tx.message = message
	return tx
}

func (tx *TopicMessageSubmitTransaction) GetMessage() []byte {
	return tx.message
}

// SetMaxChunks sets the maximum amount of chunks to use to send the message
func (tx *TopicMessageSubmitTransaction) SetMaxChunks(maxChunks uint64) *TopicMessageSubmitTransaction {
	tx._RequireNotFrozen()
	tx.maxChunks = maxChunks
	return tx
}

// GetMaxChunks returns the maximum amount of chunks to use to send the message
func (tx *TopicMessageSubmitTransaction) GetMaxChunks() uint64 {
	return tx.maxChunks
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TopicMessageSubmitTransaction) Sign(privateKey PrivateKey) *TopicMessageSubmitTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TopicMessageSubmitTransaction) SignWithOperator(client *Client) (*TopicMessageSubmitTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TopicMessageSubmitTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicMessageSubmitTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TopicMessageSubmitTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicMessageSubmitTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TopicMessageSubmitTransaction) SetGrpcDeadline(deadline *time.Duration) *TopicMessageSubmitTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TopicMessageSubmitTransaction) Freeze() (*TopicMessageSubmitTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TopicMessageSubmitTransaction) FreezeWith(client *Client) (*TopicMessageSubmitTransaction, error) {
	var err error
	if tx.nodeAccountIDs._Length() == 0 {
		if client == nil {
			return tx, errNoClientOrTransactionIDOrNodeId
		}

		tx.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	tx._InitFee(client)
	err = tx.validateNetworkOnIDs(client)
	if err != nil {
		return &TopicMessageSubmitTransaction{}, err
	}
	if err := tx._InitTransactionID(client); err != nil {
		return tx, err
	}
	body := tx.build()

	chunks := uint64((len(tx.message) + (chunkSize - 1)) / chunkSize)
	if chunks > tx.maxChunks {
		return tx, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: tx.maxChunks,
		}
	}

	initialTransactionID := tx.transactionIDs._GetCurrent().(TransactionID)
	nextTransactionID, _ := TransactionIdFromString(initialTransactionID.String())

	tx.transactionIDs = _NewLockableSlice()
	tx.transactions = _NewLockableSlice()
	tx.signedTransactions = _NewLockableSlice()
	if b, ok := body.Data.(*services.TransactionBody_ConsensusSubmitMessage); ok {
		for i := 0; uint64(i) < chunks; i++ {
			start := i * chunkSize
			end := start + chunkSize

			if end > len(tx.message) {
				end = len(tx.message)
			}

			tx.transactionIDs._Push(_TransactionIDFromProtobuf(nextTransactionID._ToProtobuf()))

			b.ConsensusSubmitMessage.Message = tx.message[start:end]
			b.ConsensusSubmitMessage.ChunkInfo = &services.ConsensusMessageChunkInfo{
				InitialTransactionID: initialTransactionID._ToProtobuf(),
				Total:                int32(chunks),
				Number:               int32(i) + 1,
			}

			body.TransactionID = nextTransactionID._ToProtobuf()
			body.Data = &services.TransactionBody_ConsensusSubmitMessage{
				ConsensusSubmitMessage: b.ConsensusSubmitMessage,
			}

			for _, nodeAccountID := range tx.nodeAccountIDs.slice {
				body.NodeAccountID = nodeAccountID.(AccountID)._ToProtobuf()

				bodyBytes, err := protobuf.Marshal(body)
				if err != nil {
					return tx, errors.Wrap(err, "error serializing tx body for topic submission")
				}

				tx.signedTransactions._Push(&services.SignedTransaction{
					BodyBytes: bodyBytes,
					SigMap:    &services.SignatureMap{},
				})
			}

			validStart := *nextTransactionID.ValidStart

			*nextTransactionID.ValidStart = validStart.Add(1 * time.Nanosecond)
		}
	}

	return tx, nil
}

// SetMaxTransactionFee sets the max transaction fee for this TopicMessageSubmitTransaction.
func (tx *TopicMessageSubmitTransaction) SetMaxTransactionFee(fee Hbar) *TopicMessageSubmitTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TopicMessageSubmitTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TopicMessageSubmitTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TopicMessageSubmitTransaction.
func (tx *TopicMessageSubmitTransaction) SetTransactionMemo(memo string) *TopicMessageSubmitTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TopicMessageSubmitTransaction.
func (tx *TopicMessageSubmitTransaction) SetTransactionValidDuration(duration time.Duration) *TopicMessageSubmitTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TopicMessageSubmitTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TopicMessageSubmitTransaction.
func (tx *TopicMessageSubmitTransaction) SetTransactionID(transactionID TransactionID) *TopicMessageSubmitTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TopicMessageSubmitTransaction.
func (tx *TopicMessageSubmitTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicMessageSubmitTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TopicMessageSubmitTransaction) SetMaxRetry(count int) *TopicMessageSubmitTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TopicMessageSubmitTransaction) SetMaxBackoff(max time.Duration) *TopicMessageSubmitTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TopicMessageSubmitTransaction) SetMinBackoff(min time.Duration) *TopicMessageSubmitTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TopicMessageSubmitTransaction) SetLogLevel(level LogLevel) *TopicMessageSubmitTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TopicMessageSubmitTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	chunks := uint64((len(tx.message) + (chunkSize - 1)) / chunkSize)
	if chunks > 1 {
		return &ScheduleCreateTransaction{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: 1,
		}
	}

	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TopicMessageSubmitTransaction) getName() string {
	return "TopicMessageSubmitTransaction"
}
func (tx *TopicMessageSubmitTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TopicMessageSubmitTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: tx.buildProtoBody(),
		},
	}
}

func (tx *TopicMessageSubmitTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TopicMessageSubmitTransaction) buildProtoBody() *services.ConsensusSubmitMessageTransactionBody {
	body := &services.ConsensusSubmitMessageTransactionBody{
		Message: tx.message,
	}

	if tx.topicID != nil {
		body.TopicID = tx.topicID._ToProtobuf()
	}

	return body
}

func (tx *TopicMessageSubmitTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().SubmitMessage,
	}
}

// Execute executes the Query with the provided client
func (tx *TopicMessageSubmitTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if tx.freezeError != nil {
		return TransactionResponse{}, tx.freezeError
	}

	list, err := tx.ExecuteAll(client)

	if err != nil {
		return TransactionResponse{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return TransactionResponse{}, errNoTransactions
}

// ExecuteAll executes the all the Transactions with the provided client
func (tx *TopicMessageSubmitTransaction) ExecuteAll(
	client *Client,
) ([]TransactionResponse, error) {
	if !tx.IsFrozen() {
		_, err := tx.FreezeWith(client)
		if err != nil {
			return []TransactionResponse{}, err
		}
	}
	transactionID := tx.GetTransactionID()
	accountID := AccountID{}
	if transactionID.AccountID != nil {
		accountID = *transactionID.AccountID
	}

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(accountID) {
		tx.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	size := tx.signedTransactions._Length() / tx.nodeAccountIDs._Length()
	list := make([]TransactionResponse, size)

	for i := 0; i < size; i++ {
		resp, err := _Execute(client, tx)

		if err != nil {
			return []TransactionResponse{}, err
		}

		list[i] = resp.(TransactionResponse)
	}

	return list, nil
}

func (tx *TopicMessageSubmitTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
