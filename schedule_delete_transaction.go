package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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
	"github.com/hashgraph/hedera-sdk-go/v2/generated/services"

	"time"
)

// ScheduleDeleteTransaction Marks a schedule in the network's action queue as deleted. Must be signed by the admin key of the
// target schedule.  A deleted schedule cannot receive any additional signing keys, nor will it be
// executed.
type ScheduleDeleteTransaction struct {
	Transaction
	scheduleID *ScheduleID
}

// NewScheduleDeleteTransaction creates ScheduleDeleteTransaction which marks a schedule in the network's action queue as deleted.
// Must be signed by the admin key of the target schedule.
// A deleted schedule cannot receive any additional signing keys, nor will it be executed.
func NewScheduleDeleteTransaction() *ScheduleDeleteTransaction {
	tx := ScheduleDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _ScheduleDeleteTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *ScheduleDeleteTransaction {
	return &ScheduleDeleteTransaction{
		Transaction: tx,
		scheduleID:  _ScheduleIDFromProtobuf(pb.GetScheduleDelete().GetScheduleID()),
	}
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
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *ScheduleDeleteTransaction) SignWithOperator(client *Client) (*ScheduleDeleteTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *ScheduleDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleDeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *ScheduleDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ScheduleDeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *ScheduleDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *ScheduleDeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *ScheduleDeleteTransaction) Freeze() (*ScheduleDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *ScheduleDeleteTransaction) FreezeWith(client *Client) (*ScheduleDeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleDeleteTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *ScheduleDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ScheduleDeleteTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetTransactionMemo(memo string) *ScheduleDeleteTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleDeleteTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *ScheduleDeleteTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetTransactionID(transactionID TransactionID) *ScheduleDeleteTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this ScheduleDeleteTransaction.
func (tx *ScheduleDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleDeleteTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *ScheduleDeleteTransaction) SetMaxRetry(count int) *ScheduleDeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *ScheduleDeleteTransaction) SetMaxBackoff(max time.Duration) *ScheduleDeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *ScheduleDeleteTransaction) SetMinBackoff(min time.Duration) *ScheduleDeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *ScheduleDeleteTransaction) SetLogLevel(level LogLevel) *ScheduleDeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *ScheduleDeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *ScheduleDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

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
		Memo:                     tx.Transaction.memo,
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
		Memo:           tx.Transaction.memo,
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
func (tx *ScheduleDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
