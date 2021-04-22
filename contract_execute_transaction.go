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

// NewContractExecuteTransaction creates a ContractExecuteTransaction builder which can be
// used to construct and execute a Contract Call Transaction.
func NewContractExecuteTransaction() ContractExecuteTransaction {
	pb := &proto.ContractCallTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractCall{ContractCall: pb}

	builder := ContractExecuteTransaction{inner, pb}

	return builder
}

func contractExecuteTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) ContractExecuteTransaction {
	return ContractExecuteTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetContractCall(),
	}
}

// SetContractID sets the contract instance to call.
func (builder ContractExecuteTransaction) SetContractID(id ContractID) ContractExecuteTransaction {
	builder.pb.ContractID = id.toProto()
	return builder
}

// SetGas sets the maximum amount of gas to use for the call.
func (builder ContractExecuteTransaction) SetGas(gas uint64) ContractExecuteTransaction {
	builder.pb.Gas = int64(gas)
	return builder
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (builder ContractExecuteTransaction) SetPayableAmount(amount Hbar) ContractExecuteTransaction {
	builder.pb.Amount = int64(amount.AsTinybar())
	return builder
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (builder ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParams) ContractExecuteTransaction {
	if params == nil {
		params = NewContractFunctionParams()
	}

	builder.pb.FunctionParameters = params.build(&name)
	return builder
}

func (builder ContractExecuteTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *ContractExecuteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_ContractCall{
			ContractCall: &proto.ContractCallTransactionBody{
				ContractID:         builder.pb.GetContractID(),
				Gas:                builder.pb.GetGas(),
				Amount:             builder.pb.GetAmount(),
				FunctionParameters: builder.pb.GetFunctionParameters(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder ContractExecuteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder ContractExecuteTransaction) SetTransactionMemo(memo string) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder ContractExecuteTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder ContractExecuteTransaction) SetNodeAccountID(nodeAccountID AccountID) ContractExecuteTransaction {
	return ContractExecuteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
