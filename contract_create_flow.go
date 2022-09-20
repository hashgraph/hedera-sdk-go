package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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
	"encoding/hex"
	"time"

	"github.com/pkg/errors"
)

type ContractCreateFlow struct {
	Transaction
	bytecode                      []byte
	proxyAccountID                *AccountID
	adminKey                      *Key
	gas                           int64
	initialBalance                int64
	autoRenewPeriod               *time.Duration
	parameters                    []byte
	nodeAccountIDs                []AccountID
	createBytecode                []byte
	appendBytecode                []byte
	autoRenewAccountID            *AccountID
	maxAutomaticTokenAssociations int32
	maxChunks                     *uint64
}

func NewContractCreateFlow() *ContractCreateFlow {
	this := ContractCreateFlow{
		Transaction: _NewTransaction(),
	}

	this.SetAutoRenewPeriod(131500 * time.Minute)
	this.SetMaxTransactionFee(NewHbar(20))

	return &this
}

func (this *ContractCreateFlow) SetBytecodeWithString(bytecode string) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.bytecode, _ = hex.DecodeString(bytecode)
	return this
}

func (this *ContractCreateFlow) SetBytecode(bytecode []byte) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.bytecode = bytecode
	return this
}

func (this *ContractCreateFlow) GetBytecode() string {
	return hex.EncodeToString(this.bytecode)
}

func (this *ContractCreateFlow) SetAdminKey(adminKey Key) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.adminKey = &adminKey
	return this
}

func (this *ContractCreateFlow) GetAdminKey() Key {
	if this.adminKey != nil {
		return *this.adminKey
	}

	return PrivateKey{}
}

func (this *ContractCreateFlow) SetGas(gas int64) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.gas = gas
	return this
}

func (this *ContractCreateFlow) GetGas() int64 {
	return this.gas
}

func (this *ContractCreateFlow) SetInitialBalance(initialBalance Hbar) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.initialBalance = initialBalance.AsTinybar()
	return this
}

func (this *ContractCreateFlow) GetInitialBalance() Hbar {
	return HbarFromTinybar(this.initialBalance)
}

func (this *ContractCreateFlow) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.autoRenewPeriod = &autoRenewPeriod
	return this
}

func (this *ContractCreateFlow) GetAutoRenewPeriod() time.Duration {
	if this.autoRenewPeriod != nil {
		return *this.autoRenewPeriod
	}

	return time.Duration(0)
}

// Deprecated
func (this *ContractCreateFlow) SetProxyAccountID(proxyAccountID AccountID) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.proxyAccountID = &proxyAccountID
	return this
}

// Deprecated
func (this *ContractCreateFlow) GetProxyAccountID() AccountID {
	if this.proxyAccountID == nil {
		return AccountID{}
	}

	return *this.proxyAccountID
}

// Sets the constructor parameters
func (this *ContractCreateFlow) SetConstructorParameters(params *ContractFunctionParameters) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.parameters = params._Build(nil)
	return this
}

// Sets the constructor parameters as their raw bytes.
func (this *ContractCreateFlow) SetConstructorParametersRaw(params []byte) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.parameters = params
	return this
}

func (this *ContractCreateFlow) GetConstructorParameters() []byte {
	return this.parameters
}

// Sets the memo to be associated with this contract.
func (this *ContractCreateFlow) SetContractMemo(memo string) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.memo = memo
	return this
}

func (this *ContractCreateFlow) GetContractMemo() string {
	return this.memo
}

func (this *ContractCreateFlow) SetMaxChunks(max uint64) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.maxChunks = &max
	return this
}

func (this *ContractCreateFlow) GetMaxChunks() uint64 {
	if this.maxChunks == nil {
		return 0
	}

	return *this.maxChunks
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of this contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (this *ContractCreateFlow) SetAutoRenewAccountID(id AccountID) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.autoRenewAccountID = &id
	return this
}

func (this *ContractCreateFlow) GetAutoRenewAccountID() AccountID {
	if this.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *this.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that this contract can be automatically associated
// with (i.e., receive air-drops from).
func (this *ContractCreateFlow) SetMaxAutomaticTokenAssociations(max int32) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.maxAutomaticTokenAssociations = max
	return this
}

func (this *ContractCreateFlow) GetMaxAutomaticTokenAssociations() int32 {
	return this.maxAutomaticTokenAssociations
}

func (this *ContractCreateFlow) _SplitBytecode() *ContractCreateFlow {
	if len(this.bytecode) > 2048 {
		this.createBytecode = this.bytecode[0:2048]
		this.appendBytecode = this.bytecode[2048:]
		return this
	}

	this.createBytecode = this.bytecode
	this.appendBytecode = []byte{}
	return this
}

func (this *ContractCreateFlow) _CreateFileCreateTransaction(client *Client) *FileCreateTransaction {
	if client == nil {
		return &FileCreateTransaction{}
	}
	fileCreateTx := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(this.createBytecode)

	if len(this.nodeAccountIDs) > 0 {
		fileCreateTx.SetNodeAccountIDs(this.nodeAccountIDs)
	}

	return fileCreateTx
}

