//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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

	"github.com/stretchr/testify/require"
)

func TestUnitTokenIDFromString(t *testing.T) {
	tokID := TokenID{
		Shard: 1,
		Realm: 2,
		Token: 3,
	}

	gotTokID, err := TokenIDFromString(tokID.String())
	require.NoError(t, err)
	assert.Equal(t, tokID.Token, gotTokID.Token)
}

func TestUnitTokenIDChecksumFromString(t *testing.T) {
	id, err := TokenIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client := ClientForTestnet()
	id.ToStringWithChecksum(*client)
	sol := id.ToSolidityAddress()
	TokenIDFromSolidityAddress(sol)
	id.Validate(client)
	pb := id._ToProtobuf()
	_TokenIDFromProtobuf(pb)

	idByte := id.ToBytes()
	TokenIDFromBytes(idByte)

	id.Compare(TokenID{Token: 32})

	assert.Equal(t, id.Token, uint64(123))
}

func TestUnitTokenIDChecksumToString(t *testing.T) {
	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitTokenIDFromStringEVM(t *testing.T) {
	id, err := TokenIDFromString("0.0.434")
	require.NoError(t, err)

	require.Equal(t, "0.0.434", id.String())
}

func TestUnitTokenIDProtobuf(t *testing.T) {
	id, err := TokenIDFromString("0.0.434")
	require.NoError(t, err)

	pb := id._ToProtobuf()

	require.Equal(t, pb, &services.TokenID{
		ShardNum: 0,
		RealmNum: 0,
		TokenNum: 434,
	})

	pbFrom := _TokenIDFromProtobuf(pb)

	require.Equal(t, id, *pbFrom)
}
