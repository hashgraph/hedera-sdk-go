package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type Fee interface {
	Fee()
}

type CustomFixedFee struct {
	Amount              int64
	DenominationTokenID *TokenID
}

func customFixedFeeFromProtobuf(fixedFee *proto.FixedFee, networkName *NetworkName) CustomFixedFee {
	id := tokenIDFromProtobuf(fixedFee.DenominatingTokenId, networkName)
	return CustomFixedFee{
		Amount:              fixedFee.Amount,
		DenominationTokenID: &id,
	}
}

func (fee *CustomFixedFee) toProtobuf() *proto.FixedFee {
	var tokenID *proto.TokenID
	if fee.DenominationTokenID != nil {
		tokenID = fee.DenominationTokenID.toProtobuf()
	}

	return &proto.FixedFee{
		Amount:              fee.Amount,
		DenominatingTokenId: tokenID,
	}
}

func (fee CustomFixedFee) Fee() {}
