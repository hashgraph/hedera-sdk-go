package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TopicUpdateTransaction
// Updates all fields on a Topic that are set in the transaction.
type TopicUpdateTransaction struct {
	Transaction
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
	transaction := TopicUpdateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _TopicUpdateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TopicUpdateTransaction {
	adminKey, _ := _KeyFromProtobuf(pb.GetConsensusUpdateTopic().GetAdminKey())
	submitKey, _ := _KeyFromProtobuf(pb.GetConsensusUpdateTopic().GetSubmitKey())

	expirationTime := _TimeFromProtobuf(pb.GetConsensusUpdateTopic().GetExpirationTime())
	autoRenew := _DurationFromProtobuf(pb.GetConsensusUpdateTopic().GetAutoRenewPeriod())
	return &TopicUpdateTransaction{
		Transaction:        transaction,
		topicID:            _TopicIDFromProtobuf(pb.GetConsensusUpdateTopic().GetTopicID()),
		autoRenewAccountID: _AccountIDFromProtobuf(pb.GetConsensusUpdateTopic().GetAutoRenewAccount()),
		adminKey:           adminKey,
		submitKey:          submitKey,
		memo:               pb.GetConsensusUpdateTopic().GetMemo().Value,
		autoRenewPeriod:    &autoRenew,
		expirationTime:     &expirationTime,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *TopicUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *TopicUpdateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetTopicID sets the topic to be updated.
func (transaction *TopicUpdateTransaction) SetTopicID(topicID TopicID) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.topicID = &topicID
	return transaction
}

// GetTopicID returns the topic to be updated.
func (transaction *TopicUpdateTransaction) GetTopicID() TopicID {
	if transaction.topicID == nil {
		return TopicID{}
	}

	return *transaction.topicID
}

// SetAdminKey sets the key required to update/delete the topic. If unset, the key will not be changed.
//
// Setting the AdminKey to an empty KeyList will clear the adminKey.
func (transaction *TopicUpdateTransaction) SetAdminKey(publicKey Key) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

// GetAdminKey returns the key required to update/delete the topic.
func (transaction *TopicUpdateTransaction) GetAdminKey() (Key, error) {
	return transaction.adminKey, nil
}

// SetSubmitKey will set the key allowed to submit messages to the topic.  If unset, the key will not be changed.
//
// Setting the submitKey to an empty KeyList will clear the submitKey.
func (transaction *TopicUpdateTransaction) SetSubmitKey(publicKey Key) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.submitKey = publicKey
	return transaction
}

// GetSubmitKey returns the key allowed to submit messages to the topic.
func (transaction *TopicUpdateTransaction) GetSubmitKey() (Key, error) {
	return transaction.submitKey, nil
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction *TopicUpdateTransaction) SetTopicMemo(memo string) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	return transaction
}

// GetTopicMemo returns the short publicly visible memo about the topic.
func (transaction *TopicUpdateTransaction) GetTopicMemo() string {
	return transaction.memo
}

// SetExpirationTime sets the effective  timestamp at (and after) which all  transactions and queries
// will fail. The expirationTime may be no longer than 90 days from the  timestamp of this transaction.
func (transaction *TopicUpdateTransaction) SetExpirationTime(expiration time.Time) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

