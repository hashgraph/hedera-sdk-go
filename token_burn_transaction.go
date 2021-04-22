package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenBurnTransaction struct {
	TransactionBuilder
	pb *proto.TokenBurnTransactionBody
}

func NewTokenBurnTransaction() TokenBurnTransaction {
	pb := &proto.TokenBurnTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenBurn{TokenBurn: pb}

	builder := TokenBurnTransaction{inner, pb}

	return builder
}

func tokenBurnTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) TokenBurnTransaction {
	return TokenBurnTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetTokenBurn(),
	}
}

// The amount to burn to the Treasury Account. Amount must be a positive non-zero number represented in the lowest denomination of the token. The new supply must be lower than 2^63.
func (builder TokenBurnTransaction) SetAmount(amount uint64) TokenBurnTransaction {
	builder.pb.Amount = amount
	return builder
}

// The token for which to burn tokens. If token does not exist, transaction results in INVALID_TOKEN_ID
func (builder TokenBurnTransaction) SetTokenID(id TokenID) TokenBurnTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

func (builder TokenBurnTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *TokenBurnTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenBurn{
			TokenBurn: &proto.TokenBurnTransactionBody{
				Token:  builder.pb.GetToken(),
				Amount: builder.pb.GetAmount(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenBurnTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenBurnTransaction {
	return TokenBurnTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenBurnTransaction) SetTransactionMemo(memo string) TokenBurnTransaction {
	return TokenBurnTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenBurnTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenBurnTransaction {
	return TokenBurnTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenBurnTransaction) SetTransactionID(transactionID TransactionID) TokenBurnTransaction {
	return TokenBurnTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenBurnTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenBurnTransaction {
	return TokenBurnTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
