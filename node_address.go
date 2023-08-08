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
	"github.com/hashgraph/hedera-protobufs-go/services"
)

// NodeAddress is the address of a node on the Hedera network
type NodeAddress struct {
	PublicKey   string
	AccountID   *AccountID
	NodeID      int64
	CertHash    []byte
	Addresses   []_Endpoint
	Description string
}

func _NodeAddressFromProtobuf(nodeAd *services.NodeAddress) NodeAddress {
	address := make([]_Endpoint, 0)

	for _, end := range nodeAd.GetServiceEndpoint() {
		address = append(address, _EndpointFromProtobuf(end))
	}

	return NodeAddress{
		PublicKey:   nodeAd.GetRSA_PubKey(),
		AccountID:   _AccountIDFromProtobuf(nodeAd.GetNodeAccountId()),
		NodeID:      nodeAd.GetNodeId(),
		CertHash:    nodeAd.GetNodeCertHash(),
		Addresses:   address,
		Description: nodeAd.GetDescription(),
	}
}

func (nodeAdd *NodeAddress) _ToProtobuf() *services.NodeAddress {
	build := &services.NodeAddress{
		RSA_PubKey:      nodeAdd.PublicKey,
		NodeId:          nodeAdd.NodeID,
		NodeAccountId:   nil,
		NodeCertHash:    nodeAdd.CertHash,
		ServiceEndpoint: nil,
		Description:     nodeAdd.Description,
	}

	if nodeAdd.AccountID != nil {
		build.NodeAccountId = nodeAdd.AccountID._ToProtobuf()
	}

	serviceEndpoint := make([]*services.ServiceEndpoint, 0)
	for _, k := range nodeAdd.Addresses {
		serviceEndpoint = append(serviceEndpoint, k._ToProtobuf())
	}
	build.ServiceEndpoint = serviceEndpoint

	return build
}

// String returns a string representation of the NodeAddress
func (nodeAdd NodeAddress) String() string {
	Addresses := ""
	for _, k := range nodeAdd.Addresses {
		Addresses += k.String()
	}
	return nodeAdd.AccountID.String() + " " + Addresses + "\n" + "CertHash " + string(nodeAdd.CertHash)
}
