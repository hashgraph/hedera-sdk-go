package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Wipes the provided amount of tokens from the specified Account. Must be signed by the Token's
// Wipe key.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If Wipe Key is not present in the Token, transaction results in TOKEN_HAS_NO_WIPE_KEY.
// If the provided account is the Token's Treasury Account, transaction results in
// CANNOT_WIPE_TOKEN_TREASURY_ACCOUNT
// On success, tokens are removed from the account and the total supply of the token is decreased
// by the wiped amount.
//
// The amount provided is in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to wipe 100 tokens from account, one must provide amount of
// 10000. In order to wipe 100.55 tokens, one must provide amount of 10055.
type TokenWipeTransaction struct {
	Transaction
	tokenID   *TokenID
	accountID *AccountID
	amount    uint64
	serial    []int64
}

func NewTokenWipeTransaction() *TokenWipeTransaction {
	transaction := TokenWipeTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func _TokenWipeTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenWipeTransaction {
	return TokenWipeTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenWipe().GetToken()),
		accountID:   _AccountIDFromProtobuf(pb.GetTokenWipe().GetAccount()),
		amount:      pb.GetTokenWipe().Amount,
		serial:      pb.GetTokenWipe().GetSerialNumbers(),
	}
}

// The token for which the account will be wiped. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *TokenWipeTransaction) SetTokenID(tokenID TokenID) *TokenWipeTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

func (transaction *TokenWipeTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

// The account to be wiped
func (transaction *TokenWipeTransaction) SetAccountID(accountID AccountID) *TokenWipeTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

func (transaction *TokenWipeTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

// The amount of tokens to wipe from the specified account. Amount must be a positive non-zero
// number in the lowest denomination possible, not bigger than the token balance of the account
// (0; balance]
func (transaction *TokenWipeTransaction) SetAmount(amount uint64) *TokenWipeTransaction {
	transaction._RequireNotFrozen()
	transaction.amount = amount
	return transaction
}

func (transaction *TokenWipeTransaction) GetAmount() uint64 {
	return transaction.amount
}

func (transaction *TokenWipeTransaction) GetSerialNumbers() []int64 {
	return transaction.serial
}

func (transaction *TokenWipeTransaction) SetSerialNumbers(serial []int64) *TokenWipeTransaction {
	transaction._RequireNotFrozen()
	transaction.serial = serial
	return transaction
}

func (transaction *TokenWipeTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.accountID != nil {
		if err := transaction.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenWipeTransaction) _Build() *proto.TransactionBody {
	body := &proto.TokenWipeAccountTransactionBody{
		Amount: transaction.amount,
	}

	if len(transaction.serial) > 0 {
		body.SerialNumbers = transaction.serial
	}

	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &proto.TransactionBody_TokenWipe{
			TokenWipe: body,
		},
	}
}

func (transaction *TokenWipeTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenWipeTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.TokenWipeAccountTransactionBody{
		Amount: transaction.amount,
	}

	if len(transaction.serial) > 0 {
		body.SerialNumbers = transaction.serial
	}

	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenWipe{
			TokenWipe: body,
		},
	}, nil
}

func _TokenWipeTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().WipeTokenAccount,
	}
}

func (transaction *TokenWipeTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenWipeTransaction) Sign(
	privateKey PrivateKey,
) *TokenWipeTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenWipeTransaction) SignWithOperator(
	client *Client,
) (*TokenWipeTransaction, error) {
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
func (transaction *TokenWipeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenWipeTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenWipeTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
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
		_TokenWipeTransactionGetMethod,
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

func (transaction *TokenWipeTransaction) Freeze() (*TokenWipeTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenWipeTransaction) FreezeWith(client *Client) (*TokenWipeTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenWipeTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenWipeTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenWipeTransaction.
func (transaction *TokenWipeTransaction) SetMaxTransactionFee(fee Hbar) *TokenWipeTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenWipeTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenWipeTransaction.
func (transaction *TokenWipeTransaction) SetTransactionMemo(memo string) *TokenWipeTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenWipeTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenWipeTransaction.
func (transaction *TokenWipeTransaction) SetTransactionValidDuration(duration time.Duration) *TokenWipeTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenWipeTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenWipeTransaction.
func (transaction *TokenWipeTransaction) SetTransactionID(transactionID TransactionID) *TokenWipeTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenWipeTransaction.
func (transaction *TokenWipeTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenWipeTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenWipeTransaction) SetMaxRetry(count int) *TokenWipeTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenWipeTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenWipeTransaction {
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

func (transaction *TokenWipeTransaction) SetMaxBackoff(max time.Duration) *TokenWipeTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenWipeTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenWipeTransaction) SetMinBackoff(min time.Duration) *TokenWipeTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenWipeTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
