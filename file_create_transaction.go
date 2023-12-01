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
	tx := FileCreateTransaction{
		transaction: _NewTransaction(),
	}

	tx.SetExpirationTime(time.Now().Add(7890000 * time.Second))
	tx._SetDefaultMaxTransactionFee(NewHbar(5))
	tx.e = &tx

	return &tx
}

func _FileCreateTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *FileCreateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileCreate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileCreate().GetExpirationTime())

	resultTx := &FileCreateTransaction{
		transaction:    tx,
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileCreate().GetContents(),
		memo:           pb.GetMemo(),
	}
	resultTx.e = resultTx
	return resultTx
}

// AddKey adds a key to the internal list of keys associated with the file. All of the keys on the list must sign to
// create or modify a file, but only one of them needs to sign in order to delete the file. Each of those "keys" may
// itself be threshold key containing other keys (including other threshold keys). In other words, the behavior is an
// AND for create/modify, OR for delete. tx is useful for acting as a revocation server. If it is desired to have the
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

// SetExpirationTime sets the time at which tx file should expire (unless FileUpdateTransaction is used before then to
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

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *FileCreateTransaction) Sign(
	privateKey PrivateKey,
) *FileCreateTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *FileCreateTransaction) SignWithOperator(
	client *Client,
) (*FileCreateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *FileCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileCreateTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *FileCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileCreateTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when tx deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *FileCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *FileCreateTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *FileCreateTransaction) Freeze() (*FileCreateTransaction, error) {
	_, err := tx.transaction.Freeze()
	return tx, err
}

func (tx *FileCreateTransaction) FreezeWith(client *Client) (*FileCreateTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *FileCreateTransaction) SetMaxTransactionFee(fee Hbar) *FileCreateTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *FileCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileCreateTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx FileCreateTransaction.
func (tx *FileCreateTransaction) SetTransactionMemo(memo string) *FileCreateTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for tx FileCreateTransaction.
func (tx *FileCreateTransaction) SetTransactionValidDuration(duration time.Duration) *FileCreateTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for tx FileCreateTransaction.
func (tx *FileCreateTransaction) SetTransactionID(transactionID TransactionID) *FileCreateTransaction {
	tx._RequireNotFrozen()

	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountID sets the _Node AccountID for tx FileCreateTransaction.
func (tx *FileCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileCreateTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *FileCreateTransaction) SetMaxRetry(count int) *FileCreateTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches tx time.
func (tx *FileCreateTransaction) SetMaxBackoff(max time.Duration) *FileCreateTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *FileCreateTransaction) SetMinBackoff(min time.Duration) *FileCreateTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *FileCreateTransaction) SetLogLevel(level LogLevel) *FileCreateTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *FileCreateTransaction) getName() string {
	return "FileCreateTransaction"
}
func (tx *FileCreateTransaction) build() *services.TransactionBody {
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

	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileCreate{
			FileCreate: body,
		},
	}
}

func (tx *FileCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
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

	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_FileCreate{
			FileCreate: body,
		},
	}, nil
}
func (tx *FileCreateTransaction) buildProtoBody() *services.FileCreateTransactionBody {
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

func (tx *FileCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().CreateFile,
	}
}
func (tx *FileCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

