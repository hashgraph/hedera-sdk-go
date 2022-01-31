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

func (approval *TokenAllowance) SetTokenID(id TokenID) *TokenAllowance {
	approval.TokenID = &id
	return approval
}

func (approval *TokenAllowance) GetTokenID() TokenID {
	if approval.TokenID != nil {
		return *approval.TokenID
	}

	return TokenID{}
}

func (approval *TokenAllowance) SetSpenderAccountID(id AccountID) *TokenAllowance {
	approval.SpenderAccountID = &id
	return approval
}

func (approval *TokenAllowance) GetSpenderAccountID() AccountID {
	if approval.SpenderAccountID != nil {
		return *approval.SpenderAccountID
	}

	return AccountID{}
}

func (approval *TokenAllowance) SetAmount(amount int64) *TokenAllowance {
	approval.Amount = amount
	return approval
}

func (approval *TokenAllowance) GetAmount() int64 {
	return approval.Amount
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
