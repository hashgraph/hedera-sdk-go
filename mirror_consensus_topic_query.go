package hedera

import (
	"context"
	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
	"time"
)

type MirrorConsensusTopicQuery struct {
	pb *mirror.ConsensusTopicQuery
}

type MirrorConsensusTopicResponse struct {
	ConsensusTimeStamp time.Time
	Message            []byte
	RunningHash        []byte
	SequenceNumber     uint64
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

	subClient, err := client.client.SubscribeTopic(context.TODO(), b.pb)

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
				// attempt a clean disconnect, but ignore if failed and stop listening
				_ = subClient.CloseSend()
				break
			}

			onNext(mirrorConsensusTopicResponseFromProto(resp))
		}
	}()

	return newMirrorSubscriptionHandle(subClient.CloseSend), nil
}
