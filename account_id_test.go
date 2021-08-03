package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnitAccountIDChecksumFromString(t *testing.T) {
	id, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	assert.Equal(t, id.Account, uint64(123))
}

func TestUnitAccountIDChecksumToString(t *testing.T) {
	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, id.String(), "50.150.520")
}
