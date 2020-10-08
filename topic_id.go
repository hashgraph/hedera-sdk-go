package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// TopicID is a unique identifier for a topic (used by the  service)
type TopicID struct {
	Shard uint64
	Realm uint64
	Topic uint64
}

// TopicIDFromString constructs a TopicID from a string formatted as `Shard.Realm.Topic` (for example "0.0.3")
func TopicIDFromString(s string) (TopicID, error) {
	shard, realm, num, err := idFromString(s)
	if err != nil {
		return TopicID{}, err
	}

	return TopicID{
		Shard: uint64(shard),
		Realm: uint64(realm),
		Topic: uint64(num),
	}, nil
}

// String returns the string representation of a TopicID in `Shard.Realm.Topic` (for example "0.0.3")
func (id TopicID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic)
}

func (id TopicID) toProto() *proto.TopicID {
	return &proto.TopicID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TopicNum: int64(id.Topic),
	}
}

func TopicIDFromProto(pb *proto.TopicID) TopicID {
	return TopicID{
		Shard: uint64(pb.ShardNum),
		Realm: uint64(pb.RealmNum),
		Topic: uint64(pb.TopicNum),
	}
}
