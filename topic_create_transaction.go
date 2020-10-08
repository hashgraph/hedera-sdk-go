package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// A TopicCreateTransaction is for creating a new Topic on HCS.
type TopicCreateTransaction struct {
	Transaction
	pb *proto.ConsensusCreateTopicTransactionBody
}

// NewTopicCreateTransaction creates a TopicCreateTransaction transaction which can be
// used to construct and execute a  Create Topic Transaction.
func NewTopicCreateTransaction() *TopicCreateTransaction {
	pb := &proto.ConsensusCreateTopicTransactionBody{}

	transaction := TopicCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	//transaction.SetReceiveRecordThreshold(MaxHbar)
	//transaction.SetSendRecordThreshold(MaxHbar)

	return &transaction

}

// SetAdminKey sets the key required to update or delete the topic. If unspecified, anyone can increase the topic's
// expirationTime.
func (transaction *TopicCreateTransaction) SetAdminKey(publicKey PublicKey) *TopicCreateTransaction {
	transaction.pb.AdminKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TopicCreateTransaction) GetAdminKey() (Key, error) {
	return publicKeyFromProto(transaction.pb.GetAdminKey())
}

// SetSubmitKey sets the key required for submitting messages to the topic. If unspecified, all submissions are allowed.
func (transaction *TopicCreateTransaction) SetSubmitKey(publicKey PublicKey) *TopicCreateTransaction {
	transaction.pb.SubmitKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TopicCreateTransaction) GetSubmitKey() (Key, error) {
	return publicKeyFromProto(transaction.pb.GetSubmitKey())
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction *TopicCreateTransaction) SetTopicMemo(memo string) *TopicCreateTransaction {
	transaction.pb.Memo = memo
	return transaction
}

func (transaction *TopicCreateTransaction) GetTopicMemo() string {
	return transaction.pb.GetMemo()
}

// SetAutoRenewPeriod sets the initial lifetime of the topic and the amount of time to extend the topic's lifetime
// automatically at expirationTime if the autoRenewAccount is configured and has sufficient funds.
//
// Required. Limited to a maximum of 90 days (server-side configuration which may change).
func (transaction *TopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicCreateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(period)
	return transaction
}

func (transaction *TopicCreateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProto(transaction.pb.GetAutoRenewPeriod())
}

// SetAutoRenewAccountID sets an optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount).
//
//If specified, there must be an adminKey and the autoRenewAccount must sign this transaction.
func (transaction *TopicCreateTransaction) SetAutoRenewAccountId(accountId AccountID) *TopicCreateTransaction {
	transaction.pb.AutoRenewAccount = accountId.toProtobuf()
	return transaction
}

func (transaction *TopicCreateTransaction) GetAutoRenewAccoutnId() AccountID {
	return accountIDFromProto(transaction.pb.GetAutoRenewAccount())
}

