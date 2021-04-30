package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkVersionInfoQuery_Execute(t *testing.T) {
	client := newTestClient(t, false)

	_, err := NewNetworkVersionQuery().
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
}

func TestNetworkVersionInfoQueryCost_Execute(t *testing.T) {
	client := newTestClient(t, false)

	query := NewNetworkVersionQuery()

	cost, err := query.GetCost(client)
	assert.NoError(t, err)

	_, err = query.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)
}
