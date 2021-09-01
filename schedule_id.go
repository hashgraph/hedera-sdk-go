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
	shard, realm, num, checksum, err := idFromString(data)
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
	if !id.isZero() && client != nil && client.network.networkName != nil {
		tempChecksum, err := checksumParseAddress(client.network.networkName.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule))
		if err != nil {
			return err
		}
		err = checksumVerify(tempChecksum.status)
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

func (id *ScheduleID) setNetworkWithClient(client *Client) {
	if client.network.networkName != nil {
		id.setNetwork(*client.network.networkName)
	}
}

func (id *ScheduleID) setNetwork(name NetworkName) {
	checksum := checkChecksum(name.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule))
	id.checksum = &checksum
}

// String returns the string representation of an ScheduleID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id ScheduleID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule)
}

func (id ScheduleID) ToStringWithChecksum(client Client) (string, error) {
	if client.network.networkName == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := checksumParseAddress(client.network.networkName.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Schedule))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Schedule, checksum.correctChecksum), nil
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

func scheduleIDFromProtobuf(scheduleID *proto.ScheduleID) ScheduleID {
	if scheduleID == nil {
		return ScheduleID{}
	}

	id := ScheduleID{
		Shard:    uint64(scheduleID.ShardNum),
		Realm:    uint64(scheduleID.RealmNum),
		Schedule: uint64(scheduleID.ScheduleNum),
	}

	return id
}

func (id ScheduleID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Schedule == 0
}

func (id ScheduleID) equals(other ScheduleID) bool {
	return id.Shard == other.Shard && id.Realm == other.Realm && id.Schedule == other.Schedule
}
