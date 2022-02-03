package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type _HbarTransfer struct {
	accountID  *AccountID
	Amount     Hbar
	IsApproved bool
}

type _HbarTransfers struct {
	transfers []*_HbarTransfer
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

func (transfers *_HbarTransfers) Len() int {
	return len(transfers.transfers)
}
func (transfers *_HbarTransfers) Swap(i, j int) {
	transfers.transfers[i], transfers.transfers[j] = transfers.transfers[j], transfers.transfers[i]
}

func (transfers *_HbarTransfers) Less(i, j int) bool {
	if transfers.transfers[i].accountID.Compare(*transfers.transfers[j].accountID) < 0 { //nolint
		return true
	}

	return false
}
