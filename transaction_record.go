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
	callResult         *ContractFunctionResult
	callResultIsCreate bool
}

func (record TransactionRecord) GetContractExecuteResult() (ContractFunctionResult, error) {
	if record.callResult == nil || record.callResultIsCreate {
		return ContractFunctionResult{}, fmt.Errorf("record does not contain a contract execute result")
	}

	return *record.callResult, nil
}

func (record TransactionRecord) GetContractCreateResult() (ContractFunctionResult, error) {
	if record.callResult == nil || !record.callResultIsCreate {
		return ContractFunctionResult{}, fmt.Errorf("record does not contain a contract create result")
	}

	return *record.callResult, nil
}

func transactionRecordFromProto(pb *proto.TransactionRecord) TransactionRecord {
	var transferList = make([]Transfer, len(pb.TransferList.AccountAmounts))

	for i, element := range pb.TransferList.AccountAmounts {
		transferList[i] = transferFromProto(element)
	}

	txRecord := TransactionRecord{
		Receipt:            transactionReceiptFromProto(pb.Receipt),
		TransactionHash:    pb.TransactionHash,
		ConsensusTimestamp: timeFromProto(pb.ConsensusTimestamp),
		TransactionID:      transactionIDFromProto(pb.TransactionID),
		TransactionMemo:    pb.Memo,
		TransactionFee:     HbarFromTinybar(int64(pb.TransactionFee)),
		Transfers:          transferList,
		callResultIsCreate: true,
		callResult:         nil,
	}

	if pb.GetContractCreateResult() != nil {
		result := contractFunctionResultFromProto(pb.GetContractCreateResult())
		txRecord.callResult = &result

		return txRecord
	}

	result := contractFunctionResultFromProto(pb.GetContractCallResult())
	txRecord.callResult = &result
	txRecord.callResultIsCreate = false

	return txRecord
}
