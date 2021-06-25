package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type Transfer struct {
	AccountID AccountID
	Amount    Hbar
}

func transferFromProtobuf(pb *proto.AccountAmount, networkName *NetworkName) Transfer {
	if pb == nil {
		return Transfer{}
	}
	return Transfer{
		AccountID: accountIDFromProtobuf(pb.AccountID, networkName),
		Amount:    HbarFromTinybar(pb.Amount),
	}
}

func (transfer Transfer) toProtobuf() *proto.TransferList {
	var ammounts = make([]*proto.AccountAmount, 0)
	ammounts = append(ammounts, &proto.AccountAmount{
		AccountID: transfer.AccountID.toProtobuf(),
		Amount:    transfer.Amount.AsTinybar(),
	})

	return &proto.TransferList{
		AccountAmounts: ammounts,
	}
}
