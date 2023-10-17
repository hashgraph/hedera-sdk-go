package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/wrapperspb"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// The complete record for a transaction on Hedera that has reached consensus.
// This is not-free to request and is available for 1 hour after a transaction reaches consensus.
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
	ExpectedDecimals           map[TokenID]uint32
	CallResult                 *ContractFunctionResult
	CallResultIsCreate         bool
	AssessedCustomFees         []AssessedCustomFee
	AutomaticTokenAssociations []TokenAssociation
	ParentConsensusTimestamp   time.Time
	AliasKey                   *PublicKey
	Duplicates                 []TransactionRecord
	Children                   []TransactionRecord
	// Deprecated
	HbarAllowances []HbarAllowance
	// Deprecated
	TokenAllowances []TokenAllowance
	// Deprecated
	TokenNftAllowances []TokenNftAllowance
	EthereumHash       []byte
	PaidStakingRewards map[AccountID]Hbar
	PrngBytes          []byte
	PrngNumber         *int32
	EvmAddress         []byte
}

// MarshalJSON returns the JSON representation of the TransactionRecord
func (record TransactionRecord) MarshalJSON() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	m := make(map[string]interface{})
	type TransferJSON struct {
		AccountID  string      `json:"accountId"`
		Amount     interface{} `json:"amount"`
		IsApproved bool        `json:"isApproved"`
	}
	var transfersJSON []TransferJSON
	for _, t := range record.Transfers {
		transfersJSON = append(transfersJSON, TransferJSON{
			AccountID:  t.AccountID.String(),
			Amount:     fmt.Sprint(t.Amount.AsTinybar()),
			IsApproved: t.IsApproved,
		})
	}

	// It's weird because we need it without field names to match the specification
	tokenTransfersMap := make(map[string]map[string]string)

	for tokenId, transfers := range record.TokenTransfers {
		accountIdMap := make(map[string]string)
		for _, transfer := range transfers {
			accountIdMap[transfer.AccountID.String()] = fmt.Sprint(transfer.Amount)
		}
		tokenTransfersMap[tokenId.String()] = accountIdMap
	}

	type TransferNftTokensJSON struct {
		SenderAccountID            string `json:"sender"`
		ReceiverAccountIDAccountID string `json:"recipient"`
		IsApproved                 bool   `json:"isApproved"`
		SerialNumber               int64  `json:"serial"`
	}
	var transfersNftTokenJSON []TransferNftTokensJSON
	tokenTransfersNftMap := make(map[string][]TransferNftTokensJSON)
	for k, v := range record.NftTransfers {
		for _, tt := range v {
			transfersNftTokenJSON = append(transfersNftTokenJSON, TransferNftTokensJSON{
				SenderAccountID:            tt.SenderAccountID.String(),
				ReceiverAccountIDAccountID: tt.ReceiverAccountID.String(),
				IsApproved:                 tt.IsApproved,
				SerialNumber:               tt.SerialNumber,
			})
		}
		tokenTransfersNftMap[k.String()] = transfersNftTokenJSON
	}

	m["transactionHash"] = hex.EncodeToString(record.TransactionHash)
	m["transactionId"] = record.TransactionID.String()
	m["transactionMemo"] = record.TransactionMemo
	m["transactionFee"] = fmt.Sprint(record.TransactionFee.AsTinybar())
	m["transfers"] = transfersJSON
	m["tokenTransfers"] = tokenTransfersMap
	m["nftTransfers"] = tokenTransfersNftMap
	m["expectedDecimals"] = record.ExpectedDecimals
	m["callResultIsCreate"] = record.CallResultIsCreate

	type AssessedCustomFeeJSON struct {
		FeeCollectorAccountID string   `json:"feeCollectorAccountId"`
		TokenID               string   `json:"tokenId"`
		Amount                string   `json:"amount"`
		PayerAccountIDs       []string `json:"payerAccountIds"`
	}

	customFeesJSON := make([]AssessedCustomFeeJSON, len(record.AssessedCustomFees))

	for i, fee := range record.AssessedCustomFees {
		payerAccountIDsStr := make([]string, len(fee.PayerAccountIDs))
		for j, accID := range fee.PayerAccountIDs {
			payerAccountIDsStr[j] = accID.String()
		}
		customFeesJSON[i] = AssessedCustomFeeJSON{
			FeeCollectorAccountID: fee.FeeCollectorAccountId.String(),
			TokenID:               fee.TokenID.String(),
			Amount:                fmt.Sprint(fee.Amount),
			PayerAccountIDs:       payerAccountIDsStr,
		}
	}

	m["assessedCustomFees"] = customFeesJSON

	type TokenAssociationOutput struct {
		TokenID   string `json:"tokenId"`
		AccountID string `json:"accountId"`
	}
	automaticTokenAssociations := make([]TokenAssociationOutput, len(record.AutomaticTokenAssociations))
	for i, ta := range record.AutomaticTokenAssociations {
		automaticTokenAssociations[i] = TokenAssociationOutput{
			TokenID:   ta.TokenID.String(),
			AccountID: ta.AccountID.String(),
		}
	}
	m["automaticTokenAssociations"] = automaticTokenAssociations

	consensusTime := record.ConsensusTimestamp.UTC().Format("2006-01-02T15:04:05.000Z")
	parentConsesnusTime := record.ParentConsensusTimestamp.UTC().Format("2006-01-02T15:04:05.000Z")
	m["consensusTimestamp"] = consensusTime
	m["parentConsensusTimestamp"] = parentConsesnusTime

	m["aliasKey"] = fmt.Sprint(record.AliasKey)
	m["ethereumHash"] = hex.EncodeToString(record.EthereumHash)
	type PaidStakingReward struct {
		AccountID  string `json:"accountId"`
		Amount     string `json:"amount"`
		IsApproved bool   `json:"isApproved"`
	}
	var paidStakingRewards []PaidStakingReward
	for k, reward := range record.PaidStakingRewards {
		paidStakingReward := PaidStakingReward{
			AccountID:  k.String(),
			Amount:     fmt.Sprint(reward.AsTinybar()),
			IsApproved: false,
		}
		paidStakingRewards = append(paidStakingRewards, paidStakingReward)
	}

	sort.Slice(paidStakingRewards, func(i, j int) bool {
		return paidStakingRewards[i].AccountID < paidStakingRewards[j].AccountID
	})

	m["paidStakingRewards"] = paidStakingRewards

	m["prngBytes"] = hex.EncodeToString(record.PrngBytes)
	m["prngNumber"] = record.PrngNumber
	m["evmAddress"] = hex.EncodeToString(record.EvmAddress)
	m["receipt"] = record.Receipt
	m["children"] = record.Children
	m["duplicates"] = record.Duplicates

	receiptBytes, err := record.Receipt.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var receiptInterface interface{}
	err = json.Unmarshal(receiptBytes, &receiptInterface)
	if err != nil {
		return nil, err
	}
	m["receipt"] = receiptInterface

	result, err := json.Marshal(m)
	return result, err
}

