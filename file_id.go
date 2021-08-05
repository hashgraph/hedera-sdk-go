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
	checksum *string
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
	shard, realm, num, checksum, err := idFromString(data)
	if err != nil {
		return FileID{}, err
	}

	return FileID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		File:     uint64(num),
		checksum: checksum,
	}, nil
}

func (id *FileID) Validate(client *Client) error {
	if !id.isZero() && client != nil && client.networkName != nil {
		tempChecksum, err := checksumParseAddress(client.networkName.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File))
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

func (id *FileID) setNetworkWithClient(client *Client) {
	if client.networkName != nil {
		id.setNetwork(*client.networkName)
	}
}

func (id *FileID) setNetwork(name NetworkName) {
	checksum := checkChecksum(name.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File))
	id.checksum = &checksum
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
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File)
}

func (id FileID) ToStringWithChecksum(client Client) (string, error) {
	if client.networkName == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := checksumParseAddress(client.networkName.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.File, checksum.correctChecksum), nil
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

func fileIDFromProtobuf(fileID *proto.FileID) FileID {
	if fileID == nil {
		return FileID{}
	}

	id := FileID{
		Shard: uint64(fileID.ShardNum),
		Realm: uint64(fileID.RealmNum),
		File:  uint64(fileID.FileNum),
	}

	return id
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
