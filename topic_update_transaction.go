package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// *TopicUpdateTransaction updates all fields on a Topic that are set in the transaction.
type TopicUpdateTransaction struct {
	Transaction
	pb *proto.ConsensusUpdateTopicTransactionBody
}

// NewTopicUpdateTransaction creates a *TopicUpdateTransaction transaction which can be
// used to construct and execute a  Update Topic Transaction.
func NewTopicUpdateTransaction() *TopicUpdateTransaction {
	pb := &proto.ConsensusUpdateTopicTransactionBody{}

	transaction := TopicUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)

	return &transaction
}

// SetTopicID sets the topic to be updated.
func (transaction *TopicUpdateTransaction) SetTopicID(topicID TopicID) *TopicUpdateTransaction {
	transaction.pb.TopicID = topicID.toProto()
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTopicID() TopicID {
	return TopicIDFromProto(transaction.pb.GetTopicID())
}

// SetAdminKey sets the key required to update/delete the topic. If unset, the key will not be changed.
//
// Setting the AdminKey to an empty KeyList will clear the adminKey.
func (transaction *TopicUpdateTransaction) SetAdminKey(publicKey PublicKey) *TopicUpdateTransaction {
	transaction.pb.AdminKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAdminKey() (Key, error) {
	return publicKeyFromProto(transaction.pb.GetAdminKey())
}

// SetSubmitKey will set the key allowed to submit messages to the topic.  If unset, the key will not be changed.
//
// Setting the submitKey to an empty KeyList will clear the submitKey.
func (transaction *TopicUpdateTransaction) SetSubmitKey(publicKey PublicKey) *TopicUpdateTransaction {
	transaction.pb.SubmitKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TopicUpdateTransaction) GetSubmitKey() (Key, error) {
	return publicKeyFromProto(transaction.pb.GetSubmitKey())
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction *TopicUpdateTransaction) SetTopicMemo(memo string) *TopicUpdateTransaction {
	transaction.pb.Memo = &proto.StringValue{Value: memo}
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTopicMemo() string {
	return transaction.pb.GetMemo().GetValue()
}

// SetExpirationTime sets the effective  timestamp at (and after) which all  transactions and queries
// will fail. The expirationTime may be no longer than 90 days from the  timestamp of this transaction.
func (transaction *TopicUpdateTransaction) SetExpirationTime(expiration time.Time) *TopicUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetExpirationTime() time.Time {
	return timeFromProto(transaction.pb.GetExpirationTime())
}

// SetAutoRenewPeriod sets the amount of time to extend the topic's lifetime automatically at expirationTime if the
// autoRenewAccount is configured and has funds. This is limited to a maximum of 90 days (server-sIDe configuration
// which may change).
func (transaction *TopicUpdateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicUpdateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(period)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProto(transaction.pb.GetAutoRenewPeriod())
}

// SetAutoRenewAccountID sets the optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount). If specified as the default value
// (0.0.0), the autoRenewAccount will be removed.
func (transaction *TopicUpdateTransaction) SetAutoRenewAccountID(accountID AccountID) *TopicUpdateTransaction {
	transaction.pb.AutoRenewAccount = accountID.toProtobuf()
	return transaction
}

func (transaction *TopicUpdateTransaction) GetAutoRenewAccountID() AccountID {
	return accountIDFromProto(transaction.pb.GetAutoRenewAccount())
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
	transaction.pb.AutoRenewAccount = &proto.AccountID{}

	return transaction
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func topicUpdateTransaction_getMethod(request request, channel *channel) method {
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

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *TopicUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicUpdateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.transactions); index++ {
		signature := signer(transaction.transactions[index].GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TopicUpdateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.id

	if !client.GetOperatorID().isZero() && client.GetOperatorID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorKey(),
			client.operator.signer,
		)
	}

	_, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeId,
		topicUpdateTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
		query_makePaymentTransaction,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *TopicUpdateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ConsensusUpdateTopic{
		ConsensusUpdateTopic: transaction.pb,
	}

	return true
}

func (transaction *TopicUpdateTransaction) Freeze() (*TopicUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicUpdateTransaction) FreezeWith(client *Client) (*TopicUpdateTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *TopicUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TopicUpdateTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionMemo(memo string) *TopicUpdateTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TopicUpdateTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetTransactionID(transactionID TransactionID) *TopicUpdateTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *TopicUpdateTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this TopicUpdateTransaction.
func (transaction *TopicUpdateTransaction) SetNodeID(nodeID AccountID) *TopicUpdateTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
