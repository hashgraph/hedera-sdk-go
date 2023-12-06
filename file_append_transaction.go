package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use tx file except in compliance with the License.
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

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// FileAppendTransaction appends the given contents to the end of the file. If a file is too big to create with a single
// FileCreateTransaction, then it can be created with the first part of its contents, and then appended multiple times
// to create the entire file.
type FileAppendTransaction struct {
	Transaction
	maxChunks uint64
	contents  []byte
	fileID    *FileID
	chunkSize int
}

// NewFileAppendTransaction creates a FileAppendTransaction transaction which can be
// used to construct and execute a File Append transaction.
func NewFileAppendTransaction() *FileAppendTransaction {
	tx := FileAppendTransaction{
		Transaction: _NewTransaction(),
		maxChunks:   20,
		contents:    make([]byte, 0),
		chunkSize:   2048,
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(5))
	tx.e = &tx

	return &tx
}

func _FileAppendTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *FileAppendTransaction {
	resultTx := &FileAppendTransaction{
		Transaction: tx,
		maxChunks:   20,
		contents:    pb.GetFileAppend().GetContents(),
		chunkSize:   2048,
		fileID:      _FileIDFromProtobuf(pb.GetFileAppend().GetFileID()),
	}
	resultTx.e = resultTx
	return resultTx
}

// SetFileID sets the FileID of the file to which the bytes are appended to.
func (tx *FileAppendTransaction) SetFileID(fileID FileID) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID of the file to which the bytes are appended to.
func (tx *FileAppendTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// SetMaxChunkSize Sets maximum amount of chunks append function can create
func (tx *FileAppendTransaction) SetMaxChunkSize(size int) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.chunkSize = size
	return tx
}

// GetMaxChunkSize returns maximum amount of chunks append function can create
func (tx *FileAppendTransaction) GetMaxChunkSize() int {
	return tx.chunkSize
}

// SetMaxChunks sets the maximum number of chunks that can be created
func (tx *FileAppendTransaction) SetMaxChunks(size uint64) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.maxChunks = size
	return tx
}

// GetMaxChunks returns the maximum number of chunks that can be created
func (tx *FileAppendTransaction) GetMaxChunks() uint64 {
	return tx.maxChunks
}

// SetContents sets the bytes to append to the contents of the file.
func (tx *FileAppendTransaction) SetContents(contents []byte) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.contents = contents
	return tx
}

// GetContents returns the bytes to append to the contents of the file.
func (tx *FileAppendTransaction) GetContents() []byte {
	return tx.contents
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *FileAppendTransaction) Sign(
	privateKey PrivateKey,
) *FileAppendTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *FileAppendTransaction) SignWithOperator(
	client *Client,
) (*FileAppendTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.Transaction.SignWithOperator(client)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *FileAppendTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileAppendTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *FileAppendTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileAppendTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when tx deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *FileAppendTransaction) SetGrpcDeadline(deadline *time.Duration) *FileAppendTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *FileAppendTransaction) Freeze() (*FileAppendTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *FileAppendTransaction) FreezeWith(client *Client) (*FileAppendTransaction, error) {
	if tx.IsFrozen() {
		return tx, nil
	}

	if tx.nodeAccountIDs._Length() == 0 {
		if client == nil {
			return tx, errNoClientOrTransactionIDOrNodeId
		}

		tx.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	tx._InitFee(client)
	err := tx.validateNetworkOnIDs(client)
	if err != nil {
		return &FileAppendTransaction{}, err
	}
	if err := tx._InitTransactionID(client); err != nil {
		return tx, err
	}
	body := tx.build()

	chunks := uint64((len(tx.contents) + (tx.chunkSize - 1)) / tx.chunkSize)
	if chunks > tx.maxChunks {
		return tx, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: tx.maxChunks,
		}
	}

	nextTransactionID := tx.transactionIDs._GetCurrent().(TransactionID)

	tx.transactionIDs = _NewLockableSlice()
	tx.transactions = _NewLockableSlice()
	tx.signedTransactions = _NewLockableSlice()

	if b, ok := body.Data.(*services.TransactionBody_FileAppend); ok {
		for i := 0; uint64(i) < chunks; i++ {
			start := i * tx.chunkSize
			end := start + tx.chunkSize

			if end > len(tx.contents) {
				end = len(tx.contents)
			}

			tx.transactionIDs._Push(_TransactionIDFromProtobuf(nextTransactionID._ToProtobuf()))
			if err != nil {
				panic(err)
			}
			b.FileAppend.Contents = tx.contents[start:end]

			body.TransactionID = nextTransactionID._ToProtobuf()
			body.Data = &services.TransactionBody_FileAppend{
				FileAppend: b.FileAppend,
			}

			for _, nodeAccountID := range tx.GetNodeAccountIDs() {
				body.NodeAccountID = nodeAccountID._ToProtobuf()

				bodyBytes, err := protobuf.Marshal(body)
				if err != nil {
					return tx, errors.Wrap(err, "error serializing body for file append")
				}

				tx.signedTransactions._Push(&services.SignedTransaction{
					BodyBytes: bodyBytes,
					SigMap:    &services.SignatureMap{},
				})
			}

			validStart := *nextTransactionID.ValidStart

			*nextTransactionID.ValidStart = validStart.Add(1 * time.Nanosecond)
		}
	}

	return tx, nil
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *FileAppendTransaction) SetMaxTransactionFee(fee Hbar) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *FileAppendTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx FileAppendTransaction.
func (tx *FileAppendTransaction) SetTransactionMemo(memo string) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for tx FileAppendTransaction.
func (tx *FileAppendTransaction) SetTransactionValidDuration(duration time.Duration) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for tx FileAppendTransaction.
func (tx *FileAppendTransaction) SetTransactionID(transactionID TransactionID) *FileAppendTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountID sets the _Node AccountID for tx FileAppendTransaction.
func (tx *FileAppendTransaction) SetNodeAccountIDs(nodeAccountIDs []AccountID) *FileAppendTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeAccountIDs)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *FileAppendTransaction) SetMaxRetry(count int) *FileAppendTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches tx time.
func (tx *FileAppendTransaction) SetMaxBackoff(max time.Duration) *FileAppendTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *FileAppendTransaction) SetMinBackoff(min time.Duration) *FileAppendTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *FileAppendTransaction) SetLogLevel(level LogLevel) *FileAppendTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *FileAppendTransaction) getName() string {
	return "FileAppendTransaction"
}
func (tx *FileAppendTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.fileID != nil {
		if err := tx.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *FileAppendTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileAppend{
			FileAppend: tx.buildProtoBody(),
		},
	}
}

func (tx *FileAppendTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileAppend{
			FileAppend: tx.buildProtoBody(),
		},
	}, nil
}
func (tx *FileAppendTransaction) buildProtoBody() *services.FileAppendTransactionBody {
	body := &services.FileAppendTransactionBody{
		Contents: tx.contents,
	}

	if tx.fileID != nil {
		body.FileID = tx.fileID._ToProtobuf()
	}

	return body
}

func (tx *FileAppendTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().AppendContent,
	}
}

func (tx *FileAppendTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
