package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type AccountAmount struct {
	AccountID AccountID
	Amount    int64
}

func accountAmountFromProto(pb *proto.AccountAmount) AccountAmount {
	return AccountAmount{
		AccountID: accountIDFromProto(pb.AccountID),
		Amount:    pb.Amount,
	}
}
