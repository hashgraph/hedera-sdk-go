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

// NewAccountUpdateTransaction creates an AccountUpdateTransaction builder which can be used to construct and
// execute a Crypto Update Transaction.
func NewAccountUpdateTransaction() AccountUpdateTransaction {
	pb := &proto.CryptoUpdateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_CryptoUpdateAccount{CryptoUpdateAccount: pb}

	builder := AccountUpdateTransaction{inner, pb}

	return builder
}

// SetAccountID sets the account ID which is being updated in this transaction
func (builder AccountUpdateTransaction) SetAccountID(id AccountID) AccountUpdateTransaction {
	builder.pb.AccountIDToUpdate = id.toProto()
	return builder
}

// SetKey sets the new key for the account being updated. The transaction must be signed by both the old key (from
// before the change) and the new key.
//
//The old key must sign for security. The new key must sign as a safeguard to avoid accidentally changing to an invalid
// key, and then having no way to recover.
func (builder AccountUpdateTransaction) SetKey(publicKey PublicKey) AccountUpdateTransaction {
	builder.pb.Key = publicKey.toProto()
	return builder
}

// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is unset, is an
// invalid account, or is an account that isn't a node, then this account is automatically proxy staked to a node chosen
// by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking, or if it
// is not currently running a node, then it will behave as if proxyAccountID was unset.
func (builder AccountUpdateTransaction) SetProxyAccountID(id AccountID) AccountUpdateTransaction {
	builder.pb.ProxyAccountID = id.toProto()
	return builder
}

// SetAutoRenewPeriod sets the duration in which it will automatically extend the expiration period. If it doesn't have
// enough balance, it extends as long as possible. If the balance is empty when it expires, then it is deleted.
func (builder AccountUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) AccountUpdateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return builder
}

// SetExpirationTime sets the new expiration time to extend to (ignored if equal to or before the current one) When
// extending the expiration date, the cost is affected by the size of the list of attached claims, and of the keys
// associated with the claims and the account.
func (builder AccountUpdateTransaction) SetExpirationTime(expiration time.Time) AccountUpdateTransaction {
	builder.pb.ExpirationTime = timeToProto(expiration)
	return builder
}

// SetReceiverSignatureRequired sets the receiverSigRequired flag on the account.
func (builder AccountUpdateTransaction) SetReceiverSignatureRequired(required bool) AccountUpdateTransaction {
	builder.pb.ReceiverSigRequiredField = &proto.CryptoUpdateTransactionBody_ReceiverSigRequired{
		ReceiverSigRequired: required,
	}
	return builder
}

// SetSendRecordThreshold sets the threshold amount for which an account record is created for any send/withdraw
// transaction
func (builder AccountUpdateTransaction) SetSendRecordThreshold(threshold Hbar) AccountUpdateTransaction {
	builder.pb.SendRecordThresholdField = &proto.CryptoUpdateTransactionBody_SendRecordThreshold{
		SendRecordThreshold: uint64(threshold.AsTinybar()),
	}
	return builder
}

// SetReceiveRecordThreshold sets the threshold amount for which an account record is created for any receive/deposit
// transaction
func (builder AccountUpdateTransaction) SetReceiveRecordThreshold(threshold Hbar) AccountUpdateTransaction {
	builder.pb.ReceiveRecordThresholdField = &proto.CryptoUpdateTransactionBody_ReceiveRecordThreshold{
		ReceiveRecordThreshold: uint64(threshold.AsTinybar()),
	}
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder AccountUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) AccountUpdateTransaction {
	return AccountUpdateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder AccountUpdateTransaction) SetTransactionMemo(memo string) AccountUpdateTransaction {
	return AccountUpdateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder AccountUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) AccountUpdateTransaction {
	return AccountUpdateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder AccountUpdateTransaction) SetTransactionID(transactionID TransactionID) AccountUpdateTransaction {
	return AccountUpdateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder AccountUpdateTransaction) SetNodeAccountID(nodeAccountID AccountID) AccountUpdateTransaction {
	return AccountUpdateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
