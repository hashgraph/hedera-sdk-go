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
	"errors"
	"fmt"
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
	transaction := ScheduleCreateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction._SetDefaultMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _ScheduleCreateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *ScheduleCreateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetScheduleCreate().GetAdminKey())
	var expirationTime time.Time
	if pb.GetScheduleCreate().GetExpirationTime() != nil {
		expirationTime = _TimeFromProtobuf(pb.GetScheduleCreate().GetExpirationTime())
	}

	return &ScheduleCreateTransaction{
		Transaction:     transaction,
		payerAccountID:  _AccountIDFromProtobuf(pb.GetScheduleCreate().GetPayerAccountID()),
		adminKey:        key,
		schedulableBody: pb.GetScheduleCreate().GetScheduledTransactionBody(),
		memo:            pb.GetScheduleCreate().GetMemo(),
		expirationTime:  &expirationTime,
		waitForExpiry:   pb.GetScheduleCreate().WaitForExpiry,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *ScheduleCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *ScheduleCreateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetPayerAccountID Sets an optional id of the account to be charged the service fee for the scheduled transaction at
// the consensus time that it executes (if ever); defaults to the ScheduleCreate payer if not
// given
func (transaction *ScheduleCreateTransaction) SetPayerAccountID(payerAccountID AccountID) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.payerAccountID = &payerAccountID

	return transaction
}

// GetPayerAccountID returns the optional id of the account to be charged the service fee for the scheduled transaction
func (transaction *ScheduleCreateTransaction) GetPayerAccountID() AccountID {
	if transaction.payerAccountID == nil {
		return AccountID{}
	}

	return *transaction.payerAccountID
}

// SetAdminKey Sets an optional Hedera key which can be used to sign a ScheduleDelete and remove the schedule
func (transaction *ScheduleCreateTransaction) SetAdminKey(key Key) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.adminKey = key

	return transaction
}

// SetExpirationTime Sets an optional timestamp for specifying when the transaction should be evaluated for execution and then expire.
// Defaults to 30 minutes after the transaction's consensus timestamp.
func (transaction *ScheduleCreateTransaction) SetExpirationTime(time time.Time) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.expirationTime = &time

	return transaction
}

// GetExpirationTime returns the optional timestamp for specifying when the transaction should be evaluated for execution and then expire.
func (transaction *ScheduleCreateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

// SetWaitForExpiry
// When set to true, the transaction will be evaluated for execution at expiration_time instead
// of when all required signatures are received.
// When set to false, the transaction will execute immediately after sufficient signatures are received
// to sign the contained transaction. During the initial ScheduleCreate transaction or via ScheduleSign transactions.
// Defaults to false.
func (transaction *ScheduleCreateTransaction) SetWaitForExpiry(wait bool) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.waitForExpiry = wait

	return transaction
}

// GetWaitForExpiry returns true if the transaction will be evaluated for execution at expiration_time instead
// of when all required signatures are received.
func (transaction *ScheduleCreateTransaction) GetWaitForExpiry() bool {
	return transaction.waitForExpiry
}

func (transaction *ScheduleCreateTransaction) _SetSchedulableTransactionBody(txBody *services.SchedulableTransactionBody) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.schedulableBody = txBody

	return transaction
}

// GetAdminKey returns the optional Hedera key which can be used to sign a ScheduleDelete and remove the schedule
func (transaction *ScheduleCreateTransaction) GetAdminKey() *Key {
	if transaction.adminKey == nil {
		return nil
	}
	return &transaction.adminKey
}

// SetScheduleMemo Sets an optional memo with a UTF-8 encoding of no more than 100 bytes which does not contain the zero byte.
func (transaction *ScheduleCreateTransaction) SetScheduleMemo(memo string) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo

	return transaction
}

// GetScheduleMemo returns the optional memo with a UTF-8 encoding of no more than 100 bytes which does not contain the zero byte.
func (transaction *ScheduleCreateTransaction) GetScheduleMemo() string {
	return transaction.memo
}

// SetScheduledTransaction Sets the scheduled transaction
func (transaction *ScheduleCreateTransaction) SetScheduledTransaction(tx ITransaction) (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := tx._ConstructScheduleProtobuf()
	if err != nil {
		return transaction, err
	}

	transaction.schedulableBody = scheduled
	return transaction, nil
}

func (transaction *ScheduleCreateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.payerAccountID != nil {
		if err := transaction.payerAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *ScheduleCreateTransaction) _Build() *services.TransactionBody {
	body := &services.ScheduleCreateTransactionBody{
		Memo:          transaction.memo,
		WaitForExpiry: transaction.waitForExpiry,
	}

	if transaction.payerAccountID != nil {
		body.PayerAccountID = transaction.payerAccountID._ToProtobuf()
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.schedulableBody != nil {
		body.ScheduledTransactionBody = transaction.schedulableBody
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ScheduleCreate{
			ScheduleCreate: body,
		},
	}
}

func (transaction *ScheduleCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleCreateTransaction`")
}
func _ScheduleCreateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().CreateSchedule,
	}
}

func (transaction *ScheduleCreateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ScheduleCreateTransaction) Sign(
	privateKey PrivateKey,
) *ScheduleCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *ScheduleCreateTransaction) SignWithOperator(
	client *Client,
) (*ScheduleCreateTransaction, error) {
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
func (transaction *ScheduleCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleCreateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ScheduleCreateTransaction) Execute(
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
		_ScheduleCreateTransactionGetMethod,
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
		TransactionID:          transaction.GetTransactionID(),
		NodeID:                 resp.(TransactionResponse).NodeID,
		Hash:                   resp.(TransactionResponse).Hash,
		ScheduledTransactionId: transaction.GetTransactionID(),
	}, nil
}

func (transaction *ScheduleCreateTransaction) Freeze() (*ScheduleCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ScheduleCreateTransaction) FreezeWith(client *Client) (*ScheduleCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ScheduleCreateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	// transaction.transactionIDs[0] = transaction.transactionIDs[0].SetScheduled(true)

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ScheduleCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ScheduleCreateTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *ScheduleCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *ScheduleCreateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetTransactionMemo(memo string) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *ScheduleCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetTransactionID(transactionID TransactionID) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *ScheduleCreateTransaction) SetMaxRetry(count int) *ScheduleCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *ScheduleCreateTransaction) SetMaxBackoff(max time.Duration) *ScheduleCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *ScheduleCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *ScheduleCreateTransaction) SetMinBackoff(min time.Duration) *ScheduleCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *ScheduleCreateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *ScheduleCreateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ScheduleCreateTransaction:%d", timestamp.UnixNano())
}

func (transaction *ScheduleCreateTransaction) SetLogLevel(level LogLevel) *ScheduleCreateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
