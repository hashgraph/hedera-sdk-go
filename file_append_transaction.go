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

// NewFileAppendTransaction creates a FileAppendTransaction builder which can be
// used to construct and execute a File Append Transaction.
func NewFileAppendTransaction() FileAppendTransaction {
	pb := &proto.FileAppendTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileAppend{FileAppend: pb}

	builder := FileAppendTransaction{inner, pb}

	return builder
}

func fileAppendTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) FileAppendTransaction {
	return FileAppendTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetFileAppend(),
	}
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

func (builder FileAppendTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *FileAppendTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_FileAppend{
			FileAppend: &proto.FileAppendTransactionBody{
				FileID:   builder.pb.GetFileID(),
				Contents: builder.pb.GetContents(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder FileAppendTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder FileAppendTransaction) SetTransactionMemo(memo string) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder FileAppendTransaction) SetTransactionValidDuration(validDuration time.Duration) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder FileAppendTransaction) SetTransactionID(transactionID TransactionID) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder FileAppendTransaction) SetNodeAccountID(nodeAccountID AccountID) FileAppendTransaction {
	return FileAppendTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
