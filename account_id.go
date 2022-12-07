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
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// AccountID is the ID for a Hedera account
type AccountID struct {
	Shard           uint64
	Realm           uint64
	Account         uint64
	AliasKey        *PublicKey
	AliasEvmAddress *[]byte
	checksum        *string
}

type _AccountIDs struct { //nolint
	accountIDs []AccountID
}

// AccountIDFromString constructs an AccountID from a string formatted as
// `Shard.Realm.Account` (for example "0.0.3")
func AccountIDFromString(data string) (AccountID, error) {
	shard, realm, num, checksum, alias, aliasEvmAddress, err := _AccountIDFromString(data)
	if err != nil {
		return AccountID{}, err
	}

	if num == -1 {
		if alias != nil {
			return AccountID{
				Shard:           uint64(shard),
				Realm:           uint64(realm),
				Account:         0,
				AliasKey:        alias,
				AliasEvmAddress: nil,
				checksum:        checksum,
			}, nil
		}

		return AccountID{
			Shard:           uint64(shard),
			Realm:           uint64(realm),
			Account:         0,
			AliasKey:        nil,
			AliasEvmAddress: aliasEvmAddress,
			checksum:        checksum,
		}, nil
	}

	return AccountID{
		Shard:           uint64(shard),
		Realm:           uint64(realm),
		Account:         uint64(num),
		AliasKey:        nil,
		AliasEvmAddress: nil,
		checksum:        checksum,
	}, nil
}

func AccountIDFromEvmAddress(shard uint64, realm uint64, evmAddress string) (AccountID, error) {
	temp, err := hex.DecodeString(evmAddress)
	if err != nil {
		return AccountID{}, err
	}
	return AccountID{
		Shard:           shard,
		Realm:           realm,
		Account:         0,
		AliasEvmAddress: &temp,
		checksum:        nil,
	}, nil
}

// AccountIDFromSolidityAddress constructs an AccountID from a string
// representation of a _Solidity address
func AccountIDFromSolidityAddress(s string) (AccountID, error) {
	shard, realm, account, err := _IdFromSolidityAddress(s)
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:    shard,
		Realm:    realm,
		Account:  account,
		checksum: nil,
	}, nil
}

func (id *AccountID) ValidateChecksum(client *Client) error {
	if id.AliasKey != nil {
		return errors.New("Account ID contains alias key, unable to validate")
	}
	if !id._IsZero() && client != nil && client.network.ledgerID != nil {
		var tempChecksum _ParseAddressResult
		var err error
		tempChecksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
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
func (id *AccountID) Validate(client *Client) error {
	if id.AliasKey != nil {
		return errors.New("Account ID contains alias key, unable to validate")
	}
	if !id._IsZero() && client != nil && client.network.ledgerID != nil {
		tempChecksum, err := _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
		if err != nil {
			return err
		}
		err = _ChecksumVerify(tempChecksum.status)
		if err != nil {
			return err
		}

		if tempChecksum.correctChecksum != *id.checksum {
			return errNetworkMismatch
		}
	}

	return nil
}

// String returns the string representation of an AccountID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id AccountID) String() string {
	if id.AliasKey != nil {
		return fmt.Sprintf("%d.%d.%s", id.Shard, id.Realm, id.AliasKey.String())
	} else if id.AliasEvmAddress != nil {
		return fmt.Sprintf("%d.%d.%s", id.Shard, id.Realm, hex.EncodeToString(*id.AliasEvmAddress))
	}

	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account)
}

// ToStringWithChecksum returns the string representation of an AccountID in
// `Shard.Realm.Account-checksum` (for example "0.0.3-sdaf")
func (id AccountID) ToStringWithChecksum(client *Client) (string, error) {
	if id.AliasKey != nil {
		return "", errors.New("Account ID contains alias key, unable get checksum")
	}
	if client.GetNetworkName() == nil && client.GetLedgerID() == nil {
		return "", errNetworkNameMissing
	}
	var checksum _ParseAddressResult
	var err error
	if client.network.ledgerID != nil {
		checksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Account, checksum.correctChecksum), nil
}

// GetChecksum Retrieve just the checksum
func (id AccountID) GetChecksum() *string {
	return id.checksum
}

// ToSolidityAddress returns the string representation of the AccountID as a
// _Solidity address.
func (id AccountID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.Account)
}

func (id AccountID) _ToProtobuf() *services.AccountID {
	resultID := &services.AccountID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
	}
	if id.AliasKey != nil {
		data, _ := protobuf.Marshal(id.AliasKey._ToProtoKey())
		resultID.Account = &services.AccountID_Alias{
			Alias: data,
		}

		return resultID
	} else if id.AliasEvmAddress != nil {
		resultID.Account = &services.AccountID_Alias{
			Alias: *id.AliasEvmAddress,
		}

		return resultID
	}

	resultID.Account = &services.AccountID_AccountNum{
		AccountNum: int64(id.Account),
	}

	return resultID
}

