package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type AccountDeleteTransaction struct {
	TransactionBuilder
	pb *proto.CryptoDeleteTransactionBody
}

func NewAccountDeleteTransaction() AccountDeleteTransaction {
	pb := &proto.CryptoDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_CryptoDelete{CryptoDelete: pb}

	builder := AccountDeleteTransaction{inner, pb}

	return builder
}

func (builder AccountDeleteTransaction) SetDeleteAccountID(id AccountID) AccountDeleteTransaction {
	builder.pb.DeleteAccountID = id.toProto()
	return builder
}

func (builder AccountDeleteTransaction) SetTransferAccountID(id AccountID) AccountDeleteTransaction {
	builder.pb.TransferAccountID = id.toProto()
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder AccountDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder AccountDeleteTransaction) SetTransactionMemo(memo string) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder AccountDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder AccountDeleteTransaction) SetTransactionID(transactionID TransactionID) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder AccountDeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) AccountDeleteTransaction {
	return AccountDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
