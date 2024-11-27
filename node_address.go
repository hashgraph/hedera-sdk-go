package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// NodeAddress is the address of a node on the Hiero network
type NodeAddress struct {
	PublicKey   string
	AccountID   *AccountID
	NodeID      int64
	CertHash    []byte
	Addresses   []Endpoint
	Description string
}

func _NodeAddressFromProtobuf(nodeAd *services.NodeAddress) NodeAddress {
	address := make([]Endpoint, 0)

	for _, end := range nodeAd.GetServiceEndpoint() {
		address = append(address, EndpointFromProtobuf(end))
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
