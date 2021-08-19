package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type endpoint struct {
	address ipv4Address
	port    int32
}

func endpointFromProtobuf(serviceEndpoint *proto.ServiceEndpoint) endpoint {
	port := serviceEndpoint.GetPort()

	if port == 0 || port == 50111 {
		port = 50211
	}

	return endpoint{
		address: ipv4AddressFromProtobuf(serviceEndpoint.GetIpAddressV4()),
		port:    port,
	}
}

func (endpoint *endpoint) toProtobuf() *proto.ServiceEndpoint {
	return &proto.ServiceEndpoint{
		IpAddressV4: endpoint.address.toProtobuf(),
		Port:        endpoint.port,
	}
}

func (endpoint *endpoint) String() string {
	return endpoint.address.String() + ":" + fmt.Sprintf("%d", endpoint.port)
}
