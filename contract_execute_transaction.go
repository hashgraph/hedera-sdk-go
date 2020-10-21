package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

// ContractExecuteTransaction calls a function of the given smart contract instance, giving it ContractFuncionParams as
// its inputs. it can use the given amount of gas, and any unspent gas will be refunded to the paying account.
//
// If this function stores information, it is charged gas to store it. There is a fee in hbars to maintain that storage
// until the expiration time, and that fee is added as part of the transaction fee.
//
// For a cheaper but more limited method to call functions, see ContractCallQuery.
type ContractExecuteTransaction struct {
	Transaction
	pb *proto.ContractCallTransactionBody
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call Transaction.
func NewContractExecuteTransaction() *ContractExecuteTransaction {
	pb := &proto.ContractCallTransactionBody{}

	transaction := ContractExecuteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

// SetContractID sets the contract instance to call.
func (transaction *ContractExecuteTransaction) SetContractID(ID ContractID) *ContractExecuteTransaction {
	transaction.pb.ContractID = ID.toProtobuf()
	return transaction
}

func (transaction ContractExecuteTransaction) GetContractID() ContractID {
	return contractIDFromProtobuf(transaction.pb.GetContractID())
}

// SetGas sets the maximum amount of gas to use for the call.
func (transaction *ContractExecuteTransaction) SetGas(gas uint64) *ContractExecuteTransaction {
	transaction.pb.Gas = int64(gas)
	return transaction
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (transaction *ContractExecuteTransaction) SetPayableAmount(amount Hbar) *ContractExecuteTransaction {
	transaction.pb.Amount = amount.AsTinybar()
	return transaction
}

func (transaction ContractExecuteTransaction) GetPayableAmount() uint64 {
	return uint64(transaction.pb.Gas)
}

func (transaction *ContractExecuteTransaction) SetFunctionParameters(params []byte) *ContractExecuteTransaction {
	transaction.pb.FunctionParameters = params
	return transaction
}

func (transaction *ContractExecuteTransaction) GetFunctionParameters() []byte {
	return transaction.pb.GetFunctionParameters()
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (transaction *ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParameters) *ContractExecuteTransaction {
	if params == nil {
		params = NewContractFunctionParameters()
	}

	transaction.pb.FunctionParameters = params.build(&name)
	return transaction
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func contractExecuteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getContract().ContractCallMethod,
	}
}

func (transaction *ContractExecuteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractExecuteTransaction) Sign(
	privateKey PrivateKey,
) *ContractExecuteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ContractExecuteTransaction) SignWithOperator(
	client *Client,
) (*ContractExecuteTransaction, error) {
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
func (transaction *ContractExecuteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractExecuteTransaction {
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
func (transaction *ContractExecuteTransaction) Execute(
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

	resp, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeId,
		contractExecuteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.id,
		NodeID:        resp.transaction.NodeID,
	}, nil
}

func (transaction *ContractExecuteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ContractCall{
		ContractCall: transaction.pb,
	}

	return true
}

func (transaction *ContractExecuteTransaction) Freeze() (*ContractExecuteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractExecuteTransaction) FreezeWith(client *Client) (*ContractExecuteTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *ContractExecuteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetMaxTransactionFee(fee Hbar) *ContractExecuteTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionMemo(memo string) *ContractExecuteTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractExecuteTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) *ContractExecuteTransaction {
	transaction.id = transactionID
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeAccountID sets the node AccountID for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractExecuteTransaction {
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}
