package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// A ConsensusTopicDeleteTransaction is for deleting a topic on HCS.
type TopicDeleteTransaction struct {
	Transaction
	topicID *TopicID
}

// NewConsensusTopicDeleteTransaction creates a ConsensusTopicDeleteTransaction transaction which can be used to construct
// and execute a Consensus Delete Topic Transaction.
func NewTopicDeleteTransaction() *TopicDeleteTransaction {
	transaction := TopicDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _TopicDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TopicDeleteTransaction {
	return &TopicDeleteTransaction{
		Transaction: transaction,
		topicID:     _TopicIDFromProtobuf(pb.GetConsensusDeleteTopic().GetTopicID()),
	}
}

func (transaction *TopicDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *TopicDeleteTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetTopicID sets the topic IDentifier.
func (transaction *TopicDeleteTransaction) SetTopicID(topicID TopicID) *TopicDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.topicID = &topicID
	return transaction
}

func (transaction *TopicDeleteTransaction) GetTopicID() TopicID {
	if transaction.topicID == nil {
		return TopicID{}
	}

	return *transaction.topicID
}

func (transaction *TopicDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.topicID != nil {
		if err := transaction.topicID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TopicDeleteTransaction) _Build() *services.TransactionBody {
	body := &services.ConsensusDeleteTopicTransactionBody{}
	if transaction.topicID != nil {
		body.TopicID = transaction.topicID._ToProtobuf()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: body,
		},
	}
}

func (transaction *TopicDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TopicDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ConsensusDeleteTopicTransactionBody{}
	if transaction.topicID != nil {
		body.TopicID = transaction.topicID._ToProtobuf()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: body,
		},
	}, nil
}

func _TopicDeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().DeleteTopic,
	}
}

func (transaction *TopicDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TopicDeleteTransaction) Sign(
	privateKey PrivateKey,
) *TopicDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TopicDeleteTransaction) SignWithOperator(
	client *Client,
) (*TopicDeleteTransaction, error) {
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
func (transaction *TopicDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TopicDeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if transaction.lockError != nil {
		return TransactionResponse{}, transaction.lockError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	var transactionID TransactionID
	if transaction.transactionIDs._Length() > 0 {
		switch t := transaction.transactionIDs._Get(transaction.nextTransactionIndex).(type) { //nolint
		case TransactionID:
			transactionID = t
		}
	}

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
		_TopicDeleteTransactionGetMethod,
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

func (transaction *TopicDeleteTransaction) Freeze() (*TopicDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicDeleteTransaction) FreezeWith(client *Client) (*TopicDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TopicDeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// SetMaxTransactionFee sets the max transaction fee for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetMaxTransactionFee(fee Hbar) *TopicDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TopicDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TopicDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TopicDeleteTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// SetTransactionMemo sets the memo for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionMemo(memo string) *TopicDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// SetTransactionValidDuration sets the valid duration for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *TopicDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// SetTransactionID sets the TransactionID for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionID(transactionID TransactionID) *TopicDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TopicDeleteTransaction) SetMaxRetry(count int) *TopicDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TopicDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicDeleteTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockedSlice()
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
		_, err := transaction.signedTransactions._Set(index, temp)
		if err != nil {
			transaction.lockError = err
		}
	}

	return transaction
}

func (transaction *TopicDeleteTransaction) SetMaxBackoff(max time.Duration) *TopicDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TopicDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TopicDeleteTransaction) SetMinBackoff(min time.Duration) *TopicDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TopicDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TopicDeleteTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TopicDeleteTransaction:%d", timestamp.UnixNano())
}
