package hedera

import (
	"crypto/rand"
	"math"
	"math/big"
	"reflect"
	"time"
)

type _ManagedNetwork struct {
	network                map[string][]_IManagedNode
	nodes                  []_IManagedNode
	goodNodes              []_IManagedNode
	maxNodeAttempts        int
	minBackoff             time.Duration
	maxBackoff             time.Duration
	maxNodesPerTransaction *int
	ledgerID               *LedgerID
	transportSecurity      bool
	verifyCertificate      bool
	nodeMinReadmitPeriod   time.Duration
	nodeMaxReadmitPeriod   time.Duration
	earliestReadmitTime    int64
}

func _NewManagedNetwork() _ManagedNetwork {
	return _ManagedNetwork{
		network:                make(map[string][]_IManagedNode),
		nodes:                  make([]_IManagedNode, 0),
		goodNodes:              make([]_IManagedNode, 0),
		maxNodeAttempts:        -1,
		minBackOff:             8 * time.Second,
		maxBackOff:             1000 * 60 * 60 * time.Millisecond,
		maxNodesPerTransaction: nil,
		ledgerID:               nil,
		transportSecurity:      false,
		verifyCertificate:      false,
		nodeMinReadmitPeriod:   8 * time.Second,
		nodeMaxReadmitPeriod:   1000 * 60 * 60 * time.Millisecond,
		earliestReadmitTime:    time.Now().UnixNano() + 1000*60*60*time.Millisecond.Milliseconds(),
	}
}

func (network *_ManagedNetwork) _SetNetwork(net map[string]_IManagedNode) error {
	for _, index := range network._GetNodesToRemove(net) {
		node := network.nodes[index]

		network._RemoveNodeFromNetwork(node)
		_ = node._Close()
		if index == len(network.nodes)-1 {
			network.nodes = network.nodes[:index]
		} else {
			network.nodes = append(network.nodes[:index], network.nodes[index+1:]...)
		}
	}

	for _, index := range network._GetGoodNodesToRemove(net) {
		if index == len(network.goodNodes)-1 {
			network.goodNodes = network.goodNodes[:index]
		} else {
			network.goodNodes = append(network.goodNodes[:index], network.goodNodes[index+1:]...)
		}
	}

	for url, originalNode := range net {
		switch n := originalNode.(type) {
		case *_Node:
			nodesForKey := network._GetNodesForKey(n.accountID.String())
			if !_AddressIsInNodeList(url, *nodesForKey) {
				*nodesForKey = append(*nodesForKey, n)
				network.goodNodes = append(network.goodNodes, n)
				network.nodes = append(network.nodes, n)
				network.network[n.accountID.String()] = *nodesForKey
			}
		case *_MirrorNode:
			nodesForKey := network._GetNodesForKey(url)
			if !_AddressIsInNodeList(url, *nodesForKey) {
				*nodesForKey = append(*nodesForKey, n)
				network.goodNodes = append(network.goodNodes, n)
				network.nodes = append(network.nodes, n)
				network.network[url] = *nodesForKey
			}
		}
	}

	for i := range network.goodNodes {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		network.goodNodes[i], network.goodNodes[j.Int64()] = network.goodNodes[j.Int64()], network.goodNodes[i]
	}

	return nil
}

