//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationNetworkVersionInfoQueryCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewNetworkVersionQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationNetworkVersionInfoQueryGetCost(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	query := NewNetworkVersionQuery().SetNodeAccountIDs(env.NodeAccountIDs)

	cost, err := query.GetCost(env.Client)
	require.NoError(t, err)

	_, err = query.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
