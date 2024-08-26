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
	address []byte
}

// Leaving this due to backwards compatibility.
func Ipv4AddressFromProtobuf(address []byte) IPv4Address {
	return Ipv4AddressFromBytes(address)
}

func Ipv4AddressFromBytes(address []byte) IPv4Address {
	return IPv4Address{
		address: address,
	}
}

func (ip *IPv4Address) _ToProtobuf() []byte {
	return ip.address
}

func (ip *IPv4Address) String() string {
	return string(ip.address)
}
