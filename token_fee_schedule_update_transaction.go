package hedera

import (
	"time"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenFeeScheduleUpdateTransaction struct {
	Transaction
	tokenID    *TokenID
	customFees []Fee
}

func NewTokenFeeScheduleUpdateTransaction() *TokenFeeScheduleUpdateTransaction {
	transaction := TokenFeeScheduleUpdateTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func TokenFeeScheduleUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenFeeScheduleUpdateTransaction {
	customFees := make([]Fee, 0)

	for _, fee := range pb.GetTokenFeeScheduleUpdate().GetCustomFees() {
		customFees = append(customFees, customFeeFromProtobuf(fee))
	}

	return TokenFeeScheduleUpdateTransaction{
		Transaction: transaction,
		tokenID:     tokenIDFromProtobuf(pb.GetTokenFeeScheduleUpdate().TokenId),
		customFees:  customFees,
	}
}

// The account to be associated with the provided tokens
func (transaction *TokenFeeScheduleUpdateTransaction) SetTokenID(tokenID TokenID) *TokenFeeScheduleUpdateTransaction {
	transaction.requireNotFrozen()
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
	transaction.requireNotFrozen()
	transaction.customFees = fees
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetCustomFees() []Fee {
	return transaction.customFees
}

func (transaction *TokenFeeScheduleUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.Validate(client); err != nil {
			return err
		}
	}

	for _, customFee := range transaction.customFees {
		if err := customFee.validateNetworkOnIDs(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenFeeScheduleUpdateTransaction) build() *proto.TransactionBody {
	body := &proto.TokenFeeScheduleUpdateTransactionBody{}
	if !transaction.tokenID.isZero() {
		body.TokenId = transaction.tokenID.toProtobuf()
	}

	if len(transaction.customFees) > 0 {
		for _, customFee := range transaction.customFees {
			if body.CustomFees == nil {
				body.CustomFees = make([]*proto.CustomFee, 0)
			}
			body.CustomFees = append(body.CustomFees, customFee.toProtobuf())
		}
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_TokenFeeScheduleUpdate{
			TokenFeeScheduleUpdate: body,
		},
	}
}

func (transaction *TokenFeeScheduleUpdateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleSignTransaction")
}

func _TokenFeeScheduleUpdateTransactionGetMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().UpdateTokenFeeSchedule,
	}
}

func (transaction *TokenFeeScheduleUpdateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
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
func (transaction *TokenFeeScheduleUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenFeeScheduleUpdateTransaction {
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
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenFeeScheduleUpdateTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TokenFeeScheduleUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionMemo(memo string) *TokenFeeScheduleUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenFeeScheduleUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionID(transactionID TransactionID) *TokenFeeScheduleUpdateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenFeeScheduleUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) SetMaxRetry(count int) *TokenFeeScheduleUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenFeeScheduleUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenFeeScheduleUpdateTransaction {
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
