package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

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
