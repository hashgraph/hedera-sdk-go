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
func (q *AccountBalanceQuery) SetGrpcDeadline(deadline *time.Duration) *AccountBalanceQuery {
	q.query.SetGrpcDeadline(deadline)
	return q
}

// SetAccountID sets the AccountID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (q *AccountBalanceQuery) SetAccountID(accountID AccountID) *AccountBalanceQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the AccountID for which you wish to query the balance.
func (q *AccountBalanceQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}

	return *q.accountID
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (q *AccountBalanceQuery) SetContractID(contractID ContractID) *AccountBalanceQuery {
	q.contractID = &contractID
	return q
}

// GetContractID returns the ContractID for which you wish to query the balance.
func (q *AccountBalanceQuery) GetContractID() ContractID {
	if q.contractID == nil {
		return ContractID{}
	}

	return *q.contractID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (q *AccountBalanceQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	q.timestamp = time.Now()
	q.paymentTransactions = make([]*services.Transaction, 0)

	pb := q.build()
	pb.CryptogetAccountBalance.Header = q.pbHeader
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

	cost := int64(resp.(*services.Response).GetCryptogetAccountBalance().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (q *AccountBalanceQuery) Execute(client *Client) (AccountBalance, error) {
	if client == nil {
		return AccountBalance{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return AccountBalance{}, err
	}

	q.timestamp = time.Now()

	q.paymentTransactions = make([]*services.Transaction, 0)

	pb := q.build()
	pb.CryptogetAccountBalance.Header = q.pbHeader
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
		return AccountBalance{}, err
	}

	return _AccountBalanceFromProtobuf(resp.(*services.Response).GetCryptogetAccountBalance()), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *AccountBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountBalanceQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *AccountBalanceQuery) SetQueryPayment(paymentAmount Hbar) *AccountBalanceQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountBalanceQuery.
func (q *AccountBalanceQuery) SetNodeAccountIDs(accountID []AccountID) *AccountBalanceQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *AccountBalanceQuery) SetMaxRetry(count int) *AccountBalanceQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountBalanceQuery) SetMaxBackoff(max time.Duration) *AccountBalanceQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *AccountBalanceQuery) SetMinBackoff(min time.Duration) *AccountBalanceQuery {
	q.query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountBalanceQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountBalanceQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *AccountBalanceQuery) SetLogLevel(level LogLevel) *AccountBalanceQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *AccountBalanceQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().CryptoGetBalance,
	}
}

func (q *AccountBalanceQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode),
	}
}

func (q *AccountBalanceQuery) getName() string {
	return "AccountBalanceQuery"
}

func (q *AccountBalanceQuery) build() *services.Query_CryptogetAccountBalance {
	pb := services.CryptoGetAccountBalanceQuery{Header: &services.QueryHeader{}}

	if q.accountID != nil {
		pb.BalanceSource = &services.CryptoGetAccountBalanceQuery_AccountID{
			AccountID: q.accountID._ToProtobuf(),
		}
	}

	if q.contractID != nil {
		pb.BalanceSource = &services.CryptoGetAccountBalanceQuery_ContractID{
			ContractID: q.contractID._ToProtobuf(),
		}
	}

	return &services.Query_CryptogetAccountBalance{
		CryptogetAccountBalance: &pb,
	}
}

func (q *AccountBalanceQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.accountID != nil {
		if err := q.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if q.contractID != nil {
		if err := q.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *AccountBalanceQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode)
}
