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

// NetworkVersionInfoQuery is the query to be executed that would return the current version of the network's protobuf and services.
type NetworkVersionInfoQuery struct {
	query
}

// NewNetworkVersionQuery creates a NetworkVersionInfoQuery builder which can be used to construct and execute a
// Network Get Version Info Query containing the current version of the network's protobuf and services.
func NewNetworkVersionQuery() *NetworkVersionInfoQuery {
	header := services.QueryHeader{}
	result := NetworkVersionInfoQuery{
		query: _NewQuery(true, &header),
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *NetworkVersionInfoQuery) SetGrpcDeadline(deadline *time.Duration) *NetworkVersionInfoQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *NetworkVersionInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}
	pb := services.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &services.NetworkGetVersionInfoQuery{},
	}
	pb.NetworkGetVersionInfo.Header = this.pbHeader

	this.pb = &services.Query{
		Query: &pb,
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

	cost := int64(resp.(*services.Response).GetNetworkGetVersionInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (this *NetworkVersionInfoQuery) Execute(client *Client) (NetworkVersionInfo, error) {
	if client == nil || client.operator == nil {
		return NetworkVersionInfo{}, errNoClientProvided
	}

	var err error

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
			return NetworkVersionInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return NetworkVersionInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "NetworkVersionInfo",
			}
		}

		cost = actualCost
	}

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return NetworkVersionInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return NetworkVersionInfo{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := services.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &services.NetworkGetVersionInfoQuery{},
	}
	pb.NetworkGetVersionInfo.Header = this.pbHeader

	this.pb = &services.Query{
		Query: &pb,
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
		return NetworkVersionInfo{}, err
	}

	return _NetworkVersionInfoFromProtobuf(resp.(*services.Response).GetNetworkGetVersionInfo()), err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *NetworkVersionInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *NetworkVersionInfoQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *NetworkVersionInfoQuery) SetQueryPayment(paymentAmount Hbar) *NetworkVersionInfoQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this NetworkVersionInfoQuery.
func (this *NetworkVersionInfoQuery) SetNodeAccountIDs(accountID []AccountID) *NetworkVersionInfoQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *NetworkVersionInfoQuery) SetMaxRetry(count int) *NetworkVersionInfoQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *NetworkVersionInfoQuery) SetMaxBackoff(max time.Duration) *NetworkVersionInfoQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *NetworkVersionInfoQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *NetworkVersionInfoQuery) SetMinBackoff(min time.Duration) *NetworkVersionInfoQuery {
	this.query.SetMaxBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *NetworkVersionInfoQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *NetworkVersionInfoQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("NetworkVersionInfoQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *NetworkVersionInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *NetworkVersionInfoQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *NetworkVersionInfoQuery) SetLogLevel(level LogLevel) *NetworkVersionInfoQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *NetworkVersionInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetNetwork().GetVersionInfo,
	}
}

func (this *NetworkVersionInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (this *NetworkVersionInfoQuery) getName() string {
	return "NetworkVersionInfoQuery"
}
func (this *NetworkVersionInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode)
}
