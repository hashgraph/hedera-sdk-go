package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenTransfer struct {
	AccountID AccountID
	Amount    int64
}

func NewTokenTransfer(accountID AccountID, amount int64) TokenTransfer {
	return TokenTransfer{
		AccountID: accountID,
		Amount:    amount,
	}
}

func tokenTransferFromProtobuf(pb *proto.AccountAmount) TokenTransfer {
	if pb == nil {
		return TokenTransfer{}
	}
	return TokenTransfer{
		AccountID: accountIDFromProtobuf(pb.AccountID),
		Amount:    pb.Amount,
	}
}

func (transfer *TokenTransfer) toProtobuf() *proto.AccountAmount {
	return &proto.AccountAmount{
		AccountID: transfer.AccountID.toProtobuf(),
		Amount:    transfer.Amount,
	}
}

func (transfer TokenTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenTransferFromBytes(data []byte) (TokenTransfer, error) {
	if data == nil {
		return TokenTransfer{}, errByteArrayNull
	}
	pb := proto.AccountAmount{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenTransfer{}, err
	}

	return tokenTransferFromProtobuf(&pb), nil
}
