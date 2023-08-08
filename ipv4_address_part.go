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
)

type _IPv4AddressPart struct {
	left  byte
	right byte
}

func (ip *_IPv4AddressPart) SetLeft(left byte) *_IPv4AddressPart {
	ip.left = left
	return ip
}

func (ip *_IPv4AddressPart) SetRight(right byte) *_IPv4AddressPart {
	ip.right = right
	return ip
}

func (ip *_IPv4AddressPart) GetLeft() byte {
	return ip.left
}

func (ip *_IPv4AddressPart) GetRight() byte {
	return ip.right
}

func (ip *_IPv4AddressPart) String() string {
	return fmt.Sprintf("%d.%d", uint(ip.left), uint(ip.right))
}
