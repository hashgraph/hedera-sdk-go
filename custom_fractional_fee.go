package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type CustomFractionalFee struct {
	CustomFee
	Numerator     int64
	Denominator   int64
	MinimumAmount int64
	MaximumAmount int64
}

func customFractionalFeeFromProtobuf(fractionalFee *proto.FractionalFee, fee CustomFee) CustomFractionalFee {
	return CustomFractionalFee{
		CustomFee:     fee,
		Numerator:     fractionalFee.FractionalAmount.Numerator,
		Denominator:   fractionalFee.FractionalAmount.Denominator,
		MinimumAmount: fractionalFee.MinimumAmount,
		MaximumAmount: fractionalFee.MaximumAmount,
	}
}

func (fee CustomFractionalFee) validateNetworkOnIDs(client *Client) error {
	if fee.FeeCollectorAccountID != nil {
		if err := fee.FeeCollectorAccountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (fee CustomFractionalFee) toProtobuf() *proto.CustomFee {
	var FeeCollectorAccountID *proto.AccountID
	if fee.FeeCollectorAccountID != nil {
		FeeCollectorAccountID = fee.CustomFee.FeeCollectorAccountID.toProtobuf()
	}

	return &proto.CustomFee{
		Fee: &proto.CustomFee_FractionalFee{
			FractionalFee: &proto.FractionalFee{
				FractionalAmount: &proto.Fraction{
					Numerator:   fee.Numerator,
					Denominator: fee.Denominator,
				},
				MinimumAmount: fee.MinimumAmount,
				MaximumAmount: fee.MaximumAmount,
			},
		},
		FeeCollectorAccountId: FeeCollectorAccountID,
	}
}

func (fee CustomFractionalFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}
