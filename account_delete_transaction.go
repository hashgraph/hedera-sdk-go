package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountDeleteTransaction creates a new account. After the account is created, the AccountID for it is in the receipt,
// or by asking for a Record of the transaction to be created, and retrieving that. The account can then automatically
// generate records for large transfers into it or out of it, which each last for 25 hours. Records are generated for
// any transfer that exceeds the thresholds given here. This account is charged hbar for each record generated, so the
// thresholds are useful for limiting Record generation to happen only for large transactions.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0,
// with a null key. Future versions of the API will support multiple realms and multiple shards.
type AccountDeleteTransaction struct {
	Transaction
	transferAccountID *AccountID
	deleteAccountID   *AccountID
}

func _AccountDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *AccountDeleteTransaction {
	return &AccountDeleteTransaction{
		Transaction:       transaction,
		transferAccountID: _AccountIDFromProtobuf(pb.GetCryptoDelete().GetTransferAccountID()),
		deleteAccountID:   _AccountIDFromProtobuf(pb.GetCryptoDelete().GetDeleteAccountID()),
	}
}

func NewAccountDeleteTransaction() *AccountDeleteTransaction {
	transaction := AccountDeleteTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

// SetNodeAccountID sets the _Node AccountID for this AccountCreateTransaction.
func (transaction *AccountDeleteTransaction) SetAccountID(accountID AccountID) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.deleteAccountID = &accountID
	return transaction
}

func (transaction *AccountDeleteTransaction) GetAccountID() AccountID {
	if transaction.deleteAccountID == nil {
		return AccountID{}
	}

	return *transaction.deleteAccountID
}

// SetTransferAccountID sets the AccountID which will receive all remaining hbars.
func (transaction *AccountDeleteTransaction) SetTransferAccountID(transferAccountID AccountID) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.transferAccountID = &transferAccountID
	return transaction
}

func (transaction *AccountDeleteTransaction) GetTransferAccountID(transferAccountID AccountID) AccountID {
	if transaction.transferAccountID == nil {
		return AccountID{}
	}

	return *transaction.transferAccountID
}

func (transaction *AccountDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.deleteAccountID != nil {
		if err := transaction.deleteAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.transferAccountID != nil {
		if err := transaction.transferAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *AccountDeleteTransaction) _Build() *services.TransactionBody {
	body := &services.CryptoDeleteTransactionBody{}

	if transaction.transferAccountID != nil {
		body.TransferAccountID = transaction.transferAccountID._ToProtobuf()
	}

	if transaction.deleteAccountID != nil {
		body.DeleteAccountID = transaction.deleteAccountID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoDelete{
			CryptoDelete: body,
		},
	}
}

func (transaction *AccountDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.CryptoDeleteTransactionBody{}

	if transaction.transferAccountID != nil {
		body.TransferAccountID = transaction.transferAccountID._ToProtobuf()
	}

	if transaction.deleteAccountID != nil {
		body.DeleteAccountID = transaction.deleteAccountID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoDelete{
			CryptoDelete: body,
		},
	}, nil
}

func _AccountDeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().CryptoDelete,
	}
}

func (transaction *AccountDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountDeleteTransaction) Sign(
	privateKey PrivateKey,
) *AccountDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *AccountDeleteTransaction) SignWithOperator(
	client *Client,
) (*AccountDeleteTransaction, error) {
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
func (transaction *AccountDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountDeleteTransaction) Execute(
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
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_AccountDeleteTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
			Hash:          make([]byte, 0),
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

func (transaction *AccountDeleteTransaction) Freeze() (*AccountDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountDeleteTransaction) FreezeWith(client *Client) (*AccountDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &AccountDeleteTransaction{}, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *AccountDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetMaxTransactionFee(fee Hbar) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *AccountDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionMemo(memo string) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *AccountDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionID(transactionID TransactionID) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *AccountDeleteTransaction) SetMaxRetry(count int) *AccountDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *AccountDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountDeleteTransaction {
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

func (transaction *AccountDeleteTransaction) SetMaxBackoff(max time.Duration) *AccountDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *AccountDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *AccountDeleteTransaction) SetMinBackoff(min time.Duration) *AccountDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *AccountDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
