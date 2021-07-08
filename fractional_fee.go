package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

type FractionalFee struct {
	Numerator     int64
	Denominator   int64
	MinimumAmount int64
	MaximumAmount int64
}

func fractionalFeeFromProtobuf(fractionalFee *proto.FractionalFee) FractionalFee {
	return FractionalFee{
		Numerator:     fractionalFee.FractionalAmount.Numerator,
		Denominator:   fractionalFee.FractionalAmount.Denominator,
		MinimumAmount: fractionalFee.MinimumAmount,
		MaximumAmount: fractionalFee.MaximumAmount,
	}
}

func (fee FractionalFee) toProtobuf() *proto.FractionalFee {
	return &proto.FractionalFee{
		FractionalAmount: &proto.Fraction{
			Numerator:   fee.Numerator,
			Denominator: fee.Denominator,
		},
		MinimumAmount: fee.MinimumAmount,
		MaximumAmount: fee.MaximumAmount,
	}
}

func (fee FractionalFee) Fee() {}
