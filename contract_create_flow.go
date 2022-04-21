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
}

func NewContractCreateFlow() *ContractCreateFlow {
	transaction := ContractCreateFlow{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(131500 * time.Minute)
	transaction.SetMaxTransactionFee(NewHbar(20))

	return &transaction
}

func (transaction *ContractCreateFlow) SetBytecodeWithString(bytecode string) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.bytecode, _ = hex.DecodeString(bytecode)
	return transaction
}

func (transaction *ContractCreateFlow) SetBytecode(bytecode []byte) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.bytecode = bytecode
	return transaction
}

func (transaction *ContractCreateFlow) GetBytecode() string {
	return hex.EncodeToString(transaction.bytecode)
}

func (transaction *ContractCreateFlow) SetAdminKey(adminKey Key) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.adminKey = &adminKey
	return transaction
}

func (transaction *ContractCreateFlow) GetAdminKey() Key {
	if transaction.adminKey != nil {
		return *transaction.adminKey
	}

	return PrivateKey{}
}

func (transaction *ContractCreateFlow) SetGas(gas int64) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.gas = gas
	return transaction
}

func (transaction *ContractCreateFlow) GetGas() int64 {
	return transaction.gas
}

func (transaction *ContractCreateFlow) SetInitialBalance(initialBalance Hbar) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.initialBalance = initialBalance.AsTinybar()
	return transaction
}

func (transaction *ContractCreateFlow) GetInitialBalance() Hbar {
	return HbarFromTinybar(transaction.initialBalance)
}

func (transaction *ContractCreateFlow) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (transaction *ContractCreateFlow) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

func (transaction *ContractCreateFlow) SetProxyAccountID(proxyAccountID AccountID) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.proxyAccountID = &proxyAccountID
	return transaction
}

func (transaction *ContractCreateFlow) GetProxyAccountID() AccountID {
	if transaction.proxyAccountID == nil {
		return AccountID{}
	}

	return *transaction.proxyAccountID
}

// Sets the constructor parameters
func (transaction *ContractCreateFlow) SetConstructorParameters(params *ContractFunctionParameters) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.parameters = params._Build(nil)
	return transaction
}

// Sets the constructor parameters as their raw bytes.
func (transaction *ContractCreateFlow) SetConstructorParametersRaw(params []byte) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.parameters = params
	return transaction
}

func (transaction *ContractCreateFlow) GetConstructorParameters() []byte {
	return transaction.parameters
}

// Sets the memo to be associated with this contract.
func (transaction *ContractCreateFlow) SetContractMemo(memo string) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	return transaction
}

func (transaction *ContractCreateFlow) GetContractMemo() string {
	return transaction.memo
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of this contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (transaction *ContractCreateFlow) SetAutoRenewAccountID(id AccountID) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.autoRenewAccountID = &id
	return transaction
}

func (transaction *ContractCreateFlow) GetAutoRenewAccountID() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that this contract can be automatically associated
// with (i.e., receive air-drops from).
func (transaction *ContractCreateFlow) SetMaxAutomaticTokenAssociations(max int32) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.maxAutomaticTokenAssociations = max
	return transaction
}

func (transaction *ContractCreateFlow) GetMaxAutomaticTokenAssociations() int32 {
	return transaction.maxAutomaticTokenAssociations
}

func (transaction *ContractCreateFlow) _SplitBytecode() *ContractCreateFlow {
	if len(transaction.bytecode) > 2048 {
		transaction.createBytecode = transaction.bytecode[0:2048]
		transaction.appendBytecode = transaction.bytecode[2048:]
		return transaction
	}

	transaction.createBytecode = transaction.bytecode
	transaction.appendBytecode = []byte{}
	return transaction
}

func (transaction *ContractCreateFlow) _CreateFileCreateTransaction(client *Client) *FileCreateTransaction {
	if client == nil {
		return &FileCreateTransaction{}
	}
	fileCreateTx := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(transaction.createBytecode)

	if len(transaction.nodeAccountIDs) > 0 {
		fileCreateTx.SetNodeAccountIDs(transaction.nodeAccountIDs)
	}

	return fileCreateTx
}

func (transaction *ContractCreateFlow) _CreateFileAppendTransaction(fileID FileID) *FileAppendTransaction {
	fileAppendTx := NewFileAppendTransaction().
		SetFileID(fileID).
		SetContents(transaction.appendBytecode)

	if len(transaction.nodeAccountIDs) > 0 {
		fileAppendTx.SetNodeAccountIDs(transaction.nodeAccountIDs)
	}

	return fileAppendTx
}

