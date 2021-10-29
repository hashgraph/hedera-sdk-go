package hedera

import (
	"fmt"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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

func (id *ScheduleID) Validate(client *Client) error {
	if !id._IsZero() && client != nil && client.GetNetworkName() != nil {
		tempChecksum, err := _ChecksumParseAddress(client.GetNetworkName()._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule))
		if err != nil {
			return err
		}
		err = _ChecksumVerify(tempChecksum.status)
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

// String returns the string representation of an ScheduleID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id ScheduleID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule)
}

func (id ScheduleID) ToStringWithChecksum(client Client) (string, error) {
	if client.GetNetworkName() == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := _ChecksumParseAddress(client.GetNetworkName()._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Schedule, checksum.correctChecksum), nil
}

func (id ScheduleID) _ToProtobuf() *proto.ScheduleID {
	return &proto.ScheduleID{
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

func _ScheduleIDFromProtobuf(scheduleID *proto.ScheduleID) *ScheduleID {
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
