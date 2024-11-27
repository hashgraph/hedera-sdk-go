//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func (this *_ManagedNetwork) removeNodeFromHealthyNodes(nodeToRemove _IManagedNode) {
	newHealthyNodes := make([]_IManagedNode, 0)

	for _, node := range this.healthyNodes {
		if node != nodeToRemove {
			newHealthyNodes = append(newHealthyNodes, node)
		}
	}

	this.healthyNodes = newHealthyNodes
}

func newMockNodes() map[string]_IManagedNode {
	address1, _ := _ManagedNodeAddressFromString("node1:50211")
	address2, _ := _ManagedNodeAddressFromString("node1:50212")
	address3, _ := _ManagedNodeAddressFromString("node2:50211")

	return map[string]_IManagedNode{
		"node1:50211": &mockManagedNode{address: address1, healthy: true},
		"node1:50212": &mockManagedNode{address: address2, healthy: true},
		"node2:50211": &mockManagedNode{address: address3, healthy: true},
	}
}

type mockManagedNode struct {
	address                    *_ManagedNodeAddress
	currentBackoff             time.Duration
	lastUsed                   time.Time
	useCount                   int64
	minBackoff                 time.Duration
	maxBackoff                 time.Duration
	badGrpcStatusCount         int64
	readmitTime                *time.Time
	healthy                    bool
	minBackoffCalled           bool
	maxBackoffCalled           bool
	setVerifyCertificateCalled bool
	toSecureCalled             bool
	toInsecureCalled           bool
}
type mockManagedNodeWithError struct {
	mockManagedNode
}

func (m *mockManagedNodeWithError) _Close() error {
	return errors.New("closing error")
}

func (m *mockManagedNode) _GetAddress() string {
	return m.address._String()
}

func (m *mockManagedNode) _GetKey() string {
	return m.address._String()
}

func (m *mockManagedNode) _IsHealthy() bool {
	return m.healthy
}

func (m *mockManagedNode) _Close() error {
	return nil
}

func (m *mockManagedNode) _DecreaseBackoff() {

}

func (m *mockManagedNode) _IncreaseBackoff() {
	m.healthy = false
}

func (m *mockManagedNode) _ResetBackoff() {
	// No need to implement this for the test
}

func (m *mockManagedNode) _GetReadmitTime() *time.Time {
	return m.readmitTime
}

func (m *mockManagedNode) _GetAttempts() int64 {
	return 0
}

func (m *mockManagedNode) _GetLastUsed() time.Time {
	return time.Now()
}

func (m *mockManagedNode) _GetManagedNode() *_ManagedNode {
	return nil
}

func (m *mockManagedNode) _ToSecure() _IManagedNode {
	m.toSecureCalled = true
	return m
}

func (m *mockManagedNode) _ToInsecure() _IManagedNode {
	m.toInsecureCalled = true
	return m
}

func (m *mockManagedNode) _SetMinBackoff(minBackoff time.Duration) {
	m.minBackoff = minBackoff
	m.minBackoffCalled = true
}

func (m *mockManagedNode) _SetMaxBackoff(maxBackoff time.Duration) {
	m.maxBackoff = maxBackoff
	m.maxBackoffCalled = true
}

func (m *mockManagedNode) _GetMinBackoff() time.Duration {
	return time.Duration(0)
}

func (m *mockManagedNode) _GetMaxBackoff() time.Duration {
	return time.Duration(0)
}

func (m *mockManagedNode) _Wait() time.Duration {
	return time.Duration(0)
}

func (m *mockManagedNode) _GetUseCount() int64 {
	return 0
}

func (m *mockManagedNode) _SetVerifyCertificate(verify bool) {
	m.setVerifyCertificateCalled = true
}

func (m *mockManagedNode) _GetVerifyCertificate() bool {
	return true
}

func (m *mockManagedNode) _InUse() {
}

