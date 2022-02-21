package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TokenAllowance struct {
	TokenID          *TokenID
	SpenderAccountID *AccountID
	Amount           int64
}

func NewTokenAllowance(tokenID TokenID, spender AccountID, amount int64) TokenAllowance {
	return TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spender,
		Amount:           amount,
	}
}

func _TokenAllowanceFromProtobuf(pb *services.TokenAllowance) TokenAllowance {
	body := TokenAllowance{
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

func (approval *TokenAllowance) _ToProtobuf() *services.TokenAllowance {
	body := &services.TokenAllowance{
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
