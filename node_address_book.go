package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type _NodeAddressBook struct {
	nodeAddresses []_NodeAddress
}

func nodeAddressBookFromProtobuf(book *proto.NodeAddressBook) _NodeAddressBook {
	addresses := make([]_NodeAddress, 0)

	for _, k := range book.NodeAddress {
		addresses = append(addresses, nodeAddressFromProtobuf(k))
	}

	return _NodeAddressBook{
		nodeAddresses: addresses,
	}
}

func (book _NodeAddressBook) toProtobuf() *proto.NodeAddressBook {
	addresses := make([]*proto.NodeAddress, 0)

	for _, k := range book.nodeAddresses {
		addresses = append(addresses, k.toProtobuf())
	}

	return &proto.NodeAddressBook{
		NodeAddress: addresses,
	}
}

func (book _NodeAddressBook) ToBytes() []byte {
	data, err := protobuf.Marshal(book.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func nodeAddressBookFromBytes(data []byte) (_NodeAddressBook, error) {
	if data == nil {
		return _NodeAddressBook{}, errByteArrayNull
	}
	pb := proto.NodeAddressBook{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return _NodeAddressBook{}, err
	}

	derivedBytes := nodeAddressBookFromProtobuf(&pb)

	return derivedBytes, nil
}
