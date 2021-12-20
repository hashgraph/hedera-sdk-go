package hedera

import (
	"fmt"
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
	ParentConsensusTimestamp   time.Time
	AliasKey                   *PublicKey
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

func _TransactionRecordFromProtobuf(pb *services.TransactionRecord) TransactionRecord {
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

	tokenAssociation := make([]TokenAssociation, 0)
	for _, association := range pb.AutomaticTokenAssociations {
		tokenAssociation = append(tokenAssociation, tokenAssociationFromProtobuf(association))
	}

	var alias *PublicKey
	if len(pb.Alias) != 0 {
		pbKey := services.Key{}
		_ = protobuf.Unmarshal(pb.Alias, &pbKey)
		initialKey, _ := _KeyFromProtobuf(&pbKey)
		switch t2 := initialKey.(type) { //nolint
		case PublicKey:
			alias = &t2
		}
	}

	txRecord := TransactionRecord{
		Receipt:                    _TransactionReceiptFromProtobuf(pb.Receipt),
		TransactionHash:            pb.TransactionHash,
		ConsensusTimestamp:         _TimeFromProtobuf(pb.ConsensusTimestamp),
		TransactionID:              _TransactionIDFromProtobuf(pb.TransactionID),
		TransactionMemo:            pb.Memo,
		TransactionFee:             HbarFromTinybar(int64(pb.TransactionFee)),
		Transfers:                  accountTransfers,
		TokenTransfers:             tokenTransfers,
		NftTransfers:               nftTransfers,
		CallResultIsCreate:         true,
		AssessedCustomFees:         assessedCustomFees,
		AutomaticTokenAssociations: tokenAssociation,
		ParentConsensusTimestamp:   _TimeFromProtobuf(pb.ParentConsensusTimestamp),
		AliasKey:                   alias,
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

func (record TransactionRecord) _ToProtobuf() (*services.TransactionRecord, error) {
	var amounts = make([]*services.AccountAmount, 0)
	for _, amount := range record.Transfers {
		amounts = append(amounts, &services.AccountAmount{
			AccountID: amount.AccountID._ToProtobuf(),
			Amount:    amount.Amount.tinybar,
		})
	}

	var transferList = services.TransferList{
		AccountAmounts: amounts,
	}

	var tokenTransfers = make([]*services.TokenTransferList, 0)

	for tokenID, tokenTransfer := range record.TokenTransfers {
		tokenTemp := make([]*services.AccountAmount, 0)

		for _, accountAmount := range tokenTransfer {
			tokenTemp = append(tokenTemp, accountAmount._ToProtobuf())
		}

		tokenTransfers = append(tokenTransfers, &services.TokenTransferList{
			Token:     tokenID._ToProtobuf(),
			Transfers: tokenTemp,
		})
	}

	for tokenID, nftTransfers := range record.NftTransfers {
		nftTemp := make([]*services.NftTransfer, 0)

		for _, nftTransfer := range nftTransfers {
			nftTemp = append(nftTemp, nftTransfer._ToProtobuf())
		}

		tokenTransfers = append(tokenTransfers, &services.TokenTransferList{
			Token:        tokenID._ToProtobuf(),
			NftTransfers: nftTemp,
		})
	}

	assessedCustomFees := make([]*services.AssessedCustomFee, 0)
	for _, fee := range record.AssessedCustomFees {
		assessedCustomFees = append(assessedCustomFees, fee._ToProtobuf())
	}

	tokenAssociation := make([]*services.TokenAssociation, 0)
	for _, association := range record.AutomaticTokenAssociations {
		tokenAssociation = append(tokenAssociation, association.toProtobuf())
	}

	var alias []byte
	if record.AliasKey != nil {
		alias, _ = protobuf.Marshal(record.AliasKey._ToProtoKey())
	}

	var tRecord = services.TransactionRecord{
		Receipt:         record.Receipt._ToProtobuf(),
		TransactionHash: record.TransactionHash,
		ConsensusTimestamp: &services.Timestamp{
			Seconds: int64(record.ConsensusTimestamp.Second()),
			Nanos:   int32(record.ConsensusTimestamp.Nanosecond()),
		},
		TransactionID:              record.TransactionID._ToProtobuf(),
		Memo:                       record.TransactionMemo,
		TransactionFee:             uint64(record.TransactionFee.AsTinybar()),
		TransferList:               &transferList,
		TokenTransferLists:         tokenTransfers,
		AssessedCustomFees:         assessedCustomFees,
		AutomaticTokenAssociations: tokenAssociation,
		ParentConsensusTimestamp: &services.Timestamp{
			Seconds: int64(record.ParentConsensusTimestamp.Second()),
			Nanos:   int32(record.ParentConsensusTimestamp.Nanosecond()),
		},
		Alias: alias,
	}

	var err error
	if record.CallResultIsCreate {
		var choice, err = record.GetContractCreateResult()

		if err != nil {
			return nil, err
		}

		tRecord.Body = &services.TransactionRecord_ContractCreateResult{
			ContractCreateResult: choice._ToProtobuf(),
		}
	} else {
		var choice, err = record.GetContractExecuteResult()

		if err != nil {
			return nil, err
		}

		tRecord.Body = &services.TransactionRecord_ContractCallResult{
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
	pb := services.TransactionRecord{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionRecord{}, err
	}

	return _TransactionRecordFromProtobuf(&pb), nil
}
