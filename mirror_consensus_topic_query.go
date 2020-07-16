package hedera

import (
	"context"
	"sync"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
)

type ConsensusMessageMetadata struct {
	ConsensusTimeStamp time.Time
	RunningHash        []byte
	SequenceNumber     uint64
	ContentSize        uint64
}

type MirrorConsensusTopicQuery struct {
	pb *mirror.ConsensusTopicQuery
}

type MirrorConsensusTopicResponse struct {
	ConsensusTimeStamp time.Time
	Message            []byte
	RunningHash        []byte
	SequenceNumber     uint64
	Contents           []byte
	Metadata           []ConsensusMessageMetadata
}

func NewMirrorConsensusTopicQuery() *MirrorConsensusTopicQuery {
	pb := &mirror.ConsensusTopicQuery{}

	return &MirrorConsensusTopicQuery{pb}

}

func (b *MirrorConsensusTopicQuery) SetTopicID(topicID ConsensusTopicID) *MirrorConsensusTopicQuery {
	b.pb.TopicID = topicID.toProto()

	return b
}

func (b *MirrorConsensusTopicQuery) SetStartTime(time time.Time) *MirrorConsensusTopicQuery {
	b.pb.ConsensusStartTime = timeToProto(time)

	return b
}

func (b *MirrorConsensusTopicQuery) SetEndTime(time time.Time) *MirrorConsensusTopicQuery {
	b.pb.ConsensusEndTime = timeToProto(time)

	return b
}

func (b *MirrorConsensusTopicQuery) SetLimit(limit uint64) *MirrorConsensusTopicQuery {
	b.pb.Limit = limit

	return b
}

func mirrorConsensusTopicResponseFromProto(r *mirror.ConsensusTopicResponse) MirrorConsensusTopicResponse {
	resp := MirrorConsensusTopicResponse{
		ConsensusTimeStamp: timeFromProto(r.ConsensusTimestamp),
		Message:            r.Message,
		RunningHash:        r.RunningHash,
		SequenceNumber:     r.SequenceNumber,
		Contents:           r.Message,
		Metadata:           make([]ConsensusMessageMetadata, 1),
	}

	resp.Metadata = append(resp.Metadata, ConsensusMessageMetadata{
		ConsensusTimeStamp: resp.ConsensusTimeStamp,
		RunningHash:        resp.RunningHash,
		SequenceNumber:     resp.SequenceNumber,
		ContentSize:        uint64(len(r.Message)),
	})

	return resp
}

func mirrorConsensusTopicResponseFromChunkedProto(message []*mirror.ConsensusTopicResponse) MirrorConsensusTopicResponse {
	length := len(message)
	size := uint64(0)
	metadata := make([]ConsensusMessageMetadata, length)
	messages := make([][]byte, length)

	for _, m := range message {
		metadata[m.ChunkInfo.Number-1] = ConsensusMessageMetadata{
			ConsensusTimeStamp: timeFromProto(m.ConsensusTimestamp),
			RunningHash:        m.RunningHash,
			SequenceNumber:     m.SequenceNumber,
			ContentSize:        uint64(len(m.Message)),
		}

		messages[m.ChunkInfo.Number-1] = m.Message
		size += uint64(len(m.Message))
	}

	final_message := make([]byte, size)
	for _, m := range messages {
		final_message = append(final_message, m...)
	}

	return MirrorConsensusTopicResponse{
		ConsensusTimeStamp: metadata[length-1].ConsensusTimeStamp,
		RunningHash:        metadata[length-1].RunningHash,
		SequenceNumber:     metadata[length-1].SequenceNumber,
		Contents:           final_message,
		Metadata:           metadata,
	}
}

func (b *MirrorConsensusTopicQuery) Subscribe(
	client MirrorClient,
	onNext func(MirrorConsensusTopicResponse),
	onError func(error),
) (MirrorSubscriptionHandle, error) {
	if b.pb.TopicID == nil {
		return MirrorSubscriptionHandle{}, newErrLocalValidationf("topic ID was not provided")
	}

	ctx, cancel := context.WithCancel(context.TODO())

	subClient, err := client.client.SubscribeTopic(ctx, b.pb)

	messages := sync.Map{}
	messagesMutex := sync.Mutex{}

	if err != nil {
		return MirrorSubscriptionHandle{}, err
	}

	go func() {
		for {
			resp, err := subClient.Recv()

			if err != nil {
				if onError != nil {
					onError(err)
				}
				cancel()
				break
			}

			if resp.ChunkInfo == nil {
				onNext(mirrorConsensusTopicResponseFromProto(resp))
			} else {
				messagesMutex.Lock()
				txID := transactionIDFromProto(resp.ChunkInfo.InitialTransactionID)
				messageI, _ := messages.LoadOrStore(txID, make([]*mirror.ConsensusTopicResponse, resp.ChunkInfo.Total))
				message := messageI.([]*mirror.ConsensusTopicResponse)
				message = append(message, resp)

				if int32(len(message)) == resp.ChunkInfo.Total {
					messages.Delete(txID)
					messagesMutex.Unlock()
					onNext(mirrorConsensusTopicResponseFromChunkedProto(message))
				} else {
					messagesMutex.Unlock()
				}

			}
		}
	}()

	return newMirrorSubscriptionHandle(cancel), nil
}

//
// The QueryPayment functions are not required for this query
//
