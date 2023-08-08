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

// Undelete a file or smart contract that was deleted by AdminDelete.
// Can only be done with a Hedera admin.
type SystemUndeleteTransaction struct {
	Transaction
	contractID *ContractID
	fileID     *FileID
}

// NewSystemUndeleteTransaction creates a SystemUndeleteTransaction transaction which can be
// used to construct and execute a System Undelete Transaction.
func NewSystemUndeleteTransaction() *SystemUndeleteTransaction {
	transaction := SystemUndeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _SystemUndeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *SystemUndeleteTransaction {
	return &SystemUndeleteTransaction{
		Transaction: transaction,
		contractID:  _ContractIDFromProtobuf(pb.GetSystemUndelete().GetContractID()),
		fileID:      _FileIDFromProtobuf(pb.GetSystemUndelete().GetFileID()),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *SystemUndeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *SystemUndeleteTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetContractID sets the ContractID of the contract whose deletion is being undone.
func (transaction *SystemUndeleteTransaction) SetContractID(contractID ContractID) *SystemUndeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.contractID = &contractID
	return transaction
}

// GetContractID returns the ContractID of the contract whose deletion is being undone.
func (transaction *SystemUndeleteTransaction) GetContractID() ContractID {
	if transaction.contractID == nil {
		return ContractID{}
	}

	return *transaction.contractID
}

// SetFileID sets the FileID of the file whose deletion is being undone.
func (transaction *SystemUndeleteTransaction) SetFileID(fileID FileID) *SystemUndeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.fileID = &fileID
	return transaction
}

// GetFileID returns the FileID of the file whose deletion is being undone.
func (transaction *SystemUndeleteTransaction) GetFileID() FileID {
	if transaction.fileID == nil {
		return FileID{}
	}

	return *transaction.fileID
}

func (transaction *SystemUndeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.contractID != nil {
		if err := transaction.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.fileID != nil {
		if err := transaction.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *SystemUndeleteTransaction) _Build() *services.TransactionBody {
	body := &services.SystemUndeleteTransactionBody{}
	if transaction.contractID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_ContractID{
			ContractID: transaction.contractID._ToProtobuf(),
		}
	}

	if transaction.fileID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_FileID{
			FileID: transaction.fileID._ToProtobuf(),
		}
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_SystemUndelete{
			SystemUndelete: body,
		},
	}
}

func (transaction *SystemUndeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *SystemUndeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.SystemUndeleteTransactionBody{}
	if transaction.contractID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_ContractID{
			ContractID: transaction.contractID._ToProtobuf(),
		}
	}

	if transaction.fileID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_FileID{
			FileID: transaction.fileID._ToProtobuf(),
		}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_SystemUndelete{
			SystemUndelete: body,
		},
	}, nil
}

func _SystemUndeleteTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	if channel._GetContract() == nil {
		return _Method{
			transaction: channel._GetFile().SystemUndelete,
		}
	}

	return _Method{
		transaction: channel._GetContract().SystemUndelete,
	}
}

func (transaction *SystemUndeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *SystemUndeleteTransaction) Sign(
	privateKey PrivateKey,
) *SystemUndeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *SystemUndeleteTransaction) SignWithOperator(
	client *Client,
) (*SystemUndeleteTransaction, error) {
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
func (transaction *SystemUndeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *SystemUndeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *SystemUndeleteTransaction) Execute(
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
		_SystemUndeleteTransactionGetMethod,
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

func (transaction *SystemUndeleteTransaction) Freeze() (*SystemUndeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *SystemUndeleteTransaction) FreezeWith(client *Client) (*SystemUndeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &SystemUndeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *SystemUndeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *SystemUndeleteTransaction) SetMaxTransactionFee(fee Hbar) *SystemUndeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *SystemUndeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *SystemUndeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *SystemUndeleteTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *SystemUndeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionMemo(memo string) *SystemUndeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *SystemUndeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionValidDuration(duration time.Duration) *SystemUndeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionID(transactionID TransactionID) *SystemUndeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *SystemUndeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *SystemUndeleteTransaction) SetMaxRetry(count int) *SystemUndeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *SystemUndeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *SystemUndeleteTransaction {
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

// SetMaxBackoff sets the maximum amount of time to wait between retries.
func (transaction *SystemUndeleteTransaction) SetMaxBackoff(max time.Duration) *SystemUndeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *SystemUndeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *SystemUndeleteTransaction) SetMinBackoff(min time.Duration) *SystemUndeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *SystemUndeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *SystemUndeleteTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("SystemUndeleteTransaction:%d", timestamp.UnixNano())
}

func (transaction *SystemUndeleteTransaction) SetLogLevel(level LogLevel) *SystemUndeleteTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
