package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenGrantKycTransaction struct {
	TransactionBuilder
	pb *proto.TokenGrantKycTransactionBody
}

func NewTokenGrantKycTransaction() TokenGrantKycTransaction {
	pb := &proto.TokenGrantKycTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenGrantKyc{TokenGrantKyc: pb}

	builder := TokenGrantKycTransaction{inner, pb}

	return builder
}

func tokenGrantKycTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) TokenGrantKycTransaction {
	return TokenGrantKycTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetTokenGrantKyc(),
	}
}

// The token for which this account will be granted KYC. If token does not exist, transaction results in INVALID_TOKEN_ID
func (builder TokenGrantKycTransaction) SetTokenID(id TokenID) TokenGrantKycTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

// The account to be KYCed
func (builder TokenGrantKycTransaction) SetAccountID(id AccountID) TokenGrantKycTransaction {
	builder.pb.Account = id.toProto()
	return builder
}

func (builder TokenGrantKycTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *TokenGrantKycTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenGrantKyc{
			TokenGrantKyc: &proto.TokenGrantKycTransactionBody{
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
func (builder TokenGrantKycTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenGrantKycTransaction {
	return TokenGrantKycTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenGrantKycTransaction) SetTransactionMemo(memo string) TokenGrantKycTransaction {
	return TokenGrantKycTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenGrantKycTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenGrantKycTransaction {
	return TokenGrantKycTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenGrantKycTransaction) SetTransactionID(transactionID TransactionID) TokenGrantKycTransaction {
	return TokenGrantKycTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenGrantKycTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenGrantKycTransaction {
	return TokenGrantKycTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
