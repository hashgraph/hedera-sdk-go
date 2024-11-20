package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TopicUpdateTransaction
// Updates all fields on a Topic that are set in the transaction.
type TopicUpdateTransaction struct {
	*Transaction[*TopicUpdateTransaction]
	topicID            *TopicID
	autoRenewAccountID *AccountID
	adminKey           Key
	submitKey          Key
	memo               string
	autoRenewPeriod    *time.Duration
	expirationTime     *time.Time
}

// NewTopicUpdateTransaction creates a TopicUpdateTransaction transaction which
// updates all fields on a Topic that are set in the transaction.
func NewTopicUpdateTransaction() *TopicUpdateTransaction {
	tx := &TopicUpdateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx.SetAutoRenewPeriod(7890000 * time.Second)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _TopicUpdateTransactionFromProtobuf(tx Transaction[*TopicUpdateTransaction], pb *services.TransactionBody) TopicUpdateTransaction {
	adminKey, _ := _KeyFromProtobuf(pb.GetConsensusUpdateTopic().GetAdminKey())
	submitKey, _ := _KeyFromProtobuf(pb.GetConsensusUpdateTopic().GetSubmitKey())

	expirationTime := _TimeFromProtobuf(pb.GetConsensusUpdateTopic().GetExpirationTime())
	autoRenew := _DurationFromProtobuf(pb.GetConsensusUpdateTopic().GetAutoRenewPeriod())
	topicCreateTransaction := TopicUpdateTransaction{
		topicID:            _TopicIDFromProtobuf(pb.GetConsensusUpdateTopic().GetTopicID()),
		autoRenewAccountID: _AccountIDFromProtobuf(pb.GetConsensusUpdateTopic().GetAutoRenewAccount()),
		adminKey:           adminKey,
		submitKey:          submitKey,
		memo:               pb.GetConsensusUpdateTopic().GetMemo().Value,
		autoRenewPeriod:    &autoRenew,
		expirationTime:     &expirationTime,
	}

	tx.childTransaction = &topicCreateTransaction
	topicCreateTransaction.Transaction = &tx
	return topicCreateTransaction
}

// SetTopicID sets the topic to be updated.
func (tx *TopicUpdateTransaction) SetTopicID(topicID TopicID) *TopicUpdateTransaction {
	tx._RequireNotFrozen()
	tx.topicID = &topicID
	return tx
}

// GetTopicID returns the topic to be updated.
func (tx *TopicUpdateTransaction) GetTopicID() TopicID {
	if tx.topicID == nil {
		return TopicID{}
	}

	return *tx.topicID
}

// SetAdminKey sets the key required to update/delete the topic. If unset, the key will not be changed.
//
// Setting the AdminKey to an empty KeyList will clear the adminKey.
func (tx *TopicUpdateTransaction) SetAdminKey(publicKey Key) *TopicUpdateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = publicKey
	return tx
}

// GetAdminKey returns the key required to update/delete the topic.
func (tx *TopicUpdateTransaction) GetAdminKey() (Key, error) {
	return tx.adminKey, nil
}

// SetSubmitKey will set the key allowed to submit messages to the topic.  If unset, the key will not be changed.
//
// Setting the submitKey to an empty KeyList will clear the submitKey.
func (tx *TopicUpdateTransaction) SetSubmitKey(publicKey Key) *TopicUpdateTransaction {
	tx._RequireNotFrozen()
	tx.submitKey = publicKey
	return tx
}

// GetSubmitKey returns the key allowed to submit messages to the topic.
func (tx *TopicUpdateTransaction) GetSubmitKey() (Key, error) {
	return tx.submitKey, nil
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (tx *TopicUpdateTransaction) SetTopicMemo(memo string) *TopicUpdateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo
	return tx
}

// GetTopicMemo returns the short publicly visible memo about the topic.
func (tx *TopicUpdateTransaction) GetTopicMemo() string {
	return tx.memo
}

// SetExpirationTime sets the effective  timestamp at (and after) which all  transactions and queries
// will fail. The expirationTime may be no longer than 90 days from the  timestamp of this transaction.
func (tx *TopicUpdateTransaction) SetExpirationTime(expiration time.Time) *TopicUpdateTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expiration
	return tx
}

// GetExpirationTime returns the effective  timestamp at (and after) which all transactions and queries will fail.
func (tx *TopicUpdateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}

	return time.Time{}
}

// SetAutoRenewPeriod sets the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds. This is limited to a maximum of 90 days (server-sIDe configuration
// which may change).
func (tx *TopicUpdateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicUpdateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &period
	return tx
}

// GetAutoRenewPeriod returns the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds.
func (tx *TopicUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return *tx.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetAutoRenewAccountID sets the optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount). If specified as the default value
// (0.0.0), the autoRenewAccount will be removed.
func (tx *TopicUpdateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) *TopicUpdateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewAccountID = &autoRenewAccountID
	return tx
}

// GetAutoRenewAccountID returns the optional account to be used at the topic's expirationTime to extend the life of the
// topic.
func (tx *TopicUpdateTransaction) GetAutoRenewAccountID() AccountID {
	if tx.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *tx.autoRenewAccountID
}

// ClearTopicMemo explicitly clears any memo on the topic by sending an empty string as the memo
func (tx *TopicUpdateTransaction) ClearTopicMemo() *TopicUpdateTransaction {
	return tx.SetTopicMemo("")
}

// ClearAdminKey explicitly clears any admin key on the topic by sending an empty key list as the key
func (tx *TopicUpdateTransaction) ClearAdminKey() *TopicUpdateTransaction {
	return tx.SetAdminKey(PublicKey{nil, nil})
}

// ClearSubmitKey explicitly clears any submit key on the topic by sending an empty key list as the key
func (tx *TopicUpdateTransaction) ClearSubmitKey() *TopicUpdateTransaction {
	return tx.SetSubmitKey(PublicKey{nil, nil})
}

// ClearAutoRenewAccountID explicitly clears any auto renew account ID on the topic by sending an empty accountID
func (tx *TopicUpdateTransaction) ClearAutoRenewAccountID() *TopicUpdateTransaction {
	tx.autoRenewAccountID = &AccountID{}
	return tx
}

// ----------- Overridden functions ----------------

func (tx TopicUpdateTransaction) getName() string {
	return "TopicUpdateTransaction"
}

func (tx TopicUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.topicID != nil {
		if err := tx.topicID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.autoRenewAccountID != nil {
		if err := tx.autoRenewAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TopicUpdateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: tx.buildProtoBody(),
		},
	}
}

func (tx TopicUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TopicUpdateTransaction) buildProtoBody() *services.ConsensusUpdateTopicTransactionBody {
	body := &services.ConsensusUpdateTopicTransactionBody{
		Memo: &wrapperspb.StringValue{Value: tx.memo},
	}

	if tx.topicID != nil {
		body.TopicID = tx.topicID._ToProtobuf()
	}

	if tx.autoRenewAccountID != nil {
		body.AutoRenewAccount = tx.autoRenewAccountID._ToProtobuf()
	}

	if tx.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*tx.autoRenewPeriod)
	}

	if tx.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*tx.expirationTime)
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.submitKey != nil {
		body.SubmitKey = tx.submitKey._ToProtoKey()
	}

	return body
}

func (tx TopicUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().UpdateTopic,
	}
}

func (tx TopicUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TopicUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
