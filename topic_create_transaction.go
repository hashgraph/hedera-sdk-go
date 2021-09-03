package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// A TopicCreateTransaction is for creating a new Topic on HCS.
type TopicCreateTransaction struct {
	Transaction
	autoRenewAccountID *AccountID
	adminKey           Key
	submitKey          Key
	memo               string
	autoRenewPeriod    *time.Duration
}

// NewTopicCreateTransaction creates a TopicCreateTransaction transaction which can be
// used to construct and execute a  Create Topic Transaction.
func NewTopicCreateTransaction() *TopicCreateTransaction {
	transaction := TopicCreateTransaction{
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(2))

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	// transaction.SetReceiveRecordThreshold(MaxHbar)
	// transaction.SetSendRecordThreshold(MaxHbar)

	return &transaction
}

func topicCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TopicCreateTransaction {
	adminKey, _ := keyFromProtobuf(pb.GetConsensusCreateTopic().GetAdminKey())
	submitKey, _ := keyFromProtobuf(pb.GetConsensusCreateTopic().GetSubmitKey())

	autoRenew := durationFromProtobuf(pb.GetConsensusCreateTopic().GetAutoRenewPeriod())
	return TopicCreateTransaction{
		Transaction:        transaction,
		autoRenewAccountID: accountIDFromProtobuf(pb.GetConsensusCreateTopic().GetAutoRenewAccount()),
		adminKey:           adminKey,
		submitKey:          submitKey,
		memo:               pb.GetContractCreateInstance().GetMemo(),
		autoRenewPeriod:    &autoRenew,
	}
}

// SetAdminKey sets the key required to update or delete the topic. If unspecified, anyone can increase the topic's
// expirationTime.
func (transaction *TopicCreateTransaction) SetAdminKey(publicKey Key) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

func (transaction *TopicCreateTransaction) GetAdminKey() (Key, error) {
	return transaction.adminKey, nil
}

// SetSubmitKey sets the key required for submitting messages to the topic. If unspecified, all submissions are allowed.
func (transaction *TopicCreateTransaction) SetSubmitKey(publicKey Key) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.submitKey = publicKey
	return transaction
}

func (transaction *TopicCreateTransaction) GetSubmitKey() (Key, error) {
	return transaction.submitKey, nil
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction *TopicCreateTransaction) SetTopicMemo(memo string) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.memo = memo
	return transaction
}

func (transaction *TopicCreateTransaction) GetTopicMemo() string {
	return transaction.memo
}

// SetAutoRenewPeriod sets the initial lifetime of the topic and the amount of time to extend the topic's lifetime
// automatically at expirationTime if the autoRenewAccount is configured and has sufficient funds.
//
// Required. Limited to a maximum of 90 days (server-sIDe configuration which may change).
func (transaction *TopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.autoRenewPeriod = &period
	return transaction
}

func (transaction *TopicCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetAutoRenewAccountID sets an optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount).
//
// If specified, there must be an adminKey and the autoRenewAccount must sign this transaction.
func (transaction *TopicCreateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.autoRenewAccountID = &autoRenewAccountID
	return transaction
}

func (transaction *TopicCreateTransaction) GetAutoRenewAccountID() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

func (transaction *TopicCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.autoRenewAccountID != nil {
		if err := transaction.autoRenewAccountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TopicCreateTransaction) build() *proto.TransactionBody {
	body := &proto.ConsensusCreateTopicTransactionBody{
		Memo: transaction.memo,
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if !transaction.autoRenewAccountID.isZero() {
		body.AutoRenewAccount = transaction.autoRenewAccountID.toProtobuf()
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey.toProtoKey()
	}

	if transaction.submitKey != nil {
		body.SubmitKey = transaction.submitKey.toProtoKey()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: body,
		},
	}
}

func (transaction *TopicCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TopicCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.ConsensusCreateTopicTransactionBody{
		Memo: transaction.memo,
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if !transaction.autoRenewAccountID.isZero() {
		body.AutoRenewAccount = transaction.autoRenewAccountID.toProtobuf()
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey.toProtoKey()
	}

	if transaction.submitKey != nil {
		body.SubmitKey = transaction.submitKey.toProtoKey()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: body,
		},
	}, nil
}

func _TopicCreateTransactionGetMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getTopic().CreateTopic,
	}
}

func (transaction *TopicCreateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TopicCreateTransaction) Sign(
	privateKey PrivateKey,
) *TopicCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TopicCreateTransaction) SignWithOperator(
	client *Client,
) (*TopicCreateTransaction, error) {
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
func (transaction *TopicCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicCreateTransaction {
	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TopicCreateTransaction) Execute(
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
		_TopicCreateTransactionGetMethod,
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

func (transaction *TopicCreateTransaction) Freeze() (*TopicCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicCreateTransaction) FreezeWith(client *Client) (*TopicCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TopicCreateTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TopicCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetMaxTransactionFee(fee Hbar) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TopicCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetTransactionMemo(memo string) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TopicCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetTransactionValidDuration(duration time.Duration) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TopicCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetTransactionID(transactionID TransactionID) *TopicCreateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this TopicCreateTransaction.
func (transaction *TopicCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TopicCreateTransaction) SetMaxRetry(count int) *TopicCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TopicCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicCreateTransaction {
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

func (transaction *TopicCreateTransaction) SetMaxBackoff(max time.Duration) *TopicCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TopicCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TopicCreateTransaction) SetMinBackoff(min time.Duration) *TopicCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TopicCreateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
