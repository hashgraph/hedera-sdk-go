<<<<<<< Updated upstream
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

func (approval *TokenNftAllowance) SetTokenID(id TokenID) *TokenNftAllowance {
	approval.TokenID = &id
	return approval
}

func (approval *TokenNftAllowance) GetTokenID() TokenID {
	if approval.TokenID != nil {
		return *approval.TokenID
	}

	return TokenID{}
}

func (approval *TokenNftAllowance) SetSpenderAccountID(id AccountID) *TokenNftAllowance {
	approval.SpenderAccountID = &id
	return approval
}

func (approval *TokenNftAllowance) GetSpenderAccountID() AccountID {
	if approval.SpenderAccountID != nil {
		return *approval.SpenderAccountID
	}

	return AccountID{}
}

func (approval *TokenNftAllowance) SetSerialNumbers(serials []int64) *TokenNftAllowance {
	approval.SerialNumbers = serials
	return approval
}

func (approval *TokenNftAllowance) GetSerialNumbers() []int64 {
	return approval.SerialNumbers
}

func (approval *TokenNftAllowance) SetApprovedForAll(approvedForAll bool) *TokenNftAllowance {
	approval.ApprovedForAll = approvedForAll
	return approval
}

func (approval *TokenNftAllowance) GetApprovedForAll() bool {
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
||||||| constructed merge base
=======
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

func (approval *TokenNftAllowance) _SetTokenID(id TokenID) *TokenNftAllowance {
	approval.TokenID = &id
	return approval
}

func (approval *TokenNftAllowance) _GetTokenID() TokenID {
	if approval.TokenID != nil {
		return *approval.TokenID
	}

	return TokenID{}
}

func (approval *TokenNftAllowance) _SetSpenderAccountID(id AccountID) *TokenNftAllowance {
	approval.SpenderAccountID = &id
	return approval
}

func (approval *TokenNftAllowance) _GetSpenderAccountID() AccountID {
	if approval.SpenderAccountID != nil {
		return *approval.SpenderAccountID
	}

	return AccountID{}
}

func (approval *TokenNftAllowance) _SetSerialNumbers(serials []int64) *TokenNftAllowance {
	approval.SerialNumbers = serials
	return approval
}

func (approval *TokenNftAllowance) _GetSerialNumbers() []int64 {
	return approval.SerialNumbers
}

func (approval *TokenNftAllowance) _SetApprovedForAll(approvedForAll bool) *TokenNftAllowance {
	approval.ApprovedForAll = approvedForAll
	return approval
}

func (approval *TokenNftAllowance) _GetApprovedForAll() bool {
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
>>>>>>> Stashed changes
