package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TopicInfoQuery struct {
	Query
	pb      *proto.ConsensusGetTopicInfoQuery
	topicID TopicID
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//  Get Topic Info Query.
func NewTopicInfoQuery() *TopicInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.ConsensusGetTopicInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_ConsensusGetTopicInfo{
		ConsensusGetTopicInfo: &pb,
	}

	return &TopicInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (query *TopicInfoQuery) SetTopicID(id TopicID) *TopicInfoQuery {
	query.topicID = id
	return query
}

func (query *TopicInfoQuery) GetTopicID() TopicID {
	return query.topicID
}

func (query *TopicInfoQuery) validateChecksums(client *Client) error {
	var err error
	err = query.topicID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *TopicInfoQuery) build() *TopicInfoQuery {
	if !query.topicID.isZero() {
		query.pb.TopicID = query.topicID.toProtobuf()
	}

	return query
}

func (query *TopicInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return Hbar{}, err
	}

	query.pbHeader.Payment = paymentTransaction
	query.pbHeader.ResponseType = proto.ResponseType_COST_ANSWER
	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	err = query.validateChecksums(client)
	if err != nil {
		return Hbar{}, err
	}

	query.build()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		topicInfoQuery_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		topicInfoQuery_getMethod,
		topicInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetConsensusGetTopicInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func topicInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetConsensusGetTopicInfo().Header.NodeTransactionPrecheckCode))
}

func topicInfoQuery_mapStatusError(_ request, response response, _ *NetworkName) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetConsensusGetTopicInfo().Header.NodeTransactionPrecheckCode),
	}
}

func topicInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getTopic().GetTopicInfo,
	}
}

// Execute executes the TopicInfoQuery using the provided client
func (query *TopicInfoQuery) Execute(client *Client) (TopicInfo, error) {
	if client == nil || client.operator == nil {
		return TopicInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateChecksums(client)
	if err != nil {
		return TopicInfo{}, err
	}

	query.build()

	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.maxQueryPayment
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return TopicInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TopicInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TopicInfoQuery",
			}
		}

		cost = actualCost
	}

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return TopicInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		topicInfoQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		topicInfoQuery_getMethod,
		topicInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return TopicInfo{}, err
	}

	return topicInfoFromProtobuf(resp.query.GetConsensusGetTopicInfo().TopicInfo, client.networkName)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TopicInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TopicInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *TopicInfoQuery) SetQueryPayment(paymentAmount Hbar) *TopicInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *TopicInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TopicInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TopicInfoQuery) SetMaxRetry(count int) *TopicInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}
