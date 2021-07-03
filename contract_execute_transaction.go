package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

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
	pb         *proto.ContractCallTransactionBody
	contractID ContractID
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call Transaction.
func NewContractExecuteTransaction() *ContractExecuteTransaction {
	pb := &proto.ContractCallTransactionBody{}

	transaction := ContractExecuteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func contractExecuteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractExecuteTransaction {
	return ContractExecuteTransaction{
		Transaction: transaction,
		pb:          pb.GetContractCall(),
		contractID:  contractIDFromProtobuf(pb.GetContractCall().GetContractID(), nil),
	}
}

// SetContractID sets the contract instance to call.
func (transaction *ContractExecuteTransaction) SetContractID(id ContractID) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.contractID = id
	return transaction
}

func (transaction ContractExecuteTransaction) GetContractID() ContractID {
	return transaction.contractID
}

// SetGas sets the maximum amount of gas to use for the call.
func (transaction *ContractExecuteTransaction) SetGas(gas uint64) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Gas = int64(gas)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetGas() uint64 {
	return uint64(transaction.pb.Gas)
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (transaction *ContractExecuteTransaction) SetPayableAmount(amount Hbar) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Amount = amount.AsTinybar()
	return transaction
}

func (transaction ContractExecuteTransaction) GetPayableAmount() Hbar {
	return HbarFromTinybar(transaction.pb.Amount)
}

//Sets the function parameters
func (transaction *ContractExecuteTransaction) SetFunctionParameters(params []byte) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FunctionParameters = params
	return transaction
}

func (transaction *ContractExecuteTransaction) GetFunctionParameters() []byte {
	return transaction.pb.GetFunctionParameters()
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (transaction *ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParameters) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	if params == nil {
		params = NewContractFunctionParameters()
	}

	transaction.pb.FunctionParameters = params.build(&name)
	return transaction
}

func (transaction *ContractExecuteTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.contractID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *ContractExecuteTransaction) build() *ContractExecuteTransaction {
	if !transaction.contractID.isZero() {
		transaction.pb.ContractID = transaction.contractID.toProtobuf()
	}

	return transaction
}

func (transaction *ContractExecuteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractExecuteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_ContractCall{
			ContractCall: &proto.ContractCallTransactionBody{
				ContractID:         transaction.pb.GetContractID(),
				Gas:                transaction.pb.GetGas(),
				Amount:             transaction.pb.GetAmount(),
				FunctionParameters: transaction.pb.GetFunctionParameters(),
			},
		},
	}, nil
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
func (transaction *ContractExecuteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractExecuteTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractExecuteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
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
		contractExecuteTransaction_getMethod,
		transaction_mapStatusError,
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
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &ContractExecuteTransaction{}, err
	}
	transaction.build()

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
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionMemo(memo string) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) *ContractExecuteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the node AccountID for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ContractExecuteTransaction) SetMaxRetry(count int) *ContractExecuteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *ContractExecuteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractExecuteTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
