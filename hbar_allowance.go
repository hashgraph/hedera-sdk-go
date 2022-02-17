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
