//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitContractIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := ContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	sol := id.ToSolidityAddress()
	ContractIDFromSolidityAddress(sol)
	err = id.Validate(client)
	require.Error(t, err)
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
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitContractIDFromStringEVM(t *testing.T) {
	t.Parallel()

	id, err := ContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, "0.0.0011223344556677889900112233445577889900", id.String())
}

func TestUnitContractIDProtobuf(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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

func TestUnitContractIDPopulateFailForWrongMirrorHost(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := ContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateContract(client)
	require.Error(t, err)
}

func TestUnitContractIDPopulateFailWithNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := ContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateContract(client)
	require.Error(t, err)
}
