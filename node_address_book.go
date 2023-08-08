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
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// NodeAddressBook is the address book for the nodes on the Hedera network
type NodeAddressBook struct {
	NodeAddresses []NodeAddress
}

func _NodeAddressBookFromProtobuf(book *services.NodeAddressBook) NodeAddressBook {
	addresses := make([]NodeAddress, 0)

	for _, k := range book.NodeAddress {
		addresses = append(addresses, _NodeAddressFromProtobuf(k))
	}

	return NodeAddressBook{
		NodeAddresses: addresses,
	}
}

func (book NodeAddressBook) _ToProtobuf() *services.NodeAddressBook {
	addresses := make([]*services.NodeAddress, 0)

	for _, k := range book.NodeAddresses {
		addresses = append(addresses, k._ToProtobuf())
	}

	return &services.NodeAddressBook{
		NodeAddress: addresses,
	}
}

// ToBytes returns the byte representation of the NodeAddressBook
func (book NodeAddressBook) ToBytes() []byte {
	data, err := protobuf.Marshal(book._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func (book NodeAddressBook) _ToMap() (result map[AccountID]NodeAddress) {
	result = map[AccountID]NodeAddress{}

	for _, node := range book.NodeAddresses {
		if node.AccountID == nil {
			continue
		}

		result[*node.AccountID] = node
	}

	return result
}

// NodeAddressBookFromBytes returns the NodeAddressBook from a raw byte array
func NodeAddressBookFromBytes(data []byte) (NodeAddressBook, error) {
	if data == nil {
		return NodeAddressBook{}, errByteArrayNull
	}
	pb := services.NodeAddressBook{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return NodeAddressBook{}, err
	}

	derivedBytes := _NodeAddressBookFromProtobuf(&pb)

	return derivedBytes, nil
}
