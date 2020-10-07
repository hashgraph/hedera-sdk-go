package hedera

import (
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
	Transaction
	pb *proto.ContractCallTransactionBody
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call Transaction.
func NewContractExecuteTransaction() *ContractExecuteTransaction {
	pb := &proto.ContractCallTransactionBody{}

	transaction := ContractExecuteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

// SetContractID sets the contract instance to call.
func (transaction *ContractExecuteTransaction) SetContractID(id ContractID) *ContractExecuteTransaction {
	transaction.pb.ContractID = id.toProto()
	return transaction
}

func (transaction ContractExecuteTransaction) GetContractID(id ContractID)  ContractID{
	return contractIDFromProto(transaction.pb.GetContractID())
}

// SetGas sets the maximum amount of gas to use for the call.
func (transaction *ContractExecuteTransaction) SetGas(gas uint64) *ContractExecuteTransaction {
	transaction.pb.Gas = int64(gas)
	return transaction
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (transaction *ContractExecuteTransaction) SetPayableAmount(amount Hbar) *ContractExecuteTransaction {
	transaction.pb.Amount = amount.AsTinybar()
	return transaction
}

func (transaction ContractExecuteTransaction) GetPayableAmount()  uint64{
	return uint64(transaction.pb.Gas)
}

func (transaction *ContractExecuteTransaction) SetFunctionParameters(params []byte) *ContractExecuteTransaction {
	transaction.pb.FunctionParameters = params
	return transaction
}

func (transaction *ContractExecuteTransaction) GetFunctionParameters() []byte {
	return transaction.pb.GetFunctionParameters()
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
//func (transaction *ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParameters) *ContractExecuteTransaction {
//	if params == nil {
//		params = NewContractFunctionParams()
//	}
//
//	transaction.pb.FunctionParameters = params.build(&name)
//	return transaction
//}
