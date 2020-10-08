package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FileDeleteTransaction struct {
	Transaction
	pb *proto.FileDeleteTransactionBody
}

func NewFileDeleteTransaction() *FileDeleteTransaction {
	pb := &proto.FileDeleteTransactionBody{}

	transaction := FileDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *FileDeleteTransaction) SetFileID(fileId FileID) *FileDeleteTransaction {
	transaction.pb.FileID = fileId.toProto()
	return transaction
}

func (transaction *FileDeleteTransaction) GetFileID() FileID {
	return fileIDFromProto(transaction.pb.GetFileID())
}
