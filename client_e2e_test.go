//+build all e2e

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationClientPingAllGoodNetwork(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	env.Client.SetMaxNodeAttempts(1)
	env.Client.PingAll()

	net := env.Client.GetNetwork()

	keys := make([]string, len(net))
	val := make([]AccountID, len(net))
	i := 0
	for st, n := range net {
		keys[i] = st
		val[i] = n
		i++
	}

	_, err := NewAccountBalanceQuery().
		SetAccountID(val[0]).
		Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func DisabledTestIntegrationClientPingAllBadNetwork(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)

	tempClient := _NewClient(env.Client.GetNetwork(), env.Client.GetMirrorNetwork(), *env.Client.GetNetworkName())
	tempClient.SetOperator(env.OperatorID, env.OperatorKey)

	tempClient.SetMaxNodeAttempts(1)
	tempClient.SetMaxNodesPerTransaction(2)
	tempClient.SetMaxAttempts(3)
	net := tempClient.GetNetwork()
	assert.True(t, len(net) > 1)

	keys := make([]string, len(net))
	val := make([]AccountID, len(net))
	i := 0
	for st, n := range net {
		keys[i] = st
		val[i] = n
		i++
	}

	tempNet := make(map[string]AccountID, 2)
	tempNet["in.process.ew:3123"] = val[0]
	tempNet[keys[1]] = val[1]

	err := tempClient.SetNetwork(tempNet)
	require.NoError(t, err)

	tempClient.PingAll()

	net = tempClient.GetNetwork()
	i = 0
	for st, n := range net {
		keys[i] = st
		val[i] = n
		i++
	}

	_, err = NewAccountBalanceQuery().
		SetAccountID(val[0]).
		Execute(tempClient)
	require.NoError(t, err)

	assert.Equal(t, 1, len(tempClient.GetNetwork()))

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
