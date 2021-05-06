package hedera

import (
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

const chunkSize = 1024

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

// NewConsensusMessageSubmitTransaction creates a ConsensusMessageSubmitTransaction builder which can be used to
// construct and execute a Consensus Submit Message Transaction.
func NewConsensusMessageSubmitTransaction() ConsensusMessageSubmitTransaction {
	pb := &proto.ConsensusSubmitMessageTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusSubmitMessage{ConsensusSubmitMessage: pb}

	builder := ConsensusMessageSubmitTransaction{inner, pb, 20, nil, ConsensusTopicID{}, TransactionID{}, 0, 0, false}

	return builder
}

func topicMessageSubmitTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) ConsensusMessageSubmitTransaction {
	tx := ConsensusMessageSubmitTransaction{
		TransactionBuilder:   transactionBuilder,
		pb:                   pb.GetConsensusSubmitMessage(),
		maxChunks:            20,
		message:              pb.GetConsensusSubmitMessage().GetMessage(),
		topicID:              consensusTopicIDFromProto(pb.GetConsensusSubmitMessage().GetTopicID()),
		initialTransactionID: transactionIDFromProto(pb.GetConsensusSubmitMessage().GetChunkInfo().InitialTransactionID),
		total:                pb.GetConsensusSubmitMessage().GetChunkInfo().GetTotal(),
		number:               pb.GetConsensusSubmitMessage().GetChunkInfo().GetNumber(),
		chunkInfoSet:         true,
	}

	return tx
}

// SetTopicID sets the topic to submit the message to.
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

func (builder ConsensusMessageSubmitTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *ConsensusMessageSubmitTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: &proto.ConsensusSubmitMessageTransactionBody{
				TopicID:   builder.topicID.toProto(),
				Message:   builder.message,
				ChunkInfo: &proto.ConsensusMessageChunkInfo{
					InitialTransactionID: builder.initialTransactionID.toProto(),
					Total:                builder.total,
					Number:               builder.number,
				},
			},
		},
	}, nil
}

func (builder ConsensusMessageSubmitTransaction) Execute(client *Client) (TransactionID, error) {
	txs, err := builder.Build(client)
	if err != nil {
		return TransactionID{}, err
	}

	return txs.Execute(client)
}

func (builder ConsensusMessageSubmitTransaction) ExecuteAll(client *Client) ([]TransactionID, error) {
	txs, err := builder.Build(client)
	if err != nil {
		return nil, err
	}

	return txs.ExecuteAll(client)
}

func (builder ConsensusMessageSubmitTransaction) Build(client *Client) (TransactionList, error) {
	// If chunk info  is set then we aren't going to chunk the message
	// Set all the required fields and return a list of 1
	if builder.chunkInfoSet {
		builder.pb.TopicID = builder.topicID.toProto()
		builder.pb.Message = builder.message
		builder.pb.ChunkInfo = &proto.ConsensusMessageChunkInfo{
			InitialTransactionID: builder.initialTransactionID.toProto(),
			Number:               builder.number,
			Total:                builder.total,
		}

		// FIXME: really have no idea why this is needed @daniel
		builder.TransactionBuilder.pb.Data = &proto.TransactionBody_ConsensusSubmitMessage{ConsensusSubmitMessage: builder.pb}

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

	var initialTransactionID TransactionID
	if builder.TransactionBuilder.pb.TransactionID != nil {
		initialTransactionID = transactionIDFromProto(builder.TransactionBuilder.pb.TransactionID)
	} else {
		initialTransactionID = NewTransactionID(client.GetOperatorID())
	}

	nextTransactionID := initialTransactionID

	for i := 0; uint64(i) < chunks; i += 1 {
		start := i * chunkSize
		end := start + chunkSize

		if end > len(builder.message) {
			end = len(builder.message)
		}

		transactionBuilder := NewConsensusMessageSubmitTransaction()
		transactionBuilder.TransactionBuilder.pb = protobuf.Clone(builder.TransactionBuilder.pb).(*proto.TransactionBody)

		transaction, err := transactionBuilder.
			SetMessage(builder.message[start:end]).
			SetTransactionID(nextTransactionID).
			SetTopicID(builder.topicID).
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

// SetTransactionMemo sets the memo for this Transaction.
func (builder ConsensusMessageSubmitTransaction) SetTransactionMemo(memo string) ConsensusMessageSubmitTransaction {
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

// SetTransactionValidDuration sets the valid duration for this Transaction.
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

// SetTransactionID sets the TransactionID for this Transaction.
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

// SetNodeAccountID sets the node AccountID for this Transaction.
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
