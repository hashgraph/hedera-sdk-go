package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type TokenBalanceMap struct {
	balances map[string]uint64
}

func (tokenBalances *TokenBalanceMap) Get(tokenID TokenID) uint64 {
	return tokenBalances.balances[TokenID{
		Shard: tokenID.Shard,
		Realm: tokenID.Realm,
		Token: tokenID.Token,
	}.String()]
}

func tokenBalanceMapFromProtobuf(pb []*services.TokenBalance, _ *NetworkName) TokenBalanceMap {
	balances := make(map[string]uint64, 0)

	for _, tokenBalance := range pb {
		balances[tokenIDFromProtobuf(tokenBalance.TokenId, nil).String()] = tokenBalance.Balance
	}

	return TokenBalanceMap{balances}
}