// GetExpirationTime returns the effective  timestamp at (and after) which all transactions and queries will fail.
func (transaction *TopicUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

// SetAutoRenewPeriod sets the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds. This is limited to a maximum of 90 days (server-sIDe configuration
// which may change).
func (transaction *TopicUpdateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &period
	return transaction
}

// GetAutoRenewPeriod returns the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds.
func (transaction *TopicUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetAutoRenewAccountID sets the optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount). If specified as the default value
// (0.0.0), the autoRenewAccount will be removed.
func (transaction *TopicUpdateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewAccountID = &autoRenewAccountID
	return transaction
}

// GetAutoRenewAccountID returns the optional account to be used at the topic's expirationTime to extend the life of the
// topic.
func (transaction *TopicUpdateTransaction) GetAutoRenewAccountID() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

// ClearTopicMemo explicitly clears any memo on the topic by sending an empty string as the memo
func (transaction *TopicUpdateTransaction) ClearTopicMemo() *TopicUpdateTransaction {
	return transaction.SetTopicMemo("")
}

// ClearAdminKey explicitly clears any admin key on the topic by sending an empty key list as the key
func (transaction *TopicUpdateTransaction) ClearAdminKey() *TopicUpdateTransaction {
	return transaction.SetAdminKey(PublicKey{nil, nil})
}

// ClearSubmitKey explicitly clears any submit key on the topic by sending an empty key list as the key
func (transaction *TopicUpdateTransaction) ClearSubmitKey() *TopicUpdateTransaction {
	return transaction.SetSubmitKey(PublicKey{nil, nil})
}

// ClearAutoRenewAccountID explicitly clears any auto renew account ID on the topic by sending an empty accountID
func (transaction *TopicUpdateTransaction) ClearAutoRenewAccountID() *TopicUpdateTransaction {
	transaction.autoRenewAccountID = &AccountID{}
	return transaction
}

func (transaction *TopicUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.topicID != nil {
		if err := transaction.topicID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.autoRenewAccountID != nil {
		if err := transaction.autoRenewAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TopicUpdateTransaction) _Build() *services.TransactionBody {
	body := &services.ConsensusUpdateTopicTransactionBody{
		Memo: &wrapperspb.StringValue{Value: transaction.memo},
	}

	if transaction.topicID != nil {
		body.TopicID = transaction.topicID._ToProtobuf()
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccount = transaction.autoRenewAccountID._ToProtobuf()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.submitKey != nil {
		body.SubmitKey = transaction.submitKey._ToProtoKey()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: body,
		},
	}
}

func (transaction *TopicUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TopicUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ConsensusUpdateTopicTransactionBody{
		Memo: &wrapperspb.StringValue{Value: transaction.memo},
	}

	if transaction.topicID != nil {
		body.TopicID = transaction.topicID._ToProtobuf()
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccount = transaction.autoRenewAccountID._ToProtobuf()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.submitKey != nil {
		body.SubmitKey = transaction.submitKey._ToProtoKey()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: body,
		},
	}, nil
}

func _TopicUpdateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().UpdateTopic,
	}
}

func (transaction *TopicUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TopicUpdateTransaction) Sign(
	privateKey PrivateKey,
) *TopicUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *TopicUpdateTransaction) SignWithOperator(
	client *Client,
) (*TopicUpdateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *TopicUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicUpdateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TopicUpdateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TopicUpdateTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
		transaction.maxBackoff,
		transaction.minBackoff,
		transaction.maxRetry,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID:  transaction.GetTransactionID(),
			NodeID:         resp.(TransactionResponse).NodeID,
			ValidateStatus: true,
		}, err
	}

	return TransactionResponse{
		TransactionID:  transaction.GetTransactionID(),
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
	}, nil
}

func (transaction *TopicUpdateTransaction) Freeze() (*TopicUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicUpdateTransaction) FreezeWith(client *Client) (*TopicUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TopicUpdateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TopicUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TopicUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TopicUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TopicUpdateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this	TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionMemo(memo string) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *TopicUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionID(transactionID TransactionID) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *TopicUpdateTransaction) SetMaxRetry(count int) *TopicUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *TopicUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicUpdateTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := transaction.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
	}

	return transaction
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *TopicUpdateTransaction) SetMaxBackoff(max time.Duration) *TopicUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *TopicUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *TopicUpdateTransaction) SetMinBackoff(min time.Duration) *TopicUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *TopicUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TopicUpdateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TopicUpdateTransaction:%d", timestamp.UnixNano())
}

func (transaction *TopicUpdateTransaction) SetLogLevel(level LogLevel) *TopicUpdateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
