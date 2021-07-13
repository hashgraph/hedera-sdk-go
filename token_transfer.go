package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"
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

func tokenTransferFromProtobuf(pb *services.AccountAmount, networkName *NetworkName) TokenTransfer {
	if pb == nil {
		return TokenTransfer{}
	}
	return TokenTransfer{
		AccountID: accountIDFromProtobuf(pb.AccountID, networkName),
		Amount:    pb.Amount,
	}
}

func (transfer *TokenTransfer) toProtobuf() *services.AccountAmount {
	return &services.AccountAmount{
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
	pb := services.AccountAmount{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenTransfer{}, err
	}

	return tokenTransferFromProtobuf(&pb, nil), nil
}
