package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type TransactionReceipt struct {
	Status                  Status
	ExchangeRate            *ExchangeRate
	TopicID                 *TopicID
	FileID                  *FileID
	ContractID              *ContractID
	AccountID               *AccountID
	TokenID                 *TokenID
	TopicSequenceNumber     uint64
	TopicRunningHash        []byte
	TopicRunningHashVersion uint64
	TotalSupply             uint64
	ScheduleID              *ScheduleID
	ScheduledTransactionID  *TransactionID
	SerialNumbers           []int64
	Duplicates              []TransactionReceipt
	Children                []TransactionReceipt
	TransactionID           *TransactionID
}

func _TransactionReceiptFromProtobuf(protoResponse *services.TransactionGetReceiptResponse, transactionID *TransactionID) TransactionReceipt {
	if protoResponse == nil {
		return TransactionReceipt{}
	}
	protoReceipt := protoResponse.GetReceipt()
	if protoReceipt == nil {
		return TransactionReceipt{}
	}
	var accountID *AccountID
	if protoReceipt.AccountID != nil {
		accountID = _AccountIDFromProtobuf(protoReceipt.AccountID)
	}

	var contractID *ContractID
	if protoReceipt.ContractID != nil {
		contractID = _ContractIDFromProtobuf(protoReceipt.ContractID)
	}

	var fileID *FileID
	if protoReceipt.FileID != nil {
		fileID = _FileIDFromProtobuf(protoReceipt.FileID)
	}

	var topicID *TopicID
	if protoReceipt.TopicID != nil {
		topicID = _TopicIDFromProtobuf(protoReceipt.TopicID)
	}

	var rate *ExchangeRate
	if protoReceipt.ExchangeRate != nil {
		exchangeRateValue := _ExchangeRateFromProtobuf(protoReceipt.ExchangeRate.GetCurrentRate())
		rate = &exchangeRateValue
	}

	var topicSequenceHash []byte
	if protoReceipt.TopicRunningHash != nil {
		topicHash := protoReceipt.TopicRunningHash
		topicSequenceHash = topicHash
	}

	var tokenID *TokenID
	if protoReceipt.TokenID != nil {
		tokenID = _TokenIDFromProtobuf(protoReceipt.TokenID)
	}

	var scheduleID *ScheduleID
	if protoReceipt.ScheduleID != nil {
		scheduleID = _ScheduleIDFromProtobuf(protoReceipt.ScheduleID)
	}

	var scheduledTransactionID *TransactionID
	if protoReceipt.ScheduledTransactionID != nil {
		scheduledTransactionIDValue := _TransactionIDFromProtobuf(protoReceipt.ScheduledTransactionID)
		scheduledTransactionID = &scheduledTransactionIDValue
	}

	childReceipts := make([]TransactionReceipt, 0)
	if len(protoResponse.ChildTransactionReceipts) > 0 {
		for _, r := range protoResponse.ChildTransactionReceipts {
			childReceipts = append(childReceipts, _TransactionReceiptFromProtobuf(&services.TransactionGetReceiptResponse{Receipt: r}, transactionID))
		}
	}

	duplicateReceipts := make([]TransactionReceipt, 0)
	if len(protoResponse.DuplicateTransactionReceipts) > 0 {
		for _, r := range protoResponse.DuplicateTransactionReceipts {
			duplicateReceipts = append(duplicateReceipts, _TransactionReceiptFromProtobuf(&services.TransactionGetReceiptResponse{Receipt: r}, transactionID))
		}
	}

	return TransactionReceipt{
		Status:                  Status(protoReceipt.Status),
		ExchangeRate:            rate,
		TopicID:                 topicID,
		FileID:                  fileID,
		ContractID:              contractID,
		AccountID:               accountID,
		TopicSequenceNumber:     protoReceipt.TopicSequenceNumber,
		TopicRunningHash:        topicSequenceHash,
		TopicRunningHashVersion: protoReceipt.TopicRunningHashVersion,
		TokenID:                 tokenID,
		TotalSupply:             protoReceipt.NewTotalSupply,
		ScheduleID:              scheduleID,
		ScheduledTransactionID:  scheduledTransactionID,
		SerialNumbers:           protoReceipt.SerialNumbers,
		Children:                childReceipts,
		Duplicates:              duplicateReceipts,
		TransactionID:           transactionID,
	}
}

func (receipt TransactionReceipt) _ToProtobuf() *services.TransactionGetReceiptResponse {
	receiptFinal := services.TransactionReceipt{
		Status:                  services.ResponseCodeEnum(receipt.Status),
		TopicSequenceNumber:     receipt.TopicSequenceNumber,
		TopicRunningHash:        receipt.TopicRunningHash,
		TopicRunningHashVersion: receipt.TopicRunningHashVersion,
		NewTotalSupply:          receipt.TotalSupply,
		SerialNumbers:           receipt.SerialNumbers,
	}

	if receipt.ExchangeRate != nil {
		receiptFinal.ExchangeRate = &services.ExchangeRateSet{
			CurrentRate: receipt.ExchangeRate._ToProtobuf(),
			NextRate:    receipt.ExchangeRate._ToProtobuf(),
		}
	}

	if receipt.TopicID != nil {
		receiptFinal.TopicID = receipt.TopicID._ToProtobuf()
	}

	if receipt.FileID != nil {
		receiptFinal.FileID = receipt.FileID._ToProtobuf()
	}

	if receipt.ContractID != nil {
		receiptFinal.ContractID = receipt.ContractID._ToProtobuf()
	}

	if receipt.AccountID != nil {
		receiptFinal.AccountID = receipt.AccountID._ToProtobuf()
	}

	if receipt.TokenID != nil {
		receiptFinal.TokenID = receipt.TokenID._ToProtobuf()
	}

	if receipt.ScheduleID != nil {
		receiptFinal.ScheduleID = receipt.ScheduleID._ToProtobuf()
	}

	if receipt.ScheduledTransactionID != nil {
		receiptFinal.ScheduledTransactionID = receipt.ScheduledTransactionID._ToProtobuf()
	}

	childReceipts := make([]*services.TransactionReceipt, 0)
	if len(receipt.Children) > 0 {
		for _, r := range receipt.Children {
			childReceipts = append(childReceipts, r._ToProtobuf().GetReceipt())
		}
	}

	duplicateReceipts := make([]*services.TransactionReceipt, 0)
	if len(receipt.Duplicates) > 0 {
		for _, r := range receipt.Duplicates {
			duplicateReceipts = append(duplicateReceipts, r._ToProtobuf().GetReceipt())
		}
	}

	return &services.TransactionGetReceiptResponse{
		Receipt:                      &receiptFinal,
		ChildTransactionReceipts:     childReceipts,
		DuplicateTransactionReceipts: duplicateReceipts,
	}
}

func (receipt TransactionReceipt) ValidateStatus(shouldValidate bool) error {
	if shouldValidate && receipt.Status != StatusSuccess {
		if receipt.TransactionID != nil {
			return _NewErrHederaReceiptStatus(*receipt.TransactionID, receipt.Status)
		}
		return _NewErrHederaReceiptStatus(TransactionID{}, receipt.Status)
	}

	return nil
}

func (receipt TransactionReceipt) ToBytes() []byte {
	data, err := protobuf.Marshal(receipt._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TransactionReceiptFromBytes(data []byte) (TransactionReceipt, error) {
	if data == nil {
		return TransactionReceipt{}, errByteArrayNull
	}
	pb := services.TransactionGetReceiptResponse{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionReceipt{}, err
	}

	return _TransactionReceiptFromProtobuf(&pb, nil), nil
}
