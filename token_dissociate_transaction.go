package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TokenDissociateTransaction struct {
	Transaction
	accountID *AccountID
	tokens    []TokenID
}

func NewTokenDissociateTransaction() *TokenDissociateTransaction {
	transaction := TokenDissociateTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _TokenDissociateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) TokenDissociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenDissociate().Tokens {
		if tokenID := _TokenIDFromProtobuf(token); tokenID != nil {
			tokens = append(tokens, *tokenID)
		}
	}

	return TokenDissociateTransaction{
		Transaction: transaction,
		accountID:   _AccountIDFromProtobuf(pb.GetTokenDissociate().GetAccount()),
		tokens:      tokens,
	}
}

// The account to be dissociated with the provided tokens
func (transaction *TokenDissociateTransaction) SetAccountID(accountID AccountID) *TokenDissociateTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

func (transaction *TokenDissociateTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

// The tokens to be dissociated with the provided account
func (transaction *TokenDissociateTransaction) SetTokenIDs(ids ...TokenID) *TokenDissociateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokens = make([]TokenID, len(ids))

	for i, tokenID := range ids {
		transaction.tokens[i] = tokenID
	}

	return transaction
}

func (transaction *TokenDissociateTransaction) AddTokenID(id TokenID) *TokenDissociateTransaction {
	transaction._RequireNotFrozen()
	if transaction.tokens == nil {
		transaction.tokens = make([]TokenID, 0)
	}

	transaction.tokens = append(transaction.tokens, id)

	return transaction
}

func (transaction *TokenDissociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(transaction.tokens))

	for i, tokenID := range transaction.tokens {
		tokenIDs[i] = tokenID
	}

	return tokenIDs
}

func (transaction *TokenDissociateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.accountID != nil {
		if err := transaction.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, tokenID := range transaction.tokens {
		if err := tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenDissociateTransaction) _Build() *services.TransactionBody {
	body := &services.TokenDissociateTransactionBody{}
	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}

	if len(transaction.tokens) > 0 {
		for _, tokenID := range transaction.tokens {
			if body.Tokens == nil {
				body.Tokens = make([]*services.TokenID, 0)
			}
			body.Tokens = append(body.Tokens, tokenID._ToProtobuf())
		}
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenDissociate{
			TokenDissociate: body,
		},
	}
}

func (transaction *TokenDissociateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenDissociateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.TokenDissociateTransactionBody{}
	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}

	if len(transaction.tokens) > 0 {
		for _, tokenID := range transaction.tokens {
			if body.Tokens == nil {
				body.Tokens = make([]*services.TokenID, 0)
			}
			body.Tokens = append(body.Tokens, tokenID._ToProtobuf())
		}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenDissociate{
			TokenDissociate: body,
		},
	}, nil
}

func _TokenDissociateTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DissociateTokens,
	}
}

func (transaction *TokenDissociateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenDissociateTransaction) Sign(
	privateKey PrivateKey,
) *TokenDissociateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenDissociateTransaction) SignWithOperator(
	client *Client,
) (*TokenDissociateTransaction, error) {
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
func (transaction *TokenDissociateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenDissociateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenDissociateTransaction) Execute(
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
		_TokenDissociateTransactionGetMethod,
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

func (transaction *TokenDissociateTransaction) Freeze() (*TokenDissociateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenDissociateTransaction) FreezeWith(client *Client) (*TokenDissociateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenDissociateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenDissociateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetMaxTransactionFee(fee Hbar) *TokenDissociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenDissociateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetTransactionMemo(memo string) *TokenDissociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenDissociateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenDissociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenDissociateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetTransactionID(transactionID TransactionID) *TokenDissociateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenDissociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenDissociateTransaction) SetMaxRetry(count int) *TokenDissociateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenDissociateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenDissociateTransaction {
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

func (transaction *TokenDissociateTransaction) SetMaxBackoff(max time.Duration) *TokenDissociateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenDissociateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenDissociateTransaction) SetMinBackoff(min time.Duration) *TokenDissociateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenDissociateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
