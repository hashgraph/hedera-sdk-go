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
	Transaction
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
	tx := FileUpdateTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(5))
	return &tx
}

func _FileUpdateTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *FileUpdateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileUpdate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileUpdate().GetExpirationTime())

	return &FileUpdateTransaction{
		Transaction:    tx,
		fileID:         _FileIDFromProtobuf(pb.GetFileUpdate().GetFileID()),
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileUpdate().GetContents(),
		memo:           pb.GetFileUpdate().GetMemo().Value,
	}
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

// ----- Required Interfaces ------- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *FileUpdateTransaction) Sign(
	privateKey PrivateKey,
) *FileUpdateTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *FileUpdateTransaction) SignWithOperator(
	client *Client,
) (*FileUpdateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *FileUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileUpdateTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *FileUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileUpdateTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when tx deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *FileUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *FileUpdateTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *FileUpdateTransaction) Freeze() (*FileUpdateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *FileUpdateTransaction) FreezeWith(client *Client) (*FileUpdateTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *FileUpdateTransaction) SetMaxTransactionFee(fee Hbar) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *FileUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this FileUpdateTransaction.
func (tx *FileUpdateTransaction) SetTransactionMemo(memo string) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this FileUpdateTransaction.
func (tx *FileUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this FileUpdateTransaction.
func (tx *FileUpdateTransaction) SetTransactionID(transactionID TransactionID) *FileUpdateTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountID sets the _Node AccountID for this FileUpdateTransaction.
func (tx *FileUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *FileUpdateTransaction) SetMaxRetry(count int) *FileUpdateTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *FileUpdateTransaction) SetMaxBackoff(max time.Duration) *FileUpdateTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *FileUpdateTransaction) SetMinBackoff(min time.Duration) *FileUpdateTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *FileUpdateTransaction) SetLogLevel(level LogLevel) *FileUpdateTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *FileUpdateTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *FileUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *FileUpdateTransaction) getName() string {
	return "FileUpdateTransaction"
}
func (tx *FileUpdateTransaction) validateNetworkOnIDs(client *Client) error {
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
func (tx *FileUpdateTransaction) build() *services.TransactionBody {
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
func (tx *FileUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileUpdate{
			FileUpdate: tx.buildProtoBody(),
		},
	}, nil
}
func (tx *FileUpdateTransaction) buildProtoBody() *services.FileUpdateTransactionBody {
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
func (tx *FileUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().UpdateFile,
	}
}
func (tx *FileUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
