package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type CustomRoyaltyFee struct {
	CustomFee
	Numerator   int64
	Denominator int64
	FallbackFee *CustomFixedFee
}

func NewCustomRoyaltyFee() *CustomRoyaltyFee {
	return &CustomRoyaltyFee{
		CustomFee:   CustomFee{},
		Numerator:   0,
		Denominator: 0,
		FallbackFee: nil,
	}
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

func (fee *CustomRoyaltyFee) GetFeeCollectorAccountID() AccountID {
	if fee.FeeCollectorAccountID != nil {
		return *fee.FeeCollectorAccountID
	}

	return AccountID{}
}

func (fee *CustomRoyaltyFee) GetNumerator() int64 {
	return fee.Numerator
}

func (fee *CustomRoyaltyFee) GetDenominator() int64 {
	return fee.Denominator
}

func (fee *CustomRoyaltyFee) GetFallbackFee() CustomFixedFee {
	if fee.FallbackFee != nil {
		return *fee.FallbackFee
	}

	return CustomFixedFee{}
}

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

func (fee CustomRoyaltyFee) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums || fee.FallbackFee == nil {
		return nil
	}

	return fee.FallbackFee._ValidateNetworkOnIDs(client)
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
