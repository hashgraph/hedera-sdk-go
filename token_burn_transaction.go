package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
	tokenID *TokenID
	amount  uint64
	serial  []int64
}

func NewTokenBurnTransaction() *TokenBurnTransaction {
	transaction := TokenBurnTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _TokenBurnTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenBurnTransaction {
	return &TokenBurnTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenBurn().Token),
		amount:      pb.GetTokenBurn().GetAmount(),
		serial:      pb.GetTokenBurn().GetSerialNumbers(),
	}
}

func (transaction *TokenBurnTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenBurnTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// The token for which to burn tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *TokenBurnTransaction) SetTokenID(tokenID TokenID) *TokenBurnTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

func (transaction *TokenBurnTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

// The amount to burn from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (transaction *TokenBurnTransaction) SetAmount(amount uint64) *TokenBurnTransaction {
	transaction._RequireNotFrozen()
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
	transaction._RequireNotFrozen()
	if transaction.serial == nil {
		transaction.serial = make([]int64, 0)
	}
	transaction.serial = append(transaction.serial, serial)
	return transaction
}

func (transaction *TokenBurnTransaction) SetSerialNumbers(serial []int64) *TokenBurnTransaction {
	transaction._RequireNotFrozen()
	transaction.serial = serial
	return transaction
}

func (transaction *TokenBurnTransaction) GetSerialNumbers() []int64 {
	return transaction.serial
}

func (transaction *TokenBurnTransaction) _ValidateNetworkOnIDs(client *Client) error {
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

func (transaction *TokenBurnTransaction) _Build() *services.TransactionBody {
	body := &services.TokenBurnTransactionBody{
		Amount: transaction.amount,
	}

	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	if transaction.serial != nil {
		body.SerialNumbers = transaction.serial
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenBurn{
			TokenBurn: body,
		},
	}
}

func (transaction *TokenBurnTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenBurnTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.TokenBurnTransactionBody{
		Amount: transaction.amount,
	}

	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	if transaction.serial != nil {
		body.SerialNumbers = transaction.serial
	}
	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenBurn{
			TokenBurn: body,
		},
	}, nil
}

func _TokenBurnTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().BurnToken,
	}
}

func (transaction *TokenBurnTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
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
func (transaction *TokenBurnTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenBurnTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
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
		_TokenBurnTransactionGetMethod,
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

func (transaction *TokenBurnTransaction) Freeze() (*TokenBurnTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenBurnTransaction) FreezeWith(client *Client) (*TokenBurnTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenBurnTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenBurnTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetMaxTransactionFee(fee Hbar) *TokenBurnTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TokenBurnTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenBurnTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TokenBurnTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *TokenBurnTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionMemo(memo string) *TokenBurnTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenBurnTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionValidDuration(duration time.Duration) *TokenBurnTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenBurnTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetTransactionID(transactionID TransactionID) *TokenBurnTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenBurnTransaction.
func (transaction *TokenBurnTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenBurnTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenBurnTransaction) SetMaxRetry(count int) *TokenBurnTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenBurnTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenBurnTransaction {
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

func (transaction *TokenBurnTransaction) SetMaxBackoff(max time.Duration) *TokenBurnTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenBurnTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenBurnTransaction) SetMinBackoff(min time.Duration) *TokenBurnTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenBurnTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TokenBurnTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TokenBurnTransaction:%d", timestamp.UnixNano())
}
