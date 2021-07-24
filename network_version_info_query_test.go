package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationNetworkVersionInfoQueryCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewNetworkVersionQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationNetworkVersionInfoQueryGetCost(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	query := NewNetworkVersionQuery().SetNodeAccountIDs(env.NodeAccountIDs)

	cost, err := query.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = query.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
