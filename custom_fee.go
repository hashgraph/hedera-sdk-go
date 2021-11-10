package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type CustomFee struct {
	FeeCollectorAccountID *AccountID
}

func _CustomFeeFromProtobuf(customFee *services.CustomFee) Fee {
	if customFee == nil {
		return nil
	}

	var id *AccountID
	if customFee.FeeCollectorAccountId != nil {
		id = _AccountIDFromProtobuf(customFee.FeeCollectorAccountId)
	}

	fee := CustomFee{
		FeeCollectorAccountID: id,
	}

	switch t := customFee.Fee.(type) {
	case *services.CustomFee_FixedFee:
		return _CustomFixedFeeFromProtobuf(t.FixedFee, fee)
	case *services.CustomFee_FractionalFee:
		return _CustomFractionalFeeFromProtobuf(t.FractionalFee, fee)
	case *services.CustomFee_RoyaltyFee:
		return _CustomRoyaltyFeeFromProtobuf(t.RoyaltyFee, fee)
	}

	return nil
}

func (fee *CustomFee) SetFeeCollectorAccountID(id AccountID) *CustomFee {
	fee.FeeCollectorAccountID = &id
	return fee
}

func (fee *CustomFee) GetFeeCollectorAccountID() AccountID {
	return *fee.FeeCollectorAccountID
}

func CustomFeeFromBytes(data []byte) (Fee, error) {
	if data == nil {
		return nil, errByteArrayNull
	}
	pb := services.CustomFee{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return nil, err
	}

	return _CustomFeeFromProtobuf(&pb), nil
}
