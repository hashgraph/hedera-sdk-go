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
	transaction
	transferAccountID *AccountID
	deleteAccountID   *AccountID
}

func _AccountDeleteTransactionFromProtobuf(transaction transaction, pb *services.TransactionBody) *AccountDeleteTransaction {
	resultTx := &AccountDeleteTransaction{
		transaction:       transaction,
		transferAccountID: _AccountIDFromProtobuf(pb.GetCryptoDelete().GetTransferAccountID()),
		deleteAccountID:   _AccountIDFromProtobuf(pb.GetCryptoDelete().GetDeleteAccountID()),
	}
	resultTx.e = resultTx
	return resultTx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *AccountDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountDeleteTransaction {
	transaction.transaction.SetGrpcDeadline(deadline)
	return transaction
}

// NewAccountDeleteTransaction creates AccountDeleteTransaction which marks an account as deleted, moving all its current hbars to another account. It will remain in
// the ledger, marked as deleted, until it expires. Transfers into it a deleted account fail. But a
// deleted account can still have its expiration extended in the normal way.
func NewAccountDeleteTransaction() *AccountDeleteTransaction {
	this := AccountDeleteTransaction{
		transaction: _NewTransaction(),
	}

	this.e = &this
	this._SetDefaultMaxTransactionFee(NewHbar(2))

	return &this
}

// SetNodeAccountID sets the _Node AccountID for this AccountDeleteTransaction.
func (this *AccountDeleteTransaction) SetAccountID(accountID AccountID) *AccountDeleteTransaction {
	this._RequireNotFrozen()
	this.deleteAccountID = &accountID
	return this
}

// GetAccountID returns the AccountID which will be deleted.
func (this *AccountDeleteTransaction) GetAccountID() AccountID {
	if this.deleteAccountID == nil {
		return AccountID{}
	}

	return *this.deleteAccountID
}

// SetTransferAccountID sets the AccountID which will receive all remaining hbars.
func (this *AccountDeleteTransaction) SetTransferAccountID(transferAccountID AccountID) *AccountDeleteTransaction {
	this._RequireNotFrozen()
	this.transferAccountID = &transferAccountID
	return this
}

// GetTransferAccountID returns the AccountID which will receive all remaining hbars.
func (this *AccountDeleteTransaction) GetTransferAccountID() AccountID {
	if this.transferAccountID == nil {
		return AccountID{}
	}

	return *this.transferAccountID
}

func (this *AccountDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

// Sign uses the provided privateKey to sign the transaction.
func (this *AccountDeleteTransaction) Sign(
	privateKey PrivateKey,
) *AccountDeleteTransaction {
	this.transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *AccountDeleteTransaction) SignWithOperator(
	client *Client,
) (*AccountDeleteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *AccountDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountDeleteTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *AccountDeleteTransaction) Freeze() (*AccountDeleteTransaction, error) {
	_, err := this.transaction.Freeze()
	return this, err
}

func (this *AccountDeleteTransaction) FreezeWith(client *Client) (*AccountDeleteTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountDeleteTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountDeleteTransaction) SetMaxTransactionFee(fee Hbar) *AccountDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *AccountDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *AccountDeleteTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountDeleteTransaction.
func (this *AccountDeleteTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountDeleteTransaction.
func (this *AccountDeleteTransaction) SetTransactionMemo(memo string) *AccountDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *AccountDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountDeleteTransaction.
func (this *AccountDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *AccountDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	AccountDeleteTransaction.
func (this *AccountDeleteTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountDeleteTransaction.
func (this *AccountDeleteTransaction) SetTransactionID(transactionID TransactionID) *AccountDeleteTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountDeleteTransaction.
func (this *AccountDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *AccountDeleteTransaction) SetMaxRetry(count int) *AccountDeleteTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *AccountDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountDeleteTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *AccountDeleteTransaction) SetMaxBackoff(max time.Duration) *AccountDeleteTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *AccountDeleteTransaction) SetMinBackoff(min time.Duration) *AccountDeleteTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *AccountDeleteTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountDeleteTransaction:%d", timestamp.UnixNano())
}

func (this *AccountDeleteTransaction) SetLogLevel(level LogLevel) *AccountDeleteTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *AccountDeleteTransaction) getName() string {
	return "AccountDeleteTransaction"
}
func (this *AccountDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.deleteAccountID != nil {
		if err := this.deleteAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if this.transferAccountID != nil {
		if err := this.transferAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *AccountDeleteTransaction) build() *services.TransactionBody {
	body := &services.CryptoDeleteTransactionBody{}

	if this.transferAccountID != nil {
		body.TransferAccountID = this.transferAccountID._ToProtobuf()
	}

	if this.deleteAccountID != nil {
		body.DeleteAccountID = this.deleteAccountID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoDelete{
			CryptoDelete: body,
		},
	}
}

func (this *AccountDeleteTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.CryptoDeleteTransactionBody{}

	if this.transferAccountID != nil {
		body.TransferAccountID = this.transferAccountID._ToProtobuf()
	}

	if this.deleteAccountID != nil {
		body.DeleteAccountID = this.deleteAccountID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoDelete{
			CryptoDelete: body,
		},
	}, nil
}

func (this *AccountDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().ApproveAllowances,
	}
}
