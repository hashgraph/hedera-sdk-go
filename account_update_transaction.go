package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type AccountUpdateTransaction struct {
	Transaction
	pb *proto.CryptoUpdateTransactionBody
}

func NewAccountUpdateTransaction() *AccountUpdateTransaction {
	pb := &proto.CryptoUpdateTransactionBody{}

	transaction := AccountUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func (transaction *AccountUpdateTransaction) SetKey(publicKey PublicKey) *AccountUpdateTransaction {
	transaction.pb.Key = publicKey.toProtoKey()
	return transaction
}

func (transaction *AccountUpdateTransaction) GetKey() (Key, error) {
	return publicKeyFromProto(transaction.pb.GetKey())
}

func (transaction *AccountUpdateTransaction ) SetAccountId(accountId AccountID) *AccountUpdateTransaction {
	transaction.pb.AccountIDToUpdate = accountId.toProtobuf()
	return transaction
}

func (transaction *AccountUpdateTransaction) GetAccountId() AccountID {
	return accountIDFromProto(transaction.pb.GetAccountIDToUpdate())
}

func (transaction *AccountUpdateTransaction ) SetReceiverSignatureRequired(receiverSignatureRequired bool) *AccountUpdateTransaction {
	transaction.pb.GetReceiverSigRequiredWrapper().Value = receiverSignatureRequired
	return transaction
}

func (transaction *AccountUpdateTransaction) GetReceiverSignatureRequired()  bool {
	return transaction.pb.GetReceiverSigRequiredWrapper().GetValue()
}

func (transaction *AccountUpdateTransaction ) SetProxyAccountId(proxyAccountId AccountID) *AccountUpdateTransaction {
	transaction.pb.ProxyAccountID = proxyAccountId.toProtobuf()
	return transaction
}

func (transaction *AccountUpdateTransaction) GetProxyAccountId() AccountID {
	return accountIDFromProto(transaction.pb.GetProxyAccountID())
}

func (transaction *AccountUpdateTransaction ) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountUpdateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProto(transaction.pb.GetAutoRenewPeriod())
}

func (transaction *AccountUpdateTransaction ) SetExpirationTime(expirationTime time.Time) *AccountUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expirationTime)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetExpirationTime() time.Time {
	return timeFromProto(transaction.pb.ExpirationTime)
}
