package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type Transfer struct {
	AccountID AccountID
	Amount    Hbar
}

func transferFromProto(pb *proto.AccountAmount) Transfer {
	return Transfer{
		AccountID: accountIDFromProto(pb.AccountID),
		Amount:    HbarFromTinybar(pb.Amount),
	}
}
