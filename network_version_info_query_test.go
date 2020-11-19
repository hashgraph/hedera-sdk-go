package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkVersionInfoQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	_, err := NewNetworkVersionQuery().Execute(client)
	assert.Nil(t, err)
}
