package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// TokenNftAllowance is a struct to encapsulate the nft methods for token allowance's.
type TokenNftAllowance struct {
	TokenID           *TokenID
	SpenderAccountID  *AccountID
	OwnerAccountID    *AccountID
	SerialNumbers     []int64
	AllSerials        bool
	DelegatingSpender *AccountID
}

// NewTokenNftAllowance creates a TokenNftAllowance delegate for the given tokenID, owner, spender, serialNumbers, approvedForAll, and delegatingSpender
func NewTokenNftAllowance(tokenID TokenID, owner AccountID, spender AccountID, serialNumbers []int64, approvedForAll bool, delegatingSpender AccountID) TokenNftAllowance {
	return TokenNftAllowance{
		TokenID:           &tokenID,
		SpenderAccountID:  &spender,
		OwnerAccountID:    &owner,
		SerialNumbers:     serialNumbers,
		AllSerials:        approvedForAll,
		DelegatingSpender: &delegatingSpender,
	}
}

func _TokenNftAllowanceFromProtobuf(pb *services.NftAllowance) TokenNftAllowance {
	body := TokenNftAllowance{
		AllSerials:    pb.ApprovedForAll.GetValue(),
		SerialNumbers: pb.SerialNumbers,
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

	if pb.DelegatingSpender != nil {
		body.DelegatingSpender = _AccountIDFromProtobuf(pb.DelegatingSpender)
	}

	return body
}

func _TokenNftWipeAllowanceProtobuf(pb *services.NftRemoveAllowance) TokenNftAllowance {
	body := TokenNftAllowance{
		SerialNumbers: pb.SerialNumbers,
	}

	if pb.TokenId != nil {
		body.TokenID = _TokenIDFromProtobuf(pb.TokenId)
	}

	if pb.Owner != nil {
		body.SpenderAccountID = _AccountIDFromProtobuf(pb.Owner)
	}

	return body
}

func (approval *TokenNftAllowance) _ToProtobuf() *services.NftAllowance {
	body := &services.NftAllowance{
		ApprovedForAll: &wrapperspb.BoolValue{Value: approval.AllSerials},
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

	if approval.DelegatingSpender != nil {
		body.DelegatingSpender = approval.DelegatingSpender._ToProtobuf()
	}

	return body
}

func (approval *TokenNftAllowance) _ToWipeProtobuf() *services.NftRemoveAllowance {
	body := &services.NftRemoveAllowance{
		SerialNumbers: approval.SerialNumbers,
	}

	if approval.OwnerAccountID != nil {
		body.Owner = approval.OwnerAccountID._ToProtobuf()
	}

	if approval.TokenID != nil {
		body.TokenId = approval.TokenID._ToProtobuf()
	}

	return body
}

// String returns a string representation of the TokenNftAllowance
func (approval *TokenNftAllowance) String() string {
	var owner string
	var spender string
	var token string
	var serials string

	if approval.OwnerAccountID != nil {
		owner = approval.OwnerAccountID.String()
	}

	if approval.SpenderAccountID != nil {
		spender = approval.SpenderAccountID.String()
	}

	if approval.TokenID != nil {
		token = approval.TokenID.String()
	}

	for _, serial := range approval.SerialNumbers {
		serials += fmt.Sprintf("%d, ", serial)
	}

	return fmt.Sprintf("OwnerAccountID: %s, SpenderAccountID: %s, TokenID: %s, Serials: %s, ApprovedForAll: %t", owner, spender, token, serials, approval.AllSerials)
}
