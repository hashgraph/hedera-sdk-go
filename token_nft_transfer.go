package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type TokenNftTransfer struct {
	SenderAccountID   AccountID
	ReceiverAccountID AccountID
	SerialNumber      int64
}

func nftTransferFromProtobuf(pb *proto.NftTransfer) TokenNftTransfer {
	if pb == nil {
		return TokenNftTransfer{}
	}

	senderAccountID := AccountID{}
	if pb.SenderAccountID != nil {
		senderAccountID = *accountIDFromProtobuf(pb.SenderAccountID)
	}

	receiverAccountID := AccountID{}
	if pb.ReceiverAccountID != nil {
		receiverAccountID = *accountIDFromProtobuf(pb.ReceiverAccountID)
	}

	return TokenNftTransfer{
		SenderAccountID:   senderAccountID,
		ReceiverAccountID: receiverAccountID,
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
