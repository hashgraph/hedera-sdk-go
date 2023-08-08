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
	transaction := EthereumTransaction{
		Transaction: _NewTransaction(),
	}

	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _EthereumTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *EthereumTransaction {
	return &EthereumTransaction{
		Transaction:   transaction,
		ethereumData:  pb.GetEthereumTransaction().EthereumData,
		callData:      _FileIDFromProtobuf(pb.GetEthereumTransaction().CallData),
		MaxGasAllowed: pb.GetEthereumTransaction().MaxGasAllowance,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *EthereumTransaction) SetGrpcDeadline(deadline *time.Duration) *EthereumTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetEthereumData
// The raw Ethereum transaction (RLP encoded type 0, 1, and 2). Complete
// unless the callData field is set.
func (transaction *EthereumTransaction) SetEthereumData(data []byte) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.ethereumData = data
	return transaction
}

// GetEthereumData returns the raw Ethereum transaction (RLP encoded type 0, 1, and 2).
func (transaction *EthereumTransaction) GetEthereumData() []byte {
	return transaction.ethereumData
}

// Deprecated
func (transaction *EthereumTransaction) SetCallData(file FileID) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.callData = &file
	return transaction
}

// SetCallDataFileID sets the file ID containing the call data.
func (transaction *EthereumTransaction) SetCallDataFileID(file FileID) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.callData = &file
	return transaction
}

// GetCallData
// For large transactions (for example contract create) this is the callData
// of the ethereumData. The data in the ethereumData will be re-written with
// the callData element as a zero length string with the original contents in
// the referenced file at time of execution. The ethereumData will need to be
// "rehydrated" with the callData for signature validation to pass.
func (transaction *EthereumTransaction) GetCallData() FileID {
	if transaction.callData != nil {
		return *transaction.callData
	}

	return FileID{}
}

// SetMaxGasAllowed
// The maximum amount, in tinybars, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (transaction *EthereumTransaction) SetMaxGasAllowed(gas int64) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.MaxGasAllowed = gas
	return transaction
}

// SetMaxGasAllowanceHbar sets the maximum amount, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (transaction *EthereumTransaction) SetMaxGasAllowanceHbar(gas Hbar) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.MaxGasAllowed = gas.AsTinybar()
	return transaction
}

// GetMaxGasAllowed returns the maximum amount, that the payer of the hedera transaction
// is willing to pay to complete the transaction.
func (transaction *EthereumTransaction) GetMaxGasAllowed() int64 {
	return transaction.MaxGasAllowed
}

func (transaction *EthereumTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.callData != nil {
		if err := transaction.callData.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *EthereumTransaction) _Build() *services.TransactionBody {
	body := &services.EthereumTransactionBody{
		EthereumData:    transaction.ethereumData,
		MaxGasAllowance: transaction.MaxGasAllowed,
	}

	if transaction.callData != nil {
		body.CallData = transaction.callData._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionID:            transaction.transactionID._ToProtobuf(),
		TransactionFee:           transaction.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		Memo:                     transaction.Transaction.memo,
		Data: &services.TransactionBody_EthereumTransaction{
			EthereumTransaction: body,
		},
	}
}

func (transaction *EthereumTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `EthereumTransaction")
}

func _EthereumTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().CallEthereum,
	}
}

func (transaction *EthereumTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *EthereumTransaction) Sign(
	privateKey PrivateKey,
) *EthereumTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *EthereumTransaction) SignWithOperator(
	client *Client,
) (*EthereumTransaction, error) {
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
func (transaction *EthereumTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *EthereumTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *EthereumTransaction) Execute(
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

	if transaction.grpcDeadline == nil {
		transaction.grpcDeadline = client.requestTimeout
	}

	resp, err := _Execute(
		client,
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_EthereumTransactionGetMethod,
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

func (transaction *EthereumTransaction) Freeze() (*EthereumTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *EthereumTransaction) FreezeWith(client *Client) (*EthereumTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	body := transaction._Build()
	if err != nil {
		return &EthereumTransaction{}, err
	}

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *EthereumTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *EthereumTransaction) SetMaxTransactionFee(fee Hbar) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *EthereumTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *EthereumTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this EthereumTransaction.
func (transaction *EthereumTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this EthereumTransaction.
func (transaction *EthereumTransaction) SetTransactionMemo(memo string) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *EthereumTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this EthereumTransaction.
func (transaction *EthereumTransaction) SetTransactionValidDuration(duration time.Duration) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	EthereumTransaction.
func (transaction *EthereumTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this EthereumTransaction.
func (transaction *EthereumTransaction) SetTransactionID(transactionID TransactionID) *EthereumTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this EthereumTransaction.
func (transaction *EthereumTransaction) SetNodeAccountIDs(nodeID []AccountID) *EthereumTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *EthereumTransaction) SetMaxRetry(count int) *EthereumTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *EthereumTransaction) AddSignature(publicKey PublicKey, signature []byte) *EthereumTransaction {
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
func (transaction *EthereumTransaction) SetMaxBackoff(max time.Duration) *EthereumTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *EthereumTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *EthereumTransaction) SetMinBackoff(min time.Duration) *EthereumTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *EthereumTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *EthereumTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("EthereumTransaction:%d", timestamp.UnixNano())
}
