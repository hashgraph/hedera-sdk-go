package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type GrantedHbarAllowance struct {
	SpenderAccountID *AccountID
	Amount           int64
}

func NewGrantedHbarAllowance(ownerAccountID AccountID, spenderAccountID AccountID, amount int64) GrantedHbarAllowance {
	return GrantedHbarAllowance{
		SpenderAccountID: &spenderAccountID,
		Amount:           amount,
	}
}

func _GrantedHbarAllowanceFromProtobuf(pb *services.GrantedCryptoAllowance) GrantedHbarAllowance {
	return GrantedHbarAllowance{
		SpenderAccountID: _AccountIDFromProtobuf(pb.Spender),
		Amount:           pb.Amount,
	}
}

func (approval *GrantedHbarAllowance) _ToProtobuf() *services.GrantedCryptoAllowance {
	body := &services.GrantedCryptoAllowance{
		Amount: approval.Amount,
	}

	if approval.SpenderAccountID != nil {
		body.Spender = approval.SpenderAccountID._ToProtobuf()
	}

	return body
}
