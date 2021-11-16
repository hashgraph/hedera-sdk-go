package hedera

import (
	"math"
	"sort"
	"time"
)

type _ManagedNetwork struct {
	network                map[string]_IManagedNode
	nodes                  []_IManagedNode
	maxNodeAttempts        int
	nodeWaitTime           time.Duration
	maxNodesPerTransaction *int
	networkName            *NetworkName
	transportSecurity      bool
	verifyCertificate      bool
}

func _NewManagedNetwork() _ManagedNetwork {
	return _ManagedNetwork{
		network:                make(map[string]_IManagedNode),
		nodes:                  make([]_IManagedNode, 0),
		maxNodeAttempts:        -1,
		nodeWaitTime:           250 * time.Millisecond,
		maxNodesPerTransaction: nil,
		transportSecurity:      false,
	}
}

func (network *_ManagedNetwork) _SetNetwork(net map[string]_IManagedNode) error {
	if len(network.nodes) == 0 {
		for url, originalNode := range net {
			network.nodes = append(network.nodes, originalNode)
			network.network[url] = originalNode
		}

		return nil
	}

	for url, node := range network._GetNodesToRemove(net) {
		for i, n := range network.nodes {
			if n._GetAddress() == node._GetAddress() {
				network.nodes = append(network.nodes[:i], network.nodes[i+1:]...)
			}
		}
		_ = node._Close()
		delete(network.network, url)
	}

	for url, originalNode := range net {
		if !network._CheckNetworkContainsEntry(originalNode) {
			network.nodes = append(network.nodes, originalNode)
			network.network[url] = originalNode
		}
	}

	return nil
}

func (network *_ManagedNetwork) _GetNumberOfNodesForTransaction() int {
	if network.maxNodesPerTransaction != nil {
		return int(math.Min(float64(*network.maxNodesPerTransaction), float64(len(network.nodes))))
	}

	return (len(network.nodes) + 3 - 1) / 3
}

func (network *_ManagedNetwork) _SetMaxNodesPerTransaction(max int) {
	network.maxNodesPerTransaction = &max
}

func (network *_ManagedNetwork) _SetMaxNodeAttempts(max int) {
	network.maxNodeAttempts = max
}

func (network *_ManagedNetwork) _GetMaxNodeAttempts() int {
	return network.maxNodeAttempts
}

func (network *_ManagedNetwork) _SetNodeWaitTime(waitTime time.Duration) {
	network.nodeWaitTime = waitTime
	for _, nod := range network.nodes {
		if nod != nil {
			nod._SetMinBackoff(waitTime.Milliseconds())
		}
	}
}

func (network *_ManagedNetwork) _GetNodeWaitTime() time.Duration {
	return network.nodeWaitTime
}

func (network *_ManagedNetwork) _GetNetworkName() *NetworkName {
	return network.networkName
}

func (network *_ManagedNetwork) _SetNetworkName(net NetworkName) *_ManagedNetwork {
	network.networkName = &net

	return network
}

func (network *_ManagedNetwork) _Close() error {
	for _, conn := range network.nodes {
		err := conn._Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (network *_ManagedNetwork) _SetTransportSecurity(transportSecurity bool) {
	if network.transportSecurity != transportSecurity {
		_ = network._Close()
		network.network = make(map[string]_IManagedNode)

		for i, node := range network.nodes {
			if transportSecurity {
				node = node._ToSecure()
			} else {
				node = node._ToInsecure()
			}
			network.nodes[i] = node
			switch n := node.(type) {
			case *_Node:
				network.network[n.accountID.String()] = n._ToSecure()
			case *_MirrorNode:
				network.network[n._GetAddress()] = n
			}
		}
	}

	network.transportSecurity = transportSecurity
}

func (network *_ManagedNetwork) _GetNumberOfMostHealthyNodes(count int32) []_IManagedNode {
	sort.Sort(_Nodes{nodes: network.nodes})
	err := network._RemoveDeadNodes()
	if err != nil {
		return []_IManagedNode{}
	}

	nodes := make([]_IManagedNode, 0)

	size := math.Min(float64(count), float64(len(network.nodes)))
	for i := 0; i < int(size); i++ {
		nodes = append(nodes, network.nodes[i])
	}

	return nodes
}

func (network *_ManagedNetwork) _RemoveDeadNodes() error {
	if network.maxNodeAttempts > 0 {
		for i := len(network.nodes) - 1; i >= 0; i-- {
			node := network.nodes[i]
			if node == nil {
				panic("found nil in network list")
			}

			if node._GetAttempts() >= int64(network.maxNodeAttempts) {
				err := node._Close()
				if err != nil {
					return err
				}
				network._RemoveNodeFromNetwork(node)
				network.nodes = append(network.nodes[:i], network.nodes[i+1:]...)
			}
		}
	}

	return nil
}

func (network *_ManagedNetwork) _RemoveNodeFromNetwork(node _IManagedNode) {
	switch n := node.(type) {
	case *_Node:
		delete(network.network, n.accountID.String())
	case *_MirrorNode:
		delete(network.network, n._GetAddress())
	}
}

func (network *_ManagedNetwork) _CheckNetworkContainsEntry(node _IManagedNode) bool {
	for _, n := range network.network {
		if n._GetAddress() == node._GetAddress() {
			return true
		}
	}

	return false
}

func (network *_ManagedNetwork) _GetNodesToRemove(net map[string]_IManagedNode) map[string]_IManagedNode {
	nodes := make(map[string]_IManagedNode)

	var b bool
	for url, node := range network.network {
		b = false
		for _, n := range net {
			if n._GetAddress() == node._GetAddress() {
				b = true
			}
		}
		if !b {
			nodes[url] = node
		}
	}

	return nodes
}

func (network *_ManagedNetwork) _SetVerifyCertificate(verify bool) *_ManagedNetwork {
	if len(network.nodes) > 0 {
		for i, node := range network.nodes {
			switch s := node.(type) { //nolint
			case *_Node:
				s._SetCertificateVerification(verify)
				network.nodes[i] = s
			}
		}
	}

	network.verifyCertificate = verify

	return network
}

func (network *_ManagedNetwork) _GetVerifyCertificate() bool {
	if len(network.nodes) > 0 {
		for _, node := range network.nodes {
			switch s := node.(type) { //nolint
			case *_Node:
				return s._GetCertificateVerification()
			}
		}
	}

	return false
}
