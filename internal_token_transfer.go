package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type _TokenTransfer struct {
	Transfers        []*_HbarTransfer
	ExpectedDecimals *uint32
}

func _TokenTransferPrivateFromProtobuf(pb *services.TokenTransferList) *_TokenTransfer {
	if pb == nil {
		return &_TokenTransfer{}
	}

	var decimals *uint32
	if pb.ExpectedDecimals != nil {
		temp := pb.ExpectedDecimals.GetValue()
		decimals = &temp
	}

	return &_TokenTransfer{
		Transfers:        _HbarTransferFromProtobuf(pb.Transfers),
		ExpectedDecimals: decimals,
	}
}

func (transfer *_TokenTransfer) _ToProtobuf() []*services.AccountAmount {
	transfers := make([]*services.AccountAmount, 0)
	for _, t := range transfer.Transfers {
		transfers = append(transfers, &services.AccountAmount{
			AccountID:  t.accountID._ToProtobuf(),
			Amount:     t.Amount.AsTinybar(),
			IsApproval: t.IsApproved,
		})
	}
	return transfers
}
