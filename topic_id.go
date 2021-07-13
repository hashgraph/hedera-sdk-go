package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TopicID is a unique identifier for a topic (used by the  service)
type TopicID struct {
	Shard    uint64
	Realm    uint64
	Topic    uint64
	checksum *string
}

// TopicIDFromString constructs a TopicID from a string formatted as `Shard.Realm.Topic` (for example "0.0.3")
func TopicIDFromString(data string) (TopicID, error) {
	shard, realm, num, checksum, err := idFromString(data)
	if err != nil {
		return TopicID{}, err
	}

	return TopicID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Topic:    uint64(num),
		checksum: checksum,
	}, nil
}

func (id *TopicID) Validate(client *Client) error {
	if !id.isZero() && client != nil && client.networkName != nil {
		tempChecksum, err := checksumParseAddress(client.networkName.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic))
		if err != nil {
			return err
		}
		err = checksumVerify(tempChecksum.status)
		if err != nil {
			return err
		}
		if id.checksum == nil {
			id.checksum = &tempChecksum.correctChecksum
			return nil
		}
		if tempChecksum.correctChecksum != *id.checksum {
			return errNetworkMismatch
		}
	}

	return nil
}

func (id *TopicID) setNetworkWithClient(client *Client) {
	if client.networkName != nil {
		id.setNetwork(*client.networkName)
	}
}

func (id *TopicID) setNetwork(name NetworkName) {
	checksum := checkChecksum(name.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic))
	id.checksum = &checksum
}
func (id TopicID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Topic == 0
}

// String returns the string representation of a TopicID in `Shard.Realm.Topic` (for example "0.0.3")
func (id TopicID) String() string {
	if id.checksum == nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Topic, *id.checksum)
}

func (id TopicID) toProtobuf() *services.TopicID {
	return &services.TopicID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TopicNum: int64(id.Topic),
	}
}

func topicIDFromProtobuf(topicID *services.TopicID, networkName *NetworkName) TopicID {
	if topicID == nil {
		return TopicID{}
	}

	id := TopicID{
		Shard: uint64(topicID.ShardNum),
		Realm: uint64(topicID.RealmNum),
		Topic: uint64(topicID.TopicNum),
	}

	if networkName != nil {
		id.setNetwork(*networkName)
	}

	return id
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
	pb := services.TopicID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TopicID{}, err
	}

	return topicIDFromProtobuf(&pb, nil), nil
}
