package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type CustomFractionalFee struct {
	CustomFee
	Numerator        int64
	Denominator      int64
	MinimumAmount    int64
	MaximumAmount    int64
	AssessmentMethod FeeAssessmentMethod
}

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

func (fee *CustomFractionalFee) SetFeeCollectorAccountID(id AccountID) *CustomFractionalFee {
	fee.FeeCollectorAccountID = &id
	return fee
}

func (fee *CustomFractionalFee) SetNumerator(numerator int64) *CustomFractionalFee {
	fee.Numerator = numerator
	return fee
}

func (fee *CustomFractionalFee) SetDenominator(denominator int64) *CustomFractionalFee {
	fee.Denominator = denominator
	return fee
}

func (fee *CustomFractionalFee) SetMin(min int64) *CustomFractionalFee {
	fee.MinimumAmount = min
	return fee
}

func (fee *CustomFractionalFee) SetMax(max int64) *CustomFractionalFee {
	fee.MaximumAmount = max
	return fee
}

func (fee *CustomFractionalFee) GetFeeCollectorAccountID() AccountID {
	if fee.FeeCollectorAccountID != nil {
		return *fee.FeeCollectorAccountID
	}

	return AccountID{}
}

func (fee *CustomFractionalFee) GetNumerator() int64 {
	return fee.Numerator
}

func (fee *CustomFractionalFee) GetDenominator() int64 {
	return fee.Denominator
}

func (fee *CustomFractionalFee) GetMin() int64 {
	return fee.MinimumAmount
}

func (fee *CustomFractionalFee) GetMax() int64 {
	return fee.MaximumAmount
}

func (fee *CustomFractionalFee) GetAssessmentMethod() FeeAssessmentMethod {
	return fee.AssessmentMethod
}

func _CustomFractionalFeeFromProtobuf(fractionalFee *proto.FractionalFee, fee CustomFee) CustomFractionalFee {
	return CustomFractionalFee{
		CustomFee:        fee,
		Numerator:        fractionalFee.FractionalAmount.Numerator,
		Denominator:      fractionalFee.FractionalAmount.Denominator,
		MinimumAmount:    fractionalFee.MinimumAmount,
		MaximumAmount:    fractionalFee.MaximumAmount,
		AssessmentMethod: FeeAssessmentMethod(fractionalFee.NetOfTransfers),
	}
}

func (fee CustomFractionalFee) _ValidateNetworkOnIDs(client *Client) error {
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

func (fee CustomFractionalFee) _ToProtobuf() *proto.CustomFee {
	var FeeCollectorAccountID *proto.AccountID
	if fee.FeeCollectorAccountID != nil {
		FeeCollectorAccountID = fee.CustomFee.FeeCollectorAccountID._ToProtobuf()
	}

	return &proto.CustomFee{
		Fee: &proto.CustomFee_FractionalFee{
			FractionalFee: &proto.FractionalFee{
				FractionalAmount: &proto.Fraction{
					Numerator:   fee.Numerator,
					Denominator: fee.Denominator,
				},
				MinimumAmount:  fee.MinimumAmount,
				MaximumAmount:  fee.MaximumAmount,
				NetOfTransfers: bool(fee.AssessmentMethod),
			},
		},
		FeeCollectorAccountId: FeeCollectorAccountID,
	}
}

func (fee CustomFractionalFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func (fee CustomFractionalFee) String() string {
	return fmt.Sprintf("feeCollectorAccountID: %s, numerator: %d, denominator: %d, min: %d, Max: %d, assessmentMethod: %s", fee.FeeCollectorAccountID.String(), fee.Numerator, fee.Denominator, fee.MinimumAmount, fee.MaximumAmount, fee.AssessmentMethod.String())
}
