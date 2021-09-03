package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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

func transactionReceiptFromProtobuf(protoReceipt *proto.TransactionReceipt) TransactionReceipt {
	if protoReceipt == nil {
		return TransactionReceipt{}
	}
	var accountID *AccountID
	if protoReceipt.AccountID != nil {
		accountID = accountIDFromProtobuf(protoReceipt.AccountID)
	}

	var contractID *ContractID
	if protoReceipt.ContractID != nil {
		contractID = contractIDFromProtobuf(protoReceipt.ContractID)
	}

	var fileID *FileID
	if protoReceipt.FileID != nil {
		fileID = fileIDFromProtobuf(protoReceipt.FileID)
	}

	var topicID *TopicID
	if protoReceipt.TopicID != nil {
		topicID = topicIDFromProtobuf(protoReceipt.TopicID)
	}

	var rate *ExchangeRate
	if protoReceipt.ExchangeRate != nil {
		exchangeRateValue := exchangeRateFromProtobuf(protoReceipt.ExchangeRate.GetCurrentRate())
		rate = &exchangeRateValue
	}

	var topicSequenceHash []byte
	if protoReceipt.TopicRunningHash != nil {
		topicHash := protoReceipt.TopicRunningHash
		topicSequenceHash = topicHash
	}

	var tokenID *TokenID
	if protoReceipt.TokenID != nil {
		tokenID = tokenIDFromProtobuf(protoReceipt.TokenID)
	}

	var scheduleID *ScheduleID
	if protoReceipt.ScheduleID != nil {
		scheduleID = scheduleIDFromProtobuf(protoReceipt.ScheduleID)
	}

	var scheduledTransactionID *TransactionID
	if protoReceipt.ScheduledTransactionID != nil {
		scheduledTransactionIDValue := transactionIDFromProtobuf(protoReceipt.ScheduledTransactionID)
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

func (receipt TransactionReceipt) toProtobuf() *proto.TransactionReceipt {
	return &proto.TransactionReceipt{
		Status:     proto.ResponseCodeEnum(receipt.Status),
		AccountID:  receipt.AccountID.toProtobuf(),
		FileID:     receipt.FileID.toProtobuf(),
		ContractID: receipt.ContractID.toProtobuf(),
		ExchangeRate: &proto.ExchangeRateSet{
			CurrentRate: receipt.ExchangeRate.toProtobuf(),
			NextRate:    receipt.ExchangeRate.toProtobuf(),
		},
		TopicID:                 receipt.TopicID.toProtobuf(),
		TopicSequenceNumber:     receipt.TopicSequenceNumber,
		TopicRunningHash:        receipt.TopicRunningHash,
		TopicRunningHashVersion: receipt.TopicRunningHashVersion,
		TokenID:                 receipt.TokenID.toProtobuf(),
		NewTotalSupply:          receipt.TotalSupply,
		ScheduleID:              receipt.ScheduleID.toProtobuf(),
		ScheduledTransactionID:  receipt.ScheduledTransactionID.toProtobuf(),
		SerialNumbers:           receipt.SerialNumbers,
	}
}

func (receipt TransactionReceipt) ToBytes() []byte {
	data, err := protobuf.Marshal(receipt.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TransactionReceiptFromBytes(data []byte) (TransactionReceipt, error) {
	if data == nil {
		return TransactionReceipt{}, errByteArrayNull
	}
	pb := proto.TransactionReceipt{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionReceipt{}, err
	}

	return transactionReceiptFromProtobuf(&pb), nil
}
