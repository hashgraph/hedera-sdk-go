package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type CustomFee struct {
	Fee                   Fee
	FeeCollectorAccountID *AccountID
}

func customFeeFromProtobuf(customFee *proto.CustomFee, networkName *NetworkName) CustomFee {
	if customFee == nil {
		return CustomFee{}
	}

	var fee Fee
	switch t := customFee.Fee.(type) {
	case *proto.CustomFee_FixedFee:
		fee = customFixedFeeFromProtobuf(t.FixedFee, networkName)
	case *proto.CustomFee_FractionalFee:
		fee = customFractionalFeeFromProtobuf(t.FractionalFee)
	}

	var accountID AccountID
	if customFee.FeeCollectorAccountId != nil {
		accountID = accountIDFromProtobuf(customFee.FeeCollectorAccountId, networkName)
	}

	return CustomFee{
		Fee:                   fee,
		FeeCollectorAccountID: &accountID,
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

func (fee *CustomFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func CustomFeeFromBytes(data []byte) (CustomFee, error) {
	if data == nil {
		return CustomFee{}, errByteArrayNull
	}
	pb := proto.CustomFee{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return CustomFee{}, err
	}

	return customFeeFromProtobuf(&pb, nil), nil
}
