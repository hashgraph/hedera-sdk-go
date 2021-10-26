package hedera

import (
	"errors"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type LiveHashDeleteTransaction struct {
	Transaction
	accountID *AccountID
	hash      []byte
}

func NewLiveHashDeleteTransaction() *LiveHashDeleteTransaction {
	transaction := LiveHashDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _LiveHashDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) LiveHashDeleteTransaction {
	return LiveHashDeleteTransaction{
		Transaction: transaction,
		accountID:   _AccountIDFromProtobuf(pb.GetCryptoDeleteLiveHash().GetAccountOfLiveHash()),
		hash:        pb.GetCryptoDeleteLiveHash().LiveHashToDelete,
	}
}

func (transaction *LiveHashDeleteTransaction) SetHash(hash []byte) *LiveHashDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.hash = hash
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetHash() []byte {
	return transaction.hash
}

func (transaction *LiveHashDeleteTransaction) SetAccountID(accountID AccountID) *LiveHashDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

func (transaction *LiveHashDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
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

func (transaction *LiveHashDeleteTransaction) _Build() *proto.TransactionBody {
	body := &proto.CryptoDeleteLiveHashTransactionBody{}

	if transaction.accountID != nil {
		body.AccountOfLiveHash = transaction.accountID._ToProtobuf()
	}

	if transaction.hash != nil {
		body.LiveHashToDelete = transaction.hash
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &proto.TransactionBody_CryptoDeleteLiveHash{
			CryptoDeleteLiveHash: body,
		},
	}
}

func (transaction *LiveHashDeleteTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashAddTransaction`")
}

func _LiveHashDeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().DeleteLiveHash,
	}
}

func (transaction *LiveHashDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *LiveHashDeleteTransaction) Sign(
	privateKey PrivateKey,
) *LiveHashDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *LiveHashDeleteTransaction) SignWithOperator(
	client *Client,
) (*LiveHashDeleteTransaction, error) {
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
func (transaction *LiveHashDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *LiveHashDeleteTransaction) Execute(
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
		_TransactionMakeRequest(_Request{
			transaction: &transaction.Transaction,
		}),
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_LiveHashDeleteTransactionGetMethod,
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

func (transaction *LiveHashDeleteTransaction) Freeze() (*LiveHashDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *LiveHashDeleteTransaction) FreezeWith(client *Client) (*LiveHashDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &LiveHashDeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *LiveHashDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetMaxTransactionFee(fee Hbar) *LiveHashDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionMemo(memo string) *LiveHashDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionID(transactionID TransactionID) *LiveHashDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) SetMaxRetry(count int) *LiveHashDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *LiveHashDeleteTransaction {
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

func (transaction *LiveHashDeleteTransaction) SetMaxBackoff(max time.Duration) *LiveHashDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *LiveHashDeleteTransaction) SetMinBackoff(min time.Duration) *LiveHashDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
