package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type CustomFee struct {
	FeeCollectorAccountID *AccountID
}

func customFeeFromProtobuf(customFee *proto.CustomFee) Fee {
	if customFee == nil {
		return nil
	}

	var id *AccountID
	if customFee.FeeCollectorAccountId != nil {
		id_ := accountIDFromProtobuf(customFee.FeeCollectorAccountId)
		id = &id_
	}

	fee := CustomFee{
		FeeCollectorAccountID: id,
	}

	switch t := customFee.Fee.(type) {
	case *proto.CustomFee_FixedFee:
		return customFixedFeeFromProtobuf(t.FixedFee, fee)
	case *proto.CustomFee_FractionalFee:
		return customFractionalFeeFromProtobuf(t.FractionalFee, fee)
	}

	return nil
}

func CustomFeeFromBytes(data []byte) (Fee, error) {
	if data == nil {
		return nil, errByteArrayNull
	}
	pb := proto.CustomFee{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return nil, err
	}

	return customFeeFromProtobuf(&pb), nil
}