func TestUnitManagedNetworkSetGet(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)
	ledgerId, err := LedgerIDFromString("mainnet")
	require.NoError(t, err)
	mn._SetLedgerID(*ledgerId)
	mn._SetMaxNodeAttempts(10)
	mn._SetMinBackoff(1 * time.Second)
	mn._SetMaxBackoff(2 * time.Second)
	mn._SetMaxNodesPerTransaction(3)
	mn._SetTransportSecurity(false)
	mn._SetVerifyCertificate(false)
	mn._SetMinNodeReadmitPeriod(4 * time.Second)
	mn._SetMaxNodeReadmitPeriod(5 * time.Second)

	require.Equal(t, 10, mn.maxNodeAttempts)
	require.Equal(t, 1*time.Second, mn._GetMinBackoff())
	require.Equal(t, 2*time.Second, mn._GetMaxBackoff())
	require.Equal(t, 3, *mn.maxNodesPerTransaction)
	require.Equal(t, ledgerId, mn._GetLedgerID())
	require.False(t, mn.transportSecurity)
	require.False(t, mn._GetVerifyCertificate())
	require.Equal(t, 4*time.Second, mn.minNodeReadmitPeriod)
	require.Equal(t, 5*time.Second, mn.maxNodeReadmitPeriod)
	for _, node := range mockNodes {
		mockNode, ok := node.(*mockManagedNode)
		require.True(t, ok, "node should be of type *mockManagedNode")

		require.True(t, mockNode.minBackoffCalled, "minBackoffCalled should be true")
		require.True(t, mockNode.maxBackoffCalled, "maxBackoffCalled should be true")
		require.True(t, mockNode.setVerifyCertificateCalled, "setVerifyCertificateCalled should be true")
		// Should not be called, as those are false by default
		require.False(t, mockNode.toSecureCalled, "toSecureCalled should be false")
		require.False(t, mockNode.toInsecureCalled, "toInsecureCalled should be false")
	}
	mn._SetTransportSecurity(true)
	mn._SetVerifyCertificate(true)
	for _, node := range mockNodes {
		mockNode, ok := node.(*mockManagedNode)
		require.True(t, ok, "node should be of type *mockManagedNode")

		require.True(t, mockNode.setVerifyCertificateCalled, "setVerifyCertificateCalled should be true")
		require.True(t, mockNode.toSecureCalled, "toSecureCalled should be true")
	}
	mn._SetTransportSecurity(false)
	for _, node := range mockNodes {
		mockNode, ok := node.(*mockManagedNode)
		require.True(t, ok, "node should be of type *mockManagedNode")
		require.True(t, mockNode.toInsecureCalled, "toInsecureCalled should be true")
	}
}

func TestUnitNewManagedNetwork(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()

	require.NotNil(t, mn.network)
	require.NotNil(t, mn.nodes)
	require.NotNil(t, mn.healthyNodes)
	require.Equal(t, -1, mn._GetMaxNodeAttempts())
	require.Equal(t, 8*time.Second, mn.minBackoff)
	require.Equal(t, 1*time.Hour, mn.maxBackoff)
	require.Nil(t, mn.maxNodesPerTransaction)
	require.Nil(t, mn.ledgerID)
	require.False(t, mn.transportSecurity)
	require.False(t, mn.verifyCertificate)
	require.Equal(t, 8*time.Second, mn._GetMinNodeReadmitPeriod())
	require.Equal(t, 1*time.Hour, mn._GetMaxNodeReadmitPeriod())
}

func TestUnitSetNetwork(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Check if the nodes are properly set in the _ManagedNetwork
	require.Equal(t, 3, len(mn.nodes))
	for _, node := range mn.nodes {
		require.Contains(t, mockNodes, node._GetAddress())
	}

	// Check if the healthy nodes are properly set in the _ManagedNetwork
	require.Equal(t, 3, len(mn.healthyNodes))
	for _, node := range mn.healthyNodes {
		require.Contains(t, mockNodes, node._GetAddress())
	}

	mockNodes["node1:50211"].(*mockManagedNode).healthy = false
	err = mn._SetNetwork(mockNodes)
	require.NoError(t, err)
	// Check if only the healthy nodes are properly set in the _ManagedNetwork
	require.Equal(t, 2, len(mn.healthyNodes))
	for _, node := range mn.healthyNodes {
		require.True(t, node._IsHealthy())
	}

	// Check if the nodes are properly set in the _ManagedNetwork
	require.Equal(t, 3, len(mn.nodes))
	for _, node := range mn.nodes {
		require.Contains(t, mockNodes, node._GetAddress())
	}
}

func TestUnitSetNetworkWithErorr(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	address4, _ := _ManagedNodeAddressFromString("node1:50213")
	mockNodesWithError := map[string]_IManagedNode{
		"node1:50213": &mockManagedNodeWithError{
			mockManagedNode: mockManagedNode{address: address4, healthy: true},
		},
	}

	err := mn._SetNetwork(mockNodesWithError)
	require.NoError(t, err)
	// Add a new node, should error, because existing node return an error on close
	err = mn._SetNetwork(mockNodes)
	require.Error(t, err)
}

