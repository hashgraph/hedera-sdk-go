package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ConsensusTopicCreateTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusCreateTopicTransactionBody
}

func NewConsensusTopicCreateTransaction() ConsensusTopicCreateTransaction {
	pb := &proto.ConsensusCreateTopicTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusCreateTopic{pb}

	builder := ConsensusTopicCreateTransaction{inner, pb}

	return builder
}

func (builder ConsensusTopicCreateTransaction) SetAdminKey(publicKey Ed25519PublicKey) ConsensusTopicCreateTransaction {
	builder.pb.AdminKey = publicKey.toProto()
	return builder
}

func (builder ConsensusTopicCreateTransaction) SetSubmitKey(publicKey Ed25519PublicKey) ConsensusTopicCreateTransaction {
	builder.pb.SubmitKey = publicKey.toProto()
	return builder
}

func (builder ConsensusTopicCreateTransaction) SetValidStart(start time.Time) ConsensusTopicCreateTransaction {
	builder.pb.ValidStartTime = timeToProto(start)
	return builder
}

func (builder ConsensusTopicCreateTransaction) SetTopicMemo(memo string) ConsensusTopicCreateTransaction {
	builder.pb.Memo = memo
	return builder
}

func (builder ConsensusTopicCreateTransaction) SetExpirationTime(expiration time.Time) ConsensusTopicCreateTransaction {
	builder.pb.ExpirationTime = timeToProto(expiration)
	return builder
}

func (builder ConsensusTopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) ConsensusTopicCreateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(period)
	return builder
}

func (builder ConsensusTopicCreateTransaction) SetAutoRenewAccount(id AccountID) ConsensusTopicCreateTransaction {
	builder.pb.AutoRenewAccount = id.toProto()
	return builder
}

func (builder ConsensusTopicCreateTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ConsensusTopicCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder ConsensusTopicCreateTransaction) SetMemo(memo string) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder ConsensusTopicCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder ConsensusTopicCreateTransaction) SetTransactionID(transactionID TransactionID) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ConsensusTopicCreateTransaction) SetNodeAccountID(nodeAccountID AccountID) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
