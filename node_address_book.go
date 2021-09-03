package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type _NodeAddressBook struct {
	nodeAddresses []_NodeAddress
}

func _NodeAddressBookFromProtobuf(book *proto.NodeAddressBook) _NodeAddressBook {
	addresses := make([]_NodeAddress, 0)

	for _, k := range book.NodeAddress {
		addresses = append(addresses, _NodeAddressFromProtobuf(k))
	}

	return _NodeAddressBook{
		nodeAddresses: addresses,
	}
}

func (book _NodeAddressBook) _ToProtobuf() *proto.NodeAddressBook {
	addresses := make([]*proto.NodeAddress, 0)

	for _, k := range book.nodeAddresses {
		addresses = append(addresses, k._ToProtobuf())
	}

	return &proto.NodeAddressBook{
		NodeAddress: addresses,
	}
}

func (book _NodeAddressBook) ToBytes() []byte {
	data, err := protobuf.Marshal(book._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func _NodeAddressBookFromBytes(data []byte) (_NodeAddressBook, error) {
	if data == nil {
		return _NodeAddressBook{}, errByteArrayNull
	}
	pb := proto.NodeAddressBook{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return _NodeAddressBook{}, err
	}

	derivedBytes := _NodeAddressBookFromProtobuf(&pb)

	return derivedBytes, nil
}
