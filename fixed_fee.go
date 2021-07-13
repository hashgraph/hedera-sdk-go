package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type Fee interface {
	Fee()
}

type CustomFixedFee struct {
	Amount              int64
	DenominationTokenID *TokenID
}

func customFixedFeeFromProtobuf(fixedFee *services.FixedFee, networkName *NetworkName) CustomFixedFee {
	id := tokenIDFromProtobuf(fixedFee.DenominatingTokenId, networkName)
	return CustomFixedFee{
		Amount:              fixedFee.Amount,
		DenominationTokenID: &id,
	}
}

func (fee *CustomFixedFee) toProtobuf() *services.FixedFee {
	var tokenID *services.TokenID
	if fee.DenominationTokenID != nil {
		tokenID = fee.DenominationTokenID.toProtobuf()
	}

	return &services.FixedFee{
		Amount:              fee.Amount,
		DenominatingTokenId: tokenID,
	}
}

func (fee CustomFixedFee) Fee() {}
