package hedera

import (
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// *TopicUpdateTransaction updates all fields on a Topic that are set in the transaction.
type TopicUpdateTransaction struct {
	Transaction
	topicID            *TopicID
	autoRenewAccountID *AccountID
	adminKey           Key
	submitKey          Key
	memo               string
	autoRenewPeriod    *time.Duration
	expirationTime     *time.Time
}

// NewTopicUpdateTransaction creates a *TopicUpdateTransaction transaction which can be
// used to construct and execute a  Update Topic Transaction.
func NewTopicUpdateTransaction() *TopicUpdateTransaction {
	transaction := TopicUpdateTransaction{
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func topicUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TopicUpdateTransaction {
	adminKey, _ := keyFromProtobuf(pb.GetConsensusUpdateTopic().GetAdminKey())
	submitKey, _ := keyFromProtobuf(pb.GetConsensusUpdateTopic().GetSubmitKey())

	expirationTime := timeFromProtobuf(pb.GetConsensusUpdateTopic().GetExpirationTime())
	autoRenew := durationFromProtobuf(pb.GetConsensusUpdateTopic().GetAutoRenewPeriod())
	return TopicUpdateTransaction{
		Transaction:        transaction,
		topicID:            topicIDFromProtobuf(pb.GetConsensusUpdateTopic().GetTopicID()),
		autoRenewAccountID: accountIDFromProtobuf(pb.GetConsensusUpdateTopic().GetAutoRenewAccount()),
		adminKey:           adminKey,
		submitKey:          submitKey,
		memo:               pb.GetConsensusUpdateTopic().GetMemo().Value,
		autoRenewPeriod:    &autoRenew,
		expirationTime:     &expirationTime,
	}
}

// SetTopicID sets the topic to be updated.
func (transaction *TopicUpdateTransaction) SetTopicID(id TopicID) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.topicID = id
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTopicID() TopicID {
	if transaction.topicID == nil {
		return TopicID{}
	}

	return *transaction.topicID
}

// SetAdminKey sets the key required to update/delete the topic. If unset, the key will not be changed.
//
// Setting the AdminKey to an empty KeyList will clear the adminKey.
func (transaction *TopicUpdateTransaction) SetAdminKey(publicKey Key) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAdminKey() (Key, error) {
	return transaction.adminKey, nil
}

// SetSubmitKey will set the key allowed to submit messages to the topic.  If unset, the key will not be changed.
//
// Setting the submitKey to an empty KeyList will clear the submitKey.
func (transaction *TopicUpdateTransaction) SetSubmitKey(publicKey Key) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.submitKey = publicKey
	return transaction
}

func (transaction *TopicUpdateTransaction) GetSubmitKey() (Key, error) {
	return transaction.submitKey, nil
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction *TopicUpdateTransaction) SetTopicMemo(memo string) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.memo = memo
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTopicMemo() string {
	return transaction.memo
}

// SetExpirationTime sets the effective  timestamp at (and after) which all  transactions and queries
// will fail. The expirationTime may be no longer than 90 days from the  timestamp of this transaction.
func (transaction *TopicUpdateTransaction) SetExpirationTime(expiration time.Time) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

func (transaction *TopicUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

// SetAutoRenewPeriod sets the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds. This is limited to a maximum of 90 days (server-sIDe configuration
// which may change).
func (transaction *TopicUpdateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.autoRenewPeriod = &period
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetAutoRenewAccountID sets the optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount). If specified as the default value
// (0.0.0), the autoRenewAccount will be removed.
func (transaction *TopicUpdateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.autoRenewAccountID = &autoRenewAccountID
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAutoRenewAccountID() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

// ClearTopicMemo explicitly clears any memo on the topic by sending an empty string as the memo
func (transaction *TopicUpdateTransaction) ClearTopicMemo() *TopicUpdateTransaction {
	return transaction.SetTopicMemo("")
}

// ClearAdminKey explicitly clears any admin key on the topic by sending an empty key list as the key
func (transaction *TopicUpdateTransaction) ClearAdminKey() *TopicUpdateTransaction {
	return transaction.SetAdminKey(PublicKey{nil})
}

// ClearSubmitKey explicitly clears any submit key on the topic by sending an empty key list as the key
func (transaction *TopicUpdateTransaction) ClearSubmitKey() *TopicUpdateTransaction {
	return transaction.SetSubmitKey(PublicKey{nil})
}

// ClearAutoRenewAccountID explicitly clears any auto renew account ID on the topic by sending an empty accountID
func (transaction *TopicUpdateTransaction) ClearAutoRenewAccountID() *TopicUpdateTransaction {
	transaction.autoRenewAccountID = &AccountID{}
	return transaction
}

func (transaction *TopicUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.topicID != nil {
		if err := transaction.topicID.Validate(client); err != nil {
			return err
		}
	}

	if transaction.autoRenewAccountID != nil {
		if err := transaction.autoRenewAccountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TopicUpdateTransaction) build() *proto.TransactionBody {
	body := &proto.ConsensusUpdateTopicTransactionBody{
		Memo: &wrappers.StringValue{Value: transaction.memo},
	}

	if !transaction.topicID.isZero() {
		body.TopicID = transaction.topicID.toProtobuf()
	}

	if !transaction.autoRenewAccountID.isZero() {
		body.AutoRenewAccount = transaction.autoRenewAccountID.toProtobuf()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = timeToProtobuf(*transaction.expirationTime)
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
		Data: &proto.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: body,
		},
	}
}

func (transaction *TopicUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TopicUpdateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.ConsensusUpdateTopicTransactionBody{
		Memo: &wrappers.StringValue{Value: transaction.memo},
	}

	if !transaction.topicID.isZero() {
		body.TopicID = transaction.topicID.toProtobuf()
	}

	if !transaction.autoRenewAccountID.isZero() {
		body.AutoRenewAccount = transaction.autoRenewAccountID.toProtobuf()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = timeToProtobuf(*transaction.expirationTime)
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
		Data: &proto.SchedulableTransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: body,
		},
	}, nil
}

func _TopicUpdateTransactionGetMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getTopic().UpdateTopic,
	}
}

func (transaction *TopicUpdateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TopicUpdateTransaction) Sign(
	privateKey PrivateKey,
) *TopicUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TopicUpdateTransaction) SignWithOperator(
	client *Client,
) (*TopicUpdateTransaction, error) {
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
func (transaction *TopicUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicUpdateTransaction {
	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TopicUpdateTransaction) Execute(
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
		_TopicUpdateTransactionGetMethod,
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

func (transaction *TopicUpdateTransaction) Freeze() (*TopicUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicUpdateTransaction) FreezeWith(client *Client) (*TopicUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TopicUpdateTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TopicUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionMemo(memo string) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionID(transactionID TransactionID) *TopicUpdateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TopicUpdateTransaction) SetMaxRetry(count int) *TopicUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TopicUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicUpdateTransaction {
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

func (transaction *TopicUpdateTransaction) SetMaxBackoff(max time.Duration) *TopicUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TopicUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TopicUpdateTransaction) SetMinBackoff(min time.Duration) *TopicUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TopicUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
