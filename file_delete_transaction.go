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

	transaction := FileDeleteTransaction{inner, pb}

	return transaction
}

func (transaction FileDeleteTransaction) SetFileID(id FileID) FileDeleteTransaction {
	transaction.pb.FileID = id.toProto()
	return transaction
}

func (transaction FileDeleteTransaction) Build(client *Client) (Transaction, error) {
	return transaction.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction FileDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileDeleteTransaction {
	return FileDeleteTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction FileDeleteTransaction) SetTransactionMemo(memo string) FileDeleteTransaction {
	return FileDeleteTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction FileDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) FileDeleteTransaction {
	return FileDeleteTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction FileDeleteTransaction) SetTransactionID(transactionID TransactionID) FileDeleteTransaction {
	return FileDeleteTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction FileDeleteTransaction) SetNodeID(nodeAccountID AccountID) FileDeleteTransaction {
	return FileDeleteTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
