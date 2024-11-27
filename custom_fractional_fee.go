package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// A fractional fee transfers the specified fraction of the total value of the tokens that are being transferred
// to the specified fee-collecting account. Along with setting a custom fractional fee, you can
type CustomFractionalFee struct {
	CustomFee
	Numerator        int64
	Denominator      int64
	MinimumAmount    int64
	MaximumAmount    int64
	AssessmentMethod FeeAssessmentMethod
}

// A fractional fee transfers the specified fraction of the total value of the tokens that are being transferred
// to the specified fee-collecting account. Along with setting a custom fractional fee, you can
func NewCustomFractionalFee() *CustomFractionalFee {
	return &CustomFractionalFee{
		CustomFee:        CustomFee{},
		Numerator:        0,
		Denominator:      0,
		MinimumAmount:    0,
		MaximumAmount:    0,
		AssessmentMethod: false,
	}
}

// SetFeeCollectorAccountID sets the account ID that will receive the custom fee
func (fee *CustomFractionalFee) SetFeeCollectorAccountID(id AccountID) *CustomFractionalFee {
	fee.FeeCollectorAccountID = &id
	return fee
}

// SetNumerator sets the numerator of the fractional fee
func (fee *CustomFractionalFee) SetNumerator(numerator int64) *CustomFractionalFee {
	fee.Numerator = numerator
	return fee
}

// SetDenominator sets the denominator of the fractional fee
func (fee *CustomFractionalFee) SetDenominator(denominator int64) *CustomFractionalFee {
	fee.Denominator = denominator
	return fee
}

// SetMin sets the minimum amount of the fractional fee
func (fee *CustomFractionalFee) SetMin(min int64) *CustomFractionalFee {
	fee.MinimumAmount = min
	return fee
}

// SetMax sets the maximum amount of the fractional fee
func (fee *CustomFractionalFee) SetMax(max int64) *CustomFractionalFee {
	fee.MaximumAmount = max
	return fee
}

// GetFeeCollectorAccountID returns the account ID that will receive the custom fee
func (fee *CustomFractionalFee) GetFeeCollectorAccountID() AccountID {
	if fee.FeeCollectorAccountID != nil {
		return *fee.FeeCollectorAccountID
	}

	return AccountID{}
}

// GetNumerator returns the numerator of the fractional fee
func (fee *CustomFractionalFee) GetNumerator() int64 {
	return fee.Numerator
}

// GetDenominator returns the denominator of the fractional fee
func (fee *CustomFractionalFee) GetDenominator() int64 {
	return fee.Denominator
}

// GetMin returns the minimum amount of the fractional fee
func (fee *CustomFractionalFee) GetMin() int64 {
	return fee.MinimumAmount
}

// GetMax returns the maximum amount of the fractional fee
func (fee *CustomFractionalFee) GetMax() int64 {
	return fee.MaximumAmount
}

// GetAssessmentMethod returns the assessment method of the fractional fee
func (fee *CustomFractionalFee) GetAssessmentMethod() FeeAssessmentMethod {
	return fee.AssessmentMethod
}

// SetAssessmentMethod sets the assessment method of the fractional fee
func (fee *CustomFractionalFee) SetAssessmentMethod(feeAssessmentMethod FeeAssessmentMethod) *CustomFractionalFee {
	fee.AssessmentMethod = feeAssessmentMethod
	return fee
}

func _CustomFractionalFeeFromProtobuf(fractionalFee *services.FractionalFee, fee CustomFee) CustomFractionalFee {
	return CustomFractionalFee{
		CustomFee:        fee,
		Numerator:        fractionalFee.FractionalAmount.Numerator,
		Denominator:      fractionalFee.FractionalAmount.Denominator,
		MinimumAmount:    fractionalFee.MinimumAmount,
		MaximumAmount:    fractionalFee.MaximumAmount,
		AssessmentMethod: FeeAssessmentMethod(fractionalFee.NetOfTransfers),
	}
}

func (fee CustomFractionalFee) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
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

// SetAllCollectorsAreExempt sets the flag that indicates if all collectors are exempt from the custom fee
func (fee *CustomFractionalFee) SetAllCollectorsAreExempt(exempt bool) *CustomFractionalFee {
	fee.AllCollectorsAreExempt = exempt
	return fee
}

func (fee CustomFractionalFee) _ToProtobuf() *services.CustomFee {
	var FeeCollectorAccountID *services.AccountID
	if fee.FeeCollectorAccountID != nil {
		FeeCollectorAccountID = fee.CustomFee.FeeCollectorAccountID._ToProtobuf()
	}

	return &services.CustomFee{
		Fee: &services.CustomFee_FractionalFee{
			FractionalFee: &services.FractionalFee{
				FractionalAmount: &services.Fraction{
					Numerator:   fee.Numerator,
					Denominator: fee.Denominator,
				},
				MinimumAmount:  fee.MinimumAmount,
				MaximumAmount:  fee.MaximumAmount,
				NetOfTransfers: bool(fee.AssessmentMethod),
			},
		},
		FeeCollectorAccountId:  FeeCollectorAccountID,
		AllCollectorsAreExempt: fee.AllCollectorsAreExempt,
	}
}

// ToBytes returns a byte array representation of the CustomFractionalFee
func (fee CustomFractionalFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// String returns a string representation of the CustomFractionalFee
func (fee CustomFractionalFee) String() string {
	return fmt.Sprintf("feeCollectorAccountID: %s, numerator: %d, denominator: %d, min: %d, Max: %d, assessmentMethod: %s", fee.FeeCollectorAccountID.String(), fee.Numerator, fee.Denominator, fee.MinimumAmount, fee.MaximumAmount, fee.AssessmentMethod.String())
}
