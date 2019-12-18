package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionReceipt struct {
	Status                       Status
	accountID                    *AccountID
	contractID                   *ContractID
	fileID                       *FileID
	consensusTopicID             *ConsensusTopicID
	ConsensusTopicSequenceNumber uint64
	ConsensusTopicRunningHash    []byte
}

func (receipt TransactionReceipt) GetFileID() FileID {
	return *receipt.fileID
}

func (receipt TransactionReceipt) GetAccountID() AccountID {
	return *receipt.accountID
}

func (receipt TransactionReceipt) GetContractID() ContractID {
	return *receipt.contractID
}

func (receipt TransactionReceipt) ConsensusTopicID() ConsensusTopicID {
	return *receipt.consensusTopicID
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

	return TransactionReceipt{
		Status:                       Status(pb.Status),
		accountID:                    accountID,
		contractID:                   contractID,
		fileID:                       fileID,
		consensusTopicID:             consensusTopicID,
		ConsensusTopicSequenceNumber: pb.TopicSequenceNumber,
		ConsensusTopicRunningHash:    pb.TopicRunningHash,
	}
}
