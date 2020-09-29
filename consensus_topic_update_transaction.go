package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ConsensusTopicUpdateTransaction updates all fields on a Topic that are set in the transaction.
type ConsensusTopicUpdateTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusUpdateTopicTransactionBody
}

// NewConsensusTopicUpdateTransaction creates a ConsensusTopicUpdateTransaction transaction which can be
// used to construct and execute a Consensus Update Topic Transaction.
func NewConsensusTopicUpdateTransaction() ConsensusTopicUpdateTransaction {
	pb := &proto.ConsensusUpdateTopicTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusUpdateTopic{ConsensusUpdateTopic: pb}

	transaction := ConsensusTopicUpdateTransaction{inner, pb}

	return transaction
}

// SetTopicID sets the topic to be updated.
func (transaction ConsensusTopicUpdateTransaction) SetTopicID(id ConsensusTopicID) ConsensusTopicUpdateTransaction {
	transaction.pb.TopicID = id.toProto()
	return transaction
}

// SetAdminKey sets the key required to update/delete the topic. If unset, the key will not be changed.
//
// Setting the AdminKey to an empty KeyList will clear the adminKey.
func (transaction ConsensusTopicUpdateTransaction) SetAdminKey(publicKey PublicKey) ConsensusTopicUpdateTransaction {
	transaction.pb.AdminKey = publicKey.toProto()
	return transaction
}

// SetSubmitKey will set the key allowed to submit messages to the topic.  If unset, the key will not be changed.
//
// Setting the submitKey to an empty KeyList will clear the submitKey.
func (transaction ConsensusTopicUpdateTransaction) SetSubmitKey(publicKey PublicKey) ConsensusTopicUpdateTransaction {
	transaction.pb.SubmitKey = publicKey.toProto()
	return transaction
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction ConsensusTopicUpdateTransaction) SetTopicMemo(memo string) ConsensusTopicUpdateTransaction {
	transaction.pb.Memo = &proto.StringValue{Value: memo}
	return transaction
}

// SetExpirationTime sets the effective consensus timestamp at (and after) which all consensus transactions and queries
// will fail. The expirationTime may be no longer than 90 days from the consensus timestamp of this transaction.
func (transaction ConsensusTopicUpdateTransaction) SetExpirationTime(expiration time.Time) ConsensusTopicUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

// SetAutoRenewPeriod sets the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds. This is limited to a maximum of 90 days (server-side configuration
// which may change).
func (transaction ConsensusTopicUpdateTransaction) SetAutoRenewPeriod(period time.Duration) ConsensusTopicUpdateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(period)
	return transaction
}

// SetAutoRenewAccountID sets the optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount). If specified as the default value
// (0.0.0), the autoRenewAccount will be removed.
func (transaction ConsensusTopicUpdateTransaction) SetAutoRenewAccountID(id AccountID) ConsensusTopicUpdateTransaction {
	transaction.pb.AutoRenewAccount = id.toProto()
	return transaction
}

// ClearTopicMemo explicitly clears any memo on the topic by sending an empty string as the memo
func (transaction ConsensusTopicUpdateTransaction) ClearTopicMemo() ConsensusTopicUpdateTransaction {
	return transaction.SetTopicMemo("")
}

// ClearAdminKey explicitly clears any admin key on the topic by sending an empty key list as the key
func (transaction ConsensusTopicUpdateTransaction) ClearAdminKey() ConsensusTopicUpdateTransaction {
	return transaction.SetAdminKey(NewKeyList())
}

// ClearSubmitKey explicitly clears any submit key on the topic by sending an empty key list as the key
func (transaction ConsensusTopicUpdateTransaction) ClearSubmitKey() ConsensusTopicUpdateTransaction {
	return transaction.SetSubmitKey(NewKeyList())
}

// ClearAutoRenewAccountID explicitly clears any auto renew account ID on the topic by sending an empty accountID
func (transaction ConsensusTopicUpdateTransaction) ClearAutoRenewAccountID() ConsensusTopicUpdateTransaction {
	transaction.pb.AutoRenewAccount = &proto.AccountID{}

	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction ConsensusTopicUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction ConsensusTopicUpdateTransaction) SetTransactionMemo(memo string) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction ConsensusTopicUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction ConsensusTopicUpdateTransaction) SetTransactionID(transactionID TransactionID) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction ConsensusTopicUpdateTransaction) SetNodeID(nodeAccountID AccountID) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
