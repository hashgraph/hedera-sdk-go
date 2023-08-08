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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// NftID is the ID for a non-fungible token
type NftID struct {
	TokenID      TokenID
	SerialNumber int64
}

// NewNftID constructs a new NftID from a TokenID and a serial number
func NftIDFromString(s string) (NftID, error) {
	split := strings.Split(s, "@")
	if len(split) < 2 {
		panic(errors.New("wrong NftID format"))
	}
	shard, realm, num, checksum, err := _IdFromString(split[1])
	if err != nil {
		return NftID{}, err
	}

	serial, err := strconv.Atoi(split[0])
	if err != nil {
		return NftID{}, err
	}

	return NftID{
		TokenID: TokenID{
			Shard:    uint64(shard),
			Realm:    uint64(realm),
			Token:    uint64(num),
			checksum: checksum,
		},
		SerialNumber: int64(serial),
	}, nil
}

// Validate checks that the NftID is valid
func (id *NftID) Validate(client *Client) error {
	if !id._IsZero() && client != nil && client.network.ledgerID != nil {
		if err := id.TokenID.ValidateChecksum(client); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// String returns a string representation of the NftID
func (id NftID) String() string {
	return fmt.Sprintf("%d@%s", id.SerialNumber, id.TokenID.String())
}

// ToStringWithChecksum returns a string representation of the NftID with a checksum
func (id NftID) ToStringWithChecksum(client Client) (string, error) {
	token, err := id.TokenID.ToStringWithChecksum(client)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d@%s", id.SerialNumber, token), nil
}

func (id NftID) _ToProtobuf() *services.NftID {
	return &services.NftID{
		TokenID:      id.TokenID._ToProtobuf(),
		SerialNumber: id.SerialNumber,
	}
}

func _NftIDFromProtobuf(pb *services.NftID) NftID {
	if pb == nil {
		return NftID{}
	}

	tokenID := TokenID{}
	if pb.TokenID != nil {
		tokenID = *_TokenIDFromProtobuf(pb.TokenID)
	}

	return NftID{
		TokenID:      tokenID,
		SerialNumber: pb.SerialNumber,
	}
}

func (id NftID) _IsZero() bool {
	return id.TokenID._IsZero() && id.SerialNumber == 0
}

// ToBytes returns the byte representation of the NftID
func (id NftID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// NftIDFromBytes returns the NftID from a raw byte array
func NftIDFromBytes(data []byte) (NftID, error) {
	pb := services.NftID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return NftID{}, err
	}

	return _NftIDFromProtobuf(&pb), nil
}
