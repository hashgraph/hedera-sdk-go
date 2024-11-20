package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// Base struct for all custom fees
type CustomFee struct {
	FeeCollectorAccountID  *AccountID
	AllCollectorsAreExempt bool
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
		FeeCollectorAccountID:  id,
		AllCollectorsAreExempt: customFee.AllCollectorsAreExempt,
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

// SetFeeCollectorAccountID sets the account ID that will receive the custom fee
func (fee *CustomFee) SetFeeCollectorAccountID(id AccountID) *CustomFee {
	fee.FeeCollectorAccountID = &id
	return fee
}

// GetFeeCollectorAccountID returns the account ID that will receive the custom fee
func (fee *CustomFee) GetFeeCollectorAccountID() AccountID {
	return *fee.FeeCollectorAccountID
}

// SetAllCollectorsAreExempt sets whether or not all collectors are exempt from the custom fee
func (fee *CustomFee) SetAllCollectorsAreExempt(exempt bool) *CustomFee {
	fee.AllCollectorsAreExempt = exempt
	return fee
}

// GetAllCollectorsAreExempt returns whether or not all collectors are exempt from the custom fee
func (fee *CustomFee) GetAllCollectorsAreExempt() bool {
	return fee.AllCollectorsAreExempt
}

// CustomFeeFromBytes returns a CustomFee from a raw protobuf byte array
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
