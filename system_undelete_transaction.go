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
	"github.com/hashgraph/hedera-sdk-go/v2/proto/services"

	"time"
)

// Undelete a file or smart contract that was deleted by AdminDelete.
// Can only be done with a Hedera admin.
type SystemUndeleteTransaction struct {
	Transaction
	contractID *ContractID
	fileID     *FileID
}

// NewSystemUndeleteTransaction creates a SystemUndeleteTransaction transaction which can be
// used to construct and execute a System Undelete Transaction.
func NewSystemUndeleteTransaction() *SystemUndeleteTransaction {
	tx := SystemUndeleteTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _SystemUndeleteTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *SystemUndeleteTransaction {
	return &SystemUndeleteTransaction{
		Transaction: tx,
		contractID:  _ContractIDFromProtobuf(pb.GetSystemUndelete().GetContractID()),
		fileID:      _FileIDFromProtobuf(pb.GetSystemUndelete().GetFileID()),
	}
}

// SetContractID sets the ContractID of the contract whose deletion is being undone.
func (tx *SystemUndeleteTransaction) SetContractID(contractID ContractID) *SystemUndeleteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// GetContractID returns the ContractID of the contract whose deletion is being undone.
func (tx *SystemUndeleteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// SetFileID sets the FileID of the file whose deletion is being undone.
func (tx *SystemUndeleteTransaction) SetFileID(fileID FileID) *SystemUndeleteTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID of the file whose deletion is being undone.
func (tx *SystemUndeleteTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *SystemUndeleteTransaction) Sign(privateKey PrivateKey) *SystemUndeleteTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *SystemUndeleteTransaction) SignWithOperator(client *Client) (*SystemUndeleteTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *SystemUndeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *SystemUndeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *SystemUndeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *SystemUndeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *SystemUndeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *SystemUndeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *SystemUndeleteTransaction) Freeze() (*SystemUndeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *SystemUndeleteTransaction) FreezeWith(client *Client) (*SystemUndeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this SystemUndeleteTransaction.
func (tx *SystemUndeleteTransaction) SetMaxTransactionFee(fee Hbar) *SystemUndeleteTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *SystemUndeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *SystemUndeleteTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this SystemUndeleteTransaction.
func (tx *SystemUndeleteTransaction) SetTransactionMemo(memo string) *SystemUndeleteTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this SystemUndeleteTransaction.
func (tx *SystemUndeleteTransaction) SetTransactionValidDuration(duration time.Duration) *SystemUndeleteTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *SystemUndeleteTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this SystemUndeleteTransaction.
func (tx *SystemUndeleteTransaction) SetTransactionID(transactionID TransactionID) *SystemUndeleteTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this SystemUndeleteTransaction.
func (tx *SystemUndeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *SystemUndeleteTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *SystemUndeleteTransaction) SetMaxRetry(count int) *SystemUndeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *SystemUndeleteTransaction) SetMaxBackoff(max time.Duration) *SystemUndeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *SystemUndeleteTransaction) SetMinBackoff(min time.Duration) *SystemUndeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *SystemUndeleteTransaction) SetLogLevel(level LogLevel) *SystemUndeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *SystemUndeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *SystemUndeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *SystemUndeleteTransaction) getName() string {
	return "SystemUndeleteTransaction"
}

func (tx *SystemUndeleteTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *SystemUndeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_SystemUndelete{
			SystemUndelete: tx.buildProtoBody(),
		},
	}
}

func (tx *SystemUndeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_SystemUndelete{
			SystemUndelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *SystemUndeleteTransaction) buildProtoBody() *services.SystemUndeleteTransactionBody {
	body := &services.SystemUndeleteTransactionBody{}
	if tx.contractID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_ContractID{
			ContractID: tx.contractID._ToProtobuf(),
		}
	}

	if tx.fileID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_FileID{
			FileID: tx.fileID._ToProtobuf(),
		}
	}

	return body
}

func (tx *SystemUndeleteTransaction) getMethod(channel *_Channel) _Method {
	if channel._GetContract() == nil {
		return _Method{
			transaction: channel._GetFile().SystemUndelete,
		}
	}

	return _Method{
		transaction: channel._GetContract().SystemUndelete,
	}
}
func (tx *SystemUndeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
