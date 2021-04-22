package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenWipeTransaction struct {
	TransactionBuilder
	pb *proto.TokenWipeAccountTransactionBody
}

func NewTokenWipeTransaction() TokenWipeTransaction {
	pb := &proto.TokenWipeAccountTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenWipe{TokenWipe: pb}

	builder := TokenWipeTransaction{inner, pb}

	return builder
}

func tokenWipeTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) TokenWipeTransaction {
	return TokenWipeTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetTokenWipe(),
	}
}

// The amount of tokens to wipe from the specified account. Amount must be a positive non-zero number in the lowest denomination possible, not bigger than the token balance of the account (0; balance]
func (builder TokenWipeTransaction) SetAmount(amount uint64) TokenWipeTransaction {
	builder.pb.Amount = amount
	return builder
}

// The token for which the account will be wiped. If token does not exist, transaction results in INVALID_TOKEN_ID
func (builder TokenWipeTransaction) SetTokenID(id TokenID) TokenWipeTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

// The account to be wiped
func (builder TokenWipeTransaction) SetAccountID(id AccountID) TokenWipeTransaction {
	builder.pb.Account = id.toProto()
	return builder
}

func (builder TokenWipeTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *TokenWipeTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenWipe{
			TokenWipe: &proto.TokenWipeAccountTransactionBody{
				Token:   builder.pb.GetToken(),
				Account: builder.pb.GetAccount(),
				Amount:  builder.pb.GetAmount(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenWipeTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenWipeTransaction {
	return TokenWipeTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenWipeTransaction) SetTransactionMemo(memo string) TokenWipeTransaction {
	return TokenWipeTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenWipeTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenWipeTransaction {
	return TokenWipeTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenWipeTransaction) SetTransactionID(transactionID TransactionID) TokenWipeTransaction {
	return TokenWipeTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenWipeTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenWipeTransaction {
	return TokenWipeTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
