//go:build all || unit
// +build all unit

package hedera

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
