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
	"fmt"
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
func (this *ScheduleInfoQuery) SetGrpcDeadline(deadline *time.Duration) *ScheduleInfoQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetScheduleID Sets the id of the schedule to interrogate
func (this *ScheduleInfoQuery) SetScheduleID(scheduleID ScheduleID) *ScheduleInfoQuery {
	this.scheduleID = &scheduleID
	return this
}

// GetScheduleID returns the id of the schedule to interrogate
func (this *ScheduleInfoQuery) GetScheduleID() ScheduleID {
	if this.scheduleID == nil {
		return ScheduleID{}
	}

	return *this.scheduleID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *ScheduleInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range this.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.ScheduleGetInfo.Header = this.pbHeader

	this.pb = &services.Query{
		Query: pb,
	}

	this.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	this.paymentTransactionIDs._Advance()


	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetScheduleGetInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (this *ScheduleInfoQuery) Execute(client *Client) (ScheduleInfo, error) {
	if client == nil || client.operator == nil {
		return ScheduleInfo{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return ScheduleInfo{}, err
	}

	if !this.paymentTransactionIDs.locked {
		this.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if this.queryPayment.tinybar != 0 {
		cost = this.queryPayment
	} else {
		if this.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = this.maxQueryPayment
		}

		actualCost, err := this.GetCost(client)
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

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return ScheduleInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return ScheduleInfo{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.ScheduleGetInfo.Header = this.pbHeader
	this.pb = &services.Query{
		Query: pb,
	}

	if this.isPaymentRequired && len(this.paymentTransactions) > 0 {
		this.paymentTransactionIDs._Advance()
	}
	this.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return ScheduleInfo{}, err
	}

	return _ScheduleInfoFromProtobuf(resp.(*services.Response).GetScheduleGetInfo().ScheduleInfo), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *ScheduleInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ScheduleInfoQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *ScheduleInfoQuery) SetQueryPayment(paymentAmount Hbar) *ScheduleInfoQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this ScheduleInfoQuery.
func (this *ScheduleInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ScheduleInfoQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// GetNodeAccountIDs returns the _Node AccountID for this ScheduleInfoQuery.
func (this *ScheduleInfoQuery) GetNodeAccountIDs() []AccountID {
	return this.query.GetNodeAccountIDs()
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *ScheduleInfoQuery) SetMaxRetry(count int) *ScheduleInfoQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *ScheduleInfoQuery) SetMaxBackoff(max time.Duration) *ScheduleInfoQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *ScheduleInfoQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *ScheduleInfoQuery) SetMinBackoff(min time.Duration) *ScheduleInfoQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *ScheduleInfoQuery) GetMinBackoff() time.Duration {
	return this.GetMinBackoff()
}

func (this *ScheduleInfoQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("ScheduleInfoQuery:%d", timestamp)
}

func (this *ScheduleInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *ScheduleInfoQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *ScheduleInfoQuery) SetLogLevel(level LogLevel) *ScheduleInfoQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *ScheduleInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetSchedule().GetScheduleInfo,
	}
}

func (this *ScheduleInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetScheduleGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (this *ScheduleInfoQuery) getName() string {
	return "ScheduleInfoQuery"
}

func (this *ScheduleInfoQuery) build() *services.Query_ScheduleGetInfo {
	body := &services.ScheduleGetInfoQuery{
		Header: &services.QueryHeader{},
	}

	if this.scheduleID != nil {
		body.ScheduleID = this.scheduleID._ToProtobuf()
	}

	return &services.Query_ScheduleGetInfo{
		ScheduleGetInfo: body,
	}
}

func (this *ScheduleInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.scheduleID != nil {
		if err := this.scheduleID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *ScheduleInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetScheduleGetInfo().Header.NodeTransactionPrecheckCode)
}