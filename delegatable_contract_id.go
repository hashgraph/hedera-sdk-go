package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// DelegatableContractID is the ID for a Hedera smart contract
type DelegatableContractID struct {
	Shard    uint64
	Realm    uint64
	Contract uint64
	checksum *string
}

// DelegatableContractIDFromString constructs a DelegatableContractID from a string formatted as `Shard.Realm.Contract` (for example "0.0.3")
func DelegatableContractIDFromString(data string) (DelegatableContractID, error) {
	shard, realm, num, checksum, _, err := _IdFromString(data)
	if err != nil {
		return DelegatableContractID{}, err
	}

	return DelegatableContractID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Contract: uint64(num),
		checksum: checksum,
	}, nil
}

func (id *DelegatableContractID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil && client.network.networkName != nil {
		tempChecksum, err := _ChecksumParseAddress(client.network.networkName._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract))
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
func (id *DelegatableContractID) Validate(client *Client) error {
	if !id._IsZero() && client != nil && client.network.networkName != nil {
		tempChecksum, err := _ChecksumParseAddress(client.network.networkName._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract))
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

// DelegatableContractIDFromSolidityAddress constructs a DelegatableContractID from a string representation of a _Solidity address
func DelegatableContractIDFromSolidityAddress(s string) (DelegatableContractID, error) {
	shard, realm, contract, err := _IdFromSolidityAddress(s)
	if err != nil {
		return DelegatableContractID{}, err
	}

	return DelegatableContractID{
		Shard:    shard,
		Realm:    realm,
		Contract: contract,
	}, nil
}

// String returns the string representation of a DelegatableContractID formatted as `Shard.Realm.Contract` (for example "0.0.3")
func (id DelegatableContractID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract)
}

func (id DelegatableContractID) ToStringWithChecksum(client Client) (string, error) {
	if client.GetNetworkName() == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := _ChecksumParseAddress(client.GetNetworkName()._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Contract, checksum.correctChecksum), nil
}

// ToSolidityAddress returns the string representation of the DelegatableContractID as a _Solidity address.
func (id DelegatableContractID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.Contract)
}

func (id DelegatableContractID) _ToProtobuf() *services.ContractID {
	return &services.ContractID{
		ShardNum:    int64(id.Shard),
		RealmNum:    int64(id.Realm),
		ContractNum: int64(id.Contract),
	}
}

func _DelegatableContractIDFromProtobuf(delegatableContractID *services.ContractID) *DelegatableContractID {
	if delegatableContractID == nil {
		return nil
	}

	return &DelegatableContractID{
		Shard:    uint64(delegatableContractID.ShardNum),
		Realm:    uint64(delegatableContractID.RealmNum),
		Contract: uint64(delegatableContractID.ContractNum),
	}
}

func (id DelegatableContractID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Contract == 0
}

func (id DelegatableContractID) _ToProtoKey() *services.Key {
	return &services.Key{Key: &services.Key_DelegatableContractId{DelegatableContractId: id._ToProtobuf()}}
}

func (id DelegatableContractID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func DelegatableContractIDFromBytes(data []byte) (DelegatableContractID, error) {
	pb := services.ContractID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return DelegatableContractID{}, err
	}

	return *_DelegatableContractIDFromProtobuf(&pb), nil
}
