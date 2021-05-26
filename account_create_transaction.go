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
	pb *proto.CryptoCreateTransactionBody
}

// NewAccountCreateTransaction creates an AccountCreateTransaction transaction which can be used to construct and
// execute a Crypto Create Transaction.
func NewAccountCreateTransaction() *AccountCreateTransaction {
	pb := &proto.CryptoCreateTransactionBody{}

	transaction := AccountCreateTransaction{
		pb:          pb,
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
	return AccountCreateTransaction{
		Transaction: transaction,
		pb:          pb.GetCryptoCreateAccount(),
	}
}

// SetKey sets the key that must sign each transfer out of the account. If RecieverSignatureRequired is true, then it
// must also sign any transfer into the account.
func (transaction *AccountCreateTransaction) SetKey(key Key) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Key = key.toProtoKey()
	return transaction
}

func (transaction *AccountCreateTransaction) GetKey() (Key, error) {
	return keyFromProtobuf(transaction.pb.GetKey())
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (transaction *AccountCreateTransaction) SetInitialBalance(initialBalance Hbar) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.InitialBalance = uint64(initialBalance.AsTinybar())
	return transaction
}

func (transaction *AccountCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(int64(transaction.pb.GetInitialBalance()))
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (transaction *AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = durationToProtobuf(autoRenewPeriod)
	return transaction
}

func (transaction *AccountCreateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProtobuf(transaction.pb.AutoRenewPeriod)
}

// SetSendRecordThreshold sets the threshold amount for which an account record is created for any send/withdraw
// transaction
//
// Deprecated: No longer used by Hedera
func (transaction *AccountCreateTransaction) setSendRecordThreshold(recordThreshold Hbar) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SendRecordThreshold = uint64(recordThreshold.AsTinybar())
	return transaction
}

func (transaction *AccountCreateTransaction) getSendRecordThreshold() Hbar {
	return HbarFromTinybar(int64(transaction.pb.GetSendRecordThreshold()))
}

// SetReceiveRecordThreshold sets the threshold amount for which an account record is created for any receive/deposit
// transaction
//
// Deprecated: No longer used by Hedera
func (transaction *AccountCreateTransaction) setReceiveRecordThreshold(recordThreshold Hbar) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ReceiveRecordThreshold = uint64(recordThreshold.AsTinybar())
	return transaction
}

func (transaction *AccountCreateTransaction) getReceiveRecordThreshold() Hbar {
	return HbarFromTinybar(int64(transaction.pb.GetReceiveRecordThreshold()))
}

// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalid account, or is an account that isn't a node, then this account is automatically proxy staked to a node
// chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a node, then it will behave as if proxyAccountID was not set.
func (transaction *AccountCreateTransaction) SetProxyAccountID(id AccountID) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ProxyAccountID = id.toProtobuf()
	return transaction
}

func (transaction *AccountCreateTransaction) GetProxyAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetProxyAccountID())
}

func (transaction *AccountCreateTransaction) SetAccountMemo(memo string) *AccountCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = memo
	return transaction
}

func (transaction *AccountCreateTransaction) GetAccountMemo() string {
	return transaction.pb.GetMemo()
}

// SetReceiverSignatureRequired sets the receiverSigRequired flag. If the receiverSigRequired flag is set to true, then
// all cryptocurrency transfers must be signed by this account's key, both for transfers in and out. If it is false,
// then only transfers out have to be signed by it. This transaction must be signed by the
// payer account. If receiverSigRequired is false, then the transaction does not have to be signed by the keys in the
// keys field. If it is true, then it must be signed by them, in addition to the keys of the payer account.
func (transaction *AccountCreateTransaction) SetReceiverSignatureRequired(required bool) *AccountCreateTransaction {
	transaction.pb.ReceiverSigRequired = required
	return transaction
}

func (transaction *AccountCreateTransaction) GetReceiverSignatureRequired() bool {
	return transaction.pb.GetReceiverSigRequired()
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
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: &proto.CryptoCreateTransactionBody{
				Key:                    transaction.pb.GetKey(),
				InitialBalance:         transaction.pb.GetInitialBalance(),
				ProxyAccountID:         transaction.pb.GetProxyAccountID(),
				SendRecordThreshold:    transaction.pb.GetSendRecordThreshold(),
				ReceiveRecordThreshold: transaction.pb.GetReceiveRecordThreshold(),
				ReceiverSigRequired:    transaction.pb.GetReceiverSigRequired(),
				AutoRenewPeriod:        transaction.pb.GetAutoRenewPeriod(),
				ShardID:                transaction.pb.GetShardID(),
				RealmID:                transaction.pb.GetRealmID(),
				NewRealmAdminKey:       transaction.pb.GetNewRealmAdminKey(),
				Memo:                   transaction.pb.GetMemo(),
			},
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
		transaction_makeRequest,
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

func (transaction *AccountCreateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoCreateAccount{
		CryptoCreateAccount: transaction.pb,
	}

	return true
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

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
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
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
