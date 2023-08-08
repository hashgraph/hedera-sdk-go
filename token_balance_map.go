package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import (
	"github.com/hashgraph/hedera-protobufs-go/services"
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
