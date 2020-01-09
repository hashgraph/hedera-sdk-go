package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractID struct {
	Shard    uint64
	Realm    uint64
	Contract uint64
}

func ContractIDFromString(s string) (ContractID, error) {
	shard, realm, num, err := idFromString(s)
	if err != nil {
		return ContractID{}, err
	}

	return ContractID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Contract: uint64(num),
	}, nil
}

func ContractIDFromSolidityAddress(s string) (ContractID, error) {
	shard, realm, contract, err := idFromSolidityAddress(s)
	if err != nil {
		return ContractID{}, err
	}

	return ContractID{
		Shard:    shard,
		Realm:    realm,
		Contract: contract,
	}, nil
}

func (id ContractID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract)
}

func (id ContractID) ToSolidityAddress() string {
	return idToSolidityAddress(id.Shard, id.Realm, id.Contract)
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
