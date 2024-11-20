package hiero

// SPDX-License-Identifier: Apache-2.0
import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type TokenBalanceMap struct {
	balances map[string]uint64
}

// Get returns the balance of the given tokenID
func (tokenBalances *TokenBalanceMap) Get(tokenID TokenID) uint64 {
	return tokenBalances.balances[tokenID.String()]
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
