package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

type Fee interface {
	_ToProtobuf() *services.CustomFee
	validateNetworkOnIDs(client *Client) error
}

// A fixed fee transfers a specified amount of the token, to the specified collection account(s),
// each time a token transfer is initiated. The custom token fee does not depend on the amount of the
// token that is being transferred.
type CustomFixedFee struct {
	CustomFee
	Amount              int64
	DenominationTokenID *TokenID
}

// A fixed fee transfers a specified amount of the token, to the specified collection account(s),
// each time a token transfer is initiated. The custom token fee does not depend on the amount of the
// token that is being transferred.
func NewCustomFixedFee() *CustomFixedFee {
	return &CustomFixedFee{
		CustomFee:           CustomFee{},
		Amount:              0,
		DenominationTokenID: nil,
	}
}

func _CustomFixedFeeFromProtobuf(fixedFee *services.FixedFee, customFee CustomFee) CustomFixedFee {
	return CustomFixedFee{
		CustomFee:           customFee,
		Amount:              fixedFee.Amount,
		DenominationTokenID: _TokenIDFromProtobuf(fixedFee.DenominatingTokenId),
	}
}

func (fee CustomFixedFee) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	if fee.DenominationTokenID != nil {
		if fee.DenominationTokenID != nil {
			if err := fee.DenominationTokenID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	if fee.FeeCollectorAccountID != nil {
		if fee.FeeCollectorAccountID != nil {
			if err := fee.FeeCollectorAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	return nil
}

func (fee CustomFixedFee) _ToProtobuf() *services.CustomFee {
	var tokenID *services.TokenID
	if fee.DenominationTokenID != nil {
		tokenID = fee.DenominationTokenID._ToProtobuf()
	}

	var FeeCollectorAccountID *services.AccountID
	if fee.FeeCollectorAccountID != nil {
		FeeCollectorAccountID = fee.CustomFee.FeeCollectorAccountID._ToProtobuf()
	}

	return &services.CustomFee{
		Fee: &services.CustomFee_FixedFee{
			FixedFee: &services.FixedFee{
				Amount:              fee.Amount,
				DenominatingTokenId: tokenID,
			},
		},
		FeeCollectorAccountId:  FeeCollectorAccountID,
		AllCollectorsAreExempt: fee.AllCollectorsAreExempt,
	}
}

// SetAmount sets the amount of the fixed fee in tinybar
func (fee *CustomFixedFee) SetAmount(tinybar int64) *CustomFixedFee {
	fee.Amount = tinybar
	return fee
}

// GetAmount returns the amount of the fixed fee
func (fee *CustomFixedFee) GetAmount() Hbar {
	return NewHbar(float64(fee.Amount))
}

// SetHbarAmount sets the amount of the fixed fee in hbar
func (fee *CustomFixedFee) SetHbarAmount(hbar Hbar) *CustomFixedFee {
	fee.Amount = int64(hbar.As(HbarUnits.Tinybar))
	fee.DenominationTokenID = nil
	return fee
}

// GetHbarAmount returns the amount of the fixed fee in hbar
func (fee *CustomFixedFee) GetHbarAmount() Hbar {
	return NewHbar(float64(fee.Amount))
}

// SetDenominatingTokenToSameToken sets the denomination token ID to the same token as the fee
func (fee *CustomFixedFee) SetDenominatingTokenToSameToken() *CustomFixedFee {
	fee.DenominationTokenID = &TokenID{0, 0, 0, nil}
	return fee
}

// SetDenominatingTokenID sets the denomination token ID
func (fee *CustomFixedFee) SetDenominatingTokenID(id TokenID) *CustomFixedFee {
	fee.DenominationTokenID = &id
	return fee
}

// GetDenominatingTokenID returns the denomination token ID
func (fee *CustomFixedFee) GetDenominatingTokenID() TokenID {
	if fee.DenominationTokenID != nil {
		return *fee.DenominationTokenID
	}

	return TokenID{}
}

// SetFeeCollectorAccountID sets the account ID that will receive the custom fee
func (fee *CustomFixedFee) SetFeeCollectorAccountID(id AccountID) *CustomFixedFee {
	fee.FeeCollectorAccountID = &id
	return fee
}

// GetFeeCollectorAccountID returns the account ID that will receive the custom fee
func (fee *CustomFixedFee) GetFeeCollectorAccountID() AccountID {
	return *fee.FeeCollectorAccountID
}

// SetAllCollectorsAreExempt sets whether all collectors are exempt from the custom fee
func (fee *CustomFixedFee) SetAllCollectorsAreExempt(exempt bool) *CustomFixedFee {
	fee.AllCollectorsAreExempt = exempt
	return fee
}

// ToBytes returns the byte representation of the CustomFixedFee
func (fee *CustomFixedFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// String returns a string representation of the CustomFixedFee
func (fee *CustomFixedFee) String() string {
	if fee.DenominationTokenID != nil {
		return fmt.Sprintf("feeCollectorAccountID: %s, amount: %d, denominatingTokenID: %s", fee.FeeCollectorAccountID.String(), fee.Amount, fee.DenominationTokenID.String())
	}

	return fmt.Sprintf("feeCollectorAccountID: %s, amount: %d", fee.FeeCollectorAccountID.String(), fee.Amount)
}
