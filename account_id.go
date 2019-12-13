package hedera

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type AccountID struct {
	Shard   uint64
	Realm   uint64
	Account uint64
}

func AccountIDFromString(s string) (AccountID, error) {
	values := strings.SplitN(s, ".", 3)
	if len(values) != 3 {
		// Was not three values separated by periods
		return AccountID{}, fmt.Errorf("expected {shard}.{realm}.{num}")
	}

	shard, err := strconv.Atoi(values[0])
	if err != nil {
		return AccountID{}, err
	}

	realm, err := strconv.Atoi(values[1])
	if err != nil {
		return AccountID{}, err
	}

	num, err := strconv.Atoi(values[2])
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:   uint64(shard),
		Realm:   uint64(realm),
		Account: uint64(num),
	}, nil
}

func (id AccountID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account)
}

func (id AccountID) toProto() *proto.AccountID {
	return &proto.AccountID{
		ShardNum:   int64(id.Shard),
		RealmNum:   int64(id.Realm),
		AccountNum: int64(id.Account),
	}
}

func accountIDFromProto(pb *proto.AccountID) AccountID {
	return AccountID{
		Shard:   uint64(pb.ShardNum),
		Realm:   uint64(pb.RealmNum),
		Account: uint64(pb.AccountNum),
	}
}
