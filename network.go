package hedera

import (
	"time"
)

type _Network struct {
	_ManagedNetwork
	addressBook map[AccountID]NodeAddress
}

func _NewNetwork() _Network {
	return _Network{
		_ManagedNetwork: _NewManagedNetwork(),
		addressBook:     nil,
	}
}

func (network *_Network) SetNetwork(net map[string]AccountID) (err error) {
	newNetwork := make(map[string]_IManagedNode)

	for url, id := range net {
		node, err := _NewNode(id, url, network.minBackoff)
		if err != nil {
			return err
		}
		newNetwork[url] = node
	}

	return network._ManagedNetwork._SetNetwork(newNetwork)
}

func (network *_Network) _GetNetwork() map[string]AccountID {
	temp := make(map[string]AccountID)
	for _, node := range network._ManagedNetwork.nodes {
		switch n := node.(type) { //nolint
		case *_Node:
			temp[n._GetAddress()] = n.accountID
		}
	}

	return temp
}

func (network *_Network) _IncreaseDelay(node *_Node) {
	node._IncreaseBackoff()

	for i, n := range network.healthyNodes {
		if goodNode, ok := n.(*_Node); ok {
			if goodNode.accountID.String() == node.accountID.String() {
				if i == len(network.healthyNodes)-1 {
					network.healthyNodes = network.healthyNodes[:i]
				} else {
					network.healthyNodes = append(network.healthyNodes[:i], network.healthyNodes[i+1:]...)
				}
			}
		}
	}
}

func (network *_Network) _GetNodeForAccountID(id AccountID) (*_Node, bool) {
	node, ok := network.network[id.String()]
	return node[0].(*_Node), ok
}

func (network *_Network) _GetNode() *_Node {
	return network._ManagedNetwork._GetNode().(*_Node)
}

func (network *_Network) _GetNetworkName() *NetworkName {
	if network._ManagedNetwork._GetLedgerID() != nil {
		temp, _ := network._ManagedNetwork._GetLedgerID().ToNetworkName()
		return &temp
	}

	return nil
}

func (network *_Network) _GetLedgerID() *LedgerID {
	if network._ManagedNetwork._GetLedgerID() != nil {
		return network._ManagedNetwork._GetLedgerID()
	}

	return &LedgerID{}
}

func (network *_Network) _SetLedgerID(id LedgerID) {
	network._ManagedNetwork._SetLedgerID(id)

	if network._ManagedNetwork.transportSecurity && network.ledgerID != nil {
		switch {
		case id.IsMainnet():
			network.addressBook = mainnetAddressBook._ToMap()
		case id.IsTestnet():
			network.addressBook = testnetAddressBook._ToMap()
		case id.IsPreviewnet():
			network.addressBook = previewnetAddressBook._ToMap()
		}

		if network.addressBook != nil {
			for _, node := range network._ManagedNetwork.nodes {
				if node, ok := node.(*_Node); ok {
					temp := network.addressBook[node.accountID]
					node.addressBook = &temp
				}
			}
			for _, nodes := range network._ManagedNetwork.network {
				for _, node := range nodes {
					if node, ok := node.(*_Node); ok {
						temp := network.addressBook[node.accountID]
						node.addressBook = &temp
					}
				}
			}
		}
	}
}

func (network *_Network) _SetNetworkName(net NetworkName) {
	ledger, err := LedgerIDFromNetworkName(net)
	if err != nil {
		panic(err)
	}

	network._SetLedgerID(*ledger)
}

func (network *_Network) _GetNodeAccountIDsForExecute() ([]AccountID, error) { //nolint
	nodes := make([]AccountID, 0)

	for i := 0; i < len(network.healthyNodes); i++ {
		node := network.healthyNodes[i]
		if node, ok := node.(*_Node); ok {
			nodes = append(nodes, node.accountID)
		}
	}

	return nodes, nil
}

func (network *_Network) _SetMaxNodesPerTransaction(max int) {
	network._ManagedNetwork._SetMaxNodesPerTransaction(max)
}

func (network *_Network) _SetMaxNodeAttempts(max int) {
	network._ManagedNetwork._SetMaxNodeAttempts(max)
}

func (network *_Network) _GetMaxNodeAttempts() int {
	return network._ManagedNetwork._GetMaxNodeAttempts()
}

func (network *_Network) _SetNodeMinBackoff(backoff time.Duration) {
	network._ManagedNetwork._SetMinBackoff(backoff)
}

func (network *_Network) _SetNodeMaxBackoff(backoff time.Duration) {
	network._ManagedNetwork._SetMaxBackoff(backoff)
}

func (network *_Network) _GetNodeMinBackoff() time.Duration {
	return network._ManagedNetwork._GetMinBackoff()
}

func (network *_Network) _GetNodeMaxBackoff() time.Duration {
	return network._ManagedNetwork._GetMaxBackoff()
}

func (network *_Network) _SetTransportSecurity(transportSecurity bool) *_Network {
	network._ManagedNetwork._SetTransportSecurity(transportSecurity)

	return network
}

func (network *_Network) _SetVerifyCertificate(verify bool) *_ManagedNetwork {
	return network._ManagedNetwork._SetVerifyCertificate(verify)
}

func (network *_Network) _GetVerifyCertificate() bool {
	return network._ManagedNetwork._GetVerifyCertificate()
}

func (network *_Network) _SetNodeMinReadmitPeriod(period time.Duration) {
	network._ManagedNetwork._SetMinNodeReadmitPeriod(period)
}

func (network *_Network) _SetNodeMaxReadmitPeriod(period time.Duration) {
	network._ManagedNetwork._SetMaxNodeReadmitPeriod(period)
}

func (network *_Network) _GetNodeMinReadmitPeriod() time.Duration {
	return network._ManagedNetwork.minNodeReadmitPeriod
}

func (network *_Network) _GetNodeMaxReadmitPeriod() time.Duration {
	return network._ManagedNetwork.maxNodeReadmitPeriod
}

func (network *_Network) Close() error {
	err := network._ManagedNetwork._Close()
	if err != nil {
		return err
	}

	return nil
}
