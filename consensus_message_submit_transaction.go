package hedera

import (
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

const chunkSize = 4096

// A ConsensusMessageSubmitTransaction is used for submitting a message to HCS.
type ConsensusMessageSubmitTransaction struct {
	TransactionBuilder
	pb                   *proto.ConsensusSubmitMessageTransactionBody
	maxChunks            uint64
	message              []byte
	topicID              ConsensusTopicID
	initialTransactionID TransactionID
	total                int32
	number               int32
	chunkInfoSet         bool
}

// NewConsensusMessageSubmitTransaction creates a ConsensusMessageSubmitTransaction transaction which can be used to
// construct and execute a Consensus Submit Message Transaction.
func NewConsensusMessageSubmitTransaction() ConsensusMessageSubmitTransaction {
	pb := &proto.ConsensusSubmitMessageTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusSubmitMessage{ConsensusSubmitMessage: pb}

	transaction := ConsensusMessageSubmitTransaction{inner, pb, 10, nil, ConsensusTopicID{}, TransactionID{}, 0, 0, false}

	return transaction
}

// SetTopicID sets the topic to submit the message to.
func (transaction ConsensusMessageSubmitTransaction) SetTopicID(id ConsensusTopicID) ConsensusMessageSubmitTransaction {
	transaction.topicID = id
	return transaction
}

// SetMessage sets the message to be submitted. Max size of the Transaction (including signatures) is 4kB.
func (transaction ConsensusMessageSubmitTransaction) SetMessage(message []byte) ConsensusMessageSubmitTransaction {
	transaction.message = message
	return transaction
}

// SetMessage sets the message to be submitted. Max size of the Transaction (including signatures) is 4kB.
func (transaction ConsensusMessageSubmitTransaction) SetMaxChunks(max uint64) ConsensusMessageSubmitTransaction {
	transaction.maxChunks = max
	return transaction
}

// SetMessage sets the message to be submitted. Max size of the Transaction (including signatures) is 4kB.
func (transaction ConsensusMessageSubmitTransaction) SetChunkInfo(transactionID TransactionID, total int32, number int32) ConsensusMessageSubmitTransaction {
	transaction.initialTransactionID = transactionID
	transaction.total = total
	transaction.number = number
	transaction.chunkInfoSet = true
	return transaction
}

func (transaction ConsensusMessageSubmitTransaction) Execute(client *Client) (TransactionID, error) {
	txs, err := transaction.Build(client)
	if err != nil {
		return TransactionID{}, err
	}

	return txs.Execute(client)
}

func (transaction ConsensusMessageSubmitTransaction) ExecuteAll(client *Client) ([]TransactionID, error) {
	txs, err := transaction.Build(client)
	if err != nil {
		return nil, err
	}

	return txs.ExecuteAll(client)
}

func (transaction ConsensusMessageSubmitTransaction) Build(client *Client) (TransactionList, error) {
	// If chunk info  is set then we aren't going to chunk the message
	// Set all the required fields and return a list of 1
	if transaction.chunkInfoSet {
		transaction.pb.TopicID = transaction.topicID.toProto()
		transaction.pb.Message = transaction.message
		transaction.pb.ChunkInfo = &proto.ConsensusMessageChunkInfo{
			InitialTransactionID: transaction.initialTransactionID.toProto(),
			Number:               transaction.number,
			Total:                transaction.total,
		}

		// FIXME: really have no idea why this is needed @daniel
		transaction.TransactionBuilder.pb.Data = &proto.TransactionBody_ConsensusSubmitMessage{ConsensusSubmitMessage: transaction.pb}

		transaction, err := transaction.TransactionBuilder.Build(client)
		if err != nil {
			return TransactionList{}, err
		}

		list := TransactionList{
			List: make([]Transaction, 1),
		}

		list.List[0] = transaction
		return list, nil
	}

	chunks := uint64((len(transaction.message) + (chunkSize - 1)) / chunkSize)

	if chunks > transaction.maxChunks {
		return TransactionList{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: transaction.maxChunks,
		}
	}

	list := make([]Transaction, chunks)

	var initialTransactionID TransactionID
	if transaction.TransactionBuilder.pb.TransactionID != nil {
		initialTransactionID = transactionIDFromProto(transaction.TransactionBuilder.pb.TransactionID)
	} else {
		initialTransactionID = NewTransactionID(client.GetOperatorID())
	}

	nextTransactionID := initialTransactionID

	for i := 0; uint64(i) < chunks; i += 1 {
		start := i * chunkSize
		end := start + chunkSize

		if end > len(transaction.message) {
			end = len(transaction.message)
		}

		transactionBuilder := NewConsensusMessageSubmitTransaction()
		transactionBuilder.TransactionBuilder.pb = protobuf.Clone(transaction.TransactionBuilder.pb).(*proto.TransactionBody)

		transaction, err := transactionBuilder.
			SetMessage(transaction.message[start:end]).
			SetTransactionID(nextTransactionID).
			SetTopicID(transaction.topicID).
			SetChunkInfo(initialTransactionID, int32(chunks), int32(i)+1).
			Build(client)

		if err != nil {
			return TransactionList{}, err
		}

		list[i] = transaction.List[0]
		nextTransactionID.ValidStart = nextTransactionID.ValidStart.Add(1 * time.Nanosecond)
	}

	return TransactionList{list}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction ConsensusMessageSubmitTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{
		transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee),
		transaction.pb,
		transaction.maxChunks,
		transaction.message,
		transaction.topicID,
		transaction.initialTransactionID,
		transaction.number,
		transaction.total,
		transaction.chunkInfoSet,
	}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction ConsensusMessageSubmitTransaction) SetTransactionMemo(memo string) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo),
		transaction.pb,
		transaction.maxChunks,
		transaction.message,
		transaction.topicID,
		transaction.initialTransactionID,
		transaction.number,
		transaction.total,
		transaction.chunkInfoSet,
	}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction ConsensusMessageSubmitTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration),
		transaction.pb,
		transaction.maxChunks,
		transaction.message,
		transaction.topicID,
		transaction.initialTransactionID,
		transaction.number,
		transaction.total,
		transaction.chunkInfoSet,
	}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction ConsensusMessageSubmitTransaction) SetTransactionID(transactionID TransactionID) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID),
		transaction.pb,
		transaction.maxChunks,
		transaction.message,
		transaction.topicID,
		transaction.initialTransactionID,
		transaction.number,
		transaction.total,
		transaction.chunkInfoSet,
	}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction ConsensusMessageSubmitTransaction) SetNodeID(nodeAccountID AccountID) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID),
		transaction.pb,
		transaction.maxChunks,
		transaction.message,
		transaction.topicID,
		transaction.initialTransactionID,
		transaction.number,
		transaction.total,
		transaction.chunkInfoSet,
	}
}
