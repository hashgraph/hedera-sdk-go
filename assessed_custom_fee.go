package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type AssessedCustomFee struct {
	Amount                int64
	TokenID               *TokenID
	FeeCollectorAccountId *AccountID
}

func assessedCustomFeeFromProtobuf(assessedFee *proto.AssessedCustomFee) AssessedCustomFee {
	accountID := accountIDFromProtobuf(assessedFee.FeeCollectorAccountId)
	tokenID := tokenIDFromProtobuf(assessedFee.TokenId)

	return AssessedCustomFee{
		Amount:                assessedFee.Amount,
		TokenID:               &tokenID,
		FeeCollectorAccountId: &accountID,
	}
}

func (fee *AssessedCustomFee) toProtobuf() *proto.AssessedCustomFee {
	var tokenID *proto.TokenID
	if fee.TokenID != nil {
		tokenID = fee.TokenID.toProtobuf()
	}

	var accountID *proto.AccountID
	if fee.TokenID != nil {
		accountID = fee.FeeCollectorAccountId.toProtobuf()
	}

	return &proto.AssessedCustomFee{
		Amount:                fee.Amount,
		TokenId:               tokenID,
		FeeCollectorAccountId: accountID,
	}
}

func (fee *AssessedCustomFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func AssessedCustomFeeFromBytes(data []byte) (AssessedCustomFee, error) {
	if data == nil {
		return AssessedCustomFee{}, errByteArrayNull
	}
	pb := proto.AssessedCustomFee{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AssessedCustomFee{}, err
	}

	return assessedCustomFeeFromProtobuf(&pb), nil
}
