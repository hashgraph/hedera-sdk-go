package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// A CryptoTransferTransaction is for transferring hbar from some account balances to other account balances.
// The accounts list can contain up to 10 accounts.
// The amounts list must be the same length as the accounts list.
//
// This transaction must be signed by the
// keys for all the sending accounts, and for any receiving accounts that have receiverSigRequired == true. The
// signatures are in the same order as the accounts, skipping those accounts that don't need a signature.
//
// Deprecated: Use `TransferTransaction` instead
type CryptoTransferTransaction struct {
	TransactionBuilder
	pb *proto.CryptoTransferTransactionBody
}

// NewCryptoTransferTransaction creates a CryptoTransferTransaction builder which can be
// used to construct and execute a Crypto Transfer Transaction.
func NewCryptoTransferTransaction() CryptoTransferTransaction {
	pb := &proto.CryptoTransferTransactionBody{
		Transfers: &proto.TransferList{
			AccountAmounts: []*proto.AccountAmount{},
		},
	}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_CryptoTransfer{CryptoTransfer: pb}

	builder := CryptoTransferTransaction{inner, pb}

	return builder
}

func cryptoTransferTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) CryptoTransferTransaction {
	return CryptoTransferTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetCryptoTransfer(),
	}
}

// AddSender adds an account and the amount of hbar (as a positive value) to be sent from the sender. If any sender
// account fails to have a sufficient balance to do the withdrawal, then the entire transaction fails, and none of those
// transfers occur, though the transaction fee is still charged.
func (builder CryptoTransferTransaction) AddSender(id AccountID, amount Hbar) CryptoTransferTransaction {
	return builder.AddTransfer(id, amount.negated())
}

// AddRecipient adds a recipient account and the amount of hbar to be received from the sender(s).
func (builder CryptoTransferTransaction) AddRecipient(id AccountID, amount Hbar) CryptoTransferTransaction {
	return builder.AddTransfer(id, amount)
}

// AddTransfer adds the accountID to the internal accounts list and the amounts to the internal amounts list. Each
// negative amount is withdrawn from the corresponding account (a sender), and each positive one is added to the
// corresponding account (a receiver). The amounts list must sum to zero and there can be a maximum of 10 transfers.
//
// AddSender and AddRecipient are provided as convenience wrappers around AddTransfer.
func (builder CryptoTransferTransaction) AddTransfer(id AccountID, amount Hbar) CryptoTransferTransaction {
	builder.pb.Transfers.AccountAmounts = append(builder.pb.Transfers.AccountAmounts, &proto.AccountAmount{
		AccountID: id.toProto(),
		Amount:    amount.AsTinybar(),
	})

	return builder
}

func (builder *CryptoTransferTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_CryptoTransfer{
			CryptoTransfer: &proto.CryptoTransferTransactionBody{
				Transfers:      builder.pb.GetTransfers(),
				TokenTransfers: builder.pb.GetTokenTransfers(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder CryptoTransferTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder CryptoTransferTransaction) SetTransactionMemo(memo string) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder CryptoTransferTransaction) SetTransactionValidDuration(validDuration time.Duration) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder CryptoTransferTransaction) SetTransactionID(transactionID TransactionID) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder CryptoTransferTransaction) SetNodeAccountID(nodeAccountID AccountID) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
