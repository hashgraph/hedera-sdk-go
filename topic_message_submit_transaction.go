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
	"fmt"
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
	transaction := TopicMessageSubmitTransaction{
		Transaction: _NewTransaction(),
		maxChunks:   20,
		message:     make([]byte, 0),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _TopicMessageSubmitTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TopicMessageSubmitTransaction {
	tx := &TopicMessageSubmitTransaction{
		Transaction: transaction,
		maxChunks:   20,
		message:     pb.GetConsensusSubmitMessage().GetMessage(),
		topicID:     _TopicIDFromProtobuf(pb.GetConsensusSubmitMessage().GetTopicID()),
	}

	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *TopicMessageSubmitTransaction) SetGrpcDeadline(deadline *time.Duration) *TopicMessageSubmitTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetTopicID Sets the topic to submit message to.
func (transaction *TopicMessageSubmitTransaction) SetTopicID(topicID TopicID) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()
	transaction.topicID = &topicID
	return transaction
}

// GetTopicID returns the TopicID for this TopicMessageSubmitTransaction
func (transaction *TopicMessageSubmitTransaction) GetTopicID() TopicID {
	if transaction.topicID == nil {
		return TopicID{}
	}

	return *transaction.topicID
}

// SetMessage Sets the message to be submitted.
func (transaction *TopicMessageSubmitTransaction) SetMessage(message []byte) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()
	transaction.message = message
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetMessage() []byte {
	return transaction.message
}

// SetMaxChunks sets the maximum amount of chunks to use to send the message
func (transaction *TopicMessageSubmitTransaction) SetMaxChunks(maxChunks uint64) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()
	transaction.maxChunks = maxChunks
	return transaction
}

// GetMaxChunks returns the maximum amount of chunks to use to send the message
func (transaction *TopicMessageSubmitTransaction) GetMaxChunks() uint64 {
	return transaction.maxChunks
}

func (transaction *TopicMessageSubmitTransaction) IsFrozen() bool {
	return transaction.Transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TopicMessageSubmitTransaction) Sign(
	privateKey PrivateKey,
) *TopicMessageSubmitTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *TopicMessageSubmitTransaction) SignWithOperator(
	client *Client,
) (*TopicMessageSubmitTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *TopicMessageSubmitTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicMessageSubmitTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

func (transaction *TopicMessageSubmitTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.topicID != nil {
		if err := transaction.topicID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TopicMessageSubmitTransaction) _Build() *services.TransactionBody {
	body := &services.ConsensusSubmitMessageTransactionBody{}
	if transaction.topicID != nil {
		body.TopicID = transaction.topicID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: body,
		},
	}
}

func (transaction *TopicMessageSubmitTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	chunks := uint64((len(transaction.message) + (chunkSize - 1)) / chunkSize)

	if chunks > 1 {
		return &ScheduleCreateTransaction{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: 1,
		}
	}

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TopicMessageSubmitTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ConsensusSubmitMessageTransactionBody{
		Message: transaction.message,
	}

	if transaction.topicID != nil {
		body.TopicID = transaction.topicID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: body,
		},
	}, nil
}

