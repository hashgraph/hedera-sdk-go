package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkVersionInfoQuery_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewNetworkVersionQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.NoError(t, err)
}

func TestNetworkVersionInfoQueryCost_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	query := NewNetworkVersionQuery().SetNodeAccountIDs(env.NodeAccountIDs)

	cost, err := query.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = query.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)
}
