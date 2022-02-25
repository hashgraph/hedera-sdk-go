package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type HbarAllowance struct {
	OwnerAccountID   *AccountID
	SpenderAccountID *AccountID
	Amount           int64
}

func NewHbarAllowance(ownerAccountID AccountID, spenderAccountID AccountID, amount int64) HbarAllowance { //nolint
	return HbarAllowance{
		OwnerAccountID:   &ownerAccountID,
		SpenderAccountID: &spenderAccountID,
		Amount:           amount,
	}
}

func _HbarAllowanceFromProtobuf(pb *services.CryptoAllowance) HbarAllowance {
	return HbarAllowance{
		OwnerAccountID:   _AccountIDFromProtobuf(pb.Owner),
		SpenderAccountID: _AccountIDFromProtobuf(pb.Spender),
		Amount:           pb.Amount,
	}
}

func _HbarAllowanceFromGrantedProtobuf(pb *services.GrantedCryptoAllowance) HbarAllowance {
	return HbarAllowance{
		OwnerAccountID:   nil,
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

func (approval *HbarAllowance) _ToGrantedProtobuf() *services.GrantedCryptoAllowance {
	body := &services.GrantedCryptoAllowance{
		Amount: approval.Amount,
	}

	if approval.SpenderAccountID != nil {
		body.Spender = approval.SpenderAccountID._ToProtobuf()
	}

	return body
}
