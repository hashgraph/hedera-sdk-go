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

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
type AccountBalanceQuery struct {
	query
	accountID  *AccountID
	contractID *ContractID
	timestamp  time.Time
}

// NewAccountBalanceQuery creates an AccountBalanceQuery query which can be used to construct and execute
// an AccountBalanceQuery.
// It is recommended that you use this for creating new instances of an AccountBalanceQuery
// instead of manually creating an instance of the struct.
func NewAccountBalanceQuery() *AccountBalanceQuery {
	header := services.QueryHeader{}
	newQuery := AccountBalanceQuery{
		query: _NewQuery(false, &header),
	}
	newQuery.e = &newQuery

	return &newQuery
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *AccountBalanceQuery) SetGrpcDeadline(deadline *time.Duration) *AccountBalanceQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetAccountID sets the AccountID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (this *AccountBalanceQuery) SetAccountID(accountID AccountID) *AccountBalanceQuery {
	this.accountID = &accountID
	return this
}

// GetAccountID returns the AccountID for which you wish to query the balance.
func (this *AccountBalanceQuery) GetAccountID() AccountID {
	if this.accountID == nil {
		return AccountID{}
	}

	return *this.accountID
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (this *AccountBalanceQuery) SetContractID(contractID ContractID) *AccountBalanceQuery {
	this.contractID = &contractID
	return this
}

// GetContractID returns the ContractID for which you wish to query the balance.
func (this *AccountBalanceQuery) GetContractID() ContractID {
	if this.contractID == nil {
		return ContractID{}
	}

	return *this.contractID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *AccountBalanceQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	this.timestamp = time.Now()
	this.paymentTransactions = make([]*services.Transaction, 0)

	pb := this.build()
	pb.CryptogetAccountBalance.Header = this.pbHeader
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

	cost := int64(resp.(*services.Response).GetCryptogetAccountBalance().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (this *AccountBalanceQuery) Execute(client *Client) (AccountBalance, error) {
	if client == nil {
		return AccountBalance{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return AccountBalance{}, err
	}

	this.timestamp = time.Now()

	this.paymentTransactions = make([]*services.Transaction, 0)

	pb := this.build()
	pb.CryptogetAccountBalance.Header = this.pbHeader
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
		return AccountBalance{}, err
	}

	return _AccountBalanceFromProtobuf(resp.(*services.Response).GetCryptogetAccountBalance()), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *AccountBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountBalanceQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *AccountBalanceQuery) SetQueryPayment(paymentAmount Hbar) *AccountBalanceQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountBalanceQuery.
func (this *AccountBalanceQuery) SetNodeAccountIDs(accountID []AccountID) *AccountBalanceQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *AccountBalanceQuery) SetMaxRetry(count int) *AccountBalanceQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *AccountBalanceQuery) SetMaxBackoff(max time.Duration) *AccountBalanceQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *AccountBalanceQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *AccountBalanceQuery) SetMinBackoff(min time.Duration) *AccountBalanceQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *AccountBalanceQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *AccountBalanceQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	return fmt.Sprintf("AccountBalanceQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *AccountBalanceQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountBalanceQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *AccountBalanceQuery) SetLogLevel(level LogLevel) *AccountBalanceQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *AccountBalanceQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().CryptoGetBalance,
	}
}

func (this *AccountBalanceQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode),
	}
}

func (this *AccountBalanceQuery) getName() string {
	return "AccountBalanceQuery"
}

func (this *AccountBalanceQuery) build() *services.Query_CryptogetAccountBalance {
	pb := services.CryptoGetAccountBalanceQuery{Header: &services.QueryHeader{}}

	if this.accountID != nil {
		pb.BalanceSource = &services.CryptoGetAccountBalanceQuery_AccountID{
			AccountID: this.accountID._ToProtobuf(),
		}
	}

	if this.contractID != nil {
		pb.BalanceSource = &services.CryptoGetAccountBalanceQuery_ContractID{
			ContractID: this.contractID._ToProtobuf(),
		}
	}

	return &services.Query_CryptogetAccountBalance{
		CryptogetAccountBalance: &pb,
	}
}

func (this *AccountBalanceQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.accountID != nil {
		if err := this.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if this.contractID != nil {
		if err := this.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *AccountBalanceQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode)
}
