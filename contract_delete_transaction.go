package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractDeleteTransaction struct {
	TransactionBuilder
	pb *proto.ContractDeleteTransactionBody
}

// NewContractDeleteTransaction creates a Contract Delete Transaction transaction which can be used to construct and execute
// a Contract Delete Transaction.
func NewContractDeleteTransaction() ContractDeleteTransaction {
	pb := &proto.ContractDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractDeleteInstance{ContractDeleteInstance: pb}

	transaction := ContractDeleteTransaction{inner, pb}

	return transaction
}

// SetContractID sets the Contract ID of the Contract to be deleted by the Contract Delete Transaction
func (transaction ContractDeleteTransaction) SetContractID(id ContractID) ContractDeleteTransaction {
	transaction.pb.ContractID = id.toProto()
	return transaction
}

// SetTransferAccountID sets the Account ID which will receive remaining hbar tied to the Contract
func (transaction ContractDeleteTransaction) SetTransferAccountID(id AccountID) ContractDeleteTransaction {
	transaction.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferAccountID{
		TransferAccountID: id.toProto(),
	}

	return transaction
}

func (transaction ContractDeleteTransaction) SetTransferContractID(id ContractID) ContractDeleteTransaction {
	transaction.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferContractID{
		TransferContractID: id.toProto(),
	}

	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction ContractDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ContractDeleteTransaction {
	return ContractDeleteTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction ContractDeleteTransaction) SetTransactionMemo(memo string) ContractDeleteTransaction {
	return ContractDeleteTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction ContractDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractDeleteTransaction {
	return ContractDeleteTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) ContractDeleteTransaction {
	return ContractDeleteTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction ContractDeleteTransaction) SetNodeID(nodeAccountID AccountID) ContractDeleteTransaction {
	return ContractDeleteTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
