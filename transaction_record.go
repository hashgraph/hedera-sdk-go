package hedera

import (
	"fmt"
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransactionRecord struct {
	Receipt            TransactionReceipt
	TransactionHash    []byte
	ConsensusTimestamp time.Time
	TransactionID      TransactionID
	TransactionMemo    string
	TransactionFee     Hbar
	Transfers          []Transfer
	TokenTransfers     map[TokenID][]TokenTransfer
	NftTransfers       map[TokenID][]TokenNftTransfer
	CallResult         *ContractFunctionResult
	CallResultIsCreate bool
	AssessedCustomFees []AssessedCustomFee
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

func _TransactionRecordFromProtobuf(pb *proto.TransactionRecord) TransactionRecord {
	if pb == nil {
		return TransactionRecord{}
	}
	var accountTransfers = make([]Transfer, len(pb.TransferList.AccountAmounts))
	var tokenTransfers = make(map[TokenID][]TokenTransfer)
	var nftTransfers = make(map[TokenID][]TokenNftTransfer)

	for i, element := range pb.TransferList.AccountAmounts {
		accountTransfers[i] = _TransferFromProtobuf(element)
	}

	for _, tokenTransfer := range pb.TokenTransferLists {
		for _, nftTransfer := range tokenTransfer.NftTransfers {
			if token := _TokenIDFromProtobuf(tokenTransfer.Token); token != nil {
				nftTransfers[*token] = append(nftTransfers[*token], _NftTransferFromProtobuf(nftTransfer))
			}
		}

		for _, accountAmount := range tokenTransfer.Transfers {
			if token := _TokenIDFromProtobuf(tokenTransfer.Token); token != nil {
				tokenTransfers[*token] = append(tokenTransfers[*token], _TokenTransferFromProtobuf(accountAmount))
			}
		}
	}

	assessedCustomFees := make([]AssessedCustomFee, 0)
	for _, fee := range pb.AssessedCustomFees {
		assessedCustomFees = append(assessedCustomFees, _AssessedCustomFeeFromProtobuf(fee))
	}

	txRecord := TransactionRecord{
		Receipt:            _TransactionReceiptFromProtobuf(pb.Receipt),
		TransactionHash:    pb.TransactionHash,
		ConsensusTimestamp: _TimeFromProtobuf(pb.ConsensusTimestamp),
		TransactionID:      _TransactionIDFromProtobuf(pb.TransactionID),
		TransactionMemo:    pb.Memo,
		TransactionFee:     HbarFromTinybar(int64(pb.TransactionFee)),
		Transfers:          accountTransfers,
		TokenTransfers:     tokenTransfers,
		NftTransfers:       nftTransfers,
		CallResultIsCreate: true,
		AssessedCustomFees: assessedCustomFees,
	}

	if pb.GetContractCreateResult() != nil {
		result := _ContractFunctionResultFromProtobuf(pb.GetContractCreateResult())

		txRecord.CallResult = &result
	} else if pb.GetContractCallResult() != nil {
		result := _ContractFunctionResultFromProtobuf(pb.GetContractCallResult())

		txRecord.CallResult = &result
		txRecord.CallResultIsCreate = false
	}

	return txRecord
}

func (record TransactionRecord) _ToProtobuf() (*proto.TransactionRecord, error) {
	var amounts = make([]*proto.AccountAmount, 0)
	for _, amount := range record.Transfers {
		amounts = append(amounts, &proto.AccountAmount{
			AccountID: amount.AccountID._ToProtobuf(),
			Amount:    amount.Amount.tinybar,
		})
	}

	var transferList = proto.TransferList{
		AccountAmounts: amounts,
	}

	var tokenTransfers = make([]*proto.TokenTransferList, 0)

	for tokenID, tokenTransfer := range record.TokenTransfers {
		tokenTemp := make([]*proto.AccountAmount, 0)

		for _, accountAmount := range tokenTransfer {
			tokenTemp = append(tokenTemp, accountAmount._ToProtobuf())
		}

		tokenTransfers = append(tokenTransfers, &proto.TokenTransferList{
			Token:     tokenID._ToProtobuf(),
			Transfers: tokenTemp,
		})
	}

	for tokenID, nftTransfers := range record.NftTransfers {
		nftTemp := make([]*proto.NftTransfer, 0)

		for _, nftTransfer := range nftTransfers {
			nftTemp = append(nftTemp, nftTransfer._ToProtobuf())
		}

		tokenTransfers = append(tokenTransfers, &proto.TokenTransferList{
			Token:        tokenID._ToProtobuf(),
			NftTransfers: nftTemp,
		})
	}

	assessedCustomFees := make([]*proto.AssessedCustomFee, 0)
	for _, fee := range record.AssessedCustomFees {
		assessedCustomFees = append(assessedCustomFees, fee._ToProtobuf())
	}

	var tRecord = proto.TransactionRecord{
		Receipt:         record.Receipt._ToProtobuf(),
		TransactionHash: record.TransactionHash,
		ConsensusTimestamp: &proto.Timestamp{
			Seconds: int64(record.ConsensusTimestamp.Second()),
			Nanos:   int32(record.ConsensusTimestamp.Nanosecond()),
		},
		TransactionID:      record.TransactionID._ToProtobuf(),
		Memo:               record.TransactionMemo,
		TransactionFee:     uint64(record.TransactionFee.AsTinybar()),
		TransferList:       &transferList,
		TokenTransferLists: tokenTransfers,
		AssessedCustomFees: assessedCustomFees,
	}

	var err error
	if record.CallResultIsCreate {
		var choice, err = record.GetContractCreateResult()

		if err != nil {
			return nil, err
		}

		tRecord.Body = &proto.TransactionRecord_ContractCreateResult{
			ContractCreateResult: choice._ToProtobuf(),
		}
	} else {
		var choice, err = record.GetContractExecuteResult()

		if err != nil {
			return nil, err
		}

		tRecord.Body = &proto.TransactionRecord_ContractCallResult{
			ContractCallResult: choice._ToProtobuf(),
		}
	}

	return &tRecord, err
}

func (record TransactionRecord) ToBytes() []byte {
	rec, err := record._ToProtobuf()
	if err != nil {
		return make([]byte, 0)
	}
	data, err := protobuf.Marshal(rec)
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TransactionRecordFromBytes(data []byte) (TransactionRecord, error) {
	if data == nil {
		return TransactionRecord{}, errByteArrayNull
	}
	pb := proto.TransactionRecord{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionRecord{}, err
	}

	return _TransactionRecordFromProtobuf(&pb), nil
}
