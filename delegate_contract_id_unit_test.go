//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitDelegatableContractIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	err = id.ValidateChecksum(client)
	require.Error(t, err)
	require.Equal(t, id.Contract, uint64(123))
	strChecksum, err := id.ToStringWithChecksum(*client)
	require.NoError(t, err)
	// different checksum because of different network
	require.Equal(t, strChecksum, "0.0.123-esxsf")
}

func TestUnitDelegatableContractIDChecksumToString(t *testing.T) {
	t.Parallel()

	id := DelegatableContractID{
		Shard:    50,
		Realm:    150,
		Contract: 520,
	}
	require.Equal(t, "50.150.520", id.String())
}

func TestUnitDelegatableContractIDFromStringEVM(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, "0.0.0011223344556677889900112233445577889900", id.String())
}

func TestUnitDelegatableContractIDProtobuf(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	pb := id._ToProtobuf()

	decoded, err := hex.DecodeString("0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_EvmAddress{EvmAddress: decoded},
	})

	pbFrom := _DelegatableContractIDFromProtobuf(pb)

	require.Equal(t, id, *pbFrom)
}

func TestUnitDelegatableContractIDEvm(t *testing.T) {
	t.Parallel()

	hexString, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	id, err := DelegatableContractIDFromString(fmt.Sprintf("0.0.%s", hexString.PublicKey().String()))
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(id.EvmAddress), hexString.PublicKey().String())

	pb := id._ToProtobuf()
	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_EvmAddress{EvmAddress: id.EvmAddress},
	})

	id, err = DelegatableContractIDFromString("0.0.123")
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

func TestUnitDelegatableContractIDToFromBytes(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123")
	require.NoError(t, err)
	require.Equal(t, id.Contract, uint64(123))
	require.Nil(t, id.EvmAddress)

	idBytes := id.ToBytes()
	idFromBytes, err := DelegatableContractIDFromBytes(idBytes)
	require.NoError(t, err)
	require.Equal(t, id, idFromBytes)
}

func TestUnitDelegatableContractIDFromEvmAddress(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromEvmAddress(0, 0, "0011223344556677889900112233445566778899")
	require.NoError(t, err)
	require.Equal(t, id.Contract, uint64(0))
	require.Equal(t, id.EvmAddress, []byte{0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99})
}

func TestUnitDelegatableContractIDFromSolidityAddress(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	sol := id.ToSolidityAddress()
	idFromSolidity, err := DelegatableContractIDFromSolidityAddress(sol)
	require.NoError(t, err)
	require.Equal(t, idFromSolidity.Contract, uint64(123))
}

func TestUnitDelegatableContractIDToProtoKey(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	pb := id._ToProtoKey()
	require.Equal(t, pb.GetContractID().GetContractNum(), int64(123))
}
