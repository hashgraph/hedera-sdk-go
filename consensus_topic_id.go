package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ConsensusTopicID struct {
	Shard uint64
	Realm uint64
	Topic uint64
}

func TopicIDFromString(s string) (ConsensusTopicID, error) {
	shard, realm, num, err := idFromString(s)
	if err != nil {
		return ConsensusTopicID{}, err
	}

	return ConsensusTopicID{
		Shard: uint64(shard),
		Realm: uint64(realm),
		Topic: uint64(num),
	}, nil
}

func (id ConsensusTopicID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic)
}

func (id ConsensusTopicID) toProto() *proto.TopicID {
	return &proto.TopicID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TopicNum: int64(id.Topic),
	}
}

func consensusTopicIDFromProto(pb *proto.TopicID) ConsensusTopicID {
	return ConsensusTopicID{
		Shard: uint64(pb.ShardNum),
		Realm: uint64(pb.RealmNum),
		Topic: uint64(pb.TopicNum),
	}
}
