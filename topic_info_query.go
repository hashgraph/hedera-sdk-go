package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TopicInfoQuery struct {
	Query
	topicID TopicID
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//  Get Topic Info Query.
func NewTopicInfoQuery() *TopicInfoQuery {
	return &TopicInfoQuery{
		Query: newQuery(true),
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

func (query *TopicInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.topicID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *TopicInfoQuery) build() *proto.Query_ConsensusGetTopicInfo {
	body := &proto.ConsensusGetTopicInfoQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.topicID.isZero() {
		body.TopicID = query.topicID.toProtobuf()
	}

	return &proto.Query_ConsensusGetTopicInfo{
		ConsensusGetTopicInfo: body,
	}
}

func (query *TopicInfoQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ConsensusGetTopicInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ConsensusGetTopicInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TopicInfoQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.ConsensusGetTopicInfo.Header.Payment = paymentTransaction
	pb.ConsensusGetTopicInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TopicInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query.costQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		topicInfoQuery_shouldRetry,
		protoReq,
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

func topicInfoQuery_mapStatusError(_ request, response response) error {
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

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return TopicInfo{}, err
	}

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
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		topicInfoQuery_getMethod,
		topicInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return TopicInfo{}, err
	}

	return topicInfoFromProtobuf(resp.query.GetConsensusGetTopicInfo().TopicInfo)
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
