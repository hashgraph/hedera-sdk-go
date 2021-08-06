package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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
	tokenID TokenID
	amount  uint64
	serial  []int64
}

func NewTokenBurnTransaction() *TokenBurnTransaction {
	transaction := TokenBurnTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func tokenBurnTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenBurnTransaction {
	return TokenBurnTransaction{
		Transaction: transaction,
		tokenID:     tokenIDFromProtobuf(pb.GetTokenBurn().Token),
		amount:      pb.GetTokenBurn().GetAmount(),
		serial:      pb.GetTokenBurn().GetSerialNumbers(),
	}
}

// The token for which to burn tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *TokenBurnTransaction) SetTokenID(id TokenID) *TokenBurnTransaction {
	transaction.requireNotFrozen()
	transaction.tokenID = id
	return transaction
}

func (transaction *TokenBurnTransaction) GetTokenID() TokenID {
	return transaction.tokenID
}

// The amount to burn from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (transaction *TokenBurnTransaction) SetAmount(amount uint64) *TokenBurnTransaction {
	transaction.requireNotFrozen()
	transaction.amount = amount
	return transaction
}

// Deprecated: Use TokenBurnTransaction.GetAmount() instead.
func (transaction *TokenBurnTransaction) GetAmmount() uint64 {
	return transaction.amount
}

func (transaction *TokenBurnTransaction) GetAmount() uint64 {
	return transaction.amount
}

func (transaction *TokenBurnTransaction) SetSerialNumber(serial int64) *TokenBurnTransaction {
	transaction.requireNotFrozen()
	if transaction.serial == nil {
		transaction.serial = make([]int64, 0)
	}
	transaction.serial = append(transaction.serial, serial)
	return transaction
}

func (transaction *TokenBurnTransaction) SetSerialNumbers(serial []int64) *TokenBurnTransaction {
	transaction.requireNotFrozen()
	transaction.serial = serial
	return transaction
}

func (transaction *TokenBurnTransaction) GetSerialNumbers() []int64 {
	return transaction.serial
}

func (transaction *TokenBurnTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = transaction.tokenID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *TokenBurnTransaction) build() *proto.TransactionBody {
	body := &proto.TokenBurnTransactionBody{
		Amount: transaction.amount,
	}

	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}

	if transaction.serial != nil {
		body.SerialNumbers = transaction.serial
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_TokenBurn{
			TokenBurn: body,
		},
	}
}

func (transaction *TokenBurnTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenBurnTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.TokenBurnTransactionBody{
		Amount: transaction.amount,
	}

	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}

	if transaction.serial != nil {
		body.SerialNumbers = transaction.serial
	}
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenBurn{
			TokenBurn: body,
		},
	}, nil
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
func (transaction *TokenBurnTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenBurnTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenBurnTransaction) Execute(
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

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
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
		tokenBurnTransaction_getMethod,
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

func (transaction *TokenBurnTransaction) Freeze() (*TokenBurnTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenBurnTransaction) FreezeWith(client *Client) (*TokenBurnTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenBurnTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenBurnTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetMaxTransactionFee(fee Hbar) *TokenBurnTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenBurnTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionMemo(memo string) *TokenBurnTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenBurnTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionValidDuration(duration time.Duration) *TokenBurnTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenBurnTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionID(transactionID TransactionID) *TokenBurnTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenBurnTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenBurnTransaction) SetMaxRetry(count int) *TokenBurnTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenBurnTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenBurnTransaction {
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
