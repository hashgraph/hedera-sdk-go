package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type TokenAssociation struct {
	TokenID   *TokenID
	AccountID *AccountID
}

func tokenAssociationFromProtobuf(pb *services.TokenAssociation) TokenAssociation {
	if pb == nil {
		return TokenAssociation{}
	}

	return TokenAssociation{
		TokenID:   _TokenIDFromProtobuf(pb.TokenId),
		AccountID: _AccountIDFromProtobuf(pb.AccountId),
	}
}

func (association *TokenAssociation) toProtobuf() *services.TokenAssociation {
	var tokenID *services.TokenID
	if association.TokenID != nil {
		tokenID = association.TokenID._ToProtobuf()
	}

	var accountID *services.AccountID
	if association.AccountID != nil {
		accountID = association.AccountID._ToProtobuf()
	}

	return &services.TokenAssociation{
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
	pb := services.TokenAssociation{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenAssociation{}, err
	}

	association := tokenAssociationFromProtobuf(&pb)

	return association, nil
}
