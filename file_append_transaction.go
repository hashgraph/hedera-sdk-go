package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// FileAppendTransaction appends the given contents to the end of the file. If a file is too big to create with a single
// FileCreateTransaction, then it can be created with the first part of its contents, and then appended multiple times
// to create the entire file.
type FileAppendTransaction struct {
	*Transaction[*FileAppendTransaction]
	maxChunks uint64
	contents  []byte
	fileID    *FileID
	chunkSize int
}

// NewFileAppendTransaction creates a FileAppendTransaction transaction which can be
// used to construct and execute a File Append Transaction.
func NewFileAppendTransaction() *FileAppendTransaction {
	tx := &FileAppendTransaction{
		maxChunks: 20,
		contents:  make([]byte, 0),
		chunkSize: 2048,
	}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _FileAppendTransactionFromProtobuf(tx Transaction[*FileAppendTransaction], pb *services.TransactionBody) FileAppendTransaction {
	fileAppend := FileAppendTransaction{
		maxChunks: 20,
		contents:  pb.GetFileAppend().GetContents(),
		chunkSize: 2048,
		fileID:    _FileIDFromProtobuf(pb.GetFileAppend().GetFileID()),
	}

	tx.childTransaction = &fileAppend
	fileAppend.Transaction = &tx
	return fileAppend
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

// Execute executes the Transaction with the provided client
func (tx *FileAppendTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if tx.freezeError != nil {
		return TransactionResponse{}, tx.freezeError
	}

	list, err := tx.ExecuteAll(client)

	if err != nil {
		if len(list) > 0 {
			return TransactionResponse{
				TransactionID: tx.GetTransactionID(),
				NodeID:        list[0].NodeID,
				Hash:          make([]byte, 0),
			}, err
		}
		return TransactionResponse{
			TransactionID: tx.GetTransactionID(),
			Hash:          make([]byte, 0),
		}, err
	}

	return list[0], nil
}

// ExecuteAll executes the all the Transactions with the provided client
func (tx *FileAppendTransaction) ExecuteAll(
	client *Client,
) ([]TransactionResponse, error) {
	if client == nil || client.operator == nil {
		return []TransactionResponse{}, errNoClientProvided
	}

	if !tx.IsFrozen() {
		_, err := tx.FreezeWith(client)
		if err != nil {
			return []TransactionResponse{}, err
		}
	}

	var transactionID TransactionID
	if tx.transactionIDs._Length() > 0 {
		transactionID = tx.GetTransactionID()
	} else {
		return []TransactionResponse{}, errors.New("transactionID list is empty")
	}

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		tx.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	size := tx.signedTransactions._Length() / tx.nodeAccountIDs._Length()
	list := make([]TransactionResponse, size)

	for i := 0; i < size; i++ {
		resp, err := _Execute(client, tx)

		if err != nil {
			return list, err
		}

		list[i] = resp.(TransactionResponse)

		_, err = list[i].SetValidateStatus(false).GetReceipt(client)
		if err != nil {
			return list, err
		}
	}

	return list, nil
}

func (tx *FileAppendTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	chunks := uint64((len(tx.contents) + (tx.chunkSize - 1)) / tx.chunkSize)
	if chunks > 1 {
		return &ScheduleCreateTransaction{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: 1,
		}
	}

	return tx.Transaction.Schedule()
}

// ----------- Overridden functions ----------------

func (tx FileAppendTransaction) getName() string {
	return "FileAppendTransaction"
}
func (tx FileAppendTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx FileAppendTransaction) build() *services.TransactionBody {
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

func (tx FileAppendTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileAppend{
			FileAppend: tx.buildProtoBody(),
		},
	}, nil
}
func (tx FileAppendTransaction) buildProtoBody() *services.FileAppendTransactionBody {
	body := &services.FileAppendTransactionBody{
		Contents: tx.contents,
	}

	if tx.fileID != nil {
		body.FileID = tx.fileID._ToProtobuf()
	}

	return body
}

func (tx FileAppendTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().AppendContent,
	}
}

// TODO can be removed at some point
func (tx FileAppendTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx FileAppendTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
