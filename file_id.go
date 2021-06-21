package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// A FileID is the ID for a file on the network.
type FileID struct {
	Shard    uint64
	Realm    uint64
	File     uint64
	Checksum *string
	Network  *NetworkName
}

// FileIDForAddressBook returns the public node address book for the current network.
func FileIDForAddressBook() FileID {
	return FileID{File: 102}
}

// FileIDForFeeSchedule returns the current fee schedule for the network.
func FileIDForFeeSchedule() FileID {
	return FileID{File: 111}
}

// FileIDForExchangeRate returns the current exchange rates of HBAR to USD.
func FileIDForExchangeRate() FileID {
	return FileID{File: 112}
}

// FileIDFromString returns a FileID parsed from the given string.
// A malformatted string will cause this to return an error instead.
func FileIDFromString(data string) (FileID, error) {
	var checksum parseAddressResult
	var err error

	var networkNames = []NetworkName{
		Mainnet,
		Testnet,
		Previewnet,
	}

	var network NetworkName
	for _, name := range networkNames {
		checksum, err = checksumParseAddress(name.Network(), data)
		if err != nil {
			return FileID{}, err
		}
		if checksum.status != 1 {
			network = name
			break
		}
	}

	err = checksumVerify(checksum.status)
	if err != nil {
		return FileID{}, err
	}

	tempChecksum := checksum.correctChecksum
	return FileID{
		Shard:    uint64(checksum.num1),
		Realm:    uint64(checksum.num2),
		File:     uint64(checksum.num3),
		Checksum: &tempChecksum,
		Network:  &network,
	}, nil
}

func FileIDValidateNetworkOnIDs(id FileID, other AccountID) error {
	if !id.isZero() && !other.isZero() && id.Network != nil && other.Network != nil && *id.Network != *other.Network {
		return errNetworkMismatch
	}

	return nil
}

func FileIDFromSolidityAddress(s string) (FileID, error) {
	shard, realm, file, err := idFromSolidityAddress(s)
	if err != nil {
		return FileID{}, err
	}

	return FileID{
		Shard: shard,
		Realm: realm,
		File:  file,
	}, nil
}

func (id FileID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.File == 0
}

func (id FileID) String() string {
	checksum, err := checksumParseAddress("", fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File))
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.File, checksum.correctChecksum)
}

func (id FileID) ToSolidityAddress() string {
	return idToSolidityAddress(id.Shard, id.Realm, id.File)
}

func (id FileID) toProtobuf() *proto.FileID {
	return &proto.FileID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		FileNum:  int64(id.File),
	}
}

func fileIDFromProtobuf(pb *proto.FileID) FileID {
	if pb == nil {
		return FileID{}
	}
	return FileID{
		Shard: uint64(pb.ShardNum),
		Realm: uint64(pb.RealmNum),
		File:  uint64(pb.FileNum),
	}
}

func (id FileID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func FileIDFromBytes(data []byte) (FileID, error) {
	if data == nil {
		return FileID{}, errByteArrayNull
	}
	pb := proto.FileID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FileID{}, err
	}

	return fileIDFromProtobuf(&pb), nil
}
