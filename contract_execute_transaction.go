package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractExecuteTransaction struct {
	TransactionBuilder
	pb *proto.ContractCallTransactionBody
}

func NewContractExecuteTransaction() ContractExecuteTransaction {
	pb := &proto.ContractCallTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractCall{ContractCall: pb}

	builder := ContractExecuteTransaction{inner, pb}

	return builder
}

func (builder ContractExecuteTransaction) SetContractID(id ContractID) ContractExecuteTransaction {
	builder.pb.ContractID = id.toProto()
	return builder
}

func (builder ContractExecuteTransaction) SetGas(gas uint64) ContractExecuteTransaction {
	builder.pb.Gas = int64(gas)
	return builder
}

func (builder ContractExecuteTransaction) SetPayableAmount(amount Hbar) ContractExecuteTransaction {
	builder.pb.Amount = int64(amount.AsTinybar())
	return builder
}

func (builder ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParams) ContractExecuteTransaction {
	if params == nil {
		params = NewContractFunctionParams()
	}

	builder.pb.FunctionParameters = params.build(&name)
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ContractExecuteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder ContractExecuteTransaction) SetTransactionMemo(memo string) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder ContractExecuteTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ContractExecuteTransaction) SetNodeAccountID(nodeAccountID AccountID) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
