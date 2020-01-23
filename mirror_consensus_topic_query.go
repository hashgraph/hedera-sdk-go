package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
	"time"
)

type MirrorConsensusTopicQuery struct {
	pb *mirror.ConsensusTopicQuery
}

type MirrorConsensusTopicResponse struct {
	ConsensusTimeStamp 	time.Time
	Message 			[]byte
	RunningHash 		[]byte
	SequenceNumber 		uint64
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

/*
func (b *MirrorConsensusTopicQuery) Subscribe(client MirrorClient, onNext func(MirrorConsensusTopicResponse), onError *func()) (MirrorSubscriptionHandle, error) {
	name := methodName()

	client.conn.NewStream(context.TODO(), grpc.StreamDesc{
		StreamName:    "",
		Handler:       nil,
		ServerStreams: false,
		ClientStreams: false,
	}, )


}
*/