// GetContractExecuteResult returns the ContractFunctionResult if the transaction was a contract call
func (record TransactionRecord) GetContractExecuteResult() (ContractFunctionResult, error) {
	if record.CallResult == nil || record.CallResultIsCreate {
		return ContractFunctionResult{}, fmt.Errorf("record does not contain a contract execute result")
	}

	return *record.CallResult, nil
}

// GetContractCreateResult returns the ContractFunctionResult if the transaction was a contract create
func (record TransactionRecord) GetContractCreateResult() (ContractFunctionResult, error) {
	if record.CallResult == nil || !record.CallResultIsCreate {
		return ContractFunctionResult{}, fmt.Errorf("record does not contain a contract create result")
	}

	return *record.CallResult, nil
}

func _TransactionRecordFromProtobuf(protoResponse *services.TransactionGetRecordResponse, txID *TransactionID) TransactionRecord {
	if protoResponse == nil {
		return TransactionRecord{}
	}
	pb := protoResponse.GetTransactionRecord()
	if pb == nil {
		return TransactionRecord{}
	}
	var accountTransfers = make([]Transfer, 0)
	var tokenTransfers = make(map[TokenID][]TokenTransfer)
	var nftTransfers = make(map[TokenID][]TokenNftTransfer)
	var expectedDecimals = make(map[TokenID]uint32)

	if pb.TransferList != nil {
		for _, element := range pb.TransferList.AccountAmounts {
			accountTransfers = append(accountTransfers, _TransferFromProtobuf(element))
		}
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

		if tokenTransfer.ExpectedDecimals != nil {
			if token := _TokenIDFromProtobuf(tokenTransfer.Token); token != nil {
				expectedDecimals[*token] = tokenTransfer.ExpectedDecimals.Value
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

	paidStakingRewards := make(map[AccountID]Hbar)
	for _, aa := range pb.PaidStakingRewards {
		account := _AccountIDFromProtobuf(aa.AccountID)
		if val, ok := paidStakingRewards[*account]; ok {
			paidStakingRewards[*account] = HbarFromTinybar(val.tinybar + aa.Amount)
		}

		paidStakingRewards[*account] = HbarFromTinybar(aa.Amount)
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

	childReceipts := make([]TransactionRecord, 0)
	if len(protoResponse.ChildTransactionRecords) > 0 {
		for _, r := range protoResponse.ChildTransactionRecords {
			childReceipts = append(childReceipts, _TransactionRecordFromProtobuf(&services.TransactionGetRecordResponse{TransactionRecord: r}, txID))
		}
	}

	duplicateReceipts := make([]TransactionRecord, 0)
	if len(protoResponse.DuplicateTransactionRecords) > 0 {
		for _, r := range protoResponse.DuplicateTransactionRecords {
			duplicateReceipts = append(duplicateReceipts, _TransactionRecordFromProtobuf(&services.TransactionGetRecordResponse{TransactionRecord: r}, txID))
		}
	}

	var transactionID TransactionID
	if pb.TransactionID != nil {
		transactionID = _TransactionIDFromProtobuf(pb.TransactionID)
	}

	txRecord := TransactionRecord{
		Receipt:                    _TransactionReceiptFromProtobuf(&services.TransactionGetReceiptResponse{Receipt: pb.GetReceipt()}, txID),
		TransactionHash:            pb.TransactionHash,
		ConsensusTimestamp:         _TimeFromProtobuf(pb.ConsensusTimestamp),
		TransactionID:              transactionID,
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
		Duplicates:                 duplicateReceipts,
		Children:                   childReceipts,
		EthereumHash:               pb.EthereumHash,
		PaidStakingRewards:         paidStakingRewards,
		EvmAddress:                 pb.EvmAddress,
	}

	if w, ok := pb.Entropy.(*services.TransactionRecord_PrngBytes); ok {
		txRecord.PrngBytes = w.PrngBytes
	}

	if w, ok := pb.Entropy.(*services.TransactionRecord_PrngNumber); ok {
		txRecord.PrngNumber = &w.PrngNumber
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

func (record TransactionRecord) _ToProtobuf() (*services.TransactionGetRecordResponse, error) {
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

		bod := &services.TokenTransferList{
			Token:     tokenID._ToProtobuf(),
			Transfers: tokenTemp,
		}

		if decimal, ok := record.ExpectedDecimals[tokenID]; ok {
			bod.ExpectedDecimals = &wrapperspb.UInt32Value{Value: decimal}
		}

		tokenTransfers = append(tokenTransfers, bod)
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

	paidStakingRewards := make([]*services.AccountAmount, 0)
	for account, hbar := range record.PaidStakingRewards {
		paidStakingRewards = append(paidStakingRewards, &services.AccountAmount{
			AccountID: account._ToProtobuf(),
			Amount:    hbar.AsTinybar(),
		})
	}

	var tRecord = services.TransactionRecord{
		Receipt:         record.Receipt._ToProtobuf().GetReceipt(),
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
		Alias:              alias,
		EthereumHash:       record.EthereumHash,
		PaidStakingRewards: paidStakingRewards,
		EvmAddress:         record.EvmAddress,
	}

	if record.PrngNumber != nil {
		tRecord.Entropy = &services.TransactionRecord_PrngNumber{PrngNumber: *record.PrngNumber}
	} else if len(record.PrngBytes) > 0 {
		tRecord.Entropy = &services.TransactionRecord_PrngBytes{PrngBytes: record.PrngBytes}
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

	childReceipts := make([]*services.TransactionRecord, 0)
	if len(record.Children) > 0 {
		for _, r := range record.Children {
			temp, err := r._ToProtobuf()
			if err != nil {
				return nil, err
			}
			childReceipts = append(childReceipts, temp.GetTransactionRecord())
		}
	}

	duplicateReceipts := make([]*services.TransactionRecord, 0)
	if len(record.Duplicates) > 0 {
		for _, r := range record.Duplicates {
			temp, err := r._ToProtobuf()
			if err != nil {
				return nil, err
			}
			duplicateReceipts = append(duplicateReceipts, temp.GetTransactionRecord())
		}
	}

	return &services.TransactionGetRecordResponse{
		TransactionRecord:           &tRecord,
		ChildTransactionRecords:     childReceipts,
		DuplicateTransactionRecords: duplicateReceipts,
	}, err
}

// Validate checks that the receipt status is Success
func (record TransactionRecord) ValidateReceiptStatus(shouldValidate bool) error {
	return record.Receipt.ValidateStatus(shouldValidate)
}

// ToBytes returns the serialized bytes of a TransactionRecord
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

// TransactionRecordFromBytes returns a TransactionRecord from a raw protobuf byte array
func TransactionRecordFromBytes(data []byte) (TransactionRecord, error) {
	if data == nil {
		return TransactionRecord{}, errByteArrayNull
	}
	pb := services.TransactionGetRecordResponse{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionRecord{}, err
	}

	return _TransactionRecordFromProtobuf(&pb, nil), nil
}
