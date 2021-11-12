package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

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

func _TokenBalanceMapFromProtobuf(pb []*services.TokenBalance) TokenBalanceMap {
	balances := make(map[string]uint64)

	for _, tokenBalance := range pb {
		balances[_TokenIDFromProtobuf(tokenBalance.TokenId).String()] = tokenBalance.Balance
	}

	return TokenBalanceMap{balances}
}

func (tokenBalances *TokenBalanceMap) _ToProtobuf() []*services.TokenBalance { // nolint
	decimals := make([]*services.TokenBalance, 0)

	for s, t := range tokenBalances.balances {
		token, err := TokenIDFromString(s)
		if err != nil {
			return []*services.TokenBalance{}
		}
		decimals = append(decimals, &services.TokenBalance{
			TokenId: token._ToProtobuf(),
			Balance: t,
		})
	}

	return decimals
}
