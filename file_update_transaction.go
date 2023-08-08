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
// This transaction must be signed by all the keys in the top level of a key list (M-of-M) of the file being updated.
// If the keys themselves are being updated, then the transaction must also be signed by all the new keys. If the keys contain
// additional KeyList or ThresholdKey then M-of-M secondary KeyList or ThresholdKey signing
// requirements must be meet
func NewFileUpdateTransaction() *FileUpdateTransaction {
	transaction := FileUpdateTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _FileUpdateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *FileUpdateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileUpdate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileUpdate().GetExpirationTime())

	return &FileUpdateTransaction{
		Transaction:    transaction,
		fileID:         _FileIDFromProtobuf(pb.GetFileUpdate().GetFileID()),
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileUpdate().GetContents(),
		memo:           pb.GetFileUpdate().GetMemo().Value,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *FileUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *FileUpdateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetFileID Sets the FileID to be updated
func (transaction *FileUpdateTransaction) SetFileID(fileID FileID) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.fileID = &fileID
	return transaction
}

// GetFileID returns the FileID to be updated
func (transaction *FileUpdateTransaction) GetFileID() FileID {
	if transaction.fileID == nil {
		return FileID{}
	}

	return *transaction.fileID
}

// SetKeys Sets the new list of keys that can modify or delete the file
func (transaction *FileUpdateTransaction) SetKeys(keys ...Key) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	if transaction.keys == nil {
		transaction.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.keys = keyList

	return transaction
}

func (transaction *FileUpdateTransaction) GetKeys() KeyList {
	if transaction.keys != nil {
		return *transaction.keys
	}

	return KeyList{}
}

// SetExpirationTime Sets the new expiry time
func (transaction *FileUpdateTransaction) SetExpirationTime(expiration time.Time) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

// GetExpirationTime returns the new expiry time
func (transaction *FileUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

// SetContents Sets the new contents that should overwrite the file's current contents
func (transaction *FileUpdateTransaction) SetContents(contents []byte) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.contents = contents
	return transaction
}

// GetContents returns the new contents that should overwrite the file's current contents
func (transaction *FileUpdateTransaction) GetContents() []byte {
	return transaction.contents
}

// SetFileMemo Sets the new memo to be associated with the file (UTF-8 encoding max 100 bytes)
func (transaction *FileUpdateTransaction) SetFileMemo(memo string) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo

	return transaction
}

// GeFileMemo
// Deprecated: use GetFileMemo()
func (transaction *FileUpdateTransaction) GeFileMemo() string {
	return transaction.memo
}

func (transaction *FileUpdateTransaction) GetFileMemo() string {
	return transaction.memo
}

func (transaction *FileUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
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

func (transaction *FileUpdateTransaction) _Build() *services.TransactionBody {
	body := &services.FileUpdateTransactionBody{
		Memo: &wrapperspb.StringValue{Value: transaction.memo},
	}
	if transaction.fileID != nil {
		body.FileID = transaction.fileID._ToProtobuf()
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.keys != nil {
		body.Keys = transaction.keys._ToProtoKeyList()
	}

	if transaction.contents != nil {
		body.Contents = transaction.contents
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileUpdate{
			FileUpdate: body,
		},
	}
}

func (transaction *FileUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *FileUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.FileUpdateTransactionBody{
		Memo: &wrapperspb.StringValue{Value: transaction.memo},
	}
	if transaction.fileID != nil {
		body.FileID = transaction.fileID._ToProtobuf()
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.keys != nil {
		body.Keys = transaction.keys._ToProtoKeyList()
	}

	if transaction.contents != nil {
		body.Contents = transaction.contents
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileUpdate{
			FileUpdate: body,
		},
	}, nil
}

func _FileUpdateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().UpdateFile,
	}
}

func (transaction *FileUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *FileUpdateTransaction) Sign(
	privateKey PrivateKey,
) *FileUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *FileUpdateTransaction) SignWithOperator(
	client *Client,
) (*FileUpdateTransaction, error) {
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
func (transaction *FileUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileUpdateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *FileUpdateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_FileUpdateTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
		transaction.maxBackoff,
		transaction.minBackoff,
		transaction.maxRetry,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID:  transaction.GetTransactionID(),
			NodeID:         resp.(TransactionResponse).NodeID,
			ValidateStatus: true,
		}, err
	}

	return TransactionResponse{
		TransactionID:  transaction.GetTransactionID(),
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
	}, nil
}

func (transaction *FileUpdateTransaction) Freeze() (*FileUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FileUpdateTransaction) FreezeWith(client *Client) (*FileUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &FileUpdateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *FileUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *FileUpdateTransaction) SetMaxTransactionFee(fee Hbar) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *FileUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *FileUpdateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *FileUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionMemo(memo string) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *FileUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	FileUpdateTransaction.
func (transaction *FileUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionID(transactionID TransactionID) *FileUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *FileUpdateTransaction) SetMaxRetry(count int) *FileUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *FileUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileUpdateTransaction {
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
func (transaction *FileUpdateTransaction) SetMaxBackoff(max time.Duration) *FileUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *FileUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *FileUpdateTransaction) SetMinBackoff(min time.Duration) *FileUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *FileUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *FileUpdateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("FileUpdateTransaction:%d", timestamp.UnixNano())
}

func (transaction *FileUpdateTransaction) SetLogLevel(level LogLevel) *FileUpdateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
