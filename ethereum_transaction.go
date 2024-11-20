package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/pkg/errors"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// EthereumTransaction is used to create a EthereumTransaction transaction which can be used to construct and execute
// a Ethereum Transaction.
type EthereumTransaction struct {
	*Transaction[*EthereumTransaction]
	ethereumData  []byte
	callData      *FileID
	MaxGasAllowed int64
}

// NewEthereumTransaction creates a EthereumTransaction transaction which can be used to construct and execute
// a Ethereum Transaction.
func NewEthereumTransaction() *EthereumTransaction {
	tx := &EthereumTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _EthereumTransactionFromProtobuf(tx Transaction[*EthereumTransaction], pb *services.TransactionBody) EthereumTransaction {
	ethereumTransaction := EthereumTransaction{
		ethereumData:  pb.GetEthereumTransaction().EthereumData,
		callData:      _FileIDFromProtobuf(pb.GetEthereumTransaction().CallData),
		MaxGasAllowed: pb.GetEthereumTransaction().MaxGasAllowance,
	}

	tx.childTransaction = &ethereumTransaction
	ethereumTransaction.Transaction = &tx
	return ethereumTransaction
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
// The maximum amount, in tinybars, that the payer of the hiero transaction
// is willing to pay to complete the transaction.
func (tx *EthereumTransaction) SetMaxGasAllowed(gas int64) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.MaxGasAllowed = gas
	return tx
}

// SetMaxGasAllowanceHbar sets the maximum amount, that the payer of the hiero transaction
// is willing to pay to complete the transaction.
func (tx *EthereumTransaction) SetMaxGasAllowanceHbar(gas Hbar) *EthereumTransaction {
	tx._RequireNotFrozen()
	tx.MaxGasAllowed = gas.AsTinybar()
	return tx
}

// GetMaxGasAllowed returns the maximum amount, that the payer of the hiero transaction
// is willing to pay to complete the transaction.
func (tx *EthereumTransaction) GetMaxGasAllowed() int64 {
	return tx.MaxGasAllowed
}

// ----------- Overridden functions ----------------

func (tx EthereumTransaction) getName() string {
	return "EthereumTransaction"
}
func (tx EthereumTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx EthereumTransaction) build() *services.TransactionBody {
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

func (tx EthereumTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `EthereumTransaction`")
}

func (tx EthereumTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx EthereumTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().CallEthereum,
	}
}

func (tx EthereumTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
