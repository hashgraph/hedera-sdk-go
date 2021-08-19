package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type nodeAddressBook struct {
	nodeAddresses []nodeAddress
}

func nodeAddressBookFromProtobuf(book *proto.NodeAddressBook) nodeAddressBook {
	addresses := make([]nodeAddress, 0)

	for _, k := range book.NodeAddress {
		addresses = append(addresses, nodeAddressFromProtobuf(k))
	}

	return nodeAddressBook{
		nodeAddresses: addresses,
	}
}

func (book nodeAddressBook) toProtobuf() *proto.NodeAddressBook {
	addresses := make([]*proto.NodeAddress, 0)

	for _, k := range book.nodeAddresses {
		addresses = append(addresses, k.toProtobuf())
	}

	return &proto.NodeAddressBook{
		NodeAddress: addresses,
	}
}

func (book nodeAddressBook) ToBytes() []byte {
	data, err := protobuf.Marshal(book.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func nodeAddressBookFromBytes(data []byte) (nodeAddressBook, error) {
	if data == nil {
		return nodeAddressBook{}, errByteArrayNull
	}
	pb := proto.NodeAddressBook{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return nodeAddressBook{}, err
	}

	derivedBytes := nodeAddressBookFromProtobuf(&pb)

	return derivedBytes, nil
}
