package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// TopicID is a unique identifier for a topic (used by the  service)
type TopicID struct {
	Shard    uint64
	Realm    uint64
	Topic    uint64
	Checksum *string
	Network  *NetworkName
}

// TopicIDFromString constructs a TopicID from a string formatted as `Shard.Realm.Topic` (for example "0.0.3")
func TopicIDFromString(data string) (TopicID, error) {
	var checksum parseAddressResult
	var err error

	var networkNames = []NetworkName{
		NetworkNameMainnet,
		NetworkNameTestnet,
		NetworkNamePreviewnet,
	}

	var network NetworkName
	for _, name := range networkNames {
		checksum, err = checksumParseAddress(name.Network(), data)
		if err != nil {
			return TopicID{}, err
		}
		if checksum.status != 1 {
			network = name
			break
		}
	}

	err = checksumVerify(checksum.status)
	if err != nil {
		return TopicID{}, err
	}

	tempChecksum := checksum.correctChecksum

	return TopicID{
		Shard:    uint64(checksum.num1),
		Realm:    uint64(checksum.num2),
		Topic:    uint64(checksum.num3),
		Checksum: &tempChecksum,
		Network:  &network,
	}, nil
}

func TopicIDValidateNetworkOnIDs(id TopicID, other *Client) error {
	if !id.isZero() && other != nil && id.Network != nil && other.networkName != nil && *id.Network != *other.networkName {
		return errNetworkMismatch
	}

	return nil
}

func (id *TopicID) SetNetworkName(network NetworkName) {
	id.Network = &network
	checksum := checkChecksum(id.Network.Network(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic))
	id.Checksum = &checksum
}

func (id TopicID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Topic == 0
}

// String returns the string representation of a TopicID in `Shard.Realm.Topic` (for example "0.0.3")
func (id TopicID) String() string {
	if id.Network == nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Topic, *id.Checksum)
}

func (id TopicID) toProtobuf() *proto.TopicID {
	return &proto.TopicID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TopicNum: int64(id.Topic),
	}
}

func topicIDFromProtobuf(pb *proto.TopicID) TopicID {
	if pb == nil {
		return TopicID{}
	}
	return TopicID{
		Shard: uint64(pb.ShardNum),
		Realm: uint64(pb.RealmNum),
		Topic: uint64(pb.TopicNum),
	}
}

func (id TopicID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TopicIDFromBytes(data []byte) (TopicID, error) {
	if data == nil {
		return TopicID{}, errByteArrayNull
	}
	pb := proto.TopicID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TopicID{}, err
	}

	return topicIDFromProtobuf(&pb), nil
}
