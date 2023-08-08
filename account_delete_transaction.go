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

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountDeleteTransaction
// Mark an account as deleted, moving all its current hbars to another account. It will remain in
// the ledger, marked as deleted, until it expires. Transfers into it a deleted account fail. But a
// deleted account can still have its expiration extended in the normal way.
type AccountDeleteTransaction struct {
	Transaction
	transferAccountID *AccountID
	deleteAccountID   *AccountID
}

func _AccountDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *AccountDeleteTransaction {
	return &AccountDeleteTransaction{
		Transaction:       transaction,
		transferAccountID: _AccountIDFromProtobuf(pb.GetCryptoDelete().GetTransferAccountID()),
		deleteAccountID:   _AccountIDFromProtobuf(pb.GetCryptoDelete().GetDeleteAccountID()),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *AccountDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountDeleteTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// NewAccountDeleteTransaction creates AccountDeleteTransaction which marks an account as deleted, moving all its current hbars to another account. It will remain in
// the ledger, marked as deleted, until it expires. Transfers into it a deleted account fail. But a
// deleted account can still have its expiration extended in the normal way.
func NewAccountDeleteTransaction() *AccountDeleteTransaction {
	transaction := AccountDeleteTransaction{
		Transaction: _NewTransaction(),
	}

	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

// SetNodeAccountID sets the _Node AccountID for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetAccountID(accountID AccountID) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.deleteAccountID = &accountID
	return transaction
}

// GetAccountID returns the AccountID which will be deleted.
func (transaction *AccountDeleteTransaction) GetAccountID() AccountID {
	if transaction.deleteAccountID == nil {
		return AccountID{}
	}

	return *transaction.deleteAccountID
}

// SetTransferAccountID sets the AccountID which will receive all remaining hbars.
func (transaction *AccountDeleteTransaction) SetTransferAccountID(transferAccountID AccountID) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.transferAccountID = &transferAccountID
	return transaction
}

// GetTransferAccountID returns the AccountID which will receive all remaining hbars.
func (transaction *AccountDeleteTransaction) GetTransferAccountID() AccountID {
	if transaction.transferAccountID == nil {
		return AccountID{}
	}

	return *transaction.transferAccountID
}

func (transaction *AccountDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.deleteAccountID != nil {
		if err := transaction.deleteAccountID.ValidateChecksum(client); err != nil {
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

func (transaction *AccountDeleteTransaction) _Build() *services.TransactionBody {
	body := &services.CryptoDeleteTransactionBody{}

	if transaction.transferAccountID != nil {
		body.TransferAccountID = transaction.transferAccountID._ToProtobuf()
	}

	if transaction.deleteAccountID != nil {
		body.DeleteAccountID = transaction.deleteAccountID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoDelete{
			CryptoDelete: body,
		},
	}
}

func (transaction *AccountDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.CryptoDeleteTransactionBody{}

	if transaction.transferAccountID != nil {
		body.TransferAccountID = transaction.transferAccountID._ToProtobuf()
	}

	if transaction.deleteAccountID != nil {
		body.DeleteAccountID = transaction.deleteAccountID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoDelete{
			CryptoDelete: body,
		},
	}, nil
}

func _AccountDeleteTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().CryptoDelete,
	}
}

func (transaction *AccountDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountDeleteTransaction) Sign(
	privateKey PrivateKey,
) *AccountDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *AccountDeleteTransaction) SignWithOperator(
	client *Client,
) (*AccountDeleteTransaction, error) {
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
func (transaction *AccountDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountDeleteTransaction) Execute(
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
		_AccountDeleteTransactionGetMethod,
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
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.(TransactionResponse).NodeID,
			Hash:          make([]byte, 0),
		}, err
	}

	return TransactionResponse{
		TransactionID:  transaction.GetTransactionID(),
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
	}, nil
}

func (transaction *AccountDeleteTransaction) Freeze() (*AccountDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountDeleteTransaction) FreezeWith(client *Client) (*AccountDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &AccountDeleteTransaction{}, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountDeleteTransaction) SetMaxTransactionFee(fee Hbar) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *AccountDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *AccountDeleteTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionMemo(memo string) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *AccountDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionID(transactionID TransactionID) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *AccountDeleteTransaction) SetMaxRetry(count int) *AccountDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *AccountDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountDeleteTransaction {
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

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *AccountDeleteTransaction) SetMaxBackoff(max time.Duration) *AccountDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *AccountDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *AccountDeleteTransaction) SetMinBackoff(min time.Duration) *AccountDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *AccountDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *AccountDeleteTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountDeleteTransaction:%d", timestamp.UnixNano())
}

func (transaction *AccountDeleteTransaction) SetLogLevel(level LogLevel) *AccountDeleteTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
