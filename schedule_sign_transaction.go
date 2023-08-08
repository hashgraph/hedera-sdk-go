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
	Transaction
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
	transaction := ScheduleSignTransaction{
		Transaction: _NewTransaction(),
	}

	transaction._SetDefaultMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _ScheduleSignTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *ScheduleSignTransaction {
	return &ScheduleSignTransaction{
		Transaction: transaction,
		scheduleID:  _ScheduleIDFromProtobuf(pb.GetScheduleSign().GetScheduleID()),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *ScheduleSignTransaction) SetGrpcDeadline(deadline *time.Duration) *ScheduleSignTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetScheduleID Sets the id of the schedule to add signing keys to
func (transaction *ScheduleSignTransaction) SetScheduleID(scheduleID ScheduleID) *ScheduleSignTransaction {
	transaction._RequireNotFrozen()
	transaction.scheduleID = &scheduleID
	return transaction
}

// GetScheduleID returns the id of the schedule to add signing keys to
func (transaction *ScheduleSignTransaction) GetScheduleID() ScheduleID {
	if transaction.scheduleID == nil {
		return ScheduleID{}
	}

	return *transaction.scheduleID
}

func (transaction *ScheduleSignTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.scheduleID != nil {
		if err := transaction.scheduleID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *ScheduleSignTransaction) _Build() *services.TransactionBody {
	body := &services.ScheduleSignTransactionBody{}
	if transaction.scheduleID != nil {
		body.ScheduleID = transaction.scheduleID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ScheduleSign{
			ScheduleSign: body,
		},
	}
}

func (transaction *ScheduleSignTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleSignTransaction")
}

func _ScheduleSignTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().SignSchedule,
	}
}

func (transaction *ScheduleSignTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ScheduleSignTransaction) Sign(
	privateKey PrivateKey,
) *ScheduleSignTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *ScheduleSignTransaction) SignWithOperator(
	client *Client,
) (*ScheduleSignTransaction, error) {
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
func (transaction *ScheduleSignTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleSignTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ScheduleSignTransaction) Execute(
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
		_ScheduleSignTransactionGetMethod,
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

func (transaction *ScheduleSignTransaction) Freeze() (*ScheduleSignTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ScheduleSignTransaction) FreezeWith(client *Client) (*ScheduleSignTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ScheduleSignTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ScheduleSignTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ScheduleSignTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleSignTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *ScheduleSignTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ScheduleSignTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *ScheduleSignTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetTransactionMemo(memo string) *ScheduleSignTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *ScheduleSignTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleSignTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetTransactionID(transactionID TransactionID) *ScheduleSignTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleSignTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *ScheduleSignTransaction) SetMaxRetry(count int) *ScheduleSignTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *ScheduleSignTransaction) SetMaxBackoff(max time.Duration) *ScheduleSignTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *ScheduleSignTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *ScheduleSignTransaction) SetMinBackoff(min time.Duration) *ScheduleSignTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *ScheduleSignTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *ScheduleSignTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ScheduleSignTransaction:%d", timestamp.UnixNano())
}

func (transaction *ScheduleSignTransaction) SetLogLevel(level LogLevel) *ScheduleSignTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
