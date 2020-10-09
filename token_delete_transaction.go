package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenDeleteTransaction struct {
	TransactionBuilder
	pb *proto.TokenDeleteTransactionBody
}

func NewTokenDeleteTransaction() TokenDeleteTransaction {
	pb := &proto.TokenDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenDeletion{TokenDeletion: pb}

	builder := TokenDeleteTransaction{inner, pb}

	return builder
}

// The token to be deleted. If invalid token is specified, transaction will result in INVALID_TOKEN_ID
func (builder TokenDeleteTransaction) SetTokenID(id TokenID) TokenDeleteTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenDeleteTransaction {
	return TokenDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenDeleteTransaction) SetTransactionMemo(memo string) TokenDeleteTransaction {
	return TokenDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenDeleteTransaction {
	return TokenDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenDeleteTransaction) SetTransactionID(transactionID TransactionID) TokenDeleteTransaction {
	return TokenDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenDeleteTransaction) SetNodeTokenID(nodeAccountID AccountID) TokenDeleteTransaction {
	return TokenDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
