package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ScheduleInfoQuery struct {
	Query
	scheduleID ScheduleID
}

func NewScheduleInfoQuery() *ScheduleInfoQuery {
	return &ScheduleInfoQuery{
		Query: newQuery(true),
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
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.scheduleID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *ScheduleInfoQuery) build() *proto.Query_ScheduleGetInfo {
	body := &proto.ScheduleGetInfoQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.scheduleID.isZero() {
		body.ScheduleID = query.scheduleID.toProtobuf()
	}

	return &proto.Query_ScheduleGetInfo{
		ScheduleGetInfo: body,
	}
}

func (query *ScheduleInfoQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ScheduleGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ScheduleGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *ScheduleInfoQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.ScheduleGetInfo.Header.Payment = paymentTransaction
	pb.ScheduleGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *ScheduleInfoQuery) GetCost(client *Client) (Hbar, error) {
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
		scheduleInfoQuery_shouldRetry,
		protoReq,
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

func scheduleInfoQuery_mapStatusError(_ request, response response) error {
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
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		scheduleInfoQuery_getMethod,
		scheduleInfoQuery_mapStatusError,
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
