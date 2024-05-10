package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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

	info, err := _AccountInfoFromProtobuf(resp.GetCryptoGetInfo().AccountInfo)
	if err != nil {
		return AccountInfo{}, err
	}

	err = fetchAccountInfoTokenRelationships(obtainUrlForMirrorNode(client), q.accountID.String(), &info)
	if err != nil {
		return info, err
	}

	return info, nil
}

/*
Helper function, which queries the mirror node and if the query result has token relations, it iterates over the token
relationships and populates the appropriate field in AccountInfo object

IMPORTANT: This function will fetch the state of the data in the Mirror Node at the moment of its execution. It
is important to note that the Mirror Node currently needs 2-3 seconds to be updated with the latest data from the
consensus nodes. So if data related to token relationships is changed and a proper timeout is not introduced the
user would not get the up to date state of token relationships. This note is ONLY for token relationship data as it
is queried from the MirrorNode. Other query information arrives at the time of consensus response.
*/
func fetchAccountInfoTokenRelationships(network string, id string, info *AccountInfo) error {
	response, err := tokenReleationshipMirrorNodeQuery(network, id)
	if err != nil {
		return err
	}

	if tokens, ok := response["tokens"].([]interface{}); ok {
		for _, token := range tokens {
			fmt.Println(token)
			tr, err := TokenRelationshipFromJson(token)
			if err != nil {
				return err
			}
			info.TokenRelationships = append(info.TokenRelationships, tr)
		}
	}

	return nil
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
