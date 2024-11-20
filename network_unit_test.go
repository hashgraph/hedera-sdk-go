//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newNetworkMockNodes() map[string]AccountID {
	nodes := make(map[string]AccountID, 2)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["2.testnet.hedera.com:50211"] = AccountID{0, 0, 5, nil, nil, nil}
	nodes["3.testnet.hedera.com:50211"] = AccountID{0, 0, 6, nil, nil, nil}
	nodes["4.testnet.hedera.com:50211"] = AccountID{0, 0, 7, nil, nil, nil}
	return nodes
}

func TestUnitNetworkAddressBookGetsSet(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	network._SetTransportSecurity(true)

	ledgerID, err := LedgerIDFromString("mainnet")
	require.NoError(t, err)

	network._SetLedgerID(*ledgerID)
	require.NoError(t, err)

	require.True(t, network.addressBook != nil)
}

func TestUnitNetworkIncreaseBackoffConcurrent(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	nodes := newNetworkMockNodes()
	err := network.SetNetwork(nodes)
	require.NoError(t, err)

	node := network._GetNode()
	require.NotNil(t, node)

	numThreads := 20
	var wg sync.WaitGroup
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			network._IncreaseBackoff(node)
			wg.Done()
		}()
	}
	wg.Wait()

	require.Equal(t, len(nodes)-1, len(network.healthyNodes))
}

func TestUnitConcurrentGetNodeReadmit(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	nodes := newNetworkMockNodes()
	err := network.SetNetwork(nodes)
	network._SetMinNodeReadmitPeriod(0)
	network._SetMaxNodeReadmitPeriod(0)
	require.NoError(t, err)

	for _, node := range network.nodes {
		node._SetMaxBackoff(-1 * time.Minute)
	}

	numThreads := 3
	var wg sync.WaitGroup
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			for i := 0; i < 20; i++ {
				node := network._GetNode()
				network._IncreaseBackoff(node)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	network._ReadmitNodes()
	require.Equal(t, len(nodes), len(network.healthyNodes))
}

func TestUnitConcurrentNodeAccess(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	nodes := newNetworkMockNodes()
	err := network.SetNetwork(nodes)
	network._SetMinNodeReadmitPeriod(0)
	network._SetMaxNodeReadmitPeriod(0)
	require.NoError(t, err)

	for _, node := range network.nodes {
		node._SetMaxBackoff(-1 * time.Minute)
	}

	numThreads := 3
	var wg sync.WaitGroup
	node := network._GetNode()
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			for i := 0; i < 20; i++ {
				network._GetNode()
				network._IncreaseBackoff(node)
				node._IsHealthy()
				node._GetAttempts()
				node._GetReadmitTime()
				node._Wait()
				node._InUse()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	network._ReadmitNodes()
	require.Equal(t, len(nodes), len(network.healthyNodes))
}

func TestUnitConcurrentNodeGetChannel(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	nodes := newNetworkMockNodes()
	err := network.SetNetwork(nodes)
	require.NoError(t, err)

	numThreads := 20
	var wg sync.WaitGroup
	node := network._GetNode()
	wg.Add(numThreads)
	logger := NewLogger("", LoggerLevelError)
	for i := 0; i < numThreads; i++ {
		go func() {
			node._GetChannel(logger)
			wg.Done()
		}()
	}
	wg.Wait()
	network._ReadmitNodes()
	require.Equal(t, len(nodes), len(network.healthyNodes))
}