func (this *ContractCreateFlow) _CreateFileAppendTransaction(fileID FileID) *FileAppendTransaction {
	fileAppendTx := NewFileAppendTransaction().
		SetFileID(fileID).
		SetContents(this.appendBytecode)

	if len(this.nodeAccountIDs) > 0 {
		fileAppendTx.SetNodeAccountIDs(this.nodeAccountIDs)
	}
	if this.maxChunks != nil {
		fileAppendTx.SetMaxChunks(*this.maxChunks)
	}

	return fileAppendTx
}

func (this *ContractCreateFlow) _CreateContractCreateTransaction(fileID FileID) *ContractCreateTransaction {
	contractCreateTx := NewContractCreateTransaction().
		SetGas(uint64(this.gas)).
		SetConstructorParametersRaw(this.parameters).
		SetInitialBalance(HbarFromTinybar(this.initialBalance)).
		SetBytecodeFileID(fileID).
		SetContractMemo(this.memo)

	if len(this.nodeAccountIDs) > 0 {
		contractCreateTx.SetNodeAccountIDs(this.nodeAccountIDs)
	}

	if this.adminKey != nil {
		contractCreateTx.SetAdminKey(*this.adminKey)
	}

	if this.proxyAccountID != nil {
		contractCreateTx.SetProxyAccountID(*this.proxyAccountID)
	}

	if this.autoRenewPeriod != nil {
		contractCreateTx.SetAutoRenewPeriod(*this.autoRenewPeriod)
	}

	if this.autoRenewAccountID != nil {
		contractCreateTx.SetAutoRenewAccountID(*this.autoRenewAccountID)
	}

	if this.maxAutomaticTokenAssociations != 0 {
		contractCreateTx.SetMaxAutomaticTokenAssociations(this.maxAutomaticTokenAssociations)
	}

	return contractCreateTx
}

func (this *ContractCreateFlow) _CreateContractCreateTransactionWithBytecode() *ContractCreateTransaction {
	contractCreateTx := NewContractCreateTransaction().
		SetGas(uint64(this.gas)).
		SetConstructorParametersRaw(this.parameters).
		SetInitialBalance(HbarFromTinybar(this.initialBalance)).
		SetBytecode(this.createBytecode).
		SetContractMemo(this.memo)

	if len(this.nodeAccountIDs) > 0 {
		contractCreateTx.SetNodeAccountIDs(this.nodeAccountIDs)
	}

	if this.adminKey != nil {
		contractCreateTx.SetAdminKey(*this.adminKey)
	}

	if this.proxyAccountID != nil {
		contractCreateTx.SetProxyAccountID(*this.proxyAccountID)
	}

	if this.autoRenewPeriod != nil {
		contractCreateTx.SetAutoRenewPeriod(*this.autoRenewPeriod)
	}

	if this.autoRenewAccountID != nil {
		contractCreateTx.SetAutoRenewAccountID(*this.autoRenewAccountID)
	}

	if this.maxAutomaticTokenAssociations != 0 {
		contractCreateTx.SetMaxAutomaticTokenAssociations(this.maxAutomaticTokenAssociations)
	}

	return contractCreateTx
}

func (this *ContractCreateFlow) _CreateTransactionReceiptQuery(response TransactionResponse) *TransactionReceiptQuery {
	return NewTransactionReceiptQuery().
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		SetTransactionID(response.TransactionID)
}

func (this *ContractCreateFlow) Execute(client *Client) (TransactionResponse, error) {
	this._SplitBytecode()

	if len(this.appendBytecode) > 0 {
		fileCreateResponse, err := this._CreateFileCreateTransaction(client).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}
		fileCreateReceipt, err := this._CreateTransactionReceiptQuery(fileCreateResponse).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}
		if fileCreateReceipt.FileID == nil {
			return TransactionResponse{}, errors.New("fileID is nil")
		}
		fileID := *fileCreateReceipt.FileID

		fileAppendResponse, err := this._CreateFileAppendTransaction(fileID).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}

		_, err = this._CreateTransactionReceiptQuery(fileAppendResponse).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}

		contractCreateResponse, err := this._CreateContractCreateTransaction(fileID).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}
		_, err = this._CreateTransactionReceiptQuery(contractCreateResponse).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}

		return contractCreateResponse, nil
	}

	contractCreateResponse, err := this._CreateContractCreateTransactionWithBytecode().Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	_, err = this._CreateTransactionReceiptQuery(contractCreateResponse).
		Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	return contractCreateResponse, nil
}

func (this *ContractCreateFlow) SetNodeAccountIDs(nodeID []AccountID) *ContractCreateFlow {
	this._RequireNotFrozen()
	this.nodeAccountIDs = nodeID
	return this
}

func (this *ContractCreateFlow) GetNodeAccountIDs() []AccountID {
	return this.nodeAccountIDs
}
