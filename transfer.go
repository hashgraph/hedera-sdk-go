package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type Transfer struct {
	AccountID AccountID
	Amount    Hbar
}

func transferFromProtobuf(pb *services.AccountAmount, networkName *NetworkName) Transfer {
	if pb == nil {
		return Transfer{}
	}
	return Transfer{
		AccountID: accountIDFromProtobuf(pb.AccountID, networkName),
		Amount:    HbarFromTinybar(pb.Amount),
	}
}

func (transfer Transfer) toProtobuf() *services.TransferList {
	var ammounts = make([]*services.AccountAmount, 0)
	ammounts = append(ammounts, &services.AccountAmount{
		AccountID: transfer.AccountID.toProtobuf(),
		Amount:    transfer.Amount.AsTinybar(),
	})

	return &services.TransferList{
		AccountAmounts: ammounts,
	}
}
