package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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

type NetworkVersionInfoQuery struct {
	Query
}

func NewNetworkVersionQuery() *NetworkVersionInfoQuery {
	header := services.QueryHeader{}
	return &NetworkVersionInfoQuery{
		Query: _NewQuery(true, &header),
	}
}

func (query *NetworkVersionInfoQuery) SetGrpcDeadline(deadline *time.Duration) *NetworkVersionInfoQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

func (query *NetworkVersionInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}
	pb := services.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &services.NetworkGetVersionInfoQuery{},
	}
	pb.NetworkGetVersionInfo.Header = query.pbHeader

	query.pb = &services.Query{
		Query: &pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_NetworkVersionInfoQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_NetworkVersionInfoQueryGetMethod,
		_NetworkVersionInfoQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetNetworkGetVersionInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _NetworkVersionInfoQueryShouldRetry(logID string, _ interface{}, response interface{}) _ExecutionState {
	return _QueryShouldRetry(logID, Status(response.(*services.Response).GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode))
}

func _NetworkVersionInfoQueryMapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _NetworkVersionInfoQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetNetwork().GetVersionInfo,
	}
}

func (query *NetworkVersionInfoQuery) Execute(client *Client) (NetworkVersionInfo, error) {
	if client == nil || client.operator == nil {
		return NetworkVersionInfo{}, errNoClientProvided
	}

	var err error

	if !query.paymentTransactionIDs.locked {
		query.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
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

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return NetworkVersionInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(query.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return NetworkVersionInfo{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := services.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &services.NetworkGetVersionInfoQuery{},
	}
	pb.NetworkGetVersionInfo.Header = query.pbHeader

	query.pb = &services.Query{
		Query: &pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_NetworkVersionInfoQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_NetworkVersionInfoQueryGetMethod,
		_NetworkVersionInfoQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return NetworkVersionInfo{}, err
	}

	return _NetworkVersionInfoFromProtobuf(resp.(*services.Response).GetNetworkGetVersionInfo()), err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *NetworkVersionInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *NetworkVersionInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *NetworkVersionInfoQuery) SetQueryPayment(paymentAmount Hbar) *NetworkVersionInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *NetworkVersionInfoQuery) SetNodeAccountIDs(accountID []AccountID) *NetworkVersionInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *NetworkVersionInfoQuery) SetMaxRetry(count int) *NetworkVersionInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *NetworkVersionInfoQuery) SetMaxBackoff(max time.Duration) *NetworkVersionInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *NetworkVersionInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *NetworkVersionInfoQuery) SetMinBackoff(min time.Duration) *NetworkVersionInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *NetworkVersionInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *NetworkVersionInfoQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("NetworkVersionInfoQuery:%d", timestamp)
}

func (query *NetworkVersionInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *NetworkVersionInfoQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}
