package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

type AccountUpdateTransaction struct {
	Transaction
	accountID                 AccountID
	proxyAccountID            AccountID
	key                       Key
	receiveRecordThreshold    uint64
	sendRecordThreshold       uint64
	autoRenewPeriod           *time.Duration
	memo                      string
	receiverSignatureRequired bool
	expirationTime            *time.Time
}

func NewAccountUpdateTransaction() *AccountUpdateTransaction {
	transaction := AccountUpdateTransaction{
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func accountUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) AccountUpdateTransaction {
	key, _ := keyFromProtobuf(pb.GetCryptoUpdateAccount().GetKey())
	var sendRecordThreshold uint64
	var receiveRecordThreshold uint64
	var receiverSignatureRequired bool

	switch s := pb.GetCryptoUpdateAccount().GetSendRecordThresholdField().(type) {
	case *proto.CryptoUpdateTransactionBody_SendRecordThreshold:
		sendRecordThreshold = s.SendRecordThreshold
	case *proto.CryptoUpdateTransactionBody_SendRecordThresholdWrapper:
		sendRecordThreshold = s.SendRecordThresholdWrapper.Value
	}

	switch s := pb.GetCryptoUpdateAccount().GetReceiveRecordThresholdField().(type) {
	case *proto.CryptoUpdateTransactionBody_ReceiveRecordThreshold:
		receiveRecordThreshold = s.ReceiveRecordThreshold
	case *proto.CryptoUpdateTransactionBody_ReceiveRecordThresholdWrapper:
		receiveRecordThreshold = s.ReceiveRecordThresholdWrapper.Value
	}

	switch s := pb.GetCryptoUpdateAccount().GetReceiverSigRequiredField().(type) {
	case *proto.CryptoUpdateTransactionBody_ReceiverSigRequired:
		receiverSignatureRequired = s.ReceiverSigRequired
	case *proto.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper:
		receiverSignatureRequired = s.ReceiverSigRequiredWrapper.Value
	}

	autoRenew := durationFromProtobuf(pb.GetCryptoUpdateAccount().AutoRenewPeriod)
	expiration := timeFromProtobuf(pb.GetCryptoUpdateAccount().ExpirationTime)

	return AccountUpdateTransaction{
		Transaction:               transaction,
		accountID:                 accountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetAccountIDToUpdate()),
		proxyAccountID:            accountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetProxyAccountID()),
		key:                       key,
		receiveRecordThreshold:    receiveRecordThreshold,
		sendRecordThreshold:       sendRecordThreshold,
		autoRenewPeriod:           &autoRenew,
		memo:                      pb.GetCryptoUpdateAccount().GetMemo().Value,
		receiverSignatureRequired: receiverSignatureRequired,
		expirationTime:            &expiration,
	}
}

//Sets the new key.
func (transaction *AccountUpdateTransaction) SetKey(key Key) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.key = key
	return transaction
}

func (transaction *AccountUpdateTransaction) GetKey() (Key, error) {
	return transaction.key, nil
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
	transaction.receiverSignatureRequired = receiverSignatureRequired
	return transaction
}

func (transaction *AccountUpdateTransaction) GetReceiverSignatureRequired() bool {
	return transaction.receiverSignatureRequired
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
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (transaction *AccountUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

func (transaction *AccountUpdateTransaction) SetExpirationTime(expirationTime time.Time) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.expirationTime = &expirationTime
	return transaction
}

//Sets the new expiration time to extend to (ignored if equal to or before the current one).
func (transaction *AccountUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}
	return time.Time{}
}

func (transaction *AccountUpdateTransaction) SetAccountMemo(memo string) *AccountUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.memo = memo

	return transaction
}

func (transaction *AccountUpdateTransaction) GeAccountMemo() string {
	return transaction.memo
}

func (transaction *AccountUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
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

func (transaction *AccountUpdateTransaction) build() *proto.TransactionBody {
	body := &proto.CryptoUpdateTransactionBody{
		SendRecordThresholdField: &proto.CryptoUpdateTransactionBody_SendRecordThreshold{
			SendRecordThreshold: transaction.sendRecordThreshold,
		},
		ReceiveRecordThresholdField: &proto.CryptoUpdateTransactionBody_ReceiveRecordThreshold{
			ReceiveRecordThreshold: transaction.receiveRecordThreshold,
		},
		ReceiverSigRequiredField: &proto.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper{
			ReceiverSigRequiredWrapper: &wrapperspb.BoolValue{Value: transaction.receiverSignatureRequired},
		},
		Memo: &wrapperspb.StringValue{Value: transaction.memo},
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = timeToProtobuf(*transaction.expirationTime)
	}

	if !transaction.accountID.isZero() {
		body.AccountIDToUpdate = transaction.accountID.toProtobuf()
	}

	if !transaction.proxyAccountID.isZero() {
		body.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	if transaction.key != nil {
		body.Key = transaction.key.toProtoKey()
	}

	pb := proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
		},
	}

	return &pb
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
	body := &proto.CryptoUpdateTransactionBody{
		SendRecordThresholdField: &proto.CryptoUpdateTransactionBody_SendRecordThreshold{
			SendRecordThreshold: transaction.sendRecordThreshold,
		},
		ReceiveRecordThresholdField: &proto.CryptoUpdateTransactionBody_ReceiveRecordThreshold{
			ReceiveRecordThreshold: transaction.receiveRecordThreshold,
		},
		ReceiverSigRequiredField: &proto.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper{
			ReceiverSigRequiredWrapper: &wrapperspb.BoolValue{Value: transaction.receiverSignatureRequired},
		},
		Memo: &wrapperspb.StringValue{Value: transaction.memo},
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = timeToProtobuf(*transaction.expirationTime)
	}

	if !transaction.accountID.isZero() {
		body.AccountIDToUpdate = transaction.accountID.toProtobuf()
	}

	if !transaction.proxyAccountID.isZero() {
		body.ProxyAccountID = transaction.proxyAccountID.toProtobuf()
	}

	if transaction.key != nil {
		body.Key = transaction.key.toProtoKey()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
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
		transaction_makeRequest(request{
			transaction: &transaction.Transaction,
		}),
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
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
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
