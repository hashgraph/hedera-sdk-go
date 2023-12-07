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

// LiveHashQuery Requests a livehash associated to an account.
type LiveHashQuery struct {
	Query
	accountID *AccountID
	hash      []byte
}

// NewLiveHashQuery creates a LiveHashQuery that requests a livehash associated to an account.
func NewLiveHashQuery() *LiveHashQuery {
	header := services.QueryHeader{}
	result := LiveHashQuery{
		Query: _NewQuery(true, &header),
	}
	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *LiveHashQuery) SetGrpcDeadline(deadline *time.Duration) *LiveHashQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetAccountID Sets the AccountID to which the livehash is associated
func (q *LiveHashQuery) SetAccountID(accountID AccountID) *LiveHashQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the AccountID to which the livehash is associated
func (q *LiveHashQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}

	return *q.accountID
}

// SetHash Sets the SHA-384 data in the livehash
func (q *LiveHashQuery) SetHash(hash []byte) *LiveHashQuery {
	q.hash = hash
	return q
}

// GetHash returns the SHA-384 data in the livehash
func (q *LiveHashQuery) GetGetHash() []byte {
	return q.hash
}

// Execute executes the QueryInterface with the provided client
func (q *LiveHashQuery) Execute(client *Client) (LiveHash, error) {
	resp, err := q.Query.execute(client)

	if err != nil {
		return LiveHash{}, err
	}

	liveHash, err := _LiveHashFromProtobuf(resp.GetCryptoGetLiveHash().LiveHash)
	if err != nil {
		return LiveHash{}, err
	}

	return liveHash, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this QueryInterface.
func (q *LiveHashQuery) SetMaxQueryPayment(maxPayment Hbar) *LiveHashQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this QueryInterface.
func (q *LiveHashQuery) SetQueryPayment(paymentAmount Hbar) *LiveHashQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this LiveHashQuery.
func (q *LiveHashQuery) SetNodeAccountIDs(accountID []AccountID) *LiveHashQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *LiveHashQuery) SetMaxBackoff(max time.Duration) *LiveHashQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *LiveHashQuery) SetMinBackoff(min time.Duration) *LiveHashQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *LiveHashQuery) SetPaymentTransactionID(transactionID TransactionID) *LiveHashQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *LiveHashQuery) SetMaxRetry(count int) *LiveHashQuery {
	q.Query.SetMaxRetry(count)
	return q
}

func (q *LiveHashQuery) SetLogLevel(level LogLevel) *LiveHashQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *LiveHashQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetLiveHash,
	}
}

func (q *LiveHashQuery) getName() string {
	return "LiveHashQuery"
}

func (q *LiveHashQuery) buildQuery() *services.Query {
	body := &services.CryptoGetLiveHashQuery{
		Header: q.pbHeader,
	}
	if q.accountID != nil {
		body.AccountID = q.accountID._ToProtobuf()
	}

	if len(q.hash) > 0 {
		body.Hash = q.hash
	}

	return &services.Query{
		Query: &services.Query_CryptoGetLiveHash{
			CryptoGetLiveHash: body,
		},
	}
}

func (q *LiveHashQuery) validateNetworkOnIDs(client *Client) error {
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

func (q *LiveHashQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetCryptoGetLiveHash()
}
