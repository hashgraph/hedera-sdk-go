package hiero

// SPDX-License-Identifier: Apache-2.0

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

type _HbarTransfer struct {
	accountID  *AccountID
	Amount     Hbar
	IsApproved bool
}

func _HbarTransferFromProtobuf(pb []*services.AccountAmount) []*_HbarTransfer {
	result := make([]*_HbarTransfer, 0)
	for _, acc := range pb {
		result = append(result, &_HbarTransfer{
			accountID:  _AccountIDFromProtobuf(acc.AccountID),
			Amount:     HbarFromTinybar(acc.Amount),
			IsApproved: acc.GetIsApproval(),
		})
	}

	return result
}

func (transfer *_HbarTransfer) _ToProtobuf() *services.AccountAmount { //nolint
	var account *services.AccountID
	if transfer.accountID != nil {
		account = transfer.accountID._ToProtobuf()
	}

	return &services.AccountAmount{
		AccountID:  account,
		Amount:     transfer.Amount.AsTinybar(),
		IsApproval: transfer.IsApproved,
	}
}
