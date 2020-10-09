package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenUnfreezeTransaction struct {
	TransactionBuilder
	pb *proto.TokenUnfreezeAccountTransactionBody
}

func NewTokenUnfreezeTransaction() TokenUnfreezeTransaction {
	pb := &proto.TokenUnfreezeAccountTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenUnfreeze{TokenUnfreeze: pb}

	builder := TokenUnfreezeTransaction{inner, pb}

	return builder
}

// The token for which this account will be unfrozen. If token does not exist, transaction results in INVALID_TOKEN_ID
func (builder TokenUnfreezeTransaction) SetTokenID(id TokenID) TokenUnfreezeTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

// The account to be unfrozen
func (builder TokenUnfreezeTransaction) SetAccountID(id AccountID) TokenUnfreezeTransaction {
	builder.pb.Account = id.toProto()
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenUnfreezeTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenUnfreezeTransaction {
	return TokenUnfreezeTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenUnfreezeTransaction) SetTransactionMemo(memo string) TokenUnfreezeTransaction {
	return TokenUnfreezeTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenUnfreezeTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenUnfreezeTransaction {
	return TokenUnfreezeTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenUnfreezeTransaction) SetTransactionID(transactionID TransactionID) TokenUnfreezeTransaction {
	return TokenUnfreezeTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenUnfreezeTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenUnfreezeTransaction {
	return TokenUnfreezeTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
