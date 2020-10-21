package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type TokenBalance struct{
	TokenId *TokenID
	Balance uint64
}

func tokenBalancesFromProtobuf(pb *proto.TokenBalance) TokenBalance {
	var tokenID = tokenIDFromProtobuf(pb.TokenId)
	return TokenBalance{
		TokenId: &tokenID,
		Balance: pb.Balance,
	}
}

func (balance *TokenBalance) toProtobuf() *proto.TokenBalance {
	return &proto.TokenBalance{
		TokenId: balance.TokenId.toProtobuf(),
		Balance: balance.Balance,
	}
}
