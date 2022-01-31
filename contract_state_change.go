package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type ContractStateChange struct {
	ContractID     *ContractID
	StorageChanges []*StorageChange
}

func _ContractStateChangeFromProtobuf(pb *services.ContractStateChange) ContractStateChange {
	if pb == nil {
		return ContractStateChange{}
	}

	storageChanges := make([]*StorageChange, 0)
	for _, sc := range pb.StorageChanges {
		temp := _StorageChangeFromProtobuf(sc)
		storageChanges = append(storageChanges, &temp)
	}

	return ContractStateChange{
		ContractID:     _ContractIDFromProtobuf(pb.ContractID),
		StorageChanges: storageChanges,
	}
}

func (csc *ContractStateChange) _ToProtobuf() *services.ContractStateChange {
	if csc.ContractID == nil {
		return &services.ContractStateChange{}
	}

	storageChanges := make([]*services.StorageChange, 0)
	for _, sc := range csc.StorageChanges {
		storageChanges = append(storageChanges, sc._ToProtobuf())
	}

	return &services.ContractStateChange{
		ContractID:     csc.ContractID._ToProtobuf(),
		StorageChanges: storageChanges,
	}
}

func (csc *ContractStateChange) ToBytes() []byte {
	data, err := protobuf.Marshal(csc._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func ContractStateChangeFromBytes(data []byte) (ContractStateChange, error) {
	if data == nil {
		return ContractStateChange{}, errByteArrayNull
	}
	pb := services.ContractStateChange{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ContractStateChange{}, err
	}

	return _ContractStateChangeFromProtobuf(&pb), nil
}
