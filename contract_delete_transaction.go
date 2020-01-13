package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ContractDeleteTransaction struct {
	TransactionBuilder
	pb *proto.ContractDeleteTransactionBody
}

func NewContractDeleteTransaction() ContractDeleteTransaction {
	pb := &proto.ContractDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractDeleteInstance{ContractDeleteInstance: pb}

	builder := ContractDeleteTransaction{inner, pb}

	return builder
}

func (builder ContractDeleteTransaction) SetContractID(id ContractID) ContractDeleteTransaction {
	builder.pb.ContractID = id.toProto()
	return builder
}

func (builder ContractDeleteTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ContractDeleteTransaction) SetMaxTransactionFee(maxTransactionFee uint64) ContractDeleteTransaction {
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
