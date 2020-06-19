package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrictlyIncreasingTransactionIDValidStart(t *testing.T) {
	accountID := AccountID{Shard: 0, Realm: 0, Account: 2}
	lastTime := NewTransactionID(accountID).ValidStart

	for i := 0; i < 200; i++ {
		thisTime := NewTransactionID(accountID).ValidStart
		assert.True(t, thisTime.After(lastTime))
		lastTime = thisTime
	}
}
