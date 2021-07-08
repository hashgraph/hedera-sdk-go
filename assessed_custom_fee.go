package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type AssessedCustomFee struct {
	Amount                int64
	TokenID               *TokenID
	FeeCollectorAccountId *AccountID
}

func assessedCustomFeeFromProtobuf(assessedFee *proto.AssessedCustomFee, networkName *NetworkName) AssessedCustomFee {
	accountID := accountIDFromProtobuf(assessedFee.FeeCollectorAccountId, networkName)
	tokenID := tokenIDFromProtobuf(assessedFee.TokenId, networkName)

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
