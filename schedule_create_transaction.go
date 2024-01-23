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
	"errors"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// ScheduleCreateTransaction Creates a new schedule entity (or simply, schedule) in the network's action queue.
// Upon SUCCESS, the receipt contains the `ScheduleID` of the created schedule. A schedule
// entity includes a scheduledTransactionBody to be executed.
// When the schedule has collected enough signing Ed25519 keys to satisfy the schedule's signing
// requirements, the schedule can be executed.
type ScheduleCreateTransaction struct {
	Transaction
	payerAccountID  *AccountID
	adminKey        Key
	schedulableBody *services.SchedulableTransactionBody
	memo            string
	expirationTime  *time.Time
	waitForExpiry   bool
}

// NewScheduleCreateTransaction creates ScheduleCreateTransaction which creates a new schedule entity (or simply, schedule) in the network's action queue.
// Upon SUCCESS, the receipt contains the `ScheduleID` of the created schedule. A schedule
// entity includes a scheduledTransactionBody to be executed.
// When the schedule has collected enough signing Ed25519 keys to satisfy the schedule's signing
// requirements, the schedule can be executed.
func NewScheduleCreateTransaction() *ScheduleCreateTransaction {
	tx := ScheduleCreateTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _ScheduleCreateTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *ScheduleCreateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetScheduleCreate().GetAdminKey())
	var expirationTime time.Time
	if pb.GetScheduleCreate().GetExpirationTime() != nil {
		expirationTime = _TimeFromProtobuf(pb.GetScheduleCreate().GetExpirationTime())
	}

	return &ScheduleCreateTransaction{
		Transaction:     tx,
		payerAccountID:  _AccountIDFromProtobuf(pb.GetScheduleCreate().GetPayerAccountID()),
		adminKey:        key,
		schedulableBody: pb.GetScheduleCreate().GetScheduledTransactionBody(),
		memo:            pb.GetScheduleCreate().GetMemo(),
		expirationTime:  &expirationTime,
		waitForExpiry:   pb.GetScheduleCreate().WaitForExpiry,
	}
}

// SetPayerAccountID Sets an optional id of the account to be charged the service fee for the scheduled transaction at
// the consensus time that it executes (if ever); defaults to the ScheduleCreate payer if not
// given
func (tx *ScheduleCreateTransaction) SetPayerAccountID(payerAccountID AccountID) *ScheduleCreateTransaction {
	tx._RequireNotFrozen()
	tx.payerAccountID = &payerAccountID

	return tx
}

// GetPayerAccountID returns the optional id of the account to be charged the service fee for the scheduled transaction
func (tx *ScheduleCreateTransaction) GetPayerAccountID() AccountID {
	if tx.payerAccountID == nil {
		return AccountID{}
	}

	return *tx.payerAccountID
}

// SetAdminKey Sets an optional Hedera key which can be used to sign a ScheduleDelete and remove the schedule
func (tx *ScheduleCreateTransaction) SetAdminKey(key Key) *ScheduleCreateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = key

	return tx
}

// SetExpirationTime Sets an optional timestamp for specifying when the transaction should be evaluated for execution and then expire.
// Defaults to 30 minutes after the transaction's consensus timestamp.
func (tx *ScheduleCreateTransaction) SetExpirationTime(time time.Time) *ScheduleCreateTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &time

	return tx
}

// GetExpirationTime returns the optional timestamp for specifying when the transaction should be evaluated for execution and then expire.
func (tx *ScheduleCreateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}

	return time.Time{}
}

// SetWaitForExpiry
// When set to true, the transaction will be evaluated for execution at expiration_time instead
// of when all required signatures are received.
// When set to false, the transaction will execute immediately after sufficient signatures are received
// to sign the contained transaction. During the initial ScheduleCreate transaction or via ScheduleSign transactions.
// Defaults to false.
func (tx *ScheduleCreateTransaction) SetWaitForExpiry(wait bool) *ScheduleCreateTransaction {
	tx._RequireNotFrozen()
	tx.waitForExpiry = wait

	return tx
}

// GetWaitForExpiry returns true if the transaction will be evaluated for execution at expiration_time instead
// of when all required signatures are received.
func (tx *ScheduleCreateTransaction) GetWaitForExpiry() bool {
	return tx.waitForExpiry
}

