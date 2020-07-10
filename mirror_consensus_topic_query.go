package hedera

import (
	"context"
	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
	"time"
)

type ChunkInfo struct {
    InitialTransactionID TransactionID
    Total                uint32
    Number               uint32
}

type MirrorConsensusTopicQuery struct {
	pb *mirror.ConsensusTopicQuery
}

type MirrorConsensusTopicResponse struct {
	ConsensusTimeStamp time.Time
	Message            []byte
	RunningHash        []byte
	SequenceNumber     uint64
    ChunkInfo         ChunkInfo
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
	return MirrorConsensusTopicResponse{
		ConsensusTimeStamp: timeFromProto(r.ConsensusTimestamp),
		Message:            r.Message,
		RunningHash:        r.RunningHash,
		SequenceNumber:     r.SequenceNumber,
        ChunkInfo: ChunkInfo{
            InitialTransactionID: transactionIDFromProto(r.ChunkInfo.InitialTransactionID),
            Total: uint32(r.ChunkInfo.Total),
            Number: uint32(r.ChunkInfo.Number),
        },
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

			onNext(mirrorConsensusTopicResponseFromProto(resp))
		}
	}()

	return newMirrorSubscriptionHandle(cancel), nil
}

//
// The QueryPayment functions are not required for this query
//
