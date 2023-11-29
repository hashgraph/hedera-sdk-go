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
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// FileAppendTransaction appends the given contents to the end of the file. If a file is too big to create with a single
// FileCreateTransaction, then it can be created with the first part of its contents, and then appended multiple times
// to create the entire file.
type FileAppendTransaction struct {
	transaction
	maxChunks uint64
	contents  []byte
	fileID    *FileID
	chunkSize int
}

// NewFileAppendTransaction creates a FileAppendTransaction transaction which can be
// used to construct and execute a File Append transaction.
func NewFileAppendTransaction() *FileAppendTransaction {
	this := FileAppendTransaction{
		transaction: _NewTransaction(),
		maxChunks:   20,
		contents:    make([]byte, 0),
		chunkSize:   2048,
	}
	this._SetDefaultMaxTransactionFee(NewHbar(5))
	this.e=&this

	return &this
}

func _FileAppendTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *FileAppendTransaction {
	return &FileAppendTransaction{
		transaction: this,
		maxChunks:   20,
		contents:    pb.GetFileAppend().GetContents(),
		chunkSize:   2048,
		fileID:      _FileIDFromProtobuf(pb.GetFileAppend().GetFileID()),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *FileAppendTransaction) SetGrpcDeadline(deadline *time.Duration) *FileAppendTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetFileID sets the FileID of the file to which the bytes are appended to.
func (this *FileAppendTransaction) SetFileID(fileID FileID) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.fileID = &fileID
	return this
}

// GetFileID returns the FileID of the file to which the bytes are appended to.
func (this *FileAppendTransaction) GetFileID() FileID {
	if this.fileID == nil {
		return FileID{}
	}

	return *this.fileID
}

// SetMaxChunkSize Sets maximum amount of chunks append function can create
func (this *FileAppendTransaction) SetMaxChunkSize(size int) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.chunkSize = size
	return this
}

// GetMaxChunkSize returns maximum amount of chunks append function can create
func (this *FileAppendTransaction) GetMaxChunkSize() int {
	return this.chunkSize
}

// SetMaxChunks sets the maximum number of chunks that can be created
func (this *FileAppendTransaction) SetMaxChunks(size uint64) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.maxChunks = size
	return this
}

// GetMaxChunks returns the maximum number of chunks that can be created
func (this *FileAppendTransaction) GetMaxChunks() uint64 {
	return this.maxChunks
}

// SetContents sets the bytes to append to the contents of the file.
func (this *FileAppendTransaction) SetContents(contents []byte) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.contents = contents
	return this
}

// GetContents returns the bytes to append to the contents of the file.
func (this *FileAppendTransaction) GetContents() []byte {
	return this.contents
}


func (this *FileAppendTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	chunks := uint64((len(this.contents) + (this.chunkSize - 1)) / this.chunkSize)
	if chunks > 1 {
		return &ScheduleCreateTransaction{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: 1,
		}
	}

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (this *FileAppendTransaction) IsFrozen() bool {
	return this._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (this *FileAppendTransaction) Sign(
	privateKey PrivateKey,
) *FileAppendTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *FileAppendTransaction) SignWithOperator(
	client *Client,
) (*FileAppendTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_,err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *FileAppendTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileAppendTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *FileAppendTransaction) Freeze() (*FileAppendTransaction, error) {
	return this.FreezeWith(nil)
}

func (this *FileAppendTransaction) FreezeWith(client *Client) (*FileAppendTransaction, error) {
	if this.IsFrozen() {
		return this, nil
	}

	if this.nodeAccountIDs._Length() == 0 {
		if client == nil {
			return this, errNoClientOrTransactionIDOrNodeId
		}

		this.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	this._InitFee(client)
	err := this.validateNetworkOnIDs(client)
	if err != nil {
		return &FileAppendTransaction{}, err
	}
	if err := this._InitTransactionID(client); err != nil {
		return this, err
	}
	body := this.build()

	chunks := uint64((len(this.contents) + (this.chunkSize - 1)) / this.chunkSize)
	if chunks > this.maxChunks {
		return this, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: this.maxChunks,
		}
	}

	nextTransactionID := this.transactionIDs._GetCurrent().(TransactionID)

	this.transactionIDs = _NewLockableSlice()
	this.transactions = _NewLockableSlice()
	this.signedTransactions = _NewLockableSlice()

	if b, ok := body.Data.(*services.TransactionBody_FileAppend); ok {
		for i := 0; uint64(i) < chunks; i++ {
			start := i * this.chunkSize
			end := start + this.chunkSize

			if end > len(this.contents) {
				end = len(this.contents)
			}

			this.transactionIDs._Push(_TransactionIDFromProtobuf(nextTransactionID._ToProtobuf()))
			if err != nil {
				panic(err)
			}
			b.FileAppend.Contents = this.contents[start:end]

			body.TransactionID = nextTransactionID._ToProtobuf()
			body.Data = &services.TransactionBody_FileAppend{
				FileAppend: b.FileAppend,
			}

			for _, nodeAccountID := range this.GetNodeAccountIDs() {
				body.NodeAccountID = nodeAccountID._ToProtobuf()

				bodyBytes, err := protobuf.Marshal(body)
				if err != nil {
					return this, errors.Wrap(err, "error serializing body for file append")
				}

				this.signedTransactions._Push(&services.SignedTransaction{
					BodyBytes: bodyBytes,
					SigMap:    &services.SignatureMap{},
				})
			}

			validStart := *nextTransactionID.ValidStart

			*nextTransactionID.ValidStart = validStart.Add(1 * time.Nanosecond)
		}
	}

	return this, nil
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FileAppendTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *FileAppendTransaction) SetMaxTransactionFee(fee Hbar) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *FileAppendTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *FileAppendTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this FileAppendTransaction.
func (this *FileAppendTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileAppendTransaction.
func (this *FileAppendTransaction) SetTransactionMemo(memo string) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *FileAppendTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileAppendTransaction.
func (this *FileAppendTransaction) SetTransactionValidDuration(duration time.Duration) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	FileAppendTransaction.
func (this *FileAppendTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileAppendTransaction.
func (this *FileAppendTransaction) SetTransactionID(transactionID TransactionID) *FileAppendTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountID sets the _Node AccountID for this FileAppendTransaction.
func (this *FileAppendTransaction) SetNodeAccountIDs(nodeAccountIDs []AccountID) *FileAppendTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeAccountIDs)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *FileAppendTransaction) SetMaxRetry(count int) *FileAppendTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *FileAppendTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileAppendTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *FileAppendTransaction) SetMaxBackoff(max time.Duration) *FileAppendTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *FileAppendTransaction) SetMinBackoff(min time.Duration) *FileAppendTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *FileAppendTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("FileAppendTransaction:%d", timestamp.UnixNano())
}

func (this *FileAppendTransaction) SetLogLevel(level LogLevel) *FileAppendTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *FileAppendTransaction) getName() string {
	return "FileAppendTransaction"
}
func (this *FileAppendTransaction) validateNetworkOnIDs(client *Client) error {
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

func (this *FileAppendTransaction) build() *services.TransactionBody {
	body := &services.FileAppendTransactionBody{}
	if this.fileID != nil {
		body.FileID = this.fileID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileAppend{
			FileAppend: body,
		},
	}
}

func (this *FileAppendTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.FileAppendTransactionBody{
		Contents: this.contents,
	}

	if this.fileID != nil {
		body.FileID = this.fileID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_FileAppend{
			FileAppend: body,
		},
	}, nil
}

func (this *FileAppendTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().AppendContent,
	}
}