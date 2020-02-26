package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ConsensusTopicID is a unique identifier for a topic (used by the consensus service)
type ConsensusTopicID struct {
	Shard uint64
	Realm uint64
	Topic uint64
}

// TopicIDFromString constructs a TopicID from a string formatted as `Shard.Realm.Topic` (for example "0.0.3")
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

// String returns the string representation of a TopicID in `Shard.Realm.Topic` (for example "0.0.3")
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
