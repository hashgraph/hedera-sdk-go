package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type CustomFee struct {
	Fee                   Fee
	FeeCollectorAccountID *AccountID
}

func customFeeFromProtobuf(customFee *proto.CustomFee, networkName *NetworkName) CustomFee {
	var fee Fee
	switch t := customFee.Fee.(type) {
	case *proto.CustomFee_FixedFee:
		fee = fixedFeeFromProtobuf(t.FixedFee, networkName)
	case *proto.CustomFee_FractionalFee:
		fee = fractionalFeeFromProtobuf(t.FractionalFee)
	}

	id := accountIDFromProtobuf(customFee.FeeCollectorAccountId, networkName)

	return CustomFee{
		Fee:                   fee,
		FeeCollectorAccountID: &id,
	}
}

func (fee *CustomFee) toProtobuf() *proto.CustomFee {
	var accountID *proto.AccountID
	if fee.FeeCollectorAccountID != nil {
		accountID = fee.FeeCollectorAccountID.toProtobuf()
	}

	customFee := &proto.CustomFee{
		FeeCollectorAccountId: accountID,
	}

	switch t := fee.Fee.(type) {
	case CustomFractionalFee:
		customFee.Fee = &proto.CustomFee_FractionalFee{FractionalFee: t.toProtobuf()}
	case CustomFixedFee:
		customFee.Fee = &proto.CustomFee_FixedFee{FixedFee: t.toProtobuf()}
	}

	return customFee
}
