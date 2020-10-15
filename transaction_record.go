package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionRecord struct {
	Receipt            TransactionReceipt
	TransactionHash    []byte
	ConsensusTimestamp time.Time
	TransactionID      TransactionID
	TransactionMemo    string
	TransactionFee     Hbar
	Transfers          []Transfer
	CallResult         *ContractFunctionResult
	CallResultIsCreate bool
}

func newTransactionRecord(
	receipt TransactionReceipt, transactionHash []byte,
	consensusTimestamp time.Time, transactionID TransactionID,
	transactionMemo string, transactionFee Hbar,
	transfers []Transfer, CallResult *ContractFunctionResult,
	CallResultIsCreate bool) TransactionRecord {

	record := TransactionRecord{
		Receipt:            receipt,
		TransactionHash:    transactionHash,
		ConsensusTimestamp: consensusTimestamp,
		TransactionID:      transactionID,
		TransactionMemo:    transactionMemo,
		TransactionFee:     transactionFee,
		Transfers:          transfers,
		CallResult:         CallResult,
		CallResultIsCreate: CallResultIsCreate,
	}

	return record

}

func (record TransactionRecord) GetContractExecuteResult() (ContractFunctionResult, error) {
	if record.CallResult == nil || record.CallResultIsCreate {
		return ContractFunctionResult{}, fmt.Errorf("record does not contain a contract execute result")
	}

	return *record.CallResult, nil
}

func (record TransactionRecord) GetContractCreateResult() (ContractFunctionResult, error) {
	if record.CallResult == nil || !record.CallResultIsCreate {
		return ContractFunctionResult{}, fmt.Errorf("record does not contain a contract create result")
	}

	return *record.CallResult, nil
}

func TransactionRecordFromProtobuf(pb *proto.TransactionRecord) TransactionRecord {
	var transferList = make([]Transfer, len(pb.TransferList.AccountAmounts))

	for i, element := range pb.TransferList.AccountAmounts {
		transferList[i] = transferFromProto(element)
	}

	txRecord := TransactionRecord{
		Receipt:            transactionReceiptFromProtobuf(pb.Receipt),
		TransactionHash:    pb.TransactionHash,
		ConsensusTimestamp: timeFromProto(pb.ConsensusTimestamp),
		TransactionID:      transactionIDFromProto(pb.TransactionID),
		TransactionMemo:    pb.Memo,
		TransactionFee:     HbarFromTinybar(int64(pb.TransactionFee)),
		Transfers:          transferList,
		CallResultIsCreate: true,
		CallResult:         nil,
	}

	if pb.GetContractCreateResult() != nil {
		result := contractFunctionResultFromProto(pb.GetContractCreateResult())

		txRecord.CallResult = &result
	} else if pb.GetContractCallResult() != nil {
		result := contractFunctionResultFromProto(pb.GetContractCallResult())

		txRecord.CallResult = &result
		txRecord.CallResultIsCreate = false
	}

	return txRecord
}

func (record TransactionRecord) toProtobuf() (proto.TransactionRecord, error) {
	var ammounts = make([]*proto.AccountAmount, 0)
	for _, ammount := range record.Transfers {
		ammounts = append(ammounts, &proto.AccountAmount{
			AccountID: ammount.AccountID.toProtobuf(),
			Amount:    ammount.Amount.tinybar,
		})
	}

	var transferList = proto.TransferList{
		AccountAmounts: ammounts,
	}

	var tRecord = proto.TransactionRecord{
		Receipt:         record.Receipt.toProtobuf(),
		TransactionHash: record.TransactionHash,
		ConsensusTimestamp: &proto.Timestamp{
			Seconds: int64(record.ConsensusTimestamp.Second()),
			Nanos:   int32(record.ConsensusTimestamp.Nanosecond()),
		},
		TransactionID:  record.TransactionID.toProtobuf(),
		Memo:           record.TransactionMemo,
		TransactionFee: uint64(record.TransactionFee.AsTinybar()),
		TransferList:   &transferList,
	}

	var err error
	if record.CallResultIsCreate {
		var choice, err = record.GetContractCreateResult()

		if err != nil {
			return proto.TransactionRecord{}, err
		}

		tRecord.Body = &proto.TransactionRecord_ContractCreateResult{
			ContractCreateResult: choice.toProtobuf(),
		}
	} else {
		var choice, err = record.GetContractExecuteResult()

		if err != nil {
			return proto.TransactionRecord{}, err
		}

		tRecord.Body = &proto.TransactionRecord_ContractCallResult{
			ContractCallResult: choice.toProtobuf(),
		}
	}

	return tRecord, err
}
