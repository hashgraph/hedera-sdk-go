package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type SystemDeleteTransaction struct {
	TransactionBuilder
	pb *proto.SystemDeleteTransactionBody
}

func NewSystemDeleteTransaction() SystemDeleteTransaction {
	pb := &proto.SystemDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_SystemDelete{SystemDelete: pb}

	transaction := SystemDeleteTransaction{inner, pb}

	return transaction
}

func (transaction SystemDeleteTransaction) SetExpirationTime(expiration time.Time) SystemDeleteTransaction {
	transaction.pb.ExpirationTime = &proto.TimestampSeconds{
		Seconds: expiration.Unix(),
	}
	return transaction
}

func (transaction SystemDeleteTransaction) SetContractID(ID ContractID) SystemDeleteTransaction {
	transaction.pb.Id = &proto.SystemDeleteTransactionBody_ContractID{ContractID: ID.toProto()}
	return transaction
}

func (transaction SystemDeleteTransaction) SetFileID(ID FileID) SystemDeleteTransaction {
	transaction.pb.Id = &proto.SystemDeleteTransactionBody_FileID{FileID: ID.toProto()}
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction SystemDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) SystemDeleteTransaction {
	return SystemDeleteTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction SystemDeleteTransaction) SetTransactionMemo(memo string) SystemDeleteTransaction {
	return SystemDeleteTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction SystemDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) SystemDeleteTransaction {
	return SystemDeleteTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction SystemDeleteTransaction) SetTransactionID(transactionID TransactionID) SystemDeleteTransaction {
	return SystemDeleteTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction SystemDeleteTransaction) SetNodeID(nodeAccountID AccountID) SystemDeleteTransaction {
	return SystemDeleteTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
