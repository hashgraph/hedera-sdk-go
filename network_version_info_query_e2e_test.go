//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationNetworkVersionInfoQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewNetworkVersionQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)
}

func TestIntegrationNetworkVersionInfoQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	query := NewNetworkVersionQuery().SetNodeAccountIDs(env.NodeAccountIDs)

	cost, err := query.GetCost(env.Client)
	require.NoError(t, err)

	_, err = query.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

}
