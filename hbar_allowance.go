<<<<<<< Updated upstream
package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type HbarAllowance struct {
	SpenderAccountID *AccountID
	Amount           int64
}

func NewHbarAllowance(id AccountID, amount int64) HbarAllowance {
	return HbarAllowance{
		SpenderAccountID: &id,
		Amount:           amount,
	}
}

func (approval *HbarAllowance) SetSpender(id AccountID) *HbarAllowance {
	approval.SpenderAccountID = &id
	return approval
}

func (approval *HbarAllowance) GetSpender() AccountID {
	if approval.SpenderAccountID != nil {
		return *approval.SpenderAccountID
	}

	return AccountID{}
}

func (approval *HbarAllowance) SetAmount(amount int64) *HbarAllowance {
	approval.Amount = amount
	return approval
}

func (approval *HbarAllowance) GetAmount() int64 {
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

	return body
}
||||||| constructed merge base
=======
package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type HbarAllowance struct {
	SpenderAccountID *AccountID
	Amount           int64
}

func NewHbarAllowance(id AccountID, amount int64) HbarAllowance {
	return HbarAllowance{
		SpenderAccountID: &id,
		Amount:           amount,
	}
}

func (approval *HbarAllowance) _SetSpender(id AccountID) *HbarAllowance {
	approval.SpenderAccountID = &id
	return approval
}

func (approval *HbarAllowance) _GetSpender() AccountID {
	if approval.SpenderAccountID != nil {
		return *approval.SpenderAccountID
	}

	return AccountID{}
}

func (approval *HbarAllowance) _SetAmount(amount int64) *HbarAllowance {
	approval.Amount = amount
	return approval
}

func (approval *HbarAllowance) _GetAmount() int64 {
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

	return body
}
>>>>>>> Stashed changes
