package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

const chunkSize = 4096

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

// NewConsensusMessageSubmitTransaction creates a ConsensusMessageSubmitTransaction builder which can be used to
// construct and execute a Consensus Submit Message Transaction.
func NewConsensusMessageSubmitTransaction() ConsensusMessageSubmitTransaction {
	pb := &proto.ConsensusSubmitMessageTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusSubmitMessage{pb}

	builder := ConsensusMessageSubmitTransaction{inner, pb, 10, nil, ConsensusTopicID{}, TransactionID{}, 0, 0, false}

	return builder
}

// SetTopic sets the topic to submit the message to.
func (builder ConsensusMessageSubmitTransaction) SetTopicID(id ConsensusTopicID) ConsensusMessageSubmitTransaction {
	builder.topicID = id
	return builder
}

// SetMessage sets the message to be submitted. Max size of the Transaction (including signatures) is 4kB.
func (builder ConsensusMessageSubmitTransaction) SetMessage(message []byte) ConsensusMessageSubmitTransaction {
	builder.message = message
	return builder
}

// SetMessage sets the message to be submitted. Max size of the Transaction (including signatures) is 4kB.
func (builder ConsensusMessageSubmitTransaction) SetMaxChunks(max uint64) ConsensusMessageSubmitTransaction {
	builder.maxChunks = max
	return builder
}

// SetMessage sets the message to be submitted. Max size of the Transaction (including signatures) is 4kB.
func (builder ConsensusMessageSubmitTransaction) SetChunkInfo(transactionID TransactionID, total int32, number int32) ConsensusMessageSubmitTransaction {
	builder.initialTransactionID = transactionID
	builder.total = total
	builder.number = number
	builder.chunkInfoSet = true
	return builder
}

func (builder ConsensusMessageSubmitTransaction) Build(client *Client) (TransactionList, error) {
	var transactionID TransactionID
	if builder.TransactionBuilder.pb.TransactionID != nil {
		transactionID = transactionIDFromProto(builder.TransactionBuilder.pb.TransactionID)
	} else {
        if client != nil {
            transactionID = NewTransactionID(client.GetOperatorID())
        } else {
            return TransactionList{}, fmt.Errorf("client must have an operator or set a transaction ID to build a consensus message transaction")
        }
	}

	// If chunk info  is set then we aren't going to chunk the message
	// Set all the required fields and return a list of 1
	if builder.chunkInfoSet || len(builder.message) < chunkSize {
		builder.pb.TopicID = builder.topicID.toProto()
		builder.pb.Message = builder.message
		builder.pb.ChunkInfo = &proto.ConsensusMessageChunkInfo{
			InitialTransactionID: builder.initialTransactionID.toProto(),
			Number:               builder.number,
			Total:                builder.total,
		}

		transaction, err := builder.TransactionBuilder.Build(client)
		if err != nil {
			return TransactionList{}, err
		}

		list := TransactionList{
			List: make([]Transaction, 1),
		}

		list.List[0] = transaction
		return list, nil
	}

	chunks := uint64((len(builder.message) + (chunkSize - 1)) / chunkSize)

	if chunks > builder.maxChunks {
		return TransactionList{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: builder.maxChunks,
		}
	}

	list := make([]Transaction, chunks)

	for i := 0; uint64(i) < chunks; i += 1 {
		start := i * chunkSize
		end := start + chunkSize

		if end > len(builder.message) {
			end = len(builder.message)
		}

		transaction, err := NewConsensusMessageSubmitTransaction().
			SetMessage(builder.message[start:end]).
			SetTopicID(builder.topicID).
			SetChunkInfo(transactionID, int32(chunks), int32(i)+1).
            SetTransactionID(transactionID).
            SetNodeAccountID(accountIDFromProto(builder.TransactionBuilder.pb.NodeAccountID)).
			Build(client)

        transactionID = NewTransactionIDWithValidStart(transactionID.AccountID, transactionID.ValidStart.Add(1 * time.Nanosecond))

		if err != nil {
			return TransactionList{}, err
		}

		list[i] = transaction.List[0]
	}

	return TransactionList{list}, nil
}

func (builder ConsensusMessageSubmitTransaction) Execute(client *Client) ([]TransactionID, error) {
    list, err := builder.Build(client)
    if err != nil {
        return nil, err
    }

    ids := make([]TransactionID, len(list.List))

    for i, tx := range list.List {
        result, err := tx.Execute(client);
        if err != nil {
            return nil, err
        }

        ids[i] = result
    }

    return ids, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ConsensusMessageSubmitTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{
		builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee),
		builder.pb,
		builder.maxChunks,
		builder.message,
		builder.topicID,
		builder.initialTransactionID,
		builder.number,
		builder.total,
		builder.chunkInfoSet,
	}
}

func (builder ConsensusMessageSubmitTransaction) SetMemo(memo string) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetTransactionMemo(memo),
		builder.pb,
		builder.maxChunks,
		builder.message,
		builder.topicID,
		builder.initialTransactionID,
		builder.number,
		builder.total,
		builder.chunkInfoSet,
	}
}

func (builder ConsensusMessageSubmitTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration),
		builder.pb,
		builder.maxChunks,
		builder.message,
		builder.topicID,
		builder.initialTransactionID,
		builder.number,
		builder.total,
		builder.chunkInfoSet,
	}
}

func (builder ConsensusMessageSubmitTransaction) SetTransactionID(transactionID TransactionID) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetTransactionID(transactionID),
		builder.pb,
		builder.maxChunks,
		builder.message,
		builder.topicID,
		builder.initialTransactionID,
		builder.number,
		builder.total,
		builder.chunkInfoSet,
	}
}

func (builder ConsensusMessageSubmitTransaction) SetNodeAccountID(nodeAccountID AccountID) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID),
		builder.pb,
		builder.maxChunks,
		builder.message,
		builder.topicID,
		builder.initialTransactionID,
		builder.number,
		builder.total,
		builder.chunkInfoSet,
	}
}
