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
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitContractIDChecksumFromString(t *testing.T) {
	id, err := ContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client := ClientForTestnet()
	id.ToStringWithChecksum(*client)
	sol := id.ToSolidityAddress()
	ContractIDFromSolidityAddress(sol)
	id.Validate(client)
	evmID, err := ContractIDFromEvmAddress(0, 0, "ace082947b949651c703ff0f02bc1541")
	require.NoError(t, err)
	pb := evmID._ToProtobuf()
	_ContractIDFromProtobuf(pb)

	idByte := id.ToBytes()
	ContractIDFromBytes(idByte)

	id._ToProtoKey()

	assert.Equal(t, id.Contract, uint64(123))
}

func TestUnitContractIDChecksumToString(t *testing.T) {
	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitContractIDFromStringEVM(t *testing.T) {
	id, err := ContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, "0.0.0011223344556677889900112233445577889900", id.String())
}

func TestUnitContractIDProtobuf(t *testing.T) {
	id, err := ContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	pb := id._ToProtobuf()

	decoded, err := hex.DecodeString("0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_EvmAddress{EvmAddress: decoded},
	})

	pbFrom := _ContractIDFromProtobuf(pb)

	require.Equal(t, id, *pbFrom)
}

func TestUnitContractIDEvm(t *testing.T) {
	hexString, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	id, err := ContractIDFromString(fmt.Sprintf("0.0.%s", hexString.PublicKey().String()))
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(id.EvmAddress), hexString.PublicKey().String())

	pb := id._ToProtobuf()
	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_EvmAddress{EvmAddress: id.EvmAddress},
	})

	id, err = ContractIDFromString("0.0.123")
	require.NoError(t, err)
	require.Equal(t, id.Contract, uint64(123))
	require.Nil(t, id.EvmAddress)

	pb = id._ToProtobuf()
	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_ContractNum{ContractNum: 123},
	})
}
