package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
)

type ConsensusMessage struct {
	TopicID            ConsensusTopicID
	ConsensusTimestamp time.Time
	Message            []byte
	RunningHash        []byte
	SequenceNumber     uint64
}

func NewConsensusMessage(id ConsensusTopicID, resp *mirror.ConsensusTopicResponse) ConsensusMessage {
	return ConsensusMessage{
		TopicID:            id,
		ConsensusTimestamp: timeFromProto(resp.ConsensusTimestamp),
		Message:            resp.Message,
		RunningHash:        resp.RunningHash,
		SequenceNumber:     resp.SequenceNumber,
	}
}

func (message ConsensusMessage) String() string {
	return string(message.Message)
}