func TestUnitManagedNetworkCloseWithError(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNode := &mockManagedNode{}
	mockNodeWithError := &mockManagedNodeWithError{
		mockManagedNode: *mockNode,
	}

	// Inject the node with an error into the healthyNodes slice
	mn.healthyNodes = append(mn.healthyNodes, mockNodeWithError)

	err := mn._Close()
	require.Error(t, err)
	require.Equal(t, "closing error", err.Error())
}

func TestUnitManagedNetworkSetTransportSecurityWithError(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNode := &mockManagedNode{}
	mockNodeWithError := &mockManagedNodeWithError{
		mockManagedNode: *mockNode,
	}

	// Inject the node with an error into the healthyNodes slice
	mn.healthyNodes = append(mn.healthyNodes, mockNodeWithError)

	// Attempt to set the transport security
	err := mn._SetTransportSecurity(true)
	require.Error(t, err)
	require.Equal(t, "closing error", err.Error())
}

func TestUnitSetNetwork_NodeRemoved(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Remove a node from the mockNodes map
	removedNodeKey := "node1:50211"
	removedNode := mockNodes[removedNodeKey]
	delete(mockNodes, removedNodeKey)

	err = mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Check if the node was removed from the _ManagedNetwork
	require.Equal(t, 2, len(mn.nodes))
	for _, node := range mn.nodes {
		require.NotEqual(t, removedNode._GetAddress(), node._GetAddress())
	}
}

func TestUnitSetNetwork_NodeAdded(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Add a new node to the mockNodes map
	newNodeKey := "node2:50212"
	address4, _ := _ManagedNodeAddressFromString("node2:50212")
	newNode := &mockManagedNode{address: address4, healthy: true}
	mockNodes[newNodeKey] = newNode

	err = mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Check if the new node was added to the _ManagedNetwork
	require.Equal(t, 4, len(mn.nodes))
	foundNewNode := false
	for _, node := range mn.nodes {
		if node._GetAddress() == newNode._GetAddress() {
			foundNewNode = true
		}
	}
	require.True(t, foundNewNode)
}

func TestUnitSetNetworkRemoveAllNodes(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Remove all nodes from the mockNodes map
	for key := range mockNodes {
		delete(mockNodes, key)
	}

	// Set up the new network without any nodes
	err = mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Check if there are no nodes in the _ManagedNetwork
	require.Equal(t, 0, len(mn.nodes))
}

func TestUnitReadmitNodes_NodeReadmitted(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	unhealthyNodeKey := "node1:50211"
	unhealthyNode := mockNodes[unhealthyNodeKey].(*mockManagedNode)
	mn.removeNodeFromHealthyNodes(unhealthyNode)
	unhealthyNode.healthy = true

	// Set readmit time for the unhealthy node to a time before now
	pastTime := time.Now().Add(-1 * time.Minute)
	unhealthyNode.readmitTime = &pastTime

	// Call _ReadmitNodes to readmit healthy nodes
	mn._ReadmitNodes()

	// Check if the previously unhealthy node is now in the healthyNodes list
	found := false
	for _, node := range mn.healthyNodes {
		if node._GetAddress() == unhealthyNodeKey {
			found = true
			break
		}
	}
	require.True(t, found, "node1:50211 should be present in the healthyNodes list after readmission")
}

func TestUnitReadmitNodes_NodeNotReadmitted(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	unhealthyNodeKey := "node1:50212"
	unhealthyNode := mockNodes[unhealthyNodeKey].(*mockManagedNode)
	unhealthyNode.healthy = false
	mn.removeNodeFromHealthyNodes(unhealthyNode)

	// Set readmit time for the unhealthy node to a time in the future
	futureTime := time.Now().Add(1 * time.Hour)
	unhealthyNode.readmitTime = &futureTime

	// Call _ReadmitNodes
	mn._ReadmitNodes()

	// Check if the unhealthy node is not present in the healthyNodes list
	found := false
	for _, node := range mn.healthyNodes {
		if node._GetAddress() == unhealthyNodeKey {
			found = true
			break
		}
	}
	require.False(t, found, "node1:50212 should not be present in the healthyNodes list since its readmit time is in the future")
}

