package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// *TopicUpdateTransaction updates all fields on a Topic that are set in the transaction.
type TopicUpdateTransaction struct {
	Transaction
	pb *proto.ConsensusUpdateTopicTransactionBody
}

// NewTopicUpdateTransaction creates a *TopicUpdateTransaction transaction which can be
// used to construct and execute a  Update Topic Transaction.
func NewTopicUpdateTransaction() *TopicUpdateTransaction {
	pb := &proto.ConsensusUpdateTopicTransactionBody{}

	transaction := TopicUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)

	return &transaction
}

// SetTopicID sets the topic to be updated.
func (transaction *TopicUpdateTransaction) SetTopicID(topicId TopicID) *TopicUpdateTransaction {
	transaction.pb.TopicID = topicId.toProto()
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTopicID() TopicID {
	return TopicIDFromProto(transaction.pb.GetTopicID())
}

// SetAdminKey sets the key required to update/delete the topic. If unset, the key will not be changed.
//
// Setting the AdminKey to an empty KeyList will clear the adminKey.
func (transaction *TopicUpdateTransaction) SetAdminKey(publicKey PublicKey) *TopicUpdateTransaction {
	transaction.pb.AdminKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAdminKey() (Key, error) {
	return publicKeyFromProto(transaction.pb.GetAdminKey())
}

// SetSubmitKey will set the key allowed to submit messages to the topic.  If unset, the key will not be changed.
//
// Setting the submitKey to an empty KeyList will clear the submitKey.
func (transaction *TopicUpdateTransaction) SetSubmitKey(publicKey PublicKey) *TopicUpdateTransaction {
	transaction.pb.SubmitKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TopicUpdateTransaction) GetSubmitKey() (Key, error) {
	return publicKeyFromProto(transaction.pb.GetSubmitKey())
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction *TopicUpdateTransaction) SetTopicMemo(memo string) *TopicUpdateTransaction {
	transaction.pb.Memo = &proto.StringValue{Value: memo}
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTopicMemo() string {
	return transaction.pb.GetMemo().GetValue()
}

// SetExpirationTime sets the effective  timestamp at (and after) which all  transactions and queries
// will fail. The expirationTime may be no longer than 90 days from the  timestamp of this transaction.
func (transaction *TopicUpdateTransaction) SetExpirationTime(expiration time.Time) *TopicUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetExpirationTime() time.Time {
	return timeFromProto(transaction.pb.GetExpirationTime())
}

// SetAutoRenewPeriod sets the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds. This is limited to a maximum of 90 days (server-side configuration
// which may change).
func (transaction *TopicUpdateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicUpdateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(period)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProto(transaction.pb.GetAutoRenewPeriod())
}

// SetAutoRenewAccountID sets the optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount). If specified as the default value
// (0.0.0), the autoRenewAccount will be removed.
func (transaction *TopicUpdateTransaction) SetAutoRenewAccountId(accountId AccountID) *TopicUpdateTransaction {
	transaction.pb.AutoRenewAccount = accountId.toProtobuf()
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAutoRenewAccountId() AccountID {
	return accountIDFromProto(transaction.pb.GetAutoRenewAccount())
}

// ClearTopicMemo explicitly clears any memo on the topic by sending an empty string as the memo
func (transaction *TopicUpdateTransaction) ClearTopicMemo() *TopicUpdateTransaction {
	return transaction.SetTopicMemo("")
}

// ClearAdminKey explicitly clears any admin key on the topic by sending an empty key list as the key
func (transaction *TopicUpdateTransaction) ClearAdminKey() *TopicUpdateTransaction {
	return transaction.SetAdminKey(PublicKey{nil})
}

// ClearSubmitKey explicitly clears any submit key on the topic by sending an empty key list as the key
func (transaction *TopicUpdateTransaction) ClearSubmitKey() *TopicUpdateTransaction {
	return transaction.SetSubmitKey(PublicKey{nil})
}

// ClearAutoRenewAccountID explicitly clears any auto renew account ID on the topic by sending an empty accountID
func (transaction *TopicUpdateTransaction) ClearAutoRenewAccountID() *TopicUpdateTransaction {
	transaction.pb.AutoRenewAccount = &proto.AccountID{}

	return transaction
}
