package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// AccountCreateTransaction creates a new account. After the account is created, the AccountID for it is in the receipt,
// or by asking for a Record of the transaction to be created, and retrieving that. The account can then automatically
// generate records for large transfers into it or out of it, which each last for 25 hours. Records are generated for
// any transfer that exceeds the thresholds given here. This account is charged hbar for each record generated, so the
// thresholds are useful for limiting Record generation to happen only for large transactions.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0,
// with a null key. Future versions of the API will support multiple realms and multiple shards.
type AccountCreateTransaction struct {
	Transaction
	proxyAccountID                AccountID
	key                           Key
	initialBalance                uint64
	receiveRecordThreshold        uint64
	sendRecordThreshold           uint64
	autoRenewPeriod               *time.Duration
	memo                          string
	receiverSignatureRequired     bool
	maxAutomaticTokenAssociations uint32
}

// NewAccountCreateTransaction creates an AccountCreateTransaction transaction which can be used to construct and
// execute a Crypto Create Transaction.
func NewAccountCreateTransaction() *AccountCreateTransaction {
	transaction := AccountCreateTransaction{
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(2))

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	transaction.setReceiveRecordThreshold(MaxHbar)
	transaction.setSendRecordThreshold(MaxHbar)

	return &transaction
}

func accountCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) AccountCreateTransaction {
	key, _ := keyFromProtobuf(pb.GetCryptoCreateAccount().GetKey())
	renew := durationFromProtobuf(pb.GetCryptoCreateAccount().GetAutoRenewPeriod())
	return AccountCreateTransaction{
		Transaction:               transaction,
		proxyAccountID:            accountIDFromProtobuf(pb.GetCryptoCreateAccount().GetProxyAccountID()),
		key:                       key,
		initialBalance:            pb.GetCryptoCreateAccount().InitialBalance,
		receiveRecordThreshold:    pb.GetCryptoCreateAccount().GetReceiveRecordThreshold(),
		sendRecordThreshold:       pb.GetCryptoCreateAccount().GetSendRecordThreshold(),
		autoRenewPeriod:           &renew,
		memo:                      pb.GetCryptoCreateAccount().GetMemo(),
		receiverSignatureRequired: pb.GetCryptoCreateAccount().ReceiverSigRequired,
	}
}

// SetKey sets the key that must sign each transfer out of the account. If RecieverSignatureRequired is true, then it
// must also sign any transfer into the account.
func (transaction *AccountCreateTransaction) SetKey(key Key) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.key = key
	return transaction
}

func (transaction *AccountCreateTransaction) GetKey() (Key, error) {
	return transaction.key, nil
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (transaction *AccountCreateTransaction) SetInitialBalance(initialBalance Hbar) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.initialBalance = uint64(initialBalance.AsTinybar())
	return transaction
}

func (transaction *AccountCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(int64(transaction.initialBalance))
}

func (transaction *AccountCreateTransaction) SetMaxAutomaticTokenAssociations(max uint32) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.maxAutomaticTokenAssociations = max
	return transaction
}

func (transaction *AccountCreateTransaction) GetMaxAutomaticTokenAssociations() uint32 {
	return transaction.maxAutomaticTokenAssociations
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (transaction *AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (transaction *AccountCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetSendRecordThreshold sets the threshold amount for which an account record is created for any send/withdraw
// transaction
//
// Deprecated: No longer used by Hedera
func (transaction *AccountCreateTransaction) setSendRecordThreshold(recordThreshold Hbar) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.sendRecordThreshold = uint64(recordThreshold.AsTinybar())
	return transaction
}

func (transaction *AccountCreateTransaction) getSendRecordThreshold() Hbar {
	return HbarFromTinybar(int64(transaction.sendRecordThreshold))
}

// SetReceiveRecordThreshold sets the threshold amount for which an account record is created for any receive/deposit
// transaction
//
// Deprecated: No longer used by Hedera
func (transaction *AccountCreateTransaction) setReceiveRecordThreshold(recordThreshold Hbar) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.receiveRecordThreshold = uint64(recordThreshold.AsTinybar())
	return transaction
}

func (transaction *AccountCreateTransaction) getReceiveRecordThreshold() Hbar {
	return HbarFromTinybar(int64(transaction.receiveRecordThreshold))
}

// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalid account, or is an account that isn't a node, then this account is automatically proxy staked to a node
// chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a node, then it will behave as if proxyAccountID was not set.
func (transaction *AccountCreateTransaction) SetProxyAccountID(id AccountID) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.proxyAccountID = id
	return transaction
}

