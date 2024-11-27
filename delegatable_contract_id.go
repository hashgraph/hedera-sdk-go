package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// ContractID is the ID for a Hiero smart contract
type DelegatableContractID struct {
	Shard      uint64
	Realm      uint64
	Contract   uint64
	EvmAddress []byte
	checksum   *string
}

// DelegatableContractIDFromString constructs a DelegatableContractID from a string formatted as `Shard.Realm.Contract` (for example "0.0.3")
func DelegatableContractIDFromString(data string) (DelegatableContractID, error) {
	shard, realm, num, checksum, evm, err := _ContractIDFromString(data)
	if err != nil {
		return DelegatableContractID{}, err
	}

	if num == -1 {
		return DelegatableContractID{
			Shard:      uint64(shard),
			Realm:      uint64(realm),
			Contract:   0,
			EvmAddress: evm,
			checksum:   checksum,
		}, nil
	}

	return DelegatableContractID{
		Shard:      uint64(shard),
		Realm:      uint64(realm),
		Contract:   uint64(num),
		EvmAddress: nil,
		checksum:   checksum,
	}, nil
}

// Verify that the client has a valid checksum.
func (id *DelegatableContractID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil {
		var tempChecksum _ParseAddressResult
		var err error
		if client.network.ledgerID != nil {
			tempChecksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract))
		}
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

// DelegatableContractIDFromEvmAddress constructs a DelegatableContractID from a string representation of a _Solidity address
func DelegatableContractIDFromEvmAddress(shard uint64, realm uint64, evmAddress string) (DelegatableContractID, error) {
	temp, err := hex.DecodeString(evmAddress)
	if err != nil {
		return DelegatableContractID{}, err
	}
	return DelegatableContractID{
		Shard:      shard,
		Realm:      realm,
		Contract:   0,
		EvmAddress: temp,
		checksum:   nil,
	}, nil
}

// DelegatableContractIDFromSolidityAddress constructs a DelegatableContractID from a string representation of a _Solidity address
// Does not populate DelegatableContractID.EvmAddress
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
	if len(id.EvmAddress) > 0 {
		temp := hex.EncodeToString(id.EvmAddress)
		return fmt.Sprintf("%d.%d.%s", id.Shard, id.Realm, temp)
	}
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract)
}

// ToStringWithChecksum returns the string representation of a DelegatableContractID formatted as `Shard.Realm.Contract-Checksum` (for example "0.0.3-abcde")
func (id DelegatableContractID) ToStringWithChecksum(client Client) (string, error) {
	if id.EvmAddress != nil {
		return "", errors.New("EvmAddress doesn't support checksums")
	}
	if client.GetNetworkName() == nil && client.GetLedgerID() == nil {
		return "", errNetworkNameMissing
	}
	var checksum _ParseAddressResult
	var err error
	if client.network.ledgerID != nil {
		checksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract))
	}
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
	resultID := services.ContractID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
	}

	if id.EvmAddress != nil {
		resultID.Contract = &services.ContractID_EvmAddress{EvmAddress: id.EvmAddress}
		return &resultID
	}

	resultID.Contract = &services.ContractID_ContractNum{ContractNum: int64(id.Contract)}

	return &resultID
}

func _DelegatableContractIDFromProtobuf(contractID *services.ContractID) *DelegatableContractID {
	if contractID == nil {
		return nil
	}
	resultID := DelegatableContractID{
		Shard: uint64(contractID.ShardNum),
		Realm: uint64(contractID.RealmNum),
	}

	switch id := contractID.Contract.(type) {
	case *services.ContractID_ContractNum:
		resultID.Contract = uint64(id.ContractNum)
		resultID.EvmAddress = nil
		return &resultID
	case *services.ContractID_EvmAddress:
		resultID.EvmAddress = id.EvmAddress
		resultID.Contract = 0
		return &resultID
	default:
		return &resultID
	}
}

func (id DelegatableContractID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Contract == 0
}

func (id DelegatableContractID) _ToProtoKey() *services.Key {
	return &services.Key{Key: &services.Key_ContractID{ContractID: id._ToProtobuf()}}
}

// ToBytes returns a byte array representation of the DelegatableContractID
func (id DelegatableContractID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// DelegatableContractIDFromBytes returns a DelegatableContractID generated from a byte array
func DelegatableContractIDFromBytes(data []byte) (DelegatableContractID, error) {
	pb := services.ContractID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return DelegatableContractID{}, err
	}

	return *_DelegatableContractIDFromProtobuf(&pb), nil
}
