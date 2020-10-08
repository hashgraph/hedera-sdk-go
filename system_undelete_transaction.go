package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type SystemUndeleteTransaction struct {
	Transaction
	pb *proto.SystemUndeleteTransactionBody
}

func NewSystemUndeleteTransaction() *SystemUndeleteTransaction {
	pb := &proto.SystemUndeleteTransactionBody{}

	transaction := SystemUndeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *SystemUndeleteTransaction) SetContractID(contractId ContractID) *SystemUndeleteTransaction {
	transaction.pb.Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: contractId.toProto()}
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetContract() ContractID {
	return contractIDFromProto(transaction.pb.GetContractID())
}

func (transaction *SystemUndeleteTransaction) SetFileID(fileId FileID) *SystemUndeleteTransaction {
	transaction.pb.Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: fileId.toProto()}
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetFileID() FileID {
	return fileIDFromProto(transaction.pb.GetFileID())
}
