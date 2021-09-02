package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type TransactionFeeSchedule struct {
	RequestType RequestType
	// Deprecated use Fees
	FeeData *FeeData
	Fees    []*FeeData
}

func transactionFeeScheduleFromProtobuf(txFeeSchedule *proto.TransactionFeeSchedule) (TransactionFeeSchedule, error) {
	if txFeeSchedule == nil {
		return TransactionFeeSchedule{}, errParameterNull
	}

	feeData := make([]*FeeData, 0)

	for _, data := range txFeeSchedule.GetFees() {
		temp, err := feeDataFromProtobuf(data)
		if err != nil {
			return TransactionFeeSchedule{}, err
		}
		feeData = append(feeData, &temp)
	}

	singleFeeData, err := feeDataFromProtobuf(txFeeSchedule.GetFeeData()) // nolint
	if err != nil {
		return TransactionFeeSchedule{}, err
	}

	return TransactionFeeSchedule{
		RequestType: RequestType(txFeeSchedule.GetHederaFunctionality()),
		Fees:        feeData,
		FeeData:     &singleFeeData,
	}, nil
}

func (txFeeSchedule TransactionFeeSchedule) toProtobuf() *proto.TransactionFeeSchedule {
	feeData := make([]*proto.FeeData, 0)
	if txFeeSchedule.Fees != nil {
		for _, data := range txFeeSchedule.Fees {
			feeData = append(feeData, data.toProtobuf())
		}
	}

	var singleFee *proto.FeeData
	if txFeeSchedule.FeeData != nil {
		singleFee = txFeeSchedule.FeeData.toProtobuf()
	}

	return &proto.TransactionFeeSchedule{
		HederaFunctionality: proto.HederaFunctionality(txFeeSchedule.RequestType),
		Fees:                feeData,
		FeeData:             singleFee,
	}
}

func (txFeeSchedule TransactionFeeSchedule) ToBytes() []byte {
	data, err := protobuf.Marshal(txFeeSchedule.toProtobuf())
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
