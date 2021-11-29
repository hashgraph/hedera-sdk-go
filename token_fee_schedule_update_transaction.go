package hedera

import (
	"time"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TokenFeeScheduleUpdateTransaction struct {
	Transaction
	tokenID    *TokenID
	customFees []Fee
}

func NewTokenFeeScheduleUpdateTransaction() *TokenFeeScheduleUpdateTransaction {
	transaction := TokenFeeScheduleUpdateTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _TokenFeeScheduleUpdateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) TokenFeeScheduleUpdateTransaction {
	customFees := make([]Fee, 0)

	for _, fee := range pb.GetTokenFeeScheduleUpdate().GetCustomFees() {
		customFees = append(customFees, _CustomFeeFromProtobuf(fee))
	}

	return TokenFeeScheduleUpdateTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenFeeScheduleUpdate().TokenId),
		customFees:  customFees,
	}
}

// The account to be associated with the provided tokens
func (transaction *TokenFeeScheduleUpdateTransaction) SetTokenID(tokenID TokenID) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

func (transaction *TokenFeeScheduleUpdateTransaction) SetCustomFees(fees []Fee) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.customFees = fees
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetCustomFees() []Fee {
	return transaction.customFees
}

func (transaction *TokenFeeScheduleUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, customFee := range transaction.customFees {
		if err := customFee._ValidateNetworkOnIDs(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenFeeScheduleUpdateTransaction) _Build() *services.TransactionBody {
	body := &services.TokenFeeScheduleUpdateTransactionBody{}
	if transaction.tokenID != nil {
		body.TokenId = transaction.tokenID._ToProtobuf()
	}

	if len(transaction.customFees) > 0 {
		for _, customFee := range transaction.customFees {
			if body.CustomFees == nil {
				body.CustomFees = make([]*services.CustomFee, 0)
			}
			body.CustomFees = append(body.CustomFees, customFee._ToProtobuf())
		}
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenFeeScheduleUpdate{
			TokenFeeScheduleUpdate: body,
		},
	}
}

func (transaction *TokenFeeScheduleUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleSignTransaction")
}

func _TokenFeeScheduleUpdateTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UpdateTokenFeeSchedule,
	}
}

func (transaction *TokenFeeScheduleUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenFeeScheduleUpdateTransaction) Sign(
	privateKey PrivateKey,
) *TokenFeeScheduleUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenFeeScheduleUpdateTransaction) SignWithOperator(
	client *Client,
) (*TokenFeeScheduleUpdateTransaction, error) {
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
func (transaction *TokenFeeScheduleUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenFeeScheduleUpdateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	} else {
		transaction.transactions = make([]*services.Transaction, 0)
	}

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.signedTransactions); index++ {
		signature := signer(transaction.signedTransactions[index].GetBodyBytes())

		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenFeeScheduleUpdateTransaction) Execute(
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
		_TokenFeeScheduleUpdateTransactionGetMethod,
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

func (transaction *TokenFeeScheduleUpdateTransaction) Freeze() (*TokenFeeScheduleUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenFeeScheduleUpdateTransaction) FreezeWith(client *Client) (*TokenFeeScheduleUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenFeeScheduleUpdateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionMemo(memo string) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionID(transactionID TransactionID) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) SetMaxRetry(count int) *TokenFeeScheduleUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if len(transaction.signedTransactions) == 0 {
		return transaction
	}

	transaction.transactions = make([]*services.Transaction, 0)
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)

	for index := 0; index < len(transaction.signedTransactions); index++ {
		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) SetMaxBackoff(max time.Duration) *TokenFeeScheduleUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenFeeScheduleUpdateTransaction) SetMinBackoff(min time.Duration) *TokenFeeScheduleUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
