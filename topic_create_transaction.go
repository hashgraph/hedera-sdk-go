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

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// A TopicCreateTransaction is for creating a new Topic on HCS.
type TopicCreateTransaction struct {
	Transaction
	autoRenewAccountID *AccountID
	adminKey           Key
	submitKey          Key
	memo               string
	autoRenewPeriod    *time.Duration
}

// NewTopicCreateTransaction creates a TopicCreateTransaction transaction which can be
// used to construct and execute a  Create Topic Transaction.
func NewTopicCreateTransaction() *TopicCreateTransaction {
	transaction := TopicCreateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	// transaction.SetReceiveRecordThreshold(MaxHbar)
	// transaction.SetSendRecordThreshold(MaxHbar)

	return &transaction
}

func _TopicCreateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TopicCreateTransaction {
	adminKey, _ := _KeyFromProtobuf(pb.GetConsensusCreateTopic().GetAdminKey())
	submitKey, _ := _KeyFromProtobuf(pb.GetConsensusCreateTopic().GetSubmitKey())

	autoRenew := _DurationFromProtobuf(pb.GetConsensusCreateTopic().GetAutoRenewPeriod())
	return &TopicCreateTransaction{
		Transaction:        transaction,
		autoRenewAccountID: _AccountIDFromProtobuf(pb.GetConsensusCreateTopic().GetAutoRenewAccount()),
		adminKey:           adminKey,
		submitKey:          submitKey,
		memo:               pb.GetConsensusCreateTopic().GetMemo(),
		autoRenewPeriod:    &autoRenew,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *TopicCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *TopicCreateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetAdminKey sets the key required to update or delete the topic. If unspecified, anyone can increase the topic's
// expirationTime.
func (transaction *TopicCreateTransaction) SetAdminKey(publicKey Key) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

// GetAdminKey returns the key required to update or delete the topic
func (transaction *TopicCreateTransaction) GetAdminKey() (Key, error) {
	return transaction.adminKey, nil
}

// SetSubmitKey sets the key required for submitting messages to the topic. If unspecified, all submissions are allowed.
func (transaction *TopicCreateTransaction) SetSubmitKey(publicKey Key) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.submitKey = publicKey
	return transaction
}

// GetSubmitKey returns the key required for submitting messages to the topic
func (transaction *TopicCreateTransaction) GetSubmitKey() (Key, error) {
	return transaction.submitKey, nil
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction *TopicCreateTransaction) SetTopicMemo(memo string) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	return transaction
}

// GetTopicMemo returns the memo for this topic
func (transaction *TopicCreateTransaction) GetTopicMemo() string {
	return transaction.memo
}

// SetAutoRenewPeriod sets the initial lifetime of the topic and the amount of time to extend the topic's lifetime
// automatically at expirationTime if the autoRenewAccount is configured and has sufficient funds.
//
// Required. Limited to a maximum of 90 days (server-sIDe configuration which may change).
func (transaction *TopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &period
	return transaction
}

// GetAutoRenewPeriod returns the auto renew period for this topic
func (transaction *TopicCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetAutoRenewAccountID sets an optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount).
//
// If specified, there must be an adminKey and the autoRenewAccount must sign this transaction.
func (transaction *TopicCreateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewAccountID = &autoRenewAccountID
	return transaction
}

// GetAutoRenewAccountID returns the auto renew account ID for this topic
func (transaction *TopicCreateTransaction) GetAutoRenewAccountID() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

func (transaction *TopicCreateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.autoRenewAccountID != nil {
		if err := transaction.autoRenewAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TopicCreateTransaction) _Build() *services.TransactionBody {
	body := &services.ConsensusCreateTopicTransactionBody{
		Memo: transaction.memo,
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccount = transaction.autoRenewAccountID._ToProtobuf()
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
		Data: &services.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: body,
		},
	}
}

func (transaction *TopicCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TopicCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ConsensusCreateTopicTransactionBody{
		Memo: transaction.memo,
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccount = transaction.autoRenewAccountID._ToProtobuf()
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
		Data: &services.SchedulableTransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: body,
		},
	}, nil
}

func _TopicCreateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().CreateTopic,
	}
}

func (transaction *TopicCreateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TopicCreateTransaction) Sign(
	privateKey PrivateKey,
) *TopicCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *TopicCreateTransaction) SignWithOperator(
	client *Client,
) (*TopicCreateTransaction, error) {
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
func (transaction *TopicCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicCreateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TopicCreateTransaction) Execute(
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
		_TopicCreateTransactionGetMethod,
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

func (transaction *TopicCreateTransaction) Freeze() (*TopicCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicCreateTransaction) FreezeWith(client *Client) (*TopicCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TopicCreateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TopicCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TopicCreateTransaction) SetMaxTransactionFee(fee Hbar) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TopicCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TopicCreateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this	TopicCreateTransaction.
func (transaction *TopicCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetTransactionMemo(memo string) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *TopicCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetTransactionValidDuration(duration time.Duration) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetTransactionID(transactionID TransactionID) *TopicCreateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *TopicCreateTransaction) SetMaxRetry(count int) *TopicCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *TopicCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicCreateTransaction {
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
func (transaction *TopicCreateTransaction) SetMaxBackoff(max time.Duration) *TopicCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *TopicCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *TopicCreateTransaction) SetMinBackoff(min time.Duration) *TopicCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TopicCreateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TopicCreateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TopicCreateTransaction:%d", timestamp.UnixNano())
}

func (transaction *TopicCreateTransaction) SetLogLevel(level LogLevel) *TopicCreateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
