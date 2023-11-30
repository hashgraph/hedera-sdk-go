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

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// ScheduleSignTransaction Adds zero or more signing keys to a schedule.
// If Long Term Scheduled Transactions are enabled and wait for expiry was set to true on the
// ScheduleCreate then the transaction will always wait till it's `expiration_time` to execute.
// Otherwise, if the resulting set of signing keys satisfy the
// scheduled transaction's signing requirements, it will be executed immediately after the
// triggering ScheduleSign.
// Upon SUCCESS, the receipt includes the scheduledTransactionID to use to query
// for the record of the scheduled transaction's execution (if it occurs).
type ScheduleSignTransaction struct {
	transaction
	scheduleID *ScheduleID
}

// NewScheduleSignTransaction creates ScheduleSignTransaction which adds zero or more signing keys to a schedule.
// If Long Term Scheduled Transactions are enabled and wait for expiry was set to true on the
// ScheduleCreate then the transaction will always wait till it's `expiration_time` to execute.
// Otherwise, if the resulting set of signing keys satisfy the
// scheduled transaction's signing requirements, it will be executed immediately after the
// triggering ScheduleSign.
// Upon SUCCESS, the receipt includes the scheduledTransactionID to use to query
// for the record of the scheduled transaction's execution (if it occurs).
func NewScheduleSignTransaction() *ScheduleSignTransaction {
	tx := ScheduleSignTransaction{
		transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _ScheduleSignTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *ScheduleSignTransaction {
	resultTx := &ScheduleSignTransaction{
		transaction: tx,
		scheduleID:  _ScheduleIDFromProtobuf(pb.GetScheduleSign().GetScheduleID()),
	}
	resultTx.e = resultTx
	return resultTx
}

// SetScheduleID Sets the id of the schedule to add signing keys to
func (tx *ScheduleSignTransaction) SetScheduleID(scheduleID ScheduleID) *ScheduleSignTransaction {
	tx._RequireNotFrozen()
	tx.scheduleID = &scheduleID
	return tx
}

// GetScheduleID returns the id of the schedule to add signing keys to
func (tx *ScheduleSignTransaction) GetScheduleID() ScheduleID {
	if tx.scheduleID == nil {
		return ScheduleID{}
	}

	return *tx.scheduleID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *ScheduleSignTransaction) Sign(privateKey PrivateKey) *ScheduleSignTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *ScheduleSignTransaction) SignWithOperator(client *Client) (*ScheduleSignTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *ScheduleSignTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleSignTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *ScheduleSignTransaction) AddSignature(publicKey PublicKey, signature []byte) *ScheduleSignTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *ScheduleSignTransaction) SetGrpcDeadline(deadline *time.Duration) *ScheduleSignTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *ScheduleSignTransaction) Freeze() (*ScheduleSignTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *ScheduleSignTransaction) FreezeWith(client *Client) (*ScheduleSignTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this ScheduleSignTransaction.
func (tx *ScheduleSignTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleSignTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *ScheduleSignTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ScheduleSignTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this ScheduleSignTransaction.
func (tx *ScheduleSignTransaction) SetTransactionMemo(memo string) *ScheduleSignTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this ScheduleSignTransaction.
func (tx *ScheduleSignTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleSignTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this ScheduleSignTransaction.
func (tx *ScheduleSignTransaction) SetTransactionID(transactionID TransactionID) *ScheduleSignTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this ScheduleSignTransaction.
func (tx *ScheduleSignTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleSignTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *ScheduleSignTransaction) SetMaxRetry(count int) *ScheduleSignTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *ScheduleSignTransaction) SetMaxBackoff(max time.Duration) *ScheduleSignTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *ScheduleSignTransaction) SetMinBackoff(min time.Duration) *ScheduleSignTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *ScheduleSignTransaction) SetLogLevel(level LogLevel) *ScheduleSignTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *ScheduleSignTransaction) getName() string {
	return "ScheduleSignTransaction"
}

func (tx *ScheduleSignTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.scheduleID != nil {
		if err := tx.scheduleID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *ScheduleSignTransaction) _Build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ScheduleSign{
			ScheduleSign: tx.buildProtoBody(),
		},
	}
}

func (tx *ScheduleSignTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleSignTransaction")
}

func (tx *ScheduleSignTransaction) buildProtoBody() *services.ScheduleSignTransactionBody {
	body := &services.ScheduleSignTransactionBody{}
	if tx.scheduleID != nil {
		body.ScheduleID = tx.scheduleID._ToProtobuf()
	}

	return body
}

func (tx *ScheduleSignTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().SignSchedule,
	}
}
