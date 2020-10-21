package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// Burns tokens from the Token's treasury Account. If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below
// zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to burn 100 tokens, one must provide amount of 10000. In order
// to burn 100.55 tokens, one must provide amount of 10055.
type TokenBurnTransaction struct {
	Transaction
	pb *proto.TokenBurnTransactionBody
}

func NewTokenBurnTransaction() *TokenBurnTransaction {
	pb := &proto.TokenBurnTransactionBody{}

	transaction := TokenBurnTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

// The token for which to burn tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *TokenBurnTransaction) SetTokenID(tokenID TokenID) *TokenBurnTransaction {
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

// The amount to burn from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (transaction *TokenBurnTransaction) SetAmount(amount uint64) *TokenBurnTransaction {
	transaction.pb.Amount = amount
	return transaction
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func tokenBurnTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().BurnToken,
	}
}

func (transaction *TokenBurnTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenBurnTransaction) Sign(
	privateKey PrivateKey,
) *TokenBurnTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenBurnTransaction) SignWithOperator(
	client *Client,
) (*TokenBurnTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *TokenBurnTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenBurnTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.transactions); index++ {
		signature := signer(transaction.transactions[index].GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenBurnTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.id

	if !client.GetOperatorID().isZero() && client.GetOperatorID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorKey(),
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
		transaction_getNodeId,
		tokenBurnTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.id,
		NodeID:        resp.transaction.NodeID,
	}, nil
}

func (transaction *TokenBurnTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenBurn{
		TokenBurn: transaction.pb,
	}

	return true
}

func (transaction *TokenBurnTransaction) Freeze() (*TokenBurnTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenBurnTransaction) FreezeWith(client *Client) (*TokenBurnTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *TokenBurnTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetMaxTransactionFee(fee Hbar) *TokenBurnTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenBurnTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionMemo(memo string) *TokenBurnTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenBurnTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionValidDuration(duration time.Duration) *TokenBurnTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenBurnTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionID(transactionID TransactionID) *TokenBurnTransaction {
	transaction.id = transactionID
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *TokenBurnTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeTokenID sets the node TokenID for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenBurnTransaction {
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}
