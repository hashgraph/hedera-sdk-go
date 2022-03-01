package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type HbarAllowance struct {
	OwnerAccountID   *AccountID
	SpenderAccountID *AccountID
	Amount           int64
}

func NewHbarAllowance(ownerAccountID AccountID, spenderAccountID AccountID, amount int64) HbarAllowance {
	return HbarAllowance{
		OwnerAccountID:   &ownerAccountID,
		SpenderAccountID: &spenderAccountID,
		Amount:           amount,
	}
}

func (approval *HbarAllowance) _SetSpender(id AccountID) *HbarAllowance { //nolint
	approval.SpenderAccountID = &id
	return approval
}

func (approval *HbarAllowance) _GetSpender() AccountID { //nolint
	if approval.SpenderAccountID != nil {
		return *approval.SpenderAccountID
	}

	return AccountID{}
}

func (approval *HbarAllowance) _SetOwner(id AccountID) *HbarAllowance { //nolint
	approval.OwnerAccountID = &id
	return approval
}

func (approval *HbarAllowance) _GetOwner() AccountID { //nolint
	if approval.OwnerAccountID != nil {
		return *approval.OwnerAccountID
	}

	return AccountID{}
}

func (approval *HbarAllowance) _SetAmount(amount int64) *HbarAllowance { //nolint
	approval.Amount = amount
	return approval
}

func (approval *HbarAllowance) _GetAmount() int64 { //nolint
	return approval.Amount
}

func _HbarAllowanceFromProtobuf(pb *services.CryptoAllowance) HbarAllowance {
	return HbarAllowance{
		SpenderAccountID: _AccountIDFromProtobuf(pb.Spender),
		Amount:           pb.Amount,
	}
}

func (approval *HbarAllowance) _ToProtobuf() *services.CryptoAllowance {
	body := &services.CryptoAllowance{
		Amount: approval.Amount,
	}

	if approval.SpenderAccountID != nil {
		body.Spender = approval.SpenderAccountID._ToProtobuf()
	}

	if approval.OwnerAccountID != nil {
		body.Owner = approval.OwnerAccountID._ToProtobuf()
	}

	return body
}
