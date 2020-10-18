package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenAssociateTransaction struct {
	TransactionBuilder
	pb *proto.TokenAssociateTransactionBody
}

func NewTokenAssociateTransaction() TokenAssociateTransaction {
	pb := &proto.TokenAssociateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenAssociate{TokenAssociate: pb}

	builder := TokenAssociateTransaction{inner, pb}

	return builder
}

// The account to be associated with the provided tokens
func (builder TokenAssociateTransaction) SetAccountID(id AccountID) TokenAssociateTransaction {
	builder.pb.Account = id.toProto()
	return builder
}

// The tokens to be associated with the provided account
func (builder TokenAssociateTransaction) SetTokenIDs(ids ...TokenID) TokenAssociateTransaction {
	tokens := make([]*proto.TokenID, len(ids))
	for i, token := range ids {
		tokens[i] = token.toProto()
	}
	builder.pb.Tokens = tokens
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenAssociateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenAssociateTransaction {
	return TokenAssociateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenAssociateTransaction) SetTransactionMemo(memo string) TokenAssociateTransaction {
	return TokenAssociateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenAssociateTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenAssociateTransaction {
	return TokenAssociateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenAssociateTransaction) SetTransactionID(transactionID TransactionID) TokenAssociateTransaction {
	return TokenAssociateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenAssociateTransaction) SetNodeTokenID(nodeAccountID AccountID) TokenAssociateTransaction {
	return TokenAssociateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
