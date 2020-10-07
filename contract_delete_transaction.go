package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractDeleteTransaction struct {
	Transaction
	pb *proto.ContractDeleteTransactionBody
}

func NewContractDeleteTransaction() *ContractDeleteTransaction {
	pb := &proto.ContractDeleteTransactionBody{}

	transaction := ContractDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *ContractDeleteTransaction) SetContractId(contractId ContractID) *ContractDeleteTransaction {
	transaction.pb.ContractID = contractId.toProto()
	return transaction
}

func (transaction *ContractDeleteTransaction) GetContractId() ContractID {
	return transaction.GetContractId()
}

func (transaction *ContractDeleteTransaction) SetTransferContractId(contractId ContractID) *ContractDeleteTransaction {
	transaction.pb.Tra = contractId.toProto()
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferContractId() ContractID {
	return transaction.GetContractId()
}

