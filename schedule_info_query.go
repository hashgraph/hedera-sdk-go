package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ScheduleInfoQuery struct {
	Query
	pb *proto.ScheduleGetInfoQuery
}

func NewScheduleInfoQuery() *ScheduleInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.ScheduleGetInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_ScheduleGetInfo{
		ScheduleGetInfo: &pb,
	}

	return &ScheduleInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func (query *ScheduleInfoQuery) SetScheduleID(id ScheduleID) *ScheduleInfoQuery {
	query.pb.ScheduleID = id.toProtobuf()
	return query
}

func (query *ScheduleInfoQuery) GetScheduleID(id ScheduleID) ScheduleID {
	return scheduleIDFromProtobuf(query.pb.GetScheduleID())
}

func (query *ScheduleInfoQuery) GetCost(client *Client) (Hbar, error) {
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

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		scheduleInfoQuery_getMethod,
		scheduleInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetScheduleGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func scheduleInfoQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetScheduleGetInfo().Header.NodeTransactionPrecheckCode)
}

func scheduleInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getSchedule().GetScheduleInfo,
	}
}

func (query *ScheduleInfoQuery) Execute(client *Client) (ScheduleInfo, error) {
	if client == nil || client.operator == nil {
		return ScheduleInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		cost = client.maxQueryPayment

		actualCost, err := query.GetCost(client)
		if err != nil {
			return ScheduleInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return ScheduleInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "ScheduleInfoQuery",
			}
		}

		cost = actualCost
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return ScheduleInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		scheduleInfoQuery_getMethod,
		scheduleInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return ScheduleInfo{}, err
	}

	return scheduleInfoFromProtobuf(resp.query.GetScheduleGetInfo().ScheduleInfo), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *ScheduleInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ScheduleInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *ScheduleInfoQuery) SetQueryPayment(paymentAmount Hbar) *ScheduleInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *ScheduleInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ScheduleInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ScheduleInfoQuery) GetNodeAccountId() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

func (query *ScheduleInfoQuery) SetMaxRetry(count int) *ScheduleInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}
