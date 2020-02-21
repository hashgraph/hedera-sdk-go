package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

// FileAppendTransaction appends the given contents to the end of the file. If a file is too big to create with a single
// FileCreateTransaction, then it can be created with the first part of its contents, and then appended multiple times
// to create the entire file.
type FileAppendTransaction struct {
	TransactionBuilder
	pb *proto.FileAppendTransactionBody
}

// NewFileAppendTransaction creates a FileAppendTransaction builder which can be
// used to construct and execute a File Append Transaction.
func NewFileAppendTransaction() FileAppendTransaction {
	pb := &proto.FileAppendTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileAppend{FileAppend: pb}

	builder := FileAppendTransaction{inner, pb}

	return builder
}

// SetFileID sets the FileID of the file to which the bytes are appended to.
func (builder FileAppendTransaction) SetFileID(id FileID) FileAppendTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

// SetContents sets the bytes to append to the contents of the file.
func (builder FileAppendTransaction) SetContents(contents []byte) FileAppendTransaction {
	builder.pb.Contents = contents
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder FileAppendTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder FileAppendTransaction) SetTransactionMemo(memo string) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder FileAppendTransaction) SetTransactionValidDuration(validDuration time.Duration) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder FileAppendTransaction) SetTransactionID(transactionID TransactionID) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder FileAppendTransaction) SetNodeAccountID(nodeAccountID AccountID) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
