package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

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

func (book NodeAddressBook) ToBytes() []byte {
	data, err := protobuf.Marshal(book._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

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