// Execute executes the Query with the provided client
func (transaction *TopicMessageSubmitTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	list, err := transaction.ExecuteAll(client)

	if err != nil {
		return TransactionResponse{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return TransactionResponse{}, errNoTransactions
}

// ExecuteAll executes the all the Transactions with the provided client
func (transaction *TopicMessageSubmitTransaction) ExecuteAll(
	client *Client,
) ([]TransactionResponse, error) {
	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return []TransactionResponse{}, err
		}
	}
	transactionID := transaction.GetTransactionID()
	accountID := AccountID{}
	if transactionID.AccountID != nil {
		accountID = *transactionID.AccountID
	}

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(accountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	size := transaction.signedTransactions._Length() / transaction.nodeAccountIDs._Length()
	list := make([]TransactionResponse, size)

	for i := 0; i < size; i++ {
		resp, err := _Execute(
			client,
			&transaction.Transaction,
			_TransactionShouldRetry,
			_TransactionMakeRequest,
			_TransactionAdvanceRequest,
			_TransactionGetNodeAccountID,
			_TopicMessageSubmitTransactionGetMethod,
			_TransactionMapStatusError,
			_TransactionMapResponse,
			transaction._GetLogID(),
			transaction.grpcDeadline,
			transaction.maxBackoff,
			transaction.minBackoff,
			transaction.maxRetry,
		)

		if err != nil {
			return []TransactionResponse{}, err
		}

		list[i] = resp.(TransactionResponse)
	}

	return list, nil
}

func (transaction *TopicMessageSubmitTransaction) Freeze() (*TopicMessageSubmitTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicMessageSubmitTransaction) FreezeWith(client *Client) (*TopicMessageSubmitTransaction, error) {
	var err error
	if transaction.nodeAccountIDs._Length() == 0 {
		if client == nil {
			return transaction, errNoClientOrTransactionIDOrNodeId
		}

		transaction.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	transaction._InitFee(client)
	err = transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TopicMessageSubmitTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	chunks := uint64((len(transaction.message) + (chunkSize - 1)) / chunkSize)
	if chunks > transaction.maxChunks {
		return transaction, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: transaction.maxChunks,
		}
	}

	initialTransactionID := transaction.transactionIDs._GetCurrent().(TransactionID)
	nextTransactionID, _ := TransactionIdFromString(initialTransactionID.String())

	transaction.transactionIDs = _NewLockableSlice()
	transaction.transactions = _NewLockableSlice()
	transaction.signedTransactions = _NewLockableSlice()
	if b, ok := body.Data.(*services.TransactionBody_ConsensusSubmitMessage); ok {
		for i := 0; uint64(i) < chunks; i++ {
			start := i * chunkSize
			end := start + chunkSize

			if end > len(transaction.message) {
				end = len(transaction.message)
			}

			transaction.transactionIDs._Push(_TransactionIDFromProtobuf(nextTransactionID._ToProtobuf()))

			b.ConsensusSubmitMessage.Message = transaction.message[start:end]
			b.ConsensusSubmitMessage.ChunkInfo = &services.ConsensusMessageChunkInfo{
				InitialTransactionID: initialTransactionID._ToProtobuf(),
				Total:                int32(chunks),
				Number:               int32(i) + 1,
			}

			body.TransactionID = nextTransactionID._ToProtobuf()
			body.Data = &services.TransactionBody_ConsensusSubmitMessage{
				ConsensusSubmitMessage: b.ConsensusSubmitMessage,
			}

			for _, nodeAccountID := range transaction.nodeAccountIDs.slice {
				body.NodeAccountID = nodeAccountID.(AccountID)._ToProtobuf()

				bodyBytes, err := protobuf.Marshal(body)
				if err != nil {
					return transaction, errors.Wrap(err, "error serializing transaction body for topic submission")
				}

				transaction.signedTransactions._Push(&services.SignedTransaction{
					BodyBytes: bodyBytes,
					SigMap:    &services.SignatureMap{},
				})
			}

			validStart := *nextTransactionID.ValidStart

			*nextTransactionID.ValidStart = validStart.Add(1 * time.Nanosecond)
		}
	}

	return transaction, nil
}

func _TopicMessageSubmitTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().SubmitMessage,
	}
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TopicMessageSubmitTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TopicMessageSubmitTransaction) SetMaxTransactionFee(fee Hbar) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TopicMessageSubmitTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TopicMessageSubmitTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this	TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetTransactionMemo(memo string) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *TopicMessageSubmitTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetTransactionValidDuration(duration time.Duration) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetTransactionID(transactionID TransactionID) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicMessageSubmitTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *TopicMessageSubmitTransaction) SetMaxRetry(count int) *TopicMessageSubmitTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *TopicMessageSubmitTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicMessageSubmitTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := transaction.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
	}

	return transaction
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *TopicMessageSubmitTransaction) SetMaxBackoff(max time.Duration) *TopicMessageSubmitTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *TopicMessageSubmitTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *TopicMessageSubmitTransaction) SetMinBackoff(min time.Duration) *TopicMessageSubmitTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *TopicMessageSubmitTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TopicMessageSubmitTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TopicMessageSubmitTransaction:%d", timestamp.UnixNano())
}

func (transaction *TopicMessageSubmitTransaction) SetLogLevel(level LogLevel) *TopicMessageSubmitTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
