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

// AccountStakersQuery gets all of the accounts that are proxy staking to this account. For each of  them, the amount
// currently staked will be given. This is not yet implemented, but will be in a future version of the API.
type AccountStakersQuery struct {
	query
	accountID *AccountID
}

// NewAccountStakersQuery creates an AccountStakersQuery query which can be used to construct and execute
// an AccountStakersQuery.
//
// It is recommended that you use this for creating new instances of an AccountStakersQuery
// instead of manually creating an instance of the struct.
func NewAccountStakersQuery() *AccountStakersQuery {
	header := services.QueryHeader{}
	result := AccountStakersQuery{
		query: _NewQuery(true, &header),
	}

	//	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *AccountStakersQuery) SetGrpcDeadline(deadline *time.Duration) *AccountStakersQuery {
	q.query.SetGrpcDeadline(deadline)
	return q
}

// SetAccountID sets the Account ID for which the stakers should be retrieved
func (q *AccountStakersQuery) SetAccountID(accountID AccountID) *AccountStakersQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the AccountID for this AccountStakersQuery.
func (q *AccountStakersQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}

	return *q.accountID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (q *AccountStakersQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	err := q.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}
	if q.query.nodeAccountIDs.locked {
		for range q.nodeAccountIDs.slice {
			paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
			if err != nil {
				return Hbar{}, err
			}
			q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
		}
	}

	pb := q.build()
	pb.CryptoGetProxyStakers.Header = q.pbHeader

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

	cost := int64(resp.(*services.Response).GetCryptoGetProxyStakers().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func (q *AccountStakersQuery) Execute(client *Client) ([]Transfer, error) {
	if client == nil || client.operator == nil {
		return []Transfer{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return []Transfer{}, err
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
			return []Transfer{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []Transfer{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "AccountStakersQuery",
			}
		}

		cost = actualCost
	}

	q.paymentTransactions = make([]*services.Transaction, 0)

	if q.nodeAccountIDs.locked {
		err = q._QueryGeneratePayments(client, cost)
		if err != nil {
			return []Transfer{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(q.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return []Transfer{}, err
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.CryptoGetProxyStakers.Header = q.pbHeader
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
		return []Transfer{}, err
	}

	var stakers = make([]Transfer, len(resp.(*services.Response).GetCryptoGetProxyStakers().Stakers.ProxyStaker))

	// TODO: This is wrong, q _Method shold return `[]ProxyStaker` not `[]Transfer`
	for i, element := range resp.(*services.Response).GetCryptoGetProxyStakers().Stakers.ProxyStaker {
		id := _AccountIDFromProtobuf(element.AccountID)
		accountID := AccountID{}

		if id == nil {
			accountID = *id
		}

		stakers[i] = Transfer{
			AccountID: accountID,
			Amount:    HbarFromTinybar(element.Amount),
		}
	}

	return stakers, err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *AccountStakersQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountStakersQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *AccountStakersQuery) SetQueryPayment(paymentAmount Hbar) *AccountStakersQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountStakersQuery.
func (q *AccountStakersQuery) SetNodeAccountIDs(accountID []AccountID) *AccountStakersQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *AccountStakersQuery) SetMaxRetry(count int) *AccountStakersQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountStakersQuery) SetMaxBackoff(max time.Duration) *AccountStakersQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *AccountStakersQuery) SetMinBackoff(min time.Duration) *AccountStakersQuery {
	q.query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountStakersQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountStakersQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *AccountStakersQuery) SetLogLevel(level LogLevel) *AccountStakersQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *AccountStakersQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetStakersByAccountID,
	}
}

func (q *AccountStakersQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetCryptoGetProxyStakers().Header.NodeTransactionPrecheckCode),
	}
}

func (q *AccountStakersQuery) getName() string {
	return "AccountStakersQuery"
}

func (q *AccountStakersQuery) build() *services.Query_CryptoGetProxyStakers {
	pb := services.Query_CryptoGetProxyStakers{
		CryptoGetProxyStakers: &services.CryptoGetStakersQuery{
			Header: &services.QueryHeader{},
		},
	}

	if q.accountID != nil {
		pb.CryptoGetProxyStakers.AccountID = q.accountID._ToProtobuf()
	}

	return &pb
}

func (q *AccountStakersQuery) validateNetworkOnIDs(client *Client) error {
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

func (q *AccountStakersQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetCryptoGetProxyStakers().Header.NodeTransactionPrecheckCode)
}
