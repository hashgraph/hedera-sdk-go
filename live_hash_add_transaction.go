package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"

	"time"
)

type LiveHashAddTransaction struct {
	Transaction
	accountID *AccountID
	hash      []byte
	keys      *KeyList
	duration  *time.Duration
}

func NewLiveHashAddTransaction() *LiveHashAddTransaction {
	transaction := LiveHashAddTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _LiveHashAddTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) LiveHashAddTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.GetKeys())
	duration := _DurationFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.Duration)

	return LiveHashAddTransaction{
		Transaction: transaction,
		accountID:   _AccountIDFromProtobuf(pb.GetCryptoAddLiveHash().GetLiveHash().GetAccountId()),
		hash:        pb.GetCryptoAddLiveHash().LiveHash.Hash,
		keys:        &keys,
		duration:    &duration,
	}
}

func (transaction *LiveHashAddTransaction) SetHash(hash []byte) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.hash = hash
	return transaction
}

func (transaction *LiveHashAddTransaction) GetHash() []byte {
	return transaction.hash
}

func (transaction *LiveHashAddTransaction) SetKeys(keys ...Key) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	if transaction.keys == nil {
		transaction.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.keys = keyList

	return transaction
}

func (transaction *LiveHashAddTransaction) GetKeys() KeyList {
	if transaction.keys != nil {
		return *transaction.keys
	}

	return KeyList{}
}

func (transaction *LiveHashAddTransaction) SetDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.duration = &duration
	return transaction
}

func (transaction *LiveHashAddTransaction) GetDuration() time.Duration {
	if transaction.duration != nil {
		return *transaction.duration
	}

	return time.Duration(0)
}

func (transaction *LiveHashAddTransaction) SetAccountID(accountID AccountID) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

func (transaction *LiveHashAddTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

func (transaction *LiveHashAddTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.accountID != nil {
		if err := transaction.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *LiveHashAddTransaction) _Build() *services.TransactionBody {
	body := &services.CryptoAddLiveHashTransactionBody{
		LiveHash: &services.LiveHash{},
	}

	if transaction.accountID != nil {
		body.LiveHash.AccountId = transaction.accountID._ToProtobuf()
	}

	if transaction.duration != nil {
		body.LiveHash.Duration = _DurationToProtobuf(*transaction.duration)
	}

	if transaction.keys != nil {
		body.LiveHash.Keys = transaction.keys._ToProtoKeyList()
	}

	if transaction.hash != nil {
		body.LiveHash.Hash = transaction.hash
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoAddLiveHash{
			CryptoAddLiveHash: body,
		},
	}
}

func (transaction *LiveHashAddTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashAddTransaction`")
}

func _LiveHashAddTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().AddLiveHash,
	}
}

func (transaction *LiveHashAddTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *LiveHashAddTransaction) Sign(
	privateKey PrivateKey,
) *LiveHashAddTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *LiveHashAddTransaction) SignWithOperator(
	client *Client,
) (*LiveHashAddTransaction, error) {
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
func (transaction *LiveHashAddTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashAddTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *LiveHashAddTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
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
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_LiveHashAddTransactionGetMethod,
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

func (transaction *LiveHashAddTransaction) Freeze() (*LiveHashAddTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *LiveHashAddTransaction) FreezeWith(client *Client) (*LiveHashAddTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &LiveHashAddTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *LiveHashAddTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetMaxTransactionFee(fee Hbar) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionMemo(memo string) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionID(transactionID TransactionID) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *LiveHashAddTransaction) SetMaxRetry(count int) *LiveHashAddTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *LiveHashAddTransaction) AddSignature(publicKey PublicKey, signature []byte) *LiveHashAddTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if len(transaction.signedTransactions) == 0 {
		return transaction
	}

	transaction.transactions = make([]*services.Transaction, 0)
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

func (transaction *LiveHashAddTransaction) SetMaxBackoff(max time.Duration) *LiveHashAddTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *LiveHashAddTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *LiveHashAddTransaction) SetMinBackoff(min time.Duration) *LiveHashAddTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *LiveHashAddTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
