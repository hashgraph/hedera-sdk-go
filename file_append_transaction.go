package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// FileAppendTransaction appends the given contents to the end of the file. If a file is too big to create with a single
// FileCreateTransaction, then it can be created with the first part of its contents, and then appended multiple times
// to create the entire file.
type FileAppendTransaction struct {
	TransactionBuilder
	pb *proto.FileAppendTransactionBody
}

// NewFileAppendTransaction creates a FileAppendTransaction transaction which can be
// used to construct and execute a File Append Transaction.
func NewFileAppendTransaction() FileAppendTransaction {
	pb := &proto.FileAppendTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileAppend{FileAppend: pb}

	transaction := FileAppendTransaction{inner, pb}

	return transaction
}

// SetFileID sets the FileID of the file to which the bytes are appended to.
func (transaction FileAppendTransaction) SetFileID(id FileID) FileAppendTransaction {
	transaction.pb.FileID = id.toProto()
	return transaction
}

// SetContents sets the bytes to append to the contents of the file.
func (transaction FileAppendTransaction) SetContents(contents []byte) FileAppendTransaction {
	transaction.pb.Contents = contents
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction FileAppendTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileAppendTransaction {
	return FileAppendTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction FileAppendTransaction) SetTransactionMemo(memo string) FileAppendTransaction {
	return FileAppendTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction FileAppendTransaction) SetTransactionValidDuration(validDuration time.Duration) FileAppendTransaction {
	return FileAppendTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction FileAppendTransaction) SetTransactionID(transactionID TransactionID) FileAppendTransaction {
	return FileAppendTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction FileAppendTransaction) SetNodeID(nodeAccountID AccountID) FileAppendTransaction {
	return FileAppendTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
