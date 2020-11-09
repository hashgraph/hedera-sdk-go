package hedera

import (
	"google.golang.org/grpc"
	"sort"
	"time"
)

type network struct {
	networkNodeIds []node
	network        map[AccountID]node

	lastSortedNodeAccountIDs int64
	nextNodeIndex            uint
}

func newNetwork(net map[string]AccountID) network {
	newNetwork := make(map[AccountID]node, len(net))
	for node, accountID := range net {
		newNetwork[accountID] = newNode(accountID, node)
	}

	networkForReturn := network{
		networkNodeIds:           make([]node, 0),
		network:                  newNetwork,
		lastSortedNodeAccountIDs: time.Now().UTC().UnixNano(),
	}

	return networkForReturn
}

func (network *network) SetNetwork(net map[string]AccountID) *network {
	for node, id := range net {
		network.networkNodeIds = append(network.networkNodeIds, newNode(id, node))
		network.network[id] = newNode(id, node)
	}

	return network
}

func (network *network) getNodeAccountIDsForExecute() []AccountID {
	if network.lastSortedNodeAccountIDs+1000 < time.Now().UTC().UnixNano() {
		sort.Sort(nodes{nodes: network.networkNodeIds})
		network.lastSortedNodeAccountIDs = time.Now().UTC().UnixNano()
	}

	slice := network.networkNodeIds[0:network.getNumberOfNodesForTransaction()]
	var accountIDs = make([]AccountID, len(slice))
	for i, id := range slice {
		accountIDs[i] = id.accountID
	}

	return accountIDs
}

func (network *network) getNumberOfNodesForTransaction() int {
	return (len(network.networkNodeIds) + 3 - 1) / 3
}

func (network *network) getNextNode() AccountID {
	nodeID := network.networkNodeIds[network.nextNodeIndex]
	network.nextNodeIndex = (network.nextNodeIndex + 1) % uint(len(network.networkNodeIds))

	return nodeID.accountID
}

func (network *network) getChannel(id AccountID) (*channel, error) {
	var i int
	var node node
	for i, node = range network.networkNodeIds {
		if node.accountID == id && node.channel != nil {
			return network.networkNodeIds[i].channel, nil
		}
	}

	conn, err := grpc.Dial(network.network[id].address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ch := newChannel(conn)
	network.networkNodeIds[i].channel = &ch
	return &ch, nil
}

func (network *network) Close() error {
	for _, conn := range network.networkNodeIds {
		if conn.channel != nil {
			err := conn.channel.client.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
