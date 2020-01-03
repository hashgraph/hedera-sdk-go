package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type FileDeleteTransaction struct {
	TransactionBuilder
	pb *proto.FileDeleteTransactionBody
}

func NewFileDeleteTransaction() FileDeleteTransaction {
	pb := &proto.FileDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileDelete{pb}

	builder := FileDeleteTransaction{inner, pb}

	return builder
}

func (builder FileDeleteTransaction) SetFileID(id FileID) FileDeleteTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder FileDeleteTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder FileDeleteTransaction) SetMaxTransactionFee(maxTransactionFee uint64) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder FileDeleteTransaction) SetTransactionMemo(memo string) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder FileDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder FileDeleteTransaction) SetTransactionID(transactionID TransactionID) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder FileDeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
