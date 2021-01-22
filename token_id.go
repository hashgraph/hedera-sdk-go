package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenID struct {
	Shard uint64
	Realm uint64
	Token uint64
}

func tokenIDFromProtobuf(tokenID *proto.TokenID) TokenID {
	return TokenID{
		Shard: uint64(tokenID.ShardNum),
		Realm: uint64(tokenID.RealmNum),
		Token: uint64(tokenID.TokenNum),
	}
}

func (id *TokenID) toProtobuf() *proto.TokenID {
	return &proto.TokenID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TokenNum: int64(id.Token),
	}
}

func (id TokenID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token)
}

func (id TokenID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenIDFromBytes(data []byte) (TokenID, error) {
	pb := proto.TokenID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenID{}, err
	}

	return tokenIDFromProtobuf(&pb), nil
}
