package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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

	"github.com/hashgraph/hedera-protobufs-go/services"
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

func (id *FileID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil && client.network.ledgerID != nil {
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
			temp, _ := client.network.ledgerID.ToNetworkName()
			return errors.New(fmt.Sprintf("network mismatch or wrong checksum given, given checksum: %s, correct checksum %s, network: %s",
				*id.checksum,
				tempChecksum.correctChecksum,
				temp))
		}
	}

	return nil
}

// Deprecated
func (id *FileID) Validate(client *Client) error {
	if !id._IsZero() && client != nil && client.GetNetworkName() != nil {
		tempChecksum, err := _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File))
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
			temp, _ := client.network.ledgerID.ToNetworkName()
			return errors.New(fmt.Sprintf("network mismatch or wrong checksum given, given checksum: %s, correct checksum %s, network: %s",
				*id.checksum,
				tempChecksum.correctChecksum,
				temp))
		}
	}

	return nil
}

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

func (id FileID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File)
}

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

func (id FileID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
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

	return *_FileIDFromProtobuf(&pb), nil
}
