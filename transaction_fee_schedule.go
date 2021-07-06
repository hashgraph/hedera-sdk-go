package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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

func transactionFeeScheduleFromProtobuf(txFeeSchedule *proto.TransactionFeeSchedule) (TransactionFeeSchedule, error) {
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

func (txFeeSchedule TransactionFeeSchedule) toProtobuf() *proto.TransactionFeeSchedule {
	var feeData *proto.FeeData
	if txFeeSchedule.FeeData != nil {
		feeData = txFeeSchedule.FeeData.toProtobuf()
	}

	return &proto.TransactionFeeSchedule{
		HederaFunctionality: proto.HederaFunctionality(txFeeSchedule.RequestType),
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
	pb := proto.TransactionFeeSchedule{}
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
