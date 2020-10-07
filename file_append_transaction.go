package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// FileAppendTransaction appends the given contents to the end of the file. If a file is too big to create with a single
// FileCreateTransaction, then it can be created with the first part of its contents, and then appended multiple times
// to create the entire file.
type FileAppendTransaction struct {
	Transaction
	pb *proto.FileAppendTransactionBody
}

// NewFileAppendTransaction creates a FileAppendTransaction transaction which can be
// used to construct and execute a File Append Transaction.
func NewFileAppendTransaction() *FileAppendTransaction {
	pb := &proto.FileAppendTransactionBody{}

	transaction := FileAppendTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

// SetFileID sets the FileID of the file to which the bytes are appended to.
func (transaction *FileAppendTransaction) SetFileID(id FileID) *FileAppendTransaction {
	transaction.pb.FileID = id.toProto()
	return transaction
}

func (transaction *FileAppendTransaction) GetFileID() FileID {
	return fileIDFromProto(transaction.pb.GetFileID())
}

// SetContents sets the bytes to append to the contents of the file.
func (transaction *FileAppendTransaction) SetContents(contents []byte) *FileAppendTransaction {
	transaction.pb.Contents = contents
	return transaction
}

func (transaction *FileAppendTransaction) GetContents() []byte {
	return transaction.pb.GetContents()
}

