package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FileDeleteTransaction struct {
	TransactionBuilder
	pb *proto.FileDeleteTransactionBody
}

func NewFileDeleteTransaction() FileDeleteTransaction {
	pb := &proto.FileDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileDelete{FileDelete: pb}

	builder := FileDeleteTransaction{inner, pb}

	return builder
}

func fileDeleteTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) FileDeleteTransaction {
	return FileDeleteTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetFileDelete(),
	}
}

func (builder FileDeleteTransaction) SetFileID(id FileID) FileDeleteTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder FileDeleteTransaction) Build(client *Client) (Transaction, error) {
	return builder.TransactionBuilder.Build(client)
}

func (builder FileDeleteTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *FileDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_FileDelete{
			FileDelete: &proto.FileDeleteTransactionBody{
				FileID: builder.pb.GetFileID(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder FileDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder FileDeleteTransaction) SetTransactionMemo(memo string) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder FileDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder FileDeleteTransaction) SetTransactionID(transactionID TransactionID) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder FileDeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) FileDeleteTransaction {
	return FileDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
