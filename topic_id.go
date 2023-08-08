package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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

// Verify that the client has a valid checksum.
func (id *TopicID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil {
		tempChecksum, err := _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic))
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
			networkName := NetworkNameOther
			if client.network.ledgerID != nil {
				networkName, _ = client.network.ledgerID.ToNetworkName()
			}
			return errors.New(fmt.Sprintf("network mismatch or wrong checksum given, given checksum: %s, correct checksum %s, network: %s",
				*id.checksum,
				tempChecksum.correctChecksum,
				networkName))
		}
	}

	return nil
}

// Deprecated - use ValidateChecksum instead
func (id *TopicID) Validate(client *Client) error {
	return id.ValidateChecksum(client)
}

func (id TopicID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Topic == 0
}

// String returns the string representation of a TopicID in `Shard.Realm.Topic` (for example "0.0.3")
func (id TopicID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic)
}

// ToStringWithChecksum returns the string representation of a TopicID in `Shard.Realm.Topic-Checksum` (for example "0.0.3-abcde")
func (id TopicID) ToStringWithChecksum(client Client) (string, error) {
	if client.GetNetworkName() == nil && client.GetLedgerID() == nil {
		return "", errNetworkNameMissing
	}
	var checksum _ParseAddressResult
	var err error
	if client.network.ledgerID != nil {
		checksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Topic))
	}
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

// ToBytes returns a byte array representation of the TopicID
func (id TopicID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TopicIDFromBytes constructs a TopicID from a byte array
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

// TopicIDFromSolidityAddress constructs an TopicID from a string
// representation of a _Solidity address
func TopicIDFromSolidityAddress(s string) (TopicID, error) {
	shard, realm, topic, err := _IdFromSolidityAddress(s)
	if err != nil {
		return TopicID{}, err
	}

	return TopicID{
		Shard:    shard,
		Realm:    realm,
		Topic:    topic,
		checksum: nil,
	}, nil
}

// ToSolidityAddress returns the string representation of the TopicID as a
// _Solidity address.
func (id TopicID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.Topic)
}
