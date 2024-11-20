//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenIDFromString(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	id, err := TokenIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
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
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitTokenIDFromStringEVM(t *testing.T) {
	t.Parallel()

	id, err := TokenIDFromString("0.0.434")
	require.NoError(t, err)

	require.Equal(t, "0.0.434", id.String())
}

func TestUnitTokenIDProtobuf(t *testing.T) {
	t.Parallel()

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
