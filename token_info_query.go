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

// TokenInfoQuery Used get information about Token instance
type TokenInfoQuery struct {
	query
	tokenID *TokenID
}

// NewTokenInfoQuery creates a TokenInfoQuery which is used get information about Token instance
func NewTokenInfoQuery() *TokenInfoQuery {
	header := services.QueryHeader{}
	result := TokenInfoQuery{
		query: _NewQuery(true, &header),
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *TokenInfoQuery) SetGrpcDeadline(deadline *time.Duration) *TokenInfoQuery {
	q.query.SetGrpcDeadline(deadline)
	return q
}

// SetTokenID Sets the topic to retrieve info about (the parameters and running state of).
func (q *TokenInfoQuery) SetTokenID(tokenID TokenID) *TokenInfoQuery {
	q.tokenID = &tokenID
	return q
}

// GetTokenID returns the TokenID for this TokenInfoQuery
func (q *TokenInfoQuery) GetTokenID() TokenID {
	if q.tokenID == nil {
		return TokenID{}
	}

	return *q.tokenID
}

// Execute executes the TopicInfoQuery using the provided client
func (q *TokenInfoQuery) Execute(client *Client) (TokenInfo, error) {
	resp, err := q.query.execute(client)

	if err != nil {
		return TokenInfo{}, err
	}

	info := _TokenInfoFromProtobuf(resp.GetTokenGetInfo().TokenInfo)

	return info, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *TokenInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenInfoQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *TokenInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenInfoQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenInfoQuery.
func (q *TokenInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TokenInfoQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *TokenInfoQuery) SetMaxRetry(count int) *TokenInfoQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *TokenInfoQuery) SetMaxBackoff(max time.Duration) *TokenInfoQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *TokenInfoQuery) SetMinBackoff(min time.Duration) *TokenInfoQuery {
	q.query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *TokenInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *TokenInfoQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *TokenInfoQuery) SetLogLevel(level LogLevel) *TokenInfoQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *TokenInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetToken().GetTokenInfo,
	}
}

func (q *TokenInfoQuery) getName() string {
	return "TokenInfoQuery"
}

func (q *TokenInfoQuery) buildQuery() *services.Query {
	body := &services.TokenGetInfoQuery{
		Header: q.pbHeader,
	}
	if q.tokenID != nil {
		body.Token = q.tokenID._ToProtobuf()
	}

	return &services.Query{
		Query: &services.Query_TokenGetInfo{
			TokenGetInfo: body,
		},
	}
}

func (q *TokenInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.tokenID != nil {
		if err := q.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *TokenInfoQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetTokenGetInfo()
}
