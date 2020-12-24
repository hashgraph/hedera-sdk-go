package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// A TopicCreateTransaction is for creating a new Topic on HCS.
type TopicCreateTransaction struct {
	Transaction
	pb *proto.ConsensusCreateTopicTransactionBody
}

// NewTopicCreateTransaction creates a TopicCreateTransaction transaction which can be
// used to construct and execute a  Create Topic Transaction.
func NewTopicCreateTransaction() *TopicCreateTransaction {
	pb := &proto.ConsensusCreateTopicTransactionBody{}

	transaction := TopicCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(2))

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	//transaction.SetReceiveRecordThreshold(MaxHbar)
	//transaction.SetSendRecordThreshold(MaxHbar)

	return &transaction

}

func topicCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TopicCreateTransaction {
	return TopicCreateTransaction{
		Transaction: transaction,
		pb:          pb.GetConsensusCreateTopic(),
	}
}

// SetAdminKey sets the key required to update or delete the topic. If unspecified, anyone can increase the topic's
// expirationTime.
func (transaction *TopicCreateTransaction) SetAdminKey(publicKey Key) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AdminKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TopicCreateTransaction) GetAdminKey() (Key, error) {
	return keyFromProtobuf(transaction.pb.GetAdminKey())
}

// SetSubmitKey sets the key required for submitting messages to the topic. If unspecified, all submissions are allowed.
func (transaction *TopicCreateTransaction) SetSubmitKey(publicKey Key) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SubmitKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TopicCreateTransaction) GetSubmitKey() (Key, error) {
	return keyFromProtobuf(transaction.pb.GetSubmitKey())
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (transaction *TopicCreateTransaction) SetTopicMemo(memo string) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = memo
	return transaction
}

func (transaction *TopicCreateTransaction) GetTopicMemo() string {
	return transaction.pb.GetMemo()
}

// SetAutoRenewPeriod sets the initial lifetime of the topic and the amount of time to extend the topic's lifetime
// automatically at expirationTime if the autoRenewAccount is configured and has sufficient funds.
//
// Required. Limited to a maximum of 90 days (server-sIDe configuration which may change).
func (transaction *TopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = durationToProtobuf(period)
	return transaction
}

func (transaction *TopicCreateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProtobuf(transaction.pb.GetAutoRenewPeriod())
}

// SetAutoRenewAccountID sets an optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount).
//
//If specified, there must be an adminKey and the autoRenewAccount must sign this transaction.
func (transaction *TopicCreateTransaction) SetAutoRenewAccountID(accountID AccountID) *TopicCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewAccount = accountID.toProtobuf()
	return transaction
}

func (transaction *TopicCreateTransaction) GetAutoRenewAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAutoRenewAccount())
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func topicCreateTransaction_getMethod(request request, channel *channel) method {
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
func (transaction *TopicCreateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
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

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
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
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		topicCreateTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()

	return TransactionResponse{
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
}

func (transaction *TopicCreateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ConsensusCreateTopic{
		ConsensusCreateTopic: transaction.pb,
	}

	return true
}

func (transaction *TopicCreateTransaction) Freeze() (*TopicCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicCreateTransaction) FreezeWith(client *Client) (*TopicCreateTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
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

func (transaction *TopicCreateTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
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
