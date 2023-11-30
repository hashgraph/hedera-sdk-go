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

// FileCreateTransaction creates a new file, containing the given contents.  It is referenced by its FileID, and does
// not have a filename, so it is important to get and hold onto the FileID. After the file is created, the FileID for
// it can be found in the receipt, or retrieved with a GetByKey query, or by asking for a Record of the transaction to
// be created, and retrieving that.
//
// See FileInfoQuery for more information about files.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0, with
// a null key. Future versions of the API will support multiple realms and multiple shards.
type FileCreateTransaction struct {
	transaction
	keys           *KeyList
	expirationTime *time.Time
	contents       []byte
	memo           string
}

// NewFileCreateTransaction creates a FileCreateTransaction which creates a new file, containing the given contents.  It is referenced by its FileID, and does
// not have a filename, so it is important to get and hold onto the FileID. After the file is created, the FileID for
// it can be found in the receipt, or retrieved with a GetByKey query, or by asking for a Record of the transaction to
// be created, and retrieving that.
//
// See FileInfoQuery for more information about files.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0, with
// a null key. Future versions of the API will support multiple realms and multiple shards.
func NewFileCreateTransaction() *FileCreateTransaction {
	this := FileCreateTransaction{
		transaction: _NewTransaction(),
	}

	this.SetExpirationTime(time.Now().Add(7890000 * time.Second))
	this._SetDefaultMaxTransactionFee(NewHbar(5))
	this.e = &this

	return &this
}

func _FileCreateTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *FileCreateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileCreate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileCreate().GetExpirationTime())

	resultTx := &FileCreateTransaction{
		transaction:    this,
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileCreate().GetContents(),
		memo:           pb.GetMemo(),
	}
	resultTx.e = resultTx
	return resultTx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *FileCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *FileCreateTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// AddKey adds a key to the internal list of keys associated with the file. All of the keys on the list must sign to
// create or modify a file, but only one of them needs to sign in order to delete the file. Each of those "keys" may
// itself be threshold key containing other keys (including other threshold keys). In other words, the behavior is an
// AND for create/modify, OR for delete. This is useful for acting as a revocation server. If it is desired to have the
// behavior be AND for all 3 operations (or OR for all 3), then the list should have only a single Key, which is a
// threshold key, with N=1 for OR, N=M for AND.
//
// If a file is created without adding ANY keys, the file is immutable and ONLY the
// expirationTime of the file can be changed using FileUpdateTransaction. The file contents or its keys will not be
// mutable.
func (this *FileCreateTransaction) SetKeys(keys ...Key) *FileCreateTransaction {
	this._RequireNotFrozen()
	if this.keys == nil {
		this.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	this.keys = keyList

	return this
}

func (this *FileCreateTransaction) GetKeys() KeyList {
	if this.keys != nil {
		return *this.keys
	}

	return KeyList{}
}

// SetExpirationTime sets the time at which this file should expire (unless FileUpdateTransaction is used before then to
// extend its life). The file will automatically disappear at the fileExpirationTime, unless its expiration is extended
// by another transaction before that time. If the file is deleted, then its contents will become empty and it will be
// marked as deleted until it expires, and then it will cease to exist.
func (this *FileCreateTransaction) SetExpirationTime(expiration time.Time) *FileCreateTransaction {
	this._RequireNotFrozen()
	this.expirationTime = &expiration
	return this
}

func (this *FileCreateTransaction) GetExpirationTime() time.Time {
	if this.expirationTime != nil {
		return *this.expirationTime
	}

	return time.Time{}
}

// SetContents sets the bytes that are the contents of the file (which can be empty). If the size of the file and other
// fields in the transaction exceed the max transaction size then FileCreateTransaction can be used to continue
// uploading the file.
func (this *FileCreateTransaction) SetContents(contents []byte) *FileCreateTransaction {
	this._RequireNotFrozen()
	this.contents = contents
	return this
}

// GetContents returns the bytes that are the contents of the file (which can be empty).
func (this *FileCreateTransaction) GetContents() []byte {
	return this.contents
}

// SetMemo Sets the memo associated with the file (UTF-8 encoding max 100 bytes)
func (this *FileCreateTransaction) SetMemo(memo string) *FileCreateTransaction {
	this._RequireNotFrozen()
	this.memo = memo
	return this
}

// GetMemo returns the memo associated with the file (UTF-8 encoding max 100 bytes)
func (this *FileCreateTransaction) GetMemo() string {
	return this.memo
}

func (this *FileCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

// Sign uses the provided privateKey to sign the transaction.
func (this *FileCreateTransaction) Sign(
	privateKey PrivateKey,
) *FileCreateTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *FileCreateTransaction) SignWithOperator(
	client *Client,
) (*FileCreateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *FileCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileCreateTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *FileCreateTransaction) Freeze() (*FileCreateTransaction, error) {
	_, err := this.transaction.Freeze()
	return this, err
}

func (this *FileCreateTransaction) FreezeWith(client *Client) (*FileCreateTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FileCreateTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FileCreateTransaction) SetMaxTransactionFee(fee Hbar) *FileCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *FileCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *FileCreateTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this FileCreateTransaction.
func (this *FileCreateTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileCreateTransaction.
func (this *FileCreateTransaction) SetTransactionMemo(memo string) *FileCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *FileCreateTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileCreateTransaction.
func (this *FileCreateTransaction) SetTransactionValidDuration(duration time.Duration) *FileCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	FileCreateTransaction.
func (this *FileCreateTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileCreateTransaction.
func (this *FileCreateTransaction) SetTransactionID(transactionID TransactionID) *FileCreateTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountID sets the _Node AccountID for this FileCreateTransaction.
func (this *FileCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *FileCreateTransaction) SetMaxRetry(count int) *FileCreateTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *FileCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileCreateTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *FileCreateTransaction) SetMaxBackoff(max time.Duration) *FileCreateTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *FileCreateTransaction) SetMinBackoff(min time.Duration) *FileCreateTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *FileCreateTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("FileCreateTransaction:%d", timestamp.UnixNano())
}

func (this *FileCreateTransaction) SetLogLevel(level LogLevel) *FileCreateTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *FileCreateTransaction) getName() string {
	return "FileCreateTransaction"
}
func (this *FileCreateTransaction) build() *services.TransactionBody {
	body := &services.FileCreateTransactionBody{
		Memo: this.memo,
	}

	if this.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*this.expirationTime)
	}

	if this.keys != nil {
		body.Keys = this.keys._ToProtoKeyList()
	}

	if this.contents != nil {
		body.Contents = this.contents
	}

	return &services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileCreate{
			FileCreate: body,
		},
	}
}

func (this *FileCreateTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.FileCreateTransactionBody{
		Memo: this.memo,
	}

	if this.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*this.expirationTime)
	}

	if this.keys != nil {
		body.Keys = this.keys._ToProtoKeyList()
	}

	if this.contents != nil {
		body.Contents = this.contents
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_FileCreate{
			FileCreate: body,
		},
	}, nil
}

func (this *FileCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().CreateFile,
	}
}
