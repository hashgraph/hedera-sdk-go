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

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-sdk-go/v2/generated/services"
)

// EthereumTransaction is used to create a EthereumTransaction transaction which can be used to construct and execute
// a Ethereum Transaction.
type EthereumTransaction struct {
	Transaction
	ethereumData  []byte
	callData      *FileID
	MaxGasAllowed int64
}

// NewEthereumTransaction creates a EthereumTransaction transaction which can be used to construct and execute
// a Ethereum Transaction.
func NewEthereumTransaction() *EthereumTransaction {
	tx := EthereumTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _EthereumTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *EthereumTransaction {
	return &EthereumTransaction{
		Transaction:   tx,
		ethereumData:  pb.GetEthereumTransaction().EthereumData,
		callData:      _FileIDFromProtobuf(pb.GetEthereumTransaction().CallData),
		MaxGasAllowed: pb.GetEthereumTransaction().MaxGasAllowance,
	}
}

// SetEthereumData
// The raw Ethereum transaction (RLP encoded type 0, 1, and 2). Complete
// unless the callData field is set.
func (tx *EthereumTransaction) SetEthereumData(data []byte) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.ethereumData = data
	return tx
}

// GetEthereumData returns the raw Ethereum transaction (RLP encoded type 0, 1, and 2).
func (tx *EthereumTransaction) GetEthereumData() []byte {
	return tx.ethereumData
}

// Deprecated
func (tx *EthereumTransaction) SetCallData(file FileID) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.callData = &file
	return tx
}

// SetCallDataFileID sets the file ID containing the call data.
func (tx *EthereumTransaction) SetCallDataFileID(file FileID) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.callData = &file
	return tx
}

// GetCallData
// For large transactions (for example contract create) this is the callData
// of the ethereumData. The data in the ethereumData will be re-written with
// the callData element as a zero length string with the original contents in
// the referenced file at time of execution. The ethereumData will need to be
// "rehydrated" with the callData for signature validation to pass.
func (tx *EthereumTransaction) GetCallData() FileID {
	if tx.callData != nil {
		return *tx.callData
	}

	return FileID{}
}

// SetMaxGasAllowed
// The maximum amount, in tinybars, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (tx *EthereumTransaction) SetMaxGasAllowed(gas int64) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.MaxGasAllowed = gas
	return tx
}

// SetMaxGasAllowanceHbar sets the maximum amount, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (tx *EthereumTransaction) SetMaxGasAllowanceHbar(gas Hbar) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.MaxGasAllowed = gas.AsTinybar()
	return tx
}

// GetMaxGasAllowed returns the maximum amount, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (tx *EthereumTransaction) GetMaxGasAllowed() int64 {
	return tx.MaxGasAllowed
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *EthereumTransaction) Sign(
	privateKey PrivateKey,
) *EthereumTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *EthereumTransaction) SignWithOperator(
	client *Client,
) (*EthereumTransaction, error) {
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
func (tx *EthereumTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *EthereumTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *EthereumTransaction) AddSignature(publicKey PublicKey, signature []byte) *EthereumTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when tx deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *EthereumTransaction) SetGrpcDeadline(deadline *time.Duration) *EthereumTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *EthereumTransaction) Freeze() (*EthereumTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *EthereumTransaction) FreezeWith(client *Client) (*EthereumTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *EthereumTransaction) SetMaxTransactionFee(fee Hbar) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *EthereumTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (tx *EthereumTransaction) GetRegenerateTransactionID() bool {
	return tx.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this EthereumTransaction.
func (tx *EthereumTransaction) GetTransactionMemo() string {
	return tx.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this EthereumTransaction.
func (tx *EthereumTransaction) SetTransactionMemo(memo string) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this EthereumTransaction.
func (tx *EthereumTransaction) SetTransactionValidDuration(duration time.Duration) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *EthereumTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this EthereumTransaction.
func (tx *EthereumTransaction) SetTransactionID(transactionID TransactionID) *EthereumTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this EthereumTransaction.
func (tx *EthereumTransaction) SetNodeAccountIDs(nodeID []AccountID) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *EthereumTransaction) SetMaxRetry(count int) *EthereumTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *EthereumTransaction) SetMaxBackoff(max time.Duration) *EthereumTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *EthereumTransaction) SetMinBackoff(min time.Duration) *EthereumTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *EthereumTransaction) SetLogLevel(level LogLevel) *EthereumTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *EthereumTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *EthereumTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *EthereumTransaction) getName() string {
	return "EthereumTransaction"
}
func (tx *EthereumTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.callData != nil {
		if err := tx.callData.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *EthereumTransaction) build() *services.TransactionBody {
	body := &services.EthereumTransactionBody{
		EthereumData:    tx.ethereumData,
		MaxGasAllowance: tx.MaxGasAllowed,
	}

	if tx.callData != nil {
		body.CallData = tx.callData._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionID:            tx.transactionID._ToProtobuf(),
		TransactionFee:           tx.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		Memo:                     tx.Transaction.memo,
		Data: &services.TransactionBody_EthereumTransaction{
			EthereumTransaction: body,
		},
	}
}

func (tx *EthereumTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `EthereumTransaction")
}

func (tx *EthereumTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().CallEthereum,
	}
}
