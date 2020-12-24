package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenDissociateTransaction struct {
	Transaction
	pb *proto.TokenDissociateTransactionBody
}

func NewTokenDissociateTransaction() *TokenDissociateTransaction {
	pb := &proto.TokenDissociateTransactionBody{}

	transaction := TokenDissociateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func tokenDissociateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenDissociateTransaction {
	return TokenDissociateTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenDissociate(),
	}
}

// The account to be dissociated with the provided tokens
func (transaction *TokenDissociateTransaction) SetAccountID(accountID AccountID) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Account = accountID.toProtobuf()
	return transaction
}

func (transaction *TokenDissociateTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.Account)
}

// The tokens to be dissociated with the provided account
func (transaction *TokenDissociateTransaction) SetTokenIDs(tokenIDs ...TokenID) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Tokens = make([]*proto.TokenID, len(tokenIDs))

	for i, tokenID := range tokenIDs {
		transaction.pb.Tokens[i] = tokenID.toProtobuf()
	}

	return transaction
}

func (transaction *TokenDissociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(transaction.pb.Tokens))

	for i, tokenID := range transaction.pb.Tokens {
		tokenIDs[i] = tokenIDFromProtobuf(tokenID)
	}

	return tokenIDs
}

func (transaction *TokenDissociateTransaction) AddTokenID(tokenID TokenID) *TokenDissociateTransaction {
	transaction.requireNotFrozen()
	if transaction.pb.Tokens == nil {
		transaction.pb.Tokens = make([]*proto.TokenID, 0)
	}

	transaction.pb.Tokens = append(transaction.pb.Tokens, tokenID.toProtobuf())
	return transaction
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
func (transaction *TokenDissociateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
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
		tokenDissociateTransaction_getMethod,
		transaction_mapResponseStatus,
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
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenDissociate{
		TokenDissociate: transaction.pb,
	}

	return true
}

func (transaction *TokenDissociateTransaction) Freeze() (*TokenDissociateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenDissociateTransaction) FreezeWith(client *Client) (*TokenDissociateTransaction, error) {
	transaction.initFee(client)
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

func (transaction *TokenDissociateTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
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
