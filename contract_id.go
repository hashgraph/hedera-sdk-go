package hedera

import (
	"encoding/binary"
	"encoding/hex"
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
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return ContractID{}, err
	}

	if len(bytes) != 20 {
		return ContractID{}, fmt.Errorf("Solidity address must be 20 bytes")
	}

	shard := uint64(binary.BigEndian.Uint32(bytes[0:4]))
	realm := binary.BigEndian.Uint64(bytes[4:12])
	contract := binary.BigEndian.Uint64(bytes[12:20])

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
	bytes := make([]byte, 20)
	binary.BigEndian.PutUint32(bytes[0:4], uint32(id.Shard))
	binary.BigEndian.PutUint64(bytes[4:12], id.Realm)
	binary.BigEndian.PutUint64(bytes[12:20], id.Contract)
	return hex.EncodeToString(bytes)
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
