package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type _Endpoint struct {
	address _IPv4Address
	port    int32
}

func endpointFromProtobuf(serviceEndpoint *proto.ServiceEndpoint) _Endpoint {
	port := serviceEndpoint.GetPort()

	if port == 0 || port == 50111 {
		port = 50211
	}

	return _Endpoint{
		address: ipv4AddressFromProtobuf(serviceEndpoint.GetIpAddressV4()),
		port:    port,
	}
}

func (endpoint *_Endpoint) toProtobuf() *proto.ServiceEndpoint {
	return &proto.ServiceEndpoint{
		IpAddressV4: endpoint.address.toProtobuf(),
		Port:        endpoint.port,
	}
}

func (endpoint *_Endpoint) String() string {
	return endpoint.address.String() + ":" + fmt.Sprintf("%d", endpoint.port)
}
