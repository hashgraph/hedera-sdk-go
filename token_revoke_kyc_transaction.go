package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Revokes KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Revoked
type TokenRevokeKycTransaction struct {
	Transaction
	tokenID   *TokenID
	accountID *AccountID
}

func NewTokenRevokeKycTransaction() *TokenRevokeKycTransaction {
	transaction := TokenRevokeKycTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func _TokenRevokeKycTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenRevokeKycTransaction {
	return TokenRevokeKycTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenRevokeKyc().GetToken()),
		accountID:   _AccountIDFromProtobuf(pb.GetTokenRevokeKyc().GetAccount()),
	}
}

// The token for which this account will get his KYC revoked. If token does not exist, transaction results in INVALID_TOKEN_ID
func (transaction *TokenRevokeKycTransaction) SetTokenID(tokenID TokenID) *TokenRevokeKycTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

func (transaction *TokenRevokeKycTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

// The account to be KYC Revoked
func (transaction *TokenRevokeKycTransaction) SetAccountID(accountID AccountID) *TokenRevokeKycTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

func (transaction *TokenRevokeKycTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

func (transaction *TokenRevokeKycTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.Validate(client); err != nil {
			return err
		}
	}

	if transaction.accountID != nil {
		if err := transaction.accountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenRevokeKycTransaction) _Build() *proto.TransactionBody {
	body := &proto.TokenRevokeKycTransactionBody{}
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
		Data: &proto.TransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: body,
		},
	}
}

func (transaction *TokenRevokeKycTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenRevokeKycTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.TokenRevokeKycTransactionBody{}
	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: body,
		},
	}, nil
}

func _TokenRevokeKycTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().RevokeKycFromTokenAccount,
	}
}

func (transaction *TokenRevokeKycTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenRevokeKycTransaction) Sign(
	privateKey PrivateKey,
) *TokenRevokeKycTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenRevokeKycTransaction) SignWithOperator(
	client *Client,
) (*TokenRevokeKycTransaction, error) {
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
func (transaction *TokenRevokeKycTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenRevokeKycTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenRevokeKycTransaction) Execute(
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
		_TokenRevokeKycTransactionGetMethod,
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

func (transaction *TokenRevokeKycTransaction) Freeze() (*TokenRevokeKycTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenRevokeKycTransaction) FreezeWith(client *Client) (*TokenRevokeKycTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenRevokeKycTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenRevokeKycTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenRevokeKycTransaction.
func (transaction *TokenRevokeKycTransaction) SetMaxTransactionFee(fee Hbar) *TokenRevokeKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenRevokeKycTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenRevokeKycTransaction.
func (transaction *TokenRevokeKycTransaction) SetTransactionMemo(memo string) *TokenRevokeKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenRevokeKycTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenRevokeKycTransaction.
func (transaction *TokenRevokeKycTransaction) SetTransactionValidDuration(duration time.Duration) *TokenRevokeKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenRevokeKycTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenRevokeKycTransaction.
func (transaction *TokenRevokeKycTransaction) SetTransactionID(transactionID TransactionID) *TokenRevokeKycTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenRevokeKycTransaction.
func (transaction *TokenRevokeKycTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenRevokeKycTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenRevokeKycTransaction) SetMaxRetry(count int) *TokenRevokeKycTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenRevokeKycTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenRevokeKycTransaction {
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

func (transaction *TokenRevokeKycTransaction) SetMaxBackoff(max time.Duration) *TokenRevokeKycTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenRevokeKycTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenRevokeKycTransaction) SetMinBackoff(min time.Duration) *TokenRevokeKycTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenRevokeKycTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
