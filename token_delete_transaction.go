package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Deletes an already created Token.
// If no value is given for a field, that field is left unchanged. For an immutable tokens
// (that is, a token created without an adminKey), only the expiry may be deleted. Setting any
// other field in that case will cause the transaction status to resolve to TOKEN_IS_IMMUTABlE.
type TokenDeleteTransaction struct {
	Transaction
	tokenID *TokenID
}

func NewTokenDeleteTransaction() *TokenDeleteTransaction {
	transaction := TokenDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func _TokenDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenDeleteTransaction {
	return TokenDeleteTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenDeletion().GetToken()),
	}
}

// The Token to be deleted
func (transaction *TokenDeleteTransaction) SetTokenID(tokenID TokenID) *TokenDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

func (transaction *TokenDeleteTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

func (transaction *TokenDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
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

func (transaction *TokenDeleteTransaction) _Build() *proto.TransactionBody {
	body := &proto.TokenDeleteTransactionBody{}
	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &proto.TransactionBody_TokenDeletion{
			TokenDeletion: body,
		},
	}
}

func (transaction *TokenDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenDeleteTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.TokenDeleteTransactionBody{}
	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
	}
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenDeletion{
			TokenDeletion: body,
		},
	}, nil
}

func _TokenDeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DeleteToken,
	}
}

func (transaction *TokenDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenDeleteTransaction) Sign(
	privateKey PrivateKey,
) *TokenDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenDeleteTransaction) SignWithOperator(
	client *Client,
) (*TokenDeleteTransaction, error) {
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
func (transaction *TokenDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenDeleteTransaction) Execute(
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
		_TransactionMakeRequest(_Request{
			transaction: &transaction.Transaction,
		}),
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TokenDeleteTransactionGetMethod,
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

func (transaction *TokenDeleteTransaction) Freeze() (*TokenDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenDeleteTransaction) FreezeWith(client *Client) (*TokenDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenDeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenDeleteTransaction.
func (transaction *TokenDeleteTransaction) SetMaxTransactionFee(fee Hbar) *TokenDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenDeleteTransaction.
func (transaction *TokenDeleteTransaction) SetTransactionMemo(memo string) *TokenDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenDeleteTransaction.
func (transaction *TokenDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *TokenDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenDeleteTransaction.
func (transaction *TokenDeleteTransaction) SetTransactionID(transactionID TransactionID) *TokenDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenDeleteTransaction.
func (transaction *TokenDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenDeleteTransaction) SetMaxRetry(count int) *TokenDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenDeleteTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
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
			publicKey._ToSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

func (transaction *TokenDeleteTransaction) SetMaxBackoff(max time.Duration) *TokenDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenDeleteTransaction) SetMinBackoff(min time.Duration) *TokenDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
