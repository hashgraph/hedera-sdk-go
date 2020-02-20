package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ContractDeleteTransaction struct {
	TransactionBuilder
	pb *proto.ContractDeleteTransactionBody
}

// NewContractDeleteTransaction creates a Contract Delete Transaction builder which can be used to construct and execute
// a Contract Delete Transaction.
func NewContractDeleteTransaction() ContractDeleteTransaction {
	pb := &proto.ContractDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractDeleteInstance{ContractDeleteInstance: pb}

	builder := ContractDeleteTransaction{inner, pb}

	return builder
}

// SetContractID sets the Contract ID of the Contract to be deleted by the Contract Delete Transaction
func (builder ContractDeleteTransaction) SetContractID(id ContractID) ContractDeleteTransaction {
	builder.pb.ContractID = id.toProto()
	return builder
}

// SetTransferAccountID sets the Account ID which will receive remaining hbar tied to the Contract
func (builder ContractDeleteTransaction) SetTransferAccountID(id AccountID) ContractDeleteTransaction {
	builder.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferAccountID{
		TransferAccountID: id.toProto(),
	}

	return builder
}

func (builder ContractDeleteTransaction) SetTransferContractID(id ContractID) ContractDeleteTransaction {
	builder.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferContractID{
		TransferContractID: id.toProto(),
	}

	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ContractDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ContractDeleteTransaction {
	return ContractDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder ContractDeleteTransaction) SetTransactionMemo(memo string) ContractDeleteTransaction {
	return ContractDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder ContractDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractDeleteTransaction {
	return ContractDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) ContractDeleteTransaction {
	return ContractDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ContractDeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) ContractDeleteTransaction {
	return ContractDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
