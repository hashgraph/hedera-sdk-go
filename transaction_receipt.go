package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type TransactionReceipt struct {
	Status                       Status
	AccountID                    *AccountID
	ContractID                   *ContractID
	FileID                       *FileID
	ConsensusTopicID             *ConsensusTopicID
	ConsensusTopicSequenceNumber uint64
	ConsensusTopicRunningHash    []byte
	ExpirationTime               *time.Time
}

func transactionReceiptFromResponse(response *proto.Response) TransactionReceipt {
	return transactionReceiptFromProto(response.GetTransactionGetReceipt().Receipt)
}

func transactionReceiptFromProto(pb *proto.TransactionReceipt) TransactionReceipt {
	var accountID *AccountID
	if pb.AccountID != nil {
		accountIDValue := accountIDFromProto(pb.AccountID)
		accountID = &accountIDValue
	}

	var contractID *ContractID
	if pb.ContractID != nil {
		contractIDValue := contractIDFromProto(pb.ContractID)
		contractID = &contractIDValue
	}

	var fileID *FileID
	if pb.FileID != nil {
		fileIDValue := fileIDFromProto(pb.FileID)
		fileID = &fileIDValue
	}

	var consensusTopicID *ConsensusTopicID
	if pb.TopicID != nil {
		consensusTopicIDValue := consensusTopicIDFromProto(pb.TopicID)
		consensusTopicID = &consensusTopicIDValue
	}

	var expirationTime *time.Time
	if pb.ExpirationTime != nil {
		expirationTimeValue := timeFromProto(pb.ExpirationTime)
		expirationTime = &expirationTimeValue
	}

	return TransactionReceipt{
		Status:                       Status(pb.Status),
		AccountID:                    accountID,
		ContractID:                   contractID,
		FileID:                       fileID,
		ConsensusTopicID:             consensusTopicID,
		ConsensusTopicSequenceNumber: pb.TopicSequenceNumber,
		ConsensusTopicRunningHash:    pb.TopicRunningHash,
		ExpirationTime:               expirationTime,
	}
}
