package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// A royalty fee is a fractional fee that is assessed each time the ownership of an NFT is transferred from
// person A to person B. The fee collector account ID defined in the royalty fee schedule will receive the
// royalty fee each time. The royalty fee charged is a fraction of the value exchanged for the NFT.
type CustomRoyaltyFee struct {
	CustomFee
	Numerator   int64
	Denominator int64
	FallbackFee *CustomFixedFee
}

// A royalty fee is a fractional fee that is assessed each time the ownership of an NFT is transferred from
// person A to person B. The fee collector account ID defined in the royalty fee schedule will receive the
// royalty fee each time. The royalty fee charged is a fraction of the value exchanged for the NFT.
func NewCustomRoyaltyFee() *CustomRoyaltyFee {
	return &CustomRoyaltyFee{
		CustomFee:   CustomFee{},
		Numerator:   0,
		Denominator: 0,
		FallbackFee: nil,
	}
}

// SetFeeCollectorAccountID sets the account ID that will receive the custom fee
func (fee *CustomRoyaltyFee) SetFeeCollectorAccountID(accountID AccountID) *CustomRoyaltyFee {
	fee.FeeCollectorAccountID = &accountID
	return fee
}

// SetNumerator sets the numerator of the fractional fee
func (fee *CustomRoyaltyFee) SetNumerator(numerator int64) *CustomRoyaltyFee {
	fee.Numerator = numerator
	return fee
}

// SetDenominator sets the denominator of the fractional fee
func (fee *CustomRoyaltyFee) SetDenominator(denominator int64) *CustomRoyaltyFee {
	fee.Denominator = denominator
	return fee
}

// SetFallbackFee If present, the fixed fee to assess to the NFT receiver when no fungible value is exchanged with the sender
func (fee *CustomRoyaltyFee) SetFallbackFee(fallbackFee *CustomFixedFee) *CustomRoyaltyFee {
	fee.FallbackFee = fallbackFee
	return fee
}

// GetFeeCollectorAccountID returns the account ID that will receive the custom fee
func (fee *CustomRoyaltyFee) GetFeeCollectorAccountID() AccountID {
	if fee.FeeCollectorAccountID != nil {
		return *fee.FeeCollectorAccountID
	}

	return AccountID{}
}

// GetNumerator returns the numerator of the fee
func (fee *CustomRoyaltyFee) GetNumerator() int64 {
	return fee.Numerator
}

// GetDenominator returns the denominator of the fee
func (fee *CustomRoyaltyFee) GetDenominator() int64 {
	return fee.Denominator
}

// GetFallbackFee returns the fallback fee
func (fee *CustomRoyaltyFee) GetFallbackFee() CustomFixedFee {
	if fee.FallbackFee != nil {
		return *fee.FallbackFee
	}

	return CustomFixedFee{}
}

// SetAllCollectorsAreExempt sets whether all collectors are exempt from the fee
func (fee *CustomRoyaltyFee) SetAllCollectorsAreExempt(exempt bool) *CustomRoyaltyFee {
	fee.AllCollectorsAreExempt = exempt
	return fee
}

func _CustomRoyaltyFeeFromProtobuf(royalty *services.RoyaltyFee, fee CustomFee) CustomRoyaltyFee {
	var fallback CustomFixedFee
	if royalty.FallbackFee != nil {
		fallback = _CustomFixedFeeFromProtobuf(royalty.FallbackFee, fee)
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

func (fee CustomRoyaltyFee) _ToProtobuf() *services.CustomFee {
	var fallback *services.FixedFee
	if fee.FallbackFee != nil {
		fallback = fee.FallbackFee._ToProtobuf().GetFixedFee()
	}

	var FeeCollectorAccountID *services.AccountID
	if fee.FeeCollectorAccountID != nil {
		FeeCollectorAccountID = fee.CustomFee.FeeCollectorAccountID._ToProtobuf()
	}

	return &services.CustomFee{
		Fee: &services.CustomFee_RoyaltyFee{
			RoyaltyFee: &services.RoyaltyFee{
				ExchangeValueFraction: &services.Fraction{
					Numerator:   fee.Numerator,
					Denominator: fee.Denominator,
				},
				FallbackFee: fallback,
			},
		},
		FeeCollectorAccountId:  FeeCollectorAccountID,
		AllCollectorsAreExempt: fee.AllCollectorsAreExempt,
	}
}
