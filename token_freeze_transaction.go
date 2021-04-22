package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenFreezeTransaction struct {
	TransactionBuilder
	pb *proto.TokenFreezeAccountTransactionBody
}

func NewTokenFreezeTransaction() TokenFreezeTransaction {
	pb := &proto.TokenFreezeAccountTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenFreeze{TokenFreeze: pb}

	builder := TokenFreezeTransaction{inner, pb}

	return builder
}

func tokenFreezeTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) TokenFreezeTransaction {
	return TokenFreezeTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetTokenFreeze(),
	}
}

// The token for which this account will be frozen. If token does not exist, transaction results in INVALID_TOKEN_ID
func (builder TokenFreezeTransaction) SetTokenID(id TokenID) TokenFreezeTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

// The account to be frozen
func (builder TokenFreezeTransaction) SetAccountID(id AccountID) TokenFreezeTransaction {
	builder.pb.Account = id.toProto()
	return builder
}

func (builder TokenFreezeTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *TokenFreezeTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenFreeze{
			TokenFreeze: &proto.TokenFreezeAccountTransactionBody{
				Token:   builder.pb.GetToken(),
				Account: builder.pb.GetAccount(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenFreezeTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenFreezeTransaction {
	return TokenFreezeTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenFreezeTransaction) SetTransactionMemo(memo string) TokenFreezeTransaction {
	return TokenFreezeTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenFreezeTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenFreezeTransaction {
	return TokenFreezeTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenFreezeTransaction) SetTransactionID(transactionID TransactionID) TokenFreezeTransaction {
	return TokenFreezeTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenFreezeTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenFreezeTransaction {
	return TokenFreezeTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
