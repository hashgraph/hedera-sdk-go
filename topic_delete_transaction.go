package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
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

// SetTopicID sets the topic identifier.
func (transaction *TopicDeleteTransaction) SetTopicID(id TopicID) *TopicDeleteTransaction {
	transaction.pb.TopicID = id.toProto()
	return transaction
}

func (transaction *TopicDeleteTransaction) GetTopicID() TopicID {
	return TopicIDFromProto(transaction.pb.GetTopicID())
}
