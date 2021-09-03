package hedera

import (
	"io/ioutil"
	"math"
	"sort"
	"time"
)

type _Network struct {
	network                map[string]AccountID
	nodes                  []*_Node
	networkNodes           map[AccountID]*_Node
	maxNodeAttempts        int
	nodeWaitTime           time.Duration
	maxNodesPerTransaction *int
	addressBook            map[AccountID]_NodeAddress
	networkName            *NetworkName
}

func _NewNetwork() _Network {
	return _Network{
		network:                make(map[string]AccountID),
		nodes:                  make([]*_Node, 0),
		networkNodes:           make(map[AccountID]*_Node),
		maxNodeAttempts:        -1,
		nodeWaitTime:           250 * time.Millisecond,
		maxNodesPerTransaction: nil,
		addressBook:            nil,
	}
}

func (network *_Network) SetNetwork(net map[string]AccountID) error {
	for url, id := range network.network {
		if _, ok := net[url]; !ok {
			err := network.networkNodes[id]._Close()
			if err != nil {
				return err
			}

			delete(network.networkNodes, id)
		}
	}

	for url, id := range net {
		if _, ok := network.network[url]; !ok {
			node := _NewNode(id, url, network.nodeWaitTime.Milliseconds())
			network.networkNodes[id] = &node
		}
	}

	network.nodes = make([]*_Node, len(net))
	i := 0
	for _, node := range network.networkNodes {
		network.nodes[i] = node
		i++
	}

	network.network = net

	return nil
}

func (network *_Network) _GetNodeAccountIDsForExecute() []AccountID {
	sort.Sort(_Nodes{nodes: network.nodes})

	if network.maxNodeAttempts > 0 {
		for i := 0; i < len(network.nodes); i++ {
			var nod *_Node
			if network.nodes[i] != nil {
				nod = network.nodes[i]
			} else {
				continue
			}
			if nod.attempts >= int64(network.maxNodeAttempts) {
				err := nod._Close()
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

	length := network._GetNumberOfNodesForTransaction()
	accountIDs := make([]AccountID, length)

	for i, id := range network.nodes[0:length] {
		accountIDs[i] = id.accountID
	}

	return accountIDs
}

func (network *_Network) _GetNetworkName() *NetworkName {
	return network.networkName
}

func (network *_Network) _SetNetworkName(net NetworkName) *_Network {
	network.networkName = &net

	switch net {
	case NetworkNameMainnet:
		network.addressBook = _ReadAddressBookResource("addressbook/mainnet.pb")
	case NetworkNameTestnet:
		network.addressBook = _ReadAddressBookResource("addressbook/testnet.pb")
	case NetworkNamePreviewnet:
		network.addressBook = _ReadAddressBookResource("addressbook/previewnet.pb")
	}

	if network.addressBook != nil {
		for _, nod := range network.nodes {
			temp := network.addressBook[nod.accountID]
			nod.addressBook = &temp
		}
	}

	return network
}

func _ReadAddressBookResource(ad string) map[AccountID]_NodeAddress {
	f, err := ioutil.ReadFile(ad)
	if err != nil {
		panic(err)
	}

	nodeAB, err := _NodeAddressBookFromBytes(f)
	if err != nil {
		panic(err)
	}

	resultMap := make(map[AccountID]_NodeAddress)
	for _, nodeAd := range nodeAB.nodeAddresses {
		if nodeAd.accountID == nil {
			continue
		}

		resultMap[*nodeAd.accountID] = nodeAd
	}

	return resultMap
}

func (network *_Network) _GetNumberOfNodesForTransaction() int {
	count := 0
	for _, node := range network.nodes {
		if node._IsHealthy() {
			count++
		}
	}

	if network.maxNodesPerTransaction != nil {
		return int(math.Min(float64(*network.maxNodesPerTransaction), float64(count)))
	}

	return (count + 3 - 1) / 3
}

func (network *_Network) _SetMaxNodesPerTransaction(max int) {
	network.maxNodesPerTransaction = &max
}

func (network *_Network) _SetMaxNodeAttempts(max int) {
	network.maxNodeAttempts = max
}

func (network *_Network) _GetMaxNodeAttempts() int {
	return network.maxNodeAttempts
}

func (network *_Network) _SetNodeWaitTime(waitTime time.Duration) {
	network.nodeWaitTime = waitTime
	for _, nod := range network.nodes {
		if nod != nil {
			nod._SetWaitTime(waitTime.Milliseconds())
		}
	}
}

func (network *_Network) _GetNodeWaitTime() time.Duration {
	return network.nodeWaitTime
}

func (network *_Network) Close() error {
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
