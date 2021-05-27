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
}

// ScheduleIDFromString constructs an ScheduleID from a string formatted as
// `Shard.Realm.Account` (for example "0.0.3")
func ScheduleIDFromString(s string) (ScheduleID, error) {
	shard, realm, num, err := idFromString(s)
	if err != nil {
		return ScheduleID{}, err
	}

	return ScheduleID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Schedule: uint64(num),
	}, nil
}

// String returns the string representation of an ScheduleID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id ScheduleID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule)
}

func (id ScheduleID) toProtobuf() *proto.ScheduleID {
	return &proto.ScheduleID{
		ShardNum:    int64(id.Shard),
		RealmNum:    int64(id.Realm),
		ScheduleNum: int64(id.Schedule),
	}
}

// UnmarshalJSON implements the encoding.JSON interface.
func (id *ScheduleID) UnmarshalJSON(data []byte) error {
	ScheduleID, err := ScheduleIDFromString(strings.Replace(string(data), "\"", "", 2))

	if err != nil {
		return err
	}

	id = &ScheduleID

	return nil
}

func scheduleIDFromProtobuf(pb *proto.ScheduleID) ScheduleID {
	if pb == nil {
		return ScheduleID{}
	}
	return ScheduleID{
		Shard:    uint64(pb.ShardNum),
		Realm:    uint64(pb.RealmNum),
		Schedule: uint64(pb.ScheduleNum),
	}
}

func (id ScheduleID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Schedule == 0
}

func (id ScheduleID) equals(other ScheduleID) bool {
	return id.Shard == other.Shard && id.Realm == other.Realm && id.Schedule == other.Schedule
}
