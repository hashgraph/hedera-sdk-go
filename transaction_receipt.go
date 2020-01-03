package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionReceipt struct {
	Status     Status
	accountID  *AccountID
	contractID *ContractID
	fileID     *FileID
}

func (receipt TransactionReceipt) FileID() FileID {
	return *receipt.fileID
}

func (receipt TransactionReceipt) AccountID() AccountID {
	return *receipt.accountID
}

func (receipt TransactionReceipt) ContractID() ContractID {
	return *receipt.contractID
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

	return TransactionReceipt{
		Status:     Status(pb.Status),
		accountID:  accountID,
		contractID: contractID,
		fileID:     fileID,
	}
}
