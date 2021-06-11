package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type NftDissociateTransaction struct {
	Transaction
	pb *proto.TokenDissociateTransactionBody
}

func NewNftDissociateTransaction() *NftDissociateTransaction {
	pb := &proto.TokenDissociateTransactionBody{}

	transaction := NftDissociateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func nftDissociateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenDissociateTransaction {
	return TokenDissociateTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenDissociate(),
	}
}

// The account to be dissociated with the provided tokens
func (transaction *NftDissociateTransaction) SetAccountID(accountID AccountID) *NftDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Account = accountID.toProtobuf()
	return transaction
}

func (transaction *NftDissociateTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.Account)
}

// The tokens to be dissociated with the provided account
func (transaction *NftDissociateTransaction) SetTokenIDs(tokenIDs ...TokenID) *NftDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Tokens = make([]*proto.TokenID, len(tokenIDs))

	for i, tokenID := range tokenIDs {
		transaction.pb.Tokens[i] = tokenID.toProtobuf()
	}

	return transaction
}

func (transaction *NftDissociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(transaction.pb.Tokens))

	for i, tokenID := range transaction.pb.Tokens {
		tokenIDs[i] = tokenIDFromProtobuf(tokenID)
	}

	return tokenIDs
}

func (transaction *NftDissociateTransaction) AddTokenID(tokenID TokenID) *NftDissociateTransaction {
	transaction.requireNotFrozen()
	if transaction.pb.Tokens == nil {
		transaction.pb.Tokens = make([]*proto.TokenID, 0)
	}

	transaction.pb.Tokens = append(transaction.pb.Tokens, tokenID.toProtobuf())
	return transaction
}

func (transaction *NftDissociateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftDissociateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenDissociate{
			TokenDissociate: &proto.TokenDissociateTransactionBody{
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

func nftDissociateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().DissociateTokens,
	}
}

func (transaction *NftDissociateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftDissociateTransaction) Sign(
	privateKey PrivateKey,
) *NftDissociateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftDissociateTransaction) SignWithOperator(
	client *Client,
) (*NftDissociateTransaction, error) {
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
func (transaction *NftDissociateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftDissociateTransaction {
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
func (transaction *NftDissociateTransaction) Execute(
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
		nftDissociateTransaction_getMethod,
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

func (transaction *NftDissociateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenDissociate{
		TokenDissociate: transaction.pb,
	}

	return true
}

func (transaction *NftDissociateTransaction) Freeze() (*NftDissociateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *NftDissociateTransaction) FreezeWith(client *Client) (*NftDissociateTransaction, error) {
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

func (transaction *NftDissociateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenDissociateTransaction.
func (transaction *NftDissociateTransaction) SetMaxTransactionFee(fee Hbar) *NftDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftDissociateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenDissociateTransaction.
func (transaction *NftDissociateTransaction) SetTransactionMemo(memo string) *NftDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftDissociateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenDissociateTransaction.
func (transaction *NftDissociateTransaction) SetTransactionValidDuration(duration time.Duration) *NftDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftDissociateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenDissociateTransaction.
func (transaction *NftDissociateTransaction) SetTransactionID(transactionID TransactionID) *NftDissociateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenDissociateTransaction.
func (transaction *NftDissociateTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftDissociateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftDissociateTransaction) SetMaxRetry(count int) *NftDissociateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftDissociateTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftDissociateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
