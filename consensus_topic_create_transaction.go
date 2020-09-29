package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// A ConsensusTopicCreateTransaction is for creating a new Topic on HCS.
type ConsensusTopicCreateTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusCreateTopicTransactionBody
}

// NewConsensusTopicCreateTransaction creates a ConsensusTopicCreateTransaction transaction which can be
// used to construct and execute a Consensus Create Topic Transaction.
func NewConsensusTopicCreateTransaction() ConsensusTopicCreateTransaction {
	pb := &proto.ConsensusCreateTopicTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusCreateTopic{ConsensusCreateTopic: pb}

	transaction := ConsensusTopicCreateTransaction{inner, pb}

	return transaction.SetAutoRenewPeriod(7890000 * time.Second)
}

// SetAdminKey sets the key required to update or delete the topic. If unspecified, anyone can increase the topic's
// expirationTime.
func (transaction ConsensusTopicCreateTransaction) SetAdminKey(publicKey PublicKey) ConsensusTopicCreateTransaction {
	transaction.pb.AdminKey = publicKey.toProto()
	return transaction
}

// SetSubmitKey sets the key required for submitting messages to the topic. If unspecified, all submissions are allowed.
func (transaction ConsensusTopicCreateTransaction) SetSubmitKey(publicKey PublicKey) ConsensusTopicCreateTransaction {
	transaction.pb.SubmitKey = publicKey.toProto()
	return transaction
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction ConsensusTopicCreateTransaction) SetTopicMemo(memo string) ConsensusTopicCreateTransaction {
	transaction.pb.Memo = memo
	return transaction
}

// SetAutoRenewPeriod sets the initial lifetime of the topic and the amount of time to extend the topic's lifetime
// automatically at expirationTime if the autoRenewAccount is configured and has sufficient funds.
//
// Required. Limited to a maximum of 90 days (server-side configuration which may change).
func (transaction ConsensusTopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) ConsensusTopicCreateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(period)
	return transaction
}

// SetAutoRenewAccountID sets an optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount).
//
//If specified, there must be an adminKey and the autoRenewAccount must sign this transaction.
func (transaction ConsensusTopicCreateTransaction) SetAutoRenewAccountID(id AccountID) ConsensusTopicCreateTransaction {
	transaction.pb.AutoRenewAccount = id.toProto()
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction ConsensusTopicCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction ConsensusTopicCreateTransaction) SetTransactionMemo(memo string) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction ConsensusTopicCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction ConsensusTopicCreateTransaction) SetTransactionID(transactionID TransactionID) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction ConsensusTopicCreateTransaction) SetNodeID(nodeAccountID AccountID) ConsensusTopicCreateTransaction {
	return ConsensusTopicCreateTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
