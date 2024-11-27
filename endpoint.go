package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type Endpoint struct {
	address    []byte
	port       int32
	domainName string
}

func (endpoint *Endpoint) SetAddress(address []byte) *Endpoint {
	endpoint.address = address
	return endpoint
}

func (endpoint *Endpoint) GetAddress() []byte {
	return endpoint.address
}

func (endpoint *Endpoint) SetPort(port int32) *Endpoint {
	endpoint.port = port
	return endpoint
}

func (endpoint *Endpoint) GetPort() int32 {
	return endpoint.port
}

func (endpoint *Endpoint) SetDomainName(domainName string) *Endpoint {
	endpoint.domainName = domainName
	return endpoint
}

func (endpoint *Endpoint) GetDomainName() string {
	return endpoint.domainName
}

func EndpointFromProtobuf(serviceEndpoint *services.ServiceEndpoint) Endpoint {
	port := serviceEndpoint.GetPort()

	if port == 0 || port == 50111 {
		port = 50211
	}

	return Endpoint{
		address:    serviceEndpoint.GetIpAddressV4(),
		port:       port,
		domainName: serviceEndpoint.GetDomainName(),
	}
}

func (endpoint *Endpoint) _ToProtobuf() *services.ServiceEndpoint {
	return &services.ServiceEndpoint{
		IpAddressV4: endpoint.address,
		Port:        endpoint.port,
		DomainName:  endpoint.domainName,
	}
}

func (endpoint *Endpoint) String() string {
	if endpoint.domainName != "" {
		// If domain name is populated domainName + port
		return endpoint.domainName + ":" + fmt.Sprintf("%d", endpoint.port)
	} else {
		return fmt.Sprintf("%d.%d.%d.%d:%d",
			int(endpoint.address[0])&0xFF,
			int(endpoint.address[1])&0xFF,
			int(endpoint.address[2])&0xFF,
			int(endpoint.address[3])&0xFF,
			endpoint.port)
	}
}
