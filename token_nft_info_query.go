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
	query
	nftID *NftID
}

// NewTokenNftInfoQuery creates TokenNftInfoQuery which
// gets info on a NFT for a given TokenID (of type NON_FUNGIBLE_UNIQUE) and serial number
// Applicable only to tokens of type NON_FUNGIBLE_UNIQUE.
func NewTokenNftInfoQuery() *TokenNftInfoQuery {
	header := services.QueryHeader{}
	result := TokenNftInfoQuery{
		query: _NewQuery(true, &header),
		nftID: nil,
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *TokenNftInfoQuery) SetGrpcDeadline(deadline *time.Duration) *TokenNftInfoQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetNftID Sets the ID of the NFT
func (this *TokenNftInfoQuery) SetNftID(nftID NftID) *TokenNftInfoQuery {
	this.nftID = &nftID
	return this
}

// GetNftID returns the ID of the NFT
func (this *TokenNftInfoQuery) GetNftID() NftID {
	if this.nftID == nil {
		return NftID{}
	}

	return *this.nftID
}

// Deprecated
func (this *TokenNftInfoQuery) SetTokenID(id TokenID) *TokenNftInfoQuery {
	return this
}

// Deprecated
func (this *TokenNftInfoQuery) GetTokenID() TokenID {
	return TokenID{}
}

// Deprecated
func (this *TokenNftInfoQuery) SetAccountID(id AccountID) *TokenNftInfoQuery {
	return this
}

// Deprecated
func (this *TokenNftInfoQuery) GetAccountID() AccountID {
	return AccountID{}
}

// Deprecated
func (this *TokenNftInfoQuery) SetStart(start int64) *TokenNftInfoQuery {
	return this
}

// Deprecated
func (this *TokenNftInfoQuery) GetStart() int64 {
	return 0
}

// Deprecated
func (this *TokenNftInfoQuery) SetEnd(end int64) *TokenNftInfoQuery {
	return this
}

// Deprecated
func (this *TokenNftInfoQuery) GetEnd() int64 {
	return 0
}

// Deprecated
func (this *TokenNftInfoQuery) ByNftID(id NftID) *TokenNftInfoQuery {
	this.nftID = &id
	return this
}

// Deprecated
func (this *TokenNftInfoQuery) ByTokenID(id TokenID) *TokenNftInfoQuery {
	return this
}

// Deprecated
func (this *TokenNftInfoQuery) ByAccountID(id AccountID) *TokenNftInfoQuery {
	return this
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *TokenNftInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range this.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.TokenGetNftInfo.Header = this.pbHeader

	this.pb = &services.Query{
		Query: pb,
	}

	this.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	this.paymentTransactionIDs._Advance()

	var resp interface{}
	resp, err = _Execute(
		client,
		this.e,
	)
	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetTokenGetNftInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (this *TokenNftInfoQuery) Execute(client *Client) ([]TokenNftInfo, error) {
	if client == nil || client.operator == nil {
		return []TokenNftInfo{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return []TokenNftInfo{}, err
	}

	if !this.paymentTransactionIDs.locked {
		this.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if this.queryPayment.tinybar != 0 {
		cost = this.queryPayment
	} else {
		if this.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = this.maxQueryPayment
		}

		actualCost, err := this.GetCost(client)
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

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return []TokenNftInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return []TokenNftInfo{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.TokenGetNftInfo.Header = this.pbHeader
	this.pb = &services.Query{
		Query: pb,
	}

	if this.isPaymentRequired && len(this.paymentTransactions) > 0 {
		this.paymentTransactionIDs._Advance()
	}
	this.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	var resp interface{}
	tokenInfos := make([]TokenNftInfo, 0)
	resp, err = _Execute(
		client,
		this.e,
	)

	if err != nil {
		return []TokenNftInfo{}, err
	}

	tokenInfos = append(tokenInfos, _TokenNftInfoFromProtobuf(resp.(*services.Response).GetTokenGetNftInfo().GetNft()))
	return tokenInfos, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *TokenNftInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenNftInfoQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *TokenNftInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenNftInfoQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenNftInfoQuery.
func (this *TokenNftInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TokenNftInfoQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *TokenNftInfoQuery) SetMaxRetry(count int) *TokenNftInfoQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *TokenNftInfoQuery) SetMaxBackoff(max time.Duration) *TokenNftInfoQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *TokenNftInfoQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *TokenNftInfoQuery) SetMinBackoff(min time.Duration) *TokenNftInfoQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *TokenNftInfoQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *TokenNftInfoQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("TokenNftInfoQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *TokenNftInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *TokenNftInfoQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *TokenNftInfoQuery) SetLogLevel(level LogLevel) *TokenNftInfoQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *TokenNftInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetToken().GetTokenNftInfo,
	}
}

func (this *TokenNftInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetTokenGetNftInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (this *TokenNftInfoQuery) getName() string {
	return "TokenNftInfoQuery"
}

func (this *TokenNftInfoQuery) build()*services.Query_TokenGetNftInfo {
	body := &services.TokenGetNftInfoQuery{
		Header: &services.QueryHeader{},
	}

	if this.nftID != nil {
		body.NftID = this.nftID._ToProtobuf()
	}

	return &services.Query_TokenGetNftInfo{
		TokenGetNftInfo: body,
	}
}

func (this *TokenNftInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.nftID != nil {
		if err := this.nftID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *TokenNftInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetTokenGetNftInfo().Header.NodeTransactionPrecheckCode)
}
