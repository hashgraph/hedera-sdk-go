package hedera

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type AccountUpdateTransaction struct {
	Transaction
	pb             *proto.CryptoUpdateTransactionBody
	accountID      AccountID
	proxyAccountID AccountID
}

func NewAccountUpdateTransaction() *AccountUpdateTransaction {
	pb := &proto.CryptoUpdateTransactionBody{}

	transaction := AccountUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func accountUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) AccountUpdateTransaction {
	return AccountUpdateTransaction{
		Transaction:    transaction,
		pb:             pb.GetCryptoUpdateAccount(),
		accountID:      accountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetAccountIDToUpdate(), nil),
		proxyAccountID: accountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetProxyAccountID(), nil),
	}
}

//Sets the new key.
func (transaction *AccountUpdateTransaction) SetKey(key Key) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Key = key.toProtoKey()
	return transaction
}

func (transaction *AccountUpdateTransaction) GetKey() (Key, error) {
	return keyFromProtobuf(transaction.pb.GetKey(), nil)
}

//Sets the account ID which is being updated in this transaction.
func (transaction *AccountUpdateTransaction) SetAccountID(id AccountID) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.accountID = id
	return transaction
}

func (transaction *AccountUpdateTransaction) GetAccountID() AccountID {
	return transaction.accountID
}

func (transaction *AccountUpdateTransaction) SetReceiverSignatureRequired(receiverSignatureRequired bool) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.GetReceiverSigRequiredWrapper().Value = receiverSignatureRequired
	return transaction
}

func (transaction *AccountUpdateTransaction) GetReceiverSignatureRequired() bool {
	return transaction.pb.GetReceiverSigRequiredWrapper().GetValue()
}

//Sets the ID of the account to which this account is proxy staked.
func (transaction *AccountUpdateTransaction) SetProxyAccountID(proxyAccountID AccountID) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.proxyAccountID = proxyAccountID
	return transaction
}

func (transaction *AccountUpdateTransaction) GetProxyAccountID() AccountID {
	return transaction.proxyAccountID
}

//Sets the duration in which it will automatically extend the expiration period.
func (transaction *AccountUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = durationToProtobuf(autoRenewPeriod)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProtobuf(transaction.pb.GetAutoRenewPeriod())
}

func (transaction *AccountUpdateTransaction) SetExpirationTime(expirationTime time.Time) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ExpirationTime = timeToProtobuf(expirationTime)
	return transaction
}

//Sets the new expiration time to extend to (ignored if equal to or before the current one).
func (transaction *AccountUpdateTransaction) GetExpirationTime() time.Time {
	return timeFromProtobuf(transaction.pb.ExpirationTime)
}

func (transaction *AccountUpdateTransaction) SetAccountMemo(memo string) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = &wrappers.StringValue{Value: memo}

	return transaction
}

func (transaction *AccountUpdateTransaction) GeAccountMemo() string {
	if transaction.pb.Memo != nil {
		return transaction.pb.Memo.GetValue()
	}

	return ""
}

func (transaction *AccountUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.accountID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.proxyAccountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *AccountUpdateTransaction) build() *AccountUpdateTransaction {
	if !transaction.accountID.isZero() {
		transaction.pb.AccountIDToUpdate = transaction.accountID.toProtobuf()
	}

	if !transaction.proxyAccountID.isZero() {
		transaction.pb.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	return transaction
}

func (transaction *AccountUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountUpdateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: &proto.CryptoUpdateTransactionBody{
				AccountIDToUpdate:           transaction.pb.GetAccountIDToUpdate(),
				Key:                         transaction.pb.GetKey(),
				ProxyAccountID:              transaction.pb.GetProxyAccountID(),
				ProxyFraction:               transaction.pb.GetProxyFraction(),
				SendRecordThresholdField:    transaction.pb.GetSendRecordThresholdField(),
				ReceiveRecordThresholdField: transaction.pb.GetReceiveRecordThresholdField(),
				AutoRenewPeriod:             transaction.pb.GetAutoRenewPeriod(),
				ExpirationTime:              transaction.pb.GetExpirationTime(),
				ReceiverSigRequiredField:    transaction.pb.GetReceiverSigRequiredField(),
				Memo:                        transaction.pb.GetMemo(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func accountUpdateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().UpdateAccount,
	}
}

func (transaction *AccountUpdateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountUpdateTransaction) Sign(
	privateKey PrivateKey,
) *AccountUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *AccountUpdateTransaction) SignWithOperator(
	client *Client,
) (*AccountUpdateTransaction, error) {

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
func (transaction *AccountUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountUpdateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountUpdateTransaction) Execute(
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
		accountUpdateTransaction_getMethod,
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

func (transaction *AccountUpdateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoUpdateAccount{
		CryptoUpdateAccount: transaction.pb,
	}

	return true
}

func (transaction *AccountUpdateTransaction) Freeze() (*AccountUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountUpdateTransaction) FreezeWith(client *Client) (*AccountUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &AccountUpdateTransaction{}, err
	}
	transaction.build()

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *AccountUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetMaxTransactionFee(fee Hbar) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionMemo(memo string) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionID(transactionID TransactionID) *AccountUpdateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the node AccountID for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *AccountUpdateTransaction) SetMaxRetry(count int) *AccountUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *AccountUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountUpdateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
