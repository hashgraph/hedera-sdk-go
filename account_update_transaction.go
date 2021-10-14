package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type AccountUpdateTransaction struct {
	Transaction
	accountID                     *AccountID
	proxyAccountID                *AccountID
	key                           Key
	receiveRecordThreshold        uint64
	sendRecordThreshold           uint64
	autoRenewPeriod               *time.Duration
	memo                          string
	receiverSignatureRequired     bool
	expirationTime                *time.Time
	maxAutomaticTokenAssociations uint32
}

func NewAccountUpdateTransaction() *AccountUpdateTransaction {
	transaction := AccountUpdateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _AccountUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) AccountUpdateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetCryptoUpdateAccount().GetKey())
	var sendRecordThreshold uint64
	var receiveRecordThreshold uint64
	var receiverSignatureRequired bool

	switch s := pb.GetCryptoUpdateAccount().GetSendRecordThresholdField().(type) {
	case *proto.CryptoUpdateTransactionBody_SendRecordThreshold:
		sendRecordThreshold = s.SendRecordThreshold // nolint
	case *proto.CryptoUpdateTransactionBody_SendRecordThresholdWrapper:
		sendRecordThreshold = s.SendRecordThresholdWrapper.Value // nolint
	}

	switch s := pb.GetCryptoUpdateAccount().GetReceiveRecordThresholdField().(type) {
	case *proto.CryptoUpdateTransactionBody_ReceiveRecordThreshold:
		receiveRecordThreshold = s.ReceiveRecordThreshold // nolint
	case *proto.CryptoUpdateTransactionBody_ReceiveRecordThresholdWrapper:
		receiveRecordThreshold = s.ReceiveRecordThresholdWrapper.Value // nolint
	}

	switch s := pb.GetCryptoUpdateAccount().GetReceiverSigRequiredField().(type) {
	case *proto.CryptoUpdateTransactionBody_ReceiverSigRequired:
		receiverSignatureRequired = s.ReceiverSigRequired // nolint
	case *proto.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper:
		receiverSignatureRequired = s.ReceiverSigRequiredWrapper.Value // nolint
	}

	autoRenew := _DurationFromProtobuf(pb.GetCryptoUpdateAccount().AutoRenewPeriod)
	expiration := _TimeFromProtobuf(pb.GetCryptoUpdateAccount().ExpirationTime)

	return AccountUpdateTransaction{
		Transaction:               transaction,
		accountID:                 _AccountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetAccountIDToUpdate()),
		proxyAccountID:            _AccountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetProxyAccountID()),
		key:                       key,
		receiveRecordThreshold:    receiveRecordThreshold,
		sendRecordThreshold:       sendRecordThreshold,
		autoRenewPeriod:           &autoRenew,
		memo:                      pb.GetCryptoUpdateAccount().GetMemo().Value,
		receiverSignatureRequired: receiverSignatureRequired,
		expirationTime:            &expiration,
	}
}

// Sets the new key.
func (transaction *AccountUpdateTransaction) SetKey(key Key) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.key = key
	return transaction
}

func (transaction *AccountUpdateTransaction) GetKey() (Key, error) {
	return transaction.key, nil
}

// Sets the account ID which is being updated in this transaction.
func (transaction *AccountUpdateTransaction) SetAccountID(accountID AccountID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

func (transaction *AccountUpdateTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

func (transaction *AccountUpdateTransaction) SetMaxAutomaticTokenAssociations(max uint32) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.maxAutomaticTokenAssociations = max
	return transaction
}

func (transaction *AccountUpdateTransaction) GetMaxAutomaticTokenAssociations() uint32 {
	return transaction.maxAutomaticTokenAssociations
}

func (transaction *AccountUpdateTransaction) SetReceiverSignatureRequired(receiverSignatureRequired bool) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.receiverSignatureRequired = receiverSignatureRequired
	return transaction
}

func (transaction *AccountUpdateTransaction) GetReceiverSignatureRequired() bool {
	return transaction.receiverSignatureRequired
}

// Sets the ID of the account to which this account is proxy staked.
func (transaction *AccountUpdateTransaction) SetProxyAccountID(proxyAccountID AccountID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.proxyAccountID = &proxyAccountID
	return transaction
}

func (transaction *AccountUpdateTransaction) GetProxyAccountID() AccountID {
	if transaction.proxyAccountID == nil {
		return AccountID{}
	}

	return *transaction.proxyAccountID
}

// Sets the duration in which it will automatically extend the expiration period.
func (transaction *AccountUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
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
	transaction._RequireNotFrozen()
	transaction.expirationTime = &expirationTime
	return transaction
}

// Sets the new expiration time to extend to (ignored if equal to or before the current one).
func (transaction *AccountUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}
	return time.Time{}
}

func (transaction *AccountUpdateTransaction) SetAccountMemo(memo string) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo

	return transaction
}

func (transaction *AccountUpdateTransaction) GeAccountMemo() string {
	return transaction.memo
}

func (transaction *AccountUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.accountID != nil {
		if err := transaction.accountID.Validate(client); err != nil {
			return err
		}
	}

	if transaction.proxyAccountID != nil {
		if err := transaction.proxyAccountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *AccountUpdateTransaction) _Build() *proto.TransactionBody {
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
		Memo:                          &wrapperspb.StringValue{Value: transaction.memo},
		MaxAutomaticTokenAssociations: &wrapperspb.Int32Value{Value: int32(transaction.maxAutomaticTokenAssociations)},
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.accountID != nil {
		body.AccountIDToUpdate = transaction.accountID._ToProtobuf()
	}

	if transaction.proxyAccountID != nil {
		body.ProxyAccountID = transaction.proxyAccountID._ToProtobuf()
	}

	if transaction.key != nil {
		body.Key = transaction.key._ToProtoKey()
	}

	pb := proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &proto.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
		},
	}

	body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: int32(transaction.maxAutomaticTokenAssociations)}

	return &pb
}

func (transaction *AccountUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountUpdateTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
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
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.accountID != nil {
		body.AccountIDToUpdate = transaction.accountID._ToProtobuf()
	}

	if transaction.proxyAccountID != nil {
		body.ProxyAccountID = transaction.proxyAccountID._ToProtobuf()
	}

	if transaction.key != nil {
		body.Key = transaction.key._ToProtoKey()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
		},
	}, nil
}

func _AccountUpdateTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().UpdateAccount,
	}
}

func (transaction *AccountUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
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
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

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
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
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

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest(_Request{
			transaction: &transaction.Transaction,
		}),
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_AccountUpdateTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()
	if err != nil {
		return TransactionResponse{}, err
	}

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
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &AccountUpdateTransaction{}, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *AccountUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetMaxTransactionFee(fee Hbar) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionMemo(memo string) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *AccountUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionID(transactionID TransactionID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *AccountUpdateTransaction) SetMaxRetry(count int) *AccountUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *AccountUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountUpdateTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
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
			publicKey._ToSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

func (transaction *AccountUpdateTransaction) SetMaxBackoff(max time.Duration) *AccountUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *AccountUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *AccountUpdateTransaction) SetMinBackoff(min time.Duration) *AccountUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *AccountUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
