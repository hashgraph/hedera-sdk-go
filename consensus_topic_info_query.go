package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ConsensusTopicInfoQuery struct {
	QueryBuilder
	pb *proto.ConsensusGetTopicInfoQuery
}

type ConsensusTopicInfo struct {
	Memo               string
	RunningHash        []byte
	SequenceNumber     uint64
	ExpirationTime     time.Time
	AdminKey           *Ed25519PublicKey
	SubmitKey          *Ed25519PublicKey
	AutoRenewPeriod    time.Duration
	AutoRenewAccountID *AccountID
}

// NewConsensusTopicInfoQuery creates a ConsensusTopicInfoQuery transaction which can be used to construct and execute a
// Consensus Get Topic Info Query.
func NewConsensusTopicInfoQuery() *ConsensusTopicInfoQuery {
	pb := &proto.ConsensusGetTopicInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ConsensusGetTopicInfo{ConsensusGetTopicInfo: pb}

	return &ConsensusTopicInfoQuery{inner, pb}
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (transaction *ConsensusTopicInfoQuery) SetTopicID(id ConsensusTopicID) *ConsensusTopicInfoQuery {
	transaction.pb.TopicID = id.toProto()
	return transaction
}

// Execute executes the ConsensusTopicInfoQuery using the provided client
func (transaction *ConsensusTopicInfoQuery) Execute(client *Client) (ConsensusTopicInfo, error) {
	resp, err := transaction.execute(client)
	if err != nil {
		return ConsensusTopicInfo{}, err
	}

	ti := resp.GetConsensusGetTopicInfo().TopicInfo

	consensusTopicInfo := ConsensusTopicInfo{
		Memo:            ti.GetMemo(),
		RunningHash:     ti.RunningHash,
		SequenceNumber:  ti.SequenceNumber,
		ExpirationTime:  timeFromProto(ti.ExpirationTime),
		AutoRenewPeriod: durationFromProto(ti.AutoRenewPeriod),
	}

	if adminKey := ti.AdminKey; adminKey != nil {
		consensusTopicInfo.AdminKey = &Ed25519PublicKey{
			keyData: adminKey.GetEd25519(),
		}
	}

	if submitKey := ti.SubmitKey; submitKey != nil {
		consensusTopicInfo.SubmitKey = &Ed25519PublicKey{
			keyData: submitKey.GetEd25519(),
		}
	}

	if ARAccountID := ti.AutoRenewAccount; ARAccountID != nil {
		ID := accountIDFromProto(ARAccountID)

		consensusTopicInfo.AutoRenewAccountID = &ID
	}

	return consensusTopicInfo, nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (transaction *ConsensusTopicInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ConsensusTopicInfoQuery {
	return &ConsensusTopicInfoQuery{*transaction.QueryBuilder.SetMaxQueryPayment(maxPayment), transaction.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (transaction *ConsensusTopicInfoQuery) SetQueryPayment(paymentAmount Hbar) *ConsensusTopicInfoQuery {
	return &ConsensusTopicInfoQuery{*transaction.QueryBuilder.SetQueryPayment(paymentAmount), transaction.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (transaction *ConsensusTopicInfoQuery) SetQueryPaymentTransaction(tx Transaction) *ConsensusTopicInfoQuery {
	return &ConsensusTopicInfoQuery{*transaction.QueryBuilder.SetQueryPaymentTransaction(tx), transaction.pb}
}
