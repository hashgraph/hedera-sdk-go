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
	pb := &proto.TokenAssociateTransactionBody{
		Tokens: make([]*proto.TokenID, 0),
	}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenAssociate{TokenAssociate: pb}

	builder := TokenAssociateTransaction{inner, pb}

	return builder
}

func tokenAssociateTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) TokenAssociateTransaction {
	return TokenAssociateTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetTokenAssociate(),
	}
}

// The account to be associated with the provided tokens
func (builder TokenAssociateTransaction) SetAccountID(id AccountID) TokenAssociateTransaction {
	builder.pb.Account = id.toProto()
	return builder
}

// The tokens to be associated with the provided account
func (builder TokenAssociateTransaction) AddTokenID(id TokenID) TokenAssociateTransaction {
	builder.pb.Tokens = append(builder.pb.Tokens, id.toProto())
	return builder
}

func (builder TokenAssociateTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *TokenAssociateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenAssociate{
			TokenAssociate: &proto.TokenAssociateTransactionBody{
				Account: builder.pb.GetAccount(),
				Tokens:  builder.pb.GetTokens(),
			},
		},
	}, nil
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
