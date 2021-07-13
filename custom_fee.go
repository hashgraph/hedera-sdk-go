package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type CustomFee struct {
	Fee                   Fee
	FeeCollectorAccountID *AccountID
}

func customFeeFromProtobuf(customFee *services.CustomFee, networkName *NetworkName) CustomFee {
	if customFee == nil {
		return CustomFee{}
	}

	var fee Fee
	switch t := customFee.Fee.(type) {
	case *services.CustomFee_FixedFee:
		fee = customFixedFeeFromProtobuf(t.FixedFee, networkName)
	case *services.CustomFee_FractionalFee:
		fee = customFractionalFeeFromProtobuf(t.FractionalFee)
	}

	id := accountIDFromProtobuf(customFee.FeeCollectorAccountId, networkName)

	return CustomFee{
		Fee:                   fee,
		FeeCollectorAccountID: &id,
	}
}

func (fee *CustomFee) toProtobuf() *services.CustomFee {
	var accountID *services.AccountID
	if fee.FeeCollectorAccountID != nil {
		accountID = fee.FeeCollectorAccountID.toProtobuf()
	}

	customFee := &services.CustomFee{
		FeeCollectorAccountId: accountID,
	}

	switch t := fee.Fee.(type) {
	case CustomFractionalFee:
		customFee.Fee = &services.CustomFee_FractionalFee{FractionalFee: t.toProtobuf()}
	case CustomFixedFee:
		customFee.Fee = &services.CustomFee_FixedFee{FixedFee: t.toProtobuf()}
	}

	return customFee
}
