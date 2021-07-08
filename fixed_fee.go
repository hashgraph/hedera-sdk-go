package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type Fee interface {
	Fee()
}

type FixedFee struct {
	Amount              int64
	DenominationTokenID *TokenID
}

func fixedFeeFromProtobuf(fixedFee *proto.FixedFee, networkName *NetworkName) FixedFee {
	id := tokenIDFromProtobuf(fixedFee.DenominatingTokenId, networkName)
	return FixedFee{
		Amount:              fixedFee.Amount,
		DenominationTokenID: &id,
	}
}

func (fee *FixedFee) toProtobuf() *proto.FixedFee {
	var tokenID *proto.TokenID
	if fee.DenominationTokenID != nil {
		tokenID = fee.DenominationTokenID.toProtobuf()
	}

	return &proto.FixedFee{
		Amount:              fee.Amount,
		DenominatingTokenId: tokenID,
	}
}

func (fee FixedFee) Fee() {}
