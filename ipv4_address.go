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

type _IPv4Address struct {
	network _IPv4AddressPart
	host    _IPv4AddressPart
}

func _Ipv4AddressFromProtobuf(byte []byte) _IPv4Address {
	return _IPv4Address{
		network: _IPv4AddressPart{
			left:  byte[0],
			right: byte[1],
		},
		host: _IPv4AddressPart{
			left:  byte[2],
			right: byte[3],
		},
	}
}

func (ip *_IPv4Address) _ToProtobuf() []byte {
	return []byte{ip.network.left, ip.network.right, ip.host.left, ip.host.right}
}

func (ip *_IPv4Address) String() string {
	return ip.network.String() + "." + ip.host.String()
}
