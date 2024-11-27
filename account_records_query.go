package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// AccountRecordsQuery gets all of the records for an account for any transfers into it and out of
// it, that were above the threshold, during the last 25 hours.
type AccountRecordsQuery struct {
	Query
	accountID *AccountID
}

// NewAccountRecordsQuery creates an AccountRecordsQuery Query which can be used to construct and execute
// an AccountRecordsQuery.
//
// It is recommended that you use this for creating new instances of an AccountRecordQuery
// instead of manually creating an instance of the struct.
func NewAccountRecordsQuery() *AccountRecordsQuery {
	header := services.QueryHeader{}
	return &AccountRecordsQuery{
		Query: _NewQuery(true, &header),
	}
}

// SetGrpcDeadline When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *AccountRecordsQuery) SetGrpcDeadline(deadline *time.Duration) *AccountRecordsQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetAccountID sets the account ID for which the records should be retrieved.
func (q *AccountRecordsQuery) SetAccountID(accountID AccountID) *AccountRecordsQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the account ID for which the records will be retrieved.
func (q *AccountRecordsQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}

	return *q.accountID
}

func (q *AccountRecordsQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *AccountRecordsQuery) Execute(client *Client) ([]TransactionRecord, error) {
	resp, err := q.Query.execute(client, q)
	records := make([]TransactionRecord, 0)

	if err != nil {
		return records, err
	}

	for _, element := range resp.GetCryptoGetAccountRecords().Records {
		record := _TransactionRecordFromProtobuf(&services.TransactionGetRecordResponse{TransactionRecord: element}, nil)
		records = append(records, record)
	}

	return records, err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *AccountRecordsQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountRecordsQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *AccountRecordsQuery) SetQueryPayment(paymentAmount Hbar) *AccountRecordsQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountRecordsQuery.
func (q *AccountRecordsQuery) SetNodeAccountIDs(accountID []AccountID) *AccountRecordsQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *AccountRecordsQuery) SetMaxRetry(count int) *AccountRecordsQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountRecordsQuery) SetMaxBackoff(max time.Duration) *AccountRecordsQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

func (q *AccountRecordsQuery) SetMinBackoff(min time.Duration) *AccountRecordsQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountRecordsQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountRecordsQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *AccountRecordsQuery) SetLogLevel(level LogLevel) *AccountRecordsQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *AccountRecordsQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetAccountRecords,
	}
}

func (q *AccountRecordsQuery) getName() string {
	return "AccountRecordsQuery"
}

func (q *AccountRecordsQuery) buildQuery() *services.Query {
	pb := services.Query_CryptoGetAccountRecords{
		CryptoGetAccountRecords: &services.CryptoGetAccountRecordsQuery{
			Header: q.pbHeader,
		},
	}

	if q.accountID != nil {
		pb.CryptoGetAccountRecords.AccountID = q.accountID._ToProtobuf()
	}

	return &services.Query{
		Query: &pb,
	}
}

func (q *AccountRecordsQuery) validateNetworkOnIDs(client *Client) error {
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

func (q *AccountRecordsQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetCryptoGetAccountRecords()
}
