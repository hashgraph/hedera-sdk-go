package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type LiveHashDeleteTransaction struct {
	Transaction
	pb *proto.CryptoDeleteLiveHashTransactionBody
}

func NewLiveHashDeleteTransaction() *LiveHashDeleteTransaction {
	pb := &proto.CryptoDeleteLiveHashTransactionBody{}

	transaction := LiveHashDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func (transaction *LiveHashDeleteTransaction) SetHash(hash []byte) *LiveHashDeleteTransaction {
	transaction.pb.LiveHashToDelete = hash
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetHash() []byte {
	return transaction.pb.GetLiveHashToDelete()
}

func (transaction *LiveHashDeleteTransaction ) SetAccountId(accountId AccountID) *LiveHashDeleteTransaction {
	transaction.pb.AccountOfLiveHash = accountId.toProtobuf()
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetAccountId() AccountID {
	return accountIDFromProto(transaction.pb.GetAccountOfLiveHash())
}
