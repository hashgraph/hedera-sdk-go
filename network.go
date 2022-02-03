package hedera

import (
	"io/ioutil"
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

func (network *_Network) SetNetwork(net map[string]AccountID) error {
	newNetwork := make(map[string]_IManagedNode)

	for url, id := range net {
		node := _NewNode(id, url, network._ManagedNetwork.minBackOff.Milliseconds())
		newNetwork[url] = &node
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

func (network *_Network) _GetNodeForAccountID(id AccountID) (*_Node, bool) {
	for _, node := range network.network[id.String()] {
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

	if network._ManagedNetwork.transportSecurity {
		network.addressBook = _ReadAddressBookResource("addressbook/" + id.String() + ".pb")

		if network.addressBook != nil {
			for _, nod := range network._ManagedNetwork.nodes {
				switch n := nod.(type) { //nolint
				case *_Node:
					temp := network.addressBook[n.accountID]
					n.addressBook = &temp
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

func _ReadAddressBookResource(ad string) map[AccountID]NodeAddress {
	f, err := ioutil.ReadFile(ad)
	if err != nil {
		panic(err)
	}

	nodeAB, err := NodeAddressBookFromBytes(f)
	if err != nil {
		panic(err)
	}

	resultMap := make(map[AccountID]NodeAddress)
	for _, nodeAd := range nodeAB.NodeAddresses {
		if nodeAd.AccountID == nil {
			continue
		}

		resultMap[*nodeAd.AccountID] = nodeAd
	}

	return resultMap
}

func (network *_Network) _GetNodeAccountIDsForExecute() ([]AccountID, error) { //nolint
	nodes := network._GetNumberOfMostHealthyNodes(int32(network._ManagedNetwork._GetNumberOfNodesForTransaction()))
	accountIDs := make([]AccountID, 0)

	for _, node := range nodes {
		switch n := node.(type) { //nolint
		case *_Node:
			accountIDs = append(accountIDs, n.accountID)
		}
	}

	return accountIDs, nil
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

func (network *_Network) _SetNodeMinBackoff(waitTime time.Duration) {
	network._ManagedNetwork._SetMinBackoff(waitTime)
}

func (network *_Network) _SetNodeMaxBackoff(waitTime time.Duration) {
	network._ManagedNetwork._SetMaxBackoff(waitTime)
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

func (network *_Network) Close() error {
	err := network._ManagedNetwork._Close()
	if err != nil {
		return err
	}

	return nil
}
