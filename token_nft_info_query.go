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
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TokenNftInfoQuery
// Applicable only to tokens of type NON_FUNGIBLE_UNIQUE.
// Gets info on a NFT for a given TokenID (of type NON_FUNGIBLE_UNIQUE) and serial number
type TokenNftInfoQuery struct {
	Query
	nftID *NftID
}

// NewTokenNftInfoQuery creates TokenNftInfoQuery which
// gets info on a NFT for a given TokenID (of type NON_FUNGIBLE_UNIQUE) and serial number
// Applicable only to tokens of type NON_FUNGIBLE_UNIQUE.
func NewTokenNftInfoQuery() *TokenNftInfoQuery {
	header := services.QueryHeader{}
	return &TokenNftInfoQuery{
		Query: _NewQuery(true, &header),
		nftID: nil,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (query *TokenNftInfoQuery) SetGrpcDeadline(deadline *time.Duration) *TokenNftInfoQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

// SetNftID Sets the ID of the NFT
func (query *TokenNftInfoQuery) SetNftID(nftID NftID) *TokenNftInfoQuery {
	query.nftID = &nftID
	return query
}

// GetNftID returns the ID of the NFT
func (query *TokenNftInfoQuery) GetNftID() NftID {
	if query.nftID == nil {
		return NftID{}
	}

	return *query.nftID
}

// Deprecated
func (query *TokenNftInfoQuery) SetTokenID(id TokenID) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) GetTokenID() TokenID {
	return TokenID{}
}

// Deprecated
func (query *TokenNftInfoQuery) SetAccountID(id AccountID) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) GetAccountID() AccountID {
	return AccountID{}
}

// Deprecated
func (query *TokenNftInfoQuery) SetStart(start int64) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) GetStart() int64 {
	return 0
}

// Deprecated
func (query *TokenNftInfoQuery) SetEnd(end int64) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) GetEnd() int64 {
	return 0
}

// Deprecated
func (query *TokenNftInfoQuery) ByNftID(id NftID) *TokenNftInfoQuery {
	query.nftID = &id
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) ByTokenID(id TokenID) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) ByAccountID(id AccountID) *TokenNftInfoQuery {
	return query
}

func (query *TokenNftInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.nftID != nil {
		if err := query.nftID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *TokenNftInfoQuery) _BuildByNft() *services.Query_TokenGetNftInfo {
	body := &services.TokenGetNftInfoQuery{
		Header: &services.QueryHeader{},
	}

	if query.nftID != nil {
		body.NftID = query.nftID._ToProtobuf()
	}

	return &services.Query_TokenGetNftInfo{
		TokenGetNftInfo: body,
	}
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (query *TokenNftInfoQuery) GetCost(client *Client) (Hbar, error) {
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

	pb := query._BuildByNft()
	pb.TokenGetNftInfo.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	var resp interface{}
	resp, err = _Execute(
		client,
		&query.Query,
		_TokenNftInfoQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TokenNftInfoQueryGetMethod,
		_TokenNftInfoQueryMapStatusError,
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

	cost := int64(resp.(*services.Response).GetTokenGetNftInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _TokenNftInfoQueryShouldRetry(_ interface{}, response interface{}) _ExecutionState {
	return _QueryShouldRetry(Status(response.(*services.Response).GetTokenGetNftInfo().Header.NodeTransactionPrecheckCode))
}

func _TokenNftInfoQueryMapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetTokenGetNftInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _TokenNftInfoQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetToken().GetTokenNftInfo,
	}
}

// Execute executes the Query with the provided client
func (query *TokenNftInfoQuery) Execute(client *Client) ([]TokenNftInfo, error) {
	if client == nil || client.operator == nil {
		return []TokenNftInfo{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return []TokenNftInfo{}, err
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
			return []TokenNftInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []TokenNftInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TokenNftInfo",
			}
		}

		cost = actualCost
	}

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return []TokenNftInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(query.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return []TokenNftInfo{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._BuildByNft()
	pb.TokenGetNftInfo.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	var resp interface{}
	tokenInfos := make([]TokenNftInfo, 0)
	resp, err = _Execute(
		client,
		&query.Query,
		_TokenNftInfoQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TokenNftInfoQueryGetMethod,
		_TokenNftInfoQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return []TokenNftInfo{}, err
	}

	tokenInfos = append(tokenInfos, _TokenNftInfoFromProtobuf(resp.(*services.Response).GetTokenGetNftInfo().GetNft()))
	return tokenInfos, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TokenNftInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenNftInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *TokenNftInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenNftInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenNftInfoQuery.
func (query *TokenNftInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TokenNftInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (query *TokenNftInfoQuery) SetMaxRetry(count int) *TokenNftInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (query *TokenNftInfoQuery) SetMaxBackoff(max time.Duration) *TokenNftInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (query *TokenNftInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (query *TokenNftInfoQuery) SetMinBackoff(min time.Duration) *TokenNftInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (query *TokenNftInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *TokenNftInfoQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("TokenNftInfoQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (query *TokenNftInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *TokenNftInfoQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}

func (query *TokenNftInfoQuery) SetLogLevel(level LogLevel) *TokenNftInfoQuery {
	query.Query.SetLogLevel(level)
	return query
}
