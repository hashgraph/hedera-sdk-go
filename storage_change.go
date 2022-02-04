package hedera

import (
	"math/big"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type StorageChange struct {
	Slot         *big.Int
	ValueRead    *big.Int
	ValueWritten *big.Int
}

func _StorageChangeFromProtobuf(pb *services.StorageChange) StorageChange {
	if pb == nil {
		return StorageChange{}
	}

	slot := new(big.Int)
	slot.SetBytes(pb.Slot)

	valueRead := new(big.Int)
	valueRead.SetBytes(pb.ValueRead)

	valueWritten := new(big.Int)
	valueWritten.SetBytes(pb.ValueWritten.GetValue())

	return StorageChange{
		Slot:         slot,
		ValueRead:    valueRead,
		ValueWritten: valueWritten,
	}
}

func (storageChange *StorageChange) _ToProtobuf() *services.StorageChange {
	return &services.StorageChange{
		Slot:      storageChange.Slot.Bytes(),
		ValueRead: storageChange.ValueRead.Bytes(),
		ValueWritten: &wrapperspb.BytesValue{
			Value: storageChange.ValueWritten.Bytes(),
		},
	}
}

func (storageChange *StorageChange) ToBytes() []byte {
	data, err := protobuf.Marshal(storageChange._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func StorageChangeFromBytes(data []byte) (StorageChange, error) {
	if data == nil {
		return StorageChange{}, errByteArrayNull
	}
	pb := services.StorageChange{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return StorageChange{}, err
	}

	return _StorageChangeFromProtobuf(&pb), nil
}
