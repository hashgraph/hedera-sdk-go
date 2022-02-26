package hedera

import (
	"crypto/rand"
	"math"
	"math/big"
	"sort"
	"time"
)

type _ManagedNetwork struct {
	network                map[string][]_IManagedNode
	nodes                  []_IManagedNode
	maxNodeAttempts        int
	minBackoff             time.Duration
	maxBackoff             time.Duration
	maxNodesPerTransaction *int
	ledgerID               *LedgerID
	transportSecurity      bool
	verifyCertificate      bool
}

func _NewManagedNetwork() _ManagedNetwork {
	return _ManagedNetwork{
		network:                make(map[string][]_IManagedNode),
		nodes:                  make([]_IManagedNode, 0),
		maxNodeAttempts:        -1,
		minBackoff:             250 * time.Millisecond,
		maxBackoff:             8 * time.Second,
		maxNodesPerTransaction: nil,
		ledgerID:               nil,
		transportSecurity:      false,
		verifyCertificate:      false,
	}
}

func (network *_ManagedNetwork) _SetNetwork(net map[string]_IManagedNode) error {
	for _, index := range network._GetNodesToRemove(net) {
		node := network.nodes[index]

		network._RemoveNodeFromNetwork(node)
		_ = node._Close()
		network.nodes = append(network.nodes[:index], network.nodes[index+1:]...)
	}

	for url, originalNode := range net {
		switch n := originalNode.(type) {
		case *_Node:
			nodesForKey := network._GetNodesForKey(n.accountID.String())
			if !_AddressIsInNodeList(url, *nodesForKey) {
				*nodesForKey = append(*nodesForKey, n)
				network.nodes = append(network.nodes, n)
				network.network[n.accountID.String()] = *nodesForKey
			}
		case *_MirrorNode:
			nodesForKey := network._GetNodesForKey(url)
			if !_AddressIsInNodeList(url, *nodesForKey) {
				*nodesForKey = append(*nodesForKey, n)
				network.nodes = append(network.nodes, n)
				network.network[url] = *nodesForKey
			}
		}
	}

	for i := range network.nodes {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		network.nodes[i], network.nodes[j.Int64()] = network.nodes[j.Int64()], network.nodes[i]
	}

	return nil
}

func (network *_ManagedNetwork) _GetNumberOfNodesForTransaction() int {
	if network.maxNodesPerTransaction != nil {
		return int(math.Min(float64(*network.maxNodesPerTransaction), float64(len(network.network))))
	}

	return (len(network.network) + 3 - 1) / 3
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

func (network *_ManagedNetwork) _SetMinBackoff(backoff time.Duration) {
	network.minBackoff = backoff
	for _, nod := range network.nodes {
		if nod != nil {
			nod._SetMinBackoff(backoff)
		}
	}
}

func (network *_ManagedNetwork) _GetMinBackoff() time.Duration {
	return network.minBackoff
}

func (network *_ManagedNetwork) _SetMaxBackoff(backoff time.Duration) {
	network.maxBackoff = backoff
	for _, nod := range network.nodes {
		if nod != nil {
			nod._SetMaxBackoff(backoff)
		}
	}
}

func (network *_ManagedNetwork) _GetMaxBackoff() time.Duration {
	return network.maxBackoff
}

func (network *_ManagedNetwork) _GetLedgerID() *LedgerID {
	return network.ledgerID
}

func (network *_ManagedNetwork) _SetLedgerID(id LedgerID) *_ManagedNetwork {
	network.ledgerID = &id

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

		for i, node := range network.nodes {
			if transportSecurity {
				node = node._ToSecure()
			} else {
				node = node._ToInsecure()
			}
			network.nodes[i] = node
		}

		for id := range network.network {
			nodesForKey := network._GetNodesForKey(id)
			tempArr := make([]_IManagedNode, 0)
			if transportSecurity {
				for _, tempNode := range *nodesForKey {
					temp := tempNode._ToSecure()
					tempArr = append(tempArr, temp)
				}
				switch n := tempArr[0].(type) {
				case *_Node:
					network.network[n.accountID.String()] = tempArr
				case *_MirrorNode:
					delete(network.network, id)
					network.network[n._GetAddress()] = tempArr
				}
			} else {
				switch n := tempArr[0].(type) {
				case *_Node:
					network.network[n.accountID.String()] = tempArr
				case *_MirrorNode:
					delete(network.network, id)
					network.network[n._GetAddress()] = tempArr
				}
			}
		}
	}

	network.transportSecurity = transportSecurity
}

func (network *_ManagedNetwork) _GetNumberOfMostHealthyNodes(count int32) []_IManagedNode {
	sort.Slice(network.nodes, func(i int, j int) bool {
		return _ManagedNodeCompare(network.nodes[i].(*_Node)._ManagedNode, network.nodes[j].(*_Node)._ManagedNode) < 0
	})

	err := network._RemoveDeadNodes()
	if err != nil {
		return []_IManagedNode{}
	}

	for _, networkPerNodeAccountID := range network.network {
		net := networkPerNodeAccountID
		sort.Slice(net, func(i int, j int) bool {
			return _ManagedNodeCompare(net[i].(*_Node)._ManagedNode, net[j].(*_Node)._ManagedNode) < 0
		})
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
		current := network.network[n.accountID.String()]
		if len(current) == 0 {
			delete(network.network, n.accountID.String())
			return
		}
		index := -1
		for i, n2 := range current {
			if n._GetAddress() == n2._GetAddress() {
				index = i
			}
		}
		if index != -1 {
			current = append(current[:index], current[index+1:]...)
		}
		if len(current) == 0 {
			delete(network.network, n.accountID.String())
			return
		}
		network.network[n.accountID.String()] = current
	case *_MirrorNode:
		delete(network.network, n._GetAddress())
	}
}

func _AddressIsInNodeList(addressString string, nodeArray []_IManagedNode) bool {
	for _, node := range nodeArray {
		if node._GetManagedNode()._GetAddress() == addressString {
			return true
		}
	}

	return false
}

func (network *_ManagedNetwork) _CheckNetworkContainsEntry(node _IManagedNode) bool { //nolint
	for _, n := range network.network {
		for _, nod := range n {
			if nod._GetAddress() == node._GetAddress() {
				return true
			}
		}
	}

	return false
}

func (network *_ManagedNetwork) _GetNodesForKey(key string) *[]_IManagedNode {
	if len(network.network[key]) > 0 {
		temp := network.network[key]
		return &temp
	}

	temp := make([]_IManagedNode, 0)
	network.network[key] = temp
	return &temp
}

func (network *_ManagedNetwork) _GetNodesToRemove(net map[string]_IManagedNode) []int {
	nodes := make([]int, 0)

	for i := len(network.nodes) - 1; i >= 0; i-- {
		node := network.nodes[i]

		if !_NodeIsInGivenNetwork(node, net) {
			nodes = append(nodes, i)
		}
	}

	return nodes
}

func _NodeIsInGivenNetwork(node _IManagedNode, givenNetwork map[string]_IManagedNode) bool {
	switch nodeType := node.(type) {
	case *_Node:
		for _, n := range givenNetwork {
			switch n1 := n.(type) { //nolint
			case *_Node:
				if nodeType.accountID.String() == n1.accountID.String() &&
					nodeType.address._String() == n1._GetAddress() {
					return true
				}
			}
		}
	case *_MirrorNode:
		for _, n := range givenNetwork {
			switch n1 := n.(type) { //nolint
			case *_MirrorNode:
				if nodeType.address._String() == n1._GetAddress() {
					return true
				}
			}
		}
	}

	return false
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
