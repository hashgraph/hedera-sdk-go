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
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// FileUpdateTransaction
// Modify the metadata and/or contents of a file. If a field is not set in the transaction body, the
// corresponding file attribute will be unchanged. This transaction must be signed by all the keys
// in the top level of a key list (M-of-M) of the file being updated. If the keys themselves are
// being updated, then the transaction must also be signed by all the new keys. If the keys contain
// additional KeyList or ThresholdKey then M-of-M secondary KeyList or ThresholdKey signing
// requirements must be meet
type FileUpdateTransaction struct {
	transaction
	fileID         *FileID
	keys           *KeyList
	expirationTime *time.Time
	contents       []byte
	memo           string
}

// NewFileUpdateTransaction creates a FileUpdateTransaction which modifies the metadata and/or contents of a file.
// If a field is not set in the transaction body, the corresponding file attribute will be unchanged.
// This transaction must be signed by all the keys in the top level of a key list (M-of-M) of the file being updated.
// If the keys themselves are being updated, then the transaction must also be signed by all the new keys. If the keys contain
// additional KeyList or ThresholdKey then M-of-M secondary KeyList or ThresholdKey signing
// requirements must be meet
func NewFileUpdateTransaction() *FileUpdateTransaction {
	this := FileUpdateTransaction{
		transaction: _NewTransaction(),
	}
	this._SetDefaultMaxTransactionFee(NewHbar(5))
	this.e= &this
	return &this
}

func _FileUpdateTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *FileUpdateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileUpdate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileUpdate().GetExpirationTime())

	return &FileUpdateTransaction{
		transaction:    this,
		fileID:         _FileIDFromProtobuf(pb.GetFileUpdate().GetFileID()),
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileUpdate().GetContents(),
		memo:           pb.GetFileUpdate().GetMemo().Value,
	}
}

// SetFileID Sets the FileID to be updated
func (this *FileUpdateTransaction) SetFileID(fileID FileID) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.fileID = &fileID
	return this
}

// GetFileID returns the FileID to be updated
func (this *FileUpdateTransaction) GetFileID() FileID {
	if this.fileID == nil {
		return FileID{}
	}

	return *this.fileID
}

// SetKeys Sets the new list of keys that can modify or delete the file
func (this *FileUpdateTransaction) SetKeys(keys ...Key) *FileUpdateTransaction {
	this._RequireNotFrozen()
	if this.keys == nil {
		this.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	this.keys = keyList

	return this
}

func (this *FileUpdateTransaction) GetKeys() KeyList {
	if this.keys != nil {
		return *this.keys
	}

	return KeyList{}
}

// SetExpirationTime Sets the new expiry time
func (this *FileUpdateTransaction) SetExpirationTime(expiration time.Time) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.expirationTime = &expiration
	return this
}

// GetExpirationTime returns the new expiry time
func (this *FileUpdateTransaction) GetExpirationTime() time.Time {
	if this.expirationTime != nil {
		return *this.expirationTime
	}

	return time.Time{}
}

// SetContents Sets the new contents that should overwrite the file's current contents
func (this *FileUpdateTransaction) SetContents(contents []byte) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.contents = contents
	return this
}

// GetContents returns the new contents that should overwrite the file's current contents
func (this *FileUpdateTransaction) GetContents() []byte {
	return this.contents
}

// SetFileMemo Sets the new memo to be associated with the file (UTF-8 encoding max 100 bytes)
func (this *FileUpdateTransaction) SetFileMemo(memo string) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.memo = memo

	return this
}

// GeFileMemo
// Deprecated: use GetFileMemo()
func (this *FileUpdateTransaction) GeFileMemo() string {
	return this.memo
}

func (this *FileUpdateTransaction) GetFileMemo() string {
	return this.memo
}

// ----- Required Interfaces ------- //

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *FileUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *FileUpdateTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// Sign uses the provided privateKey to sign the transaction.
func (this *FileUpdateTransaction) Sign(
	privateKey PrivateKey,
) *FileUpdateTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *FileUpdateTransaction) SignWithOperator(
	client *Client,
) (*FileUpdateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_,err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *FileUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileUpdateTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *FileUpdateTransaction) Freeze() (*FileUpdateTransaction, error) {
	_,err := this.transaction.Freeze()
	return this, err
}

func (this *FileUpdateTransaction) FreezeWith(client *Client) (*FileUpdateTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FileUpdateTransaction) SetMaxTransactionFee(fee Hbar) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *FileUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// SetTransactionMemo sets the memo for this FileUpdateTransaction.
func (this *FileUpdateTransaction) SetTransactionMemo(memo string) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// SetTransactionValidDuration sets the valid duration for this FileUpdateTransaction.
func (this *FileUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	FileUpdateTransaction.
func (this *FileUpdateTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileUpdateTransaction.
func (this *FileUpdateTransaction) SetTransactionID(transactionID TransactionID) *FileUpdateTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountID sets the _Node AccountID for this FileUpdateTransaction.
func (this *FileUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *FileUpdateTransaction) SetMaxRetry(count int) *FileUpdateTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *FileUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileUpdateTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *FileUpdateTransaction) SetMaxBackoff(max time.Duration) *FileUpdateTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *FileUpdateTransaction) SetMinBackoff(min time.Duration) *FileUpdateTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *FileUpdateTransaction) SetLogLevel(level LogLevel) *FileUpdateTransaction {
	this.transaction.SetLogLevel(level)
	return this
}
// ----------- overriden functions ----------------

func (this *FileUpdateTransaction) getName() string {
	return "FileUpdateTransaction"
}
func (this *FileUpdateTransaction) validateNetworkOnIDs(client *Client) error {
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

func (this *FileUpdateTransaction) build() *services.TransactionBody {
	body := &services.FileUpdateTransactionBody{
		Memo: &wrapperspb.StringValue{Value: this.memo},
	}
	if this.fileID != nil {
		body.FileID = this.fileID._ToProtobuf()
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
		Data: &services.TransactionBody_FileUpdate{
			FileUpdate: body,
		},
	}
}

func (this *FileUpdateTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.FileUpdateTransactionBody{
		Memo: &wrapperspb.StringValue{Value: this.memo},
	}
	if this.fileID != nil {
		body.FileID = this.fileID._ToProtobuf()
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
		Data: &services.SchedulableTransactionBody_FileUpdate{
			FileUpdate: body,
		},
	}, nil
}

func (this *FileUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().UpdateFile,
	}
}
