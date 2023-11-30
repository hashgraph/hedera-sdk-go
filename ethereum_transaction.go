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

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// EthereumTransaction is used to create a EthereumTransaction transaction which can be used to construct and execute
// a Ethereum transaction.
type EthereumTransaction struct {
	transaction
	ethereumData  []byte
	callData      *FileID
	MaxGasAllowed int64
}

// NewEthereumTransaction creates a EthereumTransaction transaction which can be used to construct and execute
// a Ethereum transaction.
func NewEthereumTransaction() *EthereumTransaction {
	this := EthereumTransaction{
		transaction: _NewTransaction(),
	}
	this.e = &this
	this._SetDefaultMaxTransactionFee(NewHbar(2))

	return &this
}

func _EthereumTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *EthereumTransaction {
	resultTx := &EthereumTransaction{
		transaction:   this,
		ethereumData:  pb.GetEthereumTransaction().EthereumData,
		callData:      _FileIDFromProtobuf(pb.GetEthereumTransaction().CallData),
		MaxGasAllowed: pb.GetEthereumTransaction().MaxGasAllowance,
	}
	resultTx.e = resultTx
	return resultTx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *EthereumTransaction) SetGrpcDeadline(deadline *time.Duration) *EthereumTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetEthereumData
// The raw Ethereum transaction (RLP encoded type 0, 1, and 2). Complete
// unless the callData field is set.
func (this *EthereumTransaction) SetEthereumData(data []byte) *EthereumTransaction {
	this._RequireNotFrozen()
	this.ethereumData = data
	return this
}

// GetEthereumData returns the raw Ethereum transaction (RLP encoded type 0, 1, and 2).
func (this *EthereumTransaction) GetEthereumData() []byte {
	return this.ethereumData
}

// Deprecated
func (this *EthereumTransaction) SetCallData(file FileID) *EthereumTransaction {
	this._RequireNotFrozen()
	this.callData = &file
	return this
}

// SetCallDataFileID sets the file ID containing the call data.
func (this *EthereumTransaction) SetCallDataFileID(file FileID) *EthereumTransaction {
	this._RequireNotFrozen()
	this.callData = &file
	return this
}

// GetCallData
// For large transactions (for example contract create) this is the callData
// of the ethereumData. The data in the ethereumData will be re-written with
// the callData element as a zero length string with the original contents in
// the referenced file at time of execution. The ethereumData will need to be
// "rehydrated" with the callData for signature validation to pass.
func (this *EthereumTransaction) GetCallData() FileID {
	if this.callData != nil {
		return *this.callData
	}

	return FileID{}
}

// SetMaxGasAllowed
// The maximum amount, in tinybars, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (this *EthereumTransaction) SetMaxGasAllowed(gas int64) *EthereumTransaction {
	this._RequireNotFrozen()
	this.MaxGasAllowed = gas
	return this
}

// SetMaxGasAllowanceHbar sets the maximum amount, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (this *EthereumTransaction) SetMaxGasAllowanceHbar(gas Hbar) *EthereumTransaction {
	this._RequireNotFrozen()
	this.MaxGasAllowed = gas.AsTinybar()
	return this
}

// GetMaxGasAllowed returns the maximum amount, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (this *EthereumTransaction) GetMaxGasAllowed() int64 {
	return this.MaxGasAllowed
}

// Sign uses the provided privateKey to sign the transaction.
func (this *EthereumTransaction) Sign(
	privateKey PrivateKey,
) *EthereumTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *EthereumTransaction) SignWithOperator(
	client *Client,
) (*EthereumTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *EthereumTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *EthereumTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *EthereumTransaction) Freeze() (*EthereumTransaction, error) {
	_, err := this.transaction.Freeze()
	return this, err
}

func (this *EthereumTransaction) FreezeWith(client *Client) (*EthereumTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *EthereumTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *EthereumTransaction) SetMaxTransactionFee(fee Hbar) *EthereumTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *EthereumTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *EthereumTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *EthereumTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this EthereumTransaction.
func (this *EthereumTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this EthereumTransaction.
func (this *EthereumTransaction) SetTransactionMemo(memo string) *EthereumTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *EthereumTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this EthereumTransaction.
func (this *EthereumTransaction) SetTransactionValidDuration(duration time.Duration) *EthereumTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	EthereumTransaction.
func (this *EthereumTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this EthereumTransaction.
func (this *EthereumTransaction) SetTransactionID(transactionID TransactionID) *EthereumTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this EthereumTransaction.
func (this *EthereumTransaction) SetNodeAccountIDs(nodeID []AccountID) *EthereumTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *EthereumTransaction) SetMaxRetry(count int) *EthereumTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *EthereumTransaction) AddSignature(publicKey PublicKey, signature []byte) *EthereumTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *EthereumTransaction) SetMaxBackoff(max time.Duration) *EthereumTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *EthereumTransaction) SetMinBackoff(min time.Duration) *EthereumTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *EthereumTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("EthereumTransaction:%d", timestamp.UnixNano())
}

// ----------- overriden functions ----------------

func (this *EthereumTransaction) getName() string {
	return "EthereumTransaction"
}
func (this *EthereumTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.callData != nil {
		if err := this.callData.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *EthereumTransaction) build() *services.TransactionBody {
	body := &services.EthereumTransactionBody{
		EthereumData:    this.ethereumData,
		MaxGasAllowance: this.MaxGasAllowed,
	}

	if this.callData != nil {
		body.CallData = this.callData._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionID:            this.transactionID._ToProtobuf(),
		TransactionFee:           this.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		Memo:                     this.transaction.memo,
		Data: &services.TransactionBody_EthereumTransaction{
			EthereumTransaction: body,
		},
	}
}

func (this *EthereumTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `EthereumTransaction")
}

func (this *EthereumTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().CallEthereum,
	}
}
