package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// SPDX-License-Identifier: Apache-2.0

// TokenRelationship is the information about a token relationship
type TokenRelationship struct {
	TokenID              TokenID
	Symbol               string
	Balance              uint64
	KycStatus            *bool
	FreezeStatus         *bool
	Decimals             uint32
	AutomaticAssociation bool
}

func _TokenRelationshipFromProtobuf(pb *services.TokenRelationship) *TokenRelationship {
	if pb == nil {
		return &TokenRelationship{}
	}

	tokenID := TokenID{}
	if pb.TokenId != nil {
		tokenID = *_TokenIDFromProtobuf(pb.TokenId)
	}

	return &TokenRelationship{
		TokenID:              tokenID,
		Symbol:               pb.Symbol,
		Balance:              pb.Balance,
		KycStatus:            _KycStatusFromProtobuf(pb.KycStatus),
		FreezeStatus:         _FreezeStatusFromProtobuf(pb.FreezeStatus),
		Decimals:             pb.Decimals,
		AutomaticAssociation: pb.AutomaticAssociation,
	}
}

func (relationship *TokenRelationship) _ToProtobuf() *services.TokenRelationship {
	var freezeStatus services.TokenFreezeStatus
	switch *relationship.FreezeStatus {
	case true:
		freezeStatus = 1
	case false:
		freezeStatus = 2
	default:
		freezeStatus = 0
	}

	var kycStatus services.TokenKycStatus
	switch *relationship.KycStatus {
	case true:
		kycStatus = 1
	case false:
		kycStatus = 2
	default:
		kycStatus = 0
	}

	return &services.TokenRelationship{
		TokenId:              relationship.TokenID._ToProtobuf(),
		Symbol:               relationship.Symbol,
		Balance:              relationship.Balance,
		KycStatus:            kycStatus,
		FreezeStatus:         freezeStatus,
		Decimals:             relationship.Decimals,
		AutomaticAssociation: relationship.AutomaticAssociation,
	}
}

func (relationship TokenRelationship) ToBytes() []byte {
	data, err := protobuf.Marshal(relationship._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func _TokenRelationshipsFromProtobuf(rels []*services.TokenRelationship) []*TokenRelationship {
	tokenRelationships := make([]*TokenRelationship, 0)
	for _, tokenRelationship := range rels {
		tokenRelationships = append(tokenRelationships, _TokenRelationshipFromProtobuf(tokenRelationship))
	}
	return tokenRelationships
}

func TokenRelationshipFromBytes(data []byte) (*TokenRelationship, error) {
	if data == nil {
		return &TokenRelationship{}, errors.New("byte array can't be null")
	}
	pb := services.TokenRelationship{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return &TokenRelationship{}, err
	}

	return _TokenRelationshipFromProtobuf(&pb), nil
}
