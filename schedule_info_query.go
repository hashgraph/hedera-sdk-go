package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ScheduleInfoQuery struct {
	Query
	scheduleID *ScheduleID
}

func NewScheduleInfoQuery() *ScheduleInfoQuery {
	return &ScheduleInfoQuery{
		Query: _NewQuery(true),
	}
}

func (query *ScheduleInfoQuery) SetScheduleID(scheduleID ScheduleID) *ScheduleInfoQuery {
	query.scheduleID = &scheduleID
	return query
}

func (query *ScheduleInfoQuery) GetScheduleID() ScheduleID {
	if query.scheduleID == nil {
		return ScheduleID{}
	}

	return *query.scheduleID
}

func (query *ScheduleInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.scheduleID != nil {
		if err := query.scheduleID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *ScheduleInfoQuery) _Build() *proto.Query_ScheduleGetInfo {
	body := &proto.ScheduleGetInfoQuery{
		Header: &proto.QueryHeader{},
	}

	if query.scheduleID != nil {
		body.ScheduleID = query.scheduleID._ToProtobuf()
	}

	return &proto.Query_ScheduleGetInfo{
		ScheduleGetInfo: body,
	}
}

func (query *ScheduleInfoQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ScheduleGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ScheduleGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *ScheduleInfoQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.ScheduleGetInfo.Header.Payment = paymentTransaction
	pb.ScheduleGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *ScheduleInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network._GetNodeAccountIDsForExecute()

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query._CostQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ScheduleInfoQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_ScheduleInfoQueryGetMethod,
		_ScheduleInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetScheduleGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _ScheduleInfoQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetScheduleGetInfo().Header.NodeTransactionPrecheckCode))
}

func _ScheduleInfoQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetScheduleGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _ScheduleInfoQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetSchedule().GetScheduleInfo,
	}
}

func (query *ScheduleInfoQuery) Execute(client *Client) (ScheduleInfo, error) {
	if client == nil || client.operator == nil {
		return ScheduleInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return ScheduleInfo{}, err
	}

	query._Build()

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

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return ScheduleInfo{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ScheduleInfoQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ScheduleInfoQueryGetMethod,
		_ScheduleInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return ScheduleInfo{}, err
	}

	return _ScheduleInfoFromProtobuf(resp.query.GetScheduleGetInfo().ScheduleInfo), nil
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

func (query *ScheduleInfoQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

func (query *ScheduleInfoQuery) SetMaxRetry(count int) *ScheduleInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *ScheduleInfoQuery) SetMaxBackoff(max time.Duration) *ScheduleInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *ScheduleInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *ScheduleInfoQuery) SetMinBackoff(min time.Duration) *ScheduleInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *ScheduleInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}
