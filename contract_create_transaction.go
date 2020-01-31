package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ContractCreateTransaction struct {
	TransactionBuilder
	pb *proto.ContractCreateTransactionBody
}

func NewContractCreateTransaction() ContractCreateTransaction {
	pb := &proto.ContractCreateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractCreateInstance{ContractCreateInstance: pb}

	builder := ContractCreateTransaction{inner, pb}

	// Default autoRenewPeriod to a value within the required range (~1/4 year)
	return builder.SetAutoRenewPeriod(131500 * time.Minute)
}

func (builder ContractCreateTransaction) SetBytecodeFileID(id FileID) ContractCreateTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder ContractCreateTransaction) SetAdminKey(publicKey Ed25519PublicKey) ContractCreateTransaction {
	builder.pb.AdminKey = publicKey.toProto()
	return builder
}

func (builder ContractCreateTransaction) SetContractMemo(memo string) ContractCreateTransaction {
	builder.pb.Memo = memo
	return builder
}

func (builder ContractCreateTransaction) SetGas(gas uint64) ContractCreateTransaction {
	builder.pb.Gas = int64(gas)
	return builder
}

func (builder ContractCreateTransaction) SetInitialBalance(initialBalance Hbar) ContractCreateTransaction {
	builder.pb.InitialBalance = initialBalance.AsTinybar()
	return builder
}

func (builder ContractCreateTransaction) SetProxyAccountID(id AccountID) ContractCreateTransaction {
	builder.pb.ProxyAccountID = id.toProto()
	return builder
}

func (builder ContractCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) ContractCreateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return builder
}

func (builder ContractCreateTransaction) SetConstructorParams(params *ContractFunctionParams) ContractCreateTransaction {
	builder.pb.ConstructorParameters = params.build(nil)
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ContractCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder ContractCreateTransaction) SetTransactionMemo(memo string) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder ContractCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder ContractCreateTransaction) SetTransactionID(transactionID TransactionID) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ContractCreateTransaction) SetNodeAccountID(nodeAccountID AccountID) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
