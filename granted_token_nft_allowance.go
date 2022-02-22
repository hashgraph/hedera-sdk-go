package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type GrantedTokenNftAllowance struct {
	TokenID          *TokenID
	SpenderAccountID *AccountID
	OwnerAccountID   *AccountID
	SerialNumbers    []int64
	ApprovedForAll   bool
}

func NewGrantedTokenNftAllowance(tokenID TokenID, owner AccountID, spender AccountID, serialNumbers []int64, approvedForAll bool) GrantedTokenNftAllowance {
	return GrantedTokenNftAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spender,
		OwnerAccountID:   &owner,
		SerialNumbers:    serialNumbers,
		ApprovedForAll:   approvedForAll,
	}
}

func _GrantedTokenNftAllowanceFromProtobuf(pb *services.GrantedNftAllowance) GrantedTokenNftAllowance {
	body := GrantedTokenNftAllowance{
		ApprovedForAll: pb.ApprovedForAll,
		SerialNumbers:  pb.SerialNumbers,
	}

	if pb.TokenId != nil {
		body.TokenID = _TokenIDFromProtobuf(pb.TokenId)
	}

	if pb.Spender != nil {
		body.SpenderAccountID = _AccountIDFromProtobuf(pb.Spender)
	}

	return body
}

func (approval *GrantedTokenNftAllowance) _ToProtobuf() *services.GrantedNftAllowance {
	body := &services.GrantedNftAllowance{
		ApprovedForAll: approval.ApprovedForAll,
		SerialNumbers:  approval.SerialNumbers,
	}

	if approval.SpenderAccountID != nil {
		body.Spender = approval.SpenderAccountID._ToProtobuf()
	}

	if approval.TokenID != nil {
		body.TokenId = approval.TokenID._ToProtobuf()
	}

	return body
}
