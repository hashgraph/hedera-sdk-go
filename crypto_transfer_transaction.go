package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type CryptoTransferTransaction struct {
	TransactionBuilder
	pb *proto.CryptoTransferTransactionBody
}

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

func (builder CryptoTransferTransaction) AddSender(id AccountID, amount uint64) CryptoTransferTransaction {
	return builder.AddTransfer(id, -int64(amount))
}

func (builder CryptoTransferTransaction) AddRecipient(id AccountID, amount uint64) CryptoTransferTransaction {
	return builder.AddTransfer(id, int64(amount))
}

func (builder CryptoTransferTransaction) AddTransfer(id AccountID, amount int64) CryptoTransferTransaction {
	builder.pb.Transfers.AccountAmounts = append(builder.pb.Transfers.AccountAmounts, &proto.AccountAmount{
		AccountID: id.toProto(),
		Amount:    amount,
	})

	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder CryptoTransferTransaction) SetMaxTransactionFee(maxTransactionFee uint64) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder CryptoTransferTransaction) SetTransactionMemo(memo string) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder CryptoTransferTransaction) SetTransactionValidDuration(validDuration time.Duration) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder CryptoTransferTransaction) SetTransactionID(transactionID TransactionID) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder CryptoTransferTransaction) SetNodeAccountID(nodeAccountID AccountID) CryptoTransferTransaction {
	return CryptoTransferTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
