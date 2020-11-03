package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
)

type TopicMessageChunk struct {
	ConsensusTimestamp time.Time
	ContentSize        uint64
	RunningHash        []byte
	SequenceNumber     uint64
}

func newTopicMessageChunk(resp *mirror.ConsensusTopicResponse) TopicMessageChunk {
	return TopicMessageChunk{
		ConsensusTimestamp: timeFromProtobuf(resp.ConsensusTimestamp),
		ContentSize:        uint64(len(resp.Message)),
		RunningHash:        resp.RunningHash,
		SequenceNumber:     resp.SequenceNumber,
	}
}
