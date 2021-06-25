package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// ContractID is the ID for a Hedera smart contract
type ContractID struct {
	Shard    uint64
	Realm    uint64
	Contract uint64
	checksum *string
}

// ContractIDFromString constructs a ContractID from a string formatted as `Shard.Realm.Contract` (for example "0.0.3")
func ContractIDFromString(data string) (ContractID, error) {
	shard, realm, num, checksum, err := idFromString(data)
	if err != nil {
		return ContractID{}, err
	}

	return ContractID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Contract: uint64(num),
		checksum: checksum,
	}, nil
}

func (id *ContractID) Validate(client *Client) error {
	if !id.isZero() && client != nil && client.networkName != nil {
		tempChecksum, err := checksumParseAddress(client.networkName.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract))
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

func (id *ContractID) setNetworkWithClient(client *Client) {
	if client.networkName != nil {
		id.setNetwork(*client.networkName)
	}
}

func (id *ContractID) setNetwork(name NetworkName) {
	checksum := checkChecksum(name.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract))
	id.checksum = &checksum
}

// ContractIDFromSolidityAddress constructs a ContractID from a string representation of a solidity address
func ContractIDFromSolidityAddress(s string) (ContractID, error) {
	shard, realm, contract, err := idFromSolidityAddress(s)
	if err != nil {
		return ContractID{}, err
	}

	return ContractID{
		Shard:    shard,
		Realm:    realm,
		Contract: contract,
	}, nil
}

// String returns the string representation of a ContractID formatted as `Shard.Realm.Contract` (for example "0.0.3")
func (id ContractID) String() string {
	if id.checksum == nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Contract, *id.checksum)
}

// ToSolidityAddress returns the string representation of the ContractID as a solidity address.
func (id ContractID) ToSolidityAddress() string {
	return idToSolidityAddress(id.Shard, id.Realm, id.Contract)
}

func (id ContractID) toProtobuf() *proto.ContractID {
	return &proto.ContractID{
		ShardNum:    int64(id.Shard),
		RealmNum:    int64(id.Realm),
		ContractNum: int64(id.Contract),
	}
}

func contractIDFromProtobuf(contractID *proto.ContractID, networkName *NetworkName) ContractID {
	if contractID == nil {
		return ContractID{}
	}

	id := ContractID{
		Shard:    uint64(contractID.ShardNum),
		Realm:    uint64(contractID.RealmNum),
		Contract: uint64(contractID.ContractNum),
	}

	if networkName != nil {
		id.setNetwork(*networkName)
	}

	return id
}

func (id ContractID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Contract == 0
}

func (id ContractID) toProtoKey() *proto.Key {
	return &proto.Key{Key: &proto.Key_ContractID{ContractID: id.toProtobuf()}}
}

func (id ContractID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func ContractIDFromBytes(data []byte) (ContractID, error) {
	pb := proto.ContractID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ContractID{}, err
	}

	return contractIDFromProtobuf(&pb, nil), nil
}
