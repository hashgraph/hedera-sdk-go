package hedera

import (
	"github.com/golang/protobuf/ptypes/wrappers"
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
	pb             *proto.ContractUpdateTransactionBody
	contractID     ContractID
	proxyAccountID AccountID
	bytecodeFileID FileID
}

// NewContractUpdateTransaction creates a ContractUpdateTransaction transaction which can be
// used to construct and execute a Contract Update Transaction.
func NewContractUpdateTransaction() *ContractUpdateTransaction {
	pb := &proto.ContractUpdateTransactionBody{}

	transaction := ContractUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func contractUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractUpdateTransaction {
	return ContractUpdateTransaction{
		Transaction:    transaction,
		pb:             pb.GetContractUpdateInstance(),
		contractID:     contractIDFromProtobuf(pb.GetContractUpdateInstance().GetContractID()),
		proxyAccountID: accountIDFromProtobuf(pb.GetContractUpdateInstance().GetProxyAccountID()),
		bytecodeFileID: fileIDFromProtobuf(pb.GetContractUpdateInstance().GetFileID()),
	}
}

// SetContractID sets The Contract ID instance to update (this can't be changed on the contract)
func (transaction *ContractUpdateTransaction) SetContractID(id ContractID) *ContractUpdateTransaction {
	transaction.contractID = id
	return transaction
}

func (transaction *ContractUpdateTransaction) GetContractID() ContractID {
	return transaction.contractID
}

// SetBytecodeFileID sets the file ID of file containing the smart contract byte code. A copy will be made and held by
// the contract instance, and have the same expiration time as the instance.
func (transaction *ContractUpdateTransaction) SetBytecodeFileID(id FileID) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.bytecodeFileID = id
	return transaction
}

func (transaction *ContractUpdateTransaction) GetBytecodeFileID() FileID {
	return transaction.bytecodeFileID
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
func (transaction *ContractUpdateTransaction) SetProxyAccountID(id AccountID) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.proxyAccountID = id
	return transaction
}

func (transaction *ContractUpdateTransaction) GetProxyAccountID() AccountID {
	return transaction.proxyAccountID
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
	if transaction.pb.GetMemoWrapper() != nil {
		transaction.pb.GetMemoWrapper().Value = memo
	} else {
		transaction.pb.MemoField = &proto.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrappers.StringValue{Value: memo},
		}
	}

	return transaction
}

func (transaction *ContractUpdateTransaction) GetContractMemo() string {
	if transaction.pb.GetMemoField() != nil {
		switch transaction.pb.GetMemoField().(type) {
		case *proto.ContractUpdateTransactionBody_Memo:
			return transaction.pb.GetMemo()
		case *proto.ContractUpdateTransactionBody_MemoWrapper:
			return transaction.pb.GetMemoWrapper().Value
		default:
			return ""
		}
	}

	return ""
}

func (transaction *ContractUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = transaction.contractID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.proxyAccountID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.bytecodeFileID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *ContractUpdateTransaction) build() *ContractUpdateTransaction {
	if !transaction.contractID.isZero() {
		transaction.pb.ContractID = transaction.contractID.toProtobuf()
	}

	if !transaction.proxyAccountID.isZero() {
		transaction.pb.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	if !transaction.bytecodeFileID.isZero() {
		transaction.pb.FileID = transaction.bytecodeFileID.toProtobuf()
	}

	return transaction
}

func (transaction *ContractUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractUpdateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: &proto.ContractUpdateTransactionBody{
				ContractID:      transaction.pb.GetContractID(),
				ExpirationTime:  transaction.pb.GetExpirationTime(),
				AdminKey:        transaction.pb.GetAdminKey(),
				ProxyAccountID:  transaction.pb.GetProxyAccountID(),
				AutoRenewPeriod: transaction.pb.GetAutoRenewPeriod(),
				FileID:          transaction.pb.GetFileID(),
				MemoField:       transaction.pb.GetMemoField(),
			},
		},
	}, nil
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
func (transaction *ContractUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractUpdateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractUpdateTransaction) Execute(
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
		contractUpdateTransaction_getMethod,
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
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &ContractUpdateTransaction{}, err
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

// SetNodeAccountID sets the node AccountID for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ContractUpdateTransaction) SetMaxRetry(count int) *ContractUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *ContractUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractUpdateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}

func (transaction *ContractUpdateTransaction) SetMaxBackoff(max time.Duration) *ContractUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *ContractUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *ContractUpdateTransaction) SetMinBackoff(min time.Duration) *ContractUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *ContractUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
