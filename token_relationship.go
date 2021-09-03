package hedera

import (
	"errors"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type TokenRelationship struct {
	TokenID              TokenID
	Symbol               string
	Balance              uint64
	KycStatus            *bool
	FreezeStatus         *bool
	Decimals             uint32
	AutomaticAssociation bool
}

func tokenRelationshipFromProtobuf(pb *proto.TokenRelationship) TokenRelationship {
	if pb == nil {
		return TokenRelationship{}
	}

	tokenID := TokenID{}
	if pb.TokenId != nil {
		tokenID = *tokenIDFromProtobuf(pb.TokenId)
	}

	return TokenRelationship{
		TokenID:              tokenID,
		Symbol:               pb.Symbol,
		Balance:              pb.Balance,
		KycStatus:            kycStatusFromProtobuf(pb.KycStatus),
		FreezeStatus:         freezeStatusFromProtobuf(pb.FreezeStatus),
		Decimals:             pb.Decimals,
		AutomaticAssociation: pb.AutomaticAssociation,
	}
}

func (relationship *TokenRelationship) toProtobuf() *proto.TokenRelationship {
	var freezeStatus proto.TokenFreezeStatus
	switch *relationship.FreezeStatus {
	case true:
		freezeStatus = 1
	case false:
		freezeStatus = 2
	default:
		freezeStatus = 0
	}

	var kycStatus proto.TokenKycStatus
	switch *relationship.KycStatus {
	case true:
		kycStatus = 1
	case false:
		kycStatus = 2
	default:
		kycStatus = 0
	}

	return &proto.TokenRelationship{
		TokenId:              relationship.TokenID.toProtobuf(),
		Symbol:               relationship.Symbol,
		Balance:              relationship.Balance,
		KycStatus:            kycStatus,
		FreezeStatus:         freezeStatus,
		Decimals:             relationship.Decimals,
		AutomaticAssociation: relationship.AutomaticAssociation,
	}
}

func (relationship TokenRelationship) ToBytes() []byte {
	data, err := protobuf.Marshal(relationship.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenRelationshipFromBytes(data []byte) (TokenRelationship, error) {
	if data == nil {
		return TokenRelationship{}, errors.New("byte array can't be null")
	}
	pb := proto.TokenRelationship{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenRelationship{}, err
	}

	return tokenRelationshipFromProtobuf(&pb), nil
}
