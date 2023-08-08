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
	Transaction
	contractID *ContractID
	gas        int64
	amount     int64
	parameters []byte
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call Transaction.
func NewContractExecuteTransaction() *ContractExecuteTransaction {
	transaction := ContractExecuteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _ContractExecuteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *ContractExecuteTransaction {
	return &ContractExecuteTransaction{
		Transaction: transaction,
		contractID:  _ContractIDFromProtobuf(pb.GetContractCall().GetContractID()),
		gas:         pb.GetContractCall().GetGas(),
		amount:      pb.GetContractCall().GetAmount(),
		parameters:  pb.GetContractCall().GetFunctionParameters(),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *ContractExecuteTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractExecuteTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetContractID sets the contract instance to call.
func (transaction *ContractExecuteTransaction) SetContractID(contractID ContractID) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.contractID = &contractID
	return transaction
}

// GetContractID returns the contract instance to call.
func (transaction *ContractExecuteTransaction) GetContractID() ContractID {
	if transaction.contractID == nil {
		return ContractID{}
	}

	return *transaction.contractID
}

// SetGas sets the maximum amount of gas to use for the call.
func (transaction *ContractExecuteTransaction) SetGas(gas uint64) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.gas = int64(gas)
	return transaction
}

// GetGas returns the maximum amount of gas to use for the call.
func (transaction *ContractExecuteTransaction) GetGas() uint64 {
	return uint64(transaction.gas)
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (transaction *ContractExecuteTransaction) SetPayableAmount(amount Hbar) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.amount = amount.AsTinybar()
	return transaction
}

// GetPayableAmount returns the amount of Hbar sent (the function must be payable if this is nonzero)
func (transaction ContractExecuteTransaction) GetPayableAmount() Hbar {
	return HbarFromTinybar(transaction.amount)
}

// SetFunctionParameters sets the function parameters
func (transaction *ContractExecuteTransaction) SetFunctionParameters(params []byte) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.parameters = params
	return transaction
}

// GetFunctionParameters returns the function parameters
func (transaction *ContractExecuteTransaction) GetFunctionParameters() []byte {
	return transaction.parameters
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (transaction *ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParameters) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	if params == nil {
		params = NewContractFunctionParameters()
	}

	transaction.parameters = params._Build(&name)
	return transaction
}

func (transaction *ContractExecuteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.contractID != nil {
		if err := transaction.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *ContractExecuteTransaction) _Build() *services.TransactionBody {
	body := services.ContractCallTransactionBody{
		Gas:                transaction.gas,
		Amount:             transaction.amount,
		FunctionParameters: transaction.parameters,
	}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractCall{
			ContractCall: &body,
		},
	}
}

func (transaction *ContractExecuteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractExecuteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := services.ContractCallTransactionBody{
		Gas:                transaction.gas,
		Amount:             transaction.amount,
		FunctionParameters: transaction.parameters,
	}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractCall{
			ContractCall: &body,
		},
	}, nil
}

func _ContractExecuteTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().ContractCallMethod,
	}
}

func (transaction *ContractExecuteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractExecuteTransaction) Sign(
	privateKey PrivateKey,
) *ContractExecuteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *ContractExecuteTransaction) SignWithOperator(
	client *Client,
) (*ContractExecuteTransaction, error) {
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
func (transaction *ContractExecuteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractExecuteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractExecuteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
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
		_ContractExecuteTransactionGetMethod,
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

func (transaction *ContractExecuteTransaction) Freeze() (*ContractExecuteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractExecuteTransaction) FreezeWith(client *Client) (*ContractExecuteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ContractExecuteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ContractExecuteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ContractExecuteTransaction) SetMaxTransactionFee(fee Hbar) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *ContractExecuteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *ContractExecuteTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionMemo(memo string) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *ContractExecuteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractExecuteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *ContractExecuteTransaction) SetMaxRetry(count int) *ContractExecuteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *ContractExecuteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractExecuteTransaction {
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
func (transaction *ContractExecuteTransaction) SetMaxBackoff(max time.Duration) *ContractExecuteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *ContractExecuteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *ContractExecuteTransaction) SetMinBackoff(min time.Duration) *ContractExecuteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *ContractExecuteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *ContractExecuteTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ContractExecuteTransaction:%d", timestamp.UnixNano())
}

func (transaction *ContractExecuteTransaction) SetLogLevel(level LogLevel) *ContractExecuteTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
