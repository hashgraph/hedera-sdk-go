package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
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

func (transaction *ContractDeleteTransaction) SetContractID(contractID ContractID) *ContractDeleteTransaction {
	transaction.pb.ContractID = contractID.toProto()
	return transaction
}

func (transaction *ContractDeleteTransaction) GetContractID() ContractID {
	return contractIDFromProto(transaction.pb.GetContractID())
}

func (transaction *ContractDeleteTransaction) SetTransferContractID(contractID ContractID) *ContractDeleteTransaction {
	transaction.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferContractID{
		TransferContractID: contractID.toProto(),
	}

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferContractID() ContractID {
	return contractIDFromProto(transaction.pb.GetTransferContractID())
}

func (transaction *ContractDeleteTransaction) SetTransferAccountID(accountID AccountID) *ContractDeleteTransaction {
	transaction.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferAccountID{
		TransferAccountID: accountID.toProtobuf(),
	}

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferAccountID() AccountID {
	return accountIDFromProto(transaction.pb.GetTransferAccountID())
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func contractDeleteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getContract().DeleteContract,
	}
}

func (transaction *ContractDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractDeleteTransaction) Sign(
	privateKey PrivateKey,
) *ContractDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ContractDeleteTransaction) SignWithOperator(
	client *Client,
) (*ContractDeleteTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *ContractDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractDeleteTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.transactions); index++ {
		signature := signer(transaction.transactions[index].GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractDeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.id

	if !client.GetOperatorID().isZero() && client.GetOperatorID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorKey(),
			client.operator.signer,
		)
	}

	_, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeId,
		contractDeleteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *ContractDeleteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ContractDeleteInstance{
		ContractDeleteInstance: transaction.pb,
	}

	return true
}

func (transaction *ContractDeleteTransaction) Freeze() (*ContractDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractDeleteTransaction) FreezeWith(client *Client) (*ContractDeleteTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *ContractDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ContractDeleteTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionMemo(memo string) *ContractDeleteTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractDeleteTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) *ContractDeleteTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetNodeID(nodeID AccountID) *ContractDeleteTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
