package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type CustomFractionalFee struct {
	Numerator     int64
	Denominator   int64
	MinimumAmount int64
	MaximumAmount int64
}

func fractionalFeeFromProtobuf(fractionalFee *proto.FractionalFee) CustomFractionalFee {
	return CustomFractionalFee{
		Numerator:     fractionalFee.FractionalAmount.Numerator,
		Denominator:   fractionalFee.FractionalAmount.Denominator,
		MinimumAmount: fractionalFee.MinimumAmount,
		MaximumAmount: fractionalFee.MaximumAmount,
	}
}

func (fee CustomFractionalFee) toProtobuf() *proto.FractionalFee {
	return &proto.FractionalFee{
		FractionalAmount: &proto.Fraction{
			Numerator:   fee.Numerator,
			Denominator: fee.Denominator,
		},
		MinimumAmount: fee.MinimumAmount,
		MaximumAmount: fee.MaximumAmount,
	}
}

func (fee CustomFractionalFee) Fee() {}
