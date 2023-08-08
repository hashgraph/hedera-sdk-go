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

// ContractDeleteTransaction marks a contract as deleted and transfers its remaining hBars, if any, to a
// designated receiver. After a contract is deleted, it can no longer be called.
type ContractDeleteTransaction struct {
	Transaction
	contractID        *ContractID
	transferContactID *ContractID
	transferAccountID *AccountID
	permanentRemoval  bool
}

// NewContractDeleteTransaction creates ContractDeleteTransaction which marks a contract as deleted and transfers its remaining hBars, if any, to a
// designated receiver. After a contract is deleted, it can no longer be called.
func NewContractDeleteTransaction() *ContractDeleteTransaction {
	transaction := ContractDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _ContractDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *ContractDeleteTransaction {
	return &ContractDeleteTransaction{
		Transaction:       transaction,
		contractID:        _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetContractID()),
		transferContactID: _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferContractID()),
		transferAccountID: _AccountIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferAccountID()),
		permanentRemoval:  pb.GetContractDeleteInstance().GetPermanentRemoval(),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *ContractDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractDeleteTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// Sets the contract ID which will be deleted.
func (transaction *ContractDeleteTransaction) SetContractID(contractID ContractID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.contractID = &contractID
	return transaction
}

// Returns the contract ID which will be deleted.
func (transaction *ContractDeleteTransaction) GetContractID() ContractID {
	if transaction.contractID == nil {
		return ContractID{}
	}

	return *transaction.contractID
}

// Sets the contract ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) SetTransferContractID(transferContactID ContractID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.transferContactID = &transferContactID
	return transaction
}

// Returns the contract ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) GetTransferContractID() ContractID {
	if transaction.transferContactID == nil {
		return ContractID{}
	}

	return *transaction.transferContactID
}

// Sets the account ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) SetTransferAccountID(accountID AccountID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.transferAccountID = &accountID

	return transaction
}

// Returns the account ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) GetTransferAccountID() AccountID {
	if transaction.transferAccountID == nil {
		return AccountID{}
	}

	return *transaction.transferAccountID
}

// SetPermanentRemoval
// If set to true, means this is a "synthetic" system transaction being used to
// alert mirror nodes that the contract is being permanently removed from the ledger.
// IMPORTANT: User transactions cannot set this field to true, as permanent
// removal is always managed by the ledger itself. Any ContractDeleteTransaction
// submitted to HAPI with permanent_removal=true will be rejected with precheck status
// PERMANENT_REMOVAL_REQUIRES_SYSTEM_INITIATION.
func (transaction *ContractDeleteTransaction) SetPermanentRemoval(remove bool) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.permanentRemoval = remove

	return transaction
}

// GetPermanentRemoval returns true if this is a "synthetic" system transaction.
func (transaction *ContractDeleteTransaction) GetPermanentRemoval() bool {
	return transaction.permanentRemoval
}

func (transaction *ContractDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.contractID != nil {
		if err := transaction.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.transferContactID != nil {
		if err := transaction.transferContactID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.transferAccountID != nil {
		if err := transaction.transferAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *ContractDeleteTransaction) _Build() *services.TransactionBody {
	body := &services.ContractDeleteTransactionBody{
		PermanentRemoval: transaction.permanentRemoval,
	}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	if transaction.transferContactID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: transaction.transferContactID._ToProtobuf(),
		}
	}

	if transaction.transferAccountID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: transaction.transferAccountID._ToProtobuf(),
		}
	}

	pb := services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: body,
		},
	}

	return &pb
}

func (transaction *ContractDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ContractDeleteTransactionBody{
		PermanentRemoval: transaction.permanentRemoval,
	}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	if transaction.transferContactID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: transaction.transferContactID._ToProtobuf(),
		}
	}

	if transaction.transferAccountID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: transaction.transferAccountID._ToProtobuf(),
		}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: body,
		},
	}, nil
}

func _ContractDeleteTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().DeleteContract,
	}
}

func (transaction *ContractDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractDeleteTransaction) Sign(
	privateKey PrivateKey,
) *ContractDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *ContractDeleteTransaction) SignWithOperator(
	client *Client,
) (*ContractDeleteTransaction, error) {
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
func (transaction *ContractDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractDeleteTransaction) Execute(
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
		_ContractDeleteTransactionGetMethod,
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

func (transaction *ContractDeleteTransaction) Freeze() (*ContractDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractDeleteTransaction) FreezeWith(client *Client) (*ContractDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ContractDeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ContractDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ContractDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *ContractDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *ContractDeleteTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionMemo(memo string) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *ContractDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *ContractDeleteTransaction) SetMaxRetry(count int) *ContractDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *ContractDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractDeleteTransaction {
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
func (transaction *ContractDeleteTransaction) SetMaxBackoff(max time.Duration) *ContractDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *ContractDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *ContractDeleteTransaction) SetMinBackoff(min time.Duration) *ContractDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *ContractDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *ContractDeleteTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ContractDeleteTransaction:%d", timestamp.UnixNano())
}

func (transaction *ContractDeleteTransaction) SetLogLevel(level LogLevel) *ContractDeleteTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
