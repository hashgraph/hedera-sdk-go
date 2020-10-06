package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type CryptoTransferTransaction struct {
	Transaction
	pb *proto.CryptoTransferTransactionBody
}

func NewCryptoTransferTransaction() *CryptoTransferTransaction {
	pb := &proto.CryptoTransferTransactionBody{
		Transfers: &proto.TransferList{
			AccountAmounts: []*proto.AccountAmount{},
		},
	}

	transaction := CryptoTransferTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func (transaction *CryptoTransferTransaction) AddTransfer(accountId AccountID, amount int64) *CryptoTransferTransaction {
	transaction.pb.Transfers.AccountAmounts = append(transaction.pb.Transfers.AccountAmounts, &proto.AccountAmount{AccountID: accountId.toProtobuf(), Amount: amount})
	return transaction
}

func (transaction *CryptoTransferTransaction) AddSender(accountId AccountID, amount int64) *CryptoTransferTransaction {
	return transaction.AddTransfer(accountId, amount)
}

func (transaction *CryptoTransferTransaction) AddRecipient(accountId AccountID, amount int64) *CryptoTransferTransaction {
	return transaction.AddTransfer(accountId, amount)
}
