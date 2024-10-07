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

// NewContractCreateFlow creates a new ContractCreateFlow transaction builder object.
func NewContractCreateFlow() *ContractCreateFlow {
	this := ContractCreateFlow{
		Transaction: _NewTransaction(),
	}

	this.SetAutoRenewPeriod(131500 * time.Minute)
	this.SetMaxTransactionFee(NewHbar(20))

	return &this
}

// SetBytecodeWithString sets the bytecode of the contract in hex-encoded string format.
func (tx *ContractCreateFlow) SetBytecodeWithString(bytecode string) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.bytecode, _ = hex.DecodeString(bytecode)
	return tx
}

// SetBytecode sets the bytecode of the contract in raw bytes.
func (tx *ContractCreateFlow) SetBytecode(bytecode []byte) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.bytecode = bytecode
	return tx
}

// GetBytecode returns the hex-encoded bytecode of the contract.
func (tx *ContractCreateFlow) GetBytecode() string {
	return hex.EncodeToString(tx.bytecode)
}

// Sets the state of the instance and its fields can be modified arbitrarily if this key signs a transaction
// to modify it. If this is null, then such modifications are not possible, and there is no administrator
// that can override the normal operation of this smart contract instance. Note that if it is created with no
// admin keys, then there is no administrator to authorize changing the admin keys, so
// there can never be any admin keys for that instance.
func (tx *ContractCreateFlow) SetAdminKey(adminKey Key) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.adminKey = &adminKey
	return tx
}

// GetAdminKey returns the admin key of the contract.
func (tx *ContractCreateFlow) GetAdminKey() Key {
	if tx.adminKey != nil {
		return *tx.adminKey
	}

	return PrivateKey{}
}

// SetGas sets the gas to run the constructor.
func (tx *ContractCreateFlow) SetGas(gas int64) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.gas = gas
	return tx
}

// GetGas returns the gas to run the constructor.
func (tx *ContractCreateFlow) GetGas() int64 {
	return tx.gas
}

// SetInitialBalance sets the initial number of hbars to put into the cryptocurrency account
// associated with and owned by the smart contract.
func (tx *ContractCreateFlow) SetInitialBalance(initialBalance Hbar) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.initialBalance = initialBalance.AsTinybar()
	return tx
}

// GetInitialBalance returns the initial number of hbars to put into the cryptocurrency account
// associated with and owned by the smart contract.
func (tx *ContractCreateFlow) GetInitialBalance() Hbar {
	return HbarFromTinybar(tx.initialBalance)
}

// SetAutoRenewPeriod sets the period that the instance will charge its account every this many seconds to renew.
func (tx *ContractCreateFlow) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &autoRenewPeriod
	return tx
}

// GetAutoRenewPeriod returns the period that the instance will charge its account every this many seconds to renew.
func (tx *ContractCreateFlow) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return *tx.autoRenewPeriod
	}

	return time.Duration(0)
}

// Deprecated
func (tx *ContractCreateFlow) SetProxyAccountID(proxyAccountID AccountID) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.proxyAccountID = &proxyAccountID
	return tx
}

// Deprecated
func (tx *ContractCreateFlow) GetProxyAccountID() AccountID {
	if tx.proxyAccountID == nil {
		return AccountID{}
	}

	return *tx.proxyAccountID
}

// Sets the constructor parameters
func (tx *ContractCreateFlow) SetConstructorParameters(params *ContractFunctionParameters) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.parameters = params._Build(nil)
	return tx
}

// Sets the constructor parameters as their raw bytes.
func (tx *ContractCreateFlow) SetConstructorParametersRaw(params []byte) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.parameters = params
	return tx
}

func (tx *ContractCreateFlow) GetConstructorParameters() []byte {
	return tx.parameters
}

// Sets the memo to be associated with this contract.
func (tx *ContractCreateFlow) SetContractMemo(memo string) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.memo = memo
	return tx
}

// Gets the memo to be associated with this contract.
func (tx *ContractCreateFlow) GetContractMemo() string {
	return tx.memo
}

// SetMaxChunks sets the maximum number of chunks that the contract bytecode can be split into.
func (tx *ContractCreateFlow) SetMaxChunks(max uint64) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.maxChunks = &max
	return tx
}

// GetMaxChunks returns the maximum number of chunks that the contract bytecode can be split into.
func (tx *ContractCreateFlow) GetMaxChunks() uint64 {
	if tx.maxChunks == nil {
		return 0
	}

	return *tx.maxChunks
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of this contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (tx *ContractCreateFlow) SetAutoRenewAccountID(id AccountID) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.autoRenewAccountID = &id
	return tx
}

