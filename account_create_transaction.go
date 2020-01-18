package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type AccountCreateTransaction struct {
	TransactionBuilder
	pb *proto.CryptoCreateTransactionBody
}

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

func (builder AccountCreateTransaction) SetKey(publicKey PublicKey) AccountCreateTransaction {
	builder.pb.Key = publicKey.toProto()
	return builder
}

func (builder AccountCreateTransaction) SetInitialBalance(initialBalance Hbar) AccountCreateTransaction {
	builder.pb.InitialBalance = uint64(initialBalance.AsTinybar())
	return builder
}

func (builder AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) AccountCreateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return builder
}

func (builder AccountCreateTransaction) SetSendRecordThreshold(recordThreshold Hbar) AccountCreateTransaction {
	builder.pb.SendRecordThreshold = uint64(recordThreshold.AsTinybar())
	return builder
}

func (builder AccountCreateTransaction) SetReceiveRecordThreshold(recordThreshold Hbar) AccountCreateTransaction {
	builder.pb.ReceiveRecordThreshold = uint64(recordThreshold.AsTinybar())
	return builder
}

func (builder AccountCreateTransaction) SetProxyAccountID(id AccountID) AccountCreateTransaction {
	builder.pb.ProxyAccountID = id.toProto()
	return builder
}

func (builder AccountCreateTransaction) SetReceiverSignatureRequired(required bool) AccountCreateTransaction {
	builder.pb.ReceiverSigRequired = required
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder AccountCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder AccountCreateTransaction) SetTransactionMemo(memo string) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder AccountCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder AccountCreateTransaction) SetTransactionID(transactionID TransactionID) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder AccountCreateTransaction) SetNodeAccountID(nodeAccountID AccountID) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
