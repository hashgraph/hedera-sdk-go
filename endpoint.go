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
