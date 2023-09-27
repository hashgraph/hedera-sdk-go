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
	Transaction
	maxChunks uint64
	contents  []byte
	fileID    *FileID
	chunkSize int
}

// NewFileAppendTransaction creates a FileAppendTransaction transaction which can be
// used to construct and execute a File Append Transaction.
func NewFileAppendTransaction() *FileAppendTransaction {
	transaction := FileAppendTransaction{
		Transaction: _NewTransaction(),
		maxChunks:   20,
		contents:    make([]byte, 0),
		chunkSize:   2048,
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _FileAppendTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *FileAppendTransaction {
	return &FileAppendTransaction{
		Transaction: transaction,
		maxChunks:   20,
		contents:    pb.GetFileAppend().GetContents(),
		chunkSize:   2048,
		fileID:      _FileIDFromProtobuf(pb.GetFileAppend().GetFileID()),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *FileAppendTransaction) SetGrpcDeadline(deadline *time.Duration) *FileAppendTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetFileID sets the FileID of the file to which the bytes are appended to.
func (transaction *FileAppendTransaction) SetFileID(fileID FileID) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.fileID = &fileID
	return transaction
}

// GetFileID returns the FileID of the file to which the bytes are appended to.
func (transaction *FileAppendTransaction) GetFileID() FileID {
	if transaction.fileID == nil {
		return FileID{}
	}

	return *transaction.fileID
}

// SetMaxChunkSize Sets maximum amount of chunks append function can create
func (transaction *FileAppendTransaction) SetMaxChunkSize(size int) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.chunkSize = size
	return transaction
}

// GetMaxChunkSize returns maximum amount of chunks append function can create
func (transaction *FileAppendTransaction) GetMaxChunkSize() int {
	return transaction.chunkSize
}

// SetMaxChunks sets the maximum number of chunks that can be created
func (transaction *FileAppendTransaction) SetMaxChunks(size uint64) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.maxChunks = size
	return transaction
}

// GetMaxChunks returns the maximum number of chunks that can be created
func (transaction *FileAppendTransaction) GetMaxChunks() uint64 {
	return transaction.maxChunks
}

// SetContents sets the bytes to append to the contents of the file.
func (transaction *FileAppendTransaction) SetContents(contents []byte) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.contents = contents
	return transaction
}

// GetContents returns the bytes to append to the contents of the file.
func (transaction *FileAppendTransaction) GetContents() []byte {
	return transaction.contents
}

func (transaction *FileAppendTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.fileID != nil {
		if err := transaction.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *FileAppendTransaction) _Build() *services.TransactionBody {
	body := &services.FileAppendTransactionBody{}
	if transaction.fileID != nil {
		body.FileID = transaction.fileID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileAppend{
			FileAppend: body,
		},
	}
}

func (transaction *FileAppendTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	chunks := uint64((len(transaction.contents) + (transaction.chunkSize - 1)) / transaction.chunkSize)
	if chunks > 1 {
		return &ScheduleCreateTransaction{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: 1,
		}
	}

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *FileAppendTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.FileAppendTransactionBody{
		Contents: transaction.contents,
	}

	if transaction.fileID != nil {
		body.FileID = transaction.fileID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileAppend{
			FileAppend: body,
		},
	}, nil
}

func _FileAppendTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().AppendContent,
	}
}

func (transaction *FileAppendTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *FileAppendTransaction) Sign(
	privateKey PrivateKey,
) *FileAppendTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *FileAppendTransaction) SignWithOperator(
	client *Client,
) (*FileAppendTransaction, error) {
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
func (transaction *FileAppendTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileAppendTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *FileAppendTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	list, err := transaction.ExecuteAll(client)

	if err != nil {
		if len(list) > 0 {
			return TransactionResponse{
				TransactionID: transaction.GetTransactionID(),
				NodeID:        list[0].NodeID,
				Hash:          make([]byte, 0),
			}, err
		}
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			Hash:          make([]byte, 0),
		}, err
	}

	return list[0], nil
}

