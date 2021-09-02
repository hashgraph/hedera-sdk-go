package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type Transfer struct {
	AccountID AccountID
	Amount    Hbar
}

func transferFromProtobuf(pb *proto.AccountAmount) Transfer {
	if pb == nil {
		return Transfer{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *accountIDFromProtobuf(pb.AccountID)
	}

	return Transfer{
		AccountID: accountID,
		Amount:    HbarFromTinybar(pb.Amount),
	}
}

func (transfer Transfer) toProtobuf() *proto.TransferList { // nolint
	var ammounts = make([]*proto.AccountAmount, 0)
	ammounts = append(ammounts, &proto.AccountAmount{
		AccountID: transfer.AccountID.toProtobuf(),
		Amount:    transfer.Amount.AsTinybar(),
	})

	return &proto.TransferList{
		AccountAmounts: ammounts,
	}
}
