package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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

type IPv4Address struct {
	network IPv4AddressPart
	host    IPv4AddressPart
}

func Ipv4AddressFromProtobuf(byte []byte) IPv4Address {
	return IPv4Address{
		network: IPv4AddressPart{
			left:  byte[0],
			right: byte[1],
		},
		host: IPv4AddressPart{
			left:  byte[2],
			right: byte[3],
		},
	}
}

func (ip *IPv4Address) SetNetwork(left byte, right byte) *IPv4Address {
	ip.network.left = left
	ip.network.right = right
	return ip
}

func (ip *IPv4Address) SetHost(left byte, right byte) *IPv4Address {
	ip.host.left = left
	ip.host.right = right
	return ip
}

func (ip *IPv4Address) _ToProtobuf() []byte {
	return []byte{ip.network.left, ip.network.right, ip.host.left, ip.host.right}
}

func (ip *IPv4Address) String() string {
	return ip.network.String() + "." + ip.host.String()
}
