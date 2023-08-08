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

func TestUnitTokenDecimalMapGet(t *testing.T) {
	t.Parallel()

	tokenDecimals := TokenDecimalMap{
		decimals: map[string]uint64{
			"0.0.123": 9,
			"0.0.124": 10,
		},
	}

	assert.Equal(t, uint64(9), tokenDecimals.Get(TokenID{Shard: 0, Realm: 0, Token: 123}))
	assert.Equal(t, uint64(10), tokenDecimals.Get(TokenID{Shard: 0, Realm: 0, Token: 124}))
}

func TestUnitTokenDecimalMapToProtobuf(t *testing.T) {
	t.Parallel()

	tokenDecimals := TokenDecimalMap{
		decimals: map[string]uint64{
			"0.0.123": 9,
			"0.0.124": 10,
		},
	}

	decimals := tokenDecimals._ToProtobuf()

	assert.Equal(t, 2, len(decimals))

	// The order of the decimals is not guaranteed
	for _, dec := range decimals {
		switch dec.TokenId.TokenNum {
		case 123:
			assert.Equal(t, uint32(9), dec.Decimals)
		case 124:
			assert.Equal(t, uint32(10), dec.Decimals)
		default:
			t.Errorf("Unexpected TokenID: %v", dec.TokenId.String())
		}
	}
}

func TestUnitTokenDecimalMapFromProtobuf(t *testing.T) {
	t.Parallel()

	decimals := make([]*services.TokenBalance, 0)
	decimals = append(decimals, &services.TokenBalance{
		TokenId:  &services.TokenID{ShardNum: 0, RealmNum: 0, TokenNum: 123},
		Decimals: uint32(9),
	})
	decimals = append(decimals, &services.TokenBalance{
		TokenId:  &services.TokenID{ShardNum: 0, RealmNum: 0, TokenNum: 124},
		Decimals: uint32(10),
	})

	tokenDecimals := _TokenDecimalMapFromProtobuf(decimals)

	assert.Equal(t, uint64(9), tokenDecimals.Get(TokenID{Shard: 0, Realm: 0, Token: 123}))
	assert.Equal(t, uint64(10), tokenDecimals.Get(TokenID{Shard: 0, Realm: 0, Token: 124}))
}

func TestUnitTokenDecimalMapFromProtobufEmpty(t *testing.T) {
	t.Parallel()

	tokenDecimals := TokenDecimalMap{
		decimals: map[string]uint64{
			"0.123":   9, // invalid token
			"0.0.124": 10,
		},
	}
	pb := tokenDecimals._ToProtobuf()
	assert.Equal(t, []*services.TokenBalance{}, pb)
}
