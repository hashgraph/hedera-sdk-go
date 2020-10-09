package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenRevokeKycTransaction struct {
	TransactionBuilder
	pb *proto.TokenRevokeKycTransactionBody
}

func NewTokenRevokeKycTransaction() TokenRevokeKycTransaction {
	pb := &proto.TokenRevokeKycTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenRevokeKyc{TokenRevokeKyc: pb}

	builder := TokenRevokeKycTransaction{inner, pb}

	return builder
}

// The token for which this account will get his KYC revoked. If token does not exist, transaction results in INVALID_TOKEN_ID
func (builder TokenRevokeKycTransaction) SetTokenID(id TokenID) TokenRevokeKycTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

// The account to be KYC Revoked
func (builder TokenRevokeKycTransaction) SetAccountID(id AccountID) TokenRevokeKycTransaction {
	builder.pb.Account = id.toProto()
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenRevokeKycTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenRevokeKycTransaction {
	return TokenRevokeKycTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenRevokeKycTransaction) SetTransactionMemo(memo string) TokenRevokeKycTransaction {
	return TokenRevokeKycTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenRevokeKycTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenRevokeKycTransaction {
	return TokenRevokeKycTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenRevokeKycTransaction) SetTransactionID(transactionID TransactionID) TokenRevokeKycTransaction {
	return TokenRevokeKycTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenRevokeKycTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenRevokeKycTransaction {
	return TokenRevokeKycTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
