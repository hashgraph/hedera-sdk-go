//+build all unit

package hedera

import (
	"testing"

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
