package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ContractCreateTransaction struct {
	Transaction
	pb             *proto.ContractCreateTransactionBody
	byteCodeFileID FileID
	proxyAccountID AccountID
}

func NewContractCreateTransaction() *ContractCreateTransaction {
	pb := &proto.ContractCreateTransactionBody{}

	transaction := ContractCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(131500 * time.Minute)
	transaction.SetMaxTransactionFee(NewHbar(20))

	return &transaction
}

func contractCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractCreateTransaction {
	return ContractCreateTransaction{
		Transaction:    transaction,
		pb:             pb.GetContractCreateInstance(),
		byteCodeFileID: fileIDFromProtobuf(pb.GetContractCreateInstance().GetFileID(), nil),
		proxyAccountID: accountIDFromProtobuf(pb.GetContractCreateInstance().GetProxyAccountID(), nil),
	}
}

func (transaction *ContractCreateTransaction) SetBytecodeFileID(bytecodeFileID FileID) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.byteCodeFileID = bytecodeFileID
	return transaction
}

func (transaction *ContractCreateTransaction) GetBytecodeFileID() FileID {
	return transaction.byteCodeFileID
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
	return keyFromProtobuf(transaction.pb.GetAdminKey(), nil)
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
func (transaction *ContractCreateTransaction) SetProxyAccountID(id AccountID) *ContractCreateTransaction {
	transaction.requireNotFrozen()
	transaction.proxyAccountID = id
	return transaction
}

func (transaction *ContractCreateTransaction) GetProxyAccountID() AccountID {
	return transaction.proxyAccountID
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

func (transaction *ContractCreateTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.byteCodeFileID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.proxyAccountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *ContractCreateTransaction) build() *ContractCreateTransaction {
	if !transaction.byteCodeFileID.isZero() {
		transaction.pb.FileID = transaction.byteCodeFileID.toProtobuf()
	}

	if !transaction.proxyAccountID.isZero() {
		transaction.pb.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	return transaction
}

func (transaction *ContractCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_ContractCreateInstance{
			ContractCreateInstance: &proto.ContractCreateTransactionBody{
				FileID:                transaction.pb.GetFileID(),
				AdminKey:              transaction.pb.GetAdminKey(),
				Gas:                   transaction.pb.GetGas(),
				InitialBalance:        transaction.pb.GetInitialBalance(),
				ProxyAccountID:        transaction.pb.GetProxyAccountID(),
				AutoRenewPeriod:       transaction.pb.GetAutoRenewPeriod(),
				ConstructorParameters: transaction.pb.GetConstructorParameters(),
				ShardID:               transaction.pb.GetShardID(),
				RealmID:               transaction.pb.GetRealmID(),
				NewRealmAdminKey:      transaction.pb.GetNewRealmAdminKey(),
				Memo:                  transaction.pb.GetMemo(),
			},
		},
	}, nil
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
func (transaction *ContractCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractCreateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractCreateTransaction) Execute(
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
		contractCreateTransaction_getMethod,
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
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &ContractCreateTransaction{}, err
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

func (transaction *ContractCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractCreateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
