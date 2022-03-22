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
	healthyNodes           []_IManagedNode
	maxNodeAttempts        int
	minBackoff             time.Duration
	maxBackoff             time.Duration
	maxNodesPerTransaction *int
	ledgerID               *LedgerID
	transportSecurity      bool
	verifyCertificate      bool
	minNodeReadmitPeriod   time.Duration
	maxNodeReadmitPeriod   time.Duration
	earliestReadmitTime    time.Time
}

func _NewManagedNetwork() _ManagedNetwork {
	return _ManagedNetwork{
		network:                make(map[string][]_IManagedNode),
		nodes:                  make([]_IManagedNode, 0),
		healthyNodes:           make([]_IManagedNode, 0),
		maxNodeAttempts:        -1,
		minBackoff:             8 * time.Second,
		maxBackoff:             1 * time.Hour,
		maxNodesPerTransaction: nil,
		ledgerID:               nil,
		transportSecurity:      false,
		verifyCertificate:      false,
		minNodeReadmitPeriod:   8 * time.Second,
		maxNodeReadmitPeriod:   1 * time.Hour,
		earliestReadmitTime:    time.Unix(math.MaxInt64, math.MaxInt32),
	}
}

func (this *_ManagedNetwork) _SetNetwork(network map[string]_IManagedNode) error {
	newNodes := make([]_IManagedNode, len(this.nodes))
	newHealthyNodes := make([]_IManagedNode, len(this.nodes))
	newNetwork := map[string][]_IManagedNode{}
	newNodeKeys := map[string]bool{}
	newNodeValues := map[string]bool{}

	// Remove nodes from the old this which do not belong to the new this
	for _, index := range this._GetNodesToRemove(network) {
		node := this.nodes[index]

		this._RemoveNodeFromNetwork(node)
		_ = node._Close()
		if index == len(this.nodes)-1 {
			this.nodes = this.nodes[:index]
		} else {
			this.nodes = append(this.nodes[:index], this.nodes[index+1:]...)
		}
	}

	// Copy all the nodes into the `newNodes` list
	copy(newNodes, this.nodes)

	for _, node := range newNodes {
		newNodeKeys[node._GetKey()] = true
		newNodeValues[node._GetAddress()] = true
	}

	for key, value := range network {
		_, keyOk := newNodeKeys[key]
		_, valueOk := newNodeValues[value._GetAddress()]

		if keyOk && valueOk {
			continue
		}

		newNodes = append(newNodes, value)
	}

	for _, node := range newNodes {
		newHealthyNodes = append(newHealthyNodes, node)

		value, ok := this.network[node._GetKey()]
		if !ok {
			value = []_IManagedNode{}
		}
		value = append(value, node)
		this.network[node._GetKey()] = value
	}

	this.nodes = newNodes
	this.network = newNetwork
	this.healthyNodes = newHealthyNodes
	this.ledgerID = nil

	return nil
}

func (this *_ManagedNetwork) _ReadmitNodes() {
	now := time.Now()

	nextEarliestReadmitTime := time.Now().Add(this.maxNodeReadmitPeriod)

	for _, node := range this.nodes {
		if node._GetReadmitTime().After(now) && node._GetReadmitTime().Before(nextEarliestReadmitTime) {
			nextEarliestReadmitTime = node._GetReadmitTime()
		}
	}

	if nextEarliestReadmitTime.Before(now.Add(this.minNodeReadmitPeriod)) {
		nextEarliestReadmitTime = now.Add(this.minNodeReadmitPeriod)
	}

outer:
	for _, node := range this.nodes {
		for _, node2 := range this.healthyNodes {
			if reflect.DeepEqual(node, node2) {
				continue outer
			}
		}

		if node._GetReadmitTime().Before(now) {
			this.healthyNodes = append(this.healthyNodes, node)
		}
	}
}

func (this *_ManagedNetwork) _GetNumberOfNodesForTransaction() int {
	this._ReadmitNodes()
	if this.maxNodesPerTransaction != nil {
		return int(math.Min(float64(*this.maxNodesPerTransaction), float64(len(this.network))))
	}

	return (len(this.network) + 3 - 1) / 3
}

func (this *_ManagedNetwork) _SetMaxNodesPerTransaction(max int) {
	this.maxNodesPerTransaction = &max
}

func (this *_ManagedNetwork) _SetMaxNodeAttempts(max int) {
	this.maxNodeAttempts = max
}

func (this *_ManagedNetwork) _GetMaxNodeAttempts() int {
	return this.maxNodeAttempts
}

func (this *_ManagedNetwork) _SetNodeMinReadmitPeriod(min time.Duration) {
	this.minNodeReadmitPeriod = min
	this.earliestReadmitTime = time.Now().Add(this.minNodeReadmitPeriod)
}

func (this *_ManagedNetwork) _GetNodeMinReadmitPeriod() time.Duration {
	return this.minNodeReadmitPeriod
}

func (this *_ManagedNetwork) _SetNodeMaxReadmitPeriod(max time.Duration) {
	this.maxNodeReadmitPeriod = max
}

func (this *_ManagedNetwork) _GetNodeMaxReadmitPeriod() time.Duration {
	return this.maxNodeReadmitPeriod
}

