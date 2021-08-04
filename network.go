package hedera

import (
	"math"
	"sort"
	"time"
)

type network struct {
	network                map[string]AccountID
	nodes                  []*node
	networkNodes           map[AccountID]*node
	maxNodeAttempts        int
	nodeWaitTime           time.Duration
	maxNodesPerTransaction *int
}

func newNetwork() network {
	return network{
		network:                make(map[string]AccountID),
		nodes:                  make([]*node, 0),
		networkNodes:           make(map[AccountID]*node),
		maxNodeAttempts:        -1,
		nodeWaitTime:           250 * time.Millisecond,
		maxNodesPerTransaction: nil,
	}
}

func (network *network) SetNetwork(net map[string]AccountID) error {
	for url, id := range network.network {
		if _, ok := net[url]; !ok {
			err := network.networkNodes[id].close()
			if err != nil {
				return err
			}

			delete(network.networkNodes, id)
		}
	}

	for url, id := range net {
		if _, ok := network.network[url]; !ok {
			node := newNode(id, url, network.nodeWaitTime.Milliseconds())
			network.networkNodes[id] = &node
		}
	}

	network.nodes = make([]*node, len(net))
	i := 0
	for _, node := range network.networkNodes {
		network.nodes[i] = node
		i++
	}

	network.network = net

	return nil
}

func (network *network) getNodeAccountIDsForExecute() []AccountID {
	sort.Sort(nodes{nodes: network.nodes})

	if network.maxNodeAttempts > 0 {
		for i := 0; i < len(network.nodes); i++ {
			var nod *node
			if network.nodes[i] != nil {
				nod = network.nodes[i]
			} else {
				continue
			}
			if nod.attempts >= int64(network.maxNodeAttempts) {
				err := nod.close()
				if err != nil {
					panic(err)
				}
				network.nodes = append(network.nodes[:i], network.nodes[i+1:]...)
				delete(network.network, nod.address)
				delete(network.networkNodes, nod.accountID)
				i--
			}
		}
	}

	length := network.getNumberOfNodesForTransaction()
	accountIDs := make([]AccountID, length)

	for i, id := range network.nodes[0:length] {
		accountIDs[i] = id.accountID
	}

	return accountIDs
}

func (network *network) getNumberOfNodesForTransaction() int {
	count := 0
	for _, node := range network.nodes {
		if node.isHealthy() {
			count++
		}
	}

	if network.maxNodesPerTransaction != nil {
		return int(math.Min(float64(*network.maxNodesPerTransaction), float64(count)))
	}

	return (count + 3 - 1) / 3
}

func (network *network) setMaxNodesPerTransaction(max int) {
	network.maxNodesPerTransaction = &max
}

func (network *network) setMaxNodeAttempts(max int) {
	network.maxNodeAttempts = max
}

func (network *network) getMaxNodeAttempts() int {
	return network.maxNodeAttempts
}

func (network *network) setNodeWaitTime(waitTime time.Duration) {
	network.nodeWaitTime = waitTime
	for _, nod := range network.nodes {
		if nod != nil {
			nod.setWaitTime(waitTime.Milliseconds())
		}
	}
}

func (network *network) getNodeWaitTime() time.Duration {
	return network.nodeWaitTime
}

func (network *network) Close() error {
	for _, conn := range network.nodes {
		if conn.channel != nil {
			err := conn.channel.client.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
