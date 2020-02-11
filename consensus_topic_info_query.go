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
	AdminKey           Ed25519PublicKey
	SubmitKey          Ed25519PublicKey
	AutoRenewPeriod    time.Duration
	AutoRenewAccountID AccountID
}

func NewConsensusTopicInfoQuery() *ConsensusTopicInfoQuery {
	pb := &proto.ConsensusGetTopicInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ConsensusGetTopicInfo{pb}

	return &ConsensusTopicInfoQuery{inner, pb}
}

func (builder *ConsensusTopicInfoQuery) SetTopicID(id ConsensusTopicID) *ConsensusTopicInfoQuery {
	builder.pb.TopicID = id.toProto()
	return builder
}

func (builder *ConsensusTopicInfoQuery) Execute(client *Client) (ConsensusTopicInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return ConsensusTopicInfo{}, err
	}

	return ConsensusTopicInfo{
		Memo:           resp.GetConsensusGetTopicInfo().TopicInfo.Memo,
		RunningHash:    resp.GetConsensusGetTopicInfo().TopicInfo.RunningHash,
		SequenceNumber: resp.GetConsensusGetTopicInfo().TopicInfo.SequenceNumber,
		ExpirationTime: timeFromProto(resp.GetConsensusGetTopicInfo().TopicInfo.ExpirationTime),
		AdminKey: Ed25519PublicKey{
			keyData: resp.GetConsensusGetTopicInfo().TopicInfo.AdminKey.GetEd25519(),
		},
		SubmitKey: Ed25519PublicKey{
			keyData: resp.GetConsensusGetTopicInfo().TopicInfo.SubmitKey.GetEd25519(),
		},
		AutoRenewPeriod:    durationFromProto(resp.GetConsensusGetTopicInfo().TopicInfo.AutoRenewPeriod),
		AutoRenewAccountID: accountIDFromProto(resp.GetConsensusGetTopicInfo().TopicInfo.AutoRenewAccount),
	}, nil
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
