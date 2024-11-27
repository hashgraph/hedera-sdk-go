package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// A FileID is the ID for a file on the _Network.
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
	shard, realm, num, checksum, err := _IdFromString(data)
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

// Verify that the client has a valid checksum.
func (id *FileID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil {
		var tempChecksum _ParseAddressResult
		var err error
		tempChecksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File))
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
func (id *FileID) Validate(client *Client) error {
	return id.ValidateChecksum(client)
}

// FileIDFromSolidityAddress returns a FileID parsed from the given solidity address.
func FileIDFromSolidityAddress(s string) (FileID, error) {
	shard, realm, file, err := _IdFromSolidityAddress(s)
	if err != nil {
		return FileID{}, err
	}

	return FileID{
		Shard: shard,
		Realm: realm,
		File:  file,
	}, nil
}

func (id FileID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.File == 0
}

// String returns the string representation of a FileID in the format used within protobuf.
func (id FileID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File)
}

// ToStringWithChecksum returns the string representation of a FileId with checksum.
func (id FileID) ToStringWithChecksum(client Client) (string, error) {
	if client.GetNetworkName() == nil && client.GetLedgerID() == nil {
		return "", errNetworkNameMissing
	}
	var checksum _ParseAddressResult
	var err error
	if client.network.ledgerID != nil {
		checksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File))
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.File, checksum.correctChecksum), nil
}

// ToSolidityAddress returns the string representation of a FileID in the format used by Solidity.
func (id FileID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.File)
}

func (id FileID) _ToProtobuf() *services.FileID {
	return &services.FileID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		FileNum:  int64(id.File),
	}
}

func _FileIDFromProtobuf(fileID *services.FileID) *FileID {
	if fileID == nil {
		return nil
	}

	return &FileID{
		Shard: uint64(fileID.ShardNum),
		Realm: uint64(fileID.RealmNum),
		File:  uint64(fileID.FileNum),
	}
}

// ToBytes returns a byte array representation of the FileID
func (id FileID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// FileIDFromBytes returns a FileID from a byte array
func FileIDFromBytes(data []byte) (FileID, error) {
	if data == nil {
		return FileID{}, errByteArrayNull
	}
	pb := services.FileID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FileID{}, err
	}

	return *_FileIDFromProtobuf(&pb), nil
}
