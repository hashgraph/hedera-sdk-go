package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type nodeAddress struct {
	publicKey   string
	accountID   *AccountID
	nodeID      int64
	certHash    []byte
	addresses   []endpoint
	description string
	stake       int64
}

func nodeAddressFromProtobuf(nodeAd *proto.NodeAddress) nodeAddress {
	address := make([]endpoint, 0)

	if len(nodeAd.GetIpAddress()) > 0 {
		address = append(address, endpointFromProtobuf(
			&proto.ServiceEndpoint{
				IpAddressV4: nodeAd.GetIpAddress(),
				Port:        nodeAd.GetPortno(),
			}))
	}

	for _, end := range nodeAd.GetServiceEndpoint() {
		address = append(address, endpointFromProtobuf(end))
	}

	var account AccountID
	if nodeAd.GetNodeAccountId() != nil {
		account = accountIDFromProtobuf(nodeAd.GetNodeAccountId())
	}

	return nodeAddress{
		publicKey:   nodeAd.GetRSA_PubKey(),
		accountID:   &account,
		nodeID:      nodeAd.GetNodeId(),
		certHash:    nodeAd.GetNodeCertHash(),
		addresses:   address,
		description: nodeAd.GetDescription(),
		stake:       nodeAd.GetStake(),
	}
}

func (nodeAdd *nodeAddress) toProtobuf() *proto.NodeAddress {
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
		build.NodeAccountId = nodeAdd.accountID.toProtobuf()
	}

	serviceEndpoint := make([]*proto.ServiceEndpoint, 0)
	for _, k := range nodeAdd.addresses {
		serviceEndpoint = append(serviceEndpoint, k.toProtobuf())
	}
	build.ServiceEndpoint = serviceEndpoint

	return build
}

func (nodeAdd nodeAddress) String() string {
	addresses := ""
	for _, k := range nodeAdd.addresses {
		addresses = addresses + k.String()
	}
	return nodeAdd.accountID.String() + " " + addresses + "\n" + "certHash " + string(nodeAdd.certHash)
}
