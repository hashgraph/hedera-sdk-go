package hedera

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConsensusMessageChunk struct {
	ConsensusTimestamp time.Time
	RunningHash        []byte
	SequenceNumber     uint64
	ContentSize        uint64
}

type MirrorConsensusTopicQuery struct {
	pb *mirror.ConsensusTopicQuery
}

type MirrorConsensusTopicResponse struct {
	ConsensusTimestamp time.Time
	Message            []byte
	RunningHash        []byte
	SequenceNumber     uint64
	Chunks             []ConsensusMessageChunk
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
		ConsensusTimestamp: timeFromProto(r.ConsensusTimestamp),
		Message:            r.Message,
		RunningHash:        r.RunningHash,
		SequenceNumber:     r.SequenceNumber,
		Chunks:             make([]ConsensusMessageChunk, 0, 1),
	}

	resp.Chunks = append(resp.Chunks, ConsensusMessageChunk{
		ConsensusTimestamp: resp.ConsensusTimestamp,
		RunningHash:        resp.RunningHash,
		SequenceNumber:     resp.SequenceNumber,
		ContentSize:        uint64(len(r.Message)),
	})

	return resp
}

func mirrorConsensusTopicResponseFromChunkedProto(message []*mirror.ConsensusTopicResponse) MirrorConsensusTopicResponse {
	length := len(message)
	size := uint64(0)
	chunks := make([]ConsensusMessageChunk, length)
	messages := make([][]byte, length)

	for _, m := range message {
		chunks[m.ChunkInfo.Number-1] = ConsensusMessageChunk{
			ConsensusTimestamp: timeFromProto(m.ConsensusTimestamp),
			RunningHash:        m.RunningHash,
			SequenceNumber:     m.SequenceNumber,
			ContentSize:        uint64(len(m.Message)),
		}

		messages[m.ChunkInfo.Number-1] = m.Message
		size += uint64(len(m.Message))
	}

	finalMessage := make([]byte, 0, size)

	for _, m := range messages {
		finalMessage = append(finalMessage, m...)
	}

	return MirrorConsensusTopicResponse{
		ConsensusTimestamp: timeFromProto(message[length-1].ConsensusTimestamp),
		RunningHash:        message[length-1].RunningHash,
		SequenceNumber:     message[length-1].SequenceNumber,
		Message:            finalMessage,
		Chunks:             chunks,
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

	handle := newMirrorSubscriptionHandle(cancel)

	messages := sync.Map{}
	messagesMutex := sync.Mutex{}

	go func() {
		var subClient mirror.ConsensusService_SubscribeTopicClient
		attempt := 0
		shouldRetry := true

		for {
			if shouldRetry {
				subClient, _ = client.client.SubscribeTopic(ctx, b.pb)
			}

			resp, err := subClient.Recv()

			code := status.Code(err)

			if err != nil {
				if attempt >= 10 && !shouldRetry && (code == codes.NotFound || code == codes.Unavailable) {
					if onError != nil {
						onError(err)
					}
					cancel()
					break
				} else if shouldRetry && code != codes.NotFound && code != codes.Unavailable {
					delay := 250.0 * math.Pow(2.0, float64(attempt))
					time.Sleep(time.Duration(delay) * time.Millisecond)
					attempt += 1
					continue
				} else {
					if onError != nil {
						onError(err)
					}
					cancel()
					break
				}
			}

			shouldRetry = false

			if resp.ChunkInfo == nil {
				onNext(mirrorConsensusTopicResponseFromProto(resp))
			} else {
				messagesMutex.Lock()
				txID := transactionIDFromProto(resp.ChunkInfo.InitialTransactionID)
				messageI, _ := messages.LoadOrStore(txID, make([]*mirror.ConsensusTopicResponse, 0, resp.ChunkInfo.Total))

				message := messageI.([]*mirror.ConsensusTopicResponse)
				message = append(message, resp)

				messages.Store(txID, message)

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

	return handle, nil
}

//
// The QueryPayment functions are not required for this query
//
