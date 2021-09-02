package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenDecimalMap struct {
	decimals map[string]uint64
}

func (tokenDecimals *TokenDecimalMap) Get(tokenID TokenID) uint64 {
	return tokenDecimals.decimals[TokenID{
		Shard: tokenID.Shard,
		Realm: tokenID.Realm,
		Token: tokenID.Token,
	}.String()]
}

func tokenDecimalMapFromProtobuf(pb []*proto.TokenBalance) TokenDecimalMap {
	decimals := make(map[string]uint64)

	for _, tokenDecimal := range pb {
		decimals[tokenIDFromProtobuf(tokenDecimal.TokenId).String()] = uint64(tokenDecimal.Decimals)
	}

	return TokenDecimalMap{decimals}
}

func (tokenDecimals TokenDecimalMap) toProtobuf() []*proto.TokenBalance { // nolint
	decimals := make([]*proto.TokenBalance, 0)

	for s, t := range tokenDecimals.decimals {
		token, err := TokenIDFromString(s)
		if err != nil {
			return []*proto.TokenBalance{}
		}
		decimals = append(decimals, &proto.TokenBalance{
			TokenId:  token.toProtobuf(),
			Decimals: uint32(t),
		})
	}

	return decimals
}
