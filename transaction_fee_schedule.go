package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

type TransactionFeeSchedule struct {
	RequestType RequestType
	// Deprecated: use Fees
	FeeData *FeeData
	Fees    []*FeeData
}

func _TransactionFeeScheduleFromProtobuf(txFeeSchedule *services.TransactionFeeSchedule) (TransactionFeeSchedule, error) {
	if txFeeSchedule == nil {
		return TransactionFeeSchedule{}, errParameterNull
	}

	feeData := make([]*FeeData, 0)

	for _, d := range txFeeSchedule.GetFees() {
		temp, err := _FeeDataFromProtobuf(d)
		if err != nil {
			return TransactionFeeSchedule{}, err
		}
		feeData = append(feeData, &temp)
	}

	return TransactionFeeSchedule{
		RequestType: RequestType(txFeeSchedule.GetHederaFunctionality()),
		Fees:        feeData,
	}, nil
}

func (txFeeSchedule TransactionFeeSchedule) _ToProtobuf() *services.TransactionFeeSchedule {
	feeData := make([]*services.FeeData, 0)
	if txFeeSchedule.Fees != nil {
		for _, data := range txFeeSchedule.Fees {
			feeData = append(feeData, data._ToProtobuf())
		}
	}

	var singleFee *services.FeeData
	if txFeeSchedule.FeeData != nil {
		singleFee = txFeeSchedule.FeeData._ToProtobuf()
	}

	return &services.TransactionFeeSchedule{
		HederaFunctionality: services.HederaFunctionality(txFeeSchedule.RequestType),
		Fees:                feeData,
		FeeData:             singleFee,
	}
}

func (txFeeSchedule TransactionFeeSchedule) ToBytes() []byte {
	data, err := protobuf.Marshal(txFeeSchedule._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func (txFeeSchedule TransactionFeeSchedule) String() string {
	str := ""
	for _, dat := range txFeeSchedule.Fees {
		str = str + dat.String() + ", "
	}

	if txFeeSchedule.FeeData != nil {
		return fmt.Sprintf("RequestType: %s, Feedata: %s", txFeeSchedule.RequestType.String(), txFeeSchedule.FeeData.String())
	}

	return fmt.Sprintf("RequestType: %s, Feedata: %s", txFeeSchedule.RequestType.String(), str)
}
