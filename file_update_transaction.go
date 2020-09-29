package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FileUpdateTransaction struct {
	TransactionBuilder
	pb *proto.FileUpdateTransactionBody
}

func NewFileUpdateTransaction() FileUpdateTransaction {
	pb := &proto.FileUpdateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileUpdate{FileUpdate: pb}

	transaction := FileUpdateTransaction{inner, pb}

	return transaction
}

func (transaction FileUpdateTransaction) SetFileID(id FileID) FileUpdateTransaction {
	transaction.pb.FileID = id.toProto()
	return transaction
}

func (transaction FileUpdateTransaction) AddKey(publicKey PublicKey) FileUpdateTransaction {
	if transaction.pb.Keys == nil {
		transaction.pb.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}

	transaction.pb.Keys.Keys = append(transaction.pb.Keys.Keys, publicKey.toProto())

	return transaction
}

func (transaction FileUpdateTransaction) SetExpirationTime(expiration time.Time) FileUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

func (transaction FileUpdateTransaction) SetContents(contents []byte) FileUpdateTransaction {
	transaction.pb.Contents = contents
	return transaction
}

func (transaction FileUpdateTransaction) Build(client *Client) (Transaction, error) {
	return transaction.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction FileUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileUpdateTransaction {
	return FileUpdateTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction FileUpdateTransaction) SetTransactionMemo(memo string) FileUpdateTransaction {
	return FileUpdateTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction FileUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) FileUpdateTransaction {
	return FileUpdateTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction FileUpdateTransaction) SetTransactionID(transactionID TransactionID) FileUpdateTransaction {
	return FileUpdateTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction FileUpdateTransaction) SetNodeID(nodeAccountID AccountID) FileUpdateTransaction {
	return FileUpdateTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
