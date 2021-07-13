package hedera

import (
	"errors"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TokenRelationship struct {
	TokenID      TokenID
	Symbol       string
	Balance      uint64
	KycStatus    *bool
	FreezeStatus *bool
	Decimals     uint32
}

func tokenRelationshipFromProtobuf(pb *services.TokenRelationship, networkName *NetworkName) TokenRelationship {
	if pb == nil {
		return TokenRelationship{}
	}
	return TokenRelationship{
		TokenID:      tokenIDFromProtobuf(pb.GetTokenId(), networkName),
		Symbol:       pb.Symbol,
		Balance:      pb.Balance,
		KycStatus:    kycStatusFromProtobuf(pb.KycStatus),
		FreezeStatus: freezeStatusFromProtobuf(pb.FreezeStatus),
		Decimals:     pb.Decimals,
	}
}

func (relationship *TokenRelationship) toProtobuf() *services.TokenRelationship {
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
		TokenId:      relationship.TokenID.toProtobuf(),
		Symbol:       relationship.Symbol,
		Balance:      relationship.Balance,
		KycStatus:    kycStatus,
		FreezeStatus: freezeStatus,
		Decimals:     relationship.Decimals,
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
	pb := services.TokenRelationship{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenRelationship{}, err
	}

	return tokenRelationshipFromProtobuf(&pb, nil), nil
}
