//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationClientCanExecuteSerializedTransactionFromAnotherClient(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	client2 := ClientForNetwork(env.Client.GetNetwork())
	client2.SetOperator(env.OperatorID, env.OperatorKey)

	tx, err := NewTransferTransaction().AddHbarTransfer(env.OperatorID, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(1)).SetNodeAccountIDs([]AccountID{{Account: 3}}).FreezeWith(env.Client)
	require.NoError(t, err)
	txBytes, err := tx.ToBytes()
	FromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)
	txFromBytes, ok := FromBytes.(TransferTransaction)
	require.True(t, ok)
	resp, err := txFromBytes.Execute(client2)
	require.NoError(t, err)
	reciept, err := resp.SetValidateStatus(true).GetReceipt(client2)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, reciept.Status)
}

func TestIntegrationClientCanFailGracefullyWhenDoesNotHaveNodeOfAnotherClient(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Get one of the nodes of the network from the original client
	var address string
	for key := range env.Client.GetNetwork() {
		address = key
		break
	}
	// Use that node to create a network for the second client but with a different node account id
	var network = map[string]AccountID{
		address: {Account: 99},
	}

	client2 := ClientForNetwork(network)
	client2.SetOperator(env.OperatorID, env.OperatorKey)

	// Create a transaction with a node using original client
	tx, err := NewTransferTransaction().AddHbarTransfer(env.OperatorID, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(1)).SetNodeAccountIDs([]AccountID{{Account: 3}}).FreezeWith(env.Client)
	require.NoError(t, err)
	txBytes, err := tx.ToBytes()
	FromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)
	txFromBytes, ok := FromBytes.(TransferTransaction)
	require.True(t, ok)

	// Try to execute it with the second client, which does not have the node
	_, err = txFromBytes.Execute(client2)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Invalid node AccountID was set for transaction: 0.0.3")
}

func DisabledTestIntegrationClientPingAllBadNetwork(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	netwrk := _NewNetwork()
	netwrk.SetNetwork(env.Client.GetNetwork())

	tempClient := _NewClient(netwrk, env.Client.GetMirrorNetwork(), env.Client.GetLedgerID(), true)
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

}

func TestClientInitWithMirrorNetwork(t *testing.T) {
	t.Parallel()
	mirrorNetworkString := "testnet.mirrornode.hedera.com:443"
	client, err := ClientForMirrorNetwork([]string{mirrorNetworkString})
	require.NoError(t, err)

	mirrorNetwork := client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])
	assert.NotEmpty(t, client.GetNetwork())
}