func (network *_ManagedNetwork) _ReadmitNodes() {
	now := time.Now().UnixNano()

	if network.earliestReadmitTime <= now {
		nextEarliestReadmitTime := int64(math.MaxInt64)
		searchForNextEarliestReadmitTime := true

	outer:
		for _, node := range network.nodes {
			for _, node2 := range network.goodNodes {
				if searchForNextEarliestReadmitTime && node._GetReadmitTime() > now {
					nextEarliestReadmitTime = int64(math.Min(float64(node._GetReadmitTime()), float64(nextEarliestReadmitTime)))
				}

				if reflect.DeepEqual(node, node2) {
					continue outer
				}
			}

			searchForNextEarliestReadmitTime = false

			if node._GetReadmitTime() <= now {
				network.goodNodes = append(network.goodNodes, node)
			}
		}

		network.earliestReadmitTime = int64(math.Min(
			math.Max(float64(nextEarliestReadmitTime), float64(network.nodeMinReadmitPeriod)),
			float64(network.nodeMaxReadmitPeriod)))
	}
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

func (network *_ManagedNetwork) _SetNodeMinReadmitPeriod(min time.Duration) {
	network.nodeMinReadmitPeriod = min
	network.earliestReadmitTime = time.Now().UnixNano() + network.nodeMinReadmitPeriod.Nanoseconds()
}

func (network *_ManagedNetwork) _GetNodeMinReadmitPeriod() time.Duration {
	return network.nodeMinReadmitPeriod
}

func (network *_ManagedNetwork) _SetNodeMaxReadmitPeriod(max time.Duration) {
	network.nodeMaxReadmitPeriod = max
}

func (network *_ManagedNetwork) _GetNodeMaxReadmitPeriod() time.Duration {
	return network.nodeMaxReadmitPeriod
}

func (network *_ManagedNetwork) _SetMinBackoff(waitTime time.Duration) {
	network.minBackOff = waitTime
	for _, nod := range network.goodNodes {
		if nod != nil {
			nod._SetMinBackoff(backoff)
		}
	}
}

func (network *_ManagedNetwork) _GetNode() _IManagedNode {
	network._ReadmitNodes()

	if len(network.goodNodes) == 0 {
		panic("failed to find a healthy working node")
	}

	bg := big.NewInt(int64(len(network.goodNodes)))

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	return network.goodNodes[n.Int64()]
}

func (network *_ManagedNetwork) _GetMinBackoff() time.Duration {
	return network.minBackoff
}

func (network *_ManagedNetwork) _SetMaxBackoff(waitTime time.Duration) {
	network.maxBackOff = waitTime
	for _, nod := range network.goodNodes {
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
	for _, conn := range network.goodNodes {
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

		for i, node := range network.goodNodes {
			if transportSecurity {
				node = node._ToSecure()
			} else {
				node = node._ToInsecure()
			}
			network.goodNodes[i] = node
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

//func (network *_ManagedNetwork) _GetNumberOfMostHealthyNodes(count int32) []_IManagedNode {
//	sort.Sort(_Nodes{nodes: network.goodNodes})
//	err := network._RemoveDeadNodes()
//	for _, n := range network.network {
//		sort.Sort(_Nodes{nodes: n})
//	}
//	if err != nil {
//		return []_IManagedNode{}
//	}
//
//	nodes := make([]_IManagedNode, 0)
//
//	size := math.Min(float64(count), float64(len(network.goodNodes)))
//	for i := 0; i < int(size); i++ {
//		nodes = append(nodes, network.goodNodes[i])
//	}
//
//	return nodes
//}
//
//func (network *_ManagedNetwork) _RemoveDeadNodes() error {
//	if network.maxNodeAttempts > 0 {
//		for i := len(network.goodNodes) - 1; i >= 0; i-- {
//			node := network.goodNodes[i]
//			if node == nil {
//				panic("found nil in network list")
//			}
//
//			if node._GetAttempts() >= int64(network.maxNodeAttempts) {
//				err := node._Close()
//				if err != nil {
//					return err
//				}
//				network._RemoveNodeFromNetwork(node)
//				network.goodNodes = append(network.goodNodes[:i], network.goodNodes[i+1:]...)
//			}
//		}
//	}
//
//	return nil
//}
//
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
			if len(current)-1 == index {
				current = current[:index]
			} else {
				current = append(current[:index], current[index+1:]...)
			}
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

func (network *_ManagedNetwork) _GetGoodNodesToRemove(net map[string]_IManagedNode) []int {
	nodes := make([]int, 0)

	for i := len(network.goodNodes) - 1; i >= 0; i-- {
		node := network.goodNodes[i]

		if !_NodeIsInGivenNetwork(node, net) {
			nodes = append(nodes, i)
		}
	}

	return nodes
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
	if len(network.goodNodes) > 0 {
		for i, node := range network.goodNodes {
			switch s := node.(type) { //nolint
			case *_Node:
				s._SetCertificateVerification(verify)
				network.goodNodes[i] = s
			}
		}
	}

	network.verifyCertificate = verify

	return network
}

func (network *_ManagedNetwork) _GetVerifyCertificate() bool {
	if len(network.goodNodes) > 0 {
		for _, node := range network.goodNodes {
			switch s := node.(type) { //nolint
			case *_Node:
				return s._GetCertificateVerification()
			}
		}
	}

	return false
}
