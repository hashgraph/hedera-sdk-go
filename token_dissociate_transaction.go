package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"time"
)

type TokenDissociateTransaction struct {
	Transaction
	pb        *services.TokenDissociateTransactionBody
	accountID AccountID
	tokens    []TokenID
}

func NewTokenDissociateTransaction() *TokenDissociateTransaction {
	pb := &services.TokenDissociateTransactionBody{}

	transaction := TokenDissociateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func tokenDissociateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) TokenDissociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenDissociate().Tokens {
		tokens = append(tokens, tokenIDFromProtobuf(token, nil))
	}

	return TokenDissociateTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenDissociate(),
		accountID:   accountIDFromProtobuf(pb.GetTokenDissociate().GetAccount(), nil),
		tokens:      tokens,
	}
}

// The account to be dissociated with the provided tokens
func (transaction *TokenDissociateTransaction) SetAccountID(id AccountID) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.accountID = id
	return transaction
}

func (transaction *TokenDissociateTransaction) GetAccountID() AccountID {
	return transaction.accountID
}

// The tokens to be dissociated with the provided account
func (transaction *TokenDissociateTransaction) SetTokenIDs(ids ...TokenID) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.tokens = make([]TokenID, len(ids))

	for i, tokenID := range ids {
		transaction.tokens[i] = tokenID
	}

	return transaction
}

func (transaction *TokenDissociateTransaction) AddTokenID(id TokenID) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
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

func (transaction *TokenDissociateTransaction) validateNetworkOnIDs(client *Client) error {
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

func (transaction *TokenDissociateTransaction) build() *TokenDissociateTransaction {
	if !transaction.accountID.isZero() {
		transaction.pb.Account = transaction.accountID.toProtobuf()
	}

	if len(transaction.tokens) > 0 {
		for _, tokenID := range transaction.tokens {
			if transaction.pb.Tokens == nil {
				transaction.pb.Tokens = make([]*services.TokenID, 0)
			}
			transaction.pb.Tokens = append(transaction.pb.Tokens, tokenID.toProtobuf())
		}
	}

	return transaction
}

func (transaction *TokenDissociateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenDissociateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	transaction.build()
	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &services.SchedulableTransactionBody_TokenDissociate{
			TokenDissociate: &services.TokenDissociateTransactionBody{
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

func tokenDissociateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().DissociateTokens,
	}
}

func (transaction *TokenDissociateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
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
func (transaction *TokenDissociateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenDissociateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
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
		tokenDissociateTransaction_getMethod,
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

func (transaction *TokenDissociateTransaction) onFreeze(
	pbBody *services.TransactionBody,
) bool {
	pbBody.Data = &services.TransactionBody_TokenDissociate{
		TokenDissociate: transaction.pb,
	}

	return true
}

func (transaction *TokenDissociateTransaction) Freeze() (*TokenDissociateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenDissociateTransaction) FreezeWith(client *Client) (*TokenDissociateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenDissociateTransaction{}, err
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

func (transaction *TokenDissociateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetMaxTransactionFee(fee Hbar) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenDissociateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetTransactionMemo(memo string) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenDissociateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenDissociateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetTransactionID(transactionID TransactionID) *TokenDissociateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenDissociateTransaction.
func (transaction *TokenDissociateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenDissociateTransaction) SetMaxRetry(count int) *TokenDissociateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenDissociateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenDissociateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
