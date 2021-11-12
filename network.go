package hedera

import (
	"io/ioutil"
	"time"
)

type _Network struct {
	_ManagedNetwork
	addressBook map[AccountID]_NodeAddress
}

func _NewNetwork() _Network {
	return _Network{
		_ManagedNetwork: _NewManagedNetwork(),
		addressBook:     nil,
	}
}

func (network *_Network) SetNetwork(net map[string]AccountID) error {
	newNetwork := make(map[string]_IManagedNode)
	for url, id := range net {
		node := _NewNode(id, url, network._ManagedNetwork.nodeWaitTime.Milliseconds())
		newNetwork[id.String()] = &node
	}

	return network._ManagedNetwork._SetNetwork(newNetwork)
}

func (network *_Network) _GetNetwork() map[string]AccountID {
	temp := make(map[string]AccountID)
	for node, _ := range network._ManagedNetwork.nodes {
		switch n := node.(type) { //nolint
		case *_Node:
			temp[n._GetAddress()] = n.accountID
		}
	}

	return temp
}

func (network *_Network) _GetNodeForAccountID(id AccountID) (*_Node, bool) {
	for node, _ := range network._ManagedNetwork.nodes {
		switch n := node.(type) { //nolint
		case *_Node:
			if n.accountID.String() == id.String() {
				return n, true
			}
		}
	}

	return &_Node{}, false
}

func (network *_Network) _GetNetworkName() *NetworkName {
	return network._ManagedNetwork._GetNetworkName()
}

func (network *_Network) _SetNetworkName(net NetworkName) {
	network._ManagedNetwork._SetNetworkName(net)

	if network._ManagedNetwork.transportSecurity {
		network.addressBook = _ReadAddressBookResource("addressbook/" + net.String() + ".pb")

		if network.addressBook != nil {
			for nod, _ := range network._ManagedNetwork.nodes {
				switch n := nod.(type) { //nolint
				case *_Node:
					temp := network.addressBook[n.accountID]
					n.addressBook = &temp
				}
			}
		}
	}
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

func (network *_Network) _GetNodeAccountIDsForExecute() []AccountID {
	err := network._RemoveDeadNodes()
	if err != nil {
		panic(err)
	}

	length := network._ManagedNetwork._GetNumberOfNodesForTransaction()
	accountIDs := make([]AccountID, 0)

	i := 0
	for n, _ := range network._ManagedNetwork.nodes {
		switch nod := n.(type) { //nolint
		case *_Node:
			accountIDs = append(accountIDs, nod.accountID)
		}
		i++
		if i == length {
			break
		}
	}

	return accountIDs
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

func (network *_Network) _SetNodeWaitTime(waitTime time.Duration) {
	network._ManagedNetwork._SetNodeWaitTime(waitTime)
}

func (network *_Network) _GetNodeWaitTime() time.Duration {
	return network._ManagedNetwork._GetNodeWaitTime()
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

func (network *_Network) Close() error {
	err := network._ManagedNetwork._Close()
	if err != nil {
		return err
	}

	return nil
}
