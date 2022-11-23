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

// LiveHashQuery Requests a livehash associated to an account.
type LiveHashQuery struct {
	Query
	accountID *AccountID
	hash      []byte
}

// NewLiveHashQuery creates a LiveHashQuery that requests a livehash associated to an account.
func NewLiveHashQuery() *LiveHashQuery {
	header := services.QueryHeader{}
	return &LiveHashQuery{
		Query: _NewQuery(true, &header),
	}
}

func (query *LiveHashQuery) SetGrpcDeadline(deadline *time.Duration) *LiveHashQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

// SetAccountID Sets the AccountID to which the livehash is associated
func (query *LiveHashQuery) SetAccountID(accountID AccountID) *LiveHashQuery {
	query.accountID = &accountID
	return query
}

func (query *LiveHashQuery) GetAccountID() AccountID {
	if query.accountID == nil {
		return AccountID{}
	}

	return *query.accountID
}

// SetHash Sets the SHA-384 data in the livehash
func (query *LiveHashQuery) SetHash(hash []byte) *LiveHashQuery {
	query.hash = hash
	return query
}

func (query *LiveHashQuery) GetGetHash() []byte {
	return query.hash
}

func (query *LiveHashQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.accountID != nil {
		if err := query.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *LiveHashQuery) _Build() *services.Query_CryptoGetLiveHash {
	body := &services.CryptoGetLiveHashQuery{
		Header: &services.QueryHeader{},
	}
	if query.accountID != nil {
		body.AccountID = query.accountID._ToProtobuf()
	}

	if len(query.hash) > 0 {
		body.Hash = query.hash
	}

	return &services.Query_CryptoGetLiveHash{
		CryptoGetLiveHash: body,
	}
}

func (query *LiveHashQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range query.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.CryptoGetLiveHash.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_LiveHashQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_LiveHashQueryGetMethod,
		_LiveHashQueryMapStatusError,
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

	cost := int64(resp.(*services.Response).GetCryptoGetLiveHash().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _LiveHashQueryShouldRetry(logID string, _ interface{}, response interface{}) _ExecutionState {
	return _QueryShouldRetry(logID, Status(response.(*services.Response).GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode))
}

func _LiveHashQueryMapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode),
	}
}

func _LiveHashQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetLiveHash,
	}
}

func (query *LiveHashQuery) Execute(client *Client) (LiveHash, error) {
	if client == nil || client.operator == nil {
		return LiveHash{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return LiveHash{}, err
	}

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
			return LiveHash{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return LiveHash{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "LiveHashQuery",
			}
		}

		cost = actualCost
	}

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return LiveHash{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(query.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return LiveHash{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.CryptoGetLiveHash.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_LiveHashQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_LiveHashQueryGetMethod,
		_LiveHashQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return LiveHash{}, err
	}

	liveHash, err := _LiveHashFromProtobuf(resp.(*services.Response).GetCryptoGetLiveHash().LiveHash)
	if err != nil {
		return LiveHash{}, err
	}

	return liveHash, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *LiveHashQuery) SetMaxQueryPayment(maxPayment Hbar) *LiveHashQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *LiveHashQuery) SetQueryPayment(paymentAmount Hbar) *LiveHashQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *LiveHashQuery) SetNodeAccountIDs(accountID []AccountID) *LiveHashQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *LiveHashQuery) SetMaxBackoff(max time.Duration) *LiveHashQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *LiveHashQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *LiveHashQuery) SetMinBackoff(min time.Duration) *LiveHashQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *LiveHashQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *LiveHashQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("LiveHashQuery:%d", timestamp)
}

func (query *LiveHashQuery) SetPaymentTransactionID(transactionID TransactionID) *LiveHashQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}

func (query *LiveHashQuery) SetMaxRetry(count int) *LiveHashQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *LiveHashQuery) GetMaxRetry() int {
	return query.Query.maxRetry
}
