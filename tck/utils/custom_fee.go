package utils

// SPDX-License-Identifier: Apache-2.0

import (
	"strconv"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

func ParseCustomFees(paramFees []param.CustomFee) ([]hiero.Fee, error) {
	var customFeeList []hiero.Fee
	for _, customFee := range paramFees {
		// Handle Fixed Fee
		if customFee.FixedFee != nil {
			fee := hiero.NewCustomFixedFee()
			fee.SetAmount(ParseIntFromOptional(&customFee.FixedFee.Amount))
			feeCollector, err := hiero.AccountIDFromString(customFee.FeeCollectorAccountId)
			if err != nil {
				return nil, err
			}
			fee.SetFeeCollectorAccountID(feeCollector)
			fee.SetAllCollectorsAreExempt(*customFee.FeeCollectorsExempt)
			if customFee.FixedFee.DenominatingTokenId != nil {
				tokenId, err := hiero.TokenIDFromString(*customFee.FixedFee.DenominatingTokenId)
				if err != nil {
					return nil, err
				}
				fee.SetDenominatingTokenID(tokenId)
			}
			customFeeList = append(customFeeList, fee)
		}

		// Handle Fractional Fee
		if customFee.FractionalFee != nil {
			fee := hiero.NewCustomFractionalFee()
			fee.SetNumerator(ParseIntFromOptional(&customFee.FractionalFee.Numerator))
			fee.SetDenominator(ParseIntFromOptional(&customFee.FractionalFee.Denominator))
			fee.SetMin(ParseIntFromOptional(&customFee.FractionalFee.MinimumAmount))
			fee.SetMax(ParseIntFromOptional(&customFee.FractionalFee.MaximumAmount))
			feeCollector, err := hiero.AccountIDFromString(customFee.FeeCollectorAccountId)
			if err != nil {
				return nil, err
			}
			fee.SetFeeCollectorAccountID(feeCollector)
			fee.SetAllCollectorsAreExempt(*customFee.FeeCollectorsExempt)
			customFeeList = append(customFeeList, fee)
		}

		// Handle Royalty Fee
		if customFee.RoyaltyFee != nil {
			fee := hiero.NewCustomRoyaltyFee()
			fee.SetNumerator(ParseIntFromOptional(&customFee.RoyaltyFee.Numerator))
			fee.SetDenominator(ParseIntFromOptional(&customFee.RoyaltyFee.Denominator))
			feeCollector, err := hiero.AccountIDFromString(customFee.FeeCollectorAccountId)
			if err != nil {
				return nil, err
			}
			fee.SetFeeCollectorAccountID(feeCollector)
			fee.SetAllCollectorsAreExempt(*customFee.FeeCollectorsExempt)

			if customFee.RoyaltyFee.FallbackFee != nil {
				fallback := hiero.NewCustomFixedFee()
				fallback.SetAmount(ParseIntFromOptional(&customFee.RoyaltyFee.FallbackFee.Amount))
				if customFee.RoyaltyFee.FallbackFee.DenominatingTokenId != nil {
					tokenId, err := hiero.TokenIDFromString(*customFee.RoyaltyFee.FallbackFee.DenominatingTokenId)
					if err != nil {
						return nil, err
					}
					fallback.SetDenominatingTokenID(tokenId)
				}
				fee.SetFallbackFee(fallback)
			}
			customFeeList = append(customFeeList, fee)
		}
	}
	return customFeeList, nil
}

func ParseIntFromOptional(value *string) int64 {
	if value == nil {
		return 0
	}
	parsed, _ := strconv.ParseInt(*value, 10, 64)
	return parsed
}
