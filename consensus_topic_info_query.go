package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
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

// NewConsensusTopicInfoQuery creates a ConsensusTopicInfoQuery builder which can be used to construct and execute a
// Consensus Get Topic Info Query.
func NewConsensusTopicInfoQuery() *ConsensusTopicInfoQuery {
	pb := &proto.ConsensusGetTopicInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ConsensusGetTopicInfo{pb}

	return &ConsensusTopicInfoQuery{inner, pb}
}

// SetTopicId sets the topic to retrieve info about (the parameters and running state of).
func (builder *ConsensusTopicInfoQuery) SetTopicID(id ConsensusTopicID) *ConsensusTopicInfoQuery {
	builder.pb.TopicID = id.toProto()
	return builder
}

// Execute executes the ConsensusTopicInfoQuery using the provided client
func (builder *ConsensusTopicInfoQuery) Execute(client *Client) (ConsensusTopicInfo, error) {
	resp, err := builder.execute(client)
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

func (builder *ConsensusTopicInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ConsensusTopicInfoQuery {
	return &ConsensusTopicInfoQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *ConsensusTopicInfoQuery) SetQueryPayment(paymentAmount Hbar) *ConsensusTopicInfoQuery {
	return &ConsensusTopicInfoQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *ConsensusTopicInfoQuery) SetQueryPaymentTransaction(tx Transaction) *ConsensusTopicInfoQuery {
	return &ConsensusTopicInfoQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
