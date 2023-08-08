package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/mirror"
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
