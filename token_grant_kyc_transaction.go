package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// Grants KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Granted.
type TokenGrantKycTransaction struct {
	Transaction
	tokenID   *TokenID
	accountID *AccountID
}

func NewTokenGrantKycTransaction() *TokenGrantKycTransaction {
	transaction := TokenGrantKycTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func _TokenGrantKycTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenGrantKycTransaction {
	return &TokenGrantKycTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenGrantKyc().GetToken()),
		accountID:   _AccountIDFromProtobuf(pb.GetTokenGrantKyc().GetAccount()),
	}
}

func (transaction *TokenGrantKycTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenGrantKycTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// The token for which this account will be granted KYC. If token does not exist, transaction results in INVALID_TOKEN_ID
func (transaction *TokenGrantKycTransaction) SetTokenID(tokenID TokenID) *TokenGrantKycTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

// The account to be KYCed
func (transaction *TokenGrantKycTransaction) SetAccountID(accountID AccountID) *TokenGrantKycTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

func (transaction *TokenGrantKycTransaction) _ValidateNetworkOnIDs(client *Client) error {
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

func (transaction *TokenGrantKycTransaction) _Build() *services.TransactionBody {
	body := &services.TokenGrantKycTransactionBody{}
	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: body,
		},
	}
}

func (transaction *TokenGrantKycTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenGrantKycTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.TokenGrantKycTransactionBody{}
	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenGrantKyc{
			TokenGrantKyc: body,
		},
	}, nil
}

func _TokenGrantKycTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().GrantKycToTokenAccount,
	}
}

func (transaction *TokenGrantKycTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenGrantKycTransaction) Sign(
	privateKey PrivateKey,
) *TokenGrantKycTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenGrantKycTransaction) SignWithOperator(
	client *Client,
) (*TokenGrantKycTransaction, error) {
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
func (transaction *TokenGrantKycTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenGrantKycTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenGrantKycTransaction) Execute(
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

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

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
		_TokenGrantKycTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
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

func (transaction *TokenGrantKycTransaction) Freeze() (*TokenGrantKycTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenGrantKycTransaction) FreezeWith(client *Client) (*TokenGrantKycTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenGrantKycTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenGrantKycTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetMaxTransactionFee(fee Hbar) *TokenGrantKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TokenGrantKycTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenGrantKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TokenGrantKycTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *TokenGrantKycTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionMemo(memo string) *TokenGrantKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionValidDuration(duration time.Duration) *TokenGrantKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionID(transactionID TransactionID) *TokenGrantKycTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenGrantKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenGrantKycTransaction) SetMaxRetry(count int) *TokenGrantKycTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenGrantKycTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenGrantKycTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := transaction.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
	}

	return transaction
}

func (transaction *TokenGrantKycTransaction) SetMaxBackoff(max time.Duration) *TokenGrantKycTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenGrantKycTransaction) SetMinBackoff(min time.Duration) *TokenGrantKycTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TokenGrantKycTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TokenGrantKycTransaction:%d", timestamp.UnixNano())
}
