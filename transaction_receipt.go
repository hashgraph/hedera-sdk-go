package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// TransactionReceipt is the consensus result for a transaction which is returned from a TransactionReceiptQuery.
type TransactionReceipt struct {
	// Status is the consensus status of the receipt's transaction. It might be unknown or have Failed.
	Status                       Status
	accountID                    *AccountID
	contractID                   *ContractID
	fileID                       *FileID
	consensusTopicID             *ConsensusTopicID
	consensusTopicSequenceNumber uint64
	consensusTopicRunningHash    []byte
}

// GetFileID returns the FileID associated with the receipt's transaction or else panics no FileID exists
func (receipt TransactionReceipt) GetFileID() FileID {
	return *receipt.fileID
}

// TryGetFileID returns the FileID associated with the receipt's transaction or else returns an err if no FileID exists.
func (receipt TransactionReceipt) TryGetFileID() (FileID, error) {
	if receipt.fileID == nil {
		return FileID{}, fmt.Errorf("no file id exists on this receipt")
	}

	return receipt.GetFileID(), nil
}

// GetAccountID returns the AccountID associated with the receipt's transaction or else panics if no AccountID exists
func (receipt TransactionReceipt) GetAccountID() AccountID {
	return *receipt.accountID
}

// TryGetAccountID returns the AccountID associated with the receipt's transaction or else returns an error
// if no AccountID exists
func (receipt TransactionReceipt) TryGetAccountID() (AccountID, error) {
	if receipt.accountID == nil {
		return AccountID{}, fmt.Errorf("no account id exists on this receipt")
	}

	return receipt.GetAccountID(), nil
}

// GetContractID returns the ContractID associated with the receipt's transaction or else panics if no ContractID exists
func (receipt TransactionReceipt) GetContractID() ContractID {
	return *receipt.contractID
}

// TryGetContractID returns the ContractID associated with the receipt's transaction or else returns an error
// if no ContractID exists
func (receipt TransactionReceipt) TryGetContractID() (ContractID, error) {
	if receipt.contractID == nil {
		return ContractID{}, fmt.Errorf("no contract id exists on this receipt")
	}

	return receipt.GetContractID(), nil
}

// GetConsensusTopicID returns the ConsensusTopicID associated with the receipt's transaction or else panics
// if no ConsensusTopicID exists
func (receipt TransactionReceipt) GetConsensusTopicID() ConsensusTopicID {
	return *receipt.consensusTopicID
}

// TryGetConsensusTopicID returns the ConsensusTopicID associated with the receipt's transaction or else
// returns an error if no ConsensusTopicID exists
func (receipt TransactionReceipt) TryGetConsensusTopicID() (ConsensusTopicID, error) {
	if receipt.consensusTopicID == nil {
		return ConsensusTopicID{}, fmt.Errorf("no consensus id exists on this receipt")
	}
	return receipt.GetConsensusTopicID(), nil
}

// GetConsensusTopicSequenceNumber returns the topic sequence number associated with the
// Consensus Topic. However, if a ConsensusTopicID does not exist on the receipt it will return
// potentially invalid values.
func (receipt TransactionReceipt) GetConsensusTopicSequenceNumber() uint64 {
	return receipt.consensusTopicSequenceNumber
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
	return receipt.consensusTopicRunningHash
}

// TryGetConsensusTopicRunningHash checks if the receipt contains a ConsensusTopicID. If the
// ConsensusTopicID exists it will return the running hash associated with the consensus Topic.
// Otherwise, an error will be returned.
func (receipt TransactionReceipt) TryGetConsensusTopicRunningHash() ([]byte, error) {
	if _, err := receipt.TryGetConsensusTopicID(); err != nil {
		return []byte{}, err
	}
	return receipt.consensusTopicRunningHash, nil
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
		consensusTopicSequenceNumber: pb.TopicSequenceNumber,
		consensusTopicRunningHash:    pb.TopicRunningHash,
	}
}
