package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenRelationship struct {
	TokenID      TokenID
	Symbol       string
	Balance      uint64
	KycStatus    *bool
	FreezeStatus *bool
}

func tokenRelationshipFromProtobuf(pb *proto.TokenRelationship) TokenRelationship {
	return TokenRelationship{
		TokenID:      tokenIDFromProtobuf(pb.GetTokenId()),
		Symbol:       pb.Symbol,
		Balance:      pb.Balance,
		KycStatus:    kycStatusFromProtobuf(pb.KycStatus),
		FreezeStatus: freezeStatusFromProtobuf(pb.FreezeStatus),
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
		TokenId:      relationship.TokenID.toProtobuf(),
		Symbol:       relationship.Symbol,
		Balance:      relationship.Balance,
		KycStatus:    kycStatus,
		FreezeStatus: freezeStatus,
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
	pb := proto.TokenRelationship{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenRelationship{}, err
	}

	return tokenRelationshipFromProtobuf(&pb), nil
}
