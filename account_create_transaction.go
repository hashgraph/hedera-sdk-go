package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// AccountCreateTransaction creates a new account. After the account is created, the AccountID for it is in the receipt,
// or by asking for a Record of the transaction to be created, and retrieving that. The account can then automatically
// generate records for large transfers into it or out of it, which each last for 25 hours. Records are generated for
// any transfer that exceeds the thresholds given here. This account is charged hbar for each record generated, so the
// thresholds are useful for limiting Record generation to happen only for large transactions.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0,
// with a null key. Future versions of the API will support multiple realms and multiple shards.
type AccountCreateTransaction struct {
	TransactionBuilder
	pb *proto.CryptoCreateTransactionBody
}

// NewAccountCreateTransaction creates an AccountCreateTransaction builder which can be used to construct and
// execute a Crypto Create Transaction.
func NewAccountCreateTransaction() AccountCreateTransaction {
	pb := &proto.CryptoCreateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_CryptoCreateAccount{CryptoCreateAccount: pb}

	builder := AccountCreateTransaction{inner, pb}
	builder.SetAutoRenewPeriod(7890000 * time.Second)

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	builder.SetReceiveRecordThreshold(MaxHbar)
	builder.SetSendRecordThreshold(MaxHbar)

	return builder
}

// SetKey sets the key that must sign each transfer out of the account. If RecieverSignatureRequired is true, then it
// must also sign any transfer into the account.
func (builder AccountCreateTransaction) SetKey(publicKey PublicKey) AccountCreateTransaction {
	builder.pb.Key = publicKey.toProto()
	return builder
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (builder AccountCreateTransaction) SetInitialBalance(initialBalance Hbar) AccountCreateTransaction {
	builder.pb.InitialBalance = uint64(initialBalance.AsTinybar())
	return builder
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (builder AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) AccountCreateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return builder
}

// SetSendRecordThreshold sets the threshold amount for which an account record is created for any send/withdraw
// transaction
func (builder AccountCreateTransaction) SetSendRecordThreshold(recordThreshold Hbar) AccountCreateTransaction {
	builder.pb.SendRecordThreshold = uint64(recordThreshold.AsTinybar())
	return builder
}

// SetReceiveRecordThreshold sets the threshold amount for which an account record is created for any receive/deposit
// transaction
func (builder AccountCreateTransaction) SetReceiveRecordThreshold(recordThreshold Hbar) AccountCreateTransaction {
	builder.pb.ReceiveRecordThreshold = uint64(recordThreshold.AsTinybar())
	return builder
}

// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalid account, or is an account that isn't a node, then this account is automatically proxy staked to a node
// chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a node, then it will behave as if proxyAccountID was not set.
func (builder AccountCreateTransaction) SetProxyAccountID(id AccountID) AccountCreateTransaction {
	builder.pb.ProxyAccountID = id.toProto()
	return builder
}

// SetReceiverSignatureRequired sets the receiverSigRequired flag. If the receiverSigRequired flag is set to true, then
// all cryptocurrency transfers must be signed by this account's key, both for transfers in and out. If it is false,
// then only transfers out have to be signed by it. This transaction must be signed by the
// payer account. If receiverSigRequired is false, then the transaction does not have to be signed by the keys in the
// keys field. If it is true, then it must be signed by them, in addition to the keys of the payer account.
func (builder AccountCreateTransaction) SetReceiverSignatureRequired(required bool) AccountCreateTransaction {
	builder.pb.ReceiverSigRequired = required
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder AccountCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder AccountCreateTransaction) SetTransactionMemo(memo string) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder AccountCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder AccountCreateTransaction) SetTransactionID(transactionID TransactionID) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder AccountCreateTransaction) SetNodeAccountID(nodeAccountID AccountID) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