// ExecuteAll executes the all the Transactions with the provided client
func (transaction *FileAppendTransaction) ExecuteAll(
	client *Client,
) ([]TransactionResponse, error) {
	if client == nil || client.operator == nil {
		return []TransactionResponse{}, errNoClientProvided
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return []TransactionResponse{}, err
		}
	}

	var transactionID TransactionID
	if transaction.transactionIDs._Length() > 0 {
		transactionID = transaction.GetTransactionID()
	} else {
		return []TransactionResponse{}, errors.New("transactionID list is empty")
	}

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	size := transaction.signedTransactions._Length() / transaction.nodeAccountIDs._Length()
	list := make([]TransactionResponse, size)

	for i := 0; i < size; i++ {
		resp, err := _Execute(
			client,
			&transaction.Transaction,
			_TransactionShouldRetry,
			_TransactionMakeRequest,
			_TransactionAdvanceRequest,
			_TransactionGetNodeAccountID,
			_FileAppendTransactionGetMethod,
			_TransactionMapStatusError,
			_TransactionMapResponse,
			transaction._GetLogID(),
			transaction.grpcDeadline,
			transaction.maxBackoff,
			transaction.minBackoff,
			transaction.maxRetry,
		)

		if err != nil {
			return list, err
		}

		list[i] = resp.(TransactionResponse)

		_, err = NewTransactionReceiptQuery().
			SetNodeAccountIDs([]AccountID{resp.(TransactionResponse).NodeID}).
			SetTransactionID(resp.(TransactionResponse).TransactionID).
			Execute(client)
		if err != nil {
			return list, err
		}
	}

	return list, nil
}

func (transaction *FileAppendTransaction) Freeze() (*FileAppendTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FileAppendTransaction) FreezeWith(client *Client) (*FileAppendTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}

	if transaction.nodeAccountIDs._Length() == 0 {
		if client == nil {
			return transaction, errNoClientOrTransactionIDOrNodeId
		}

		transaction.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &FileAppendTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	chunks := uint64((len(transaction.contents) + (transaction.chunkSize - 1)) / transaction.chunkSize)
	if chunks > transaction.maxChunks {
		return transaction, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: transaction.maxChunks,
		}
	}

	nextTransactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

	transaction.transactionIDs = _NewLockableSlice()
	transaction.transactions = _NewLockableSlice()
	transaction.signedTransactions = _NewLockableSlice()

	if b, ok := body.Data.(*services.TransactionBody_FileAppend); ok {
		for i := 0; uint64(i) < chunks; i++ {
			start := i * transaction.chunkSize
			end := start + transaction.chunkSize

			if end > len(transaction.contents) {
				end = len(transaction.contents)
			}

			transaction.transactionIDs._Push(_TransactionIDFromProtobuf(nextTransactionID._ToProtobuf()))
			if err != nil {
				panic(err)
			}
			b.FileAppend.Contents = transaction.contents[start:end]

			body.TransactionID = nextTransactionID._ToProtobuf()
			body.Data = &services.TransactionBody_FileAppend{
				FileAppend: b.FileAppend,
			}

			for _, nodeAccountID := range transaction.GetNodeAccountIDs() {
				body.NodeAccountID = nodeAccountID._ToProtobuf()

				bodyBytes, err := protobuf.Marshal(body)
				if err != nil {
					return transaction, errors.Wrap(err, "error serializing body for file append")
				}

				transaction.signedTransactions._Push(&services.SignedTransaction{
					BodyBytes: bodyBytes,
					SigMap:    &services.SignatureMap{},
				})
			}

			validStart := *nextTransactionID.ValidStart

			*nextTransactionID.ValidStart = validStart.Add(1 * time.Nanosecond)
		}
	}

	return transaction, nil
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *FileAppendTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *FileAppendTransaction) SetMaxTransactionFee(fee Hbar) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *FileAppendTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *FileAppendTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this FileAppendTransaction.
func (transaction *FileAppendTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetTransactionMemo(memo string) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *FileAppendTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetTransactionValidDuration(duration time.Duration) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	FileAppendTransaction.
func (transaction *FileAppendTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetTransactionID(transactionID TransactionID) *FileAppendTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetNodeAccountIDs(nodeAccountIDs []AccountID) *FileAppendTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeAccountIDs)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *FileAppendTransaction) SetMaxRetry(count int) *FileAppendTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *FileAppendTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileAppendTransaction {
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
func (transaction *FileAppendTransaction) SetMaxBackoff(max time.Duration) *FileAppendTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *FileAppendTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *FileAppendTransaction) SetMinBackoff(min time.Duration) *FileAppendTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *FileAppendTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *FileAppendTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("FileAppendTransaction:%d", timestamp.UnixNano())
}

func (transaction *FileAppendTransaction) SetLogLevel(level LogLevel) *FileAppendTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
