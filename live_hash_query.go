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

// LiveHashQuery Requests a livehash associated to an account.
type LiveHashQuery struct {
	query
	accountID *AccountID
	hash      []byte
}

// NewLiveHashQuery creates a LiveHashQuery that requests a livehash associated to an account.
func NewLiveHashQuery() *LiveHashQuery {
	header := services.QueryHeader{}
	result := LiveHashQuery{
		query: _NewQuery(true, &header),
	}
	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *LiveHashQuery) SetGrpcDeadline(deadline *time.Duration) *LiveHashQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetAccountID Sets the AccountID to which the livehash is associated
func (this *LiveHashQuery) SetAccountID(accountID AccountID) *LiveHashQuery {
	this.accountID = &accountID
	return this
}

// GetAccountID returns the AccountID to which the livehash is associated
func (this *LiveHashQuery) GetAccountID() AccountID {
	if this.accountID == nil {
		return AccountID{}
	}

	return *this.accountID
}

// SetHash Sets the SHA-384 data in the livehash
func (this *LiveHashQuery) SetHash(hash []byte) *LiveHashQuery {
	this.hash = hash
	return this
}

// GetHash returns the SHA-384 data in the livehash
func (this *LiveHashQuery) GetGetHash() []byte {
	return this.hash
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *LiveHashQuery) GetCost(client *Client) (Hbar, error) {
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
	pb.CryptoGetLiveHash.Header = this.pbHeader

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

	cost := int64(resp.(*services.Response).GetCryptoGetLiveHash().Header.Cost)
	return HbarFromTinybar(cost), nil
}
// Execute executes the Query with the provided client
func (this *LiveHashQuery) Execute(client *Client) (LiveHash, error) {
	if client == nil || client.operator == nil {
		return LiveHash{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return LiveHash{}, err
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

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return LiveHash{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return LiveHash{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.CryptoGetLiveHash.Header = this.pbHeader
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
		return LiveHash{}, err
	}

	liveHash, err := _LiveHashFromProtobuf(resp.(*services.Response).GetCryptoGetLiveHash().LiveHash)
	if err != nil {
		return LiveHash{}, err
	}

	return liveHash, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *LiveHashQuery) SetMaxQueryPayment(maxPayment Hbar) *LiveHashQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *LiveHashQuery) SetQueryPayment(paymentAmount Hbar) *LiveHashQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this LiveHashQuery.
func (this *LiveHashQuery) SetNodeAccountIDs(accountID []AccountID) *LiveHashQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *LiveHashQuery) SetMaxBackoff(max time.Duration) *LiveHashQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *LiveHashQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *LiveHashQuery) SetMinBackoff(min time.Duration) *LiveHashQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *LiveHashQuery) GetMinBackoff() time.Duration {
	return *this.query.GetGrpcDeadline()
}

func (this *LiveHashQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("LiveHashQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *LiveHashQuery) SetPaymentTransactionID(transactionID TransactionID) *LiveHashQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *LiveHashQuery) SetMaxRetry(count int) *LiveHashQuery {
	this.query.SetMaxRetry(count)
	return this
}

// GetMaxRetry returns the max number of errors before execution will fail.
func (this *LiveHashQuery) GetMaxRetry() int {
	return this.query.maxRetry
}

func (this *LiveHashQuery) SetLogLevel(level LogLevel) *LiveHashQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *LiveHashQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetLiveHash,
	}
}

func (this *LiveHashQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode),
	}
}

func (this *LiveHashQuery) getName() string {
	return "LiveHashQuery"
}

func (this *LiveHashQuery) build() *services.Query_CryptoGetLiveHash {
	body := &services.CryptoGetLiveHashQuery{
		Header: &services.QueryHeader{},
	}
	if this.accountID != nil {
		body.AccountID = this.accountID._ToProtobuf()
	}

	if len(this.hash) > 0 {
		body.Hash = this.hash
	}

	return &services.Query_CryptoGetLiveHash{
		CryptoGetLiveHash: body,
	}
}

func (this *LiveHashQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.accountID != nil {
		if err := this.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *LiveHashQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode)
}