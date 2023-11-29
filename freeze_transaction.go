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

type FreezeTransaction struct {
	transaction
	startTime  time.Time
	endTime    time.Time
	fileID     *FileID
	fileHash   []byte
	freezeType FreezeType
}

func NewFreezeTransaction() *FreezeTransaction {
	this := FreezeTransaction{
		transaction: _NewTransaction(),
	}

	this._SetDefaultMaxTransactionFee(NewHbar(2))
	this.e= &this

	return &this
}

func _FreezeTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *FreezeTransaction {
	startTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(pb.GetFreeze().GetStartHour()), int(pb.GetFreeze().GetStartMin()), // nolint
		0, time.Now().Nanosecond(), time.Now().Location(),
	)

	endTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(pb.GetFreeze().GetEndHour()), int(pb.GetFreeze().GetEndMin()), // nolint
		0, time.Now().Nanosecond(), time.Now().Location(),
	)

	return &FreezeTransaction{
		transaction: this,
		startTime:   startTime,
		endTime:     endTime,
		fileID:      _FileIDFromProtobuf(pb.GetFreeze().GetUpdateFile()),
		fileHash:    pb.GetFreeze().FileHash,
	}
}

func (this *FreezeTransaction) SetStartTime(startTime time.Time) *FreezeTransaction {
	this._RequireNotFrozen()
	this.startTime = startTime
	return this
}

func (this *FreezeTransaction) GetStartTime() time.Time {
	return this.startTime
}

// Deprecated
func (this *FreezeTransaction) SetEndTime(endTime time.Time) *FreezeTransaction {
	this._RequireNotFrozen()
	this.endTime = endTime
	return this
}

// Deprecated
func (this *FreezeTransaction) GetEndTime() time.Time {
	return this.endTime
}

func (this *FreezeTransaction) SetFileID(id FileID) *FreezeTransaction {
	this._RequireNotFrozen()
	this.fileID = &id
	return this
}

func (this *FreezeTransaction) GetFileID() *FileID {
	return this.fileID
}

func (this *FreezeTransaction) SetFreezeType(freezeType FreezeType) *FreezeTransaction {
	this._RequireNotFrozen()
	this.freezeType = freezeType
	return this
}

func (this *FreezeTransaction) GetFreezeType() FreezeType {
	return this.freezeType
}

func (this *FreezeTransaction) SetFileHash(hash []byte) *FreezeTransaction {
	this._RequireNotFrozen()
	this.fileHash = hash
	return this
}

func (this *FreezeTransaction) GetFileHash() []byte {
	return this.fileHash
}

func (this *FreezeTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *FreezeTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (this *FreezeTransaction) Sign(
	privateKey PrivateKey,
) *FreezeTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *FreezeTransaction) SignWithOperator(
	client *Client,
) (*FreezeTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_,err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *FreezeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FreezeTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *FreezeTransaction) Freeze() (*FreezeTransaction, error) {
	_,err := this.transaction.Freeze()
	return this, err
}

func (this *FreezeTransaction) FreezeWith(client *Client) (*FreezeTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FreezeTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FreezeTransaction) SetMaxTransactionFee(fee Hbar) *FreezeTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *FreezeTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FreezeTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *FreezeTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this FreezeTransaction.
func (this *FreezeTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FreezeTransaction.
func (this *FreezeTransaction) SetTransactionMemo(memo string) *FreezeTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *FreezeTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FreezeTransaction.
func (this *FreezeTransaction) SetTransactionValidDuration(duration time.Duration) *FreezeTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	FreezeTransaction.
func (this *FreezeTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FreezeTransaction.
func (this *FreezeTransaction) SetTransactionID(transactionID TransactionID) *FreezeTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountID sets the _Node AccountID for this FreezeTransaction.
func (this *FreezeTransaction) SetNodeAccountIDs(nodeID []AccountID) *FreezeTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *FreezeTransaction) SetMaxRetry(count int) *FreezeTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *FreezeTransaction) AddSignature(publicKey PublicKey, signature []byte) *FreezeTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *FreezeTransaction) SetMaxBackoff(max time.Duration) *FreezeTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *FreezeTransaction) SetMinBackoff(min time.Duration) *FreezeTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *FreezeTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("FreezeTransaction:%d", timestamp.UnixNano())
}

func (this *FreezeTransaction) SetLogLevel(level LogLevel) *FreezeTransaction {
	this.transaction.SetLogLevel(level)
	return this
}
// ----------- overriden functions ----------------

func (this *FreezeTransaction) getName() string {
	return "FreezeTransaction"
}
func (this *FreezeTransaction) build() *services.TransactionBody {
	body := &services.FreezeTransactionBody{
		FileHash:   this.fileHash,
		StartTime:  _TimeToProtobuf(this.startTime),
		FreezeType: services.FreezeType(this.freezeType),
	}

	if this.fileID != nil {
		body.UpdateFile = this.fileID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_Freeze{
			Freeze: body,
		},
	}
}
func (this *FreezeTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.FreezeTransactionBody{
		FileHash:   this.fileHash,
		StartTime:  _TimeToProtobuf(this.startTime),
		FreezeType: services.FreezeType(this.freezeType),
	}

	if this.fileID != nil {
		body.UpdateFile = this.fileID._ToProtobuf()
	}
	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_Freeze{
			Freeze: body,
		},
	}, nil
}

func (this *FreezeTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFreeze().Freeze,
	}
}
