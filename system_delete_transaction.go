package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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

	"github.com/hashgraph/hedera-sdk-go/v2/proto/services"
)

// Delete a file or smart contract - can only be done with a Hedera admin.
// When it is deleted, it immediately disappears from the system as seen by the user,
// but is still stored internally until the expiration time, at which time it
// is truly and permanently deleted.
// Until that time, it can be undeleted by the Hedera admin.
// When a smart contract is deleted, the cryptocurrency account within it continues
// to exist, and is not affected by the expiration time here.
type SystemDeleteTransaction struct {
	Transaction
	contractID     *ContractID
	fileID         *FileID
	expirationTime *time.Time
}

// NewSystemDeleteTransaction creates a SystemDeleteTransaction transaction which can be
// used to construct and execute a System Delete Transaction.
func NewSystemDeleteTransaction() *SystemDeleteTransaction {
	tx := SystemDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _SystemDeleteTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *SystemDeleteTransaction {
	expiration := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(),
		int(pb.GetSystemDelete().ExpirationTime.Seconds), time.Now().Nanosecond(), time.Now().Location(),
	)
	return &SystemDeleteTransaction{
		Transaction:    tx,
		contractID:     _ContractIDFromProtobuf(pb.GetSystemDelete().GetContractID()),
		fileID:         _FileIDFromProtobuf(pb.GetSystemDelete().GetFileID()),
		expirationTime: &expiration,
	}
}

// SetExpirationTime sets the time at which this transaction will expire.
func (tx *SystemDeleteTransaction) SetExpirationTime(expiration time.Time) *SystemDeleteTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expiration
	return tx
}

// GetExpirationTime returns the time at which this transaction will expire.
func (tx *SystemDeleteTransaction) GetExpirationTime() int64 {
	if tx.expirationTime != nil {
		return tx.expirationTime.Unix()
	}

	return 0
}

// SetContractID sets the ContractID of the contract which will be deleted.
func (tx *SystemDeleteTransaction) SetContractID(contractID ContractID) *SystemDeleteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// GetContractID returns the ContractID of the contract which will be deleted.
func (tx *SystemDeleteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// SetFileID sets the FileID of the file which will be deleted.
func (tx *SystemDeleteTransaction) SetFileID(fileID FileID) *SystemDeleteTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID of the file which will be deleted.
func (tx *SystemDeleteTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *SystemDeleteTransaction) Sign(privateKey PrivateKey) *SystemDeleteTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *SystemDeleteTransaction) SignWithOperator(client *Client) (*SystemDeleteTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *SystemDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *SystemDeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *SystemDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *SystemDeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *SystemDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *SystemDeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *SystemDeleteTransaction) Freeze() (*SystemDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *SystemDeleteTransaction) FreezeWith(client *Client) (*SystemDeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this SystemDeleteTransaction.
func (tx *SystemDeleteTransaction) SetMaxTransactionFee(fee Hbar) *SystemDeleteTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *SystemDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *SystemDeleteTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this SystemDeleteTransaction.
func (tx *SystemDeleteTransaction) SetTransactionMemo(memo string) *SystemDeleteTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this SystemDeleteTransaction.
func (tx *SystemDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *SystemDeleteTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *SystemDeleteTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this SystemDeleteTransaction.
func (tx *SystemDeleteTransaction) SetTransactionID(transactionID TransactionID) *SystemDeleteTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this SystemDeleteTransaction.
func (tx *SystemDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *SystemDeleteTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *SystemDeleteTransaction) SetMaxRetry(count int) *SystemDeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *SystemDeleteTransaction) SetMaxBackoff(max time.Duration) *SystemDeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *SystemDeleteTransaction) SetMinBackoff(min time.Duration) *SystemDeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *SystemDeleteTransaction) SetLogLevel(level LogLevel) *SystemDeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *SystemDeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *SystemDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *SystemDeleteTransaction) getName() string {
	return "SystemDeleteTransaction"
}

func (tx *SystemDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.contractID != nil {
		if err := tx.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.fileID != nil {
		if err := tx.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *SystemDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_SystemDelete{
			SystemDelete: tx.buildProtoBody(),
		},
	}
}

func (tx *SystemDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_SystemDelete{
			SystemDelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *SystemDeleteTransaction) buildProtoBody() *services.SystemDeleteTransactionBody {
	body := &services.SystemDeleteTransactionBody{}

	if tx.expirationTime != nil {
		body.ExpirationTime = &services.TimestampSeconds{
			Seconds: tx.expirationTime.Unix(),
		}
	}

	if tx.contractID != nil {
		body.Id = &services.SystemDeleteTransactionBody_ContractID{
			ContractID: tx.contractID._ToProtobuf(),
		}
	}

	if tx.fileID != nil {
		body.Id = &services.SystemDeleteTransactionBody_FileID{
			FileID: tx.fileID._ToProtobuf(),
		}
	}

	return body
}

func (tx *SystemDeleteTransaction) getMethod(channel *_Channel) _Method {
	if channel._GetContract() == nil {
		return _Method{
			transaction: channel._GetFile().SystemDelete,
		}
	}

	return _Method{
		transaction: channel._GetContract().SystemDelete,
	}
}

func (tx *SystemDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
