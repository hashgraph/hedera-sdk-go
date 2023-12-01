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

// AccountInfoQuery
// Get all the information about an account, including the balance. This does not get the list of
// account records.
type AccountInfoQuery struct {
	query
	accountID *AccountID
}

// NewAccountInfoQuery
// Creates an AccountInfoQuery which retrieves all the information about an account, including the balance. This does not get the list of
// account records.
func NewAccountInfoQuery() *AccountInfoQuery {
	header := services.QueryHeader{}
	result := AccountInfoQuery{
		query: _NewQuery(true, &header),
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *AccountInfoQuery) SetGrpcDeadline(deadline *time.Duration) *AccountInfoQuery {
	q.query.SetGrpcDeadline(deadline)
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

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (q *AccountInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	if q.nodeAccountIDs.locked {
		for range q.nodeAccountIDs.slice {
			paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
			if err != nil {
				return Hbar{}, err
			}
			q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.CryptoGetInfo.Header = q.pbHeader

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

	cost := int64(resp.(*services.Response).GetCryptoGetInfo().Header.Cost)

	return HbarFromTinybar(cost), nil
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountInfoQuery.
func (q *AccountInfoQuery) SetNodeAccountIDs(accountID []AccountID) *AccountInfoQuery {
	q.query.SetNodeAccountIDs(accountID)
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
	q.query.SetMaxRetry(count)
	return q
}

// Execute executes the Query with the provided client
func (q *AccountInfoQuery) Execute(client *Client) (AccountInfo, error) {
	if client == nil || client.operator == nil {
		return AccountInfo{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return AccountInfo{}, err
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
			return AccountInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return AccountInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "AccountInfoQuery",
			}
		}

		cost = actualCost
	}

	q.paymentTransactions = make([]*services.Transaction, 0)
	if q.nodeAccountIDs.locked {
		err = q._QueryGeneratePayments(client, cost)
		if err != nil {
			return AccountInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(q.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return AccountInfo{}, err
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.CryptoGetInfo.Header = q.pbHeader
	q.pb = &services.Query{
		Query: pb,
	}

	if q.isPaymentRequired && len(q.paymentTransactions) > 0 {
		q.paymentTransactionIDs._Advance()
	}
	q.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY
	resp, err := _Execute(
		client,
		q.e,
	)

	if err != nil {
		return AccountInfo{}, err
	}

	return _AccountInfoFromProtobuf(resp.(*services.Response).GetCryptoGetInfo().AccountInfo)
}

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountInfoQuery) SetMaxBackoff(max time.Duration) *AccountInfoQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *AccountInfoQuery) SetMinBackoff(min time.Duration) *AccountInfoQuery {
	q.query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountInfoQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *AccountInfoQuery) SetLogLevel(level LogLevel) *AccountInfoQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *AccountInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetAccountInfo,
	}
}

func (q *AccountInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetCryptoGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (q *AccountInfoQuery) getName() string {
	return "AccountInfoQuery"
}

func (q *AccountInfoQuery) build() *services.Query_CryptoGetInfo {
	pb := services.Query_CryptoGetInfo{
		CryptoGetInfo: &services.CryptoGetInfoQuery{
			Header: &services.QueryHeader{},
		},
	}

	if q.accountID != nil {
		pb.CryptoGetInfo.AccountID = q.accountID._ToProtobuf()
	}

	return &pb
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

func (q *AccountInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetCryptoGetInfo().Header.NodeTransactionPrecheckCode)
}
