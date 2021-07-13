package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"
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

func newTransactionReceipt(
	status Status, exchangeRate *ExchangeRate,
	topicID TopicID, fileID FileID,
	contractID ContractID, accountID AccountID,
	topicSequenceNumber uint64, topicRunningHash []byte,
	topicRunningHashVersion uint64, totalSupply uint64, scheduleId ScheduleID,
	scheduledTransactionID TransactionID, tokenID TokenID, serialNumbers []int64) TransactionReceipt {

	receipt := TransactionReceipt{
		Status:                  status,
		ExchangeRate:            exchangeRate,
		TopicID:                 &topicID,
		FileID:                  &fileID,
		ContractID:              &contractID,
		AccountID:               &accountID,
		TokenID:                 &tokenID,
		TopicSequenceNumber:     topicSequenceNumber,
		TopicRunningHash:        topicRunningHash,
		TopicRunningHashVersion: topicRunningHashVersion,
		TotalSupply:             totalSupply,
		ScheduleID:              &scheduleId,
		ScheduledTransactionID:  &scheduledTransactionID,
		SerialNumbers:           serialNumbers,
	}

	return receipt

}

func transactionReceiptFromProtobuf(protoReceipt *services.TransactionReceipt, networkName *NetworkName) TransactionReceipt {
	if protoReceipt == nil {
		return TransactionReceipt{}
	}
	var accountID *AccountID
	if protoReceipt.AccountID != nil {
		accountIDValue := accountIDFromProtobuf(protoReceipt.AccountID, networkName)
		accountID = &accountIDValue
	}

	var contractID *ContractID
	if protoReceipt.ContractID != nil {
		contractIDValue := contractIDFromProtobuf(protoReceipt.ContractID, networkName)
		contractID = &contractIDValue
	}

	var fileID *FileID
	if protoReceipt.FileID != nil {
		fileIDValue := fileIDFromProtobuf(protoReceipt.FileID, networkName)
		fileID = &fileIDValue
	}

	var topicID *TopicID
	if protoReceipt.TopicID != nil {
		topicIDValue := topicIDFromProtobuf(protoReceipt.TopicID, networkName)
		topicID = &topicIDValue
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
		id := tokenIDFromProtobuf(protoReceipt.TokenID, networkName)
		tokenID = &id
	}

	var scheduleID *ScheduleID
	if protoReceipt.ScheduleID != nil {
		scheduleIDValue := scheduleIDFromProtobuf(protoReceipt.ScheduleID, networkName)
		scheduleID = &scheduleIDValue
	}

	var scheduledTransactionID *TransactionID
	if protoReceipt.ScheduledTransactionID != nil {
		scheduledTransactionIDValue := transactionIDFromProtobuf(protoReceipt.ScheduledTransactionID, networkName)
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

func (receipt TransactionReceipt) toProtobuf() *services.TransactionReceipt {
	return &services.TransactionReceipt{
		Status:     services.ResponseCodeEnum(receipt.Status),
		AccountID:  receipt.AccountID.toProtobuf(),
		FileID:     receipt.FileID.toProtobuf(),
		ContractID: receipt.ContractID.toProtobuf(),
		ExchangeRate: &services.ExchangeRateSet{
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
	pb := services.TransactionReceipt{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionReceipt{}, err
	}

	return transactionReceiptFromProtobuf(&pb, nil), nil
}
