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

func (transfer Transfer) toProto() proto.TransferList {
	var ammounts = make([]*proto.AccountAmount, 0)
	ammounts = append(ammounts, &proto.AccountAmount{
		AccountID: transfer.AccountID.toProtobuf(),
		Amount:    transfer.Amount.AsTinybar(),
	})

	return proto.TransferList{
		AccountAmounts: ammounts,
	}
}
