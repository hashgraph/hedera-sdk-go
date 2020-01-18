package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type AccountID struct {
	Shard   uint64
	Realm   uint64
	Account uint64
}

func AccountIDFromString(s string) (AccountID, error) {
	shard, realm, num, err := idFromString(s)
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:   uint64(shard),
		Realm:   uint64(realm),
		Account: uint64(num),
	}, nil
}

func AccountIDFromSolidityAddress(s string) (AccountID, error) {
	shard, realm, account, err := idFromSolidityAddress(s)
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:   shard,
		Realm:   realm,
		Account: account,
	}, nil
}

func (id AccountID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account)
}

func (id AccountID) ToSolidityAddress() string {
	return idToSolidityAddress(id.Shard, id.Realm, id.Account)
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
