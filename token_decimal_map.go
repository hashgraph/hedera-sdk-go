package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
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

func _TokenDecimalMapFromProtobuf(pb []*services.TokenBalance) TokenDecimalMap {
	decimals := make(map[string]uint64)

	for _, tokenDecimal := range pb {
		decimals[_TokenIDFromProtobuf(tokenDecimal.TokenId).String()] = uint64(tokenDecimal.Decimals)
	}

	return TokenDecimalMap{decimals}
}

func (tokenDecimals TokenDecimalMap) _ToProtobuf() []*services.TokenBalance { // nolint
	decimals := make([]*services.TokenBalance, 0)

	for s, t := range tokenDecimals.decimals {
		token, err := TokenIDFromString(s)
		if err != nil {
			return []*services.TokenBalance{}
		}
		decimals = append(decimals, &services.TokenBalance{
			TokenId:  token._ToProtobuf(),
			Decimals: uint32(t),
		})
	}

	return decimals
}
