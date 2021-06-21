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
	Checksum *string
	Network  *NetworkName
}

// ScheduleIDFromString constructs an ScheduleID from a string formatted as
// `Shard.Realm.Account` (for example "0.0.3")
func ScheduleIDFromString(data string) (ScheduleID, error) {
	var checksum parseAddressResult
	var err error

	var networkNames = []NetworkName{
		Mainnet,
		Testnet,
		Previewnet,
	}

	var network NetworkName
	for _, name := range networkNames {
		checksum, err = checksumParseAddress(name.Network(), data)
		if err != nil {
			return ScheduleID{}, err
		}
		if checksum.status != 1 {
			network = name
			break
		}
	}

	err = checksumVerify(checksum.status)
	if err != nil {
		return ScheduleID{}, err
	}

	tempChecksum := checksum.correctChecksum

	return ScheduleID{
		Shard:    uint64(checksum.num1),
		Realm:    uint64(checksum.num2),
		Schedule: uint64(checksum.num3),
		Checksum: &tempChecksum,
		Network:  &network,
	}, nil
}

func ScheduleIDValidateNetworkOnIDs(id ScheduleID, other AccountID) error {
	if !id.isZero() && !other.isZero() && id.Network != nil && other.Network != nil && *id.Network != *other.Network {
		return errNetworkMismatch
	}

	return nil
}

// String returns the string representation of an ScheduleID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id ScheduleID) String() string {
	checksum, err := checksumParseAddress("", fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule))
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Schedule, checksum.correctChecksum)
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
