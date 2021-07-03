package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Associates the provided account with the provided tokens. Must be signed by the provided Account's key.
// If the provided account is not found, the transaction will resolve to
// INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to
// ACCOUNT_DELETED.
// If any of the provided tokens is not found, the transaction will resolve to
// INVALID_TOKEN_REF.
// If any of the provided tokens has been deleted, the transaction will resolve to
// TOKEN_WAS_DELETED.
// If an association between the provided account and any of the tokens already exists, the
// transaction will resolve to
// TOKEN_ALREADY_ASSOCIATED_TO_ACCOUNT.
// If the provided account's associations count exceed the constraint of maximum token
// associations per account, the transaction will resolve to
// TOKENS_PER_ACCOUNT_LIMIT_EXCEEDED.
// On success, associations between the provided account and tokens are made and the account is
// ready to interact with the tokens.
type TokenAssociateTransaction struct {
	Transaction
	pb        *proto.TokenAssociateTransactionBody
	accountID AccountID
	tokens    []TokenID
}

func NewTokenAssociateTransaction() *TokenAssociateTransaction {
	pb := &proto.TokenAssociateTransactionBody{}

	transaction := TokenAssociateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func tokenAssociateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenAssociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenAssociate().Tokens {
		tokens = append(tokens, tokenIDFromProtobuf(token, nil))
	}

	return TokenAssociateTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenAssociate(),
		accountID:   accountIDFromProtobuf(pb.GetTokenAssociate().GetAccount(), nil),
		tokens:      tokens,
	}
}

// The account to be associated with the provided tokens
func (transaction *TokenAssociateTransaction) SetAccountID(id AccountID) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.accountID = id
	return transaction
}

func (transaction *TokenAssociateTransaction) GetAccountID() AccountID {
	return transaction.accountID
}

// The tokens to be associated with the provided account
func (transaction *TokenAssociateTransaction) SetTokenIDs(ids ...TokenID) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.tokens = make([]TokenID, len(ids))

	for i, tokenID := range ids {
		transaction.tokens[i] = tokenID
	}

	return transaction
}

func (transaction *TokenAssociateTransaction) AddTokenID(id TokenID) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	if transaction.tokens == nil {
		transaction.tokens = make([]TokenID, 0)
	}

	transaction.tokens = append(transaction.tokens, id)

	return transaction
}

func (transaction *TokenAssociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(transaction.tokens))

	for i, tokenID := range transaction.tokens {
		tokenIDs[i] = tokenID
	}

	return tokenIDs
}

func (transaction *TokenAssociateTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.accountID.Validate(client)
	if err != nil {
		return err
	}
	for _, tokenID := range transaction.tokens {
		err = tokenID.Validate(client)
	}
	if err != nil {
		return err
	}

	return nil
}

func (transaction *TokenAssociateTransaction) build() *TokenAssociateTransaction {
	if !transaction.accountID.isZero() {
		transaction.pb.Account = transaction.accountID.toProtobuf()
	}

	if len(transaction.tokens) > 0 {
		for _, tokenID := range transaction.tokens {
			if transaction.pb.Tokens == nil {
				transaction.pb.Tokens = make([]*proto.TokenID, 0)
			}
			transaction.pb.Tokens = append(transaction.pb.Tokens, tokenID.toProtobuf())
		}
	}

	return transaction
}

func (transaction *TokenAssociateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenAssociateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenAssociate{
			TokenAssociate: &proto.TokenAssociateTransactionBody{
				Account: transaction.pb.GetAccount(),
				Tokens:  transaction.pb.GetTokens(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func tokenAssociateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().AssociateTokens,
	}
}

func (transaction *TokenAssociateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenAssociateTransaction) Sign(
	privateKey PrivateKey,
) *TokenAssociateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenAssociateTransaction) SignWithOperator(
	client *Client,
) (*TokenAssociateTransaction, error) {
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
func (transaction *TokenAssociateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenAssociateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenAssociateTransaction) Execute(
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
		tokenAssociateTransaction_getMethod,
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

func (transaction *TokenAssociateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenAssociate{
		TokenAssociate: transaction.pb,
	}

	return true
}

func (transaction *TokenAssociateTransaction) Freeze() (*TokenAssociateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenAssociateTransaction) FreezeWith(client *Client) (*TokenAssociateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenAssociateTransaction{}, err
	}
	transaction.build()

	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *TokenAssociateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetMaxTransactionFee(fee Hbar) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenAssociateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetTransactionMemo(memo string) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenAssociateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenAssociateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetTransactionID(transactionID TransactionID) *TokenAssociateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenAssociateTransaction) SetMaxRetry(count int) *TokenAssociateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenAssociateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenAssociateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
