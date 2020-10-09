package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenDissociateTransaction struct {
	TransactionBuilder
	pb *proto.TokenDissociateTransactionBody
}

func NewTokenDissociateTransaction() TokenDissociateTransaction {
	pb := &proto.TokenDissociateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenDissociate{TokenDissociate: pb}

	builder := TokenDissociateTransaction{inner, pb}

	return builder
}

// The account to be dissociated with the provided tokens
func (builder TokenDissociateTransaction) SetAccountID(id AccountID) TokenDissociateTransaction {
	builder.pb.Account = id.toProto()
	return builder
}

// The tokens to be dissociated with the provided account
func (builder TokenDissociateTransaction) AddTokenID(id TokenID) TokenDissociateTransaction {
	builder.pb.Tokens = append(builder.pb.Tokens, id.toProto())
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenDissociateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenDissociateTransaction {
	return TokenDissociateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenDissociateTransaction) SetTransactionMemo(memo string) TokenDissociateTransaction {
	return TokenDissociateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenDissociateTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenDissociateTransaction {
	return TokenDissociateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenDissociateTransaction) SetTransactionID(transactionID TransactionID) TokenDissociateTransaction {
	return TokenDissociateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenDissociateTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenDissociateTransaction {
	return TokenDissociateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