// GetAutoRenewAccountID returns the account to charge for auto-renewal of this contract.
func (tx *ContractCreateFlow) GetAutoRenewAccountID() AccountID {
	if tx.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *tx.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that this contract can be automatically associated
// with (i.e., receive air-drops from).
func (tx *ContractCreateFlow) SetMaxAutomaticTokenAssociations(max int32) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.maxAutomaticTokenAssociations = max
	return tx
}

// GetMaxAutomaticTokenAssociations returns the maximum number of tokens that this
// contract can be automatically associated with.
func (tx *ContractCreateFlow) GetMaxAutomaticTokenAssociations() int32 {
	return tx.maxAutomaticTokenAssociations
}

func (tx *ContractCreateFlow) splitBytecode() *ContractCreateFlow {
	if len(tx.bytecode) > 2048 {
		tx.createBytecode = tx.bytecode[0:2048]
		tx.appendBytecode = tx.bytecode[2048:]
		return tx
	}

	tx.createBytecode = tx.bytecode
	tx.appendBytecode = []byte{}
	return tx
}

func (tx *ContractCreateFlow) _CreateFileCreateTransaction(client *Client) *FileCreateTransaction {
	if client == nil {
		return &FileCreateTransaction{}
	}
	fileCreateTx := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(tx.createBytecode)

	if len(tx.nodeAccountIDs) > 0 {
		fileCreateTx.SetNodeAccountIDs(tx.nodeAccountIDs)
	}

	return fileCreateTx
}

func (tx *ContractCreateFlow) _CreateFileAppendTransaction(fileID FileID) *FileAppendTransaction {
	fileAppendTx := NewFileAppendTransaction().
		SetFileID(fileID).
		SetContents(tx.appendBytecode)

	if len(tx.nodeAccountIDs) > 0 {
		fileAppendTx.SetNodeAccountIDs(tx.nodeAccountIDs)
	}
	if tx.maxChunks != nil {
		fileAppendTx.SetMaxChunks(*tx.maxChunks)
	}

	return fileAppendTx
}

func (tx *ContractCreateFlow) _CreateContractCreateTransaction(fileID FileID) *ContractCreateTransaction {
	contractCreateTx := NewContractCreateTransaction().
		SetGas(uint64(tx.gas)).
		SetConstructorParametersRaw(tx.parameters).
		SetInitialBalance(HbarFromTinybar(tx.initialBalance)).
		SetBytecodeFileID(fileID).
		SetContractMemo(tx.memo)

	if len(tx.nodeAccountIDs) > 0 {
		contractCreateTx.SetNodeAccountIDs(tx.nodeAccountIDs)
	}

	if tx.adminKey != nil {
		contractCreateTx.SetAdminKey(*tx.adminKey)
	}

	if tx.proxyAccountID != nil {
		contractCreateTx.SetProxyAccountID(*tx.proxyAccountID)
	}

	if tx.autoRenewPeriod != nil {
		contractCreateTx.SetAutoRenewPeriod(*tx.autoRenewPeriod)
	}

	if tx.autoRenewAccountID != nil {
		contractCreateTx.SetAutoRenewAccountID(*tx.autoRenewAccountID)
	}

	if tx.maxAutomaticTokenAssociations != 0 {
		contractCreateTx.SetMaxAutomaticTokenAssociations(tx.maxAutomaticTokenAssociations)
	}

	return contractCreateTx
}

func (tx *ContractCreateFlow) Execute(client *Client) (TransactionResponse, error) {
	tx.splitBytecode()

	fileCreateResponse, err := tx._CreateFileCreateTransaction(client).Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	fileCreateReceipt, err := fileCreateResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	if fileCreateReceipt.FileID == nil {
		return TransactionResponse{}, errors.New("fileID is nil")
	}
	fileID := *fileCreateReceipt.FileID
	if len(tx.appendBytecode) > 0 {
		fileAppendResponse, err := tx._CreateFileAppendTransaction(fileID).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}

		_, err = fileAppendResponse.SetValidateStatus(true).GetReceipt(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}
	contractCreateResponse, err := tx._CreateContractCreateTransaction(fileID).Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	_, err = contractCreateResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	return contractCreateResponse, nil
}

// SetNodeAccountIDs sets the node AccountID for this ContractCreateFlow.
func (tx *ContractCreateFlow) SetNodeAccountIDs(nodeID []AccountID) *ContractCreateFlow {
	tx._RequireNotFrozen()
	tx.nodeAccountIDs = nodeID
	return tx
}

// GetNodeAccountIDs returns the node AccountID for this ContractCreateFlow.
func (tx *ContractCreateFlow) GetNodeAccountIDs() []AccountID {
	return tx.nodeAccountIDs
}
