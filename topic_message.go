package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
)

type TopicMessage struct {
	ConsensusTimestamp time.Time
	Contents           []byte
	RunningHash        []byte
	SequenceNumber     uint64
	Chunks             []TopicMessageChunk
	TransactionID      *TransactionID
}

func topicMessageOfSingle(resp *mirror.ConsensusTopicResponse) TopicMessage {
	return TopicMessage{
		ConsensusTimestamp: timeFromProtobuf(resp.ConsensusTimestamp),
		Contents:           resp.Message,
		RunningHash:        resp.RunningHash,
		SequenceNumber:     resp.SequenceNumber,
		Chunks:             nil,
		TransactionID:      nil,
	}
}

func topicMessageOfMany(message []*mirror.ConsensusTopicResponse) TopicMessage {
	length := len(message)
	size := uint64(0)
	chunks := make([]TopicMessageChunk, length)
	messages := make([][]byte, length)
	var transactionID *TransactionID = nil

	for _, m := range message {
		if transactionID == nil {
			value := transactionIDFromProtobuf(m.ChunkInfo.InitialTransactionID)
			transactionID = &value
		}

		chunks[m.ChunkInfo.Number-1] = newTopicMessageChunk(m)
		messages[m.ChunkInfo.Number-1] = m.Message
		size += uint64(len(m.Message))
	}

	finalMessage := make([]byte, 0, size)

	for _, m := range messages {
		finalMessage = append(finalMessage, m...)
	}

	return TopicMessage{
		ConsensusTimestamp: timeFromProtobuf(message[length-1].ConsensusTimestamp),
		RunningHash:        message[length-1].RunningHash,
		SequenceNumber:     message[length-1].SequenceNumber,
		Contents:           finalMessage,
		Chunks:             chunks,
		TransactionID:      transactionID,
	}
}
