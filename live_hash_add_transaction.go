package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type LiveHashAddTransaction struct {
	Transaction
	pb *proto.CryptoAddLiveHashTransactionBody
}

func NewLiveHashAddTransaction() *LiveHashAddTransaction {
	pb := &proto.CryptoAddLiveHashTransactionBody{}

	transaction := LiveHashAddTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func (transaction *LiveHashAddTransaction) SetHash(hash []byte) *LiveHashAddTransaction {
	transaction.pb.LiveHash.Hash = hash
	return transaction
}

func (transaction *LiveHashAddTransaction) GetHash() []byte {
	return transaction.pb.GetLiveHash().GetHash()
}

func (transaction *LiveHashAddTransaction) SetKeys(keyList KeyList) *LiveHashAddTransaction {
	transaction.pb.LiveHash.Keys = keyList.toProtoKeyList()
	return transaction
}

func (transaction *LiveHashAddTransaction) GetKeys() KeyList {
	return keyListFromProto(transaction.pb.GetLiveHash().GetKeys())
}

func (transaction *LiveHashAddTransaction) SetDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction.pb.LiveHash.Duration = durationToProto(duration)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetDuration() time.Duration {
	return durationFromProto(transaction.pb.GetLiveHash().GetDuration())
}

func (transaction *LiveHashAddTransaction ) SetAccountId(accountId AccountID) *LiveHashAddTransaction {
	transaction.pb.LiveHash.AccountId = accountId.toProtobuf()
	return transaction
}

func (transaction *LiveHashAddTransaction) GetAccountId() AccountID {
	return accountIDFromProto(transaction.pb.LiveHash.GetAccountId())
}
