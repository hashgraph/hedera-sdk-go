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
	"regexp"
	"strconv"
)

var hostAndPort = regexp.MustCompile(`^(\S+):(\d+)$`)

type _ManagedNodeAddress struct {
	address *string
	port    uint32
}

func _ManagedNodeAddressFromString(str string) (*_ManagedNodeAddress, error) {
	hostAndPortMatch := hostAndPort.FindStringSubmatch(str)

	if len(hostAndPortMatch) > 1 {
		port, err := strconv.ParseUint(hostAndPortMatch[2], 10, 64)
		if err != nil {
			return nil, err
		}

		return &_ManagedNodeAddress{
			address: &hostAndPortMatch[1],
			port:    uint32(port),
		}, nil
	}

	return nil, fmt.Errorf("failed to parse node address")
}

func (address *_ManagedNodeAddress) _IsTransportSecurity() bool {
	return address.port == 50212 || address.port == 443
}

func (address *_ManagedNodeAddress) _ToInsecure() *_ManagedNodeAddress {
	port := address.port

	switch address.port {
	case 50212:
		port = 50211
	case 443:
		port = 5600
	}

	return &_ManagedNodeAddress{
		address: address.address,
		port:    port,
	}
}

func (address *_ManagedNodeAddress) _ToSecure() *_ManagedNodeAddress {
	port := address.port

	switch port {
	case 50211:
		return &_ManagedNodeAddress{
			address: address.address,
			port:    50212,
		}
	case 5600:
		return &_ManagedNodeAddress{
			address: address.address,
			port:    443,
		}
	}

	return &_ManagedNodeAddress{
		address: address.address,
		port:    port,
	}
}

func (address *_ManagedNodeAddress) _Equals(comp _ManagedNodeAddress) bool { //nolint
	if address.address != nil && address.address == comp.address {
		if address.port == comp.port {
			return true
		}
	}

	return false
}

func (address *_ManagedNodeAddress) _String() string {
	if address.address != nil {
		return *address.address + ":" + strconv.FormatInt(int64(address.port), 10)
	}

	return ""
}
