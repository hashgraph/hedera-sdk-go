package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"crypto/rand"
	"math"
	"math/big"
	"sync"
	"time"
)

type _ManagedNetwork struct {
	network                map[string][]_IManagedNode
	nodes                  []_IManagedNode
	healthyNodes           []_IManagedNode
	healthyNodesMutex      *sync.RWMutex
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
		network:                map[string][]_IManagedNode{},
		nodes:                  []_IManagedNode{},
		healthyNodes:           []_IManagedNode{},
		healthyNodesMutex:      &sync.RWMutex{},
		maxNodeAttempts:        -1,
		minBackoff:             8 * time.Second,
		maxBackoff:             1 * time.Hour,
		maxNodesPerTransaction: nil,
		ledgerID:               nil,
		transportSecurity:      false,
		verifyCertificate:      false,
		minNodeReadmitPeriod:   8 * time.Second,
		maxNodeReadmitPeriod:   1 * time.Hour,
	}
}

func (this *_ManagedNetwork) _SetNetwork(network map[string]_IManagedNode) error {
	newNodes := make([]_IManagedNode, len(this.nodes))
	newNodeKeys := map[string]bool{}
	newNodeValues := map[string]bool{}

	// Copy all the nodes into the `newNodes` list
	copy(newNodes, this.nodes)

	// Remove nodes from the old this which do not belong to the new this
	for _, index := range _GetNodesToRemove(network, newNodes) {
		node := newNodes[index]

		if err := node._Close(); err != nil {
			return err
		}

		if index == len(newNodes)-1 {
			newNodes = newNodes[:index]
		} else {
			newNodes = append(newNodes[:index], newNodes[index+1:]...)
		}
	}

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

	newNetwork, newHealthyNodes := _CreateNetworkFromNodes(newNodes)

	this.nodes = newNodes
	this.network = newNetwork
	this.healthyNodes = newHealthyNodes

	return nil
}

func (this *_ManagedNetwork) _ReadmitNodes() {
	now := time.Now()

	this.healthyNodesMutex.Lock()
	defer this.healthyNodesMutex.Unlock()

	if this.earliestReadmitTime.Before(now) {
		nextEarliestReadmitTime := now.Add(this.maxNodeReadmitPeriod)

		for _, node := range this.nodes {
			if node._GetReadmitTime() != nil && node._GetReadmitTime().After(now) && node._GetReadmitTime().Before(nextEarliestReadmitTime) {
				nextEarliestReadmitTime = *node._GetReadmitTime()
			}
		}

		this.earliestReadmitTime = nextEarliestReadmitTime
		if this.earliestReadmitTime.Before(now.Add(this.minNodeReadmitPeriod)) {
			this.earliestReadmitTime = now.Add(this.minNodeReadmitPeriod)
		}

	outer:
		for _, node := range this.nodes {
			for _, healthyNode := range this.healthyNodes {
				if node == healthyNode {
					continue outer
				}
			}

			if node._GetReadmitTime().Before(now) {
				this.healthyNodes = append(this.healthyNodes, node)
			}
		}
	}
}

func (this *_ManagedNetwork) _GetNumberOfNodesForTransaction() int { // nolint
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

func (this *_ManagedNetwork) _SetMinNodeReadmitPeriod(min time.Duration) {
	this.minNodeReadmitPeriod = min
	this.earliestReadmitTime = time.Now().Add(this.minNodeReadmitPeriod)
}

func (this *_ManagedNetwork) _GetMinNodeReadmitPeriod() time.Duration {
	return this.minNodeReadmitPeriod
}

func (this *_ManagedNetwork) _SetMaxNodeReadmitPeriod(max time.Duration) {
	this.maxNodeReadmitPeriod = max
}

func (this *_ManagedNetwork) _GetMaxNodeReadmitPeriod() time.Duration {
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
	this.healthyNodesMutex.RLock()
	defer this.healthyNodesMutex.RUnlock()

	if len(this.healthyNodes) == 0 {
		panic("failed to find a healthy working node")
	}

	bg := big.NewInt(int64(len(this.healthyNodes)))
	index, _ := rand.Int(rand.Reader, bg)
	return this.healthyNodes[index.Int64()]
}

func (this *_ManagedNetwork) _GetMinBackoff() time.Duration {
	return this.minBackoff
}

func (this *_ManagedNetwork) _SetMaxBackoff(maxBackoff time.Duration) {
	this.maxBackoff = maxBackoff
	for _, node := range this.healthyNodes {
		node._SetMaxBackoff(maxBackoff)
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
		if err := conn._Close(); err != nil {
			return err
		}
	}

	return nil
}

func _CreateNetworkFromNodes(nodes []_IManagedNode) (network map[string][]_IManagedNode, healthyNodes []_IManagedNode) {
	healthyNodes = []_IManagedNode{}
	network = map[string][]_IManagedNode{}

	for _, node := range nodes {
		if node._IsHealthy() {
			healthyNodes = append(healthyNodes, node)
		}

		value, ok := network[node._GetKey()]
		if !ok {
			value = []_IManagedNode{}
		}
		value = append(value, node)
		network[node._GetKey()] = value
	}

	return network, healthyNodes
}

func (this *_ManagedNetwork) _SetTransportSecurity(transportSecurity bool) (err error) {
	if this.transportSecurity != transportSecurity {
		if err := this._Close(); err != nil {
			return err
		}

		newNodes := make([]_IManagedNode, len(this.nodes))

		copy(newNodes, this.nodes)

		for i, node := range newNodes {
			if transportSecurity {
				newNodes[i] = node._ToSecure()
			} else {
				newNodes[i] = node._ToInsecure()
			}
		}

		newNetwork, newHealthyNodes := _CreateNetworkFromNodes(newNodes)

		this.nodes = newNodes
		this.healthyNodes = newHealthyNodes
		this.network = newNetwork
	}

	this.transportSecurity = transportSecurity
	return nil
}

func _GetNodesToRemove(network map[string]_IManagedNode, nodes []_IManagedNode) []int {
	nodeIndices := []int{}

	for i := len(nodes) - 1; i >= 0; i-- {
		if _, ok := network[nodes[i]._GetKey()]; !ok {
			nodeIndices = append(nodeIndices, i)
		}
	}

	return nodeIndices
}

func (this *_ManagedNetwork) _SetVerifyCertificate(verify bool) *_ManagedNetwork {
	for _, node := range this.nodes {
		node._SetVerifyCertificate(verify)
	}

	this.verifyCertificate = verify
	return this
}

func (this *_ManagedNetwork) _GetVerifyCertificate() bool {
	return this.verifyCertificate
}
