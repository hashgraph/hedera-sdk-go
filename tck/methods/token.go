package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"strconv"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

type TokenService struct {
	sdkService *SDKService
}

func (t *TokenService) SetSdkService(service *SDKService) {
	t.sdkService = service
}

// CreateToken jRPC method for createToken
func (t *TokenService) CreateToken(_ context.Context, params param.CreateTokenParams) (*response.TokenResponse, error) {

	transaction := hiero.NewTokenCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.AdminKey != nil {
		key, err := utils.GetKeyFromString(*params.AdminKey)
		if err != nil {
			return nil, err
		}
		transaction.SetAdminKey(key)
	}

	if params.KycKey != nil {
		key, err := utils.GetKeyFromString(*params.KycKey)
		if err != nil {
			return nil, err
		}
		transaction.SetKycKey(key)
	}

	if params.FreezeKey != nil {
		key, err := utils.GetKeyFromString(*params.FreezeKey)
		if err != nil {
			return nil, err
		}
		transaction.SetFreezeKey(key)
	}

	if params.WipeKey != nil {
		key, err := utils.GetKeyFromString(*params.WipeKey)
		if err != nil {
			return nil, err
		}
		transaction.SetWipeKey(key)
	}

	if params.PauseKey != nil {
		key, err := utils.GetKeyFromString(*params.PauseKey)
		if err != nil {
			return nil, err
		}
		transaction.SetPauseKey(key)
	}

	if params.MetadataKey != nil {
		key, err := utils.GetKeyFromString(*params.MetadataKey)
		if err != nil {
			return nil, err
		}
		transaction.SetMetadataKey(key)
	}

	if params.SupplyKey != nil {
		key, err := utils.GetKeyFromString(*params.SupplyKey)
		if err != nil {
			return nil, err
		}
		transaction.SetSupplyKey(key)
	}

	if params.FeeScheduleKey != nil {
		key, err := utils.GetKeyFromString(*params.FeeScheduleKey)
		if err != nil {
			return nil, err
		}
		transaction.SetFeeScheduleKey(key)
	}

	if params.Name != nil {
		transaction.SetTokenName(*params.Name)
	}
	if params.Symbol != nil {
		transaction.SetTokenSymbol(*params.Symbol)
	}
	if params.Decimals != nil {
		transaction.SetDecimals(uint(*params.Decimals))
	}
	if params.Memo != nil {
		transaction.SetTokenMemo(*params.Memo)
	}
	if params.TokenType != nil {
		if *params.TokenType == "ft" {
			transaction.SetTokenType(hiero.TokenTypeFungibleCommon)
		} else if *params.TokenType == "nft" {
			transaction.SetTokenType(hiero.TokenTypeNonFungibleUnique)
		} else {
			return nil, response.InvalidParams.WithData("Invalid token type")
		}
	}
	if params.SupplyType != nil {
		if *params.SupplyType == "finite" {
			transaction.SetSupplyType(hiero.TokenSupplyTypeFinite)
		} else if *params.SupplyType == "infinite" {
			transaction.SetSupplyType(hiero.TokenSupplyTypeInfinite)
		} else {
			return nil, response.InvalidParams.WithData("Invalid supply type")
		}
	}
	if params.MaxSupply != nil {
		maxSupply, err := strconv.ParseInt(*params.MaxSupply, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetMaxSupply(maxSupply)
	}
	if params.InitialSupply != nil {
		initialSupply, err := strconv.ParseInt(*params.InitialSupply, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetInitialSupply(uint64(initialSupply))
	}
	if params.TreasuryAccountId != nil {
		accountID, err := hiero.AccountIDFromString(*params.TreasuryAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetTreasuryAccountID(accountID)
	}
	if params.FreezeDefault != nil {
		transaction.SetFreezeDefault(*params.FreezeDefault)
	}
	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}
	if params.AutoRenewAccountId != nil {
		autoRenewAccountId, err := hiero.AccountIDFromString(*params.AutoRenewAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewAccount(autoRenewAccountId)
	}
	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}

	if params.Metadata != nil {
		transaction.SetTokenMetadata([]byte(*params.Metadata))
	}

	if params.CommonTransactionParams != nil {
		params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
	}

	if params.CustomFees != nil {
		var customFeeList []hiero.Fee
		for _, customFee := range *params.CustomFees {
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
		transaction.SetCustomFees(customFeeList)
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{TokenId: receipt.TokenID.String(), Status: receipt.Status.String()}, nil
}

func ParseIntFromOptional(value *string) int64 {
	if value == nil {
		return 0
	}
	parsed, _ := strconv.ParseInt(*value, 10, 64)
	return parsed
}
