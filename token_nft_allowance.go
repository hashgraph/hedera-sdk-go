package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TokenNftAllowance struct {
	TokenID          *TokenID
	SpenderAccountID *AccountID
	SerialNumbers    []int64
	ApprovedForAll   bool
}

func NewTokenNftAllowance(tokenID TokenID, spender AccountID, serialNumbers []int64, approvedForAll bool) TokenNftAllowance {
	return TokenNftAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spender,
		SerialNumbers:    serialNumbers,
		ApprovedForAll:   approvedForAll,
	}
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

	if approval.TokenID != nil {
		body.TokenId = approval.TokenID._ToProtobuf()
	}

	return body
}
