package hedera

import (
	"sort"
	"time"
)

type network struct {
	network      map[string]AccountID
	nodes        []node
	networkNodes map[AccountID]node

	lastSortedNodeAccountIDs int64
}

func newNetwork() network {
	networkForReturn := network{
		network:                  make(map[string]AccountID),
		nodes:                    make([]node, 0),
		networkNodes:             make(map[AccountID]node),
		lastSortedNodeAccountIDs: time.Now().UTC().Unix(),
	}

	return networkForReturn
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
			network.networkNodes[id] = newNode(id, url)
		}
	}

	network.nodes = make([]node, len(net))
	i := 0
	for _, node := range network.networkNodes {
		network.nodes[i] = node
		i++
	}

	network.network = net

	return nil
}

func (network *network) getNodeAccountIDsForExecute() []AccountID {
	if network.lastSortedNodeAccountIDs+1 < time.Now().UTC().Unix() {
		sort.Sort(nodes{nodes: network.nodes})
		network.lastSortedNodeAccountIDs = time.Now().UTC().Unix()
	}

	slice := network.nodes[0:network.getNumberOfNodesForTransaction()]
	var accountIDs = make([]AccountID, len(slice))
	for i, id := range slice {
		accountIDs[i] = id.accountID
	}

	return accountIDs
}

func (network *network) getNumberOfNodesForTransaction() int {
	return (len(network.nodes) + 3 - 1) / 3
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
