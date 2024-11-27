package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ContractExecuteTransaction calls a function of the given smart contract instance, giving it ContractFuncionParams as
// its inputs. it can use the given amount of gas, and any unspent gas will be refunded to the paying account.
//
// If tx function stores information, it is charged gas to store it. There is a fee in hbars to maintain that storage
// until the expiration time, and that fee is added as part of the transaction fee.
//
// For a cheaper but more limited _Method to call functions, see ContractCallQuery.
type ContractExecuteTransaction struct {
	*Transaction[*ContractExecuteTransaction]
	contractID *ContractID
	gas        int64
	amount     int64
	parameters []byte
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call Transaction.
func NewContractExecuteTransaction() *ContractExecuteTransaction {
	tx := &ContractExecuteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _ContractExecuteTransactionFromProtobuf(tx Transaction[*ContractExecuteTransaction], pb *services.TransactionBody) ContractExecuteTransaction {
	contractExecuteTransaction := ContractExecuteTransaction{
		contractID: _ContractIDFromProtobuf(pb.GetContractCall().GetContractID()),
		gas:        pb.GetContractCall().GetGas(),
		amount:     pb.GetContractCall().GetAmount(),
		parameters: pb.GetContractCall().GetFunctionParameters(),
	}

	tx.childTransaction = &contractExecuteTransaction
	contractExecuteTransaction.Transaction = &tx
	return contractExecuteTransaction
}

// SetContractID sets the contract instance to call.
func (tx *ContractExecuteTransaction) SetContractID(contractID ContractID) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// GetContractID returns the contract instance to call.
func (tx *ContractExecuteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// SetGas sets the maximum amount of gas to use for the call.
func (tx *ContractExecuteTransaction) SetGas(gas uint64) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.gas = int64(gas)
	return tx
}

// GetGas returns the maximum amount of gas to use for the call.
func (tx *ContractExecuteTransaction) GetGas() uint64 {
	return uint64(tx.gas)
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (tx *ContractExecuteTransaction) SetPayableAmount(amount Hbar) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.amount = amount.AsTinybar()
	return tx
}

// GetPayableAmount returns the amount of Hbar sent (the function must be payable if this is nonzero)
func (tx *ContractExecuteTransaction) GetPayableAmount() Hbar {
	return HbarFromTinybar(tx.amount)
}

// SetFunctionParameters sets the function parameters
func (tx *ContractExecuteTransaction) SetFunctionParameters(params []byte) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.parameters = params
	return tx
}

// GetFunctionParameters returns the function parameters
func (tx *ContractExecuteTransaction) GetFunctionParameters() []byte {
	return tx.parameters
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (tx *ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParameters) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	if params == nil {
		params = NewContractFunctionParameters()
	}

	tx.parameters = params._Build(&name)
	return tx
}

// ----------- Overridden functions ----------------

func (tx ContractExecuteTransaction) getName() string {
	return "ContractExecuteTransaction"
}
func (tx ContractExecuteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.contractID != nil {
		if err := tx.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx ContractExecuteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractCall{
			ContractCall: tx.buildProtoBody(),
		},
	}
}

func (tx ContractExecuteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractCall{
			ContractCall: tx.buildProtoBody(),
		},
	}, nil
}

func (tx ContractExecuteTransaction) buildProtoBody() *services.ContractCallTransactionBody {
	body := &services.ContractCallTransactionBody{
		Gas:                tx.gas,
		Amount:             tx.amount,
		FunctionParameters: tx.parameters,
	}

	if tx.contractID != nil {
		body.ContractID = tx.contractID._ToProtobuf()
	}

	return body
}

func (tx ContractExecuteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().ContractCallMethod,
	}
}

func (tx ContractExecuteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx ContractExecuteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
