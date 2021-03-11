package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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
}

func newTransactionReceipt(
	status Status, exchangeRate *ExchangeRate,
	topicID TopicID, fileID FileID,
	contractID ContractID, accountID AccountID,
	topicSequenceNumber uint64, topicRunningHash []byte,
	topicRunningHashVersion uint64, totalSupply uint64, scheduleId ScheduleID,
	scheduledTransactionID TransactionID, tokenID TokenID) TransactionReceipt {

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
	}

	return receipt

}

func transactionReceiptFromProtobuf(protoReceipt *proto.TransactionReceipt) TransactionReceipt {
	var accountID *AccountID
	if protoReceipt.AccountID != nil {
		accountIDValue := accountIDFromProtobuf(protoReceipt.AccountID)
		accountID = &accountIDValue
	}

	var contractID *ContractID
	if protoReceipt.ContractID != nil {
		contractIDValue := contractIDFromProtobuf(protoReceipt.ContractID)
		contractID = &contractIDValue
	}

	var fileID *FileID
	if protoReceipt.FileID != nil {
		fileIDValue := fileIDFromProtobuf(protoReceipt.FileID)
		fileID = &fileIDValue
	}

	var topicID *TopicID
	if protoReceipt.TopicID != nil {
		topicIDValue := topicIDFromProtobuf(protoReceipt.TopicID)
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
		id := tokenIDFromProtobuf(protoReceipt.TokenID)
		tokenID = &id
	}

	var scheduleID *ScheduleID
	if protoReceipt.ScheduleID != nil {
		scheduleIDValue := scheduleIDFromProtobuf(protoReceipt.ScheduleID)
		scheduleID = &scheduleIDValue
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
	pb := proto.TransactionReceipt{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionReceipt{}, err
	}

	return transactionReceiptFromProtobuf(&pb), nil
}
