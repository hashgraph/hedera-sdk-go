package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type SystemUndeleteTransaction struct {
	TransactionBuilder
	pb *proto.SystemUndeleteTransactionBody
}

func NewSystemUndeleteTransaction() SystemUndeleteTransaction {
	pb := &proto.SystemUndeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_SystemUndelete{SystemUndelete: pb}

	transaction := SystemUndeleteTransaction{inner, pb}

	return transaction
}

func (transaction SystemUndeleteTransaction) SetContractID(ID ContractID) SystemUndeleteTransaction {
	transaction.pb.Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: ID.toProto()}
	return transaction
}

func (transaction SystemUndeleteTransaction) SetFileID(ID FileID) SystemUndeleteTransaction {
	transaction.pb.Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: ID.toProto()}
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction SystemUndeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction SystemUndeleteTransaction) SetTransactionMemo(memo string) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction SystemUndeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction SystemUndeleteTransaction) SetTransactionID(transactionID TransactionID) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction SystemUndeleteTransaction) SetNodeID(nodeAccountID AccountID) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
