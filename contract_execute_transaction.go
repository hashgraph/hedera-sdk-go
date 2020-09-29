package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ContractExecuteTransaction calls a function of the given smart contract instance, giving it ContractFuncionParams as
// its inputs. it can use the given amount of gas, and any unspent gas will be refunded to the paying account.
//
// If this function stores information, it is charged gas to store it. There is a fee in hbars to maintain that storage
// until the expiration time, and that fee is added as part of the transaction fee.
//
// For a cheaper but more limited method to call functions, see ContractCallQuery.
type ContractExecuteTransaction struct {
	TransactionBuilder
	pb *proto.ContractCallTransactionBody
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call Transaction.
func NewContractExecuteTransaction() ContractExecuteTransaction {
	pb := &proto.ContractCallTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractCall{ContractCall: pb}

	transaction := ContractExecuteTransaction{inner, pb}

	return transaction
}

// SetContractID sets the contract instance to call.
func (transaction ContractExecuteTransaction) SetContractID(id ContractID) ContractExecuteTransaction {
	transaction.pb.ContractID = id.toProto()
	return transaction
}

// SetGas sets the maximum amount of gas to use for the call.
func (transaction ContractExecuteTransaction) SetGas(gas uint64) ContractExecuteTransaction {
	transaction.pb.Gas = int64(gas)
	return transaction
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (transaction ContractExecuteTransaction) SetPayableAmount(amount Hbar) ContractExecuteTransaction {
	transaction.pb.Amount = int64(amount.AsTinybar())
	return transaction
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (transaction ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParams) ContractExecuteTransaction {
	if params == nil {
		params = NewContractFunctionParams()
	}

	transaction.pb.FunctionParameters = params.build(&name)
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction ContractExecuteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ContractExecuteTransaction {
	return ContractExecuteTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction ContractExecuteTransaction) SetTransactionMemo(memo string) ContractExecuteTransaction {
	return ContractExecuteTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction ContractExecuteTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractExecuteTransaction {
	return ContractExecuteTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) ContractExecuteTransaction {
	return ContractExecuteTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction ContractExecuteTransaction) SetNodeID(nodeAccountID AccountID) ContractExecuteTransaction {
	return ContractExecuteTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
