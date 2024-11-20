package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/pkg/errors"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

const chunkSize = 1024

// TopicMessageSubmitTransaction
// Sends a message/messages to the Topic ID
type TopicMessageSubmitTransaction struct {
	*Transaction[*TopicMessageSubmitTransaction]
	maxChunks uint64
	message   []byte
	topicID   *TopicID
}

// NewTopicMessageSubmitTransaction createsTopicMessageSubmitTransaction which
// sends a message/messages to the Topic ID
func NewTopicMessageSubmitTransaction() *TopicMessageSubmitTransaction {
	tx := &TopicMessageSubmitTransaction{
		maxChunks: 20,
		message:   make([]byte, 0),
	}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _TopicMessageSubmitTransactionFromProtobuf(tx Transaction[*TopicMessageSubmitTransaction], pb *services.TransactionBody) TopicMessageSubmitTransaction {
	topicMessageSubmitTransaction := TopicMessageSubmitTransaction{
		maxChunks: 20,
		message:   pb.GetConsensusSubmitMessage().GetMessage(),
		topicID:   _TopicIDFromProtobuf(pb.GetConsensusSubmitMessage().GetTopicID()),
	}

	tx.childTransaction = &topicMessageSubmitTransaction
	topicMessageSubmitTransaction.Transaction = &tx
	return topicMessageSubmitTransaction
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
// The message should be a byte array or a string
// If other types are provided, it will not set the value
func (tx *TopicMessageSubmitTransaction) SetMessage(message interface{}) *TopicMessageSubmitTransaction {
	tx._RequireNotFrozen()
	switch m := message.(type) {
	case string:
		tx.message = []byte(m)
	case []byte:
		tx.message = m
	}
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

func (tx *TopicMessageSubmitTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	chunks := uint64((len(tx.message) + (chunkSize - 1)) / chunkSize)
	if chunks > 1 {
		return &ScheduleCreateTransaction{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: 1,
		}
	}

	return tx.Transaction.Schedule()
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

// ----------- Overridden functions ----------------

func (tx TopicMessageSubmitTransaction) getName() string {
	return "TopicMessageSubmitTransaction"
}
func (tx TopicMessageSubmitTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx TopicMessageSubmitTransaction) build() *services.TransactionBody {
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

func (tx TopicMessageSubmitTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TopicMessageSubmitTransaction) buildProtoBody() *services.ConsensusSubmitMessageTransactionBody {
	body := &services.ConsensusSubmitMessageTransactionBody{
		Message: tx.message,
	}

	if tx.topicID != nil {
		body.TopicID = tx.topicID._ToProtobuf()
	}

	return body
}

func (tx TopicMessageSubmitTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().SubmitMessage,
	}
}

func (tx TopicMessageSubmitTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TopicMessageSubmitTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
