package hedera

import (
	"time"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ConsensusTopicUpdateTransaction updates all fields on a Topic that are set in the transaction.
type ConsensusTopicUpdateTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusUpdateTopicTransactionBody
}

// NewConsensusTopicUpdateTransaction creates a ConsensusTopicUpdateTransaction builder which can be
// used to construct and execute a Consensus Update Topic Transaction.
func NewConsensusTopicUpdateTransaction() ConsensusTopicUpdateTransaction {
	pb := &proto.ConsensusUpdateTopicTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusUpdateTopic{pb}

	builder := ConsensusTopicUpdateTransaction{inner, pb}

	return builder
}

// SetTopicID sets the topic to be updated.
func (builder ConsensusTopicUpdateTransaction) SetTopicID(id ConsensusTopicID) ConsensusTopicUpdateTransaction {
	builder.pb.TopicID = id.toProto()
	return builder
}

// SetAdminKey sets the key required to update/delete the topic. If unset, the key will not be changed.
//
// Setting the AdminKey to an empty KeyList will clear the adminKey.
func (builder ConsensusTopicUpdateTransaction) SetAdminKey(publicKey PublicKey) ConsensusTopicUpdateTransaction {
	builder.pb.AdminKey = publicKey.toProto()
	return builder
}

// SetSubmitKey will set the key allowed to submit messages to the topic.  If unset, the key will not be changed.
//
// Setting the submitKey to an empty KeyList will clear the submitKey.
func (builder ConsensusTopicUpdateTransaction) SetSubmitKey(publicKey PublicKey) ConsensusTopicUpdateTransaction {
	builder.pb.SubmitKey = publicKey.toProto()
	return builder
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (builder ConsensusTopicUpdateTransaction) SetTopicMemo(memo string) ConsensusTopicUpdateTransaction {
	builder.pb.Memo = &wrappers.StringValue{Value: memo}
	return builder
}

// SetExpirationTime sets the effective consensus timestamp at (and after) which all consensus transactions and queries
// will fail. The expirationTime may be no longer than 90 days from the consensus timestamp of this transaction.
func (builder ConsensusTopicUpdateTransaction) SetExpirationTime(expiration time.Time) ConsensusTopicUpdateTransaction {
	builder.pb.ExpirationTime = timeToProto(expiration)
	return builder
}

// SetAutoRenewPeriod sets the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds. This is limited to a maximum of 90 days (server-side configuration
// which may change).
func (builder ConsensusTopicUpdateTransaction) SetAutoRenewPeriod(period time.Duration) ConsensusTopicUpdateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(period)
	return builder
}

// SetAutoRenewAccountID sets the optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount). If specified as the default value
// (0.0.0), the autoRenewAccount will be removed.
func (builder ConsensusTopicUpdateTransaction) SetAutoRenewAccountID(id AccountID) ConsensusTopicUpdateTransaction {
	builder.pb.AutoRenewAccount = id.toProto()
	return builder
}

// ClearTopicMemo explicitly clears any memo on the topic by sending an empty string as the memo
func (builder ConsensusTopicUpdateTransaction) ClearTopicMemo() ConsensusTopicUpdateTransaction {
	return builder.SetTopicMemo("")
}

// ClearAdminKey explicitly clears any admin key on the topic by sending an empty key list as the key
func (builder ConsensusTopicUpdateTransaction) ClearAdminKey() ConsensusTopicUpdateTransaction {
	return builder.SetAdminKey(NewKeyList())
}

// ClearSubmitKey explicitly clears any submit key on the topic by sending an empty key list as the key
func (builder ConsensusTopicUpdateTransaction) ClearSubmitKey() ConsensusTopicUpdateTransaction {
	return builder.SetSubmitKey(NewKeyList())
}

// ClearAutoRenewAccountID explicitly clears any auto renew account ID on the topic by sending an empty accountID
func (builder ConsensusTopicUpdateTransaction) ClearAutoRenewAccountID() ConsensusTopicUpdateTransaction {
	builder.pb.AutoRenewAccount = &proto.AccountID{}

	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder ConsensusTopicUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder ConsensusTopicUpdateTransaction) SetTransactionMemo(memo string) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder ConsensusTopicUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder ConsensusTopicUpdateTransaction) SetTransactionID(transactionID TransactionID) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ConsensusTopicUpdateTransaction) SetNodeAccountID(nodeAccountID AccountID) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
