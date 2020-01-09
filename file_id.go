package hedera

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FileID struct {
	Shard uint64
	Realm uint64
	File  uint64
}

func FileIDFromString(s string) (FileID, error) {
	shard, realm, num, err := idFromString(s)
	if err != nil {
		return FileID{}, err
	}

	return FileID{
		Shard: uint64(shard),
		Realm: uint64(realm),
		File:  uint64(num),
	}, nil
}

func FileIDFromSolidityAddress(s string) (FileID, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return FileID{}, err
	}

	if len(bytes) != 20 {
		return FileID{}, fmt.Errorf("Solidity address must be 20 bytes")
	}

	shard := uint64(binary.BigEndian.Uint32(bytes[0:4]))
	realm := binary.BigEndian.Uint64(bytes[4:12])
	file := binary.BigEndian.Uint64(bytes[12:20])

	return FileID{
		Shard: shard,
		Realm: realm,
		File:  file,
	}, nil
}

func (id FileID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.File)
}

func (id FileID) ToSolidityAddress() string {
	bytes := make([]byte, 20)
	binary.BigEndian.PutUint32(bytes[0:4], uint32(id.Shard))
	binary.BigEndian.PutUint64(bytes[4:12], id.Realm)
	binary.BigEndian.PutUint64(bytes[12:20], id.File)
	return hex.EncodeToString(bytes)
}

func (id FileID) toProto() *proto.FileID {
	return &proto.FileID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		FileNum:  int64(id.File),
	}
}

func fileIDFromProto(pb *proto.FileID) FileID {
	return FileID{
		Shard: uint64(pb.ShardNum),
		Realm: uint64(pb.RealmNum),
		File:  uint64(pb.FileNum),
	}
}
