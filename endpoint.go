package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type _Endpoint struct {
	address _IPv4Address
	port    int32
}

func _EndpointFromProtobuf(serviceEndpoint *services.ServiceEndpoint) _Endpoint {
	port := serviceEndpoint.GetPort()

	if port == 0 || port == 50111 {
		port = 50211
	}

	return _Endpoint{
		address: _Ipv4AddressFromProtobuf(serviceEndpoint.GetIpAddressV4()),
		port:    port,
	}
}

func (endpoint *_Endpoint) _ToProtobuf() *services.ServiceEndpoint {
	return &services.ServiceEndpoint{
		IpAddressV4: endpoint.address._ToProtobuf(),
		Port:        endpoint.port,
	}
}

func (endpoint *_Endpoint) String() string {
	return endpoint.address.String() + ":" + fmt.Sprintf("%d", endpoint.port)
}
