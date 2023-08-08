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

// TopicMessage is a message from a topic}
type TopicMessage struct {
	ConsensusTimestamp time.Time
	Contents           []byte
	RunningHash        []byte
	SequenceNumber     uint64
	Chunks             []TopicMessageChunk
	TransactionID      *TransactionID
}

func _TopicMessageOfSingle(resp *mirror.ConsensusTopicResponse) TopicMessage {
	return TopicMessage{
		ConsensusTimestamp: _TimeFromProtobuf(resp.ConsensusTimestamp),
		Contents:           resp.Message,
		RunningHash:        resp.RunningHash,
		SequenceNumber:     resp.SequenceNumber,
		Chunks:             nil,
		TransactionID:      nil,
	}
}

func _TopicMessageOfMany(message []*mirror.ConsensusTopicResponse) TopicMessage {
	length := len(message)
	size := uint64(0)
	chunks := make([]TopicMessageChunk, length)
	messages := make([][]byte, length)
	var transactionID *TransactionID = nil

	for _, m := range message {
		if transactionID == nil {
			value := _TransactionIDFromProtobuf(m.ChunkInfo.InitialTransactionID)
			transactionID = &value
		}

		chunks[m.ChunkInfo.Number-1] = _NewTopicMessageChunk(m)
		messages[m.ChunkInfo.Number-1] = m.Message
		size += uint64(len(m.Message))
	}

	finalMessage := make([]byte, 0, size)

	for _, m := range messages {
		finalMessage = append(finalMessage, m...)
	}

	return TopicMessage{
		ConsensusTimestamp: _TimeFromProtobuf(message[length-1].ConsensusTimestamp),
		RunningHash:        message[length-1].RunningHash,
		SequenceNumber:     message[length-1].SequenceNumber,
		Contents:           finalMessage,
		Chunks:             chunks,
		TransactionID:      transactionID,
	}
}
