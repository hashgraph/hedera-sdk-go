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
	contractID      ContractID
	proxyAccountID  AccountID
	bytecodeFileID  FileID
	adminKey        Key
	gas             int64
	initialBalance  int64
	autoRenewPeriod *time.Duration
	expirationTime  *time.Time
	parameters      []byte
	memo            string
}

// NewContractUpdateTransaction creates a ContractUpdateTransaction transaction which can be
// used to construct and execute a Contract Update Transaction.
func NewContractUpdateTransaction() *ContractUpdateTransaction {
	transaction := ContractUpdateTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func contractUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractUpdateTransaction {
	key, _ := keyFromProtobuf(pb.GetContractUpdateInstance().AdminKey)
	autoRenew := durationFromProtobuf(pb.GetContractUpdateInstance().GetAutoRenewPeriod())
	expiration := timeFromProtobuf(pb.GetContractUpdateInstance().GetExpirationTime())
	var memo string

	switch m := pb.GetContractUpdateInstance().GetMemoField().(type) {
	case *proto.ContractUpdateTransactionBody_Memo:
		memo = m.Memo
	case *proto.ContractUpdateTransactionBody_MemoWrapper:
		memo = m.MemoWrapper.Value
	}

	return ContractUpdateTransaction{
		Transaction:     transaction,
		contractID:      contractIDFromProtobuf(pb.GetContractUpdateInstance().GetContractID()),
		proxyAccountID:  accountIDFromProtobuf(pb.GetContractUpdateInstance().GetProxyAccountID()),
		bytecodeFileID:  fileIDFromProtobuf(pb.GetContractUpdateInstance().GetFileID()),
		adminKey:        key,
		autoRenewPeriod: &autoRenew,
		expirationTime:  &expiration,
		memo:            memo,
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
	transaction.adminKey = publicKey
	return transaction
}

func (transaction *ContractUpdateTransaction) GetAdminKey() (Key, error) {
	return transaction.adminKey, nil
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
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (transaction *ContractUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetExpirationTime extends the expiration of the instance and its account to the provIDed time. If the time provIDed
// is the current or past time, then there will be no effect.
func (transaction *ContractUpdateTransaction) SetExpirationTime(expiration time.Time) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

func (transaction *ContractUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

// SetContractMemo sets the memo associated with the contract (max 100 bytes)
func (transaction *ContractUpdateTransaction) SetContractMemo(memo string) *ContractUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.memo = memo
	//if transaction.pb.GetMemoWrapper() != nil {
	//	transaction.pb.GetMemoWrapper().Value = memo
	//} else {
	//	transaction.pb.MemoField = &proto.ContractUpdateTransactionBody_MemoWrapper{
	//		MemoWrapper: &wrappers.StringValue{Value: memo},
	//	}
	//}

	return transaction
}

func (transaction *ContractUpdateTransaction) GetContractMemo() string {
	return transaction.memo
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

func (transaction *ContractUpdateTransaction) build() *proto.TransactionBody {
	body := &proto.ContractUpdateTransactionBody{}

	if transaction.expirationTime != nil {
		body.ExpirationTime = timeToProtobuf(*transaction.expirationTime)
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey.toProtoKey()
	}

	if !transaction.contractID.isZero() {
		body.ContractID = transaction.contractID.toProtobuf()
	}

	if !transaction.proxyAccountID.isZero() {
		body.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	if !transaction.bytecodeFileID.isZero() {
		body.FileID = transaction.bytecodeFileID.toProtobuf()
	}

	if body.GetMemoWrapper() != nil {
		body.GetMemoWrapper().Value = transaction.memo
	} else {
		body.MemoField = &proto.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrappers.StringValue{Value: transaction.memo},
		}
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: body,
		},
	}
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
	body := &proto.ContractUpdateTransactionBody{}

	if transaction.expirationTime != nil {
		body.ExpirationTime = timeToProtobuf(*transaction.expirationTime)
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey.toProtoKey()
	}

	if !transaction.contractID.isZero() {
		body.ContractID = transaction.contractID.toProtobuf()
	}

	if !transaction.proxyAccountID.isZero() {
		body.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	if !transaction.bytecodeFileID.isZero() {
		body.FileID = transaction.bytecodeFileID.toProtobuf()
	}

	if body.GetMemoWrapper() != nil {
		body.GetMemoWrapper().Value = transaction.memo
	} else {
		body.MemoField = &proto.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrappers.StringValue{Value: transaction.memo},
		}
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey.toProtoKey()
	}

	if !transaction.contractID.isZero() {
		body.ContractID = transaction.contractID.toProtobuf()
	}

	if !transaction.proxyAccountID.isZero() {
		body.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	if !transaction.bytecodeFileID.isZero() {
		body.FileID = transaction.bytecodeFileID.toProtobuf()
	}

	if body.GetMemoWrapper() != nil {
		body.GetMemoWrapper().Value = transaction.memo
	} else {
		body.MemoField = &proto.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrappers.StringValue{Value: transaction.memo},
		}
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: body,
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
		transaction_makeRequest(request{
			transaction: &transaction.Transaction,
		}),
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
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
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
	transaction.requireOneNodeAccountID()

	if !transaction.isFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	if len(transaction.signedTransactions) == 0 {
		return transaction
	}

	transaction.transactions = make([]*proto.Transaction, 0)
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)

	for index := 0; index < len(transaction.signedTransactions); index++ {
		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	//transaction.signedTransactions[0].SigMap.SigPair = append(transaction.signedTransactions[0].SigMap.SigPair, publicKey.toSignaturePairProtobuf(signature))
	return transaction
}
