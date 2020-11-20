package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

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
