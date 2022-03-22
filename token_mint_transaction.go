package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func _TokenMintTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenMintTransaction {
	return &TokenMintTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenMint().GetToken()),
		amount:      pb.GetTokenMint().GetAmount(),
		meta:        pb.GetTokenMint().GetMetadata(),
	}
}

func (transaction *TokenMintTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenMintTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// The token for which to mint tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *TokenMintTransaction) SetTokenID(tokenID TokenID) *TokenMintTransaction {
	transaction._RequireNotFrozen()
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
	transaction._RequireNotFrozen()
	transaction.amount = amount
	return transaction
}

func (transaction *TokenMintTransaction) GetAmount() uint64 {
	return transaction.amount
}

func (transaction *TokenMintTransaction) SetMetadatas(meta [][]byte) *TokenMintTransaction {
	transaction._RequireNotFrozen()
	transaction.meta = meta
	return transaction
}

func (transaction *TokenMintTransaction) SetMetadata(meta []byte) *TokenMintTransaction {
	transaction._RequireNotFrozen()
	if transaction.meta == nil {
		transaction.meta = make([][]byte, 0)
	}
	transaction.meta = append(transaction.meta, [][]byte{meta}...)
	return transaction
}

func (transaction *TokenMintTransaction) GetMetadatas() [][]byte {
	return transaction.meta
}

func (transaction *TokenMintTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenMintTransaction) _Build() *services.TransactionBody {
	body := &services.TokenMintTransactionBody{
		Amount: transaction.amount,
	}

	if transaction.meta != nil {
		body.Metadata = transaction.meta
	}

	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenMint{
			TokenMint: body,
		},
	}
}

func (transaction *TokenMintTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenMintTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.TokenMintTransactionBody{
		Amount: transaction.amount,
	}

	if transaction.meta != nil {
		body.Metadata = transaction.meta
	}

	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}
	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenMint{
			TokenMint: body,
		},
	}, nil
}

func _TokenMintTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().MintToken,
	}
}

func (transaction *TokenMintTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
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
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

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
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
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

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TokenMintTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
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
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenMintTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenMintTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetMaxTransactionFee(fee Hbar) *TokenMintTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TokenMintTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenMintTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TokenMintTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *TokenMintTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetTransactionMemo(memo string) *TokenMintTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenMintTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetTransactionValidDuration(duration time.Duration) *TokenMintTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenMintTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetTransactionID(transactionID TransactionID) *TokenMintTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenMintTransaction.
func (transaction *TokenMintTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenMintTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenMintTransaction) SetMaxRetry(count int) *TokenMintTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenMintTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenMintTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := transaction.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
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

func (transaction *TokenMintTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TokenMintTransaction:%d", timestamp.UnixNano())
}
