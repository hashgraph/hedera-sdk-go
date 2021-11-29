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
}

func _TransactionReceiptFromProtobuf(protoReceipt *services.TransactionReceipt) TransactionReceipt {
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
	}
}

func (receipt TransactionReceipt) _ToProtobuf() *services.TransactionReceipt {
	return &services.TransactionReceipt{
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
	pb := services.TransactionReceipt{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionReceipt{}, err
	}

	return _TransactionReceiptFromProtobuf(&pb), nil
}
