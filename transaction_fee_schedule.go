package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TransactionFeeSchedule struct {
	RequestType RequestType
	FeeData     *FeeData
}

func newTransactionFeeSchedule() TransactionFeeSchedule {
	return TransactionFeeSchedule{
		RequestType: RequestTypeNone,
		FeeData:     nil,
	}
}

func transactionFeeScheduleFromProtobuf(txFeeSchedule *services.TransactionFeeSchedule) (TransactionFeeSchedule, error) {
	if txFeeSchedule == nil {
		return TransactionFeeSchedule{}, errParameterNull
	}

	feeData, err := feeDataFromProtobuf(txFeeSchedule.GetFeeData())
	if err != nil {
		return TransactionFeeSchedule{}, err
	}

	return TransactionFeeSchedule{
		RequestType: RequestType(txFeeSchedule.GetHederaFunctionality()),
		FeeData:     &feeData,
	}, nil
}

func (txFeeSchedule TransactionFeeSchedule) toProtobuf() *services.TransactionFeeSchedule {
	var feeData *services.FeeData
	if txFeeSchedule.FeeData != nil {
		feeData = txFeeSchedule.FeeData.toProtobuf()
	}

	return &services.TransactionFeeSchedule{
		HederaFunctionality: services.HederaFunctionality(txFeeSchedule.RequestType),
		FeeData:             feeData,
	}
}

func (txFeeSchedule TransactionFeeSchedule) ToBytes() []byte {
	data, err := protobuf.Marshal(txFeeSchedule.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func transactionFeeScheduleFromBytes(data []byte) (TransactionFeeSchedule, error) {
	if data == nil {
		return TransactionFeeSchedule{}, errByteArrayNull
	}
	pb := services.TransactionFeeSchedule{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionFeeSchedule{}, err
	}

	info, err := transactionFeeScheduleFromProtobuf(&pb)
	if err != nil {
		return TransactionFeeSchedule{}, err
	}

	return info, nil
}

func (txFeeSchedule TransactionFeeSchedule) String() string {
	return fmt.Sprintf("RequestType: %s, Feedata: %s", txFeeSchedule.RequestType.String(), txFeeSchedule.FeeData.String())
}