func TestUnitReadmitNodes_UpdateEarliestReadmitTime(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Make a node unhealthy and set its readmit time to a future time, before the minNodeReadmitPeriod
	unhealthyNodeKey := "node1:50211"
	unhealthyNode := mockNodes[unhealthyNodeKey].(*mockManagedNode)
	unhealthyNode.healthy = false
	mn.removeNodeFromHealthyNodes(unhealthyNode)

	futureReadmitTime := time.Now().Add(3 * time.Second) // Assuming minNodeReadmitPeriod is greater than 3 seconds
	unhealthyNode.readmitTime = &futureReadmitTime

	// Call _ReadmitNodes
	mn._ReadmitNodes()

	// Check if the unhealthy node is not present in the healthyNodes list
	found := false
	for _, node := range mn.healthyNodes {
		if node._GetAddress() == unhealthyNodeKey {
			found = true
			break
		}
	}
	require.False(t, found, "node1:50211 should not be present in the healthyNodes list since its readmit time is in the future")

	// Check if the earliestReadmitTime is updated to now.Add(this.minNodeReadmitPeriod)
	require.WithinDuration(t, futureReadmitTime.Add(mn.minNodeReadmitPeriod), mn.earliestReadmitTime, 5*time.Second)
}

func TestUnitGetNumberOfNodesForTransaction_Default(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	numNodes := mn._GetNumberOfNodesForTransaction()

	// Default behavior: (len(this.network) + 3 - 1) / 3
	expectedNumNodes := (len(mockNodes) + 3 - 1) / 3
	require.Equal(t, expectedNumNodes, numNodes)
}

func TestUnitGetNumberOfNodesForTransaction_MaxNodesPerTransaction(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	maxNodes := 2
	mn._SetMaxNodesPerTransaction(maxNodes)

	numNodes := mn._GetNumberOfNodesForTransaction()

	// If maxNodesPerTransaction is set, the number of nodes should be the minimum of maxNodesPerTransaction and the number of nodes in the network
	expectedNumNodes := int(math.Min(float64(maxNodes), float64(len(mockNodes))))
	require.Equal(t, expectedNumNodes, numNodes)
}

func TestUnitGetNumberOfNodesForTransaction_MaxNodesGreaterThanNetworkSize(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	maxNodes := 9
	mn._SetMaxNodesPerTransaction(maxNodes)

	numNodes := mn._GetNumberOfNodesForTransaction()

	expectedNumNodes := int(math.Min(float64(maxNodes), float64(len(mockNodes))))
	require.Equal(t, expectedNumNodes, numNodes)
}

func TestUnitGetNumberOfNodesForTransaction_MaxNodesNotSet(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	numNodes := mn._GetNumberOfNodesForTransaction()
	// 1/3 of the network size
	require.Equal(t, 1, numNodes)
}

func TestUnitGetNode(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Ensure that there are healthy nodes in the network
	require.NotEqual(t, 0, len(mn.healthyNodes))

	// Get a random node from the managed network
	node := mn._GetNode()

	// Check if the returned node is not nil
	require.NotNil(t, node)

	// Check if the returned node is one of the healthy nodes in the managed network
	found := false
	for _, healthyNode := range mn.healthyNodes {
		if node._GetAddress() == healthyNode._GetAddress() {
			found = true
			break
		}
	}
	require.True(t, found, "The returned node should be one of the healthy nodes in the managed network")
}

func TestUnitGetNodePanicNoHealthyNodes(t *testing.T) {
	t.Parallel()

	mn := _NewManagedNetwork()
	mockNodes := newMockNodes()
	err := mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Mark all nodes as unhealthy and set their readmit time in the future
	for _, node := range mockNodes {
		node.(*mockManagedNode).healthy = false
		readmitTime := time.Now().Add(1 * time.Minute)
		node.(*mockManagedNode).readmitTime = &readmitTime
	}

	// Update the network with unhealthy nodes
	err = mn._SetNetwork(mockNodes)
	require.NoError(t, err)

	// Ensure that there are no healthy nodes in the network
	require.Equal(t, 0, len(mn.healthyNodes))

	// Check if calling _GetNode() panics when there are no healthy nodes
	defer func() {
		if r := recover(); r != nil {
			panicValue, ok := r.(string)
			require.True(t, ok, "Panic value should be a string")
			require.Equal(t, "failed to find a healthy working node", panicValue)
		}
	}()

	mn._GetNode()
	require.Fail(t, "Expected _GetNode to panic")
}
