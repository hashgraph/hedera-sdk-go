package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// Transfer is a transfer of hbars or tokens from one account to another
type Transfer struct {
	AccountID  AccountID
	Amount     Hbar
	IsApproved bool
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
		AccountID:  accountID,
		Amount:     HbarFromTinybar(pb.Amount),
		IsApproved: pb.GetIsApproval(),
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
