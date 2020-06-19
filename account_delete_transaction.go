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

// NewAccountDeleteTransaction creates an AccountDeleteTransaction builder which can be used to construct and execute
// a Crypto Delete Transaction.
func NewAccountDeleteTransaction() AccountDeleteTransaction {
	pb := &proto.CryptoDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_CryptoDelete{CryptoDelete: pb}

	builder := AccountDeleteTransaction{inner, pb}

	return builder
}

// SetDeleteAccountID sets the account ID which should be deleted.
func (builder AccountDeleteTransaction) SetDeleteAccountID(id AccountID) AccountDeleteTransaction {
	builder.pb.DeleteAccountID = id.toProto()
	return builder
}

// SetTransferAccountID sets the account ID which will receive all remaining hbars.
func (builder AccountDeleteTransaction) SetTransferAccountID(id AccountID) AccountDeleteTransaction {
	builder.pb.TransferAccountID = id.toProto()
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder AccountDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder AccountDeleteTransaction) SetTransactionMemo(memo string) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder AccountDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder AccountDeleteTransaction) SetTransactionID(transactionID TransactionID) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder AccountDeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
