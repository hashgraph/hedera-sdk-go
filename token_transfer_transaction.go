package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenTransferTransaction struct {
	TransactionBuilder
	pb             *proto.TokenTransfersTransactionBody
	tokenIdIndexes map[string]int
}

// NewTokenTransferTransaction creates a TokenTransferTransaction builder which can be
// used to construct and execute a Token Transfers Transaction.
func NewTokenTransferTransaction() TokenTransferTransaction {
	pb := &proto.TokenTransfersTransactionBody{
		TokenTransfers: make([]*proto.TokenTransferList, 0),
	}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenTransfers{TokenTransfers: pb}

	builder := TokenTransferTransaction{inner, pb, make(map[string]int)}

	return builder
}

func (builder TokenTransferTransaction) AddSender(tokenID TokenID, accountID AccountID, amount uint64) TokenTransferTransaction {
	return builder.AddTransfers(tokenID, accountID, -int64(amount))
}

func (builder TokenTransferTransaction) AddRecipient(tokenID TokenID, accountID AccountID, amount uint64) TokenTransferTransaction {
	return builder.AddTransfers(tokenID, accountID, int64(amount))
}

func (builder TokenTransferTransaction) AddTransfers(tokenID TokenID, accountID AccountID, amount int64) TokenTransferTransaction {
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
		Amount:    amount,
	})

	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenTransferTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenTransferTransaction {
	return TokenTransferTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb, builder.tokenIdIndexes}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenTransferTransaction) SetTransactionMemo(memo string) TokenTransferTransaction {
	return TokenTransferTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb, builder.tokenIdIndexes}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenTransferTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenTransferTransaction {
	return TokenTransferTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb, builder.tokenIdIndexes}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenTransferTransaction) SetTransactionID(transactionID TransactionID) TokenTransferTransaction {
	return TokenTransferTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb, builder.tokenIdIndexes}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder TokenTransferTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenTransferTransaction {
	return TokenTransferTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb, builder.tokenIdIndexes}
}
