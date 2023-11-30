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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// A TopicCreateTransaction is for creating a new Topic on HCS.
type TopicCreateTransaction struct {
	transaction
	autoRenewAccountID *AccountID
	adminKey           Key
	submitKey          Key
	memo               string
	autoRenewPeriod    *time.Duration
}

// NewTopicCreateTransaction creates a TopicCreateTransaction transaction which can be
// used to construct and execute a  Create Topic transaction.
func NewTopicCreateTransaction() *TopicCreateTransaction {
	tx := TopicCreateTransaction{
		transaction: _NewTransaction(),
	}

	tx.e = &tx
	tx.SetAutoRenewPeriod(7890000 * time.Second)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive tx takes place for this new account.
	// This should be an explicit ask.
	// tx.SetReceiveRecordThreshold(MaxHbar)
	// tx.SetSendRecordThreshold(MaxHbar)

	return &tx
}

func _TopicCreateTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *TopicCreateTransaction {
	adminKey, _ := _KeyFromProtobuf(pb.GetConsensusCreateTopic().GetAdminKey())
	submitKey, _ := _KeyFromProtobuf(pb.GetConsensusCreateTopic().GetSubmitKey())

	autoRenew := _DurationFromProtobuf(pb.GetConsensusCreateTopic().GetAutoRenewPeriod())
	resultTx := &TopicCreateTransaction{
		transaction:        tx,
		autoRenewAccountID: _AccountIDFromProtobuf(pb.GetConsensusCreateTopic().GetAutoRenewAccount()),
		adminKey:           adminKey,
		submitKey:          submitKey,
		memo:               pb.GetConsensusCreateTopic().GetMemo(),
		autoRenewPeriod:    &autoRenew,
	}
	resultTx.e = resultTx
	return resultTx
}

// SetAdminKey sets the key required to update or delete the topic. If unspecified, anyone can increase the topic's
// expirationTime.
func (tx *TopicCreateTransaction) SetAdminKey(publicKey Key) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = publicKey
	return tx
}

// GetAdminKey returns the key required to update or delete the topic
func (tx *TopicCreateTransaction) GetAdminKey() (Key, error) {
	return tx.adminKey, nil
}

// SetSubmitKey sets the key required for submitting messages to the topic. If unspecified, all submissions are allowed.
func (tx *TopicCreateTransaction) SetSubmitKey(publicKey Key) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.submitKey = publicKey
	return tx
}

// GetSubmitKey returns the key required for submitting messages to the topic
func (tx *TopicCreateTransaction) GetSubmitKey() (Key, error) {
	return tx.submitKey, nil
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (tx *TopicCreateTransaction) SetTopicMemo(memo string) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo
	return tx
}

// GetTopicMemo returns the memo for this topic
func (tx *TopicCreateTransaction) GetTopicMemo() string {
	return tx.memo
}

// SetAutoRenewPeriod sets the initial lifetime of the topic and the amount of time to extend the topic's lifetime
// automatically at expirationTime if the autoRenewAccount is configured and has sufficient funds.
//
// Required. Limited to a maximum of 90 days (server-sIDe configuration which may change).
func (tx *TopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &period
	return tx
}

// GetAutoRenewPeriod returns the auto renew period for this topic
func (tx *TopicCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return *tx.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetAutoRenewAccountID sets an optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount).
//
// If specified, there must be an adminKey and the autoRenewAccount must sign this transaction.
func (tx *TopicCreateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewAccountID = &autoRenewAccountID
	return tx
}

// GetAutoRenewAccountID returns the auto renew account ID for this topic
func (tx *TopicCreateTransaction) GetAutoRenewAccountID() AccountID {
	if tx.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *tx.autoRenewAccountID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TopicCreateTransaction) Sign(privateKey PrivateKey) *TopicCreateTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TopicCreateTransaction) SignWithOperator(client *Client) (*TopicCreateTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TopicCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicCreateTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TopicCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicCreateTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TopicCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *TopicCreateTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TopicCreateTransaction) Freeze() (*TopicCreateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TopicCreateTransaction) FreezeWith(client *Client) (*TopicCreateTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TopicCreateTransaction.
func (tx *TopicCreateTransaction) SetMaxTransactionFee(fee Hbar) *TopicCreateTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TopicCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TopicCreateTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TopicCreateTransaction.
func (tx *TopicCreateTransaction) SetTransactionMemo(memo string) *TopicCreateTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TopicCreateTransaction.
func (tx *TopicCreateTransaction) SetTransactionValidDuration(duration time.Duration) *TopicCreateTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TopicCreateTransaction.
func (tx *TopicCreateTransaction) SetTransactionID(transactionID TransactionID) *TopicCreateTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TopicCreateTransaction.
func (tx *TopicCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicCreateTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TopicCreateTransaction) SetMaxRetry(count int) *TopicCreateTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TopicCreateTransaction) SetMaxBackoff(max time.Duration) *TopicCreateTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TopicCreateTransaction) SetMinBackoff(min time.Duration) *TopicCreateTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *TopicCreateTransaction) SetLogLevel(level LogLevel) *TopicCreateTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *TopicCreateTransaction) getName() string {
	return "TopicCreateTransaction"
}

func (tx *TopicCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.autoRenewAccountID != nil {
		if err := tx.autoRenewAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TopicCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: tx.buildProtoBody(),
		},
	}
}

func (tx *TopicCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TopicCreateTransaction) buildProtoBody() *services.ConsensusCreateTopicTransactionBody {
	body := &services.ConsensusCreateTopicTransactionBody{
		Memo: tx.memo,
	}

	if tx.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*tx.autoRenewPeriod)
	}

	if tx.autoRenewAccountID != nil {
		body.AutoRenewAccount = tx.autoRenewAccountID._ToProtobuf()
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.submitKey != nil {
		body.SubmitKey = tx.submitKey._ToProtoKey()
	}

	return body
}

func (tx *TopicCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().CreateTopic,
	}
}
