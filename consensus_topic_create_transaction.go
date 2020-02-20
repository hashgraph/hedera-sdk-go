package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ConsensusTopicCreateTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusCreateTopicTransactionBody
}

// NewConsensusTopicCreateTransaction creates a ConsensusTopicCreateTransaction builder which can be
// used to construct and execute a Consensus Create Topic Transaction.
func NewConsensusTopicCreateTransaction() ConsensusTopicCreateTransaction {
	pb := &proto.ConsensusCreateTopicTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusCreateTopic{pb}

	builder := ConsensusTopicCreateTransaction{inner, pb}

	return builder.SetAutoRenewPeriod(7890000 * time.Second)
}

// SetAdminKey sets the key required to update or delete the topic. If unspecified, anyone can increase the topic's
// expirationTime.
func (builder ConsensusTopicCreateTransaction) SetAdminKey(publicKey Ed25519PublicKey) ConsensusTopicCreateTransaction {
	builder.pb.AdminKey = publicKey.toProto()
	return builder
}

// SetSubmitKey sets the key required for submitting messages to the topic. If unspecified, all submissions are allowed.
func (builder ConsensusTopicCreateTransaction) SetSubmitKey(publicKey Ed25519PublicKey) ConsensusTopicCreateTransaction {
	builder.pb.SubmitKey = publicKey.toProto()
	return builder
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (builder ConsensusTopicCreateTransaction) SetTopicMemo(memo string) ConsensusTopicCreateTransaction {
	builder.pb.Memo = memo
	return builder
}

// SetAutoRenewPeriod sets the initial lifetime of the topic and the amount of time to extend the topic's lifetime
// automatically at expirationTime if the autoRenewAccount is configured and has sufficient funds.
//
// Required. Limited to a maximum of 90 days (server-side configuration which may change).
func (builder ConsensusTopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) ConsensusTopicCreateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(period)
	return builder
}

// SetAutoRenewAccountID sets an optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount).
//
//If specified, there must be an adminKey and the autoRenewAccount must sign this transaction.
func (builder ConsensusTopicCreateTransaction) SetAutoRenewAccountID(id AccountID) ConsensusTopicCreateTransaction {
	builder.pb.AutoRenewAccount = id.toProto()
	return builder
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
