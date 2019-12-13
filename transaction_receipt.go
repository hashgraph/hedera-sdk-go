package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type TransactionReceipt struct {
	// TODO: Make the status enum look nicer in Go
	Status proto.ResponseCodeEnum

	AccountID                    *AccountID
	ContractID                   *ContractID
	FileID                       *FileID
	ConsensusTopicID             *ConsensusTopicID
	ConsensusTopicSequenceNumber uint64
	ConsensusTopicRunningHash    []byte
	ExpirationTime               *time.Time
}

func transactionReceiptFromResponse(response *proto.Response) TransactionReceipt {
	pb := response.GetTransactionGetReceipt()

	var accountID *AccountID
	if pb.Receipt.AccountID != nil {
		accountIDValue := accountIDFromProto(pb.Receipt.AccountID)
		accountID = &accountIDValue
	}

	var contractID *ContractID
	if pb.Receipt.ContractID != nil {
		contractIDValue := contractIDFromProto(pb.Receipt.ContractID)
		contractID = &contractIDValue
	}

	var fileID *FileID
	if pb.Receipt.FileID != nil {
		fileIDValue := fileIDFromProto(pb.Receipt.FileID)
		fileID = &fileIDValue
	}

	var consensusTopicID *ConsensusTopicID
	if pb.Receipt.TopicID != nil {
		consensusTopicIDValue := consensusTopicIDFromProto(pb.Receipt.TopicID)
		consensusTopicID = &consensusTopicIDValue
	}

	var expirationTime *time.Time
	if pb.Receipt.ExpirationTime != nil {
		expirationTimeValue := timeFromProto(pb.Receipt.ExpirationTime)
		expirationTime = &expirationTimeValue
	}

	return TransactionReceipt{
		Status:                       pb.Receipt.Status,
		AccountID:                    accountID,
		ContractID:                   contractID,
		FileID:                       fileID,
		ConsensusTopicID:             consensusTopicID,
		ConsensusTopicSequenceNumber: pb.Receipt.TopicSequenceNumber,
		ConsensusTopicRunningHash:    pb.Receipt.TopicRunningHash,
		ExpirationTime:               expirationTime,
	}
}
