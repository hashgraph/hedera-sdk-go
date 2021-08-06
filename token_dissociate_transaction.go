package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenDissociateTransaction struct {
	Transaction
	accountID AccountID
	tokens    []TokenID
}

func NewTokenDissociateTransaction() *TokenDissociateTransaction {
	transaction := TokenDissociateTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func tokenDissociateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenDissociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenDissociate().Tokens {
		tokens = append(tokens, tokenIDFromProtobuf(token))
	}

	return TokenDissociateTransaction{
		Transaction: transaction,
		accountID:   accountIDFromProtobuf(pb.GetTokenDissociate().GetAccount()),
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
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
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

func (transaction *TokenDissociateTransaction) build() *proto.TransactionBody {
	body := &proto.TokenDissociateTransactionBody{}
	if !transaction.accountID.isZero() {
		body.Account = transaction.accountID.toProtobuf()
	}

	if len(transaction.tokens) > 0 {
		for _, tokenID := range transaction.tokens {
			if body.Tokens == nil {
				body.Tokens = make([]*proto.TokenID, 0)
			}
			body.Tokens = append(body.Tokens, tokenID.toProtobuf())
		}
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_TokenDissociate{
			TokenDissociate: body,
		},
	}
}

func (transaction *TokenDissociateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenDissociateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.TokenDissociateTransactionBody{}
	if !transaction.accountID.isZero() {
		body.Account = transaction.accountID.toProtobuf()
	}

	if len(transaction.tokens) > 0 {
		for _, tokenID := range transaction.tokens {
			if body.Tokens == nil {
				body.Tokens = make([]*proto.TokenID, 0)
			}
			body.Tokens = append(body.Tokens, tokenID.toProtobuf())
		}
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenDissociate{
			TokenDissociate: body,
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
		transaction_makeRequest(request{
			transaction: &transaction.Transaction,
		}),
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
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
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
	transaction.requireOneNodeAccountID()

	if !transaction.isFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
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
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	//transaction.signedTransactions[0].SigMap.SigPair = append(transaction.signedTransactions[0].SigMap.SigPair, publicKey.toSignaturePairProtobuf(signature))
	return transaction
}
