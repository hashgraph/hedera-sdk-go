package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// ContractUpdateTransaction is used to modify a smart contract instance to have the given parameter values. Any nil
// field is ignored (left unchanged). If only the contractInstanceExpirationTime is being modified, then no signature is
// needed on this transaction other than for the account paying for the transaction itself. But if any of the other
// fields are being modified, then it must be signed by the adminKey. The use of adminKey is not currently supported in
// this API, but in the future will be implemented to allow these fields to be modified, and also to make modifications
// to the state of the instance. If the contract is created with no admin key, then none of the fields can be changed
// that need an admin signature, and therefore no admin key can ever be added. So if there is no admin key, then things
// like the bytecode are immutable. But if there is an admin key, then they can be changed.
//
// For example, the admin key might be a threshold key, which requires 3 of 5 binding arbitration judges to agree before
// the bytecode can be changed. This can be used to add flexibility to the management of smart contract behavior. But
// this is optional. If the smart contract is created without an admin key, then such a key can never be added, and its
// bytecode will be immutable.
type ContractUpdateTransaction struct {
	Transaction
	pb *proto.ContractUpdateTransactionBody
}

// NewContractUpdateTransaction creates a ContractUpdateTransaction transaction which can be
// used to construct and execute a Contract Update Transaction.
func NewContractUpdateTransaction() *ContractUpdateTransaction {
	pb := &proto.ContractUpdateTransactionBody{}

	transaction := ContractUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func contractUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractUpdateTransaction {
	return ContractUpdateTransaction{
		Transaction: transaction,
		pb:          pb.GetContractUpdateInstance(),
	}
}

// SetContractID sets The Contract ID instance to update (this can't be changed on the contract)
func (transaction *ContractUpdateTransaction) SetContractID(contractID ContractID) *ContractUpdateTransaction {
	transaction.pb.ContractID = contractID.toProtobuf()
	return transaction
}

func (transaction *ContractUpdateTransaction) GetContractID() ContractID {
	return contractIDFromProtobuf(transaction.pb.GetContractID())
}

// SetBytecodeFileID sets the file ID of file containing the smart contract byte code. A copy will be made and held by
// the contract instance, and have the same expiration time as the instance.
func (transaction *ContractUpdateTransaction) SetBytecodeFileID(fileID FileID) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FileID = fileID.toProtobuf()
	return transaction
}

func (transaction *ContractUpdateTransaction) GetBytecodeFileID() FileID {
	return fileIDFromProtobuf(transaction.pb.GetFileID())
}

// SetAdminKey sets the key which can be used to arbitrarily modify the state of the instance by signing a
// ContractUpdateTransaction to modify it. If the admin key was never set then such modifications are not possible,
// and there is no administrator that can overrIDe the normal operation of the smart contract instance.
func (transaction *ContractUpdateTransaction) SetAdminKey(publicKey PublicKey) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AdminKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *ContractUpdateTransaction) GetAdminKey() (Key, error) {
	return keyFromProtobuf(transaction.pb.GetAdminKey())
}

// SetProxyAccountID sets the ID of the account to which this contract is proxy staked. If proxyAccountID is left unset,
// is an invalID account, or is an account that isn't a node, then this contract is automatically proxy staked to a node
// chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking,
// or if it is not currently running a node, then it will behave as if proxyAccountID was never set.
func (transaction *ContractUpdateTransaction) SetProxyAccountID(accountID AccountID) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ProxyAccountID = accountID.toProtobuf()
	return transaction
}

func (transaction *ContractUpdateTransaction) GetProxyAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetProxyAccountID())
}

// SetAutoRenewPeriod sets the duration for which the contract instance will automatically charge its account to
// renew for.
func (transaction *ContractUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = durationToProtobuf(autoRenewPeriod)
	return transaction
}

func (transaction *ContractUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProtobuf(transaction.pb.GetAutoRenewPeriod())
}

// SetExpirationTime extends the expiration of the instance and its account to the provIDed time. If the time provIDed
// is the current or past time, then there will be no effect.
func (transaction *ContractUpdateTransaction) SetExpirationTime(expiration time.Time) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ExpirationTime = timeToProtobuf(expiration)
	return transaction
}

func (transaction *ContractUpdateTransaction) GetExpirationTime() time.Time {
	return timeFromProtobuf(transaction.pb.GetExpirationTime())
}

// SetContractMemo sets the memo associated with the contract (max 100 bytes)
func (transaction *ContractUpdateTransaction) SetContractMemo(memo string) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = memo
	return transaction
}

func (transaction *ContractUpdateTransaction) GetContractMemo() string {
	return transaction.pb.GetMemo()
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func contractUpdateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getContract().UpdateContract,
	}
}

func (transaction *ContractUpdateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractUpdateTransaction) Sign(
	privateKey PrivateKey,
) *ContractUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ContractUpdateTransaction) SignWithOperator(
	client *Client,
) (*ContractUpdateTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
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
func (transaction *ContractUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractUpdateTransaction {
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
func (transaction *ContractUpdateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
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
		contractUpdateTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.transactionIDs[0],
		NodeID:        resp.transaction.NodeID,
	}, nil
}

func (transaction *ContractUpdateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ContractUpdateInstance{
		ContractUpdateInstance: transaction.pb,
	}

	return true
}

func (transaction *ContractUpdateTransaction) Freeze() (*ContractUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractUpdateTransaction) FreezeWith(client *Client) (*ContractUpdateTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *ContractUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetMaxTransactionFee(fee Hbar) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetTransactionMemo(memo string) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetTransactionID(transactionID TransactionID) *ContractUpdateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *ContractUpdateTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeAccountID sets the node AccountID for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}
