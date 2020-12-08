package hedera

import (
	"github.com/pkg/errors"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ContractCreateTransaction struct {
	Transaction
	pb *proto.ContractCreateTransactionBody
}

func NewContractCreateTransaction() *ContractCreateTransaction {
	pb := &proto.ContractCreateTransactionBody{}

	transaction := ContractCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(131500 * time.Minute)

	return &transaction
}

func contractCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractCreateTransaction {
	return ContractCreateTransaction{
		Transaction: transaction,
		pb:          pb.GetContractCreateInstance(),
	}
}

func (transaction *ContractCreateTransaction) SetBytecodeFileID(bytecodeFileID FileID) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FileID = bytecodeFileID.toProtobuf()
	return transaction
}

func (transaction *ContractCreateTransaction) GetBytecodeFileID() FileID {
	return fileIDFromProtobuf(transaction.pb.FileID)
}

/**
 * Sets the state of the instance and its fields can be modified arbitrarily if this key signs a transaction
 * to modify it. If this is null, then such modifications are not possible, and there is no administrator
 * that can override the normal operation of this smart contract instance. Note that if it is created with no
 * admin keys, then there is no administrator to authorize changing the admin keys, so
 * there can never be any admin keys for that instance.
 */
func (transaction *ContractCreateTransaction) SetAdminKey(adminKey Key) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AdminKey = adminKey.toProtoKey()
	return transaction
}

func (transaction *ContractCreateTransaction) GetAdminKey() (Key, error) {
	return keyFromProtobuf(transaction.pb.GetAdminKey())
}

// Sets the gas to run the constructor.
func (transaction *ContractCreateTransaction) SetGas(gas uint64) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Gas = int64(gas)
	return transaction
}

func (transaction *ContractCreateTransaction) GetGas() uint64 {
	return uint64(transaction.pb.GetGas())
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (transaction *ContractCreateTransaction) SetInitialBalance(initialBalance Hbar) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.InitialBalance = initialBalance.AsTinybar()
	return transaction
}

// GetInitialBalance gets the initial number of Hbar in the account
func (transaction *ContractCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(transaction.pb.GetInitialBalance())
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (transaction *ContractCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = durationToProtobuf(autoRenewPeriod)
	return transaction
}

func (transaction *ContractCreateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProtobuf(transaction.pb.GetAutoRenewPeriod())
}

// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalID account, or is an account that isn't a node, then this account is automatically proxy staked to a node
// chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a node, then it will behave as if proxyAccountID was not set.
func (transaction *ContractCreateTransaction) SetProxyAccountID(ID AccountID) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ProxyAccountID = ID.toProtobuf()
	return transaction
}

func (transaction *ContractCreateTransaction) GetProxyAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.ProxyAccountID)
}

//Sets the constructor parameters
func (transaction *ContractCreateTransaction) SetConstructorParameters(params *ContractFunctionParameters) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ConstructorParameters = params.build(nil)
	return transaction
}

//Sets the constructor parameters as their raw bytes.
func (transaction *ContractCreateTransaction) SetConstructorParametersRaw(params []byte) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ConstructorParameters = params
	return transaction
}

func (transaction *ContractCreateTransaction) GetConstructorParameters() []byte {
	return transaction.pb.ConstructorParameters
}

//Sets the memo to be associated with this contract.
func (transaction *ContractCreateTransaction) SetContractMemo(memo string) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = memo
	return transaction
}

func (transaction *ContractCreateTransaction) GetContractMemo() string {
	return transaction.pb.GetMemo()
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func contractCreateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getContract().CreateContract,
	}
}

func (transaction *ContractCreateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractCreateTransaction) Sign(
	privateKey PrivateKey,
) *ContractCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ContractCreateTransaction) SignWithOperator(
	client *Client,
) (*ContractCreateTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client == nil {
		return nil, errors.Wrap(errNoClientProvided, "for SignWithOperator")
	} else if client.operator == nil {
		return nil, errors.Wrap(errClientOperatorSigning, "for SignWithOperator")
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, errors.Wrap(err, "FreezeWith in SignWithOperator")
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *ContractCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractCreateTransaction {
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
func (transaction *ContractCreateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
		return TransactionResponse{}, errors.Wrap(errNoClientProvided, "for Execution")
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, errors.Wrap(err, "FreezeWith in Execute")
		}
	}

	transactionID := transaction.transactionIDs[0]

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
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
		contractCreateTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, errors.Wrap(err, "execution error")
	}

	return TransactionResponse{
		TransactionID: transaction.transactionIDs[0],
		NodeID:        resp.transaction.NodeID,
	}, nil
}

func (transaction *ContractCreateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ContractCreateInstance{
		ContractCreateInstance: transaction.pb,
	}

	return true
}

func (transaction *ContractCreateTransaction) Freeze() (*ContractCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractCreateTransaction) FreezeWith(client *Client) (*ContractCreateTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, errors.Wrap(err, "initTransactionID in ContractCreateTransaction.FreezeWith")
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *ContractCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetMaxTransactionFee(fee Hbar) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetTransactionMemo(memo string) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetTransactionValidDuration(duration time.Duration) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetTransactionID(transactionID TransactionID) *ContractCreateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *ContractCreateTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeAccountIDs sets the node AccountID for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ContractCreateTransaction) SetMaxRetry(count int) *ContractCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}