func (transaction *ContractCreateFlow) _CreateContractCreateTransaction(fileID FileID) *ContractCreateTransaction {
	contractCreateTx := NewContractCreateTransaction().
		SetGas(uint64(transaction.gas)).
		SetConstructorParametersRaw(transaction.parameters).
		SetInitialBalance(HbarFromTinybar(transaction.initialBalance)).
		SetBytecodeFileID(fileID).
		SetContractMemo(transaction.memo)

	if len(transaction.nodeAccountIDs) > 0 {
		contractCreateTx.SetNodeAccountIDs(transaction.nodeAccountIDs)
	}

	if transaction.adminKey != nil {
		contractCreateTx.SetAdminKey(*transaction.adminKey)
	}

	if transaction.proxyAccountID != nil {
		contractCreateTx.SetProxyAccountID(*transaction.proxyAccountID)
	}

	if transaction.autoRenewPeriod != nil {
		contractCreateTx.SetAutoRenewPeriod(*transaction.autoRenewPeriod)
	}

	if transaction.autoRenewAccountID != nil {
		contractCreateTx.SetAutoRenewAccountID(*transaction.autoRenewAccountID)
	}

	if transaction.maxAutomaticTokenAssociations != 0 {
		contractCreateTx.SetMaxAutomaticTokenAssociations(transaction.maxAutomaticTokenAssociations)
	}

	return contractCreateTx
}

func (transaction *ContractCreateFlow) _CreateContractCreateTransactionWithBytecode() *ContractCreateTransaction {
	contractCreateTx := NewContractCreateTransaction().
		SetGas(uint64(transaction.gas)).
		SetConstructorParametersRaw(transaction.parameters).
		SetInitialBalance(HbarFromTinybar(transaction.initialBalance)).
		SetBytecode(transaction.createBytecode).
		SetContractMemo(transaction.memo)

	if len(transaction.nodeAccountIDs) > 0 {
		contractCreateTx.SetNodeAccountIDs(transaction.nodeAccountIDs)
	}

	if transaction.adminKey != nil {
		contractCreateTx.SetAdminKey(*transaction.adminKey)
	}

	if transaction.proxyAccountID != nil {
		contractCreateTx.SetProxyAccountID(*transaction.proxyAccountID)
	}

	if transaction.autoRenewPeriod != nil {
		contractCreateTx.SetAutoRenewPeriod(*transaction.autoRenewPeriod)
	}

	if transaction.autoRenewAccountID != nil {
		contractCreateTx.SetAutoRenewAccountID(*transaction.autoRenewAccountID)
	}

	if transaction.maxAutomaticTokenAssociations != 0 {
		contractCreateTx.SetMaxAutomaticTokenAssociations(transaction.maxAutomaticTokenAssociations)
	}

	return contractCreateTx
}

func (transaction *ContractCreateFlow) _CreateContractCreateTransactionWithInitcode() *ContractCreateTransaction {
	contractCreateTx := NewContractCreateTransaction().
		SetGas(uint64(transaction.gas)).
		SetConstructorParametersRaw(transaction.parameters).
		SetInitialBalance(HbarFromTinybar(transaction.initialBalance)).
		SetInitcode(transaction.createBytecode).
		SetContractMemo(transaction.memo)

	if len(transaction.nodeAccountIDs) > 0 {
		contractCreateTx.SetNodeAccountIDs(transaction.nodeAccountIDs)
	}

	if transaction.adminKey != nil {
		contractCreateTx.SetAdminKey(*transaction.adminKey)
	}

	if transaction.proxyAccountID != nil {
		contractCreateTx.SetProxyAccountID(*transaction.proxyAccountID)
	}

	if transaction.autoRenewPeriod != nil {
		contractCreateTx.SetAutoRenewPeriod(*transaction.autoRenewPeriod)
	}

	return contractCreateTx
}

func (transaction *ContractCreateFlow) _CreateTransactionReceiptQuery(response TransactionResponse) *TransactionReceiptQuery {
	return NewTransactionReceiptQuery().
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		SetTransactionID(response.TransactionID)
}

func (transaction *ContractCreateFlow) Execute(client *Client) (TransactionResponse, error) {
	transaction._SplitBytecode()

	if len(transaction.appendBytecode) > 0 {
		fileCreateResponse, err := transaction._CreateFileCreateTransaction(client).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}
		fileCreateReceipt, err := transaction._CreateTransactionReceiptQuery(fileCreateResponse).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}
		if fileCreateReceipt.FileID == nil {
			return TransactionResponse{}, errors.New("fileID is nil")
		}
		fileID := *fileCreateReceipt.FileID

		fileAppendResponse, err := transaction._CreateFileAppendTransaction(fileID).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}
		_, err = transaction._CreateTransactionReceiptQuery(fileAppendResponse).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}

		contractCreateResponse, err := transaction._CreateContractCreateTransaction(fileID).Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}
		_, err = transaction._CreateTransactionReceiptQuery(contractCreateResponse).
			Execute(client)
		if err != nil {
			return TransactionResponse{}, err
		}

		return contractCreateResponse, nil
	}

	contractCreateResponse, err := transaction._CreateContractCreateTransactionWithBytecode().Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	_, err = transaction._CreateTransactionReceiptQuery(contractCreateResponse).
		Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	return contractCreateResponse, nil
}

func (transaction *ContractCreateFlow) SetNodeAccountIDs(nodeID []AccountID) *ContractCreateFlow {
	transaction._RequireNotFrozen()
	transaction.nodeAccountIDs = nodeID
	return transaction
}

func (transaction *ContractCreateFlow) GetNodeAccountIDs() []AccountID {
	return transaction.nodeAccountIDs
}
