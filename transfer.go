package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type Transfer struct {
	AccountID AccountID
	Amount    Hbar
}

func _TransferFromProtobuf(pb *services.AccountAmount) Transfer {
	if pb == nil {
		return Transfer{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	return Transfer{
		AccountID: accountID,
		Amount:    HbarFromTinybar(pb.Amount),
	}
}

func (transfer Transfer) _ToProtobuf() *services.TransferList { // nolint
	var ammounts = make([]*services.AccountAmount, 0)
	ammounts = append(ammounts, &services.AccountAmount{
		AccountID: transfer.AccountID._ToProtobuf(),
		Amount:    transfer.Amount.AsTinybar(),
	})

	return &services.TransferList{
		AccountAmounts: ammounts,
	}
}
