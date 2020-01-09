package hedera

import (
	"encoding/binary"
	"encoding/hex"
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
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return AccountID{}, err
	}

	if len(bytes) != 20 {
		return AccountID{}, fmt.Errorf("Solidity address must be 20 bytes")
	}

	shard := uint64(binary.BigEndian.Uint32(bytes[0:4]))
	realm := binary.BigEndian.Uint64(bytes[4:12])
	account := binary.BigEndian.Uint64(bytes[12:20])

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
	bytes := make([]byte, 20)
	binary.BigEndian.PutUint32(bytes[0:4], uint32(id.Shard))
	binary.BigEndian.PutUint64(bytes[4:12], id.Realm)
	binary.BigEndian.PutUint64(bytes[12:20], id.Account)
	return hex.EncodeToString(bytes)
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
