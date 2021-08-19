package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type CustomRoyaltyFee struct {
	CustomFee
	Numerator   int64
	Denominator int64
	FallbackFee *CustomFixedFee
}

func customRoyaltyFeeFromProtobuf(royalty *proto.RoyaltyFee, fee CustomFee) CustomRoyaltyFee {
	var fallback CustomFixedFee
	if royalty.FallbackFee != nil {
		fallback = customFixedFeeFromProtobuf(royalty.FallbackFee, fee)
	}
	return CustomRoyaltyFee{
		CustomFee:   fee,
		Numerator:   royalty.ExchangeValueFraction.Numerator,
		Denominator: royalty.ExchangeValueFraction.Denominator,
		FallbackFee: &fallback,
	}
}

func (fee CustomRoyaltyFee) validateNetworkOnIDs(client *Client) error {
	return fee.FallbackFee.validateNetworkOnIDs(client)
}

func (fee CustomRoyaltyFee) toProtobuf() *proto.CustomFee {
	var fallback proto.FixedFee
	if fee.FallbackFee != nil {
		if fee.FallbackFee.DenominationTokenID != nil {
			fallback.DenominatingTokenId = fee.FallbackFee.DenominationTokenID.toProtobuf()
		}
		fallback.Amount = fee.FallbackFee.Amount
	}

	var FeeCollectorAccountID *proto.AccountID
	if fee.FeeCollectorAccountID != nil {
		FeeCollectorAccountID = fee.CustomFee.FeeCollectorAccountID.toProtobuf()
	}

	return &proto.CustomFee{
		Fee: &proto.CustomFee_RoyaltyFee{
			RoyaltyFee: &proto.RoyaltyFee{
				ExchangeValueFraction: &proto.Fraction{
					Numerator:   fee.Numerator,
					Denominator: fee.Denominator,
				},
				FallbackFee: &fallback,
			},
		},
		FeeCollectorAccountId: FeeCollectorAccountID,
	}
}
