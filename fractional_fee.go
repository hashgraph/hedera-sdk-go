package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type CustomFractionalFee struct {
	Numerator     int64
	Denominator   int64
	MinimumAmount int64
	MaximumAmount int64
}

func customFractionalFeeFromProtobuf(fractionalFee *services.FractionalFee) CustomFractionalFee {
	return CustomFractionalFee{
		Numerator:     fractionalFee.FractionalAmount.Numerator,
		Denominator:   fractionalFee.FractionalAmount.Denominator,
		MinimumAmount: fractionalFee.MinimumAmount,
		MaximumAmount: fractionalFee.MaximumAmount,
	}
}

func (fee CustomFractionalFee) toProtobuf() *services.FractionalFee {
	return &services.FractionalFee{
		FractionalAmount: &services.Fraction{
			Numerator:   fee.Numerator,
			Denominator: fee.Denominator,
		},
		MinimumAmount: fee.MinimumAmount,
		MaximumAmount: fee.MaximumAmount,
	}
}

func (fee CustomFractionalFee) Fee() {}
