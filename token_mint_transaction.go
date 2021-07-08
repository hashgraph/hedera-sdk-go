package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Mints tokens from the Token's treasury Account. If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below
// zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to mint 100 tokens, one must provide amount of 10000. In order
// to mint 100.55 tokens, one must provide amount of 10055.
type TokenMintTransaction struct {
	Transaction
	pb      *proto.TokenMintTransactionBody
	tokenID TokenID
}

func NewTokenMintTransaction() *TokenMintTransaction {
	pb := &proto.TokenMintTransactionBody{}

	transaction := TokenMintTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func tokenMintTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenMintTransaction {
	return TokenMintTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenMint(),
		tokenID:     tokenIDFromProtobuf(pb.GetTokenMint().GetToken(), nil),
	}
}

// The token for which to mint tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *TokenMintTransaction) SetTokenID(id TokenID) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.tokenID = id
	return transaction
}

func (transaction *TokenMintTransaction) GetTokenID() TokenID {
	return transaction.tokenID
}

// The amount to mint from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (transaction *TokenMintTransaction) SetAmount(amount uint64) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Amount = amount
	return transaction
}

func (transaction *TokenMintTransaction) GetAmount() uint64 {
	return transaction.pb.GetAmount()
}

func (transaction *TokenMintTransaction) SetMetadatas(meta [][]byte) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Metadata = meta
	return transaction
}

func (transaction *TokenMintTransaction) SetMetadata(meta []byte) *TokenMintTransaction {
	transaction.requireNotFrozen()
	if transaction.pb.Metadata == nil {
		transaction.pb.Metadata = make([][]byte, 0)
	}
	transaction.pb.Metadata = append(transaction.pb.Metadata, [][]byte{meta}...)
	return transaction
}

func (transaction *TokenMintTransaction) GetMetadatas() [][]byte {
	return transaction.pb.GetMetadata()
}

func (transaction *TokenMintTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.tokenID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *TokenMintTransaction) build() *TokenMintTransaction {
	if !transaction.tokenID.isZero() {
		transaction.pb.Token = transaction.tokenID.toProtobuf()
	}

	return transaction
}

func (transaction *TokenMintTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenMintTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenMint{
			TokenMint: &proto.TokenMintTransactionBody{
				Token:  transaction.pb.GetToken(),
				Amount: transaction.pb.GetAmount(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func tokenMintTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().MintToken,
	}
}

func (transaction *TokenMintTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenMintTransaction) Sign(
	privateKey PrivateKey,
) *TokenMintTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenMintTransaction) SignWithOperator(
	client *Client,
) (*TokenMintTransaction, error) {
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
func (transaction *TokenMintTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenMintTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenMintTransaction) Execute(
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
		tokenMintTransaction_getMethod,
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

func (transaction *TokenMintTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenMint{
		TokenMint: transaction.pb,
	}

	return true
}

func (transaction *TokenMintTransaction) Freeze() (*TokenMintTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenMintTransaction) FreezeWith(client *Client) (*TokenMintTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenMintTransaction{}, err
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

func (transaction *TokenMintTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetMaxTransactionFee(fee Hbar) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenMintTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetTransactionMemo(memo string) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenMintTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetTransactionValidDuration(duration time.Duration) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenMintTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetTransactionID(transactionID TransactionID) *TokenMintTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenMintTransaction) SetMaxRetry(count int) *TokenMintTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenMintTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenMintTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
