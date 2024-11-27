package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ScheduleCreateTransaction Creates a new schedule entity (or simply, schedule) in the network's action queue.
// Upon SUCCESS, the receipt contains the `ScheduleID` of the created schedule. A schedule
// entity includes a scheduledTransactionBody to be executed.
// When the schedule has collected enough signing Ed25519 keys to satisfy the schedule's signing
// requirements, the schedule can be executed.
type ScheduleCreateTransaction struct {
	*Transaction[*ScheduleCreateTransaction]
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
	tx := &ScheduleCreateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _ScheduleCreateTransactionFromProtobuf(tx Transaction[*ScheduleCreateTransaction], pb *services.TransactionBody) ScheduleCreateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetScheduleCreate().GetAdminKey())
	var expirationTime time.Time
	if pb.GetScheduleCreate().GetExpirationTime() != nil {
		expirationTime = _TimeFromProtobuf(pb.GetScheduleCreate().GetExpirationTime())
	}

	scheduleCreateTransaction := ScheduleCreateTransaction{
		payerAccountID:  _AccountIDFromProtobuf(pb.GetScheduleCreate().GetPayerAccountID()),
		adminKey:        key,
		schedulableBody: pb.GetScheduleCreate().GetScheduledTransactionBody(),
		memo:            pb.GetScheduleCreate().GetMemo(),
		expirationTime:  &expirationTime,
		waitForExpiry:   pb.GetScheduleCreate().WaitForExpiry,
	}

	tx.childTransaction = &scheduleCreateTransaction
	scheduleCreateTransaction.Transaction = &tx
	return scheduleCreateTransaction
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

// SetAdminKey Sets an optional Hiero key which can be used to sign a ScheduleDelete and remove the schedule
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

// GetAdminKey returns the optional Hiero key which can be used to sign a ScheduleDelete and remove the schedule
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
func (tx *ScheduleCreateTransaction) SetScheduledTransaction(scheduledTx TransactionInterface) (*ScheduleCreateTransaction, error) {
	tx._RequireNotFrozen()

	scheduled, err := scheduledTx.constructScheduleProtobuf()
	if err != nil {
		return tx, err
	}

	tx.schedulableBody = scheduled
	return tx, nil
}

// ----------- Overridden functions ----------------

func (tx ScheduleCreateTransaction) getName() string {
	return "ScheduleCreateTransaction"
}

func (tx ScheduleCreateTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx ScheduleCreateTransaction) build() *services.TransactionBody {
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

func (tx ScheduleCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleCreateTransaction`")
}

func (tx ScheduleCreateTransaction) buildProtoBody() *services.ScheduleCreateTransactionBody {
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

func (tx ScheduleCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().CreateSchedule,
	}
}

func (tx ScheduleCreateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx ScheduleCreateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
