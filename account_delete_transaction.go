package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// AccountDeleteTransaction marks an account as deleted, moving all its current hbars to another account. It will remain
// in the ledger, marked as deleted, until it expires. Transfers into it a deleted account fail. But a deleted account
// can still have its expiration extended in the normal way.
type AccountDeleteTransaction struct {
	TransactionBuilder
	pb *proto.CryptoDeleteTransactionBody
}

// NewAccountDeleteTransaction creates an AccountDeleteTransaction transaction which can be used to construct and execute
// a Crypto Delete Transaction.
func NewAccountDeleteTransaction() AccountDeleteTransaction {
	pb := &proto.CryptoDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_CryptoDelete{CryptoDelete: pb}

	transaction := AccountDeleteTransaction{inner, pb}

	return transaction
}

// SetDeleteAccountID sets the account ID which should be deleted.
func (transaction AccountDeleteTransaction) SetDeleteAccountID(id AccountID) AccountDeleteTransaction {
	transaction.pb.DeleteAccountID = id.toProto()
	return transaction
}

// SetTransferAccountID sets the account ID which will receive all remaining hbars.
func (transaction AccountDeleteTransaction) SetTransferAccountID(id AccountID) AccountDeleteTransaction {
	transaction.pb.TransferAccountID = id.toProto()
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction AccountDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) AccountDeleteTransaction {
	return AccountDeleteTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction AccountDeleteTransaction) SetTransactionMemo(memo string) AccountDeleteTransaction {
	return AccountDeleteTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction AccountDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) AccountDeleteTransaction {
	return AccountDeleteTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction AccountDeleteTransaction) SetTransactionID(transactionID TransactionID) AccountDeleteTransaction {
	return AccountDeleteTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction AccountDeleteTransaction) SetNodeID(nodeAccountID AccountID) AccountDeleteTransaction {
	return AccountDeleteTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
