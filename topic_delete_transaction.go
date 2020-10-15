package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

// A ConsensusTopicDeleteTransaction is for deleting a topic on HCS.
type TopicDeleteTransaction struct {
	Transaction
	pb *proto.ConsensusDeleteTopicTransactionBody
}

// NewConsensusTopicDeleteTransaction creates a ConsensusTopicDeleteTransaction transaction which can be used to construct
// and execute a Consensus Delete Topic Transaction.
func NewConsensusTopicDeleteTransaction() *TopicDeleteTransaction {
	pb := &proto.ConsensusDeleteTopicTransactionBody{}

	transaction := TopicDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

// SetTopicID sets the topic IDentifier.
func (transaction *TopicDeleteTransaction) SetTopicID(ID TopicID) *TopicDeleteTransaction {
	transaction.pb.TopicID = ID.toProto()
	return transaction
}

func (transaction *TopicDeleteTransaction) GetTopicID() TopicID {
	return TopicIDFromProto(transaction.pb.GetTopicID())
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func topicDeleteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getTopic().DeleteTopic,
	}
}

func (transaction *TopicDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
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
func (transaction *TopicDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicDeleteTransaction {
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
func (transaction *TopicDeleteTransaction) Execute(
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
		topicDeleteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *TopicDeleteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ConsensusDeleteTopic{
		ConsensusDeleteTopic: transaction.pb,
	}

	return true
}

func (transaction *TopicDeleteTransaction) Freeze() (*TopicDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TopicDeleteTransaction) FreezeWith(client *Client) (*TopicDeleteTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *TopicDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetMaxTransactionFee(fee Hbar) *TopicDeleteTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TopicDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionMemo(memo string) *TopicDeleteTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TopicDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *TopicDeleteTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TopicDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionID(transactionID TransactionID) *TopicDeleteTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *TopicDeleteTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetNodeID(nodeID AccountID) *TopicDeleteTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
