package hedera

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ConsensusTopicID struct {
	Shard uint64
	Realm uint64
	Topic uint64
}

func ConsensusTopicIDFromString(s string) (ConsensusTopicID, error) {
	values := strings.SplitN(s, ".", 3)
	if len(values) != 3 {
		// Was not three values separated by periods
		return ConsensusTopicID{}, fmt.Errorf("expected {shard}.{realm}.{num}")
	}

	shard, err := strconv.Atoi(values[0])
	if err != nil {
		return ConsensusTopicID{}, err
	}

	realm, err := strconv.Atoi(values[1])
	if err != nil {
		return ConsensusTopicID{}, err
	}

	num, err := strconv.Atoi(values[2])
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
