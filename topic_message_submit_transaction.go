package hedera

import (
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

const chunkSize = 4096

type TopicMessageSubmitTransaction struct {
	Transaction
	pb                  *proto.ConsensusSubmitMessageTransactionBody
	maxChunks           uint64
	message             []byte
	chunkedTransactions []*singleTopicMessageSubmitTransaction
}

func NewTopicMessageSubmitTransaction() *TopicMessageSubmitTransaction {
	transaction := TopicMessageSubmitTransaction{
		pb:                  &proto.ConsensusSubmitMessageTransactionBody{},
		Transaction:         newTransaction(),
		maxChunks:           10,
		message:             make([]byte, 0),
		chunkedTransactions: make([]*singleTopicMessageSubmitTransaction, 0),
	}

	return &transaction
}

func topicMessageSubmitTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TopicMessageSubmitTransaction {
	tx := TopicMessageSubmitTransaction{
		Transaction:         transaction,
		pb:                  pb.GetConsensusSubmitMessage(),
		maxChunks:           10,
		message:             make([]byte, 0),
		chunkedTransactions: make([]*singleTopicMessageSubmitTransaction, 0),
	}

	for _, _ = range transaction.transactions {
		singleTx := singleTopicMessageSubmitTransactionFromProtobuf(transaction, pb)
		tx.chunkedTransactions = append(tx.chunkedTransactions, &singleTx)
	}

	return tx
}

func (transaction *TopicMessageSubmitTransaction) SetTopicID(ID TopicID) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.pb.TopicID = ID.toProtobuf()
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetTopicID() TopicID {
	return topicIDFromProtobuf(transaction.pb.GetTopicID())
}

func (transaction *TopicMessageSubmitTransaction) SetMessage(message []byte) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.message = message
	return transaction
}

func (transaction *TopicMessageSubmitTransaction) GetMessage() []byte {
	return transaction.pb.GetMessage()
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
	return len(transaction.chunkedTransactions) > 0
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

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
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
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for _, tx := range transaction.chunkedTransactions {
		tx.SignWith(publicKey, signer)
	}

	return transaction
}

func (transaction *TopicMessageSubmitTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if len(transaction.Transaction.GetNodeAccountIDs()) == 0 {
		transaction.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	list, err := transaction.ExecuteAll(client)

	return list[0], err
}

// ExecuteAll executes the all the Transactions with the provided client
func (transaction *TopicMessageSubmitTransaction) ExecuteAll(
	client *Client,
) ([]TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.transactionIDs[0]

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	list := make([]TransactionResponse, len(transaction.chunkedTransactions))

	for i, tx := range transaction.chunkedTransactions {
		resp, err := tx.Execute(client)
		if err != nil {
			return []TransactionResponse{}, err
		}

		list[i] = resp
	}

	return list, nil
}

func (transaction *TopicMessageSubmitTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ConsensusSubmitMessage{
		ConsensusSubmitMessage: transaction.pb,
	}

	return true
}

func (transaction *TopicMessageSubmitTransaction) Freeze() (*TopicMessageSubmitTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicMessageSubmitTransaction) FreezeWith(client *Client) (*TopicMessageSubmitTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	if transaction.pb.ChunkInfo != nil || len(transaction.message) < chunkSize {
		tx, err := newSingleTopicMessageSubmitTransaction(
			protobuf.Clone(transaction.Transaction.pbBody).(*proto.TransactionBody),
			protobuf.Clone(transaction.pb).(*proto.ConsensusSubmitMessageTransactionBody),
			transaction.message,
			transaction.pb.ChunkInfo,
		).FreezeWith(client)
		if err != nil {
			return transaction, err
		}

		transaction.chunkedTransactions = append(transaction.chunkedTransactions, tx)

		return transaction, err
	}

	chunks := uint64((len(transaction.message) + (chunkSize - 1)) / chunkSize)
	if chunks > transaction.maxChunks {
		return transaction, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: transaction.maxChunks,
		}
	}

	var initialTransactionID TransactionID
	if transaction.Transaction.pbBody.TransactionID != nil {
		initialTransactionID = transaction.GetTransactionID()
	} else {
		initialTransactionID = TransactionIDGenerate(client.GetOperatorAccountID())
	}

	nextTransactionID := initialTransactionID

	for i := 0; uint64(i) < chunks; i += 1 {
		start := i * chunkSize
		end := start + chunkSize

		if end > len(transaction.message) {
			end = len(transaction.message)
		}

		transaction.SetTransactionID(nextTransactionID)

		tx, err := newSingleTopicMessageSubmitTransaction(
			protobuf.Clone(transaction.Transaction.pbBody).(*proto.TransactionBody),
			protobuf.Clone(transaction.pb).(*proto.ConsensusSubmitMessageTransactionBody),
			transaction.message[start:end],
			&proto.ConsensusMessageChunkInfo{
				InitialTransactionID: initialTransactionID.toProtobuf(),
				Total:                int32(chunks),
				Number:               int32(i) + 1,
			},
		).FreezeWith(client)
		if err != nil {
			return transaction, err
		}

		transaction.chunkedTransactions = append(transaction.chunkedTransactions, tx)
		nextTransactionID.ValidStart = nextTransactionID.ValidStart.Add(1 * time.Nanosecond)
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
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

func (transaction *TopicMessageSubmitTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeAccountID sets the node AccountID for this TopicMessageSubmitTransaction.
func (transaction *TopicMessageSubmitTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicMessageSubmitTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}
