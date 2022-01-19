package hedera

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
}

func _TransactionReceiptFromProtobuf(protoResponse *services.TransactionGetReceiptResponse) TransactionReceipt {
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
			childReceipts = append(childReceipts, _TransactionReceiptFromProtobuf(&services.TransactionGetReceiptResponse{Receipt: r}))
		}
	}

	duplicateReceipts := make([]TransactionReceipt, 0)
	if len(protoResponse.DuplicateTransactionReceipts) > 0 {
		for _, r := range protoResponse.DuplicateTransactionReceipts {
			duplicateReceipts = append(duplicateReceipts, _TransactionReceiptFromProtobuf(&services.TransactionGetReceiptResponse{Receipt: r}))
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
	}
}

func (receipt TransactionReceipt) _ToProtobuf() *services.TransactionGetReceiptResponse {
	receiptFinal := services.TransactionReceipt{
		Status:     services.ResponseCodeEnum(receipt.Status),
		AccountID:  receipt.AccountID._ToProtobuf(),
		FileID:     receipt.FileID._ToProtobuf(),
		ContractID: receipt.ContractID._ToProtobuf(),
		ExchangeRate: &services.ExchangeRateSet{
			CurrentRate: receipt.ExchangeRate._ToProtobuf(),
			NextRate:    receipt.ExchangeRate._ToProtobuf(),
		},
		TopicID:                 receipt.TopicID._ToProtobuf(),
		TopicSequenceNumber:     receipt.TopicSequenceNumber,
		TopicRunningHash:        receipt.TopicRunningHash,
		TopicRunningHashVersion: receipt.TopicRunningHashVersion,
		TokenID:                 receipt.TokenID._ToProtobuf(),
		NewTotalSupply:          receipt.TotalSupply,
		ScheduleID:              receipt.ScheduleID._ToProtobuf(),
		ScheduledTransactionID:  receipt.ScheduledTransactionID._ToProtobuf(),
		SerialNumbers:           receipt.SerialNumbers,
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

	return _TransactionReceiptFromProtobuf(&pb), nil
}
