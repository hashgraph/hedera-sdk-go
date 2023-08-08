//go:build all || unit
// +build all unit

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
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/assert"
)

func TestUnitTokenBalanceMapGet(t *testing.T) {
	t.Parallel()

	var tokenBalances TokenBalanceMap
	tokenBalances.balances = make(map[string]uint64)
	tokenBalances.balances["0.0.123"] = 100

	assert.Equal(t, uint64(100), tokenBalances.Get(TokenID{Shard: 0, Realm: 0, Token: 123}))
}

func TestUnitTokenBalanceMapProtobuf(t *testing.T) {
	t.Parallel()

	var tokenBalances TokenBalanceMap
	tokenBalances.balances = make(map[string]uint64)
	tokenBalances.balances["0.0.123"] = 100

	pb := tokenBalances._ToProtobuf()
	tokenBalances2 := _TokenBalanceMapFromProtobuf(pb)

	assert.Equal(t, tokenBalances.balances, tokenBalances2.balances)
}

func TestUnitTokenBalanceMapEmpty(t *testing.T) {
	t.Parallel()

	var tokenBalances TokenBalanceMap
	tokenBalances.balances = make(map[string]uint64)

	// Breaks token, err := TokenIDFromString(s)
	tokenBalances.balances["0.123"] = 100

	pb := tokenBalances._ToProtobuf()

	// test that we get an empty array back
	assert.Equal(t, []*services.TokenBalance{}, pb)
}
