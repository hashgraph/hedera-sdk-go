package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type NodeAddress struct {
	PublicKey   string
	AccountID   *AccountID
	NodeID      int64
	CertHash    []byte
	Addresses   []_Endpoint
	Description string
	Stake       int64
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
		Stake:       nodeAd.GetStake(),
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
		Stake:           nodeAdd.Stake,
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

func (nodeAdd NodeAddress) String() string {
	Addresses := ""
	for _, k := range nodeAdd.Addresses {
		Addresses += k.String()
	}
	return nodeAdd.AccountID.String() + " " + Addresses + "\n" + "CertHash " + string(nodeAdd.CertHash)
}