func (tx *ScheduleCreateTransaction) _SetSchedulableTransactionBody(txBody *services.SchedulableTransactionBody) *ScheduleCreateTransaction {
	tx._RequireNotFrozen()
	tx.schedulableBody = txBody

	return tx
}

// GetAdminKey returns the optional Hedera key which can be used to sign a ScheduleDelete and remove the schedule
func (tx *ScheduleCreateTransaction) GetAdminKey() *Key {
	if tx.adminKey == nil {
		return nil
	}
	return &tx.adminKey
}

// SetScheduleMemo Sets an optional memo with a UTF-8 encoding of no more than 100 bytes which does not contain the zero byte.
func (tx *ScheduleCreateTransaction) SetScheduleMemo(memo string) *ScheduleCreateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo

	return tx
}

// GetScheduleMemo returns the optional memo with a UTF-8 encoding of no more than 100 bytes which does not contain the zero byte.
func (tx *ScheduleCreateTransaction) GetScheduleMemo() string {
	return tx.memo
}

// SetScheduledTransaction Sets the scheduled transaction
func (tx *ScheduleCreateTransaction) SetScheduledTransaction(scheduledTx ITransaction) (*ScheduleCreateTransaction, error) {
	tx._RequireNotFrozen()

	scheduled, err := scheduledTx._ConstructScheduleProtobuf()
	if err != nil {
		return tx, err
	}

	tx.schedulableBody = scheduled
	return tx, nil
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *ScheduleCreateTransaction) Sign(privateKey PrivateKey) *ScheduleCreateTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *ScheduleCreateTransaction) SignWithOperator(client *Client) (*ScheduleCreateTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *ScheduleCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleCreateTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *ScheduleCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *ScheduleCreateTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// SetGrpcDeadline When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *ScheduleCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *ScheduleCreateTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *ScheduleCreateTransaction) Freeze() (*ScheduleCreateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *ScheduleCreateTransaction) FreezeWith(client *Client) (*ScheduleCreateTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee for this ScheduleCreateTransaction.
func (tx *ScheduleCreateTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleCreateTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *ScheduleCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ScheduleCreateTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this ScheduleCreateTransaction.
func (tx *ScheduleCreateTransaction) SetTransactionMemo(memo string) *ScheduleCreateTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this ScheduleCreateTransaction.
func (tx *ScheduleCreateTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleCreateTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *ScheduleCreateTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this ScheduleCreateTransaction.
func (tx *ScheduleCreateTransaction) SetTransactionID(transactionID TransactionID) *ScheduleCreateTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this ScheduleCreateTransaction.
func (tx *ScheduleCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleCreateTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *ScheduleCreateTransaction) SetMaxRetry(count int) *ScheduleCreateTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *ScheduleCreateTransaction) SetMaxBackoff(max time.Duration) *ScheduleCreateTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *ScheduleCreateTransaction) SetMinBackoff(min time.Duration) *ScheduleCreateTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *ScheduleCreateTransaction) SetLogLevel(level LogLevel) *ScheduleCreateTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *ScheduleCreateTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *ScheduleCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *ScheduleCreateTransaction) getName() string {
	return "ScheduleCreateTransaction"
}

func (tx *ScheduleCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.payerAccountID != nil {
		if err := tx.payerAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *ScheduleCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ScheduleCreate{
			ScheduleCreate: tx.buildProtoBody(),
		},
	}
}

func (tx *ScheduleCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleCreateTransaction`")
}

func (tx *ScheduleCreateTransaction) buildProtoBody() *services.ScheduleCreateTransactionBody {
	body := &services.ScheduleCreateTransactionBody{
		Memo:          tx.memo,
		WaitForExpiry: tx.waitForExpiry,
	}

	if tx.payerAccountID != nil {
		body.PayerAccountID = tx.payerAccountID._ToProtobuf()
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.schedulableBody != nil {
		body.ScheduledTransactionBody = tx.schedulableBody
	}

	if tx.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*tx.expirationTime)
	}

	return body
}

func (tx *ScheduleCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().CreateSchedule,
	}
}

func (tx *ScheduleCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
