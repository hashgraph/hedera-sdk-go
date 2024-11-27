package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// AccountInfoQuery
// Get all the information about an account, including the balance. This does not get the list of
// account records.
type AccountInfoQuery struct {
	Query
	accountID *AccountID
}

// NewAccountInfoQuery
// Creates an AccountInfoQuery which retrieves all the information about an account, including the balance. This does not get the list of
// account records.
func NewAccountInfoQuery() *AccountInfoQuery {
	header := services.QueryHeader{}
	return &AccountInfoQuery{
		Query: _NewQuery(true, &header),
	}
}

func (q *AccountInfoQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *AccountInfoQuery) Execute(client *Client) (AccountInfo, error) {
	resp, err := q.execute(client, q)

	if err != nil {
		return AccountInfo{}, err
	}

	return _AccountInfoFromProtobuf(resp.GetCryptoGetInfo().AccountInfo)
}

// SetGrpcDeadline When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *AccountInfoQuery) SetGrpcDeadline(deadline *time.Duration) *AccountInfoQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetAccountID sets the AccountID for this AccountInfoQuery.
func (q *AccountInfoQuery) SetAccountID(accountID AccountID) *AccountInfoQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the AccountID for this AccountInfoQuery.
func (q *AccountInfoQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}

	return *q.accountID
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountInfoQuery.
func (q *AccountInfoQuery) SetNodeAccountIDs(accountID []AccountID) *AccountInfoQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetQueryPayment sets the Hbar payment to pay the _Node a fee for handling this query
func (q *AccountInfoQuery) SetQueryPayment(queryPayment Hbar) *AccountInfoQuery {
	q.queryPayment = queryPayment
	return q
}

// SetMaxQueryPayment sets the maximum payment allowable for this query.
func (q *AccountInfoQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *AccountInfoQuery {
	q.maxQueryPayment = queryMaxPayment
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *AccountInfoQuery) SetMaxRetry(count int) *AccountInfoQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountInfoQuery) SetMaxBackoff(max time.Duration) *AccountInfoQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *AccountInfoQuery) SetMinBackoff(min time.Duration) *AccountInfoQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountInfoQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *AccountInfoQuery) SetLogLevel(level LogLevel) *AccountInfoQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *AccountInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetAccountInfo,
	}
}

func (q *AccountInfoQuery) getName() string {
	return "AccountInfoQuery"
}

func (q *AccountInfoQuery) buildQuery() *services.Query {
	pbQuery := services.Query_CryptoGetInfo{
		CryptoGetInfo: &services.CryptoGetInfoQuery{
			Header: q.pbHeader,
		},
	}

	if q.accountID != nil {
		pbQuery.CryptoGetInfo.AccountID = q.accountID._ToProtobuf()
	}

	return &services.Query{
		Query: &pbQuery,
	}
}

func (q *AccountInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.accountID != nil {
		if err := q.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *AccountInfoQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetCryptoGetInfo()
}
