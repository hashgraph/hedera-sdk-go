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
	tokenID *TokenID
	amount  uint64
	meta    [][]byte
}

func NewTokenMintTransaction() *TokenMintTransaction {
	transaction := TokenMintTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func tokenMintTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenMintTransaction {
	return TokenMintTransaction{
		Transaction: transaction,
		tokenID:     tokenIDFromProtobuf(pb.GetTokenMint().GetToken()),
		amount:      pb.GetTokenMint().GetAmount(),
		meta:        pb.GetTokenMint().GetMetadata(),
	}
}

// The token for which to mint tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *TokenMintTransaction) SetTokenID(tokenID TokenID) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

func (transaction *TokenMintTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

// The amount to mint from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (transaction *TokenMintTransaction) SetAmount(amount uint64) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.amount = amount
	return transaction
}

func (transaction *TokenMintTransaction) GetAmount() uint64 {
	return transaction.amount
}

func (transaction *TokenMintTransaction) SetMetadatas(meta [][]byte) *TokenMintTransaction {
	transaction.requireNotFrozen()
	transaction.meta = meta
	return transaction
}

func (transaction *TokenMintTransaction) SetMetadata(meta []byte) *TokenMintTransaction {
	transaction.requireNotFrozen()
	if transaction.meta == nil {
		transaction.meta = make([][]byte, 0)
	}
	transaction.meta = append(transaction.meta, [][]byte{meta}...)
	return transaction
}

func (transaction *TokenMintTransaction) GetMetadatas() [][]byte {
	return transaction.meta
}

func (transaction *TokenMintTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenMintTransaction) build() *proto.TransactionBody {
	body := &proto.TokenMintTransactionBody{
		Amount: transaction.amount,
	}

	if transaction.meta != nil {
		body.Metadata = transaction.meta
	}

	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_TokenMint{
			TokenMint: body,
		},
	}
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
	body := &proto.TokenMintTransactionBody{
		Amount: transaction.amount,
	}

	if transaction.meta != nil {
		body.Metadata = transaction.meta
	}

	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenMint{
			TokenMint: body,
		},
	}, nil
}

func _TokenMintTransactionGetMethod(request request, channel *channel) method {
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
		_TransactionShouldRetry,
		_TransactionMakeRequest(request{
			transaction: &transaction.Transaction,
		}),
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TokenMintTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()
	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
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
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
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
	transaction.requireOneNodeAccountID()

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

	return transaction
}

func (transaction *TokenMintTransaction) SetMaxBackoff(max time.Duration) *TokenMintTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenMintTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenMintTransaction) SetMinBackoff(min time.Duration) *TokenMintTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenMintTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
