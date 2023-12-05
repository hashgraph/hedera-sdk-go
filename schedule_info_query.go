package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// ScheduleInfoQuery Gets information about a schedule in the network's action queue.
type ScheduleInfoQuery struct {
	query
	scheduleID *ScheduleID
}

// NewScheduleInfoQuery creates ScheduleInfoQuery which gets information about a schedule in the network's action queue.
func NewScheduleInfoQuery() *ScheduleInfoQuery {
	header := services.QueryHeader{}
	result := ScheduleInfoQuery{
		query: _NewQuery(true, &header),
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *ScheduleInfoQuery) SetGrpcDeadline(deadline *time.Duration) *ScheduleInfoQuery {
	q.query.SetGrpcDeadline(deadline)
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

// Execute executes the Query with the provided client
func (q *ScheduleInfoQuery) Execute(client *Client) (ScheduleInfo, error) {
	resp, err := q.query.execute(client)

	if err != nil {
		return ScheduleInfo{}, err
	}

	return _ScheduleInfoFromProtobuf(resp.GetScheduleGetInfo().ScheduleInfo), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ScheduleInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ScheduleInfoQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ScheduleInfoQuery) SetQueryPayment(paymentAmount Hbar) *ScheduleInfoQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ScheduleInfoQuery.
func (q *ScheduleInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ScheduleInfoQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ScheduleInfoQuery) SetMaxRetry(count int) *ScheduleInfoQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ScheduleInfoQuery) SetMaxBackoff(max time.Duration) *ScheduleInfoQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ScheduleInfoQuery) SetMinBackoff(min time.Duration) *ScheduleInfoQuery {
	q.query.SetMinBackoff(min)
	return q
}

func (q *ScheduleInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *ScheduleInfoQuery {
	q.query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *ScheduleInfoQuery) SetLogLevel(level LogLevel) *ScheduleInfoQuery {
	q.query.SetLogLevel(level)
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
