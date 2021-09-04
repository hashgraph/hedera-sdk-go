package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

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

func _TokenBalanceMapFromProtobuf(pb []*proto.TokenBalance) TokenBalanceMap {
	balances := make(map[string]uint64)

	for _, tokenBalance := range pb {
		balances[_TokenIDFromProtobuf(tokenBalance.TokenId).String()] = tokenBalance.Balance
	}

	return TokenBalanceMap{balances}
}

func (tokenBalances *TokenBalanceMap) _ToProtobuf() []*proto.TokenBalance { // nolint
	decimals := make([]*proto.TokenBalance, 0)

	for s, t := range tokenBalances.balances {
		token, err := TokenIDFromString(s)
		if err != nil {
			return []*proto.TokenBalance{}
		}
		decimals = append(decimals, &proto.TokenBalance{
			TokenId: token._ToProtobuf(),
			Balance: t,
		})
	}

	return decimals
}
