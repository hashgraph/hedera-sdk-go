//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
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
