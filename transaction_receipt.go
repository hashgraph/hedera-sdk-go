package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionReceipt struct {
	Status Status
	ExchangeRate *ExchangeRate
	TopicID *TopicID
	FileID	*FileID
	ContractID *ContractID
	AccountID *AccountID
	TopicSequenceNumber uint64
	TopicRunningHash []byte
	TopicRunningHashVersion uint64
}

func newTransactionReceipt(
	status Status, exchangeRate *ExchangeRate,
	topicID TopicID, fileID	FileID,
	contractID ContractID, accountID AccountID,
	topicSequenceNumber uint64, topicRunningHash []byte,
	topicRunningHashVersion uint64 ) *TransactionReceipt {

	receipt := TransactionReceipt{
		Status: status,
		ExchangeRate: exchangeRate,
		TopicID: &topicID,
		FileID: &fileID,
		ContractID: &contractID,
		AccountID: &accountID,
		TopicSequenceNumber: topicSequenceNumber,
		TopicRunningHash: topicRunningHash,
		TopicRunningHashVersion: topicRunningHashVersion,
	}

	return &receipt

}

// GetFileID returns the FileID associated with the receipt's transaction or else panics no FileID exists
func (receipt *TransactionReceipt) GetFileID() FileID {
	return *receipt.FileID
}

// TryGetFileID returns the FileID associated with the receipt's transaction or else returns an err if no FileID exists.
func (receipt TransactionReceipt) TryGetFileID() (FileID, error) {
	if receipt.FileID == nil {
		return FileID{}, fmt.Errorf("no file id exists on this receipt")
	}

	return receipt.GetFileID(), nil
}

// GetAccountID returns the AccountID associated with the receipt's transaction or else panics if no AccountID exists
func (receipt TransactionReceipt) GetAccountID() AccountID {
	return *receipt.AccountID
}

// TryGetAccountID returns the AccountID associated with the receipt's transaction or else returns an error
// if no AccountID exists
func (receipt TransactionReceipt) TryGetAccountID() (AccountID, error) {
	if receipt.AccountID == nil {
		return AccountID{}, fmt.Errorf("no account id exists on this receipt")
	}

	return receipt.GetAccountID(), nil
}

// GetContractID returns the ContractID associated with the receipt's transaction or else panics if no ContractID exists
func (receipt TransactionReceipt) GetContractID() ContractID {
	return *receipt.ContractID
}

// TryGetContractID returns the ContractID associated with the receipt's transaction or else returns an error
// if no ContractID exists
func (receipt TransactionReceipt) TryGetContractID() (ContractID, error) {
	if receipt.ContractID == nil {
		return ContractID{}, fmt.Errorf("no contract id exists on this receipt")
	}

	return receipt.GetContractID(), nil
}

// GetConsensusTopicID returns the ConsensusTopicID associated with the receipt's transaction or else panics
// if no ConsensusTopicID exists
func (receipt TransactionReceipt) GetConsensusTopicID() TopicID {
	return *receipt.TopicID
}

// TryGetConsensusTopicID returns the ConsensusTopicID associated with the receipt's transaction or else
// returns an error if no ConsensusTopicID exists
func (receipt TransactionReceipt) TryGetConsensusTopicID() (TopicID, error) {
	if receipt.TopicID == nil {
		return TopicID{}, fmt.Errorf("no consensus id exists on this receipt")
	}
	return receipt.GetConsensusTopicID(), nil
}

// GetConsensusTopicSequenceNumber returns the topic sequence number associated with the
// Consensus Topic. However, if a ConsensusTopicID does not exist on the receipt it will return
// potentially invalid values.
func (receipt TransactionReceipt) GetConsensusTopicSequenceNumber() uint64 {
	return receipt.TopicSequenceNumber
}

// TryGetConsensusTopicSequenceNumber checks if the receipt contains a ConsensusTopicID. If
// the ConsensusTopicID exists it will return the ConsensusTopicSequenceNumber. Otherwise an
// error will be returned.
func (receipt TransactionReceipt) TryGetConsensusTopicSequenceNumber() (uint64, error) {
	if _, err := receipt.TryGetConsensusTopicID(); err != nil {
		return 0, err
	}
	return receipt.GetConsensusTopicSequenceNumber(), nil
}

// GetConsensusTopicRunningHash returns the running hash associated with the Consensus Topic.
// However, if a ConsensusTopicID does not exist on the receipt it will return potentially
// invalid values (likely an empty slice).
func (receipt TransactionReceipt) GetConsensusTopicRunningHash() []byte {
	return receipt.TopicRunningHash
}

// TryGetConsensusTopicRunningHash checks if the receipt contains a ConsensusTopicID. If the
// ConsensusTopicID exists it will return the running hash associated with the consensus Topic.
// Otherwise, an error will be returned.
func (receipt TransactionReceipt) TryGetConsensusTopicRunningHash() ([]byte, error) {
	if _, err := receipt.TryGetConsensusTopicID(); err != nil {
		return []byte{}, err
	}
	return receipt.TopicRunningHash, nil
}


func transactionReceiptFromProtobuf(protoReceipt proto.TransactionReceipt) TransactionReceipt{
	var accountID *AccountID
	if protoReceipt.AccountID != nil {
		accountIDValue := accountIDFromProto(protoReceipt.AccountID)
		accountID = &accountIDValue
	}

	var contractID *ContractID
	if protoReceipt.ContractID != nil {
		contractIDValue := contractIDFromProto(protoReceipt.ContractID)
		contractID = &contractIDValue
	}

	var fileID *FileID
	if protoReceipt.FileID != nil {
		fileIDValue := fileIDFromProto(protoReceipt.FileID)
		fileID = &fileIDValue
	}

	var topicID *TopicID
	if protoReceipt.TopicID != nil {
		topicIDValue := TopicIDFromProto(protoReceipt.TopicID)
		topicID = &topicIDValue
	}

	var rate *ExchangeRate
	if protoReceipt.ExchangeRate != nil {
		exchangeRateValue := exchangeRateFromProtobuf(protoReceipt.ExchangeRate.GetCurrentRate())
		rate = &exchangeRateValue
	}

	return TransactionReceipt{
		Status:              Status(protoReceipt.Status),
		ExchangeRate:      	 rate,
		TopicID:             topicID,
		FileID:              fileID,
		ContractID:          contractID,
		AccountID:           accountID,
		TopicSequenceNumber: protoReceipt.TopicSequenceNumber,
		TopicRunningHash:    protoReceipt.TopicRunningHash,
		TopicRunningHashVersion: protoReceipt.TopicRunningHashVersion,
	}
}

func (receipt *TransactionReceipt) toProtobuf() proto.TransactionReceipt{
	return proto.TransactionReceipt{
		Status:     proto.ResponseCodeEnum(receipt.Status),
		AccountID:  receipt.AccountID.toProtobuf(),
		FileID:     receipt.FileID.toProto(),
		ContractID: receipt.ContractID.toProto(),
		ExchangeRate: &proto.ExchangeRateSet{
			CurrentRate: receipt.ExchangeRate.toProtobuf(),
			NextRate:    receipt.ExchangeRate.toProtobuf(),
		},
		TopicID:                 receipt.TopicID.toProto(),
		TopicSequenceNumber:     receipt.TopicSequenceNumber,
		TopicRunningHash:        receipt.TopicRunningHash,
		TopicRunningHashVersion: receipt.TopicRunningHashVersion,
		TokenId:                 nil,
	}
}
