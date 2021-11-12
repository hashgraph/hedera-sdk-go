package hedera

import (
	"math"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type _ManagedNetwork struct {
	network                map[string]_IManagedNode
	nodes                  map[_IManagedNode][]_ManagedNode
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
		nodes:                  make(map[_IManagedNode][]_ManagedNode),
		maxNodeAttempts:        -1,
		nodeWaitTime:           250 * time.Millisecond,
		maxNodesPerTransaction: nil,
		transportSecurity:      false,
	}
}

func (network *_ManagedNetwork) _SetNetwork(net map[string]_IManagedNode) error {
	if len(network.nodes) == 0 {
		for url, originalNode := range net {
			network.nodes[originalNode] = append(network.nodes[originalNode], _NewManagedNode(url, network.nodeWaitTime.Milliseconds()))
			network.network[url] = originalNode
		}

		return nil
	}

	for url, node := range network._GetNodesToRemove(net) {
		for i, _ := range network.nodes {
			if i._GetAddress() == node._GetAddress() {
				delete(network.nodes, i)
			}
		}
		_ = node._Close()
		delete(network.network, url)
	}

	temp := make(map[_IManagedNode][]_ManagedNode)

	for url, node := range net {
		temp[node] = append(temp[node], _NewManagedNode(url, network.nodeWaitTime.Milliseconds()))
	}

	for node, urlArr := range temp {
		_, ok := network.nodes[node]
		if ok {
			tempArr := make([]_ManagedNode, 0)
			for i := 0; i < len(urlArr); i++ {
				for j := 0; j < len(network.nodes[node]); j++ {
					if urlArr[i]._GetAddress() == network.nodes[node][j]._GetAddress() {
						tempArr = append(tempArr, urlArr[i])
					}
				}
			}
			network.nodes[node] = tempArr
		} else {
			network.nodes[node] = urlArr
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
	for nod, _ := range network.nodes {
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
	for conn, _ := range network.nodes {
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

		for node, urls := range network.nodes {
			if transportSecurity {
				node = node._ToSecure()
			} else {
				node = node._ToInsecure()
			}
			network.nodes[node] = urls
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
	temp := make([]_IManagedNode, 0)
	for n, _ := range network.nodes {
		temp = append(temp, n)
	}
	sort.Sort(_Nodes{nodes: temp})
	err := network._RemoveDeadNodes()
	if err != nil {
		return []_IManagedNode{}
	}

	nodes := make([]_IManagedNode, 0)

	size := math.Min(float64(count), float64(len(network.nodes)))
	i := 0
	for n, _ := range network.nodes {
		nodes = append(nodes, n)
		i++
		if float64(i) == size {
			break
		}
	}

	return nodes
}

func (network *_ManagedNetwork) _RemoveDeadNodes() error {
	if network.maxNodeAttempts > 0 {
		for n, _ := range network.nodes {
			if n == nil {
				return errors.New("null nodes can't be removed")
			}

			if n._GetAttempts() >= int64(network.maxNodeAttempts) {
				err := n._Close()
				if err != nil {
					return err
				}
				network._RemoveNodeFromNetwork(n)
				delete(network.nodes, n)
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
		switch nodeOut := node.(type) {
		case *_Node:
			b = false
			for _, n := range net {
				switch nodeIn := n.(type) {
				case *_Node:
					if nodeIn.accountID.String() == nodeOut.accountID.String() {
						b = true
					}
				}
			}
			if !b {
				nodes[url] = node
			}
		}
	}

	return nodes
}

func (network *_ManagedNetwork) _SetVerifyCertificate(verify bool) *_ManagedNetwork {
	if len(network.nodes) > 0 {
		for node, urlArray := range network.nodes {
			switch s := node.(type) { //nolint
			case *_Node:
				s._SetCertificateVerification(verify)
				network.nodes[s] = urlArray
			}
		}
	}

	network.verifyCertificate = verify

	return network
}

func (network *_ManagedNetwork) _GetVerifyCertificate() bool {
	if len(network.nodes) > 0 {
		for node, _ := range network.nodes {
			switch s := node.(type) { //nolint
			case *_Node:
				return s._GetCertificateVerification()
			}
		}
	}

	return false
}
