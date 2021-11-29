package hedera

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
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
	shard, realm, num, checksum, err := _IdFromString(data)
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

func (id *TopicID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil && client.network.networkName != nil {
		tempChecksum, err := _ChecksumParseAddress(client.network.networkName._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic))
		if err != nil {
			return err
		}
		err = _ChecksumVerify(tempChecksum.status)
		if err != nil {
			return err
		}
		if id.checksum == nil {
			return errChecksumMissing
		}
		if tempChecksum.correctChecksum != *id.checksum {
			return errors.New(fmt.Sprintf("network mismatch or wrong checksum given, given checksum: %s, correct checksum %s, network: %s",
				*id.checksum,
				tempChecksum.correctChecksum,
				*client.network.networkName))
		}
	}

	return nil
}

// Deprecated
func (id *TopicID) Validate(client *Client) error {
	if !id._IsZero() && client != nil && client.GetNetworkName() != nil {
		tempChecksum, err := _ChecksumParseAddress(client.GetNetworkName()._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic))
		if err != nil {
			return err
		}
		err = _ChecksumVerify(tempChecksum.status)
		if err != nil {
			return err
		}
		if id.checksum == nil {
			return errChecksumMissing
		}
		if tempChecksum.correctChecksum != *id.checksum {
			return errors.New(fmt.Sprintf("network mismatch or wrong checksum given, given checksum: %s, correct checksum %s, network: %s",
				*id.checksum,
				tempChecksum.correctChecksum,
				*client.network.networkName))
		}
	}

	return nil
}

func (id TopicID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Topic == 0
}

// String returns the string representation of a TopicID in `Shard.Realm.Topic` (for example "0.0.3")
func (id TopicID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic)
}

func (id TopicID) ToStringWithChecksum(client Client) (string, error) {
	if client.GetNetworkName() == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := _ChecksumParseAddress(client.GetNetworkName()._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Topic, checksum.correctChecksum), nil
}

func (id TopicID) _ToProtobuf() *services.TopicID {
	return &services.TopicID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TopicNum: int64(id.Topic),
	}
}

func _TopicIDFromProtobuf(topicID *services.TopicID) *TopicID {
	if topicID == nil {
		return nil
	}

	return &TopicID{
		Shard: uint64(topicID.ShardNum),
		Realm: uint64(topicID.RealmNum),
		Topic: uint64(topicID.TopicNum),
	}
}

func (id TopicID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
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

	return *_TopicIDFromProtobuf(&pb), nil
}
