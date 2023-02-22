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

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountIDChecksumFromString(t *testing.T) {
	id, err := AccountIDFromString("0.0.123-rmkyk")

	client := ClientForTestnet()
	id.ToStringWithChecksum(client)
	id.GetChecksum()
	sol := id.ToSolidityAddress()
	AccountIDFromSolidityAddress(sol)
	id.Validate(client)
	evmID, err := AccountIDFromEvmAddress(0, 0, "ace082947b949651c703ff0f02bc1541")
	require.NoError(t, err)
	pb := evmID._ToProtobuf()
	_AccountIDFromProtobuf(pb)

	idByte := id.ToBytes()
	AccountIDFromBytes(idByte)

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	alias := key.ToAccountID(0, 0)
	pb = alias._ToProtobuf()
	_AccountIDFromProtobuf(pb)

	require.NoError(t, err)
	assert.Equal(t, id.Account, uint64(123))
}

func TestUnitAccountIDChecksumToString(t *testing.T) {
	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitAccountIDFromStringAlias(t *testing.T) {
	key, err := GeneratePrivateKey()
	require.NoError(t, err)
	id, err := AccountIDFromString("0.0." + key.PublicKey().String())
	require.NoError(t, err)
	id2 := key.ToAccountID(0, 0)

	assert.Equal(t, id.String(), id2.String())
}

func TestUnitChecksum(t *testing.T) {
	id, err := LedgerIDFromString("01")
	require.NoError(t, err)
	ad1, err := _ChecksumParseAddress(id, "0.0.3")
	require.NoError(t, err)
	id, err = LedgerIDFromString("10")
	require.NoError(t, err)
	ad2, err := _ChecksumParseAddress(id, "0.0.3")
	require.NoError(t, err)

	require.NotEqual(t, ad1.correctChecksum, ad2.correctChecksum)
}

func TestUnitAccountIDEvm(t *testing.T) {
	id, err := AccountIDFromString("0.0.0011223344556677889900112233445566778899")
	require.NoError(t, err)

	require.Equal(t, id.String(), "0.0.0011223344556677889900112233445566778899")
}
func TestUnitAccountIDEvmAddressOnly0x(t *testing.T) {
	id, err := AccountIDFromString("0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)

	require.Equal(t, id.String(), "0.0.1234567890abcdef1234567890abcdef12345678")
}

func TestUnitAccountIDEvmAddressOnlyWithout0x(t *testing.T) {
	id, err := AccountIDFromString("1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)

	require.Equal(t, id.String(), "0.0.1234567890abcdef1234567890abcdef12345678")
}
func TestUnitAccountIDEvmPublicAddress0x(t *testing.T) {
	id, err := AccountIDFromEvmPublicAddress("0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)

	require.Equal(t, id.String(), "0.0.1234567890abcdef1234567890abcdef12345678")
}
func TestUnitAccountIDEvmPublicAddressWithout0x(t *testing.T) {
	id, err := AccountIDFromEvmPublicAddress("0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)

	require.Equal(t, id.String(), "0.0.1234567890abcdef1234567890abcdef12345678")
}
