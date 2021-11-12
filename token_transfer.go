package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
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

func _TokenTransferFromProtobuf(pb *services.AccountAmount) TokenTransfer {
	if pb == nil {
		return TokenTransfer{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	return TokenTransfer{
		AccountID: accountID,
		Amount:    pb.Amount,
	}
}

func (transfer *TokenTransfer) _ToProtobuf() *services.AccountAmount {
	return &services.AccountAmount{
		AccountID: transfer.AccountID._ToProtobuf(),
		Amount:    transfer.Amount,
	}
}

func (transfer TokenTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer._ToProtobuf())
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

	return _TokenTransferFromProtobuf(&pb), nil
}

func (transfer TokenTransfer) String() string {
	return fmt.Sprintf("accountID: %s, amount: %d", transfer.AccountID.String(), transfer.Amount)
}
