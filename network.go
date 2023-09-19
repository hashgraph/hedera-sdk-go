package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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

// SetNetwork sets the network to the given map of node addresses.
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

func (network *_Network) _IncreaseBackoff(node *_Node) {
	network.healthyNodesMutex.Lock()
	defer network.healthyNodesMutex.Unlock()
	node._IncreaseBackoff()

	index := -1
	for i, healthyNode := range network.healthyNodes {
		if node == healthyNode {
			index = i
			break
		}
	}

	if index >= 0 && index < len(network.healthyNodes) {
		network.healthyNodes = append(network.healthyNodes[:index], network.healthyNodes[index+1:]...)
	}
}

func (network *_Network) _GetNodeForAccountID(id AccountID) (*_Node, bool) {
	node, ok := network.network[id.String()]
	if !ok || node == nil {
		return nil, false
	}
	return node[0].(*_Node), ok
}

func (network *_Network) _GetNode() *_Node {
	return network._ManagedNetwork._GetNode().(*_Node)
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

func (network *_Network) _GetNodeAccountIDsForExecute() []AccountID { //nolint
	nodes := make([]AccountID, 0)
	nodesForTransaction := network._GetNumberOfNodesForTransaction()

	network.healthyNodesMutex.RLock()
	defer network.healthyNodesMutex.RUnlock()
	for i := 0; i < nodesForTransaction; i++ {
		nodes = append(nodes, network.healthyNodes[i].(*_Node).accountID)
	}

	return nodes
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
	_ = network._ManagedNetwork._SetTransportSecurity(transportSecurity)
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

// Close closes the network.
func (network *_Network) Close() error {
	err := network._ManagedNetwork._Close()
	if err != nil {
		return err
	}

	return nil
}

func _NetworkForMainnet(nodeAddresses map[AccountID]NodeAddress) *_Network {
	network := _NewNetwork()
	network.addressBook = nodeAddresses
	network._SetLedgerID(*NewLedgerIDMainnet())
	_ = network.SetNetwork(network._ToNet())
	return &network
}

func _NetworkForTestnet(nodeAddresses map[AccountID]NodeAddress) *_Network {
	network := _NewNetwork()
	network.addressBook = nodeAddresses
	network._SetLedgerID(*NewLedgerIDTestnet())
	_ = network.SetNetwork(network._ToNet())
	return &network
}

func _NetworkForPreviewnet(nodeAddresses map[AccountID]NodeAddress) *_Network {
	network := _NewNetwork()
	network.addressBook = nodeAddresses
	network._SetLedgerID(*NewLedgerIDPreviewnet())
	_ = network.SetNetwork(network._ToNet())
	return &network
}

func (network *_Network) _SetNetworkFromAddressBook(addressBook NodeAddressBook) {
	network.addressBook = addressBook._ToMap()
	_ = network.SetNetwork(network._ToNet())
}

func (network *_Network) _ToNet() map[string]AccountID {
	newNetwork := make(map[string]AccountID)
	for accountID, node := range network.addressBook {
		for _, address := range node.Addresses {
			newNetwork[address.String()] = accountID
		}
	}
	return newNetwork
}
