package hedera

import (
	"sort"
)

type network struct {
	network      map[string]AccountID
	nodes        []*node
	networkNodes map[AccountID]*node
}

func newNetwork() network {
	return network{
		network:                  make(map[string]AccountID),
		nodes:                    make([]*node, 0),
		networkNodes:             make(map[AccountID]*node),
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
			node := newNode(id, url)
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

	return (count + 3 - 1) / 3
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