func (transaction *AccountCreateTransaction) GetProxyAccountID() AccountID {
	return transaction.proxyAccountID
}

func (transaction *AccountCreateTransaction) SetAccountMemo(memo string) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.memo = memo
	return transaction
}

func (transaction *AccountCreateTransaction) GetAccountMemo() string {
	return transaction.memo
}

func (transaction *AccountCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	return transaction.proxyAccountID.Validate(client)
}

func (transaction *AccountCreateTransaction) build() *proto.TransactionBody {
	body := &proto.CryptoCreateTransactionBody{
		InitialBalance:                transaction.initialBalance,
		SendRecordThreshold:           transaction.receiveRecordThreshold,
		ReceiveRecordThreshold:        transaction.sendRecordThreshold,
		ReceiverSigRequired:           transaction.receiverSignatureRequired,
		Memo:                          transaction.memo,
		MaxAutomaticTokenAssociations: transaction.maxAutomaticTokenAssociations,
	}

	if transaction.key != nil {
		body.Key = transaction.key.toProtoKey()
	}

	if !transaction.proxyAccountID.isZero() {
		body.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	return &proto.TransactionBody{
		TransactionID:            transaction.transactionID.toProtobuf(),
		TransactionFee:           transaction.transactionFee,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		Memo:                     transaction.Transaction.memo,
		Data: &proto.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: body,
		},
	}
}

// SetReceiverSignatureRequired sets the receiverSigRequired flag. If the receiverSigRequired flag is set to true, then
// all cryptocurrency transfers must be signed by this account's key, both for transfers in and out. If it is false,
// then only transfers out have to be signed by it. This transaction must be signed by the
// payer account. If receiverSigRequired is false, then the transaction does not have to be signed by the keys in the
// keys field. If it is true, then it must be signed by them, in addition to the keys of the payer account.
func (transaction *AccountCreateTransaction) SetReceiverSignatureRequired(required bool) *AccountCreateTransaction {
	transaction.receiverSignatureRequired = required
	return transaction
}

func (transaction *AccountCreateTransaction) GetReceiverSignatureRequired() bool {
	return transaction.receiverSignatureRequired
}

func (transaction *AccountCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.CryptoCreateTransactionBody{
		InitialBalance:         transaction.initialBalance,
		SendRecordThreshold:    transaction.receiveRecordThreshold,
		ReceiveRecordThreshold: transaction.sendRecordThreshold,
		ReceiverSigRequired:    transaction.receiverSignatureRequired,
		Memo:                   transaction.memo,
	}

	if transaction.key != nil {
		body.Key = transaction.key.toProtoKey()
	}

	if !transaction.proxyAccountID.isZero() {
		body.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: body,
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func accountCreateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().CreateAccount,
	}
}

func (transaction *AccountCreateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountCreateTransaction) Sign(
	privateKey PrivateKey,
) *AccountCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *AccountCreateTransaction) SignWithOperator(
	client *Client,
) (*AccountCreateTransaction, error) {
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
func (transaction *AccountCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountCreateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountCreateTransaction) Execute(
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
		accountCreateTransaction_getMethod,
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

func (transaction *AccountCreateTransaction) Freeze() (*AccountCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountCreateTransaction) FreezeWith(client *Client) (*AccountCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction.validateNetworkOnIDs(client)
	body := transaction.build()
	if err != nil {
		return &AccountCreateTransaction{}, err
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
}

func (transaction *AccountCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetMaxTransactionFee(fee Hbar) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *AccountCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionMemo(memo string) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *AccountCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionID(transactionID TransactionID) *AccountCreateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the node AccountID for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *AccountCreateTransaction) SetMaxRetry(count int) *AccountCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *AccountCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountCreateTransaction {
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

func (transaction *AccountCreateTransaction) SetMaxBackoff(max time.Duration) *AccountCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *AccountCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *AccountCreateTransaction) SetMinBackoff(min time.Duration) *AccountCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *AccountCreateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
