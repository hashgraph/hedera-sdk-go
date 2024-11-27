package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// FileUpdateTransaction
// Modify the metadata and/or contents of a file. If a field is not set in the transaction body, the
// corresponding file attribute will be unchanged. This transaction must be signed by all the keys
// in the top level of a key list (M-of-M) of the file being updated. If the keys themselves are
// being updated, then the transaction must also be signed by all the new keys. If the keys contain
// additional KeyList or ThresholdKey then M-of-M secondary KeyList or ThresholdKey signing
// requirements must be meet
type FileUpdateTransaction struct {
	*Transaction[*FileUpdateTransaction]
	fileID         *FileID
	keys           *KeyList
	expirationTime *time.Time
	contents       []byte
	memo           string
}

// NewFileUpdateTransaction creates a FileUpdateTransaction which modifies the metadata and/or contents of a file.
// If a field is not set in the transaction body, the corresponding file attribute will be unchanged.
// tx transaction must be signed by all the keys in the top level of a key list (M-of-M) of the file being updated.
// If the keys themselves are being updated, then the transaction must also be signed by all the new keys. If the keys contain
// additional KeyList or ThresholdKey then M-of-M secondary KeyList or ThresholdKey signing
// requirements must be meet
func NewFileUpdateTransaction() *FileUpdateTransaction {
	tx := &FileUpdateTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))
	return tx
}

func _FileUpdateTransactionFromProtobuf(tx Transaction[*FileUpdateTransaction], pb *services.TransactionBody) FileUpdateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileUpdate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileUpdate().GetExpirationTime())

	fileUpdateTransaction := FileUpdateTransaction{
		fileID:         _FileIDFromProtobuf(pb.GetFileUpdate().GetFileID()),
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileUpdate().GetContents(),
		memo:           pb.GetFileUpdate().GetMemo().Value,
	}

	tx.childTransaction = &fileUpdateTransaction
	fileUpdateTransaction.Transaction = &tx
	return fileUpdateTransaction
}

// SetFileID Sets the FileID to be updated
func (tx *FileUpdateTransaction) SetFileID(fileID FileID) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID to be updated
func (tx *FileUpdateTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// SetKeys Sets the new list of keys that can modify or delete the file
func (tx *FileUpdateTransaction) SetKeys(keys ...Key) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	if tx.keys == nil {
		tx.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	tx.keys = keyList

	return tx
}

func (tx *FileUpdateTransaction) GetKeys() KeyList {
	if tx.keys != nil {
		return *tx.keys
	}

	return KeyList{}
}

// SetExpirationTime Sets the new expiry time
func (tx *FileUpdateTransaction) SetExpirationTime(expiration time.Time) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expiration
	return tx
}

// GetExpirationTime returns the new expiry time
func (tx *FileUpdateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}

	return time.Time{}
}

// SetContents Sets the new contents that should overwrite the file's current contents
func (tx *FileUpdateTransaction) SetContents(contents []byte) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.contents = contents
	return tx
}

// GetContents returns the new contents that should overwrite the file's current contents
func (tx *FileUpdateTransaction) GetContents() []byte {
	return tx.contents
}

// SetFileMemo Sets the new memo to be associated with the file (UTF-8 encoding max 100 bytes)
func (tx *FileUpdateTransaction) SetFileMemo(memo string) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo

	return tx
}

// GeFileMemo
// Deprecated: use GetFileMemo()
func (tx *FileUpdateTransaction) GeFileMemo() string {
	return tx.memo
}

func (tx *FileUpdateTransaction) GetFileMemo() string {
	return tx.memo
}

// ----------- Overridden functions ----------------

func (tx FileUpdateTransaction) getName() string {
	return "FileUpdateTransaction"
}
func (tx FileUpdateTransaction) validateNetworkOnIDs(client *Client) error {
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
func (tx FileUpdateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileUpdate{
			FileUpdate: tx.buildProtoBody(),
		},
	}
}
func (tx FileUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileUpdate{
			FileUpdate: tx.buildProtoBody(),
		},
	}, nil
}
func (tx FileUpdateTransaction) buildProtoBody() *services.FileUpdateTransactionBody {
	body := &services.FileUpdateTransactionBody{
		Memo: &wrapperspb.StringValue{Value: tx.memo},
	}
	if tx.fileID != nil {
		body.FileID = tx.fileID._ToProtobuf()
	}

	if tx.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*tx.expirationTime)
	}

	if tx.keys != nil {
		body.Keys = tx.keys._ToProtoKeyList()
	}

	if tx.contents != nil {
		body.Contents = tx.contents
	}

	return body
}

func (tx FileUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().UpdateFile,
	}
}

func (tx FileUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx FileUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
