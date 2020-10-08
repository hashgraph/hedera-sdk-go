package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type SystemDeleteTransaction struct {
	Transaction
	pb *proto.SystemDeleteTransactionBody
}

func NewSystemDeleteTransaction() *SystemDeleteTransaction {
	pb := &proto.SystemDeleteTransactionBody{}

	transaction := SystemDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *SystemDeleteTransaction) SetExpirationTime(expiration time.Time) *SystemDeleteTransaction {
	transaction.pb.ExpirationTime = &proto.TimestampSeconds{
		Seconds: expiration.Unix(),
	}
	return transaction
}

func (transaction *SystemDeleteTransaction) GetExpirationTime() int64 {
	return transaction.pb.GetExpirationTime().Seconds
}

func (transaction *SystemDeleteTransaction) SetContractID(contractId ContractID) *SystemDeleteTransaction {
	transaction.pb.Id = &proto.SystemDeleteTransactionBody_ContractID{ContractID: contractId.toProto()}
	return transaction
}

func (transaction *SystemDeleteTransaction) GetContract() ContractID {
	return contractIDFromProto(transaction.pb.GetContractID())
}

func (transaction *SystemDeleteTransaction) SetFileID(fileId FileID) *SystemDeleteTransaction {
	transaction.pb.Id = &proto.SystemDeleteTransactionBody_FileID{FileID: fileId.toProto()}
	return transaction
}

func (transaction *SystemDeleteTransaction) GetFileID() FileID {
	return fileIDFromProto(transaction.pb.GetFileID())
}
