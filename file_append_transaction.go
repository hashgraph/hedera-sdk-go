package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type FileAppendTransaction struct {
	TransactionBuilder
	pb *proto.FileAppendTransactionBody
}

func NewFileAppendTransaction() FileAppendTransaction {
	pb := &proto.FileAppendTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileAppend{FileAppend: pb}

	builder := FileAppendTransaction{inner, pb}

	return builder
}

func (builder FileAppendTransaction) SetFileID(id FileID) FileAppendTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder FileAppendTransaction) SetContents(contents []byte) FileAppendTransaction {
	builder.pb.Contents = contents
	return builder
}

func (builder FileAppendTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder FileAppendTransaction) SetMaxTransactionFee(maxTransactionFee uint64) FileAppendTransaction {
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
