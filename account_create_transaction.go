package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"math"
	"time"
)

type AccountCreateTransaction struct {
	TransactionBuilder
	pb *proto.CryptoCreateTransactionBody
}

func NewAccountCreateTransaction() AccountCreateTransaction {
	pb := &proto.CryptoCreateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_CryptoCreateAccount{pb}

	builder := AccountCreateTransaction{inner, pb}
	builder.SetAutoRenewPeriod(7890000 * time.Second)

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	builder.SetReceiveRecordThreshold(uint64(math.MaxInt64))
	builder.SetSendRecordThreshold(uint64(math.MaxInt64))

	return builder
}

func (builder AccountCreateTransaction) SetKey(publicKey Ed25519PublicKey) AccountCreateTransaction {
	builder.pb.Key = publicKey.toProto()
	return builder
}

func (builder AccountCreateTransaction) SetInitialBalance(tinyBars uint64) AccountCreateTransaction {
	builder.pb.InitialBalance = tinyBars
	return builder
}

func (builder AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) AccountCreateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return builder
}

func (builder AccountCreateTransaction) SetSendRecordThreshold(recordThreshold uint64) AccountCreateTransaction {
	builder.pb.SendRecordThreshold = recordThreshold
	return builder
}

func (builder AccountCreateTransaction) SetReceiveRecordThreshold(recordThreshold uint64) AccountCreateTransaction {
	builder.pb.ReceiveRecordThreshold = recordThreshold
	return builder
}

func (builder AccountCreateTransaction) Build(client *Client) Transaction {
	// If a shard/realm is not set, it is inferred from the Operator on the Client

	if builder.pb.ShardID == nil {
		builder.pb.ShardID = &proto.ShardID{
			ShardNum: int64(client.operator.accountID.Shard),
		}
	}

	if builder.pb.RealmID == nil {
		builder.pb.RealmID = &proto.RealmID{
			ShardNum: int64(client.operator.accountID.Shard),
			RealmNum: int64(client.operator.accountID.Realm),
		}
	}

	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder AccountCreateTransaction) SetMaxTransactionFee(maxTransactionFee uint64) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder AccountCreateTransaction) SetMemo(memo string) AccountCreateTransaction {
	return AccountCreateTransaction{builder.TransactionBuilder.SetMemo(memo), builder.pb}
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
