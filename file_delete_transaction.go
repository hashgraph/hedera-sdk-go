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
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// FileDeleteTransaction Deletes the given file. After deletion, it will be marked as deleted and will have no contents.
// But information about it will continue to exist until it expires. A list of keys was given when
// the file was created. All the top level keys on that list must sign transactions to create or
// modify the file, but any single one of the top level keys can be used to delete the file. tx
// transaction must be signed by 1-of-M KeyList keys. If keys contains additional KeyList or
// ThresholdKey then 1-of-M secondary KeyList or ThresholdKey signing requirements must be meet.
type FileDeleteTransaction struct {
	Transaction
	fileID *FileID
}

// NewFileDeleteTransaction creates a FileDeleteTransaction which deletes the given file. After deletion,
// it will be marked as deleted and will have no contents.
// But information about it will continue to exist until it expires. A list of keys was given when
// the file was created. All the top level keys on that list must sign transactions to create or
// modify the file, but any single one of the top level keys can be used to delete the file. tx
// transaction must be signed by 1-of-M KeyList keys. If keys contains additional KeyList or
// ThresholdKey then 1-of-M secondary KeyList or ThresholdKey signing requirements must be meet.
func NewFileDeleteTransaction() *FileDeleteTransaction {
	tx := FileDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _FileDeleteTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *FileDeleteTransaction {
	resultTx := &FileDeleteTransaction{
		Transaction: tx,
		fileID:      _FileIDFromProtobuf(pb.GetFileDelete().GetFileID()),
	}
	return resultTx
}

// SetFileID Sets the FileID of the file to be deleted
func (tx *FileDeleteTransaction) SetFileID(fileID FileID) *FileDeleteTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID of the file to be deleted
func (tx *FileDeleteTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *FileDeleteTransaction) Sign(
	privateKey PrivateKey,
) *FileDeleteTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *FileDeleteTransaction) SignWithOperator(
	client *Client,
) (*FileDeleteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *FileDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileDeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *FileDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileDeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when tx deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *FileDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *FileDeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *FileDeleteTransaction) Freeze() (*FileDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *FileDeleteTransaction) FreezeWith(client *Client) (*FileDeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *FileDeleteTransaction) GetMaxTransactionFee() Hbar {
	return tx.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *FileDeleteTransaction) SetMaxTransactionFee(fee Hbar) *FileDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *FileDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FileDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx FileDeleteTransaction.
func (tx *FileDeleteTransaction) SetTransactionMemo(memo string) *FileDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for tx FileDeleteTransaction.
func (tx *FileDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *FileDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for tx FileDeleteTransaction.
func (tx *FileDeleteTransaction) SetTransactionID(transactionID TransactionID) *FileDeleteTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountID sets the _Node AccountID for tx FileDeleteTransaction.
func (tx *FileDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *FileDeleteTransaction) SetMaxRetry(count int) *FileDeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches tx time.
func (tx *FileDeleteTransaction) SetMaxBackoff(max time.Duration) *FileDeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *FileDeleteTransaction) SetMinBackoff(min time.Duration) *FileDeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *FileDeleteTransaction) SetLogLevel(level LogLevel) *FileDeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *FileDeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *FileDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *FileDeleteTransaction) getName() string {
	return "FileDeleteTransaction"
}
func (tx *FileDeleteTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *FileDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileDelete{
			FileDelete: tx.buildProtoBody(),
		},
	}
}

func (tx *FileDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileDelete{
			FileDelete: tx.buildProtoBody(),
		},
	}, nil
}
func (tx *FileDeleteTransaction) buildProtoBody() *services.FileDeleteTransactionBody {
	body := &services.FileDeleteTransactionBody{}
	if tx.fileID != nil {
		body.FileID = tx.fileID._ToProtobuf()
	}
	return body
}

func (tx *FileDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().DeleteFile,
	}
}
func (tx *FileDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
