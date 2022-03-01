package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TokenAllowance struct {
	TokenID          *TokenID
	SpenderAccountID *AccountID
	OwnerAccountID   *AccountID
	Amount           int64
}

func NewTokenAllowance(tokenID TokenID, owner AccountID, spender AccountID, amount int64) TokenAllowance {
	return TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spender,
		OwnerAccountID:   &owner,
		Amount:           amount,
	}
}

func (approval *TokenAllowance) _SetTokenID(id TokenID) *TokenAllowance { //nolint
	approval.TokenID = &id
	return approval
}

func (approval *TokenAllowance) _GetTokenID() TokenID { //nolint
	if approval.TokenID != nil {
		return *approval.TokenID
	}

	return TokenID{}
}

func (approval *TokenAllowance) _SetSpenderAccountID(id AccountID) *TokenAllowance { //nolint
	approval.SpenderAccountID = &id
	return approval
}

func (approval *TokenAllowance) _GetSpenderAccountID() AccountID { //nolint
	if approval.SpenderAccountID != nil {
		return *approval.SpenderAccountID
	}

	return AccountID{}
}

func (approval *TokenAllowance) _SetOwnerAccountID(id AccountID) *TokenAllowance { //nolint
	approval.OwnerAccountID = &id
	return approval
}

func (approval *TokenAllowance) _GetOwnerAccountID() AccountID { //nolint
	if approval.OwnerAccountID != nil {
		return *approval.OwnerAccountID
	}

	return AccountID{}
}

func (approval *TokenAllowance) _SetAmount(amount int64) *TokenAllowance { //nolint
	approval.Amount = amount
	return approval
}

func (approval *TokenAllowance) _GetAmount() int64 { //nolint
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

	if pb.Owner != nil {
		body.OwnerAccountID = _AccountIDFromProtobuf(pb.Owner)
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

	if approval.OwnerAccountID != nil {
		body.Owner = approval.OwnerAccountID._ToProtobuf()
	}

	return body
}
