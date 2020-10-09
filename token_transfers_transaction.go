package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenTransfersTransaction struct {
	TransactionBuilder
	pb             *proto.TokenTransfersTransactionBody
	tokenIdIndexes map[string]int
}

// NewTokenTransfersTransaction creates a TokenTransfersTransaction builder which can be
// used to construct and execute a Token Transfers Transaction.
func NewTokenTransfersTransaction() TokenTransfersTransaction {
	pb := &proto.TokenTransfersTransactionBody{
		TokenTransfers: make([]*proto.TokenTransferList, 0),
	}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenTransfers{TokenTransfers: pb}

	builder := TokenTransfersTransaction{inner, pb, make(map[string]int)}

	return builder
}

func (builder TokenTransfersTransaction) AddSender(tokenID TokenID, accountID AccountID, amount Hbar) TokenTransfersTransaction {
	return builder.AddTransfers(tokenID, accountID, amount.negated())
}

func (builder TokenTransfersTransaction) AddRecipient(tokenID TokenID, accountID AccountID, amount Hbar) TokenTransfersTransaction {
	return builder.AddTransfers(tokenID, accountID, amount)
}

func (builder TokenTransfersTransaction) AddTransfers(tokenID TokenID, accountID AccountID, amount Hbar) TokenTransfersTransaction {
	index, ok := builder.tokenIdIndexes[tokenID.String()]

	if !ok {
		index := len(builder.pb.TokenTransfers)
		builder.tokenIdIndexes[tokenID.String()] = index

		builder.pb.TokenTransfers = append(builder.pb.TokenTransfers, &proto.TokenTransferList{
			Token:     tokenID.toProto(),
			Transfers: make([]*proto.AccountAmount, 0),
		})

	}

	builder.pb.TokenTransfers[index].Transfers = append(builder.pb.TokenTransfers[index].Transfers, &proto.AccountAmount{
		AccountID: accountID.toProto(),
		Amount:    amount.AsTinybar(),
	})

	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenTransfersTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenTransfersTransaction {
	return TokenTransfersTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb, builder.tokenIdIndexes}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenTransfersTransaction) SetTransactionMemo(memo string) TokenTransfersTransaction {
	return TokenTransfersTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb, builder.tokenIdIndexes}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenTransfersTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenTransfersTransaction {
	return TokenTransfersTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb, builder.tokenIdIndexes}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenTransfersTransaction) SetTransactionID(transactionID TransactionID) TokenTransfersTransaction {
	return TokenTransfersTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb, builder.tokenIdIndexes}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder TokenTransfersTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenTransfersTransaction {
	return TokenTransfersTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb, builder.tokenIdIndexes}
}
