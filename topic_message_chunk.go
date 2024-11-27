package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
)

// TopicMessageChunk is a chunk of a topic message
type TopicMessageChunk struct {
	ConsensusTimestamp time.Time
	ContentSize        uint64
	RunningHash        []byte
	SequenceNumber     uint64
}

func _NewTopicMessageChunk(resp *mirror.ConsensusTopicResponse) TopicMessageChunk {
	return TopicMessageChunk{
		ConsensusTimestamp: _TimeFromProtobuf(resp.ConsensusTimestamp),
		ContentSize:        uint64(len(resp.Message)),
		RunningHash:        resp.RunningHash,
		SequenceNumber:     resp.SequenceNumber,
	}
}
