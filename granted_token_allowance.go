package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type GrantedTokenAllowance struct {
	TokenID          *TokenID
	SpenderAccountID *AccountID
	OwnerAccountID   *AccountID
	Amount           int64
}

func NewGrantedTokenAllowance(tokenID TokenID, owner AccountID, spender AccountID, amount int64) GrantedTokenAllowance { //nolint
	return GrantedTokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spender,
		OwnerAccountID:   &owner,
		Amount:           amount,
	}
}

func _GrantedTokenAllowanceFromProtobuf(pb *services.GrantedTokenAllowance) GrantedTokenAllowance {
	body := GrantedTokenAllowance{
		Amount: pb.Amount,
	}

	if pb.TokenId != nil {
		body.TokenID = _TokenIDFromProtobuf(pb.TokenId)
	}

	if pb.Spender != nil {
		body.SpenderAccountID = _AccountIDFromProtobuf(pb.Spender)
	}

	return body
}

func (approval *GrantedTokenAllowance) _ToProtobuf() *services.GrantedTokenAllowance {
	body := &services.GrantedTokenAllowance{
		Amount: approval.Amount,
	}

	if approval.SpenderAccountID != nil {
		body.Spender = approval.SpenderAccountID._ToProtobuf()
	}

	if approval.TokenID != nil {
		body.TokenId = approval.TokenID._ToProtobuf()
	}

	return body
}
