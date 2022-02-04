package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TokenNftAllowance struct {
	TokenID          *TokenID
	SpenderAccountID *AccountID
	OwnerAccountID   *AccountID
	SerialNumbers    []int64
	ApprovedForAll   bool
}

func NewTokenNftAllowance(tokenID TokenID, owner AccountID, spender AccountID, serialNumbers []int64, approvedForAll bool) TokenNftAllowance {
	return TokenNftAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spender,
		OwnerAccountID:   &owner,
		SerialNumbers:    serialNumbers,
		ApprovedForAll:   approvedForAll,
	}
}

func (approval *TokenNftAllowance) _SetTokenID(id TokenID) *TokenNftAllowance { //nolint
	approval.TokenID = &id
	return approval
}

func (approval *TokenNftAllowance) _GetTokenID() TokenID { //nolint
	if approval.TokenID != nil {
		return *approval.TokenID
	}

	return TokenID{}
}

func (approval *TokenNftAllowance) _SetSpenderAccountID(id AccountID) *TokenNftAllowance { //nolint
	approval.SpenderAccountID = &id
	return approval
}

func (approval *TokenNftAllowance) _GetSpenderAccountID() AccountID { //nolint
	if approval.SpenderAccountID != nil {
		return *approval.SpenderAccountID
	}

	return AccountID{}
}

func (approval *TokenNftAllowance) _SetOwnerAccountID(id AccountID) *TokenNftAllowance { //nolint
	approval.OwnerAccountID = &id
	return approval
}

func (approval *TokenNftAllowance) _GetOwnerAccountID() AccountID { //nolint
	if approval.OwnerAccountID != nil {
		return *approval.OwnerAccountID
	}

	return AccountID{}
}

func (approval *TokenNftAllowance) _SetSerialNumbers(serials []int64) *TokenNftAllowance { //nolint
	approval.SerialNumbers = serials
	return approval
}

func (approval *TokenNftAllowance) _GetSerialNumbers() []int64 { //nolint
	return approval.SerialNumbers
}

func (approval *TokenNftAllowance) _SetApprovedForAll(approvedForAll bool) *TokenNftAllowance { //nolint
	approval.ApprovedForAll = approvedForAll
	return approval
}

func (approval *TokenNftAllowance) _GetApprovedForAll() bool { //nolint
	return approval.ApprovedForAll
}

func _TokenNftAllowanceFromProtobuf(pb *services.NftAllowance) TokenNftAllowance {
	body := TokenNftAllowance{
		ApprovedForAll: pb.ApprovedForAll.GetValue(),
		SerialNumbers:  pb.SerialNumbers,
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

func (approval *TokenNftAllowance) _ToProtobuf() *services.NftAllowance {
	body := &services.NftAllowance{
		ApprovedForAll: &wrapperspb.BoolValue{Value: approval.ApprovedForAll},
		SerialNumbers:  approval.SerialNumbers,
	}

	if approval.SpenderAccountID != nil {
		body.Spender = approval.SpenderAccountID._ToProtobuf()
	}

	if approval.OwnerAccountID != nil {
		body.Owner = approval.OwnerAccountID._ToProtobuf()
	}

	if approval.TokenID != nil {
		body.TokenId = approval.TokenID._ToProtobuf()
	}

	return body
}
