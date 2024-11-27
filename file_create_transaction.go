package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
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
	*Transaction[*FileCreateTransaction]
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
	tx := &FileCreateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx.SetExpirationTime(time.Now().Add(7890000 * time.Second))
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _FileCreateTransactionFromProtobuf(tx Transaction[*FileCreateTransaction], pb *services.TransactionBody) FileCreateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileCreate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileCreate().GetExpirationTime())

	fileCreateTransaction := FileCreateTransaction{
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileCreate().GetContents(),
		memo:           pb.GetMemo(),
	}
	tx.childTransaction = &fileCreateTransaction
	fileCreateTransaction.Transaction = &tx
	return fileCreateTransaction
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
func (tx *FileCreateTransaction) SetKeys(keys ...Key) *FileCreateTransaction {
	tx._RequireNotFrozen()
	if tx.keys == nil {
		tx.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	tx.keys = keyList

	return tx
}

func (tx *FileCreateTransaction) GetKeys() KeyList {
	if tx.keys != nil {
		return *tx.keys
	}

	return KeyList{}
}

// SetExpirationTime sets the time at which this file should expire (unless FileUpdateTransaction is used before then to
// extend its life). The file will automatically disappear at the fileExpirationTime, unless its expiration is extended
// by another transaction before that time. If the file is deleted, then its contents will become empty and it will be
// marked as deleted until it expires, and then it will cease to exist.
func (tx *FileCreateTransaction) SetExpirationTime(expiration time.Time) *FileCreateTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expiration
	return tx
}

func (tx *FileCreateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}

	return time.Time{}
}

// SetContents sets the bytes that are the contents of the file (which can be empty). If the size of the file and other
// fields in the transaction exceed the max transaction size then FileCreateTransaction can be used to continue
// uploading the file.
func (tx *FileCreateTransaction) SetContents(contents []byte) *FileCreateTransaction {
	tx._RequireNotFrozen()
	tx.contents = contents
	return tx
}

// GetContents returns the bytes that are the contents of the file (which can be empty).
func (tx *FileCreateTransaction) GetContents() []byte {
	return tx.contents
}

// SetMemo Sets the memo associated with the file (UTF-8 encoding max 100 bytes)
func (tx *FileCreateTransaction) SetMemo(memo string) *FileCreateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo
	return tx
}

// GetMemo returns the memo associated with the file (UTF-8 encoding max 100 bytes)
func (tx *FileCreateTransaction) GetMemo() string {
	return tx.memo
}

// ----------- Overridden functions ----------------

func (tx FileCreateTransaction) getName() string {
	return "FileCreateTransaction"
}
func (tx FileCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileCreate{
			FileCreate: tx.buildProtoBody(),
		},
	}
}

func (tx FileCreateTransaction) validateNetworkOnIDs(client *Client) error {
	return nil
}

func (tx FileCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileCreate{
			FileCreate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx FileCreateTransaction) buildProtoBody() *services.FileCreateTransactionBody {
	body := &services.FileCreateTransactionBody{
		Memo: tx.memo,
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

func (tx FileCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().CreateFile,
	}
}
func (tx FileCreateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx FileCreateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, tx)
}
