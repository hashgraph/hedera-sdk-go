//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	id.ToStringWithChecksum(client)
	id.GetChecksum()
	sol := id.ToSolidityAddress()
	AccountIDFromSolidityAddress(sol)
	err = id.Validate(client)
	require.Error(t, err)
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
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitAccountIDFromStringAlias(t *testing.T) {
	t.Parallel()

	key, err := GeneratePrivateKey()
	require.NoError(t, err)
	id, err := AccountIDFromString("0.0." + key.PublicKey().String())
	require.NoError(t, err)
	id2 := key.ToAccountID(0, 0)

	assert.Equal(t, id.String(), id2.String())
}

func TestUnitChecksum(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	id, err := AccountIDFromString("0.0.0011223344556677889900112233445566778899")
	require.NoError(t, err)

	require.Equal(t, id.String(), "0.0.0011223344556677889900112233445566778899")
}

func TestUnitAccountIDPopulateFailForWrongMirrorHost(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateAccount(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateFailWithNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateAccount(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailForWrongMirrorHost(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailWithNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailWithNoMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}
