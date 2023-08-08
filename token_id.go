package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// TokenID is the ID for a Hedera token
type TokenID struct {
	Shard    uint64
	Realm    uint64
	Token    uint64
	checksum *string
}

type _TokenIDs struct {
	tokenIDs []TokenID
}

func _TokenIDFromProtobuf(tokenID *services.TokenID) *TokenID {
	if tokenID == nil {
		return nil
	}

	return &TokenID{
		Shard: uint64(tokenID.ShardNum),
		Realm: uint64(tokenID.RealmNum),
		Token: uint64(tokenID.TokenNum),
	}
}

func (id *TokenID) _ToProtobuf() *services.TokenID {
	return &services.TokenID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TokenNum: int64(id.Token),
	}
}

// String returns a string representation of the TokenID formatted as `Shard.Realm.TokenID` (for example "0.0.3")
func (id TokenID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token)
}

// ToStringWithChecksum returns a string representation of the TokenID formatted as `Shard.Realm.TokenID-Checksum` (for example "0.0.3-abcd")
func (id TokenID) ToStringWithChecksum(client Client) (string, error) {
	if client.GetNetworkName() == nil && client.GetLedgerID() == nil {
		return "", errNetworkNameMissing
	}
	var checksum _ParseAddressResult
	var err error
	if client.network.ledgerID != nil {
		checksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Token, checksum.correctChecksum), nil
}

// ToBytes returns a byte array representation of the TokenID
func (id TokenID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TokenIDFromBytes returns a TokenID from a byte array
func TokenIDFromBytes(data []byte) (TokenID, error) {
	if data == nil {
		return TokenID{}, errByteArrayNull
	}
	pb := services.TokenID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenID{}, err
	}

	return *_TokenIDFromProtobuf(&pb), nil
}

// NftID constructs an NftID from a TokenID and a serial number
func (id *TokenID) Nft(serial int64) NftID {
	return NftID{
		TokenID:      *id,
		SerialNumber: serial,
	}
}

// TokenIDFromString constructs an TokenID from a string formatted as
// `Shard.Realm.TokenID` (for example "0.0.3")
func TokenIDFromString(data string) (TokenID, error) {
	shard, realm, num, checksum, err := _IdFromString(data)
	if err != nil {
		return TokenID{}, err
	}

	return TokenID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Token:    uint64(num),
		checksum: checksum,
	}, nil
}

// Verify that the client has a valid checksum.
func (id *TokenID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil {
		var tempChecksum _ParseAddressResult
		var err error
		tempChecksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
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
func (id *TokenID) Validate(client *Client) error {
	return id.ValidateChecksum(client)
}

// TokenIDFromSolidityAddress constructs a TokenID from a string
// representation of a _Solidity address
func TokenIDFromSolidityAddress(s string) (TokenID, error) {
	shard, realm, token, err := _IdFromSolidityAddress(s)
	if err != nil {
		return TokenID{}, err
	}

	return TokenID{
		Shard:    shard,
		Realm:    realm,
		Token:    token,
		checksum: nil,
	}, nil
}

// ToSolidityAddress returns the string representation of the TokenID as a
// _Solidity address.
func (id TokenID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.Token)
}

func (id TokenID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Token == 0
}

// Compare compares two TokenIDs
func (id TokenID) Compare(given TokenID) int {
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

	if id.Token > given.Token { //nolint
		return 1
	} else if id.Token < given.Token {
		return -1
	} else { //nolint
		return 0
	}
}

func (tokenIDs _TokenIDs) Len() int {
	return len(tokenIDs.tokenIDs)
}
func (tokenIDs _TokenIDs) Swap(i, j int) {
	tokenIDs.tokenIDs[i], tokenIDs.tokenIDs[j] = tokenIDs.tokenIDs[j], tokenIDs.tokenIDs[i]
}

func (tokenIDs _TokenIDs) Less(i, j int) bool {
	if tokenIDs.tokenIDs[i].Compare(tokenIDs.tokenIDs[j]) < 0 { //nolint
		return true
	}

	return false
}
