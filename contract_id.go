package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// ContractID is the ID for a Hiero smart contract
type ContractID struct {
	Shard      uint64
	Realm      uint64
	Contract   uint64
	EvmAddress []byte
	checksum   *string
}

// ContractIDFromString constructs a ContractID from a string formatted as `Shard.Realm.Contract` (for example "0.0.3")
func ContractIDFromString(data string) (ContractID, error) {
	shard, realm, num, checksum, evm, err := _ContractIDFromString(data)
	if err != nil {
		return ContractID{}, err
	}

	if num == -1 {
		return ContractID{
			Shard:      uint64(shard),
			Realm:      uint64(realm),
			Contract:   0,
			EvmAddress: evm,
			checksum:   checksum,
		}, nil
	}

	return ContractID{
		Shard:      uint64(shard),
		Realm:      uint64(realm),
		Contract:   uint64(num),
		EvmAddress: nil,
		checksum:   checksum,
	}, nil
}

// Verify that the client has a valid checksum.
func (id *ContractID) ValidateChecksum(client *Client) error {
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

// Deprecated - use ValidateChecksum instead
func (id *ContractID) Validate(client *Client) error {
	return id.ValidateChecksum(client)
}

// ContractIDFromEvmAddress constructs a ContractID from a string representation of an EVM address
func ContractIDFromEvmAddress(shard uint64, realm uint64, evmAddress string) (ContractID, error) {
	temp, err := hex.DecodeString(evmAddress)
	if err != nil {
		return ContractID{}, err
	}
	return ContractID{
		Shard:      shard,
		Realm:      realm,
		Contract:   0,
		EvmAddress: temp,
		checksum:   nil,
	}, nil
}

// ContractIDFromSolidityAddress constructs a ContractID from a string representation of a _Solidity address
// Does not populate ContractID.EvmAddress
// Deprecated
func ContractIDFromSolidityAddress(s string) (ContractID, error) {
	shard, realm, contract, err := _IdFromSolidityAddress(s)
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
	if len(id.EvmAddress) > 0 {
		temp := hex.EncodeToString(id.EvmAddress)
		return fmt.Sprintf("%d.%d.%s", id.Shard, id.Realm, temp)
	}
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract)
}

// ToStringWithChecksum returns the string representation of a ContractID formatted as `Shard.Realm.Contract-Checksum`
func (id ContractID) ToStringWithChecksum(client Client) (string, error) {
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

// ToSolidityAddress returns the string representation of the ContractID as a _Solidity address.
func (id ContractID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.Contract)
}

// PopulateContract gets the actual `Contract` field of the `ContractId` from the Mirror Node.
// Should be used after generating `ContractId.FromEvmAddress()` because it sets the `Contract` field to `0`
// automatically since there is no connection between the `Contract` and the `evmAddress`
func (id *ContractID) PopulateContract(client *Client) error {
	if client.mirrorNetwork == nil || len(client.GetMirrorNetwork()) == 0 {
		return errors.New("mirror node is not set")
	}
	mirrorUrl := client.GetMirrorNetwork()[0]
	index := strings.Index(mirrorUrl, ":")
	if index == -1 {
		return errors.New("invalid mirrorUrl format")
	}
	mirrorUrl = mirrorUrl[:index]
	url := fmt.Sprintf("https://%s/api/v1/contracts/%s", mirrorUrl, hex.EncodeToString(id.EvmAddress))
	if client.GetLedgerID().String() == "" {
		url = fmt.Sprintf("http://%s:5551/api/v1/contracts/%s", mirrorUrl, hex.EncodeToString(id.EvmAddress))
	}

	resp, err := http.Get(url) // #nosec
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	mirrorContractId, ok := result["contract_id"].(string)
	if !ok {
		return errors.New("unexpected response format")
	}

	numStr := mirrorContractId[strings.LastIndex(mirrorContractId, ".")+1:]
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return err
	}
	id.Contract = uint64(num)
	return nil
}

func (id ContractID) _ToProtobuf() *services.ContractID {
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

func _ContractIDFromProtobuf(contractID *services.ContractID) *ContractID {
	if contractID == nil {
		return nil
	}
	resultID := ContractID{
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

func (id ContractID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Contract == 0
}

func (id ContractID) _ToProtoKey() *services.Key {
	return &services.Key{Key: &services.Key_ContractID{ContractID: id._ToProtobuf()}}
}

// ToBytes returns a byte array representation of the ContractID
func (id ContractID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// ContractIDFromBytes returns a ContractID generated from a byte array
func ContractIDFromBytes(data []byte) (ContractID, error) {
	pb := services.ContractID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ContractID{}, err
	}

	return *_ContractIDFromProtobuf(&pb), nil
}
