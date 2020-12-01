package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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
	pb *proto.TokenGrantKycTransactionBody
}

func NewTokenGrantKycTransaction() *TokenGrantKycTransaction {
	pb := &proto.TokenGrantKycTransactionBody{}

	transaction := TokenGrantKycTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func tokenGrantKycTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenGrantKycTransaction {
	return TokenGrantKycTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenGrantKyc(),
	}
}

// The token for which this account will be granted KYC. If token does not exist, transaction results in INVALID_TOKEN_ID
func (transaction *TokenGrantKycTransaction) SetTokenID(tokenID TokenID) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTokenID() TokenID {
	return tokenIDFromProtobuf(transaction.pb.GetToken())
}

// The account to be KYCed
func (transaction *TokenGrantKycTransaction) SetAccountID(accountID AccountID) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Account = accountID.toProtobuf()
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAccount())
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func tokenGrantKycTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().GrantKycToTokenAccount,
	}
}

func (transaction *TokenGrantKycTransaction) IsFrozen() bool {
	return transaction.isFrozen()
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
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
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
func (transaction *TokenGrantKycTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
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
		tokenGrantKycTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.transactionIDs[0],
		NodeID:        resp.transaction.NodeID,
	}, nil
}

func (transaction *TokenGrantKycTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenGrantKyc{
		TokenGrantKyc: transaction.pb,
	}

	return true
}

func (transaction *TokenGrantKycTransaction) Freeze() (*TokenGrantKycTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenGrantKycTransaction) FreezeWith(client *Client) (*TokenGrantKycTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *TokenGrantKycTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetMaxTransactionFee(fee Hbar) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionMemo(memo string) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionValidDuration(duration time.Duration) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionID(transactionID TransactionID) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeTokenID sets the node TokenID for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenGrantKycTransaction) SetMaxRetry(count int) *TokenGrantKycTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}
