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
	transaction.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferContractID{
		TransferContractID: contractId.toProto(),
	}

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferContractId() ContractID {
	return transaction.GetTransferContractId()
}

func (transaction *ContractDeleteTransaction) SetTransferAccountId(accountId AccountID) *ContractDeleteTransaction {
	transaction.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferAccountID{
		TransferAccountID: accountId.toProtobuf(),
	}

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferAccountId() AccountID {
	return transaction.GetTransferAccountId()
}

