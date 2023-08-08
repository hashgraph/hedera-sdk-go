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

	"github.com/hashgraph/hedera-protobufs-go/services"

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
	transaction := ScheduleDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _ScheduleDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *ScheduleDeleteTransaction {
	return &ScheduleDeleteTransaction{
		Transaction: transaction,
		scheduleID:  _ScheduleIDFromProtobuf(pb.GetScheduleDelete().GetScheduleID()),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *ScheduleDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *ScheduleDeleteTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetScheduleID Sets the ScheduleID of the scheduled transaction to be deleted
func (transaction *ScheduleDeleteTransaction) SetScheduleID(scheduleID ScheduleID) *ScheduleDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.scheduleID = &scheduleID
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetScheduleID() ScheduleID {
	if transaction.scheduleID == nil {
		return ScheduleID{}
	}

	return *transaction.scheduleID
}

func (transaction *ScheduleDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
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

func (transaction *ScheduleDeleteTransaction) _Build() *services.TransactionBody {
	body := &services.ScheduleDeleteTransactionBody{}
	if transaction.scheduleID != nil {
		body.ScheduleID = transaction.scheduleID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ScheduleDelete{
			ScheduleDelete: body,
		},
	}
}

func (transaction *ScheduleDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *ScheduleDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ScheduleDeleteTransactionBody{}
	if transaction.scheduleID != nil {
		body.ScheduleID = transaction.scheduleID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ScheduleDelete{
			ScheduleDelete: body,
		},
	}, nil
}

func _ScheduleDeleteTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().DeleteSchedule,
	}
}

func (transaction *ScheduleDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ScheduleDeleteTransaction) Sign(
	privateKey PrivateKey,
) *ScheduleDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *ScheduleDeleteTransaction) SignWithOperator(
	client *Client,
) (*ScheduleDeleteTransaction, error) {
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
func (transaction *ScheduleDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ScheduleDeleteTransaction) Execute(
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
		_ScheduleDeleteTransactionGetMethod,
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

func (transaction *ScheduleDeleteTransaction) Freeze() (*ScheduleDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ScheduleDeleteTransaction) FreezeWith(client *Client) (*ScheduleDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ScheduleDeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ScheduleDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ScheduleDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *ScheduleDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ScheduleDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *ScheduleDeleteTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this	ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionMemo(memo string) *ScheduleDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *ScheduleDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionID(transactionID TransactionID) *ScheduleDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *ScheduleDeleteTransaction) SetMaxRetry(count int) *ScheduleDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *ScheduleDeleteTransaction) SetMaxBackoff(max time.Duration) *ScheduleDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *ScheduleDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *ScheduleDeleteTransaction) SetMinBackoff(min time.Duration) *ScheduleDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *ScheduleDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *ScheduleDeleteTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ScheduleDeleteTransaction:%d", timestamp.UnixNano())
}

func (transaction *ScheduleDeleteTransaction) SetLogLevel(level LogLevel) *ScheduleDeleteTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
