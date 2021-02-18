package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenIDFromString(t *testing.T) {
	tokID := TokenID{
		Shard: 1,
		Realm: 2,
		Token: 3,
	}

	gotTokID, err := TokenIDFromString(tokID.String())
	assert.NoError(t, err)
	assert.Equal(t, tokID, gotTokID)
}
