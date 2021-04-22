package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenMintTransaction struct {
	TransactionBuilder
	pb *proto.TokenMintTransactionBody
}

func NewTokenMintTransaction() TokenMintTransaction {
	pb := &proto.TokenMintTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenMint{TokenMint: pb}

	builder := TokenMintTransaction{inner, pb}

	return builder
}

func tokenMintTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) TokenMintTransaction {
	return TokenMintTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetTokenMint(),
	}
}

// The amount to mint to the Treasury Account. Amount must be a positive non-zero number represented in the lowest denomination of the token. The new supply must be lower than 2^63.
func (builder TokenMintTransaction) SetAmount(amount uint64) TokenMintTransaction {
	builder.pb.Amount = amount
	return builder
}

// The token for which to mint tokens. If token does not exist, transaction results in INVALID_TOKEN_ID
func (builder TokenMintTransaction) SetTokenID(id TokenID) TokenMintTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

func (builder TokenMintTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *TokenMintTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenMint{
			TokenMint: &proto.TokenMintTransactionBody{
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
func (builder TokenMintTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenMintTransaction {
	return TokenMintTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenMintTransaction) SetTransactionMemo(memo string) TokenMintTransaction {
	return TokenMintTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenMintTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenMintTransaction {
	return TokenMintTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenMintTransaction) SetTransactionID(transactionID TransactionID) TokenMintTransaction {
	return TokenMintTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenMintTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenMintTransaction {
	return TokenMintTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
