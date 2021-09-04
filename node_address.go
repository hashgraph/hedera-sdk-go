package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type _NodeAddress struct {
	publicKey   string
	accountID   *AccountID
	nodeID      int64
	certHash    []byte
	addresses   []_Endpoint
	description string
	stake       int64
}

func _NodeAddressFromProtobuf(nodeAd *proto.NodeAddress) _NodeAddress {
	address := make([]_Endpoint, 0)

	if len(nodeAd.GetIpAddress()) > 0 { // nolint
		address = append(address, _EndpointFromProtobuf(
			&proto.ServiceEndpoint{
				IpAddressV4: nodeAd.GetIpAddress(), // nolint
				Port:        nodeAd.GetPortno(),    // nolint
			}))
	}

	for _, end := range nodeAd.GetServiceEndpoint() {
		address = append(address, _EndpointFromProtobuf(end))
	}

	return _NodeAddress{
		publicKey:   nodeAd.GetRSA_PubKey(),
		accountID:   _AccountIDFromProtobuf(nodeAd.GetNodeAccountId()),
		nodeID:      nodeAd.GetNodeId(),
		certHash:    nodeAd.GetNodeCertHash(),
		addresses:   address,
		description: nodeAd.GetDescription(),
		stake:       nodeAd.GetStake(),
	}
}

func (nodeAdd *_NodeAddress) _ToProtobuf() *proto.NodeAddress {
	build := &proto.NodeAddress{
		RSA_PubKey:      nodeAdd.publicKey,
		NodeId:          nodeAdd.nodeID,
		NodeAccountId:   nil,
		NodeCertHash:    nodeAdd.certHash,
		ServiceEndpoint: nil,
		Description:     nodeAdd.description,
		Stake:           nodeAdd.stake,
	}

	if nodeAdd.accountID != nil {
		build.NodeAccountId = nodeAdd.accountID._ToProtobuf()
	}

	serviceEndpoint := make([]*proto.ServiceEndpoint, 0)
	for _, k := range nodeAdd.addresses {
		serviceEndpoint = append(serviceEndpoint, k._ToProtobuf())
	}
	build.ServiceEndpoint = serviceEndpoint

	return build
}

func (nodeAdd _NodeAddress) String() string {
	addresses := ""
	for _, k := range nodeAdd.addresses {
		addresses += k.String()
	}
	return nodeAdd.accountID.String() + " " + addresses + "\n" + "certHash " + string(nodeAdd.certHash)
}
