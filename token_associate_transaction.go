package hedera

import (
	"github.com/pkg/errors"
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
	pb *proto.TokenAssociateTransactionBody
}

func NewTokenAssociateTransaction() *TokenAssociateTransaction {
	pb := &proto.TokenAssociateTransactionBody{}

	transaction := TokenAssociateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func tokenAssociateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenAssociateTransaction {
	return TokenAssociateTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenAssociate(),
	}
}

// The account to be associated with the provided tokens
func (transaction *TokenAssociateTransaction) SetAccountID(accountID AccountID) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Account = accountID.toProtobuf()
	return transaction
}

func (transaction *TokenAssociateTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAccount())
}

// The tokens to be associated with the provided account
func (transaction *TokenAssociateTransaction) SetTokenIDs(tokenIDs ...TokenID) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Tokens = make([]*proto.TokenID, len(tokenIDs))

	for i, tokenID := range tokenIDs {
		transaction.pb.Tokens[i] = tokenID.toProtobuf()
	}

	return transaction
}

func (transaction *TokenAssociateTransaction) AddTokenID(tokenID TokenID) *TokenAssociateTransaction {
	transaction.requireNotFrozen()
	if transaction.pb.Tokens == nil {
		transaction.pb.Tokens = make([]*proto.TokenID, 0)
	}

	transaction.pb.Tokens = append(transaction.pb.Tokens, tokenID.toProtobuf())

	return transaction
}

func (transaction *TokenAssociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(transaction.pb.GetTokens()))

	for i, tokenID := range transaction.pb.GetTokens() {
		tokenIDs[i] = tokenIDFromProtobuf(tokenID)
	}

	return tokenIDs
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
		return nil, errors.Wrap(errNoClientProvided, "for SignWithOperator")
	} else if client.operator == nil {
		return nil, errors.Wrap(errClientOperatorSigning, "for SignWithOperator")
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, errors.Wrap(err, "FreezeWith in SignWithOperator")
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
	} else {
		transaction.transactions = make([]*proto.Transaction, 0)
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.signedTransactions); index++ {
		signature := signer(transaction.signedTransactions[index].GetBodyBytes())

		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenAssociateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, errors.Wrap(err, "FreezeWith in Execute")
		}
	}

	transactionID := transaction.transactionIDs[0]

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
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
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, errors.Wrap(err, "execution error")
	}

	return TransactionResponse{
		TransactionID: transaction.transactionIDs[0],
		NodeID:        resp.transaction.NodeID,
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
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, errors.Wrap(err, "initTransactionID in TokenAssociateTransaction.FreezeWith")
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

func (transaction *TokenAssociateTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
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
