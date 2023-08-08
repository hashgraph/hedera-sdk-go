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
	"strings"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// ScheduleID is the ID for a Hedera account
type ScheduleID struct {
	Shard    uint64
	Realm    uint64
	Schedule uint64
	checksum *string
}

// ScheduleIDFromString constructs an ScheduleID from a string formatted as
// `Shard.Realm.Account` (for example "0.0.3")
func ScheduleIDFromString(data string) (ScheduleID, error) {
	shard, realm, num, checksum, err := _IdFromString(data)
	if err != nil {
		return ScheduleID{}, err
	}

	return ScheduleID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Schedule: uint64(num),
		checksum: checksum,
	}, nil
}

// ValidateChecksum validates the checksum of the account ID
func (id *ScheduleID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil {
		var tempChecksum _ParseAddressResult
		var err error
		tempChecksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule))
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
func (id *ScheduleID) Validate(client *Client) error {
	return id.ValidateChecksum(client)
}

// String returns the string representation of an ScheduleID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id ScheduleID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule)
}

// ToStringWithChecksum returns the string representation of an ScheduleID in
// `Shard.Realm.Account-checksum` (for example "0.0.3-laujm")
func (id ScheduleID) ToStringWithChecksum(client Client) (string, error) {
	if client.GetNetworkName() == nil && client.GetLedgerID() == nil {
		return "", errNetworkNameMissing
	}
	var checksum _ParseAddressResult
	var err error
	if client.network.ledgerID != nil {
		checksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule))
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Schedule, checksum.correctChecksum), nil
}

func (id ScheduleID) _ToProtobuf() *services.ScheduleID {
	return &services.ScheduleID{
		ShardNum:    int64(id.Shard),
		RealmNum:    int64(id.Realm),
		ScheduleNum: int64(id.Schedule),
	}
}

// UnmarshalJSON implements the encoding.JSON interface.
func (id *ScheduleID) UnmarshalJSON(data []byte) error {
	scheduleID, err := ScheduleIDFromString(strings.Replace(string(data), "\"", "", 2))

	if err != nil {
		return err
	}

	id.Shard = scheduleID.Shard
	id.Realm = scheduleID.Realm
	id.Schedule = scheduleID.Schedule
	id.checksum = scheduleID.checksum

	return nil
}

func _ScheduleIDFromProtobuf(scheduleID *services.ScheduleID) *ScheduleID {
	if scheduleID == nil {
		return nil
	}

	return &ScheduleID{
		Shard:    uint64(scheduleID.ShardNum),
		Realm:    uint64(scheduleID.RealmNum),
		Schedule: uint64(scheduleID.ScheduleNum),
	}
}

func (id ScheduleID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Schedule == 0
}

func (id ScheduleID) _Equals(other ScheduleID) bool { // nolint
	return id.Shard == other.Shard && id.Realm == other.Realm && id.Schedule == other.Schedule
}
