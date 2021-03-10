package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

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
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func contractDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractDeleteTransaction {
	return ContractDeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetContractDeleteInstance(),
	}
}

//Sets the contract ID which should be deleted.
func (transaction *ContractDeleteTransaction) SetContractID(contractID ContractID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ContractID = contractID.toProtobuf()
	return transaction
}

func (transaction *ContractDeleteTransaction) GetContractID() ContractID {
	return contractIDFromProtobuf(transaction.pb.GetContractID())
}

//Sets the contract ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) SetTransferContractID(contractID ContractID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferContractID{
		TransferContractID: contractID.toProtobuf(),
	}

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferContractID() ContractID {
	return contractIDFromProtobuf(transaction.pb.GetTransferContractID())
}

//Sets the account ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) SetTransferAccountID(accountID AccountID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Obtainers = &proto.ContractDeleteTransactionBody_TransferAccountID{
		TransferAccountID: accountID.toProtobuf(),
	}

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetTransferAccountID())
}

func (transaction *ContractDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	txBytes, err := protobuf.Marshal(transaction.constructProtobuf())
	if err != nil {
		return &ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setTransactionBodyBytes(txBytes), nil
}

func (transaction *ContractDeleteTransaction) constructProtobuf() *proto.TransactionBody{
	return &proto.TransactionBody{
		TransactionID:            transaction.pbBody.GetTransactionID(),
		NodeAccountID:            transaction.pbBody.GetNodeAccountID(),
		TransactionFee:           transaction.pbBody.GetTransactionFee(),
		TransactionValidDuration: transaction.pbBody.GetTransactionValidDuration(),
		GenerateRecord:           transaction.pbBody.GetGenerateRecord(),
		Memo:                     transaction.pbBody.GetMemo(),
		Data: &proto.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: &proto.ContractDeleteTransactionBody{
				ContractID: transaction.pb.GetContractID(),
				Obtainers:  transaction.pb.GetObtainers(),
			},
		},
	}
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

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
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
		_, _ = transaction.Freeze()
	} else {
		transaction.transactions = make([]*proto.Transaction, 0)
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.signedTransactions); index++ {
		signature := signer(transaction.signedTransactions[index].GetBodyBytes())

		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractDeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := transaction.GetTransactionID()

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		contractDeleteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()

	return TransactionResponse{
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
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
	if transaction.IsFrozen() {
		return transaction, nil
	}
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
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionMemo(memo string) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the node AccountID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ContractDeleteTransaction) SetMaxRetry(count int) *ContractDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *ContractDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractDeleteTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
