package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// NetworkVersionInfoQuery is the query to be executed that would return the current version of the network's protobuf and services.
type NetworkVersionInfoQuery struct {
	Query
}

// NewNetworkVersionQuery creates a NetworkVersionInfoQuery builder which can be used to construct and execute a
// Network Get Version Info Query containing the current version of the network's protobuf and services.
func NewNetworkVersionQuery() *NetworkVersionInfoQuery {
	header := services.QueryHeader{}
	return &NetworkVersionInfoQuery{
		Query: _NewQuery(true, &header),
	}
}

// SetGrpcDeadline When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *NetworkVersionInfoQuery) SetGrpcDeadline(deadline *time.Duration) *NetworkVersionInfoQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

func (q *NetworkVersionInfoQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *NetworkVersionInfoQuery) Execute(client *Client) (NetworkVersionInfo, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return NetworkVersionInfo{}, err
	}

	return _NetworkVersionInfoFromProtobuf(resp.GetNetworkGetVersionInfo()), err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *NetworkVersionInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *NetworkVersionInfoQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *NetworkVersionInfoQuery) SetQueryPayment(paymentAmount Hbar) *NetworkVersionInfoQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this NetworkVersionInfoQuery.
func (q *NetworkVersionInfoQuery) SetNodeAccountIDs(accountID []AccountID) *NetworkVersionInfoQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *NetworkVersionInfoQuery) SetMaxRetry(count int) *NetworkVersionInfoQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *NetworkVersionInfoQuery) SetMaxBackoff(max time.Duration) *NetworkVersionInfoQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *NetworkVersionInfoQuery) SetMinBackoff(min time.Duration) *NetworkVersionInfoQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *NetworkVersionInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *NetworkVersionInfoQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *NetworkVersionInfoQuery) SetLogLevel(level LogLevel) *NetworkVersionInfoQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *NetworkVersionInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetNetwork().GetVersionInfo,
	}
}

func (q *NetworkVersionInfoQuery) getName() string {
	return "NetworkVersionInfoQuery"
}

func (q *NetworkVersionInfoQuery) buildQuery() *services.Query {
	pb := services.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &services.NetworkGetVersionInfoQuery{
			Header: q.pbHeader,
		},
	}

	return &services.Query{
		Query: &pb,
	}
}

func (q *NetworkVersionInfoQuery) validateNetworkOnIDs(*Client) error {
	return nil
}

func (q *NetworkVersionInfoQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetNetworkGetVersionInfo()
}
