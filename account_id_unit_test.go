//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountIDChecksumFromString(t *testing.T) {
	id, err := AccountIDFromString("0.0.123-rmkyk")
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
