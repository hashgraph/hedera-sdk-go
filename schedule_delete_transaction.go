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
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// ScheduleDeleteTransaction Marks a schedule in the network's action queue as deleted. Must be signed by the admin key of the
// target schedule.  A deleted schedule cannot receive any additional signing keys, nor will it be
// executed.
type ScheduleDeleteTransaction struct {
	transaction
	scheduleID *ScheduleID
}

// NewScheduleDeleteTransaction creates ScheduleDeleteTransaction which marks a schedule in the network's action queue as deleted.
// Must be signed by the admin key of the target schedule.
// A deleted schedule cannot receive any additional signing keys, nor will it be executed.
func NewScheduleDeleteTransaction() *ScheduleDeleteTransaction {
	tx := ScheduleDeleteTransaction{
		transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _ScheduleDeleteTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *ScheduleDeleteTransaction {
	resultTx := &ScheduleDeleteTransaction{
		transaction: tx,
		scheduleID:  _ScheduleIDFromProtobuf(pb.GetScheduleDelete().GetScheduleID()),
	}
	resultTx.e = resultTx
	return resultTx
}

// SetScheduleID Sets the ScheduleID of the scheduled transaction to be deleted
func (tx *ScheduleDeleteTransaction) SetScheduleID(scheduleID ScheduleID) *ScheduleDeleteTransaction {
	tx._RequireNotFrozen()
	tx.scheduleID = &scheduleID
	return tx
}

func (tx *ScheduleDeleteTransaction) GetScheduleID() ScheduleID {
	if tx.scheduleID == nil {
		return ScheduleID{}
	}

	return *tx.scheduleID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *ScheduleDeleteTransaction) Sign(privateKey PrivateKey) *ScheduleDeleteTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *ScheduleDeleteTransaction) SignWithOperator(client *Client) (*ScheduleDeleteTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *ScheduleDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleDeleteTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *ScheduleDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ScheduleDeleteTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *ScheduleDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *ScheduleDeleteTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *ScheduleDeleteTransaction) Freeze() (*ScheduleDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *ScheduleDeleteTransaction) FreezeWith(client *Client) (*ScheduleDeleteTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleDeleteTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *ScheduleDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ScheduleDeleteTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetTransactionMemo(memo string) *ScheduleDeleteTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleDeleteTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetTransactionID(transactionID TransactionID) *ScheduleDeleteTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleDeleteTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *ScheduleDeleteTransaction) SetMaxRetry(count int) *ScheduleDeleteTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *ScheduleDeleteTransaction) SetMaxBackoff(max time.Duration) *ScheduleDeleteTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *ScheduleDeleteTransaction) SetMinBackoff(min time.Duration) *ScheduleDeleteTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *ScheduleDeleteTransaction) SetLogLevel(level LogLevel) *ScheduleDeleteTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *ScheduleDeleteTransaction) getName() string {
	return "ScheduleDeleteTransaction"
}

func (tx *ScheduleDeleteTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *ScheduleDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ScheduleDelete{
			ScheduleDelete: tx.buildProtoBody(),
		},
	}
}

func (tx *ScheduleDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_ScheduleDelete{
			ScheduleDelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *ScheduleDeleteTransaction) buildProtoBody() *services.ScheduleDeleteTransactionBody {
	body := &services.ScheduleDeleteTransactionBody{}
	if tx.scheduleID != nil {
		body.ScheduleID = tx.scheduleID._ToProtobuf()
	}

	return body
}

func (tx *ScheduleDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().DeleteSchedule,
	}
}
