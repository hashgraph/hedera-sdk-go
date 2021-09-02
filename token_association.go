package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenAssociation struct {
	TokenID   *TokenID
	AccountID *AccountID
}

func tokenAssociationFromProtobuf(pb *proto.TokenAssociation) TokenAssociation {
	if pb == nil {
		return TokenAssociation{}
	}
	var tokenID TokenID
	if pb.TokenId != nil {
		tokenID = tokenIDFromProtobuf(pb.TokenId)
	}

	var accountID AccountID
	if pb.AccountId != nil {
		accountID = accountIDFromProtobuf(pb.AccountId)
	}

	return TokenAssociation{
		TokenID:   &tokenID,
		AccountID: &accountID,
	}
}

func (association *TokenAssociation) toProtobuf() *proto.TokenAssociation {
	var tokenID *proto.TokenID
	if association.TokenID != nil {
		tokenID = association.TokenID.toProtobuf()
	}

	var accountID *proto.AccountID
	if association.AccountID != nil {
		accountID = association.AccountID.toProtobuf()
	}

	return &proto.TokenAssociation{
		TokenId:   tokenID,
		AccountId: accountID,
	}
}

func (association *TokenAssociation) ToBytes() []byte {
	data, err := protobuf.Marshal(association.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenAssociationFromBytes(data []byte) (TokenAssociation, error) {
	if data == nil {
		return TokenAssociation{}, errByteArrayNull
	}
	pb := proto.TokenAssociation{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenAssociation{}, err
	}

	association := tokenAssociationFromProtobuf(&pb)

	return association, nil
}
