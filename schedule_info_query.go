package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ScheduleInfoQuery Gets information about a schedule in the network's action queue.
type ScheduleInfoQuery struct {
	Query
	scheduleID *ScheduleID
}

// NewScheduleInfoQuery creates ScheduleInfoQuery which gets information about a schedule in the network's action queue.
func NewScheduleInfoQuery() *ScheduleInfoQuery {
	header := services.QueryHeader{}
	return &ScheduleInfoQuery{
		Query: _NewQuery(true, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *ScheduleInfoQuery) SetGrpcDeadline(deadline *time.Duration) *ScheduleInfoQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetScheduleID Sets the id of the schedule to interrogate
func (q *ScheduleInfoQuery) SetScheduleID(scheduleID ScheduleID) *ScheduleInfoQuery {
	q.scheduleID = &scheduleID
	return q
}

// GetScheduleID returns the id of the schedule to interrogate
func (q *ScheduleInfoQuery) GetScheduleID() ScheduleID {
	if q.scheduleID == nil {
		return ScheduleID{}
	}

	return *q.scheduleID
}

func (q *ScheduleInfoQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *ScheduleInfoQuery) Execute(client *Client) (ScheduleInfo, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return ScheduleInfo{}, err
	}

	return _ScheduleInfoFromProtobuf(resp.GetScheduleGetInfo().ScheduleInfo), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ScheduleInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ScheduleInfoQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ScheduleInfoQuery) SetQueryPayment(paymentAmount Hbar) *ScheduleInfoQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ScheduleInfoQuery.
func (q *ScheduleInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ScheduleInfoQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ScheduleInfoQuery) SetMaxRetry(count int) *ScheduleInfoQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ScheduleInfoQuery) SetMaxBackoff(max time.Duration) *ScheduleInfoQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ScheduleInfoQuery) SetMinBackoff(min time.Duration) *ScheduleInfoQuery {
	q.Query.SetMinBackoff(min)
	return q
}

func (q *ScheduleInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *ScheduleInfoQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *ScheduleInfoQuery) SetLogLevel(level LogLevel) *ScheduleInfoQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *ScheduleInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetSchedule().GetScheduleInfo,
	}
}

func (q *ScheduleInfoQuery) getName() string {
	return "ScheduleInfoQuery"
}

func (q *ScheduleInfoQuery) buildQuery() *services.Query {
	body := &services.ScheduleGetInfoQuery{
		Header: q.pbHeader,
	}

	if q.scheduleID != nil {
		body.ScheduleID = q.scheduleID._ToProtobuf()
	}

	return &services.Query{
		Query: &services.Query_ScheduleGetInfo{
			ScheduleGetInfo: body,
		},
	}
}

func (q *ScheduleInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.scheduleID != nil {
		if err := q.scheduleID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *ScheduleInfoQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetScheduleGetInfo()
}
