package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TokenAllowance struct {
	TokenID          *TokenID
	SpenderAccountID *AccountID
	OwnerAccountID   *AccountID
	Amount           int64
}

func NewTokenAllowance(tokenID TokenID, owner AccountID, spender AccountID, amount int64) TokenAllowance { //nolint
	return TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spender,
		OwnerAccountID:   &owner,
		Amount:           amount,
	}
}

func _TokenAllowanceFromGrantedProtobuf(pb *services.GrantedTokenAllowance) TokenAllowance {
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

func (approval *TokenAllowance) _ToGrantedProtobuf() *services.GrantedTokenAllowance {
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

func (approval *TokenAllowance) String() string {
	var owner string
	var spender string
	var token string

	if approval.OwnerAccountID != nil {
		owner = approval.OwnerAccountID.String()
	}

	if approval.SpenderAccountID != nil {
		spender = approval.SpenderAccountID.String()
	}

	if approval.TokenID != nil {
		token = approval.TokenID.String()
	}

	return fmt.Sprintf("OwnerAccountID: %s, SpenderAccountID: %s, TokenID: %s, Amount: %s", owner, spender, token, HbarFromTinybar(approval.Amount).String())
}
