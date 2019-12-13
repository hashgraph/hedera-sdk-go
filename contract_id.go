package hedera

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractID struct {
	Shard    uint64
	Realm    uint64
	Contract uint64
}

func ContractIDFromString(s string) (ContractID, error) {
	values := strings.SplitN(s, ".", 3)
	if len(values) != 3 {
		// Was not three values separated by periods
		return ContractID{}, fmt.Errorf("expected {shard}.{realm}.{num}")
	}

	shard, err := strconv.Atoi(values[0])
	if err != nil {
		return ContractID{}, err
	}

	realm, err := strconv.Atoi(values[1])
	if err != nil {
		return ContractID{}, err
	}

	num, err := strconv.Atoi(values[2])
	if err != nil {
		return ContractID{}, err
	}

	return ContractID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Contract: uint64(num),
	}, nil
}

func (id ContractID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract)
}

func (id ContractID) toProto() *proto.ContractID {
	return &proto.ContractID{
		ShardNum:    int64(id.Shard),
		RealmNum:    int64(id.Realm),
		ContractNum: int64(id.Contract),
	}
}

func contractIDFromProto(pb *proto.ContractID) ContractID {
	return ContractID{
		Shard:    uint64(pb.ShardNum),
		Realm:    uint64(pb.RealmNum),
		Contract: uint64(pb.ContractNum),
	}
}
