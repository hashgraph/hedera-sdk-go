package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type AssessedCustomFee struct {
	Amount                int64
	TokenID               *TokenID
	FeeCollectorAccountId *AccountID
}

func assessedCustomFeeFromProtobuf(assessedFee *services.AssessedCustomFee, networkName *NetworkName) AssessedCustomFee {
	accountID := accountIDFromProtobuf(assessedFee.FeeCollectorAccountId, networkName)
	tokenID := tokenIDFromProtobuf(assessedFee.TokenId, networkName)

	return AssessedCustomFee{
		Amount:                assessedFee.Amount,
		TokenID:               &tokenID,
		FeeCollectorAccountId: &accountID,
	}
}

func (fee *AssessedCustomFee) toProtobuf() *services.AssessedCustomFee {
	var tokenID *services.TokenID
	if fee.TokenID != nil {
		tokenID = fee.TokenID.toProtobuf()
	}

	var accountID *services.AccountID
	if fee.TokenID != nil {
		accountID = fee.FeeCollectorAccountId.toProtobuf()
	}

	return &services.AssessedCustomFee{
		Amount:                fee.Amount,
		TokenId:               tokenID,
		FeeCollectorAccountId: accountID,
	}
}