// UnmarshalJSON implements the encoding.JSON interface.
func (id *AccountID) UnmarshalJSON(data []byte) error {
	accountID, err := AccountIDFromString(strings.Replace(string(data), "\"", "", 2))

	if err != nil {
		return err
	}

	*id = accountID

	return nil
}

func _AccountIDFromProtobuf(accountID *services.AccountID) *AccountID {
	if accountID == nil {
		return nil
	}
	resultAccountID := &AccountID{
		Shard: uint64(accountID.ShardNum),
		Realm: uint64(accountID.RealmNum),
	}

	switch t := accountID.Account.(type) {
	case *services.AccountID_Alias:
		pb := services.Key{}
		_ = protobuf.Unmarshal(t.Alias, &pb)
		initialKey, err := _KeyFromProtobuf(&pb)
		if err != nil && t.Alias != nil {
			resultAccountID.Account = 0
			resultAccountID.AliasEvmAddress = &t.Alias
			return resultAccountID
		}
		if evm, ok := pb.Key.(*services.Key_ECDSASecp256K1); ok && len(evm.ECDSASecp256K1) == 20 {
			resultAccountID.Account = 0
			resultAccountID.AliasEvmAddress = &evm.ECDSASecp256K1
			return resultAccountID
		}
		switch t2 := initialKey.(type) {
		case PublicKey:
			resultAccountID.Account = 0
			resultAccountID.AliasKey = &t2
			return resultAccountID
		default:
			return &AccountID{}
		}
	case *services.AccountID_AccountNum:
		resultAccountID.Account = uint64(t.AccountNum)
		resultAccountID.AliasKey = nil
		return resultAccountID
	default:
		return &AccountID{}
	}
}

func (id AccountID) IsZero() bool {
	return id._IsZero()
}

func (id AccountID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Account == 0 && id.AliasKey == nil
}

func (id AccountID) Equals(other AccountID) bool {
	return id._Equals(other)
}

func (id AccountID) _Equals(other AccountID) bool {
	initialAlias := ""
	otherAlias := ""
	if id.AliasKey != nil && other.AliasKey != nil {
		initialAlias = id.AliasKey.String()
		otherAlias = other.AliasKey.String()
	}

	return id.Shard == other.Shard && id.Realm == other.Realm && id.Account == other.Account && initialAlias == otherAlias
}

// ToBytes returns the wire-format encoding of AccountID
func (id AccountID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// AccountIDFromBytes converts wire-format encoding to Account ID
func AccountIDFromBytes(data []byte) (AccountID, error) {
	if data == nil {
		return AccountID{}, errByteArrayNull
	}
	pb := services.AccountID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AccountID{}, err
	}

	return *_AccountIDFromProtobuf(&pb), nil
}

func (id AccountID) Compare(given AccountID) int {
	if id.Shard > given.Shard { //nolint
		return 1
	} else if id.Shard < given.Shard {
		return -1
	}

	if id.Realm > given.Realm { //nolint
		return 1
	} else if id.Realm < given.Realm {
		return -1
	}

	if id.AliasKey != nil && given.AliasKey != nil {
		if id.AliasKey.String() > given.AliasKey.String() { //nolint
			return 1
		} else if id.AliasKey.String() < given.AliasKey.String() {
			return -1
		}
	}

	if id.AliasEvmAddress != nil && given.AliasEvmAddress != nil {
		originalEvmAddress := hex.EncodeToString(*id.AliasEvmAddress)
		givenEvmAddress := hex.EncodeToString(*id.AliasEvmAddress)
		if originalEvmAddress > givenEvmAddress { //nolint
			return 1
		} else if originalEvmAddress < givenEvmAddress {
			return -1
		}
	}

	if id.Account > given.Account { //nolint
		return 1
	} else if id.Account < given.Account {
		return -1
	} else {
		return 0
	}
}

func (accountIDs _AccountIDs) Len() int { //nolint
	return len(accountIDs.accountIDs)
}
func (accountIDs _AccountIDs) Swap(i, j int) { //nolint
	accountIDs.accountIDs[i], accountIDs.accountIDs[j] = accountIDs.accountIDs[j], accountIDs.accountIDs[i]
}

func (accountIDs _AccountIDs) Less(i, j int) bool { //nolint
	if accountIDs.accountIDs[i].Compare(accountIDs.accountIDs[j]) < 0 { //nolint
		return true
	}

	return false
}
