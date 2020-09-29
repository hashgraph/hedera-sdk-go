package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// AccountUpdateTransaction changes properties for the given account. Any unset field is left unchanged. This
// transaction must be signed by the existing key for this account.
type AccountUpdateTransaction struct {
	TransactionBuilder
	pb *proto.CryptoUpdateTransactionBody
}

// NewAccountUpdateTransaction creates an AccountUpdateTransaction transaction which can be used to construct and
// execute a Crypto Update Transaction.
func NewAccountUpdateTransaction() AccountUpdateTransaction {
	pb := &proto.CryptoUpdateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_CryptoUpdateAccount{CryptoUpdateAccount: pb}

	transaction := AccountUpdateTransaction{inner, pb}

	return transaction
}

// SetAccountID sets the account ID which is being updated in this transaction
func (transaction AccountUpdateTransaction) SetAccountID(id AccountID) AccountUpdateTransaction {
	transaction.pb.AccountIDToUpdate = id.toProto()
	return transaction
}

// SetKey sets the new key for the account being updated. The transaction must be signed by both the old key (from
// before the change) and the new key.
//
//The old key must sign for security. The new key must sign as a safeguard to avoid accidentally changing to an invalid
// key, and then having no way to recover.
func (transaction AccountUpdateTransaction) SetKey(publicKey PublicKey) AccountUpdateTransaction {
	transaction.pb.Key = publicKey.toProto()
	return transaction
}

// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is unset, is an
// invalid account, or is an account that isn't a node, then this account is automatically proxy staked to a node chosen
// by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking, or if it
// is not currently running a node, then it will behave as if proxyAccountID was unset.
func (transaction AccountUpdateTransaction) SetProxyAccountID(id AccountID) AccountUpdateTransaction {
	transaction.pb.ProxyAccountID = id.toProto()
	return transaction
}

// SetAutoRenewPeriod sets the duration in which it will automatically extend the expiration period. If it doesn't have
// enough balance, it extends as long as possible. If the balance is empty when it expires, then it is deleted.
func (transaction AccountUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) AccountUpdateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return transaction
}

// SetExpirationTime sets the new expiration time to extend to (ignored if equal to or before the current one) When
// extending the expiration date, the cost is affected by the size of the list of attached claims, and of the keys
// associated with the claims and the account.
func (transaction AccountUpdateTransaction) SetExpirationTime(expiration time.Time) AccountUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

// SetReceiverSignatureRequired sets the receiverSigRequired flag on the account.
func (transaction AccountUpdateTransaction) SetReceiverSignatureRequired(required bool) AccountUpdateTransaction {
	transaction.pb.ReceiverSigRequiredField = &proto.CryptoUpdateTransactionBody_ReceiverSigRequired{
		ReceiverSigRequired: required,
	}
	return transaction
}

// SetSendRecordThreshold sets the threshold amount for which an account record is created for any send/withdraw
// transaction
//
// Deprecated: No longer used by Hedera
func (transaction AccountUpdateTransaction) SetSendRecordThreshold(threshold Hbar) AccountUpdateTransaction {
	transaction.pb.SendRecordThresholdField = &proto.CryptoUpdateTransactionBody_SendRecordThreshold{
		SendRecordThreshold: uint64(threshold.AsTinybar()),
	}
	return transaction
}

// SetReceiveRecordThreshold sets the threshold amount for which an account record is created for any receive/deposit
// transaction
//
// Deprecated: No longer used by Hedera
func (transaction AccountUpdateTransaction) SetReceiveRecordThreshold(threshold Hbar) AccountUpdateTransaction {
	transaction.pb.ReceiveRecordThresholdField = &proto.CryptoUpdateTransactionBody_ReceiveRecordThreshold{
		ReceiveRecordThreshold: uint64(threshold.AsTinybar()),
	}
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction AccountUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) AccountUpdateTransaction {
	return AccountUpdateTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction AccountUpdateTransaction) SetTransactionMemo(memo string) AccountUpdateTransaction {
	return AccountUpdateTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction AccountUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) AccountUpdateTransaction {
	return AccountUpdateTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction AccountUpdateTransaction) SetTransactionID(transactionID TransactionID) AccountUpdateTransaction {
	return AccountUpdateTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction AccountUpdateTransaction) SetNodeID(nodeAccountID AccountID) AccountUpdateTransaction {
	return AccountUpdateTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
