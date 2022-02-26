package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type ScheduleInfoQuery struct {
	Query
	scheduleID *ScheduleID
}

func NewScheduleInfoQuery() *ScheduleInfoQuery {
	header := services.QueryHeader{}
	return &ScheduleInfoQuery{
		Query: _NewQuery(true, &header),
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
		if err := query.scheduleID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *ScheduleInfoQuery) _Build() *services.Query_ScheduleGetInfo {
	body := &services.ScheduleGetInfoQuery{
		Header: &services.QueryHeader{},
	}

	if query.scheduleID != nil {
		body.ScheduleID = query.scheduleID._ToProtobuf()
	}

	return &services.Query_ScheduleGetInfo{
		ScheduleGetInfo: body,
	}
}

func (query *ScheduleInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return Hbar{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range query.nodeAccountIDs {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.ScheduleGetInfo.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ScheduleInfoQueryShouldRetry,
		_CostQueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ScheduleInfoQueryGetMethod,
		_ScheduleInfoQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
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

func _ScheduleInfoQueryShouldRetry(logID string, _ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(logID, Status(response.query.GetScheduleGetInfo().Header.NodeTransactionPrecheckCode))
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

	var err error

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return ScheduleInfo{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}
	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return ScheduleInfo{}, err
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
			return ScheduleInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return ScheduleInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "ScheduleInfo",
			}
		}

		cost = actualCost
	}

	query.nextPaymentTransactionIndex = 0
	query.paymentTransactions = make([]*services.Transaction, 0)

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return ScheduleInfo{}, err
	}

	pb := query._Build()
	pb.ScheduleGetInfo.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ScheduleInfoQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ScheduleInfoQueryGetMethod,
		_ScheduleInfoQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
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

func (query *ScheduleInfoQuery) _GetLogID() string {
	timestamp := query.paymentTransactionID.ValidStart
	return fmt.Sprintf("ScheduleInfoQuery:%d", timestamp.UnixNano())
}
