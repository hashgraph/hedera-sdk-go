package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// A ConsensusTopicDeleteTransaction is for deleting a topic on HCS.
type ConsensusTopicDeleteTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusDeleteTopicTransactionBody
}

// NewConsensusTopicDeleteTransaction creates a ConsensusTopicDeleteTransaction transaction which can be used to construct
// and execute a Consensus Delete Topic Transaction.
func NewConsensusTopicDeleteTransaction() ConsensusTopicDeleteTransaction {
	pb := &proto.ConsensusDeleteTopicTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusDeleteTopic{ConsensusDeleteTopic: pb}

	transaction := ConsensusTopicDeleteTransaction{inner, pb}

	return transaction
}

// SetTopicID sets the topic identifier.
func (transaction ConsensusTopicDeleteTransaction) SetTopicID(id ConsensusTopicID) ConsensusTopicDeleteTransaction {
	transaction.pb.TopicID = id.toProto()
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction ConsensusTopicDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusTopicDeleteTransaction {
	return ConsensusTopicDeleteTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction ConsensusTopicDeleteTransaction) SetTransactionMemo(memo string) ConsensusTopicDeleteTransaction {
	return ConsensusTopicDeleteTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction ConsensusTopicDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusTopicDeleteTransaction {
	return ConsensusTopicDeleteTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction ConsensusTopicDeleteTransaction) SetTransactionID(transactionID TransactionID) ConsensusTopicDeleteTransaction {
	return ConsensusTopicDeleteTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction ConsensusTopicDeleteTransaction) SetNodeID(nodeAccountID AccountID) ConsensusTopicDeleteTransaction {
	return ConsensusTopicDeleteTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
