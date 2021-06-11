package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type NftTransfer struct {
	SenderAccountID   AccountID
	ReceiverAccountID AccountID
	SerialNumber      int64
}

func NewNftTransfer(sender AccountID, receiver AccountID, serial int64) NftTransfer {
	return NftTransfer{
		SenderAccountID:   sender,
		ReceiverAccountID: receiver,
		SerialNumber:      serial,
	}
}

func nftTransferFromProtobuf(pb *proto.NftTransfer) NftTransfer {
	if pb == nil {
		return NftTransfer{}
	}
	return NftTransfer{
		SenderAccountID:   accountIDFromProtobuf(pb.SenderAccountID),
		ReceiverAccountID: accountIDFromProtobuf(pb.ReceiverAccountID),
		SerialNumber:      pb.SerialNumber,
	}
}

func (transfer *NftTransfer) toProtobuf() *proto.NftTransfer {
	return &proto.NftTransfer{
		SenderAccountID:   transfer.SenderAccountID.toProtobuf(),
		ReceiverAccountID: transfer.ReceiverAccountID.toProtobuf(),
		SerialNumber:      transfer.SerialNumber,
	}
}

func (transfer NftTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func NftTransferFromBytes(data []byte) (NftTransfer, error) {
	if data == nil {
		return NftTransfer{}, errByteArrayNull
	}
	pb := proto.NftTransfer{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return NftTransfer{}, err
	}

	return nftTransferFromProtobuf(&pb), nil
}
