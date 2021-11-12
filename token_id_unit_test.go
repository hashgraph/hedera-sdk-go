//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenIDFromString(t *testing.T) {
	tokID := TokenID{
		Shard: 1,
		Realm: 2,
		Token: 3,
	}

	gotTokID, err := TokenIDFromString(tokID.String())
	assert.NoError(t, err)
	assert.Equal(t, tokID.Token, gotTokID.Token)
}
