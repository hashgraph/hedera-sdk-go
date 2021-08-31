package hedera

import (
	"fmt"
	"time"

	protobuf "github.com/golang/protobuf/proto"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransactionRecord struct {
	Receipt                    TransactionReceipt
	TransactionHash            []byte
	ConsensusTimestamp         time.Time
	TransactionID              TransactionID
	TransactionMemo            string
	TransactionFee             Hbar
	Transfers                  []Transfer
	TokenTransfers             map[TokenID][]TokenTransfer
	NftTransfers               map[TokenID][]TokenNftTransfer
	CallResult                 *ContractFunctionResult
	CallResultIsCreate         bool
	AssessedCustomFees         []AssessedCustomFee
	AutomaticTokenAssociations []TokenAssociation
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

func transactionRecordFromProtobuf(pb *proto.TransactionRecord) TransactionRecord {
	if pb == nil {
		return TransactionRecord{}
	}
	var accountTransfers = make([]Transfer, len(pb.TransferList.AccountAmounts))
	var tokenTransfers = make(map[TokenID][]TokenTransfer, 0)
	var nftTransfers = make(map[TokenID][]TokenNftTransfer, 0)

	for i, element := range pb.TransferList.AccountAmounts {
		accountTransfers[i] = transferFromProtobuf(element)
	}

	for _, tokenTransfer := range pb.TokenTransferLists {
		for _, nftTransfer := range tokenTransfer.NftTransfers {
			token := tokenIDFromProtobuf(tokenTransfer.Token)
			nftTransfers[token] = append(nftTransfers[token], nftTransferFromProtobuf(nftTransfer))
		}

		for _, accountAmount := range tokenTransfer.Transfers {
			token := tokenIDFromProtobuf(tokenTransfer.Token)
			tokenTransfers[token] = append(tokenTransfers[token], tokenTransferFromProtobuf(accountAmount))
		}
	}

	assessedCustomFees := make([]AssessedCustomFee, 0)
	for _, fee := range pb.AssessedCustomFees {
		assessedCustomFees = append(assessedCustomFees, assessedCustomFeeFromProtobuf(fee))
	}

	tokenAssociation := make([]TokenAssociation, 0)
	for _, association := range pb.AutomaticTokenAssociations {
		tokenAssociation = append(tokenAssociation, tokenAssociationFromProtobuf(association))
	}

	txRecord := TransactionRecord{
		Receipt:                    transactionReceiptFromProtobuf(pb.Receipt),
		TransactionHash:            pb.TransactionHash,
		ConsensusTimestamp:         timeFromProtobuf(pb.ConsensusTimestamp),
		TransactionID:              transactionIDFromProtobuf(pb.TransactionID),
		TransactionMemo:            pb.Memo,
		TransactionFee:             HbarFromTinybar(int64(pb.TransactionFee)),
		Transfers:                  accountTransfers,
		TokenTransfers:             tokenTransfers,
		NftTransfers:               nftTransfers,
		CallResultIsCreate:         true,
		AssessedCustomFees:         assessedCustomFees,
		AutomaticTokenAssociations: tokenAssociation,
	}

	if pb.GetContractCreateResult() != nil {
		result := contractFunctionResultFromProtobuf(pb.GetContractCreateResult())

		txRecord.CallResult = &result
	} else if pb.GetContractCallResult() != nil {
		result := contractFunctionResultFromProtobuf(pb.GetContractCallResult())

		txRecord.CallResult = &result
		txRecord.CallResultIsCreate = false
	}

	return txRecord
}

func (record TransactionRecord) toProtobuf() (*proto.TransactionRecord, error) {
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

	var tokenTransfers = make([]*proto.TokenTransferList, 0)

	for tokenID, tokenTransfer := range record.TokenTransfers {
		tokenTemp := make([]*proto.AccountAmount, 0)

		for _, accountAmount := range tokenTransfer {
			tokenTemp = append(tokenTemp, accountAmount.toProtobuf())
		}

		tokenTransfers = append(tokenTransfers, &proto.TokenTransferList{
			Token:     tokenID.toProtobuf(),
			Transfers: tokenTemp,
		})
	}

	for tokenID, nftTransfers := range record.NftTransfers {
		nftTemp := make([]*proto.NftTransfer, 0)

		for _, nftTransfer := range nftTransfers {
			nftTemp = append(nftTemp, nftTransfer.toProtobuf())
		}

		tokenTransfers = append(tokenTransfers, &proto.TokenTransferList{
			Token:        tokenID.toProtobuf(),
			NftTransfers: nftTemp,
		})
	}

	assessedCustomFees := make([]*proto.AssessedCustomFee, 0)
	for _, fee := range record.AssessedCustomFees {
		assessedCustomFees = append(assessedCustomFees, fee.toProtobuf())
	}

	tokenAssociation := make([]*proto.TokenAssociation, 0)
	for _, association := range record.AutomaticTokenAssociations {
		tokenAssociation = append(tokenAssociation, association.toProtobuf())
	}

	var tRecord = proto.TransactionRecord{
		Receipt:         record.Receipt.toProtobuf(),
		TransactionHash: record.TransactionHash,
		ConsensusTimestamp: &proto.Timestamp{
			Seconds: int64(record.ConsensusTimestamp.Second()),
			Nanos:   int32(record.ConsensusTimestamp.Nanosecond()),
		},
		TransactionID:              record.TransactionID.toProtobuf(),
		Memo:                       record.TransactionMemo,
		TransactionFee:             uint64(record.TransactionFee.AsTinybar()),
		TransferList:               &transferList,
		TokenTransferLists:         tokenTransfers,
		AssessedCustomFees:         assessedCustomFees,
		AutomaticTokenAssociations: tokenAssociation,
	}

	var err error
	if record.CallResultIsCreate {
		var choice, err = record.GetContractCreateResult()

		if err != nil {
			return nil, err
		}

		tRecord.Body = &proto.TransactionRecord_ContractCreateResult{
			ContractCreateResult: choice.toProtobuf(),
		}
	} else {
		var choice, err = record.GetContractExecuteResult()

		if err != nil {
			return nil, err
		}

		tRecord.Body = &proto.TransactionRecord_ContractCallResult{
			ContractCallResult: choice.toProtobuf(),
		}
	}

	return &tRecord, err
}

func (record TransactionRecord) ToBytes() []byte {
	rec, err := record.toProtobuf()
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

	return transactionRecordFromProtobuf(&pb), nil
}
