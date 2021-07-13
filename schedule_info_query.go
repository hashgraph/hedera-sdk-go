package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type ScheduleInfoQuery struct {
	Query
	pb         *services.ScheduleGetInfoQuery
	scheduleID ScheduleID
}

func NewScheduleInfoQuery() *ScheduleInfoQuery {
	header := services.QueryHeader{}
	query := newQuery(true, &header)
	pb := services.ScheduleGetInfoQuery{Header: &header}
	query.pb.Query = &services.Query_ScheduleGetInfo{
		ScheduleGetInfo: &pb,
	}

	return &ScheduleInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func (query *ScheduleInfoQuery) SetScheduleID(id ScheduleID) *ScheduleInfoQuery {
	query.scheduleID = id
	return query
}

func (query *ScheduleInfoQuery) GetScheduleID(id ScheduleID) ScheduleID {
	return query.scheduleID
}

func (query *ScheduleInfoQuery) validateNetworkOnIDs(client *Client) error {
	var err error
	err = query.scheduleID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *ScheduleInfoQuery) build() *ScheduleInfoQuery {
	if !query.scheduleID.isZero() {
		query.pb.ScheduleID = query.scheduleID.toProtobuf()
	}

	return query
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
	query.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	err = query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	query.build()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		scheduleInfoQuery_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		scheduleInfoQuery_getMethod,
		scheduleInfoQuery_mapStatusError,
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

func scheduleInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetScheduleGetInfo().Header.NodeTransactionPrecheckCode))
}

func scheduleInfoQuery_mapStatusError(_ request, response response, _ *NetworkName) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetScheduleGetInfo().Header.NodeTransactionPrecheckCode),
	}
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

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return ScheduleInfo{}, err
	}

	query.build()

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

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return ScheduleInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		scheduleInfoQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		scheduleInfoQuery_getMethod,
		scheduleInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return ScheduleInfo{}, err
	}

	return scheduleInfoFromProtobuf(resp.query.GetScheduleGetInfo().ScheduleInfo, client.networkName), nil
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