func (this *_ManagedNetwork) _SetMinBackoff(minBackoff time.Duration) {
	this.minBackoff = minBackoff
	for _, nod := range this.healthyNodes {
		if nod != nil {
			nod._SetMinBackoff(minBackoff)
		}
	}
}

func (this *_ManagedNetwork) _GetNode() _IManagedNode {
	this._ReadmitNodes()

	if len(this.healthyNodes) == 0 {
		panic("failed to find a healthy working node")
	}

	bg := big.NewInt(int64(len(this.healthyNodes)))

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	return this.healthyNodes[n.Int64()]
}

func (this *_ManagedNetwork) _GetMinBackoff() time.Duration {
	return this.minBackoff
}

func (this *_ManagedNetwork) _SetMaxBackoff(maxBackoff time.Duration) {
	this.maxBackoff = maxBackoff
	for _, nod := range this.healthyNodes {
		if nod != nil {
			nod._SetMaxBackoff(maxBackoff)
		}
	}
}

func (this *_ManagedNetwork) _GetMaxBackoff() time.Duration {
	return this.maxBackoff
}

func (this *_ManagedNetwork) _GetLedgerID() *LedgerID {
	return this.ledgerID
}

func (this *_ManagedNetwork) _SetLedgerID(id LedgerID) *_ManagedNetwork {
	this.ledgerID = &id

	return this
}

func (this *_ManagedNetwork) _Close() error {
	for _, conn := range this.healthyNodes {
		err := conn._Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *_ManagedNetwork) _SetTransportSecurity(transportSecurity bool) {
	if this.transportSecurity != transportSecurity {
		_ = this._Close()

		for i, node := range this.healthyNodes {
			if transportSecurity {
				node = node._ToSecure()
			} else {
				node = node._ToInsecure()
			}
			this.healthyNodes[i] = node
		}

		for id := range this.network {
			nodesForKey := this._GetNodesForKey(id)
			tempArr := make([]_IManagedNode, 0)
			if transportSecurity {
				for _, tempNode := range *nodesForKey {
					temp := tempNode._ToSecure()
					tempArr = append(tempArr, temp)
				}
				switch n := tempArr[0].(type) {
				case *_Node:
					this.network[n.accountID.String()] = tempArr
				case *_MirrorNode:
					delete(this.network, id)
					this.network[n._GetAddress()] = tempArr
				}
			} else {
				switch n := tempArr[0].(type) {
				case *_Node:
					this.network[n.accountID.String()] = tempArr
				case *_MirrorNode:
					delete(this.network, id)
					this.network[n._GetAddress()] = tempArr
				}
			}
		}
	}

	this.transportSecurity = transportSecurity
}

//func (network *_ManagedNetwork) _RemoveDeadNodes() error {
//	if network.maxNodeAttempts > 0 {
//		for i := len(network.healthyNodes) - 1; i >= 0; i-- {
//			node := network.healthyNodes[i]
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
//				network.healthyNodes = append(network.healthyNodes[:i], network.healthyNodes[i+1:]...)
//			}
//		}
//	}
//
//	return nil
//}
//
func (this *_ManagedNetwork) _RemoveNodeFromNetwork(node _IManagedNode) {
	switch n := node.(type) {
	case *_Node:
		current := this.network[n.accountID.String()]
		if len(current) == 0 {
			delete(this.network, n.accountID.String())
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
			delete(this.network, n.accountID.String())
			return
		}
		this.network[n.accountID.String()] = current
	case *_MirrorNode:
		delete(this.network, n._GetAddress())
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

func (this *_ManagedNetwork) _CheckNetworkContainsEntry(node _IManagedNode) bool { //nolint
	for _, n := range this.network {
		for _, nod := range n {
			if nod._GetAddress() == node._GetAddress() {
				return true
			}
		}
	}

	return false
}

func (this *_ManagedNetwork) _GetNodesForKey(key string) *[]_IManagedNode {
	if len(this.network[key]) > 0 {
		temp := this.network[key]
		return &temp
	}

	temp := make([]_IManagedNode, 0)
	this.network[key] = temp
	return &temp
}

func (this *_ManagedNetwork) _GetGoodNodesToRemove(net map[string]_IManagedNode) []int {
	nodes := make([]int, 0)

	for i := len(this.healthyNodes) - 1; i >= 0; i-- {
		node := this.healthyNodes[i]

		if !_NodeIsInGivenNetwork(node, net) {
			nodes = append(nodes, i)
		}
	}

	return nodes
}

func (this *_ManagedNetwork) _GetNodesToRemove(net map[string]_IManagedNode) []int {
	nodes := make([]int, 0)

	for i := len(this.nodes) - 1; i >= 0; i-- {
		node := this.nodes[i]

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

func (this *_ManagedNetwork) _SetVerifyCertificate(verify bool) *_ManagedNetwork {
	if len(this.healthyNodes) > 0 {
		for i, node := range this.healthyNodes {
			switch s := node.(type) { //nolint
			case *_Node:
				s._SetCertificateVerification(verify)
				this.healthyNodes[i] = s
			}
		}
	}

	this.verifyCertificate = verify

	return this
}

func (this *_ManagedNetwork) _GetVerifyCertificate() bool {
	if len(this.healthyNodes) > 0 {
		for _, node := range this.healthyNodes {
			switch s := node.(type) { //nolint
			case *_Node:
				return s._GetCertificateVerification()
			}
		}
	}

	return false
}
