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

// FileDeleteTransaction Deletes the given file. After deletion, it will be marked as deleted and will have no contents.
// But information about it will continue to exist until it expires. A list of keys was given when
// the file was created. All the top level keys on that list must sign transactions to create or
// modify the file, but any single one of the top level keys can be used to delete the file. This
// transaction must be signed by 1-of-M KeyList keys. If keys contains additional KeyList or
// ThresholdKey then 1-of-M secondary KeyList or ThresholdKey signing requirements must be meet.
type FileDeleteTransaction struct {
	transaction
	fileID *FileID
}

// NewFileDeleteTransaction creates a FileDeleteTransaction which deletes the given file. After deletion,
// it will be marked as deleted and will have no contents.
// But information about it will continue to exist until it expires. A list of keys was given when
// the file was created. All the top level keys on that list must sign transactions to create or
// modify the file, but any single one of the top level keys can be used to delete the file. This
// transaction must be signed by 1-of-M KeyList keys. If keys contains additional KeyList or
// ThresholdKey then 1-of-M secondary KeyList or ThresholdKey signing requirements must be meet.
func NewFileDeleteTransaction() *FileDeleteTransaction {
	this := FileDeleteTransaction{
		transaction: _NewTransaction(),
	}
	this._SetDefaultMaxTransactionFee(NewHbar(5))
	this.e = &this

	return &this
}

func _FileDeleteTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *FileDeleteTransaction {
	resultTx := &FileDeleteTransaction{
		transaction: this,
		fileID:      _FileIDFromProtobuf(pb.GetFileDelete().GetFileID()),
	}
	resultTx.e = resultTx
	return resultTx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *FileDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *FileDeleteTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetFileID Sets the FileID of the file to be deleted
func (this *FileDeleteTransaction) SetFileID(fileID FileID) *FileDeleteTransaction {
	this._RequireNotFrozen()
	this.fileID = &fileID
	return this
}

// GetFileID returns the FileID of the file to be deleted
func (this *FileDeleteTransaction) GetFileID() FileID {
	if this.fileID == nil {
		return FileID{}
	}

	return *this.fileID
}

func (this *FileDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

// Sign uses the provided privateKey to sign the transaction.
func (this *FileDeleteTransaction) Sign(
	privateKey PrivateKey,
) *FileDeleteTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *FileDeleteTransaction) SignWithOperator(
	client *Client,
) (*FileDeleteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *FileDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileDeleteTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *FileDeleteTransaction) Freeze() (*FileDeleteTransaction, error) {
	_, err := this.transaction.Freeze()
	return this, err
}

func (this *FileDeleteTransaction) FreezeWith(client *Client) (*FileDeleteTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FileDeleteTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FileDeleteTransaction) SetMaxTransactionFee(fee Hbar) *FileDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *FileDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *FileDeleteTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this FileDeleteTransaction.
func (this *FileDeleteTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileDeleteTransaction.
func (this *FileDeleteTransaction) SetTransactionMemo(memo string) *FileDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *FileDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileDeleteTransaction.
func (this *FileDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *FileDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	FileDeleteTransaction.
func (this *FileDeleteTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileDeleteTransaction.
func (this *FileDeleteTransaction) SetTransactionID(transactionID TransactionID) *FileDeleteTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountID sets the _Node AccountID for this FileDeleteTransaction.
func (this *FileDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *FileDeleteTransaction) SetMaxRetry(count int) *FileDeleteTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *FileDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileDeleteTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *FileDeleteTransaction) SetMaxBackoff(max time.Duration) *FileDeleteTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *FileDeleteTransaction) SetMinBackoff(min time.Duration) *FileDeleteTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *FileDeleteTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("FileDeleteTransaction:%d", timestamp.UnixNano())
}

func (this *FileDeleteTransaction) SetLogLevel(level LogLevel) *FileDeleteTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *FileDeleteTransaction) getName() string {
	return "FileDeleteTransaction"
}
func (this *FileDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.fileID != nil {
		if err := this.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *FileDeleteTransaction) build() *services.TransactionBody {
	body := &services.FileDeleteTransactionBody{}
	if this.fileID != nil {
		body.FileID = this.fileID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileDelete{
			FileDelete: body,
		},
	}
}

func (this *FileDeleteTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.FileDeleteTransactionBody{}
	if this.fileID != nil {
		body.FileID = this.fileID._ToProtobuf()
	}
	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_FileDelete{
			FileDelete: body,
		},
	}, nil
}

func (this *FileDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().DeleteFile,
	}
}
