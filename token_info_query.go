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

	//	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *TokenInfoQuery) SetGrpcDeadline(deadline *time.Duration) *TokenInfoQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetTokenID Sets the topic to retrieve info about (the parameters and running state of).
func (this *TokenInfoQuery) SetTokenID(tokenID TokenID) *TokenInfoQuery {
	this.tokenID = &tokenID
	return this
}

// GetTokenID returns the TokenID for this TokenInfoQuery
func (this *TokenInfoQuery) GetTokenID() TokenID {
	if this.tokenID == nil {
		return TokenID{}
	}

	return *this.tokenID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *TokenInfoQuery) GetCost(client *Client) (Hbar, error) {
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
	pb.TokenGetInfo.Header = this.pbHeader

	this.pb = &services.Query{
		Query: pb,
	}

	this.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	this.paymentTransactionIDs._Advance()

	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetTokenGetInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the TopicInfoQuery using the provided client
func (this *TokenInfoQuery) Execute(client *Client) (TokenInfo, error) {
	if client == nil || client.operator == nil {
		return TokenInfo{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return TokenInfo{}, err
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
			return TokenInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TokenInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TokenInfoQuery",
			}
		}

		cost = actualCost
	}

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return TokenInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return TokenInfo{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.TokenGetInfo.Header = this.pbHeader
	this.pb = &services.Query{
		Query: pb,
	}

	if this.isPaymentRequired && len(this.paymentTransactions) > 0 {
		this.paymentTransactionIDs._Advance()
	}
	this.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return TokenInfo{}, err
	}

	info := _TokenInfoFromProtobuf(resp.(*services.Response).GetTokenGetInfo().TokenInfo)

	return info, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *TokenInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenInfoQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *TokenInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenInfoQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenInfoQuery.
func (this *TokenInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TokenInfoQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *TokenInfoQuery) SetMaxRetry(count int) *TokenInfoQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *TokenInfoQuery) SetMaxBackoff(max time.Duration) *TokenInfoQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *TokenInfoQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *TokenInfoQuery) SetMinBackoff(min time.Duration) *TokenInfoQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *TokenInfoQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *TokenInfoQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("TokenInfoQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *TokenInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *TokenInfoQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *TokenInfoQuery) SetLogLevel(level LogLevel) *TokenInfoQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *TokenInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetToken().GetTokenInfo,
	}
}

func (this *TokenInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetTokenGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (this *TokenInfoQuery) getName() string {
	return "TokenInfoQuery"
}

func (this *TokenInfoQuery) build() *services.Query_TokenGetInfo {
	body := &services.TokenGetInfoQuery{
		Header: &services.QueryHeader{},
	}
	if this.tokenID != nil {
		body.Token = this.tokenID._ToProtobuf()
	}

	return &services.Query_TokenGetInfo{
		TokenGetInfo: body,
	}
}

func (this *TokenInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.tokenID != nil {
		if err := this.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *TokenInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetTokenGetInfo().Header.NodeTransactionPrecheckCode)
}
