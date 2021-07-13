package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
	"time"

	protobuf "github.com/golang/protobuf/proto"
)

const chunkSize = 1024

type TopicMessageSubmitTransaction struct {
	Transaction
	pb        *services.ConsensusSubmitMessageTransactionBody
	maxChunks uint64
	message   []byte
	topicID   TopicID
}

func NewTopicMessageSubmitTransaction() *TopicMessageSubmitTransaction {
	transaction := TopicMessageSubmitTransaction{
		pb:          &services.ConsensusSubmitMessageTransactionBody{},
		Transaction: newTransaction(),
		maxChunks:   20,
		message:     make([]byte, 0),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func topicMessageSubmitTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) TopicMessageSubmitTransaction {
	tx := TopicMessageSubmitTransaction{
		Transaction: transaction,
		pb:          pb.GetConsensusSubmitMessage(),
		maxChunks:   20,
		message:     make([]byte, 0),
		topicID:     topicIDFromProtobuf(pb.GetConsensusSubmitMessage().GetTopicID(), nil),
	}

	return tx
}

func (transaction *TopicMessageSubmitTransaction) SetTopicID(id TopicID) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.topicID = id
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetTopicID() TopicID {
	return transaction.topicID
}

func (transaction *TopicMessageSubmitTransaction) SetMessage(message []byte) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.message = message
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetMessage() []byte {
	return transaction.message
}

func (transaction *TopicMessageSubmitTransaction) SetMaxChunks(maxChunks uint64) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.maxChunks = maxChunks
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetMaxChunks() uint64 {
	return transaction.maxChunks
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (transaction *TopicMessageSubmitTransaction) IsFrozen() bool {
	return transaction.Transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TopicMessageSubmitTransaction) Sign(
	privateKey PrivateKey,
) *TopicMessageSubmitTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TopicMessageSubmitTransaction) SignWithOperator(
	client *Client,
) (*TopicMessageSubmitTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

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
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

func (transaction *TopicMessageSubmitTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.topicID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *TopicMessageSubmitTransaction) build() *TopicMessageSubmitTransaction {
	if !transaction.topicID.isZero() {
		transaction.pb.TopicID = transaction.topicID.toProtobuf()
	}

	return transaction
}

func (transaction *TopicMessageSubmitTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	chunks := uint64((len(transaction.message) + (chunkSize - 1)) / chunkSize)

	if chunks > 1 {
		return &ScheduleCreateTransaction{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: 1,
		}
	}

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TopicMessageSubmitTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	transaction.build()
	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &services.SchedulableTransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: &services.ConsensusSubmitMessageTransactionBody{
				TopicID:   transaction.pb.GetTopicID(),
				Message:   transaction.message,
				ChunkInfo: &services.ConsensusMessageChunkInfo{},
			},
		},
	}, nil
}

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
	} else {
		return TransactionResponse{}, errNoTransactions
	}
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

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(accountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	size := len(transaction.signedTransactions) / len(transaction.nodeIDs)
	list := make([]TransactionResponse, size)

	for i := 0; i < size; i++ {
		resp, err := execute(
			client,
			request{
				transaction: &transaction.Transaction,
			},
			transaction_shouldRetry,
			transaction_makeRequest,
			transaction_advanceRequest,
			transaction_getNodeAccountID,
			topicMessageSubmitTransaction_getMethod,
			transaction_mapStatusError,
			transaction_mapResponse,
		)

		if err != nil {
			return []TransactionResponse{}, err
		}

		list[i] = resp.transaction
	}

	return list, nil
}

func (transaction *TopicMessageSubmitTransaction) Freeze() (*TopicMessageSubmitTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicMessageSubmitTransaction) FreezeWith(client *Client) (*TopicMessageSubmitTransaction, error) {
	if len(transaction.nodeIDs) == 0 {
		if client == nil {
			return transaction, errNoClientOrTransactionIDOrNodeId
		} else {
			transaction.nodeIDs = client.network.getNodeAccountIDsForExecute()
		}
	}

	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TopicMessageSubmitTransaction{}, err
	}
	transaction.build()

	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	chunks := uint64((len(transaction.message) + (chunkSize - 1)) / chunkSize)
	if chunks > transaction.maxChunks {
		return transaction, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: transaction.maxChunks,
		}
	}

	initialTransactionID := transaction.GetTransactionID()
	nextTransactionID := transactionIDFromProtobuf(initialTransactionID.toProtobuf(), nil)

	transaction.transactionIDs = make([]TransactionID, 0)
	transaction.transactions = make([]*services.Transaction, 0)
	transaction.signedTransactions = make([]*services.SignedTransaction, 0)

	for i := 0; uint64(i) < chunks; i += 1 {
		start := i * chunkSize
		end := start + chunkSize

		if end > len(transaction.message) {
			end = len(transaction.message)
		}

		transaction.transactionIDs = append(transaction.transactionIDs, transactionIDFromProtobuf(nextTransactionID.toProtobuf(), nil))

		transaction.pb.Message = transaction.message[start:end]
		transaction.pb.ChunkInfo = &services.ConsensusMessageChunkInfo{
			InitialTransactionID: initialTransactionID.toProtobuf(),
			Total:                int32(chunks),
			Number:               int32(i) + 1,
		}

		transaction.pbBody.TransactionID = nextTransactionID.toProtobuf()
		transaction.pbBody.Data = &services.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: transaction.pb,
		}

		for _, nodeAccountID := range transaction.nodeIDs {
			transaction.pbBody.NodeAccountID = nodeAccountID.toProtobuf()

			bodyBytes, err := protobuf.Marshal(transaction.pbBody)
			if err != nil {
				return transaction, errors.Wrap(err, "error serializing transaction body for topic submission")
			}

			transaction.signedTransactions = append(transaction.signedTransactions, &services.SignedTransaction{
				BodyBytes: bodyBytes,
				SigMap:    &services.SignatureMap{},
			})
		}

		validStart := *nextTransactionID.ValidStart

		*nextTransactionID.ValidStart = validStart.Add(1 * time.Nanosecond)
	}

	return transaction, nil
}

func topicMessageSubmitTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getTopic().SubmitMessage,
	}
}

func (transaction *TopicMessageSubmitTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetMaxTransactionFee(fee Hbar) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetTransactionMemo(memo string) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetTransactionValidDuration(duration time.Duration) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetTransactionID(transactionID TransactionID) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) SetMaxRetry(count int) *TopicMessageSubmitTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicMessageSubmitTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
