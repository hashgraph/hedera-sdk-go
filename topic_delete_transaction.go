package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

// A ConsensusTopicDeleteTransaction is for deleting a topic on HCS.
type TopicDeleteTransaction struct {
	Transaction
	pb      *proto.ConsensusDeleteTopicTransactionBody
	topicID TopicID
}

// NewConsensusTopicDeleteTransaction creates a ConsensusTopicDeleteTransaction transaction which can be used to construct
// and execute a Consensus Delete Topic Transaction.
func NewTopicDeleteTransaction() *TopicDeleteTransaction {
	pb := &proto.ConsensusDeleteTopicTransactionBody{}

	transaction := TopicDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func topicDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TopicDeleteTransaction {
	return TopicDeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetConsensusDeleteTopic(),
	}
}

// SetTopicID sets the topic IDentifier.
func (transaction *TopicDeleteTransaction) SetTopicID(ID TopicID) *TopicDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.TopicID = ID.toProtobuf()
	return transaction
}

func (transaction *TopicDeleteTransaction) GetTopicID() TopicID {
	return transaction.topicID
}

func (transaction *TopicDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.topicID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *TopicDeleteTransaction) build() *TopicDeleteTransaction {
	if !transaction.topicID.isZero() {
		transaction.pb.TopicID = transaction.topicID.toProtobuf()
	}

	return transaction
}

func (transaction *TopicDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TopicDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: &proto.ConsensusDeleteTopicTransactionBody{
				TopicID: transaction.pb.TopicID,
			},
		},
	}, nil
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
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
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
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		topicDeleteTransaction_getMethod,
		transaction_mapStatusError,
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
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TopicDeleteTransaction{}, err
	}
	transaction.build()

	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

// SetMaxTransactionFee sets the max transaction fee for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetMaxTransactionFee(fee Hbar) *TopicDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetTransactionMemo sets the memo for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionMemo(memo string) *TopicDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// SetTransactionValidDuration sets the valid duration for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *TopicDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// SetTransactionID sets the TransactionID for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetTransactionID(transactionID TransactionID) *TopicDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this TopicDeleteTransaction.
func (transaction *TopicDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *TopicDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TopicDeleteTransaction) SetMaxRetry(count int) *TopicDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TopicDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *TopicDeleteTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
