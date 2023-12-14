//go:build all || e2e
// +build all e2e

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationClientCanExecuteSerializedTransactionFromAnotherClient(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
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

	netwrk := _NewNetwork()
	netwrk.SetNetwork(env.Client.GetNetwork())

	tempClient := _NewClient(netwrk, env.Client.GetMirrorNetwork(), env.Client.GetLedgerID())
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
