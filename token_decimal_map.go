package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

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

func tokenDecimalMapFromProtobuf(pb []*services.TokenBalance, _ *NetworkName) TokenDecimalMap {
	decimals := make(map[string]uint64, 0)

	for _, tokenDecimal := range pb {
		decimals[tokenIDFromProtobuf(tokenDecimal.TokenId, nil).String()] = uint64(tokenDecimal.Decimals)
	}

	return TokenDecimalMap{decimals}
}
