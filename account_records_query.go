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

// AccountRecordsQuery gets all of the records for an account for any transfers into it and out of
// it, that were above the threshold, during the last 25 hours.
type AccountRecordsQuery struct {
	query
	accountID *AccountID
}

// NewAccountRecordsQuery creates an AccountRecordsQuery query which can be used to construct and execute
// an AccountRecordsQuery.
//
// It is recommended that you use this for creating new instances of an AccountRecordQuery
// instead of manually creating an instance of the struct.
func NewAccountRecordsQuery() *AccountRecordsQuery {
	header := services.QueryHeader{}
	result := AccountRecordsQuery{
		query: _NewQuery(true, &header),
	}
	result.e = &result

	return &result

}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *AccountRecordsQuery) SetGrpcDeadline(deadline *time.Duration) *AccountRecordsQuery {
	q.query.SetGrpcDeadline(deadline)
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

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (q *AccountRecordsQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range q.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.CryptoGetAccountRecords.Header = q.pbHeader

	q.pb = &services.Query{
		Query: pb,
	}

	q.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	q.paymentTransactionIDs._Advance()

	resp, err := _Execute(
		client,
		q.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetCryptoGetAccountRecords().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (q *AccountRecordsQuery) Execute(client *Client) ([]TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return []TransactionRecord{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return []TransactionRecord{}, err
	}

	if !q.paymentTransactionIDs.locked {
		q.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if q.queryPayment.tinybar != 0 {
		cost = q.queryPayment
	} else {
		if q.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = q.maxQueryPayment
		}

		actualCost, err := q.GetCost(client)
		if err != nil {
			return []TransactionRecord{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []TransactionRecord{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "AccountRecordsQuery",
			}
		}

		cost = actualCost
	}

	q.paymentTransactions = make([]*services.Transaction, 0)

	if q.nodeAccountIDs.locked {
		err = q._QueryGeneratePayments(client, cost)
		if err != nil {
			return []TransactionRecord{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(q.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			if err != nil {
				return []TransactionRecord{}, err
			}
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.CryptoGetAccountRecords.Header = q.pbHeader
	q.pb = &services.Query{
		Query: pb,
	}

	records := make([]TransactionRecord, 0)

	if q.isPaymentRequired && len(q.paymentTransactions) > 0 {
		q.paymentTransactionIDs._Advance()
	}
	q.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(
		client,
		q.e,
	)

	if err != nil {
		return []TransactionRecord{}, err
	}

	for _, element := range resp.(*services.Response).GetCryptoGetAccountRecords().Records {
		record := _TransactionRecordFromProtobuf(&services.TransactionGetRecordResponse{TransactionRecord: element}, nil)
		records = append(records, record)
	}

	return records, err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *AccountRecordsQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountRecordsQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *AccountRecordsQuery) SetQueryPayment(paymentAmount Hbar) *AccountRecordsQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountRecordsQuery.
func (q *AccountRecordsQuery) SetNodeAccountIDs(accountID []AccountID) *AccountRecordsQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *AccountRecordsQuery) SetMaxRetry(count int) *AccountRecordsQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountRecordsQuery) SetMaxBackoff(max time.Duration) *AccountRecordsQuery {
	q.query.SetMaxBackoff(max)
	return q
}

func (q *AccountRecordsQuery) SetMinBackoff(min time.Duration) *AccountRecordsQuery {
	q.query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountRecordsQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountRecordsQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *AccountRecordsQuery) SetLogLevel(level LogLevel) *AccountRecordsQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *AccountRecordsQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetAccountRecords,
	}
}

func (q *AccountRecordsQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetCryptoGetAccountRecords().Header.NodeTransactionPrecheckCode),
	}
}

func (q *AccountRecordsQuery) getName() string {
	return "AccountRecordsQuery"
}

func (q *AccountRecordsQuery) build() *services.Query_CryptoGetAccountRecords {
	pb := services.Query_CryptoGetAccountRecords{
		CryptoGetAccountRecords: &services.CryptoGetAccountRecordsQuery{
			Header: &services.QueryHeader{},
		},
	}

	if q.accountID != nil {
		pb.CryptoGetAccountRecords.AccountID = q.accountID._ToProtobuf()
	}

	return &pb
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

func (q *AccountRecordsQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetCryptoGetAccountRecords().Header.NodeTransactionPrecheckCode)
}
