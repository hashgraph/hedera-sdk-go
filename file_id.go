package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"
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
	if id.checksum == nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.File, *id.checksum)
}

func (id FileID) ToSolidityAddress() string {
	return idToSolidityAddress(id.Shard, id.Realm, id.File)
}

func (id FileID) toProtobuf() *services.FileID {
	return &services.FileID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		FileNum:  int64(id.File),
	}
}

func fileIDFromProtobuf(fileID *services.FileID, networkName *NetworkName) FileID {
	if fileID == nil {
		return FileID{}
	}

	id := FileID{
		Shard: uint64(fileID.ShardNum),
		Realm: uint64(fileID.RealmNum),
		File:  uint64(fileID.FileNum),
	}

	if networkName != nil {
		id.setNetwork(*networkName)
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
	pb := services.FileID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FileID{}, err
	}

	return fileIDFromProtobuf(&pb, nil), nil
}
