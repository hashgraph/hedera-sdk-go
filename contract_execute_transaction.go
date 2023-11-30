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

	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// ContractExecuteTransaction calls a function of the given smart contract instance, giving it ContractFuncionParams as
// its inputs. it can use the given amount of gas, and any unspent gas will be refunded to the paying account.
//
// If this function stores information, it is charged gas to store it. There is a fee in hbars to maintain that storage
// until the expiration time, and that fee is added as part of the transaction fee.
//
// For a cheaper but more limited _Method to call functions, see ContractCallQuery.
type ContractExecuteTransaction struct {
	transaction
	contractID *ContractID
	gas        int64
	amount     int64
	parameters []byte
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call transaction.
func NewContractExecuteTransaction() *ContractExecuteTransaction {
	this := ContractExecuteTransaction{
		transaction: _NewTransaction(),
	}
	this._SetDefaultMaxTransactionFee(NewHbar(2))
	this.e = &this

	return &this
}

func _ContractExecuteTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *ContractExecuteTransaction {
	resultTx := &ContractExecuteTransaction{
		transaction: this,
		contractID:  _ContractIDFromProtobuf(pb.GetContractCall().GetContractID()),
		gas:         pb.GetContractCall().GetGas(),
		amount:      pb.GetContractCall().GetAmount(),
		parameters:  pb.GetContractCall().GetFunctionParameters(),
	}
	resultTx.e = resultTx
	return resultTx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *ContractExecuteTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractExecuteTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetContractID sets the contract instance to call.
func (this *ContractExecuteTransaction) SetContractID(contractID ContractID) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.contractID = &contractID
	return this
}

// GetContractID returns the contract instance to call.
func (this *ContractExecuteTransaction) GetContractID() ContractID {
	if this.contractID == nil {
		return ContractID{}
	}

	return *this.contractID
}

// SetGas sets the maximum amount of gas to use for the call.
func (this *ContractExecuteTransaction) SetGas(gas uint64) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.gas = int64(gas)
	return this
}

// GetGas returns the maximum amount of gas to use for the call.
func (this *ContractExecuteTransaction) GetGas() uint64 {
	return uint64(this.gas)
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (this *ContractExecuteTransaction) SetPayableAmount(amount Hbar) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.amount = amount.AsTinybar()
	return this
}

// GetPayableAmount returns the amount of Hbar sent (the function must be payable if this is nonzero)
func (this ContractExecuteTransaction) GetPayableAmount() Hbar {
	return HbarFromTinybar(this.amount)
}

// SetFunctionParameters sets the function parameters
func (this *ContractExecuteTransaction) SetFunctionParameters(params []byte) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.parameters = params
	return this
}

// GetFunctionParameters returns the function parameters
func (this *ContractExecuteTransaction) GetFunctionParameters() []byte {
	return this.parameters
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (this *ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParameters) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	if params == nil {
		params = NewContractFunctionParameters()
	}

	this.parameters = params._Build(&name)
	return this
}

func (this *ContractExecuteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

// Sign uses the provided privateKey to sign the transaction.
func (this *ContractExecuteTransaction) Sign(
	privateKey PrivateKey,
) *ContractExecuteTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *ContractExecuteTransaction) SignWithOperator(
	client *Client,
) (*ContractExecuteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *ContractExecuteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractExecuteTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *ContractExecuteTransaction) Freeze() (*ContractExecuteTransaction, error) {
	_, err := this.transaction.Freeze()
	return this, err
}

func (this *ContractExecuteTransaction) FreezeWith(client *Client) (*ContractExecuteTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *ContractExecuteTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *ContractExecuteTransaction) SetMaxTransactionFee(fee Hbar) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *ContractExecuteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *ContractExecuteTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ContractExecuteTransaction.
func (this *ContractExecuteTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractExecuteTransaction.
func (this *ContractExecuteTransaction) SetTransactionMemo(memo string) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *ContractExecuteTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractExecuteTransaction.
func (this *ContractExecuteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	ContractExecuteTransaction.
func (this *ContractExecuteTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractExecuteTransaction.
func (this *ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) *ContractExecuteTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractExecuteTransaction.
func (this *ContractExecuteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractExecuteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *ContractExecuteTransaction) SetMaxRetry(count int) *ContractExecuteTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *ContractExecuteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractExecuteTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *ContractExecuteTransaction) SetMaxBackoff(max time.Duration) *ContractExecuteTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *ContractExecuteTransaction) SetMinBackoff(min time.Duration) *ContractExecuteTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *ContractExecuteTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ContractExecuteTransaction:%d", timestamp.UnixNano())
}

func (this *ContractExecuteTransaction) SetLogLevel(level LogLevel) *ContractExecuteTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *ContractExecuteTransaction) getName() string {
	return "ContractExecuteTransaction"
}
func (this *ContractExecuteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.contractID != nil {
		if err := this.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *ContractExecuteTransaction) build() *services.TransactionBody {
	body := services.ContractCallTransactionBody{
		Gas:                this.gas,
		Amount:             this.amount,
		FunctionParameters: this.parameters,
	}

	if this.contractID != nil {
		body.ContractID = this.contractID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractCall{
			ContractCall: &body,
		},
	}
}

func (this *ContractExecuteTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := services.ContractCallTransactionBody{
		Gas:                this.gas,
		Amount:             this.amount,
		FunctionParameters: this.parameters,
	}

	if this.contractID != nil {
		body.ContractID = this.contractID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractCall{
			ContractCall: &body,
		},
	}, nil
}

func (this *ContractExecuteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().ContractCallMethod,
	}
}
