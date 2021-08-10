package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenNftTransfer struct {
	SenderAccountID   AccountID
	ReceiverAccountID AccountID
	SerialNumber      int64
}

func newNftTransfer(sender AccountID, receiver AccountID, serial int64) TokenNftTransfer {
	return TokenNftTransfer{
		SenderAccountID:   sender,
		ReceiverAccountID: receiver,
		SerialNumber:      serial,
	}
}

func nftTransferFromProtobuf(pb *proto.NftTransfer) TokenNftTransfer {
	if pb == nil {
		return TokenNftTransfer{}
	}
	return TokenNftTransfer{
		SenderAccountID:   accountIDFromProtobuf(pb.SenderAccountID),
		ReceiverAccountID: accountIDFromProtobuf(pb.ReceiverAccountID),
		SerialNumber:      pb.SerialNumber,
	}
}

func (transfer *TokenNftTransfer) toProtobuf() *proto.NftTransfer {
	return &proto.NftTransfer{
		SenderAccountID:   transfer.SenderAccountID.toProtobuf(),
		ReceiverAccountID: transfer.ReceiverAccountID.toProtobuf(),
		SerialNumber:      transfer.SerialNumber,
	}
}

func (transfer TokenNftTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func NftTransferFromBytes(data []byte) (TokenNftTransfer, error) {
	if data == nil {
		return TokenNftTransfer{}, errByteArrayNull
	}
	pb := proto.NftTransfer{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenNftTransfer{}, err
	}

	return nftTransferFromProtobuf(&pb), nil
}
