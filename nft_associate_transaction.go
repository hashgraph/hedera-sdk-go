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
type NftAssociateTransaction struct {
	Transaction
	pb *proto.TokenAssociateTransactionBody
}

func NewNftAssociateTransaction() *NftAssociateTransaction {
	pb := &proto.TokenAssociateTransactionBody{}

	transaction := NftAssociateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func nftAssociateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenAssociateTransaction {
	return TokenAssociateTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenAssociate(),
	}
}

// The account to be associated with the provided tokens
func (transaction *NftAssociateTransaction) SetAccountID(accountID AccountID) *NftAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Account = accountID.toProtobuf()
	return transaction
}

func (transaction *NftAssociateTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAccount())
}

// The tokens to be associated with the provided account
func (transaction *NftAssociateTransaction) SetTokenIDs(tokenIDs ...TokenID) *NftAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Tokens = make([]*proto.TokenID, len(tokenIDs))

	for i, tokenID := range tokenIDs {
		transaction.pb.Tokens[i] = tokenID.toProtobuf()
	}

	return transaction
}

func (transaction *NftAssociateTransaction) AddTokenID(tokenID TokenID) *NftAssociateTransaction {
	transaction.requireNotFrozen()
	if transaction.pb.Tokens == nil {
		transaction.pb.Tokens = make([]*proto.TokenID, 0)
	}

	transaction.pb.Tokens = append(transaction.pb.Tokens, tokenID.toProtobuf())

	return transaction
}

func (transaction *NftAssociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(transaction.pb.GetTokens()))

	for i, tokenID := range transaction.pb.GetTokens() {
		tokenIDs[i] = tokenIDFromProtobuf(tokenID)
	}

	return tokenIDs
}

func (transaction *NftAssociateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftAssociateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
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

func nftAssociateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().AssociateTokens,
	}
}

func (transaction *NftAssociateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftAssociateTransaction) Sign(
	privateKey PrivateKey,
) *NftAssociateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftAssociateTransaction) SignWithOperator(
	client *Client,
) (*NftAssociateTransaction, error) {
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
func (transaction *NftAssociateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftAssociateTransaction {
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
func (transaction *NftAssociateTransaction) Execute(
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
		nftAssociateTransaction_getMethod,
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

func (transaction *NftAssociateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenAssociate{
		TokenAssociate: transaction.pb,
	}

	return true
}

func (transaction *NftAssociateTransaction) Freeze() (*NftAssociateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *NftAssociateTransaction) FreezeWith(client *Client) (*NftAssociateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *NftAssociateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenAssociateTransaction.
func (transaction *NftAssociateTransaction) SetMaxTransactionFee(fee Hbar) *NftAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftAssociateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenAssociateTransaction.
func (transaction *NftAssociateTransaction) SetTransactionMemo(memo string) *NftAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftAssociateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenAssociateTransaction.
func (transaction *NftAssociateTransaction) SetTransactionValidDuration(duration time.Duration) *NftAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftAssociateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenAssociateTransaction.
func (transaction *NftAssociateTransaction) SetTransactionID(transactionID TransactionID) *NftAssociateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenAssociateTransaction.
func (transaction *NftAssociateTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftAssociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftAssociateTransaction) SetMaxRetry(count int) *NftAssociateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftAssociateTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftAssociateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
