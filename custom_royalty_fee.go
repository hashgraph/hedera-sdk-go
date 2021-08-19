package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type CustomRoyaltyFee struct {
	CustomFee
	Numerator   int64
	Denominator int64
	FallbackFee *CustomFixedFee
}

func (fee *CustomRoyaltyFee) SetFeeCollectorAccountID(accountID AccountID) *CustomRoyaltyFee {
	fee.FeeCollectorAccountID = &accountID
	return fee
}

func (fee *CustomRoyaltyFee) SetNumerator(numerator int64) *CustomRoyaltyFee {
	fee.Numerator = numerator
	return fee
}

func (fee *CustomRoyaltyFee) SetDenominator(denominator int64) *CustomRoyaltyFee {
	fee.Denominator = denominator
	return fee
}

func (fee *CustomRoyaltyFee) SetFallbackFee(fallbackFee *CustomFixedFee) *CustomRoyaltyFee {
	fee.FallbackFee = fallbackFee
	return fee
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
	if client == nil || !client.autoValidateChecksums || fee.FallbackFee == nil {
		return nil
	}

	return fee.FallbackFee.validateNetworkOnIDs(client)
}

func (fee CustomRoyaltyFee) toProtobuf() *proto.CustomFee {
	var fallback *proto.FixedFee
	if fee.FallbackFee != nil {
		fallback = fee.FallbackFee.toProtobuf().GetFixedFee()
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
				FallbackFee: fallback,
			},
		},
		FeeCollectorAccountId: FeeCollectorAccountID,
	}
}
